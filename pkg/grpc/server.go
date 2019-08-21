package grpc

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	queuesmetrics "github.com/danieloliveira079/laravel-queues-exporter/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ExporterServer struct {
	queuesMetricsServer *queuesMetricsServer
}

func (s *ExporterServer) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	s.queuesMetricsServer = &queuesMetricsServer{}
	queuesmetrics.RegisterQueuesMetricsServer(grpcServer, s.queuesMetricsServer)
	log.Println("GRPC ExporterServer listening on", address)
	return grpcServer.Serve(listen)
}

func (s *ExporterServer) Process(metrics []metric.Metric) {
	//log.Println("gRPC Process:", metrics)
	s.queuesMetricsServer.Process(metrics)
}
