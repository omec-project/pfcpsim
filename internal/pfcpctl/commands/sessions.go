package commands

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type ServiceCreate struct {
	Args struct {
		Count int `command:"count"`
	} `required:"yes"`

	Options struct {
		BaseID int `command:"baseID"`
		UePool string `command:"ue-pool"`
		GnBAddress string `command:"gnb-addr"`
		SDFfilter string `command:"sdf_filter"`
		Qfi uint8 `command:"qfi"`
	}
}

type ServiceModify struct {
	Count int `command:"count"`
	BaseID int `command:"baseID"`
	UePool string `command:"ue-pool"`
	GnBAddress string `command:"gnb-addr"`
}

type ServiceDelete struct {
	Count int `command:"count"`
	BaseID int `command:"baseID"`
}

type ServiceOptions struct {
	//List ServiceList `command:"list"`
	Create ServiceCreate `command:"create"`
	Modify ServiceModify `command:"modify"`
	Delete ServiceDelete `command:"delete"`
}

func RegisterSessionCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("session", "Handle sessions", "Command to create/modify/delete sessions", &ServiceOptions{})
}

func (s *ServiceCreate) Execute(args []string) error {
	fmt.Println("Selected create service")

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
