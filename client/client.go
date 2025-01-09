package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Manjeet2001/GoCache/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValueServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	setResp, err := client.Set(ctx, &pb.SetRequest{
		Key:   "test-key",
		Value: "test-value",
	})
	if err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}
	log.Printf("Set response: %v", setResp.Success)

	getResp, err := client.Get(ctx, &pb.GetRequest{
		Key: "test-key",
	})
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	log.Printf("Get response: value=%s, found=%v", getResp.Value, getResp.Found)

	deleteResp, err := client.Delete(ctx, &pb.DeleteRequest{
		Key: "test-key",
	})
	if err != nil {
		log.Fatalf("Failed to delete key: %v", err)
	}
	log.Printf("Delete response: %v", deleteResp.Success)

	addNodeResp, err := client.AddNode(ctx, &pb.NodeRequest{
		Address: "localhost:50052",
	})
	if err != nil {
		log.Fatalf("Failed to add node: %v", err)
	}
	log.Printf("Add node response: %v", addNodeResp.Success)
}
