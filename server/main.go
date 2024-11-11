package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/kanerix/mutual-exclusion/lamport"
	pb "github.com/kanerix/mutual-exclusion/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedMutualExclusionServer
	mu           sync.Mutex
	nodeID       string
	isProcessing bool
	queue        []request
	clock        *lamport.LamportClock
}

type request struct {
	nodeID    string
	timestamp uint64
}

var nodes = []string{"localhost:3001", "localhost:3002", "localhost:3003"}

func (s *server) RequestAccess(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clock.Max(req.Timestamp)

	s.queue = append(s.queue, request{
		req.NodeId,
		req.Timestamp,
	})
	fmt.Printf("[%s] Access Requested by %s\n", s.nodeID, req.NodeId)

	sort.Slice(s.queue, func(i, j int) bool {
		if s.queue[i].timestamp == s.queue[j].timestamp {
			return s.queue[i].nodeID < s.queue[j].nodeID
		}
		return s.queue[i].timestamp < s.queue[j].timestamp
	})

	if !s.isProcessing && s.queue[0].nodeID == req.NodeId {
		s.isProcessing = true
		s.queue = s.queue[1:]

		fmt.Printf("[%s] Access Granted to %s\n", s.nodeID, req.NodeId)
		return &pb.Response{Granted: true}, nil
	}

	return &pb.Response{Granted: false}, nil
}

func (s *server) ReleaseAccess(ctx context.Context, release *pb.Release) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clock.Step()

	fmt.Printf("[%s] Access Released by %s\n", s.nodeID, release.NodeId)
	if len(s.queue) > 0 {
		go s.notify(s.queue[0].nodeID)
	}

	s.isProcessing = false
	return &pb.Response{Granted: true}, nil
}

func (s *server) notify(targetNode string) {
	conn, err := grpc.NewClient(targetNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMutualExclusionClient(conn)
	_, err = client.RequestAccess(context.Background(), &pb.Request{NodeId: targetNode})

	if err != nil {
		log.Fatalf("could not request access: %v", err)
	}
}

func newServer(nodeID string) *server {
	s := &server{
		queue:  make([]request, 0),
		nodeID: nodeID,
		clock:  lamport.NewLamportClock(),
	}
	return s
}

func startServer(nodeID, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMutualExclusionServer(s, newServer(nodeID))

	go func() {
		fmt.Printf("Server %s listening on %v\n", nodeID, address)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Printf("\nStopping Server %s\n", nodeID)
	s.Stop()
}

func requestCriticalSection(nodeID string) {
	clock := lamport.NewLamportClock()

	for _, addr := range nodes {
		if addr == fmt.Sprintf("localhost:%s", nodeID) {
			continue
		}

		clock.Step()

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewMutualExclusionClient(conn)
		resp, err := client.RequestAccess(context.Background(), &pb.Request{NodeId: nodeID, Timestamp: clock.Now()})
		if err != nil || !resp.Granted {
			fmt.Printf("Node %s request denied by %s\n", nodeID, addr)
		}

		if resp.Granted {
			fmt.Printf("Node %s request granted by %s\n", nodeID, addr)
			break
		}
	}

	time.Sleep(1 * time.Second)
}

func releaseCriticalSection(nodeID string) {
	for _, addr := range nodes {
		if addr == fmt.Sprintf("localhost:%s", nodeID) {
			continue
		}

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewMutualExclusionClient(conn)
		client.ReleaseAccess(context.Background(), &pb.Release{NodeId: nodeID})

		fmt.Printf("Node %s released access\n", nodeID)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide a port for the server.")
	}

	nodeID := os.Args[1]

	go startServer(nodeID, fmt.Sprintf("localhost:%s", nodeID))

	time.Sleep(2 * time.Second)

	requestCriticalSection(nodeID)

	time.Sleep(5 * time.Second)

	releaseCriticalSection(nodeID)
}
