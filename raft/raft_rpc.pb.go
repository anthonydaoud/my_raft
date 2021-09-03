// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft_rpc.proto

/*
Package raft is a generated protocol buffer package.

It is generated from these files:
	raft_rpc.proto

It has these top-level messages:
	Ok
	RemoteNode
	StartNodeRequest
	LogEntry
	AppendEntriesRequest
	AppendEntriesReply
	RequestVoteRequest
	RequestVoteReply
	RegisterClientRequest
	RegisterClientReply
	ClientRequest
	ClientReply
*/
package raft

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A log entry in Raft can be any of any of these four types.
type CommandType int32

const (
	CommandType_CLIENT_REGISTRATION   CommandType = 0
	CommandType_INIT                  CommandType = 1
	CommandType_NOOP                  CommandType = 2
	CommandType_STATE_MACHINE_COMMAND CommandType = 3
)

var CommandType_name = map[int32]string{
	0: "CLIENT_REGISTRATION",
	1: "INIT",
	2: "NOOP",
	3: "STATE_MACHINE_COMMAND",
}
var CommandType_value = map[string]int32{
	"CLIENT_REGISTRATION":   0,
	"INIT":                  1,
	"NOOP":                  2,
	"STATE_MACHINE_COMMAND": 3,
}

func (x CommandType) String() string {
	return proto.EnumName(CommandType_name, int32(x))
}
func (CommandType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The possible responses to a client request
type ClientStatus int32

const (
	ClientStatus_OK                   ClientStatus = 0
	ClientStatus_NOT_LEADER           ClientStatus = 1
	ClientStatus_ELECTION_IN_PROGRESS ClientStatus = 2
	ClientStatus_CLUSTER_NOT_STARTED  ClientStatus = 3
	ClientStatus_REQ_FAILED           ClientStatus = 4
)

var ClientStatus_name = map[int32]string{
	0: "OK",
	1: "NOT_LEADER",
	2: "ELECTION_IN_PROGRESS",
	3: "CLUSTER_NOT_STARTED",
	4: "REQ_FAILED",
}
var ClientStatus_value = map[string]int32{
	"OK":                   0,
	"NOT_LEADER":           1,
	"ELECTION_IN_PROGRESS": 2,
	"CLUSTER_NOT_STARTED":  3,
	"REQ_FAILED":           4,
}

func (x ClientStatus) String() string {
	return proto.EnumName(ClientStatus_name, int32(x))
}
func (ClientStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// Used to represent result of an RPC call that can either be successful or
// unsuccessful.
type Ok struct {
	Ok     bool   `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
	Reason string `protobuf:"bytes,2,opt,name=reason" json:"reason,omitempty"`
}

func (m *Ok) Reset()                    { *m = Ok{} }
func (m *Ok) String() string            { return proto.CompactTextString(m) }
func (*Ok) ProtoMessage()               {}
func (*Ok) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Ok) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *Ok) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

// Represents a node in the Raft cluster
type RemoteNode struct {
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	Id   string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *RemoteNode) Reset()                    { *m = RemoteNode{} }
func (m *RemoteNode) String() string            { return proto.CompactTextString(m) }
func (*RemoteNode) ProtoMessage()               {}
func (*RemoteNode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RemoteNode) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *RemoteNode) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type StartNodeRequest struct {
	// The node sending the request to start another node in the cluster
	FromNode *RemoteNode `protobuf:"bytes,1,opt,name=fromNode" json:"fromNode,omitempty"`
	// The list of nodes in the cluster that the new node should start up with
	NodeList []*RemoteNode `protobuf:"bytes,2,rep,name=nodeList" json:"nodeList,omitempty"`
}

func (m *StartNodeRequest) Reset()                    { *m = StartNodeRequest{} }
func (m *StartNodeRequest) String() string            { return proto.CompactTextString(m) }
func (*StartNodeRequest) ProtoMessage()               {}
func (*StartNodeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StartNodeRequest) GetFromNode() *RemoteNode {
	if m != nil {
		return m.FromNode
	}
	return nil
}

func (m *StartNodeRequest) GetNodeList() []*RemoteNode {
	if m != nil {
		return m.NodeList
	}
	return nil
}

type LogEntry struct {
	// Index of log entry (first index = 1)
	Index uint64 `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	// The term that this entry was in when added
	TermId uint64 `protobuf:"varint,2,opt,name=termId" json:"termId,omitempty"`
	// Type of command associated with this entry
	Type CommandType `protobuf:"varint,3,opt,name=type,enum=raft.CommandType" json:"type,omitempty"`
	// Command associated with this log entry in the user's finite-state-machine.
	// Note that we only care about this value when type = STATE_MACHINE_COMMAND
	Command uint64 `protobuf:"varint,4,opt,name=command" json:"command,omitempty"`
	// Data associated with this log entry in the user's finite-state-machine.
	Data []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	// After processing this log entry, what ID to use when caching the
	// response. Use an empty string to not cache at all
	CacheId string `protobuf:"bytes,6,opt,name=cacheId" json:"cacheId,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LogEntry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *LogEntry) GetTermId() uint64 {
	if m != nil {
		return m.TermId
	}
	return 0
}

func (m *LogEntry) GetType() CommandType {
	if m != nil {
		return m.Type
	}
	return CommandType_CLIENT_REGISTRATION
}

func (m *LogEntry) GetCommand() uint64 {
	if m != nil {
		return m.Command
	}
	return 0
}

func (m *LogEntry) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *LogEntry) GetCacheId() string {
	if m != nil {
		return m.CacheId
	}
	return ""
}

type AppendEntriesRequest struct {
	// The leader's term
	Term uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	// Address of the leader sending this request
	Leader *RemoteNode `protobuf:"bytes,2,opt,name=leader" json:"leader,omitempty"`
	// The index of the log entry immediately preceding the new ones
	PrevLogIndex uint64 `protobuf:"varint,3,opt,name=prevLogIndex" json:"prevLogIndex,omitempty"`
	// The term of the log entry at prevLogIndex
	PrevLogTerm uint64 `protobuf:"varint,4,opt,name=prevLogTerm" json:"prevLogTerm,omitempty"`
	// The log entries the follower needs to store. Empty for heartbeat messages.
	Entries []*LogEntry `protobuf:"bytes,5,rep,name=entries" json:"entries,omitempty"`
	// The leader's commitIndex
	LeaderCommit uint64 `protobuf:"varint,6,opt,name=leaderCommit" json:"leaderCommit,omitempty"`
}

func (m *AppendEntriesRequest) Reset()                    { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()               {}
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AppendEntriesRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeader() *RemoteNode {
	if m != nil {
		return m.Leader
	}
	return nil
}

func (m *AppendEntriesRequest) GetPrevLogIndex() uint64 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogTerm() uint64 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AppendEntriesRequest) GetLeaderCommit() uint64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type AppendEntriesReply struct {
	// The current term, for leader to update itself.
	Term uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	// True if follower contained entry matching prevLogIndex and prevLogTerm.
	Success bool `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *AppendEntriesReply) Reset()                    { *m = AppendEntriesReply{} }
func (m *AppendEntriesReply) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesReply) ProtoMessage()               {}
func (*AppendEntriesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *AppendEntriesReply) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type RequestVoteRequest struct {
	// The candidate's current term Id
	Term uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	// The cadidate Id currently requesting a node to vote for it.
	Candidate *RemoteNode `protobuf:"bytes,2,opt,name=candidate" json:"candidate,omitempty"`
	// The index of the candidate's last log entry
	LastLogIndex uint64 `protobuf:"varint,3,opt,name=lastLogIndex" json:"lastLogIndex,omitempty"`
	// The term of the candidate's last log entry
	LastLogTerm uint64 `protobuf:"varint,4,opt,name=lastLogTerm" json:"lastLogTerm,omitempty"`
}

func (m *RequestVoteRequest) Reset()                    { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()               {}
func (*RequestVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RequestVoteRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteRequest) GetCandidate() *RemoteNode {
	if m != nil {
		return m.Candidate
	}
	return nil
}

func (m *RequestVoteRequest) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

type RequestVoteReply struct {
	// The current term, for candidate to update itsel
	Term uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	// True means candidate received vote
	VoteGranted bool `protobuf:"varint,2,opt,name=voteGranted" json:"voteGranted,omitempty"`
}

func (m *RequestVoteReply) Reset()                    { *m = RequestVoteReply{} }
func (m *RequestVoteReply) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteReply) ProtoMessage()               {}
func (*RequestVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *RequestVoteReply) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteReply) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

