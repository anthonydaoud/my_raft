package raft

import (
	"math"
)

// doCandidate implements the logic for a Raft node in the candidate state.
func (r *RaftNode) doCandidate() stateFunction {
	r.Out("Transitioning to CANDIDATE_STATE")
	r.State = CANDIDATE_STATE

	if r.IsShutdown {
		return nil
	}
	// initial work to ask other nodes for a vote
	r.setCurrentTerm(r.GetCurrentTerm() + 1)
	r.Leader = nil
	r.setVotedFor(r.Id)

	electionResults := make(chan bool)
	fallbackChan := make(chan bool)

	// note: the candidate code increments the term before we enter doCandidate
	go r.requestVotes(electionResults, fallbackChan, r.GetCurrentTerm())
	timeout := randomTimeout(r.config.ElectionTimeout)

	// handling all possible channels
	for {
		select {
		case <-timeout:
			// we did not win the election so we want to increment our term
			// and become a candidate again for a new election cycle
			return r.doCandidate
		case result := <-electionResults:
			if result {
				return r.doLeader
			} else {
				r.setVotedFor("")
				return r.doFollower
			}
		case goBack := <-fallbackChan:
			if goBack {
				r.setVotedFor("")
				return r.doFollower
			}
		case msg := <-r.appendEntries: // type AppendEntriesMsg
			_, goBack := r.handleAppendEntries(msg)
			if goBack {
				r.setVotedFor("")
				return r.doFollower
			}

		case msg := <-r.requestVote: // type RequestVoteMsg
			// at this point you have already voted for yourself
			goBack := r.handleCompetingRequestVote(msg)
			if goBack {
				return r.doFollower
			}
		case msg := <-r.registerClient: // type RegisterClientMsg
			r.handleCandidateRegisterClient(msg)
		case msg := <-r.clientRequest: // type ClientRequestMsg
			r.handleCandidateClientRequest(msg)
		case shutdown := <-r.gracefulExit: // type bool
			// clean up
			if shutdown {
				return nil
			}
		}
		if r.IsShutdown {
			return nil
		}
	}
}

// requestVotes is called to request votes from all other nodes. It takes in a
// channel on which the result of the vote should be sent over: true for a
// successful election, false otherwise.
func (r *RaftNode) requestVotes(electionResults chan bool, fallbackChan chan bool, currTerm uint64) {

	// send RPC to all nodes, retry if the node returns an error
	requestVoteRequest := RequestVoteRequest{
		Term:         currTerm,
		Candidate:    r.GetRemoteSelf(),
		LastLogIndex: r.getLastLogIndex(),
		LastLogTerm:  r.getLogEntry(r.getLastLogIndex()).GetTermId(),
	}
	otherNodes := r.getOtherNodes() // []RemoteNode
	electionResultChan := make(chan bool)
	for _, remoteNode := range otherNodes {
		go r.requestSingleVote(remoteNode, requestVoteRequest, electionResultChan, fallbackChan)
	}

	// collect results
	mid := int(math.Floor(float64(r.config.ClusterSize) / 2))
	successCount := 1
	for {
		select {
		case electionResult := <-electionResultChan:
			if electionResult {
				successCount++
				if successCount > mid {
					electionResults <- true
					return
				}
			}
		}

	}

	return
}

// helper to make a single vote
func (r *RaftNode) requestSingleVote(remoteNode RemoteNode, req RequestVoteRequest,
	electionResultChan chan bool, fallbackChan chan bool) {
	requestVoteReply, err := remoteNode.RequestVoteRPC(r, &req)
	if err != nil || req.Term != r.GetCurrentTerm() {
		return
	}
	if requestVoteReply.GetTerm() > r.GetCurrentTerm() {
		Debug.Printf("Vote reponse tells me to fallback")
		r.setCurrentTerm(requestVoteReply.GetTerm()) // update own term to latest term
		fallbackChan <- true
		return // do not bother sending any more vote requests
	}

	if requestVoteReply.GetVoteGranted() {
		Debug.Printf("Vote reponse good")
		electionResultChan <- true
	} else {
		Debug.Printf("Vote reponse bad")
		electionResultChan <- false
	}
}

// handleCompetingRequestVote handles an incoming vote request when the current
// node is in the candidate or leader state. It returns true if the caller
// should fall back to the follower state, false otherwise.
func (r *RaftNode) handleCompetingRequestVote(msg RequestVoteMsg) (fallback bool) {
	state := r.State
	if state != CANDIDATE_STATE && state != LEADER_STATE {
		return false // default return false, this case should not happen
	}
	request := msg.request // a pointer
	if request.GetTerm() < r.GetCurrentTerm() {
		r.voteDenied(msg)
		return false
	}
	lastLogIndex := r.getLastLogIndex()
	lastLogTerm := r.getLogEntry(lastLogIndex).GetTermId()
	if (lastLogTerm < request.LastLogTerm) ||
		(lastLogTerm == request.LastLogTerm &&
			request.LastLogIndex > lastLogIndex) {
		if r.Leader != nil {
			r.Leader = nil // make sure our leader field is nil
		}
		r.setCurrentTerm(request.GetTerm())
		r.setVotedFor(request.GetCandidate().GetId())
		r.voteGranted(msg) // writes back to chan in the msg
		return true
	}
	if (lastLogTerm == request.LastLogTerm) && (lastLogIndex == request.LastLogIndex) && (request.GetTerm() > r.GetCurrentTerm()) {
		r.voteGranted(msg)
		return true
	}
	// this is necessarily the case where request.GetTerm() == r.GetCurrentTerm()
	r.voteDenied(msg)
	return false
}


// helper to handle a client registration in the case where the raft node is a candidate
func (r *RaftNode) handleCandidateRegisterClient(msg RegisterClientMsg) {
	response := RegisterClientReply{
		Status:     ClientStatus_ELECTION_IN_PROGRESS,
		ClientId:   0, //client will ignore
		LeaderHint: r.GetRemoteSelf(),
	}
	msg.reply <- response
}


// helper to handle a client request in the case where the raft node is a candidate
func (r *RaftNode) handleCandidateClientRequest(msg ClientRequestMsg) {
	response := ClientReply{
		Status:     ClientStatus_ELECTION_IN_PROGRESS,
		Response:   "ELECTION IN PROGRESS, UNABLE TO PERFORM REQUEST", //client will ignore
		LeaderHint: r.GetRemoteSelf(),
	}
	msg.reply <- response
}


// helper to find all nodes that are not us
func (r *RaftNode) getOtherNodes() []RemoteNode {
	remoteNodes := r.GetNodeList()
	otherNodes := make([]RemoteNode, 0)
	for i := 0; i < len(remoteNodes); i++ {
		if remoteNodes[i].Id != r.Id {
			otherNodes = append(otherNodes, remoteNodes[i])
		}
	}
	return otherNodes
}
