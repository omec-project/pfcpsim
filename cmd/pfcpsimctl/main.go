// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package main

import (
	"context"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connect() (pb.PFCPSimClient, *grpc.ClientConn) {
	serverAddress := ":54321" //TODO make this configurable

	// Create an insecure gRPC Channel
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error dialing %v: %v", serverAddress, err)
	}

	return pb.NewPFCPSimClient(conn), conn
}

func main() {
	helpMsg := "'disassociate': Teardown Association \n 'associate': Setup Association \n 'create': Create Sessions  \n 'delete': Delete Sessions \n 'exit': Exit gracefully \n"
	cmd := getopt.StringLong("command", 'c', "", helpMsg)

	getopt.Parse()

	simClient, conn := connect()
	defer conn.Close()

	switch *cmd {
	case "disassociate":
		_, err := simClient.Disassociate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Errorf("Error while disassociating: %v", err)
			break
		}

		log.Info("Disassociation completed")

	case "associate":
		_, err := simClient.Associate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info("Association completed")

	case "create":
		_, err := simClient.CreateSession(context.Background(), &pb.CreateSessionRequest{
			Count: 1, //FIXME parse this from flags
		})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info("Sessions created")

	case "modify":
		_, err := simClient.ModifySession(context.Background(), &pb.ModifySessionRequest{})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info("Sessions modified")

	case "delete":
		_, err := simClient.DeleteSession(context.Background(), &pb.DeleteSessionRequest{})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info("Sessions deleted")

	default:
		log.Error("Command not recognized")
		break

	}
}
