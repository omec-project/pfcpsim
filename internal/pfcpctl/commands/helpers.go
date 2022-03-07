package commands

import (
	"fmt"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/pfcpctl/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCServerAddr = "localhost"
	gRPCServerPort = "54321"

)

func connect() (pb.PFCPSimClient, *grpc.ClientConn) {
	// Create an insecure gRPC Channel
	serverAddr := fmt.Sprintf("%s:%s", gRPCServerAddr, gRPCServerPort)
	conn, err := grpc.Dial(config.GlobalConfig.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error dialing %v: %v", serverAddr, err)
	}

	return pb.NewPFCPSimClient(conn), conn
}
