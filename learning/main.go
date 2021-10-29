package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/wu-xian/raft"
)

func main() {
	id := flag.String("id", "", "")
	addr := flag.String("addr", "", "")
	flag.Parse()
	if len(*id) == 0 {
		panic("id is empty")
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(*id)
	raftConfig.HeartbeatTimeout = 2 * time.Second
	raftConfig.ElectionTimeout = 10 * time.Second

	fsm := &learningStateMachine{}
	inmemStore := raft.NewInmemStore()
	snapStore := raft.NewInmemSnapshotStore()
	tcpAddr, err := net.ResolveTCPAddr("", *addr)
	if err != nil {
		panic(err)
	}
	tcpTrans, err := raft.NewTCPTransport(*addr, tcpAddr, 2, time.Second*10, os.Stdout)
	if err != nil {
		panic(err)
	}
	ra, err := raft.NewRaft(raftConfig, fsm, inmemStore, inmemStore, snapStore, tcpTrans)
	if err != nil {
		panic(err)
	}
	raftConfiguration, err := raft.GetConfiguration(raftConfig, fsm, inmemStore, inmemStore, snapStore, tcpTrans)
	servers := []raft.Server{
		{
			Suffrage: raft.Voter,
			ID:       raft.ServerID("1"),
			Address:  raft.ServerAddress("127.0.0.1:12379"),
		}, {
			Suffrage: raft.Voter,
			ID:       raft.ServerID("2"),
			Address:  raft.ServerAddress("127.0.0.1:22379"),
		},
	}

	raftConfiguration.Servers = servers

	ra.BootstrapCluster(raftConfiguration)
	go func() {
		for {
			<-time.After(time.Second)
			if *id == "1" {
				future := ra.Apply([]byte(time.Now().Format("2006-01-02 15:04:05")), time.Second*10)
				fmt.Println("apply future", future)
			}
			fmt.Println("leader:", ra.Leader())
		}
	}()
	fmt.Println("wait")
	ch := make(chan struct{}, 0)
	ch <- struct{}{}
}

var _ raft.FSM = &learningStateMachine{}

type learningStateMachine struct {
}

func (h *learningStateMachine) Apply(log *raft.Log) interface{} {
	fmt.Println("my apply")
	fmt.Println(string(log.Data))
	return nil
}

func (h *learningStateMachine) Snapshot() (raft.FSMSnapshot, error) {
	fmt.Println("Snapshot")
	return nil, nil
}

func (h *learningStateMachine) Restore(io.ReadCloser) error {
	fmt.Println("apRestoreple")
	return nil
}
