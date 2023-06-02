package fuzz

import (
	"math/rand"
	"testing"
	"time"

	"github.com/omec-project/pfcpsim/internal/pfcpsim/fuzz"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
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
	defer func() {
		err = sim.TerminatePFCPSim()
		require.NoError(t, err)
	}()
	err = sim.CreateSession(2,
		session.PdrMax,
		session.QerMax,
		session.FarMax,
		session.UrrMax,
		uint(0))
	if err != nil {
		require.NoError(t, err, "CreateSession failed")
	}
}

func Fuzz(f *testing.F) {
	testcases := []uint{0, 10, 84527, 156325}
	for _, tc := range testcases {
		f.Add(tc)
	}
	session.SetCheck(false)

	f.Fuzz(func(t *testing.T, seed uint) {
		time.Sleep(5 * time.Second)
		sim := fuzz.NewPfcpSimCfg("ens18", "10.10.0.59", "127.0.0.8")
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
		err = sim.CreateSession(2, rand.Intn(session.PdrMax),
			rand.Intn(session.QerMax),
			rand.Intn(session.FarMax),
			rand.Intn(session.UrrMax),
			seed)
		if err != nil {
			require.NoError(t, err, "CreateSession failed")
		}
		err = sim.ModifySession(2,
			rand.Intn(session.FarMax),
			rand.Intn(session.UrrMax),
			seed)
		if err != nil {
			require.NoError(t, err, "ModifySession failed")
		}
		err = sim.DeleteSession(2)
		if err != nil {
			require.NoError(t, err, "DeleteSession failed")
		}
	})
}
