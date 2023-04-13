// SPDX-License-Identifier: Apache-2.0
// Copyright 2023-present Ian Chen <ychen.cs10@nycu.edu.tw>
package session

import (
	"encoding/binary"

	"github.com/wmnsk/go-pfcp/ie"
)

//  Measurement Information IE bits definition
const (
	MBQE = 1 << iota // Measurement Before QoS Enforcement
	INAM             // Inactive Measurement
	RADI             // Reduced Application Detection Information
	ISTM             // Immediate Start Time Metering
	MNOP             // Measurement of Number of Packets
)

func newMeasurementInfo(info uint8) *ie.IE {
	if info == 0 {
		return nil
	}
	return ie.NewMeasurementInformation(info)
}

// Reporting Triggers IE bits definition
const (
	RPT_TRIG_PERIO = 1 << iota
	RPT_TRIG_VOLTH
	RPT_TRIG_TIMTH
	RPT_TRIG_QUHTI
	RPT_TRIG_START
	RPT_TRIG_STOPT
	RPT_TRIG_DROTH
	RPT_TRIG_LIUSA
	RPT_TRIG_VOLQU
	RPT_TRIG_TIMQU
	RPT_TRIG_ENVCL
	RPT_TRIG_MACAR
	RPT_TRIG_EVETH
	RPT_TRIG_EVEQU
	RPT_TRIG_IPMJL
	RPT_TRIG_QUVTI
	RPT_TRIG_REEMR
	RPT_TRIG_UPINT
)

type ReportingTrigger struct {
	Flags uint32
}

func NewRptTrig(rpgTrig ReportingTrigger) *ie.IE {
	if rpgTrig.Flags == 0 {
		return nil
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, rpgTrig.Flags)
	if b[2] != 0 {
		return ie.NewReportingTriggers(b[:3]...)
	}
	return ie.NewReportingTriggers(b[:2]...)
}

// Volume Threshold IE
type volumThreshold struct {
	flags uint8
	tvol  uint64
	uvol  uint64
	dvol  uint64
}

func newVolumeThreshold(vParams *volumThreshold) *ie.IE {
	if vParams == nil {
		return nil
	}
	return ie.NewVolumeThreshold(vParams.flags, vParams.tvol, vParams.uvol, vParams.dvol)
}

// Volume Measurement IE Flag bits definition
const (
	TOVOL uint8 = 1 << iota
	ULVOL
	DLVOL
	TONOP
	ULNOP
	DLNOP
)

// Volume Threshold IE
type volumeQuota struct {
	flags uint8
	tvol  uint64
	uvol  uint64
	dvol  uint64
}

func newVolumeQuota(vParams *volumeQuota) *ie.IE {
	if vParams == nil {
		return nil
	}
	return ie.NewVolumeQuota(vParams.flags, vParams.tvol, vParams.uvol, vParams.dvol)
}

func newMeasurementMethod(mParams *measurementMethodParams) *ie.IE {
	if mParams == nil {
		return nil
	}
	return ie.NewMeasurementMethod(mParams.event, mParams.volum, mParams.durat)
}
