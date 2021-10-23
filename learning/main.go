package main

import (
	"fmt"
	"io"
	"os"

	"github.com/wu-xian/raft"
)

func main() {
	// fmt.Println("a")
	// raftConfig := raft.DefaultConfig()
	// raftConfig.LocalID = "mylocal"
	// raftConfig.LogOutput = os.Stdout
	// raftStable := raft.NewInmemStore()
	// raftSnaps := raft.NewInmemSnapshotStore()
	// serverAddr := raft.NewInmemAddr()
	// _, raftTrans := raft.NewInmemTransport(serverAddr)
	// fmt.Println("local addr", raftTrans.LocalAddr())
	// fsm := &raft.MockFSM{}
	// raftConfiguration, err := raft.GetConfiguration(raftConfig, fsm, raftStable, raftStable, raftSnaps, raftTrans)
	// if err != nil {
	// 	panic(err)
	// }
	// raftConfiguration.Servers = []raft.Server{raft.Server{ID: "12", Address: "localhost:123"}}
	// err = raft.BootstrapCluster(raftConfig, raftStable, raftStable, raftSnaps, raftTrans, raftConfiguration)
	// if err != nil {
	// 	panic(err)
	// }

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = "mylocal"
	fsm := &learningStateMachine{}
	inmemStore := raft.NewInmemStore()
	snapStore := raft.NewInmemSnapshotStore()
	_, trans := raft.NewInmemTransport(raft.NewInmemAddr())
	ra, err := raft.NewRaft(raftConfig, fsm, inmemStore, inmemStore, snapStore, trans)
	if err != nil {
		panic(err)
	}

	raftConfig.LocalID = "mylocal"
	raftConfig.LogOutput = os.Stdout
	raftStable := raft.NewInmemStore()
	raftSnaps := raft.NewInmemSnapshotStore()
	serverAddr := raft.NewInmemAddr()
	_, raftTrans := raft.NewInmemTransport(serverAddr)
	raftConfiguration, err := raft.GetConfiguration(raftConfig, fsm, raftStable, raftStable, raftSnaps, raftTrans)

	future := ra.BootstrapCluster(raftConfiguration)
	fmt.Println(future)

	ch := make(chan struct{}, 0)
	ch <- struct{}{}
}

var _ raft.FSM = &learningStateMachine{}

type learningStateMachine struct {
}

func (h *learningStateMachine) Apply(log *raft.Log) interface{} {
	return nil
}

func (h *learningStateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (h *learningStateMachine) Restore(io.ReadCloser) error {
	return nil
}
