package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/chrishenyard/grpc_api_go/protos"
)

type jobServer struct {
	pb.UnimplementedJobsServer

	mu   sync.RWMutex
	jobs map[string]*pb.JobResponse
}

func newJobServer() *jobServer {
	return &jobServer{
		jobs: make(map[string]*pb.JobResponse),
	}
}

func (s *jobServer) CreateJob(ctx context.Context, req *pb.JobCreateRequest) (*pb.JobResponse, error) {
	if req.GetJobName() == "" {
		return nil, status.Error(codes.InvalidArgument, "job_name is required")
	}

	id := uuid.NewString()
	job := &pb.JobResponse{
		JobId:          id,
		JobName:        req.GetJobName(),
		JobDescription: req.GetJobDescription(),
		Status:         req.GetStatus(),
	}

	s.mu.Lock()
	s.jobs[id] = job
	s.mu.Unlock()

	return job, nil
}

func (s *jobServer) GetJob(ctx context.Context, req *pb.JobRequest) (*pb.JobResponse, error) {
	s.mu.RLock()
	job, ok := s.jobs[req.GetJobId()]
	s.mu.RUnlock()

	if !ok {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	return job, nil
}

func (s *jobServer) GetJobs(req *pb.JobListOptions, stream pb.Jobs_GetJobsServer) error {
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 100
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, job := range s.jobs {
		if count >= limit {
			break
		}
		if err := stream.Send(job); err != nil {
			return err
		}
		count++
	}

	return nil
}

func (s *jobServer) DeleteJob(ctx context.Context, req *pb.JobRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[req.GetJobId()]; !ok {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	delete(s.jobs, req.GetJobId())
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobsServer(grpcServer, newJobServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
