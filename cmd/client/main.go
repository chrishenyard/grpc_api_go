package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-jobs/protos"
)

func main() {
	log.Println("Starting gRPC client...")

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to gRPC server")

	c := pb.NewJobsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := c.CreateJob(ctx, &pb.JobCreateRequest{
		JobName:        "Human Response Coordinator",
		JobDescription: "Legacy",
		Status:         "You're number 1501 in the applicant pool",
	})
	if err != nil {
		log.Fatalf("CreateJob: %v", err)
	}
	log.Printf("Created: %+v\n", created)

	got, err := c.GetJob(ctx, &pb.JobRequest{JobId: created.GetJobId()})
	if err != nil {
		log.Fatalf("GetJob: %v", err)
	}
	log.Printf("GetJob: %+v\n", got)

	stream, err := c.GetJobs(ctx, &pb.JobListOptions{Limit: 10})
	if err != nil {
		log.Fatalf("GetJobs: %v", err)
	}

	log.Println("GetJobs stream:")
	for {
		job, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream recv: %v", err)
		}
		log.Printf("  %+v\n", job)
	}

	_, err = c.DeleteJob(ctx, &pb.JobRequest{JobId: created.GetJobId()})
	if err != nil {
		log.Fatalf("DeleteJob: %v", err)
	}
	log.Println("Deleted job.")
}
