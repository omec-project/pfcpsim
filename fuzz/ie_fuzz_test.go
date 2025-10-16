// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-present Ian Chen <ychen.cs10@nycu.edu.tw>

package fuzz

import (
	"crypto/rand"
	"flag"
	"math/big"
	"testing"
	"time"

	"github.com/omec-project/pfcpsim/internal/pfcpsim/export"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

func getRand(n int) int {
	res, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return n
	}

	return int(res.Int64())
}

var (
	iface = flag.String("iface", "eth0", "the interface name you used to establish the connection with UPF.")
	upfN3 = flag.String("upfN3", "192.168.0.5", "the N3 interface address of the UPF")
	upfN4 = flag.String("upfN4", "127.0.0.8", "the N4 interface address of the UPF")
)

func Fuzz(f *testing.F) {
	var testcases []uint
	for range 100 {
		testcases = append(testcases, uint(getRand(MaxInt)))
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	session.SetCheck(false)

	f.Fuzz(func(t *testing.T, input uint) {
		time.Sleep(5 * time.Second)

		sim := export.NewPfcpSimCfg(*iface, *upfN3, *upfN4)

		err := sim.InitPFCPSim()
		if err != nil {
			t.Fatalf("InitPFCPSim failed: %v", err)
		}

		err = sim.Associate()
		if err != nil {
			t.Fatalf("Associate failed: %v", err)
		}

		defer func() {
			err = sim.TerminatePFCPSim()
			if err != nil {
				t.Errorf("TerminatePFCPSim failed: %v", err)
			}
		}()

		err = sim.CreateSession(2, getRand(session.PdrMax),
			int(input)%session.QerMax,
			int(input)%session.FarMax,
			int(input)%session.UrrMax,
			input)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		err = sim.ModifySession(2,
			getRand(session.FarMax),
			getRand(session.UrrMax),
			input)
		if err != nil {
			t.Fatalf("ModifySession failed: %v", err)
		}

		time.Sleep(3 * time.Second)

		err = sim.DeleteSession(2)
		if err != nil {
			t.Fatalf("DeleteSession failed: %v", err)
		}
	})
}
