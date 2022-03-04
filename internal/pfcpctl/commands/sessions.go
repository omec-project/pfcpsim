package commands

import "github.com/jessevdk/go-flags"

type ServiceCreate struct {}
type ServiceModify struct {}
type ServiceDelete struct {}

type ServiceOptions struct {
	//List ServiceList `command:"list"`
	Create ServiceCreate `command:"create"`
	Modify ServiceModify `command:"modify"`
	Delete ServiceDelete `command:"delete"`
}

func RegisterSessionCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("session", "Sessions Commands", "Commands to crete/modify/delete sessions", &ServiceOptions{})
}

func (s *ServiceCreate) Execute() {

}

func (s *ServiceModify) Execute() {

}

func (s *ServiceDelete) Execute() {

}
