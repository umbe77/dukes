// Copyright (c) 2023 Robeto Ughi
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/umbe77/dukes/cache"
)

const (
	raftConnectionPoolCount = 11
	retainSnapshotCount     = 2
)

type Cluster struct {
	Storage *cache.Cache
	raft    *raft.Raft

	RaftDir  string
	RaftBind string
	inmem    bool
}

func NewStore(ch *cache.Cache, raftDir, raftBind string) *Cluster {
	return &Cluster{
		Storage:  ch,
		RaftDir:  raftDir,
		RaftBind: raftBind,
	}
}

func (s *Cluster) SetInmem() {
	s.inmem = true
}

func (s *Cluster) Start(nodeId string, isLeader bool) error {
	//raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	//setup raft communication
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, raftConnectionPoolCount, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	snapshot, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return err
	}

	var (
		logStore    raft.LogStore
		stableStore raft.StableStore
	)

	if s.inmem {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		boltDB, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
		if err != nil {
			return fmt.Errorf("new bolt store: %s", err)
		}
		logStore = boltDB
		stableStore = boltDB
	}

	//Initialize raft system
	ra, err := raft.NewRaft(config, s.Storage, logStore, stableStore, snapshot, transport)
	if err != nil {
		return err
	}
	s.raft = ra

	if isLeader {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}
	return nil
}
