package main

import (
	"github.com/Manjeet2001/GoCache/pkg/cluster"
	"github.com/Manjeet2001/GoCache/pkg/store"
	pb "github.com/Manjeet2001/GoCache/proto" // Ensure the Protobuf files are generated using 'protoc --go_out=. --go-grpc_out=.' and include a definition for 'SetRequest' message.
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedKeyValueServiceServer // Confirm "protoc" was run to generate this file
	store                                 *store.KeyValueStore
	cluster                               *cluster.ConsistentHash
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	// Determine which node should handle this key
	_, ok := s.cluster.GetNode(req.Key)
	if !ok {
		return &pb.SetResponse{Success: false, Error: "No nodes available"}, nil
	}

	// In a real implementation, this would forward to the correct node
	s.store.Set(req.Key, req.Value)
	return &pb.SetResponse{Success: true}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	_, ok := s.cluster.GetNode(req.Key)
	if !ok {
		return &pb.GetResponse{Found: false}, nil
	}

	value, found := s.store.Get(req.Key)
	return &pb.GetResponse{
		Value: value,
		Found: found,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	_, ok := s.cluster.GetNode(req.Key)
	if !ok {
		return &pb.DeleteResponse{Success: false}, nil
	}

	s.store.Delete(req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *server) AddNode(ctx context.Context, req *pb.NodeRequest) (*pb.NodeResponse, error) {
	node := cluster.Node{
		ID:      req.Address,
		Address: req.Address,
	}
	s.cluster.AddNode(node)
	return &pb.NodeResponse{Success: true}, nil
}

func (s *server) RemoveNode(ctx context.Context, req *pb.NodeRequest) (*pb.NodeResponse, error) {
	s.cluster.RemoveNode(req.Address)
	return &pb.NodeResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	store := &server{
		store:   store.NewKeyValueStore(),
		cluster: cluster.NewConsistentHash(),
	}

	// Add initial node
	store.cluster.AddNode(cluster.Node{
		ID:      "node1",
		Address: "localhost:50051",
	})

	pb.RegisterKeyValueServiceServer(s, store) // Ensure the Protobuf server registration matches the generated code
	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
