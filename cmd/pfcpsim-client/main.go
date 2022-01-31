// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/c-robinson/iplib"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
)

var (
	remotePeerAddress net.IP
	upfAddress        net.IP
	nodeBAddress      net.IP
	ueAddressPool     string

	lastUEAddress net.IP

	inputFile  string
	outputFile string

	bufferFlag   bool
	notifyCPFlag bool

	sessionCount int

	mode string

	// Emulates 5G SMF/ 4G SGW
	globalPFCPSimClient *pfcpsim.PFCPClient
)

const (
	// Values for mock-up4 environment
	defaultGNodeBAddress = "198.18.0.10"
	defaultUeAddressPool = "17.0.0.0/24"

	defaultUpfN3Address = "198.18.0.1"
)

// copyOutputToLogfile reads from Stdout and Stderr to save in a persistent file,
// provided through logfile parameter.
func copyOutputToLogfile() func() {
	f, _ := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	out := os.Stdout
	multiWriter := io.MultiWriter(out, f)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	// writes with log.Print should also write to multiWriter
	log.SetOutput(multiWriter)

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_, _ = io.Copy(multiWriter, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel. this will let multiWriter finish writing before the program exits
		_ = w.Close()
		<-exit

		_ = f.Close()
	}

}

// getLocalAddress discovers local IP by retrieving interface used for default gateway.
// Returns error if fail occurs at any stage.
func getLocalAddress() (net.IP, error) {
	cmd := "route -n get default | grep 'interface:' | grep -o '[^ ]*$'"
	cmdOutput, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	interfaceName := strings.TrimSuffix(string(cmdOutput[:]), "\n")

	itf, _ := net.InterfaceByName(interfaceName)
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil { //Verify if IP is IPV4
				ip = v.IP
			}
		}
	}

	if ip != nil {
		return ip, nil
	}

	return nil, fmt.Errorf("could not find interface: %v", interfaceName)
}

// parseArgs perform flag parsing and validation saving necessary data to global variables.
func parseArgs() {
	inputF := getopt.StringLong("input-file", 'f', "", "File to poll for input commands. Default is stdin")
	outputF := getopt.StringLong("output-file", 'o', "", "File in which copy from Stdout. Default uses only Stdout")
	remotePeer := getopt.StringLong("remote-peer-address", 'r', "127.0.0.1", "Address or hostname of the remote peer (PFCP Agent)")
	upfAddr := getopt.StringLong("upf-address", 'u', defaultUpfN3Address, "Address of the UPF (UP4)")
	ueAddrPool := getopt.StringLong("ue-address-pool", 'e', defaultUeAddressPool, "The IPv4 CIDR prefix from which UE addresses will be generated, incrementally")
	NodeBAddr := getopt.StringLong("nodeb-address", 'g', defaultGNodeBAddress, "The IPv4 of (g/e)NodeBAddress")
	mode = *getopt.Enum('m', []string{"4G", "5G"}, "4G", "Set the execution mode of pfcpsim-client to 4G or 5G. Default is 4G")

	sessionCount = *getopt.IntLong("session-count", 'c', 1, "Set the amount of sessions to create, starting from 1 (included)")
	notifyCPFlag = *getopt.BoolLong("notify-cp", 'n', "If this argument is present, downlink FARs will have the notify CP flag set to true")
	bufferFlag = *getopt.BoolLong("buffer", 'b', "If this argument is present, downlink FARs will have the buffering flag set to true")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	// Flag checks and validations

	if *outputF != "" {
		outputFile = *outputF
	}

	if *inputF != "" {
		inputFile = *inputF
	}

	if sessionCount <= 0 {
		log.Fatalf("Session count cannot be 0 or a negative number")
	}

	// IPs checks
	nodeBAddress = net.ParseIP(*NodeBAddr)
	if nodeBAddress == nil {
		log.Fatalf("Could not retrieve IP address of (g/e)NodeB")
	}

	remotePeerAddress = net.ParseIP(*remotePeer)
	if remotePeerAddress == nil {
		address, err := net.LookupHost(*remotePeer)
		if err != nil {
			log.Fatalf("Could not retrieve hostname or address for remote peer: %s", *remotePeer)
		}
		remotePeerAddress = net.ParseIP(address[0])
	}

	upfAddress = net.ParseIP(*upfAddr)
	if upfAddress == nil {
		log.Fatalf("Error while parsing UPF address")
	}

	_, _, err := net.ParseCIDR(*ueAddrPool)
	if err != nil {
		log.Fatalf("Could not parse ue address pool: %v", err)
	}
	ueAddressPool = *ueAddrPool
}

