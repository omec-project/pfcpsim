// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	"github.com/jessevdk/go-flags"
	pb "github.com/omec-project/pfcpsim/api"
	log "github.com/sirupsen/logrus"
)

type commonArgs struct {
	Count int `short:"c" long:"count" default:"1" description:"The number of sessions to create"`
	BaseID int `short:"i" long:"baseID"  default:"1" description:"The base ID to use"`
	UePool string `short:"u" long:"ue-pool" default:"17.0.0.0/24" description:"The UE pool address"`
	GnBAddress string `short:"g" long:"gnb-addr" description:"The (g/e)nodeB address"`
	SDFfilter string `short:"s" long:"sdf-filter" description:"The SDF Filter to use"`
	QFI uint8 `short:"q" long:"qfi" description:"The QFI value for QERs. Max value 64."`
}

type sessionCreate struct {
	Args struct{
		commonArgs
	}
}

type sessionModify struct {
	Args struct {
		commonArgs
		BufferFlag bool `short:"b" long:"buffer" description:"If set, downlink FARs will have the buffer flag set to true"`
		NotifyCPFlag bool `short:"n" long:"notifycp" description:"If set, downlink FARs will have the notify CP flag set to true"`
	}
}

type sessionDelete struct {
	Args struct{
		Count int `short:"c" long:"count" default:"1" description:"The number of sessions to create"`
		BaseID int `short:"i" long:"baseID"  default:"1" description:"The base ID to use"`
	}
}

type SessionOptions struct {
	Create sessionCreate `command:"create"`
	Modify sessionModify `command:"modify"`
	Delete sessionDelete `command:"delete"`
}

func RegisterSessionCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("session", "Handle sessions", "Command to create/modify/delete sessions", &SessionOptions{})
}

func (s *sessionCreate) Execute(args []string) error {
	if s.Args.QFI > 64 {
		log.Fatalf("QFI cannot be greater than 64. Provided QFI: %v", s.Args.QFI)
	}

	client := connect()

	res, err := client.CreateSession(context.Background(), &pb.CreateSessionRequest{
		Count:         int32(s.Args.Count),
		BaseID:        int32(s.Args.BaseID),
		NodeBAddress:  s.Args.GnBAddress,
		UeAddressPool: s.Args.UePool,
		SdfFilter:     s.Args.SDFfilter,
		Qfi: int32(s.Args.QFI),
	})

	if err != nil {
		log.Fatalf("Error while creating sessions: %v", err)
	}

	log.Infof(res.Message)

	return nil
}

func (s *sessionModify) Execute(args []string) error {
	client := connect()

	res, err := client.ModifySession(context.Background(), &pb.ModifySessionRequest{
		Count:         int32(s.Args.Count),
		BaseID:        int32(s.Args.BaseID),
		NodeBAddress:  s.Args.GnBAddress,
		UeAddressPool: s.Args.UePool,
		BufferFlag:    s.Args.BufferFlag,
		NotifyCPFlag:  s.Args.NotifyCPFlag,
	})

	if err != nil {
		log.Fatalf("Error while modifying sessions: %v", err)
	}

	log.Infof(res.Message)

	return nil
}

func (s *sessionDelete) Execute(args []string) error {
	client := connect()

	res, err := client.DeleteSession(context.Background(), &pb.DeleteSessionRequest{
		Count:  int32(s.Args.Count),
		BaseID: int32(s.Args.BaseID),
	})
	if err != nil {
		log.Fatalf("Error while deleting sessions: %v", err)
	}

	log.Infof(res.Message)

	return nil
}
