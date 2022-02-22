/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
)

const (
	maxPacketSize = 1500
	timeout = -1 * time.Second // sniff forever
)

func connectPFCPSim() error {
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

func isConfigured() bool {
	if upfN3Address != "" && remotePeerAddress != "" {
		return true
	}

	return false
}

func isRemotePeerConnected() bool {
	return remotePeerConnected
}

// getLocalAddress returns the first IP address of the interfaceName, if specified,
// otherwise returns the IP address of the first non-loopback interface
// Returns error if fail occurs at any stage.
func getLocalAddress(ifaceName string) (net.IP, error) {
	if ifaceName != "" {
		// Interface name is specified. Use it.
		interfaceName = ifaceName
		interfaceAddrs, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return nil, err
		}

		addrs, _ := interfaceAddrs.Addrs()
		for _, address := range addrs {
			// Check address type to be non-loopback
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP, nil
				}
			}
		}
	}
	// Interface name is not provided. Use first non-loopback one
	ifaceName, IpAddress, err := getFirstNonLoopbackInterface()
	if err != nil {
		return nil, err
	}
	interfaceName = ifaceName

	return IpAddress, nil
}

// getFirstNonLoopbackInterface returns the first non-loopback interface name and IP address.
// returns error if failure happens at any stage.
func getFirstNonLoopbackInterface() (string, net.IP, error){
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return i.Name, ipnet.IP, nil
				}
			}
		}
	}

	return "", nil, pfcpsim.NewNoValidInterfaceError()
}

// sniffer starts sniffing on interfaceName.
// if pcapPath is empty, returns immediately.
func sniffer(doneChannel chan bool) error {
	if pcapPath == "" {
		return nil
	}

	// Open output pcap file
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(pcapPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w := pcapgo.NewWriter(f)
	// Write header
	if err = w.WriteFileHeader(maxPacketSize, layers.LinkTypeEthernet); err != nil {
		return err
	}

	defer f.Close()

	// Open the device for capturing
	handle, err := pcap.OpenLive(interfaceName, maxPacketSize, true, timeout)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Start processing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	select {
	case <-doneChannel:
		return nil
	default:
		for packet := range packetSource.Packets() {
			if err = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
				return err
			}
		}
	}

	return nil
}