// readInput will cycle through user's input. if inputFile was provided as a flag, Stdin redirection is performed.
func readInput(input chan<- string) {
	if inputFile != "" {
		// Set inputFile as stdIn

		oldStdin := os.Stdin
		defer func() {
			// restore StdIN
			os.Stdin = oldStdin
		}()

		f, err := os.Open(inputFile)
		if err != nil {
			log.Errorf("Error while reading inputFile: %v", err)
		} else {
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Errorf("Error while closing file: %v", err)
				}
			}(f)

			os.Stdin = f
		}
	}

	for {
		var u string
		_, err := fmt.Scanf("%s\n", &u)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Infof("Skipping bad entry: %v", err)
			}
		}
		input <- u
	}
}

// handleUserInput spawn a goroutine cycling through user's input.
func handleUserInput() {
	userInput := make(chan string)
	go readInput(userInput)

	for {
		fmt.Println("'disassociate': Teardown Association")
		fmt.Println("'associate': Setup Association")
		fmt.Println("'create': Create Sessions ")
		fmt.Println("'modify': Modify Sessions ")
		fmt.Println("'delete': Delete Sessions ")
		fmt.Println("'exit': Exit ")
		fmt.Print("Enter service: ")

		select {
		case userAnswer := <-userInput:
			switch userAnswer {
			case "disassociate":
				log.Info("Selected teardown association")
				err := globalPFCPSimClient.TeardownAssociation()
				if err != nil {
					log.Errorf("Error while tearing down association: %v", err)
					break
				}

				log.Infoln("Teardown association completed")

			case "associate":
				log.Info("Selected setup association")
				err := globalPFCPSimClient.SetupAssociation()
				if err != nil {
					log.Errorf("Error while setting up association: %v", err)
					break
				}

				log.Infof("Setup association completed")

			case "create":
				log.Info("Selected create sessions")
				createSessions(sessionCount)

			case "modify":
				log.Info("Selected modify session")
				err := globalPFCPSimClient.ModifySessions(notifyCPFlag, bufferFlag, &nodeBAddress)
				if err != nil {
					log.Errorf("Error while modifying sessions: %v", err)
				}

				log.Info("Modified all sessions")

			case "delete":
				log.Info("Selected delete sessions")
				err := globalPFCPSimClient.DeleteAllSessions()
				if err != nil {
					log.Errorf("Error while deleting sessions: %v", err)
					break
				}

				log.Infof("Deleted all sessions")

			case "exit":
				log.Info("Shutting down")
				globalPFCPSimClient.TeardownAssociation()
				globalPFCPSimClient.DisconnectN4()
				os.Exit(0)

			default:
				log.Error("Command not found")
			}
		}
	}
}

// getNextUEAddress retrieves the next available IP address from ueAddressPool
func getNextUEAddress() net.IP {
	if lastUEAddress != nil {
		lastUEAddress = iplib.NextIP(lastUEAddress)
		return lastUEAddress
	}

	// TODO handle case net IP is full
	ueIpFromPool, _, _ := net.ParseCIDR(ueAddressPool)
	lastUEAddress = iplib.NextIP(ueIpFromPool)
	return lastUEAddress
}

