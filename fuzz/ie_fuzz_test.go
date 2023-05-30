package fuzz

import (
	"testing"

	"github.com/omec-project/pfcpsim/internal/pfcpsim/fuzz"
	"github.com/stretchr/testify/require"
)

func TestBasicFunction(t *testing.T) {
	sim := fuzz.NewPfcpSimCfg("eth0", "192.168.0.5", "127.0.0.8")
	err := sim.InitPFCPSim()
	if err != nil {
		require.NoError(t, err, "InitPFCPSim failed")
	}
	err = sim.Associate()
	if err != nil {
		require.NoError(t, err, "Associate failed")
	}
	defer sim.TerminatePFCPSim()
	err = sim.CreateSession()
	if err != nil {
		require.NoError(t, err, "CreateSession failed")
	}
}

func FuzzIE(f *testing.F) {
	testcases := []int{0}
	for _, tc := range testcases {
		f.Add(tc) // seed corpus
	}
	f.Fuzz(func(t *testing.T, seed int) {
		// TODO: generate random CreateSessionIEs
		sim := fuzz.NewPfcpSimCfg("ens18", "10.10.0.59", "127.0.0.8")
		err := sim.InitPFCPSim()
		if err != nil {
			require.NoError(t, err, "InitPFCPSim failed")
		}
		err = sim.Associate()
		if err != nil {
			require.NoError(t, err, "Associate failed")
		}
		defer sim.TerminatePFCPSim()
		err = sim.CreateSession()
		if err != nil {
			require.NoError(t, err, "CreateSession failed")
		}
	})
}
