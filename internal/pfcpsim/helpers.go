// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	log "github.com/sirupsen/logrus"
	"github.com/wmnsk/go-pfcp/ie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var notInit = errors.New("PFCP simulator is not initialized")

const (
	sdfFilterFormatWPort  = "permit out %v from %v to assigned %v-%v"
	sdfFilterFormatWOPort = "permit out %v from %v to assigned"
)

func GetSimulator() *pfcpsim.PFCPClient {
	return sim
}

func ConnectPFCPSim() error {
	if sim == nil {
		localAddr, err := getLocalAddress(interfaceName)
		if err != nil {
			return err
		}

		sim = pfcpsim.NewPFCPClient(localAddr.String())
	}

	err := sim.ConnectN4(remotePeerAddress)
	if err != nil {
		return err
	}

	remotePeerConnected = true

	return nil
}

func DisconnectPFCPSim() error {
	if sim == nil {
		return notInit
	}

	return sim.TeardownAssociation()
}

func isConfigured() bool {
	if upfN3Address != "" && remotePeerAddress != "" {
		return true
	}

	return false
}

func isRemotePeerConnected() bool {
	return remotePeerConnected
}

// isNumOfAppFiltersCorrect returns error if the number of the passed filter
// exceed the max number of supported application filters.
func isNumOfAppFiltersCorrect(filters []string) error {
	if len(filters) > SessionStep/2 {
		log.Errorf("Too many application filters: %v", filters)
		return status.Error(codes.Aborted, "Too many application filters")
	}

	return nil
}

// getLocalAddress returns the first IP address of the interfaceName, if
// specified, otherwise returns the IP address of the first non-loopback
// interface
// Returns error if fail occurs at any stage.
func getLocalAddress(interfaceName string) (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	if interfaceName != "" {
		// Interface name is specified. Use it.
		interfaceAddrs, err := net.InterfaceByName(interfaceName)
		if err != nil {
			return nil, err
		}

		addrs, err = interfaceAddrs.Addrs()
		if err != nil {
			return nil, err
		}
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, pfcpsim.NewNoValidInterfaceError()
}

// ParseAppFilter parses an application filter. Returns a tuple formed by a
// formatted SDF filter and a uint8 representing the Application QER gate status
// and a precedence. Returns error if fail occurs while validating the filter
// string.
func ParseAppFilter(filter string) (string, uint8, uint32, error) {
	if filter == "" {
		// parsing a wildcard app filter
		return "", ie.GateStatusOpen, 100, nil
	}

	result := strings.Split(filter, ":")
	if len(result) != 5 {
		return "", 0, 0, pfcpsim.NewInvalidFormatError(
			"Parser was not able to generate the correct number of arguments." +
				" Please make sure to use the right format")
	}

	proto, ipNetAddr, portRange, action, precedence := result[0], result[1], result[2], result[3], result[4]

	var gateStatus uint8

	switch action {
	case "allow":
		gateStatus = ie.GateStatusOpen
	case "deny":
		gateStatus = ie.GateStatusClosed
	default:
		return "", 0, 0, pfcpsim.NewInvalidFormatError(
			"Action. Please make sure to use 'allow' or 'deny'")
	}

	if !(proto == "ip" || proto == "udp" || proto == "tcp") {
		return "", 0, 0, pfcpsim.NewInvalidFormatError("Unsupported or unknown protocol.")
	}

	precedenceConverted, err := strconv.Atoi(precedence)
	if err != nil {
		return "", 0, 0, pfcpsim.NewInvalidFormatError("Precedence. Please make sure it is a number", err)
	}

	precedenceUint := uint32(precedenceConverted)

	if ipNetAddr != "any" {
		_, _, err := net.ParseCIDR(ipNetAddr)
		if err != nil {
			return "", 0, 0, pfcpsim.NewInvalidFormatError("IP and subnet mask.", err)
		}
	}

	if portRange != "any" {
		portList := strings.Split(portRange, "-")
		if !(len(portList) == 2) {
			return "", 0, 0, pfcpsim.NewInvalidFormatError(
				"Port range. Please make sure to use dash '-' to separate the two ports",
			)
		}

		lowerPort, err := strconv.Atoi(portList[0])
		if err != nil {
			return "", 0, 0, pfcpsim.NewInvalidFormatError("Port range.", err)
		}

		upperPort, err := strconv.Atoi(portList[1])
		if err != nil {
			return "", 0, 0, pfcpsim.NewInvalidFormatError("Port range.", err)
		}

		if lowerPort > upperPort {
			return "", 0, 0, pfcpsim.NewInvalidFormatError("Port range. Lower port is greater than upper port")
		}

		return fmt.Sprintf(sdfFilterFormatWPort, proto, ipNetAddr, lowerPort, upperPort), gateStatus, precedenceUint, nil
	} else {
		return fmt.Sprintf(sdfFilterFormatWOPort, proto, ipNetAddr), gateStatus, precedenceUint, nil
	}
}
