package raft

import (
	"math"
	"strings"
)

// doFollower implements the logic for a Raft node in the follower state.
func (r *RaftNode) doFollower() stateFunction {
	r.Out("Transitioning to FOLLOWER_STATE")
	r.State = FOLLOWER_STATEs

	timeout := randomTimeout(r.config.ElectionTimeout)
	for {
		select {
		case <-timeout:
			return r.doCandidate
		case msg := <-r.appendEntries: // type AppendEntriesMsg
			resetTimeout, fallback := r.handleAppendEntries(msg)
			if resetTimeout || fallback {
				timeout = randomTimeout(r.config.ElectionTimeout)
			}
		case msg := <-r.requestVote: // type RequestVoteMsg
			if r.handleFollowerRequestVote(msg) {
				timeout = randomTimeout(r.config.ElectionTimeout)
			}
		case msg := <-r.registerClient: // type RegisterClientMsg
			r.handleFollowerRegisterClient(msg)
		case msg := <-r.clientRequest: // type ClientRequestMsg
			r.handleFollowerClientRequest(msg)
		case shutdown := <-r.gracefulExit:
			if shutdown {
				return nil
			}
		}
		if r.IsShutdown {
			return nil
		}
	}
	return nil
}

// helper to handle a client request when in follower mode
func (r *RaftNode) handleFollowerClientRequest(msg ClientRequestMsg) {
	response := ClientReply{}
	if r.Leader == nil {
		response = ClientReply{
			Status:     ClientStatus_REQ_FAILED,
			Response:   "NOT LEADER, UNABLE TO PERFORM REQUEST", //client will ignore
			LeaderHint: nil,
		}
	} else {
		response = ClientReply{
			Status:     ClientStatus_NOT_LEADER,
			Response:   "NOT LEADER, UNABLE TO PERFORM REQUEST", //client will ignore
			LeaderHint: r.Leader,
		}
	}
	msg.reply <- response
}

// helper to handle client registration msg when in follower
func (r *RaftNode) handleFollowerRegisterClient(msg RegisterClientMsg) {
	response := RegisterClientReply{}
	if r.Leader == nil {
		response = RegisterClientReply{
			Status:     ClientStatus_REQ_FAILED,
			ClientId:   0, //client will ignore
			LeaderHint: nil,
		}
	} else {
		response = RegisterClientReply{
			Status:     ClientStatus_NOT_LEADER,
			ClientId:   0, //client will ignore
			LeaderHint: r.Leader,
		}
	}
	msg.reply <- response
}

// helper to handle a vote request when we are in follower mode
func (r *RaftNode) handleFollowerRequestVote(msg RequestVoteMsg) bool {
	if msg.request.Term < r.GetCurrentTerm() {
		Debug.Printf("Denied vote, reason: small term")
		r.voteDenied(msg)
		return false
	}
	votedFor := r.GetVotedFor()
	if strings.Compare(msg.request.Candidate.Id, votedFor) == 0 {
		//candidate re-requesting vote
		r.voteGranted(msg)
		return true
	} else if strings.Compare(votedFor, "") != 0 {
		//we've already voted for someone else
		Debug.Printf("Denied vote, reason: already voted for other")

		r.voteDenied(msg)
		return false
	} else {
		//we have yet to vote
		lastLogIndex := r.getLastLogIndex()
		lastLogEntry := r.getLogEntry(lastLogIndex)
		if msg.request.LastLogTerm > lastLogEntry.TermId {
			//will vote for node
			r.voteGranted(msg)
			return true
		} else if msg.request.LastLogTerm == lastLogEntry.TermId &&
			msg.request.LastLogIndex >= lastLogIndex {
			//will vote for node
			r.voteGranted(msg)
			return true
		} else {
			Debug.Printf("Denied vote, reason: log is older than mine")
			r.voteDenied(msg)
			return false
		}
	}
}

//helper function to grant a vote to node in request vote message
func (r *RaftNode) voteGranted(msg RequestVoteMsg) {
	r.setCurrentTerm(msg.request.Term)
	r.setVotedFor(msg.request.Candidate.Id)
	r.Leader = msg.request.Candidate
	response := RequestVoteReply{
		Term:        r.GetCurrentTerm(),
		VoteGranted: true,
	}
	msg.reply <- response
}

//helper function to deny a vote to node in request vote message
func (r *RaftNode) voteDenied(msg RequestVoteMsg) {
	response := RequestVoteReply{
		Term:        r.GetCurrentTerm(),
		VoteGranted: false,
	}
	msg.reply <- response
}

// Applies actions to state machine from lastApplied to commitIndex, updates lastApplied
func (r *RaftNode) processLogEntries() {
	for r.commitIndex >= r.lastApplied {
		toBeApplied := r.getLogEntry(r.lastApplied)
		r.processLogEntry(*toBeApplied)
		r.lastApplied += 1
	}
}

// handleAppendEntries handles an incoming AppendEntriesMsg. It is called by a
// node in a follower, candidate, or leader state. It returns two booleans:
// - resetTimeout is true if the follower node should reset the election timeout
// - fallback is true if the node should become a follower again
func (r *RaftNode) handleAppendEntries(msg AppendEntriesMsg) (resetTimeout, fallback bool) {

	request := msg.request // *AppendEntriesRequest
	reply := msg.reply     // chan AppendEntriesReply

	//valid for all 3 states, verify leader-leader interaction tho
	if request.GetTerm() < r.GetCurrentTerm() {
		response := AppendEntriesReply{
			Term:    r.GetCurrentTerm(),
			Success: false,
		}
		reply <- response
		return false, false
	}
	//can't have 2 leaders in same term as it would imply a node voted twice in same term.
	//Therefore, if you are a leader and receive AppendEntriesMsg and the term is not
	//smaller than yours, term must be greater than yours (unless bug elsewhere). Thus,
	//should fallback and perform same steps as Follower and Candidate State.

	// accept the request as leader, fallback should be true in all return cases
	r.Leader = request.GetLeader()
	r.setCurrentTerm(request.GetTerm())
	logEntry := r.getLogEntry(request.GetPrevLogIndex())
	if logEntry == nil || logEntry.GetTermId() != request.GetPrevLogTerm() {
		response := AppendEntriesReply{
			Term:    r.GetCurrentTerm(),
			Success: false,
		}
		reply <- response
		return true, true
	} else {
		//PrevLogIndex exists with PrevLogTerm matched
		newLogEntry := r.getLogEntry(request.GetPrevLogIndex() + 1)
		if logEntry != nil && newLogEntry.GetTermId() != request.GetTerm() {
			r.truncateLog(request.GetPrevLogIndex() + 1)
		}
		for _, entry := range request.GetEntries() {
			if r.getLastLogIndex() < entry.GetIndex() {
				r.appendLogEntry(*entry)
			} else {
				continue
			}
		}
		if request.GetLeaderCommit() > r.commitIndex {
			r.commitIndex = uint64(math.Min(float64(request.GetLeaderCommit()),
				float64(r.getLastLogIndex())))
		}
		r.processLogEntries()
		response := AppendEntriesReply{
			Term:    r.GetCurrentTerm(),
			Success: true,
		}
		reply <- response
		return true, true
	}
	return
}
