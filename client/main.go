package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/carlosfgti/go-grpc/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new task
	createResp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       "Complete gRPC tutorial",
		Description: "Finish writing the gRPC tutorial with Go and MongoDB",
		DueDate:     timestamppb.New(time.Now().Add(24 * time.Hour)),
	})
	if err != nil {
		log.Fatalf("Could not create task: %v", err)
	}
	fmt.Printf("Created task: %v\n", createResp)

	// Get the created task
	getResp, err := client.GetTask(ctx, &pb.GetTaskRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Could not get task: %v", err)
	}
	fmt.Printf("Got task: %v\n", getResp)

	// Update the task
	updateResp, err := client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:          createResp.Id,
		Title:       "Complete gRPC tutorial (updated)",
		Description: "Finish writing and reviewing the gRPC tutorial with Go and MongoDB",
		Completed:   true,
		DueDate:     timestamppb.New(time.Now().Add(48 * time.Hour)),
	})
	if err != nil {
		log.Fatalf("Could not update task: %v", err)
	}
	fmt.Printf("Updated task: %v\n", updateResp)

	// List tasks
	listResp, err := client.ListTasks(ctx, &pb.ListTasksRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("Could not list tasks: %v", err)
	}
	fmt.Printf("Listed %d tasks, total count: %d\n", len(listResp.Tasks), listResp.TotalCount)
	for _, task := range listResp.Tasks {
		fmt.Printf("- %s: %s\n", task.Id, task.Title)
	}

	// Delete the task
	deleteResp, err := client.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Could not delete task: %v", err)
	}
	fmt.Printf("Deleted task, success: %v\n", deleteResp.Success)
}
