// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"github.com/omec-project/pfcpsim/internal/pfcpctl/commands"
	"github.com/omec-project/pfcpsim/internal/pfcpctl/config"
)

func main() {
	parser := flags.NewNamedParser(path.Base(os.Args[0]),
		flags.HelpFlag|flags.PassDoubleDash|flags.PassAfterNonOption)
	_, err := parser.AddGroup("Global Options", "", &config.GlobalOptions)
	if err != nil {
		panic(err)
	}
	// Set server address and configure other parameters
	config.ProcessGlobalOptions()

	commands.RegisterServiceCommands(parser)
	commands.RegisterSessionCommands(parser)

	_, err = parser.ParseArgs(os.Args[1:])
	if err != nil {
		_, ok := err.(*flags.Error)
		if ok {
			realF := err.(*flags.Error)
			if realF.Type == flags.ErrHelp {
				os.Stdout.WriteString(err.Error() + "\n")
				return
			}
		}

		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
