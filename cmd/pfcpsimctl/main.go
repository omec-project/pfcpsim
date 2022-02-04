// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

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
	helpMsg := "'disassociate': Teardown Association \n 'associate': Setup Association \n 'create': Create Sessions  \n 'delete': Delete Sessions \n 'interrupt': Emulates a crash \n"
	cmd := getopt.StringLong("command", 'c', "", helpMsg)
	count := getopt.IntLong("count", 'n', 1, "The number of sessions to create/modify/delete")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	simClient, conn := connect()
	defer conn.Close()

	switch *cmd {
	case "disassociate":
		res, err := simClient.Disassociate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Errorf("Error while disassociating: %v", err)
			break
		}

		log.Info(res.Message)

	case "associate":
		res, err := simClient.Associate(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info(res.Message)

	case "create":
		res, err := simClient.CreateSession(context.Background(), &pb.CreateSessionRequest{
			Count: int32(*count),
		})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info(res.Message)

	case "modify":
		res, err := simClient.ModifySession(context.Background(), &pb.ModifySessionRequest{
			Count: int32(*count),
		})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info(res.Message)

	case "delete":
		res, err := simClient.DeleteSession(context.Background(), &pb.DeleteSessionRequest{
			Count: int32(*count),
		})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info(res.Message)

	case "interrupt":
		res, err := simClient.Interrupt(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			log.Errorf("Error while associating: %v", err)
			break
		}

		log.Info(res.Message)

	default:
		log.Error("Command not recognized")
		break

	}
}
