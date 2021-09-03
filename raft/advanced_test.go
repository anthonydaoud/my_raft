package raft

import (
	"github.com/raft/hashmachine"
	"math"
	"strconv"
	"testing"
	"time"
)

// vars to ensure sequential running of tests to prevent raftlog global state from
// interefering with testing operations
var seq1 chan bool = make(chan bool, 1)
var seq2 chan bool = make(chan bool, 1)
var seq3 chan bool = make(chan bool, 1)

//Test to ensure that functionality of raft (good inputs and function calls)
func TestBasicFunctionality(t *testing.T) {
	config := DefaultConfig()
	waitTime := 500 * time.Millisecond
	nodes, err := CreateLocalCluster(config)
	if err != nil {
		t.Errorf("Error creating nodes: %v", err)
		return
	}
	timeDelay := randomTimeout(time.Millisecond * 500)
	<-timeDelay
	leaderFound := false
	//ensure only 1 leader in cluster with log entry 0 inserted
	for !leaderFound {
		for _, node := range nodes {
			if node.State == LEADER_STATE {
				if leaderFound {
					t.Errorf("Multiple leaders found")
					return
				} else {
					leaderFound = true
				}
			}
			entry := node.getLogEntry(0)
			if entry == nil || entry.Type != CommandType_INIT {
				t.Errorf("Bad first entry")
				return
			}
		}
	}
	timeDelay = randomTimeout(time.Millisecond * 150)
	<-timeDelay
	client, err := Connect(nodes[0].GetRemoteSelf().Addr)
	if err != nil {
		t.Errorf("Client failed to connect")
		return
	}
	timeDelay = randomTimeout(time.Millisecond * 200)
	<-timeDelay
	for _, node := range nodes {
		entry := node.getLogEntry(1)
		if entry == nil || entry.Type != CommandType_CLIENT_REGISTRATION {
			t.Errorf("Client Still Not Registered")
			return
		}
	}
	err = client.SendRequest(hashmachine.HASH_CHAIN_INIT, []byte(strconv.Itoa(123)))
	if err != nil {
		t.Errorf("Client request failed")
	}
	timeDelay = randomTimeout(time.Millisecond * 2000)
	<-timeDelay
	for _, node := range nodes {
		entry := node.getLogEntry(2)
		if entry == nil || entry.Type != CommandType_STATE_MACHINE_COMMAND {
			t.Errorf("Hash Init Command Failed")
			return
		}
	}
	addRequests := 10
	for i := 0; i < addRequests; i++ {
		err = client.SendRequest(hashmachine.HASH_CHAIN_ADD, []byte(strconv.Itoa(i)))
		//wait briefly after requests
		timeDelay = randomTimeout(time.Millisecond * 1000)
		<-timeDelay
		if err != nil {
			t.Errorf("Hash Add Command Failed %v", err)
			return
			//i = int(math.Max(float64(0), float64(i-1)))
		}
	}
	timeDelay = randomTimeout(time.Millisecond * 5000)
	<-timeDelay
	for _, node := range nodes {
		for i := 3; i < 3+addRequests; i++ {
			entry := node.getLogEntry(uint64(i))
			if entry == nil || entry.Type != CommandType_STATE_MACHINE_COMMAND {
				t.Errorf("Hash Add Command Failed")
				return
			}
		}
	}
	for _, node := range nodes {
		node.IsShutdown = true
		timeDelay = randomTimeout(5 * waitTime)
		<-timeDelay
	}
	seq1 <- true
}

