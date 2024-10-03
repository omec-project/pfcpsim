// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	"github.com/jessevdk/go-flags"
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/logger"
)

type commonArgs struct {
	Count           int      `short:"c" long:"count" default:"1" description:"The number of sessions to create"`
	BaseID          int      `short:"i" long:"baseID"  default:"1" description:"The base ID to use"`
	UePool          string   `short:"u" long:"ue-pool" default:"17.0.0.0/24" description:"The UE pool address"`
	GnBAddress      string   `short:"g" long:"gnb-addr" description:"The UE pool address"`
	AppFilterString []string `short:"a" long:"app-filter" default:"ip:any:any:allow:100" description:"Specify an application filter. Format: '{ip | udp | tcp}:{IPv4 Prefix | any}:{<lower-L4-port>-<upper-L4-port> | any}:{allow | deny}:{rule-precedence}' . e.g. 'udp:10.0.0.0/8:80-88:allow:100'"`
	QFI             uint8    `short:"q" long:"qfi" description:"The QFI value for QERs. Max value 64."`
}

func (a *commonArgs) validate() {
	if a.BaseID <= 0 {
		logger.PfcpsimLog.Fatalln("baseID cannot be 0 or a negative number")
	}

	if a.Count <= 0 {
		logger.PfcpsimLog.Fatalln("count cannot be 0 or a negative number")
	}
}

type sessionCreate struct {
	Args struct {
		commonArgs
	}
}

type sessionModify struct {
	Args struct {
		commonArgs
		BufferFlag   bool `short:"b" long:"buffer" description:"If set, downlink FARs will have the buffer flag set to true"`
		NotifyCPFlag bool `short:"n" long:"notifycp" description:"Set true to have downlink FARs notify CP"`
	}
}

type sessionDelete struct {
	Args struct {
		commonArgs
	}
}

type SessionOptions struct {
	Create sessionCreate `command:"create"`
	Modify sessionModify `command:"modify"`
	Delete sessionDelete `command:"delete"`
}

func RegisterSessionCommands(parser *flags.Parser) {
	_, err := parser.AddCommand(
		"session",
		"Handle sessions",
		"Command to create/modify/delete sessions",
		&SessionOptions{},
	)
	if err != nil {
		logger.PfcpsimLog.Warnln(err)
	}
}

func (s *sessionCreate) Execute(args []string) error {
	if s.Args.QFI > 64 {
		logger.PfcpsimLog.Fatalf("qfi cannot be greater than 64. Provided qfi: %v", s.Args.QFI)
	}

	client := connect()

	defer disconnect()

	s.Args.validate()

	res, err := client.CreateSession(context.Background(), &pb.CreateSessionRequest{
		Count:         int32(s.Args.Count),
		BaseID:        int32(s.Args.BaseID),
		NodeBAddress:  s.Args.GnBAddress,
		UeAddressPool: s.Args.UePool,
		AppFilters:    s.Args.AppFilterString,
		Qfi:           int32(s.Args.QFI),
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while creating sessions: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)

	return nil
}

func (s *sessionModify) Execute(args []string) error {
	client := connect()

	defer disconnect()

	s.Args.validate()

	res, err := client.ModifySession(context.Background(), &pb.ModifySessionRequest{
		Count:         int32(s.Args.Count),
		BaseID:        int32(s.Args.BaseID),
		NodeBAddress:  s.Args.GnBAddress,
		UeAddressPool: s.Args.UePool,
		BufferFlag:    s.Args.BufferFlag,
		NotifyCPFlag:  s.Args.NotifyCPFlag,
		AppFilters:    s.Args.AppFilterString,
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while modifying sessions: %v", err)
	}

	logger.PfcpsimLog.Infof(res.Message)

	return nil
}

func (s *sessionDelete) Execute(args []string) error {
	client := connect()

	defer disconnect()

	s.Args.validate()

	res, err := client.DeleteSession(context.Background(), &pb.DeleteSessionRequest{
		Count:  int32(s.Args.Count),
		BaseID: int32(s.Args.BaseID),
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while deleting sessions: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)

	return nil
}
