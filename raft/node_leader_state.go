package raft

import (
	"math"
	"time"
)

// doLeader implements the logic for a Raft node in the leader state.
func (r *RaftNode) doLeader() stateFunction {
	r.Out("Transitioning to LEADER_STATE")

	r.State = LEADER_STATE

	r.onLeaderStart()
	fallbackChan := make(chan bool)
	sentToMajorityChan := make(chan bool)
	timeout := time.After(r.config.HeartbeatTimeout)
	go r.SendHeartbeats(fallbackChan, sentToMajorityChan)
	for {
		select {
		case <-timeout:
			go r.SendHeartbeats(fallbackChan, sentToMajorityChan)
			timeout = time.After(r.config.HeartbeatTimeout)
		case fallback := <-fallbackChan:

			if fallback {
				r.onLeaderExit()
				return r.doFollower
			}
		case sentToMajority := <-sentToMajorityChan:
			if sentToMajority {
				matchIndex := r.findFurthestValidMatchIndex()
				if matchIndex > r.commitIndex {
					r.commitIndex = matchIndex
					go r.processLogEntries()
				}
			}
		case msg := <-r.requestVote: // type RequestVoteMsg

			goBack := r.handleCompetingRequestVote(msg)
			if goBack {
				r.onLeaderExit()
				return r.doFollower
			}
		case msg := <-r.registerClient: // type RegisterClientMsg
			go r.handleLeaderRegisterClient(msg)
		case msg := <-r.clientRequest: // type ClientRequestMsg
			go r.handleLeaderClientRequest(msg)

		case msg := <-r.appendEntries: // type AppendEntriesMsg
			_, fallback := r.handleAppendEntries(msg)
			if fallback {
				r.onLeaderExit()
				return r.doFollower
			}
		case shutdown := <-r.gracefulExit:
			// clean up
			if shutdown {
				r.onLeaderExit()
				return nil
			}
		default:
		}
		if r.IsShutdown {
			return nil
		}

	}
	return nil
}

// sendHeartbeats is used by the leader to send out heartbeats to each of
// the other nodes. It returns true if the leader should fall back to the
// follower state. (This happens if we discover that we are in an old term.)
//
// If another node isn't up-to-date, then the leader should attempt to
// update them, and, if an index has made it to a quorum of nodes, commit
// up to that index. Once committed to that index, the replicated state
// machine should be given the new log entries via processLogEntry.
func (r *RaftNode) SendHeartbeats(fallbackChan, sentToMajorityChan chan bool) bool {
	otherNodes := r.getOtherNodes() // []RemoteNode
	helperFallbackChan := make(chan bool)
	failedToSend := make(chan bool)
	for _, remoteNode := range otherNodes {
		heartbeat := r.createHeatbeat(&remoteNode)
		if heartbeat == nil {
			continue
		}
		//sendSingleHeartbeat cant have remoteNode passed by ref, will refer to only 1 remote node
		go r.sendSingleHeartbeat(remoteNode, heartbeat, helperFallbackChan, failedToSend)
	}
	successCount := 1 //we are successful for ourself
	mid := int(math.Floor(float64(r.config.ClusterSize) / 2))
	failCount := 0
	for {
		select {
		case goBack := <-helperFallbackChan:
			if goBack {
				if fallbackChan != nil {
					fallbackChan <- true
				}
				return false
			} else {
				successCount++
				if successCount > mid {
					if sentToMajorityChan != nil {
						sentToMajorityChan <- true
					}
					return true
				}
			}
		case <-failedToSend:
			failCount += 1
			if failCount > mid {
				return false
			}
		}
	}
	return false

}