// Empty message represents that a client needs to send no data over to
// register itself.
type RegisterClientRequest struct {
}

func (m *RegisterClientRequest) Reset()                    { *m = RegisterClientRequest{} }
func (m *RegisterClientRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterClientRequest) ProtoMessage()               {}
func (*RegisterClientRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type RegisterClientReply struct {
	// OK if state machine registered client
	Status ClientStatus `protobuf:"varint,1,opt,name=status,enum=raft.ClientStatus" json:"status,omitempty"`
	// Unique ID for client session
	ClientId uint64 `protobuf:"varint,2,opt,name=clientId" json:"clientId,omitempty"`
	// In cases where the client contacted a non-leader, the node should
	// reply with the correct current leader.
	LeaderHint *RemoteNode `protobuf:"bytes,3,opt,name=leaderHint" json:"leaderHint,omitempty"`
}

func (m *RegisterClientReply) Reset()                    { *m = RegisterClientReply{} }
func (m *RegisterClientReply) String() string            { return proto.CompactTextString(m) }
func (*RegisterClientReply) ProtoMessage()               {}
func (*RegisterClientReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *RegisterClientReply) GetStatus() ClientStatus {
	if m != nil {
		return m.Status
	}
	return ClientStatus_OK
}

func (m *RegisterClientReply) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

func (m *RegisterClientReply) GetLeaderHint() *RemoteNode {
	if m != nil {
		return m.LeaderHint
	}
	return nil
}

type ClientRequest struct {
	// The unique client ID associated with this client session (received
	// via a previous RegisterClient call).
	ClientId uint64 `protobuf:"varint,1,opt,name=clientId" json:"clientId,omitempty"`
	// A sequence number is associated to request to avoid duplicates
	SequenceNum uint64 `protobuf:"varint,2,opt,name=sequenceNum" json:"sequenceNum,omitempty"`
	// Command to be executed on the state machine; it may affect state
	StateMachineCmd uint64 `protobuf:"varint,4,opt,name=stateMachineCmd" json:"stateMachineCmd,omitempty"`
	// Data to accompany the command to the state machine; it may affect state
	Data []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ClientRequest) Reset()                    { *m = ClientRequest{} }
func (m *ClientRequest) String() string            { return proto.CompactTextString(m) }
func (*ClientRequest) ProtoMessage()               {}
func (*ClientRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ClientRequest) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

func (m *ClientRequest) GetSequenceNum() uint64 {
	if m != nil {
		return m.SequenceNum
	}
	return 0
}

func (m *ClientRequest) GetStateMachineCmd() uint64 {
	if m != nil {
		return m.StateMachineCmd
	}
	return 0
}

func (m *ClientRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type ClientReply struct {
	// OK if state machine successfully applied command
	Status ClientStatus `protobuf:"varint,1,opt,name=status,enum=raft.ClientStatus" json:"status,omitempty"`
	// State machine output, if successful
	Response string `protobuf:"bytes,2,opt,name=response" json:"response,omitempty"`
	// In cases where the client contacted a non-leader, the node should
	// reply with the correct current leader.
	LeaderHint *RemoteNode `protobuf:"bytes,3,opt,name=leaderHint" json:"leaderHint,omitempty"`
}

func (m *ClientReply) Reset()                    { *m = ClientReply{} }
func (m *ClientReply) String() string            { return proto.CompactTextString(m) }
func (*ClientReply) ProtoMessage()               {}
func (*ClientReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ClientReply) GetStatus() ClientStatus {
	if m != nil {
		return m.Status
	}
	return ClientStatus_OK
}

func (m *ClientReply) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

func (m *ClientReply) GetLeaderHint() *RemoteNode {
	if m != nil {
		return m.LeaderHint
	}
	return nil
}

func init() {
	proto.RegisterType((*Ok)(nil), "raft.Ok")
	proto.RegisterType((*RemoteNode)(nil), "raft.RemoteNode")
	proto.RegisterType((*StartNodeRequest)(nil), "raft.StartNodeRequest")
	proto.RegisterType((*LogEntry)(nil), "raft.LogEntry")
	proto.RegisterType((*AppendEntriesRequest)(nil), "raft.AppendEntriesRequest")
	proto.RegisterType((*AppendEntriesReply)(nil), "raft.AppendEntriesReply")
	proto.RegisterType((*RequestVoteRequest)(nil), "raft.RequestVoteRequest")
	proto.RegisterType((*RequestVoteReply)(nil), "raft.RequestVoteReply")
	proto.RegisterType((*RegisterClientRequest)(nil), "raft.RegisterClientRequest")
	proto.RegisterType((*RegisterClientReply)(nil), "raft.RegisterClientReply")
	proto.RegisterType((*ClientRequest)(nil), "raft.ClientRequest")
	proto.RegisterType((*ClientReply)(nil), "raft.ClientReply")
	proto.RegisterEnum("raft.CommandType", CommandType_name, CommandType_value)
	proto.RegisterEnum("raft.ClientStatus", ClientStatus_name, ClientStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RaftRPC service

type RaftRPCClient interface {
	// Used when a node in the cluster is first starting up so it can notify a
	// leader what their listening address is.
	JoinCaller(ctx context.Context, in *RemoteNode, opts ...grpc.CallOption) (*Ok, error)
	// Once the first node in the cluster has all of the addresses for all other
	// nodes in the cluster it can then tell them to transition into Follower
	// state and start the Raft protocol.
	StartNodeCaller(ctx context.Context, in *StartNodeRequest, opts ...grpc.CallOption) (*Ok, error)
	// Invoked by leader to replicate log entries; also used as a heartbeat
	// between leaders and followers.
	AppendEntriesCaller(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesReply, error)
	// Invoked by candidate nodes to request votes from other nodes.
	RequestVoteCaller(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteReply, error)
	// Called by a client when it first starts up to register itself with the
	// Raft cluster and get a unique client id.
	RegisterClientCaller(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*RegisterClientReply, error)
	// Called by a client to make a request to a Raft node
	ClientRequestCaller(ctx context.Context, in *ClientRequest, opts ...grpc.CallOption) (*ClientReply, error)
}

type raftRPCClient struct {
	cc *grpc.ClientConn
}

func NewRaftRPCClient(cc *grpc.ClientConn) RaftRPCClient {
	return &raftRPCClient{cc}
}

func (c *raftRPCClient) JoinCaller(ctx context.Context, in *RemoteNode, opts ...grpc.CallOption) (*Ok, error) {
	out := new(Ok)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/JoinCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRPCClient) StartNodeCaller(ctx context.Context, in *StartNodeRequest, opts ...grpc.CallOption) (*Ok, error) {
	out := new(Ok)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/StartNodeCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRPCClient) AppendEntriesCaller(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesReply, error) {
	out := new(AppendEntriesReply)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/AppendEntriesCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRPCClient) RequestVoteCaller(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteReply, error) {
	out := new(RequestVoteReply)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/RequestVoteCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRPCClient) RegisterClientCaller(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*RegisterClientReply, error) {
	out := new(RegisterClientReply)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/RegisterClientCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRPCClient) ClientRequestCaller(ctx context.Context, in *ClientRequest, opts ...grpc.CallOption) (*ClientReply, error) {
	out := new(ClientReply)
	err := grpc.Invoke(ctx, "/raft.RaftRPC/ClientRequestCaller", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RaftRPC service

type RaftRPCServer interface {
	// Used when a node in the cluster is first starting up so it can notify a
	// leader what their listening address is.
	JoinCaller(context.Context, *RemoteNode) (*Ok, error)
	// Once the first node in the cluster has all of the addresses for all other
	// nodes in the cluster it can then tell them to transition into Follower
	// state and start the Raft protocol.
	StartNodeCaller(context.Context, *StartNodeRequest) (*Ok, error)
	// Invoked by leader to replicate log entries; also used as a heartbeat
	// between leaders and followers.
	AppendEntriesCaller(context.Context, *AppendEntriesRequest) (*AppendEntriesReply, error)
	// Invoked by candidate nodes to request votes from other nodes.
	RequestVoteCaller(context.Context, *RequestVoteRequest) (*RequestVoteReply, error)
	// Called by a client when it first starts up to register itself with the
	// Raft cluster and get a unique client id.
	RegisterClientCaller(context.Context, *RegisterClientRequest) (*RegisterClientReply, error)
	// Called by a client to make a request to a Raft node
	ClientRequestCaller(context.Context, *ClientRequest) (*ClientReply, error)
}

func RegisterRaftRPCServer(s *grpc.Server, srv RaftRPCServer) {
	s.RegisterService(&_RaftRPC_serviceDesc, srv)
}

func _RaftRPC_JoinCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).JoinCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/JoinCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).JoinCaller(ctx, req.(*RemoteNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRPC_StartNodeCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).StartNodeCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/StartNodeCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).StartNodeCaller(ctx, req.(*StartNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRPC_AppendEntriesCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).AppendEntriesCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/AppendEntriesCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).AppendEntriesCaller(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRPC_RequestVoteCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).RequestVoteCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/RequestVoteCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).RequestVoteCaller(ctx, req.(*RequestVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRPC_RegisterClientCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).RegisterClientCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/RegisterClientCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).RegisterClientCaller(ctx, req.(*RegisterClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRPC_ClientRequestCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRPCServer).ClientRequestCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.RaftRPC/ClientRequestCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRPCServer).ClientRequestCaller(ctx, req.(*ClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RaftRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "raft.RaftRPC",
	HandlerType: (*RaftRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinCaller",
			Handler:    _RaftRPC_JoinCaller_Handler,
		},
		{
			MethodName: "StartNodeCaller",
			Handler:    _RaftRPC_StartNodeCaller_Handler,
		},
		{
			MethodName: "AppendEntriesCaller",
			Handler:    _RaftRPC_AppendEntriesCaller_Handler,
		},
		{
			MethodName: "RequestVoteCaller",
			Handler:    _RaftRPC_RequestVoteCaller_Handler,
		},
		{
			MethodName: "RegisterClientCaller",
			Handler:    _RaftRPC_RegisterClientCaller_Handler,
		},
		{
			MethodName: "ClientRequestCaller",
			Handler:    _RaftRPC_ClientRequestCaller_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft_rpc.proto",
}

func init() { proto.RegisterFile("raft_rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 867 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x55, 0x41, 0x6f, 0xe3, 0x44,
	0x14, 0xae, 0x1d, 0x37, 0x4d, 0x5f, 0x4a, 0xd6, 0x9d, 0x74, 0xbb, 0xde, 0x70, 0x89, 0x2c, 0x21,
	0x45, 0xd5, 0xaa, 0x5a, 0x15, 0x71, 0x46, 0xc1, 0x35, 0xad, 0x21, 0xb1, 0xcb, 0xc4, 0xcb, 0xd5,
	0x1a, 0xec, 0x69, 0xd7, 0x6a, 0xe2, 0x31, 0xf6, 0x74, 0x45, 0x7e, 0x01, 0x07, 0x0e, 0x1c, 0xb9,
	0x71, 0xe5, 0xaf, 0xf1, 0x33, 0xd0, 0x8c, 0xc7, 0x59, 0x3b, 0xb8, 0x48, 0x68, 0x6f, 0xf3, 0xbe,
	0xf7, 0xde, 0xf7, 0xbe, 0xf7, 0xc5, 0x33, 0x81, 0x51, 0x41, 0xee, 0x79, 0x54, 0xe4, 0xf1, 0x65,
	0x5e, 0x30, 0xce, 0x90, 0x21, 0x62, 0xfb, 0x0d, 0xe8, 0xc1, 0x23, 0x1a, 0x81, 0xce, 0x1e, 0x2d,
	0x6d, 0xaa, 0xcd, 0x06, 0x58, 0x67, 0x8f, 0xe8, 0x1c, 0xfa, 0x05, 0x25, 0x25, 0xcb, 0x2c, 0x7d,
	0xaa, 0xcd, 0x8e, 0xb1, 0x8a, 0xec, 0xb7, 0x00, 0x98, 0x6e, 0x18, 0xa7, 0x3e, 0x4b, 0x28, 0x42,
	0x60, 0x90, 0x24, 0x29, 0x64, 0xdf, 0x31, 0x96, 0x67, 0xc1, 0x94, 0x26, 0xaa, 0x4b, 0x4f, 0x13,
	0x3b, 0x03, 0x73, 0xc5, 0x49, 0xc1, 0x45, 0x03, 0xa6, 0x3f, 0x3f, 0xd1, 0x92, 0xa3, 0x37, 0x30,
	0xb8, 0x2f, 0xd8, 0x46, 0x40, 0xb2, 0x77, 0x78, 0x65, 0x5e, 0x0a, 0x31, 0x97, 0x1f, 0xb9, 0xf1,
	0xae, 0x42, 0x54, 0x67, 0x2c, 0xa1, 0x8b, 0xb4, 0xe4, 0x96, 0x3e, 0xed, 0x75, 0x57, 0xd7, 0x15,
	0xf6, 0x5f, 0x1a, 0x0c, 0x16, 0xec, 0xc1, 0xcd, 0x78, 0xb1, 0x45, 0x67, 0x70, 0x98, 0x66, 0x09,
	0xfd, 0x45, 0x4e, 0x31, 0x70, 0x15, 0x88, 0xe5, 0x38, 0x2d, 0x36, 0x5e, 0x25, 0xd3, 0xc0, 0x2a,
	0x42, 0x5f, 0x80, 0xc1, 0xb7, 0x39, 0xb5, 0x7a, 0x53, 0x6d, 0x36, 0xba, 0x3a, 0xad, 0x86, 0x38,
	0x6c, 0xb3, 0x21, 0x59, 0x12, 0x6e, 0x73, 0x8a, 0x65, 0x1a, 0x59, 0x70, 0x14, 0x57, 0xa0, 0x65,
	0xc8, 0xfe, 0x3a, 0x14, 0x7e, 0x24, 0x84, 0x13, 0xeb, 0x70, 0xaa, 0xcd, 0x4e, 0xb0, 0x3c, 0xcb,
	0x6a, 0x12, 0xbf, 0xa7, 0x5e, 0x62, 0xf5, 0xa5, 0x29, 0x75, 0x68, 0xff, 0xad, 0xc1, 0xd9, 0x3c,
	0xcf, 0x69, 0x96, 0x08, 0xb1, 0x29, 0x2d, 0x6b, 0x7b, 0x10, 0x18, 0x42, 0x91, 0x12, 0x2d, 0xcf,
	0x68, 0x06, 0xfd, 0x35, 0x25, 0x09, 0x2d, 0xa4, 0xe6, 0x2e, 0x0b, 0x54, 0x1e, 0xd9, 0x70, 0x92,
	0x17, 0xf4, 0xc3, 0x82, 0x3d, 0x78, 0x72, 0xf5, 0x9e, 0x64, 0x69, 0x61, 0x68, 0x0a, 0x43, 0x15,
	0x87, 0x62, 0x50, 0xb5, 0x46, 0x13, 0x42, 0x33, 0x38, 0xa2, 0x95, 0x2a, 0xeb, 0x50, 0x7a, 0x3e,
	0xaa, 0x06, 0xd6, 0xd6, 0xe2, 0x3a, 0x2d, 0xe6, 0x55, 0x93, 0x85, 0x53, 0x29, 0x97, 0x5b, 0x1a,
	0xb8, 0x85, 0xd9, 0xdf, 0x00, 0xda, 0xdb, 0x34, 0x5f, 0x6f, 0x3b, 0xf7, 0xb4, 0xe0, 0xa8, 0x7c,
	0x8a, 0x63, 0x5a, 0x96, 0x72, 0xd1, 0x01, 0xae, 0x43, 0xfb, 0x4f, 0x0d, 0x90, 0x72, 0xe8, 0x47,
	0xc6, 0xe9, 0x7f, 0x99, 0x75, 0x09, 0xc7, 0x31, 0xc9, 0x92, 0x34, 0x21, 0x9c, 0x3e, 0xeb, 0xd7,
	0xc7, 0x12, 0xb9, 0x02, 0x29, 0xf9, 0xbe, 0x65, 0x4d, 0x4c, 0x58, 0xa6, 0xe2, 0xa6, 0x65, 0x0d,
	0xc8, 0xbe, 0x05, 0xb3, 0xa5, 0xef, 0xb9, 0x15, 0xa7, 0x30, 0xfc, 0xc0, 0x38, 0xbd, 0x29, 0x48,
	0xc6, 0x69, 0xa2, 0xd6, 0x6c, 0x42, 0xf6, 0x2b, 0x78, 0x89, 0xe9, 0x43, 0x5a, 0x72, 0x5a, 0x38,
	0xeb, 0x94, 0x66, 0x5c, 0xf1, 0xda, 0xbf, 0x6b, 0x30, 0xde, 0xcf, 0x88, 0x31, 0x17, 0xd0, 0x2f,
	0x39, 0xe1, 0x4f, 0xa5, 0x1c, 0x34, 0xba, 0x42, 0xea, 0xdb, 0x95, 0x25, 0x2b, 0x99, 0xc1, 0xaa,
	0x02, 0x4d, 0x60, 0x10, 0x4b, 0x7c, 0xf7, 0xfd, 0xef, 0x62, 0xf4, 0x16, 0xa0, 0xfa, 0xdd, 0x6e,
	0xd3, 0x8c, 0x4b, 0x1b, 0xba, 0x9c, 0x6b, 0xd4, 0xd8, 0xbf, 0x69, 0xf0, 0x59, 0x4b, 0x63, 0x8b,
	0x5f, 0xdb, 0xe3, 0x9f, 0xc2, 0xb0, 0x14, 0x65, 0x59, 0x4c, 0xfd, 0xa7, 0x8d, 0x1a, 0xdf, 0x84,
	0xd0, 0x0c, 0x5e, 0x08, 0x9d, 0x74, 0x49, 0xe2, 0xf7, 0x69, 0x46, 0x9d, 0x4d, 0x7d, 0xc9, 0xf6,
	0xe1, 0xae, 0xcb, 0x66, 0xff, 0xaa, 0xc1, 0xf0, 0x13, 0x7c, 0x29, 0x68, 0x99, 0xb3, 0xac, 0xa4,
	0xea, 0xf9, 0xda, 0xc5, 0xff, 0xdf, 0x97, 0x8b, 0x77, 0x30, 0x6c, 0xbc, 0x1c, 0xe8, 0x15, 0x8c,
	0x9d, 0x85, 0xe7, 0xfa, 0x61, 0x84, 0xdd, 0x1b, 0x6f, 0x15, 0xe2, 0x79, 0xe8, 0x05, 0xbe, 0x79,
	0x80, 0x06, 0x60, 0x78, 0xbe, 0x17, 0x9a, 0x9a, 0x38, 0xf9, 0x41, 0x70, 0x67, 0xea, 0xe8, 0x35,
	0xbc, 0x5c, 0x85, 0xf3, 0xd0, 0x8d, 0x96, 0x73, 0xe7, 0xd6, 0xf3, 0xdd, 0xc8, 0x09, 0x96, 0xcb,
	0xb9, 0x7f, 0x6d, 0xf6, 0x2e, 0x52, 0x38, 0x69, 0x8a, 0x47, 0x7d, 0xd0, 0x83, 0xef, 0xcd, 0x03,
	0x34, 0x02, 0xf0, 0x83, 0x30, 0x5a, 0xb8, 0xf3, 0x6b, 0x17, 0x9b, 0x1a, 0xb2, 0xe0, 0xcc, 0x5d,
	0xb8, 0x8e, 0x18, 0x12, 0x79, 0x7e, 0x74, 0x87, 0x83, 0x1b, 0xec, 0xae, 0x56, 0xa6, 0x5e, 0x29,
	0x79, 0xb7, 0x0a, 0x5d, 0x1c, 0x89, 0x8e, 0x55, 0x38, 0xc7, 0xa1, 0x7b, 0x6d, 0xf6, 0x04, 0x05,
	0x76, 0x7f, 0x88, 0xbe, 0x9d, 0x7b, 0x0b, 0xf7, 0xda, 0x34, 0xae, 0xfe, 0xe8, 0xc1, 0x11, 0x26,
	0xf7, 0x1c, 0xdf, 0x39, 0xe8, 0x02, 0xe0, 0x3b, 0x96, 0x66, 0x0e, 0x59, 0xaf, 0x69, 0x81, 0xfe,
	0xb5, 0xf9, 0x64, 0x50, 0x21, 0xc1, 0xa3, 0x7d, 0x80, 0xbe, 0x82, 0x17, 0xbb, 0x07, 0x5f, 0x35,
	0x9c, 0x57, 0xe9, 0xfd, 0xff, 0x81, 0x56, 0xdb, 0x12, 0xc6, 0xad, 0x27, 0x42, 0xb5, 0x4e, 0xaa,
	0x92, 0xae, 0x77, 0x72, 0x62, 0x75, 0xe6, 0xf2, 0xf5, 0xd6, 0x3e, 0x40, 0x37, 0x70, 0xda, 0xb8,
	0x8c, 0x8a, 0xcc, 0xaa, 0x85, 0xef, 0xbf, 0x22, 0x93, 0xf3, 0x8e, 0x4c, 0x45, 0x74, 0x07, 0x67,
	0xed, 0x1b, 0xa7, 0xb8, 0x3e, 0xaf, 0x3b, 0x3a, 0xee, 0xe9, 0xe4, 0x75, 0x77, 0xb2, 0x62, 0xfc,
	0x1a, 0xc6, 0xad, 0x6a, 0x45, 0x38, 0x6e, 0x7e, 0x9b, 0x35, 0xd1, 0x69, 0x1b, 0x94, 0x04, 0x3f,
	0xf5, 0xe5, 0xff, 0xf7, 0x97, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x43, 0x38, 0xc5, 0x93, 0xd1,
	0x07, 0x00, 0x00,
}