//A client who peforms the wrong move whenever possible
func TestBadClient(t *testing.T) {
	<-seq1
	config := DefaultConfig()
	mid := int(math.Floor(float64(config.ClusterSize) / 2))
	nodeX, err := CreateNode(0, nil, config)
	waitTime := 500 * time.Millisecond
	if err != nil {
		t.Errorf("Error in making node")
		return
	}
	_, err = Connect(nodeX.GetRemoteSelf().Addr)
	if err == nil {
		t.Errorf("Should have in error, client connect before cluster made")
		return
	}
	nodes, err := CreateLocalCluster(config)
	if err != nil {
		Error.Printf("Error creating nodes: %v", err)
		return
	}
	timeDelay := randomTimeout(waitTime)
	<-timeDelay
	client, err := Connect(nodes[0].GetRemoteSelf().Addr)
	for err != nil {
		client, err = Connect(nodes[0].GetRemoteSelf().Addr)
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	count := 0
	for _, node := range nodes {
		entry := node.getLogEntry(1)
		if entry != nil && entry.Type == CommandType_CLIENT_REGISTRATION {
			count += 1
		}
	}
	if count <= mid {
		t.Errorf("Client Still Not Registered")
		return
	}
	//verifying visually that state is proper
	nodes[0].Out(nodes[0].String())
	nodes[0].Out(nodes[0].FormatState())
	nodes[0].Out(nodes[0].FormatLogCache())
	//hash before init called
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	err = client.SendRequest(hashmachine.HASH_CHAIN_ADD, []byte(strconv.Itoa(1)))
	if err == nil {
		t.Errorf("Client request should have failed %v", err)
		return
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	err = client.SendRequest(hashmachine.HASH_CHAIN_INIT, []byte(strconv.Itoa(1)))
	if err != nil {
		t.Errorf("Client request failed")
		return
	}
	//init after hash called
	timeDelay = randomTimeout(waitTime)
	initRequests := 5
	<-timeDelay
	for i := 0; i < initRequests; i++ {
		err = client.SendRequest(hashmachine.HASH_CHAIN_INIT, []byte(strconv.Itoa(i)))
		//wait briefly after requests
		timeDelay = randomTimeout(waitTime)
		<-timeDelay
		if err == nil {
			t.Errorf("Client request should have failed")
			return
		}
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	for i := 4; i < 4+initRequests; i++ {
		count = 0
		for _, node := range nodes {
			entry := node.getLogEntry(uint64(i))
			if entry != nil && entry.Type == CommandType_STATE_MACHINE_COMMAND {
				count += 1
			}
		}
		if count <= mid {
			t.Errorf("Commands should have been stored even if bad commands")
			return
		}
	}
	//send request to follower node
	if nodes[0].State != LEADER_STATE {
		client.Leader = nodes[0].GetRemoteSelf()
	} else {
		client.Leader = nodes[1].GetRemoteSelf()
	}
	err = client.SendRequest(hashmachine.HASH_CHAIN_ADD, []byte(strconv.Itoa(50)))
	if err != nil {
		t.Errorf("Should have successfully been performed")
	}
	for _, node := range nodes {
		node.IsShutdown = true
		//node.GracefulExit()
		timeDelay = randomTimeout(5 * waitTime)
		<-timeDelay
	}
	seq2 <- true
}

//Test to ensure that after a partition, rollback occurs and consensus achieved
func TestPartitionRecovery(t *testing.T) {
	<-seq2
	config := DefaultConfig()
	config.ClusterSize = 5
	waitTime := 500 * time.Millisecond
	nodes, err := CreateLocalCluster(config)
	if err != nil {
		Error.Printf("Error creating nodes: %v", err)
		return
	}
	timeDelay := randomTimeout(waitTime)
	<-timeDelay
	client, err := Connect(nodes[0].GetRemoteSelf().Addr)
	for err != nil {
		client, err = Connect(nodes[0].GetRemoteSelf().Addr)
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	err = client.SendRequest(hashmachine.HASH_CHAIN_INIT, []byte(strconv.Itoa(123)))
	if err != nil {
		t.Errorf("Client request failed")
	}
	addRequests := 10
	for i := 0; i < addRequests; i++ {
		err = client.SendRequest(hashmachine.HASH_CHAIN_ADD, []byte(strconv.Itoa(i)))
		//wait briefly after requests
		timeDelay = randomTimeout(waitTime)
		<-timeDelay
		if err != nil {
			t.Errorf("Hash Add Command Failed %v", err)
			return
		}
	}
	timeDelay = randomTimeout(5 * waitTime)
	<-timeDelay
	//origLeader will be partitioned with 1 other node and shouldnt commit past index 12
	origLeaderIdx := -1
	for idx, node := range nodes {
		if node.State == LEADER_STATE {
			origLeaderIdx = idx
			break
		}
	}
	//put origLeader at index 0
	tmp := nodes[0]
	nodes[0] = nodes[origLeaderIdx]
	nodes[origLeaderIdx] = tmp
	//now will separate nodes 0 (leader) and 1 from 2,3,4
	for i := 2; i < 5; i++ {
		node := nodes[i]
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[0].GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[1].GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*nodes[0].GetRemoteSelf(), *node.GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*nodes[1].GetRemoteSelf(), *node.GetRemoteSelf(), false)
	}
	for i := 0; i < 2; i++ {
		node := nodes[i]
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[2].GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[3].GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[4].GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*nodes[2].GetRemoteSelf(), *node.GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*nodes[3].GetRemoteSelf(), *node.GetRemoteSelf(), false)
		node.NetworkPolicy.RegisterPolicy(*nodes[4].GetRemoteSelf(), *node.GetRemoteSelf(), false)
	}
	//while no new leader, continue waiting
	newLeaderIdx := -1
	for newLeaderIdx == -1 {
		timeDelay = randomTimeout(5 * waitTime)
		<-timeDelay
		for i := 2; i < 5; i++ {
			if nodes[i].State == LEADER_STATE {
				newLeaderIdx = i
				break
			}
		}
	}
	//put origLeader at index 2
	tmp = nodes[2]
	nodes[2] = nodes[newLeaderIdx]
	nodes[newLeaderIdx] = tmp
	//have new client attempt to connect to origCluster (should fail and be appended to origCluster)
	_, err = Connect(nodes[1].GetRemoteSelf().Addr)
	if err == nil {
		t.Errorf("Should Have Failed to connect")
		return
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	//have new client attempt to connect to newCluster (should work)
	newClient, err := Connect(nodes[4].GetRemoteSelf().Addr)
	if err != nil {
		t.Errorf("Should Have connected")
	}
	timeDelay = randomTimeout(waitTime)
	<-timeDelay
	//perform new add requests
	addRequests = 10
	for i := 0; i < addRequests; i++ {
		err = newClient.SendRequest(hashmachine.HASH_CHAIN_ADD, []byte(strconv.Itoa(i)))
		//wait briefly after requests
		timeDelay = randomTimeout(waitTime)
		<-timeDelay
		if err != nil {
			t.Errorf("Hash Add Command Failed %v", err)
			return
		}
	}
	timeDelay = randomTimeout(10 * waitTime)
	<-timeDelay
	//verify log entries are in up to index at 23 (1 connect, 10 requests more than before)
	for i := 2; i < 5; i++ {
		entry := nodes[i].getLogEntry(nodes[i].getLastLogIndex())
		if entry.Index != 23 || nodes[i].commitIndex != 23 {
			t.Errorf("Partitioned nodes failed to handle requests")
			return
		}
	}
	//also verify commit index at 12 for original Cluster
	//note: last entry is 15 cause of client reg request
	for i := 0; i < 2; i++ {
		entry := nodes[i].getLogEntry(nodes[i].getLastLogIndex())
		if nodes[i].commitIndex != 12 || entry.Index != 15 {
			t.Errorf("Original nodes have bad last log entry")
			return
		}
	}
	//restore partition and verify all nodes in cluster have same setup (using network policy)
	for i := 2; i < 5; i++ {
		node := nodes[i]
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[0].GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[1].GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*nodes[0].GetRemoteSelf(), *node.GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*nodes[1].GetRemoteSelf(), *node.GetRemoteSelf(), true)
	}
	for i := 0; i < 2; i++ {
		node := nodes[i]
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[2].GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[3].GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*node.GetRemoteSelf(), *nodes[4].GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*nodes[2].GetRemoteSelf(), *node.GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*nodes[3].GetRemoteSelf(), *node.GetRemoteSelf(), true)
		node.NetworkPolicy.RegisterPolicy(*nodes[4].GetRemoteSelf(), *node.GetRemoteSelf(), true)
	}
	timeDelay = randomTimeout(10 * waitTime)
	<-timeDelay
	for _, node := range nodes {
		entry := node.getLogEntry(node.getLastLogIndex())
		if entry.Index != 23 || node.commitIndex != 23 {
			t.Errorf("Partitioned failed to be resolved")
			return
		}
	}
	for _, node := range nodes {
		node.IsShutdown = true
		timeDelay = randomTimeout(3 * waitTime)
		<-timeDelay
	}
	seq3 <- true
}

//Test to ensure that our client can simply connect to a raft... A first test.
func TestBasicClientConnect(t *testing.T) {
	<-seq3
	ports := make([]int, 0)
	ports = append(ports, 7005)
	ports = append(ports, 7006)
	ports = append(ports, 7007)
	waitTime := 500 * time.Millisecond
	config := DefaultConfig()
	nodes, err := CreateDefinedLocalCluster(config, ports)
	if err != nil {
		t.Errorf("Failed to make cluster")
		return
	}
	timeDelay := randomTimeout(5 * waitTime)
	<-timeDelay
	_, err = Connect(nodes[0].GetRemoteSelf().Addr)
	if err != nil {
		t.Errorf("Client failed to connect")
		return
	}
	timeDelay = randomTimeout(5 * waitTime)
	<-timeDelay
	origLeaderIdx := -1
	for idx, node := range nodes {
		if node.State == LEADER_STATE {
			origLeaderIdx = idx
			break
		}
	}
	//put origLeader at index 0
	tmp := nodes[0]
	nodes[0] = nodes[origLeaderIdx]
	nodes[origLeaderIdx] = tmp
}
