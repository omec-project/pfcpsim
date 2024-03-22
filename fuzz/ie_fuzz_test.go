// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-present Ian Chen <ychen.cs10@nycu.edu.tw>

package fuzz

import (
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/omec-project/pfcpsim/internal/pfcpsim/export"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	"github.com/stretchr/testify/require"
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

// Example of a basic function test
// func TestBasicFunction(t *testing.T) {
// 	sim := export.NewPfcpSimCfg("eth0", "192.168.0.5", "127.0.0.8")
// 	err := sim.InitPFCPSim()
// 	if err != nil {
// 		require.NoError(t, err, "InitPFCPSim failed")
// 	}
// 	err = sim.Associate()
// 	if err != nil {
// 		require.NoError(t, err, "Associate failed")
// 	}
// 	defer func() {
// 		err = sim.TerminatePFCPSim()
// 		require.NoError(t, err)
// 	}()
// 	err = sim.CreateSession(2,
// 		session.PdrNoFuzz,
// 		session.QerNoFuzz,
// 		session.FarNoFuzz,
// 		session.UrrNoFuzz,
// 		uint(0))
// 	if err != nil {
// 		require.NoError(t, err, "CreateSession failed")
// 	}
// }

func Fuzz(f *testing.F) {
	var testcases []uint
	for i := 0; i < 100; i++ {
		testcases = append(testcases, uint(getRand(MaxInt)))
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	session.SetCheck(false)

	f.Fuzz(func(t *testing.T, input uint) {
		time.Sleep(5 * time.Second)

		sim := export.NewPfcpSimCfg("eth0", "192.168.0.5", "127.0.0.8")

		err := sim.InitPFCPSim()
		if err != nil {
			require.NoError(t, err, "InitPFCPSim failed")
		}

		err = sim.Associate()
		if err != nil {
			require.NoError(t, err, "Associate failed")
		}

		defer func() {
			err = sim.TerminatePFCPSim()
			require.NoError(t, err)
		}()

		err = sim.CreateSession(2, getRand(session.PdrMax),
			int(input)%session.QerMax,
			int(input)%session.FarMax,
			int(input)%session.UrrMax,
			input)
		if err != nil {
			require.NoError(t, err, "CreateSession failed")
		}

		err = sim.ModifySession(2,
			getRand(session.FarMax),
			getRand(session.UrrMax),
			input)
		if err != nil {
			require.NoError(t, err, "ModifySession failed")
		}

		time.Sleep(3 * time.Second)

		err = sim.DeleteSession(2)
		if err != nil {
			require.NoError(t, err, "DeleteSession failed")
		}
	})
}
