## Dukes Distributed Cache Service

### TODO

- [x] Memory Cache
- [ ] Use B+Tree for cache instead of a map?
- [x] TCP Server
- [x] Binary message packets
- [x] Define type for value in message in order to have typed value
- [ ] Raft Consensus
- [x] TCP Client
- [ ] How to review client for access the cluster?


#### FSM Interface
- Apply(*raft.Log) interface{}
- Snapshot() (raft.FSMSnapshot, error)
- Restore(snapshot io.ReadCloser) error