// createSessions create 'count' sessions incrementally.
// Once created, the sessions are established through PFCP client.
func createSessions(count int) {
	// SEIDs and TEIDs are generated based on number of active sessions
	baseID := globalPFCPSimClient.GetNumActiveSessions() + 1

	for i := baseID; i < ((count) + baseID); i++ {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i + 10)
		downlinkTEID := uint32(i + 11)

		uplinkFarID := uint32(i)
		downlinkFarID := uint32(i + 1)

		uplinkPdrID := uint16(i)
		dowlinkPdrID := uint16(i + 1)

		sessQerID := uint32(i + 3)

		appQerID := uint32(i)

		uplinkAppQerID := appQerID
		downlinkAppQerID := appQerID + 1

		uplinkPDR := session.NewPDRBuilder().
			WithID(uplinkPdrID).
			WithMEthod(session.Create).
			WithTEID(uplinkTEID).
			WithRulesIDs(uplinkFarID, sessQerID, uplinkAppQerID).
			WithN3Address(upfAddress.String()).
			WithSDFFilter("permit out ip from any to assigned").
			BuildPDR()

		DownlinkPDR := session.NewPDRBuilder().
			WithID(dowlinkPdrID).
			WithMEthod(session.Create).
			WithRulesIDs(downlinkFarID, sessQerID, downlinkAppQerID).
			WithPrecedence(100).
			WithUEAddress(getNextUEAddress().String()).
			WithSDFFilter("permit out ip from any to assigned").
			MarkAsDownlink().
			BuildPDR()

		uplinkFAR := session.NewFARBuilder().
			WithID(uplinkFarID).
			WithAction(session.ActionForward).
			WithMethod(session.Create).
			BuildFAR()

		downlinkFAR := session.NewFARBuilder().
			WithID(downlinkFarID).
			WithAction(session.ActionDrop).
			WithMethod(session.Create).
			WithTEID(downlinkTEID).
			WithDownlinkIP(nodeBAddress.String()).
			MarkAsDownlink().
			BuildFAR()

		sessionQER := session.NewQERBuilder().
			WithID(sessQerID).
			WithMethod(session.Create).
			WithQFI(0x09).
			WithUplinkMBR(50000).
			WithDownlinkMBR(50000).
			Build()

		appQER := session.NewQERBuilder().
			WithID(appQerID).
			WithMethod(session.Create).
			WithQFI(0x08).
			WithUplinkMBR(50000).
			WithUplinkGBR(50000).
			WithDownlinkMBR(30000).
			WithUplinkGBR(30000).
			Build()

		sess := session.NewSession()
		// add session rules to new session object
		sess.DownlinkTEID = downlinkTEID

		sess.UplinkPDRs = append(sess.UplinkPDRs, uplinkPDR)
		sess.DownlinkPDRs = append(sess.DownlinkPDRs, DownlinkPDR)

		sess.UplinkFARs = append(sess.UplinkFARs, uplinkFAR)
		sess.DownlinkFARs = append(sess.DownlinkFARs, downlinkFAR)

		sess.QERs = append(sess.QERs, sessionQER)
		sess.QERs = append(sess.QERs, appQER)

		err := globalPFCPSimClient.EstablishSession(sess)
		if err != nil {
			log.Errorf("Error while establishing sessions: %v", err)
			return
		}

		log.Infof("Created session with SEID %v", sess.LocalSEID)
	}

}

func main() {
	parseArgs()

	if outputFile != "" {
		stopLogToFile := copyOutputToLogfile()
		defer stopLogToFile()
	}

	localAddress, err := getLocalAddress()
	if err != nil {
		log.Fatalf("Error while retrieving local address: %v", err)
	}

	globalPFCPSimClient = pfcpsim.NewPFCPClient(localAddress.String())

	err = globalPFCPSimClient.ConnectN4(remotePeerAddress.String())
	if err != nil {
		log.Fatalf("Failed to connect to remote peer: %v", err)
	}

	log.Infof("PFCP client is connected")

	handleUserInput()
}
