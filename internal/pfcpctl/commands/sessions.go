package commands

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type commonArgs struct {
	Count int `short:"c" long:"count" default:"1" description:"The number of sessions to create"`
	BaseID int `short:"i" long:"baseID"  default:"1" description:"The base ID to use"`
	UePool string `short:"u" long:"ue-pool" default:"17.0.0.0/24" description:"The UE pool address"`
	GnBAddress string `short:"g" long:"gnb-addr" description:"The UE pool address"`
	SDFfilter string `short:"s" long:"sdf-filter" description:"The SDF Filter to use"`
	QFI uint8 `short:"q" long:"qfi" description:"The QFI value for QERs. Max value 64."`
}

type ServiceCreate struct {
	Args struct{
		commonArgs
	}
}

type ServiceModify struct {
	Args struct {
		commonArgs
		bufferFlag bool `short:"b" long:"buffer" description:"If set, downlink FARs will have the buffer flag set to true"`
		notifyCPFlag bool `short:"n" long:"notifyCP" description:"If set, downlink FARs will have the notify CP flag set to true"`
	}
}

type ServiceDelete struct {
	Count int
	BaseID int
}

type SessionOptions struct {
	Create ServiceCreate `command:"create"`
	Modify ServiceModify `command:"modify"`
	Delete ServiceDelete `command:"delete"`
}

func RegisterSessionCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("session", "Handle sessions", "Command to create/modify/delete sessions", &SessionOptions{})
}

func (s *ServiceCreate) Execute(args []string) error {
	fmt.Printf("Selected create service, with options:\n count: %d , baseID: %d", s.Args.Count, s.Args.BaseID)

	if s.Args.QFI > 64 {
		log.Fatalf("QFI cannot be greater than 64. Provided QFI: %v", s.Args.QFI)
	}

	return nil
}

func (s *ServiceModify) Execute(args []string) error {
	fmt.Println("Selected modify service")

	return nil
}

func (s *ServiceDelete) Execute(args []string) error {
	fmt.Println("Selected delete service")

	return nil
}