// helper to send a single heartbeat to one other raft node
func (r *RaftNode) sendSingleHeartbeat(remoteNode RemoteNode, heartbeat *AppendEntriesRequest,
	helperFallbackChan chan bool, failChan chan bool) {
	reply, err := remoteNode.AppendEntriesRPC(r, heartbeat)
	// RESOLVED: not sure what to do with error, prob nothing as it is in a heartbeat
	if err != nil {
		failChan <- false
		return
	}
	if reply.Success {
		if len(heartbeat.Entries) > 0 {
			r.leaderMutex.Lock()
			r.nextIndex[remoteNode.Id] = uint64(math.Max(float64(heartbeat.Entries[len(heartbeat.Entries)-1].Index+1),
				float64(r.nextIndex[remoteNode.Id])))
			r.matchIndex[remoteNode.Id] = r.nextIndex[remoteNode.Id] - 1
			r.leaderMutex.Unlock()
		}
		helperFallbackChan <- false
	} else {
		// if we are in this case one of two things happened
		// (1) the term of the follower is higher than our term
		// (2) we need to send some older data from our log
		if reply.Term > r.GetCurrentTerm() { // (1)
			r.setCurrentTerm(reply.Term) // downgrade our term
			helperFallbackChan <- true
		} else {
			r.leaderMutex.Lock()
			r.nextIndex[remoteNode.Id] = uint64(math.Max(float64(r.nextIndex[remoteNode.Id]-1), float64(0)))
			r.leaderMutex.Unlock()
			helperFallbackChan <- false
		}
	}
}

// helper to help with leader start
func (r *RaftNode) onLeaderStart() {
	r.processLogEntries() //updates last applied
	r.Leader = r.GetRemoteSelf()
	otherNodes := r.getOtherNodes() // []RemoteNode
	index := r.getLastLogIndex() + 1
	for _, node := range otherNodes {
		r.leaderMutex.Lock()
		r.nextIndex[node.Id] = index
		r.matchIndex[node.Id] = 0
		r.leaderMutex.Unlock()
	}
}

// helper to help with leader exit
func (r *RaftNode) onLeaderExit() {
	r.leaderMutex.Lock()
	r.nextIndex = make(map[string]uint64)
	r.matchIndex = make(map[string]uint64)
	r.leaderMutex.Unlock()
}

// helper to create a heatbeat
func (r *RaftNode) createHeatbeat(remoteNode *RemoteNode) *AppendEntriesRequest {
	r.leaderMutex.Lock()
	remoteNodeNextIndex := r.nextIndex[remoteNode.Id]
	currTerm := r.GetCurrentTerm()
	leader := r.Leader
	var prevLogTerm uint64
	if r.getLogEntry(remoteNodeNextIndex-1) != nil {
		prevLogTerm = r.getLogEntry(remoteNodeNextIndex - 1).TermId
	} else {
		r.leaderMutex.Unlock()
		return nil
	}
	entries := r.getLogEntries(remoteNodeNextIndex)
	commit := r.commitIndex
	heartbeat := AppendEntriesRequest{
		Term:         currTerm,
		Leader:       leader, // should be self
		PrevLogIndex: remoteNodeNextIndex - 1,
		PrevLogTerm:  prevLogTerm, // the term of the log entry at prevLogIndex
		Entries:      entries,     // []*LogEntry
		LeaderCommit: commit,      // leader's commit commitIndex (index of last commited entry)
	}
	r.leaderMutex.Unlock()
	return &heartbeat
}

// helper to deal with a client request when we are in the leader state
func (r *RaftNode) handleLeaderClientRequest(msg ClientRequestMsg) {
	Debug.Printf("In handle Client Request")
	entry := LogEntry{
		Index:   r.getLastLogIndex() + 1,
		TermId:  r.GetCurrentTerm(),
		Type:    CommandType_STATE_MACHINE_COMMAND,
		Command: msg.request.GetStateMachineCmd(), //doesn't matter if type!=STATE_MACHINE_COMMAND
		Data:    msg.request.GetData(),
		CacheId: createCacheId(msg.request.GetClientId(), msg.request.GetSequenceNum()),
	}
	r.appendLogEntry(entry)
	r.requestsMutex.Lock()
	r.requestsByCacheId[entry.CacheId] = msg.reply
	r.requestsMutex.Unlock()
}

