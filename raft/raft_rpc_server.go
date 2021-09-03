// Purpose: Implements functions that are invoked by clients and other nodes
// over RPC.

package raft

import "golang.org/x/net/context"

// JoinCaller is called through GRPC to execute a join request.
func (local *RaftNode) JoinCaller(ctx context.Context, r *RemoteNode) (*Ok, error) {
	// Check if the network policy prevents incoming requests from the requesting node
	if local.NetworkPolicy.IsDenied(*r, *local.GetRemoteSelf()) {
		return nil, ErrorNetworkPolicyDenied
	}

	// Defer to the local Join implementation, and marshall the results to
	// respond to the GRPC request.
	err := local.Join(r)
	return &Ok{Ok: err == nil}, err
}

// StartNodeCaller is called through GRPC to execute a start node request.
func (local *RaftNode) StartNodeCaller(ctx context.Context, req *StartNodeRequest) (*Ok, error) {
	// RESOLVED: Students should implement this method
	// Make sure to check the provided Raft node's network policy to make sure
	// we are allowed to respond to this request. Respond with
	// ErrorNetworkPolicyDenied if not.
	if local.NetworkPolicy.IsDenied(*(req.FromNode), *local.GetRemoteSelf()) {
		return nil, ErrorNetworkPolicyDenied
	}
	err := local.StartNode(req)

	return &Ok{Ok: err == nil}, err
}

// AppendEntriesCaller is called through GRPC to respond to an append entries request.
func (local *RaftNode) AppendEntriesCaller(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesReply, error) {
	// Make sure to check the provided Raft node's network policy to make sure
	// we are allowed to respond to this request. Respond with
	// ErrorNetworkPolicyDenied if not.
	if local.NetworkPolicy.IsDenied(*(req.Leader), *local.GetRemoteSelf()) {
		return nil, ErrorNetworkPolicyDenied
	}
	reply, err := local.AppendEntries(req)
	return &reply, err
}

// RequestVoteCaller is called through GRPC to respond to a vote request.
func (local *RaftNode) RequestVoteCaller(ctx context.Context, req *RequestVoteRequest) (*RequestVoteReply, error) {
	// Make sure to check the provided Raft node's network policy to make sure
	// we are allowed to respond to this request. Respond with
	// ErrorNetworkPolicyDenied if not.
	if local.NetworkPolicy.IsDenied(*(req.Candidate), *local.GetRemoteSelf()) {
		return nil, ErrorNetworkPolicyDenied
	}
	reply, err := local.RequestVote(req)
	return &reply, err
}

// RegisterClientCaller is called through GRPC to respond to a client
// registration request.
func (local *RaftNode) RegisterClientCaller(ctx context.Context, req *RegisterClientRequest) (*RegisterClientReply, error) {
	reply, err := local.RegisterClient(req)
	return &reply, err
}

// ClientRequestCaller is called through GRPC to respond to a client request.
func (local *RaftNode) ClientRequestCaller(ctx context.Context, req *ClientRequest) (*ClientReply, error) {
	reply, err := local.ClientRequest(req)
	return &reply, err
}
