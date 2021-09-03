# Raft

There are no known bugs and no extra features.

Our test cases cover the following:

### Bad Client Behavior:

	Hashing before init
	Init after successful hash
	Connecting to leader that cannot commit becase of partition

### Normal Client Behavior:
	Verify good init
	Verify hashing is correctly processed
### Tests For Followers:
	Verifying that client register is handled correnctly in all cases
	Verifying that client request is handled properly
	Shutting down nodes
	Verifying that logs are up to date after appends (correctness of logs)
	Testing that incorrect follower logs are rolled back (after healed partition)
### Tests for Leader and Candidate:
	Verifying that a correct Candidate promoted to leader
	Make sure leader handles appending entries to followers in normal cases
	Verify that after a partition there can be 2 leaders
	Verify that after a partition is healed all nodes reach a consensus

Tests should run in 3-5 mins.

Code coverage is between 65-70%

Please remove all previous raftlogs before running the following:
<go test -cover>