// helper to handle the registration of a client when in leader state
func (r *RaftNode) handleLeaderRegisterClient(msg RegisterClientMsg) {
	Debug.Printf("In handle Reg Client")
	//must guarantee commited entry before replying
	entry := LogEntry{
		Index:   r.getLastLogIndex() + 1,
		TermId:  r.GetCurrentTerm(),
		Type:    CommandType_CLIENT_REGISTRATION,
		Command: 0,
		Data:    nil,
		CacheId: "",
	}
	//send to all others
	r.appendLogEntry(entry)

	if r.SendHeartbeats(nil, nil) == true {
		response := RegisterClientReply{
			Status:     ClientStatus_OK,
			ClientId:   entry.Index,
			LeaderHint: r.Leader,
		}
		msg.reply <- response
	} else {
		response := RegisterClientReply{
			Status:     ClientStatus_REQ_FAILED,
			ClientId:   0,
			LeaderHint: r.Leader,
		}
		msg.reply <- response
	}
}

// helper to get a range of entries includes first index
func (r *RaftNode) getLogEntries(start uint64) []*LogEntry {
	//Debug.Printf("line")
	entries := make([]*LogEntry, 0)
	// if start >=, this method will return empty entries
	for i := start; i <= r.getLastLogIndex(); i++ { // last index was incremented after client request
		entries = append(entries, r.getLogEntry(i))
	}
	return entries
}

//finds the furthest N such that a majority of matchIndex[i] >= N append
//log[N].term == CurrentTerm
func (r *RaftNode) findFurthestValidMatchIndex() uint64 {
	if r.config.ClusterSize < 2 {
		return uint64(0)
	}
	tmp := make([]uint64, len(r.matchIndex))
	idx := 0
	r.leaderMutex.Lock()
	for _, value := range r.matchIndex {
		tmp[idx] = value
		idx++
	}
	r.leaderMutex.Unlock()
	bubbleSort(tmp)
	//highest match index which a majority agree with
	mid := len(tmp) / 2
	highestMatchIndex := tmp[mid]
	for highestMatchIndex > 0 {
		if r.getLogEntry(highestMatchIndex).GetTermId() == r.GetCurrentTerm() {
			return highestMatchIndex
		}
		highestMatchIndex -= 1
	}
	return highestMatchIndex
}

// helper to sort, used to determine the match index
func bubbleSort(arrayInput []uint64) {
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(arrayInput)-1; i++ {
			if arrayInput[i+1] < arrayInput[i] {
				tmp := arrayInput[i]
				arrayInput[i] = arrayInput[i+1]
				arrayInput[i+1] = tmp
				swapped = true
			}
		}
	}
}

// processLogEntry applies a single log entry to the finite state machine. It is
// called once a log entry has been replicated to a majority and committed by
// the leader. Once the entry has been applied, the leader responds to the client
// with the result, and also caches the response.
func (r *RaftNode) processLogEntry(entry LogEntry) ClientReply {
	Out.Printf("Node %v Processing log entry: %v \n", r.Id, entry)

	status := ClientStatus_OK
	response := ""
	var err error

	// Apply command on state machine
	if entry.Type == CommandType_STATE_MACHINE_COMMAND {
		response, err = r.stateMachine.ApplyCommand(entry.Command, entry.Data)
		if err != nil {
			status = ClientStatus_REQ_FAILED
			response = err.Error()
		}
	}

	// Construct reply
	reply := ClientReply{
		Status:     status,
		Response:   response,
		LeaderHint: r.GetRemoteSelf(),
	}

	// Add reply to cache
	if entry.CacheId != "" {
		r.CacheClientReply(entry.CacheId, reply)
	}

	// Send reply to client
	r.requestsMutex.Lock()
	replyChan, exists := r.requestsByCacheId[entry.CacheId]
	if exists {
		replyChan <- reply
		delete(r.requestsByCacheId, entry.CacheId)
	}
	r.requestsMutex.Unlock()

	return reply
}
