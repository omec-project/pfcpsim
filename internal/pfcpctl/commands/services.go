package commands

import (
	"context"

	"github.com/jessevdk/go-flags"
	pb "github.com/omec-project/pfcpsim/api"
	log "github.com/sirupsen/logrus"
)

type Associate struct {}
type Disassociate struct {}
type ConfigureRemoteAddresses struct {
	RemotePeerAddress string `short:"r" long:"remote-peer-addr" default:"" description:"The remote PFCP agent address."`
	N3InterfaceAddress string `short:"n" long:"n3-addr" default:"" description:"The IPv4 address of the UPF's N3 interface"`
}


type ServiceOptions struct {
	Associate Associate `command:"associate"`
	Disassociate Disassociate `command:"disassociate"`
	Configure ConfigureRemoteAddresses `command:"configure"`
}

func RegisterServiceCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("service", "configure pfcpsim", "Command to configure pfcpsim", &ServiceOptions{})
}

func (c *ConfigureRemoteAddresses) Execute(args []string) error {
	client, _ := connect()

	res, err := client.Configure(context.Background(), &pb.ConfigureRequest{
		UpfN3Address:      c.N3InterfaceAddress,
		RemotePeerAddress: c.RemotePeerAddress,
	})
	if err != nil {
		log.Fatalf("Error while configuring remote addresses: %v", err)
	}

	log.Info(res.Message)

	return nil
}

func (c *Associate) Execute(args []string) error {
	client, _ := connect()

	res, err := client.Associate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("Error while associating: %v", err)
	}

	log.Infof(res.Message)

	return nil
}

func (c *Disassociate) Execute(args []string) error {
	client, _ := connect()

	res, err := client.Associate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("Error while disassociating: %v", err)
	}

	log.Infof(res.Message)

	return nil
}
