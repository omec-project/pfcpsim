// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package main

import (
	"context"
	"os"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultgRPCServerAddress = ":54321"
)

func connect(serverAddr string) (pb.PFCPSimClient, *grpc.ClientConn) {
	// Create an insecure gRPC Channel
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error dialing %v: %v", serverAddr, err)
	}

	return pb.NewPFCPSimClient(conn), conn
}

func main() {
	// TODO improve parser
	helpMsg := "'configure': Configure Server " +
		"\n 'disassociate': Teardown Association and disconnect from remote peer" +
		"\n 'associate': Connect to remote peer and setup association " +
		"\n 'create': Create Sessions  " +
		"\n 'delete': Delete Sessions " +
		"\n"
	cmd := getopt.StringLong("command", 'c', "", helpMsg)
	count := getopt.IntLong("count", 'n', 1, "The number of sessions to create/modify/delete")
	srvAddr := getopt.StringLong("server", 's', defaultgRPCServerAddress, "The gRPC address of pfcpsim")
	baseId := getopt.IntLong("baseID", 'i', 1, "First ID used to generate all other ID fields.")
	n3Addr := getopt.StringLong("n3-addr", 'a', "", "The IPv4 address of the UPF's N3 interface")
	ueAddrPool := getopt.StringLong("ue-pool", 'u', "", "The IPv4 prefix from which UE addresses will be drawn.")
	nodeBAddr := getopt.StringLong("gnb-addr", 'g', "", "The IPv4 address of the NodeB")
	remotePeer := getopt.StringLong("remote-peer", 'r', "", "The remote PFCP Agent address")

	bufferFlag := getopt.BoolLong("buffer", 'b', "If set, downlink FARs will have the buffer flag set to true")
	notifyCPFlag := getopt.BoolLong("notifycp", 'm', "If set, downlink FARs will have the notify CP flag set to true")
	sdfFilter := getopt.StringLong("sdf-filter", 'f', "" ,"Allows to set a custom SDF filter")
	qfi := getopt.Int32Long("qfi", 'q', 0, "Allows to set a custom QFI value for QERs. Max value 64")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if *qfi > 64 {
		log.Fatalf("QFI value cannot exceed 64.")
	}

	simClient, conn := connect(*srvAddr)
	defer conn.Close()

	switch *cmd {
	case "configure":
		res, err := simClient.Configure(context.Background(), &pb.ConfigureRequest{
			UpfN3Address:      *n3Addr,
			RemotePeerAddress: *remotePeer,
		})
		if err != nil {
			log.Fatalf("Error while configuring: %v", err)
		}

		log.Info(res.Message)

	case "disassociate":
		res, err := simClient.Disassociate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Fatalf("Error while disassociating: %v", err)
		}

		log.Infof(res.Message)

	case "associate":
		res, err := simClient.Associate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Fatalf("Error while associating: %v", err)
		}

		log.Infof(res.Message)

	case "create":
		res, err := simClient.CreateSession(context.Background(), &pb.CreateSessionRequest{
			Count:         int32(*count),
			BaseID:        int32(*baseId),
			NodeBAddress:  *nodeBAddr,
			UeAddressPool: *ueAddrPool,
			SdfFilter:     *sdfFilter,
			Qfi: *qfi,
		})
		if err != nil {
			log.Fatalf("Error while associating: %v", err)
		}

		log.Infof(res.Message)

	case "modify":
		res, err := simClient.ModifySession(context.Background(), &pb.ModifySessionRequest{
			Count:         int32(*count),
			BaseID:        int32(*baseId),
			NodeBAddress:  *nodeBAddr,
			UeAddressPool: *ueAddrPool,
			BufferFlag:    *bufferFlag,
			NotifyCPFlag:  *notifyCPFlag,
		})
		if err != nil {
			log.Fatalf("Error while associating: %v", err)
		}

		log.Infof(res.Message)

	case "delete":
		res, err := simClient.DeleteSession(context.Background(), &pb.DeleteSessionRequest{
			Count:  int32(*count),
			BaseID: int32(*baseId),
		})
		if err != nil {
			log.Fatalf("Error while associating: %v", err)
		}

		log.Infof(res.Message)

	default:
		log.Fatalf("Command not recognized")

	}
}
