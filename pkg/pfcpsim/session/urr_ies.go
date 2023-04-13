// SPDX-License-Identifier: Apache-2.0
// Copyright 2023-present Ian Chen <ychen.cs10@nycu.edu.tw>
package session

import (
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
	RPT_TRIG_VOLTH // Volume Threshold
	RPT_TRIG_TIMTH // Time Threshold
	RPT_TRIG_QUHTI // Quota Holding Time
	RPT_TRIG_START // Start of Traffic
	RPT_TRIG_STOPT // Stop of Traffic
	RPT_TRIG_DROTH // Dropped DL Traffic Threshold
	RPT_TRIG_LIUSA // Linked Usage Reporting
	RPT_TRIG_VOLQU // Volume Quota
	RPT_TRIG_TIMQU // Time Quota
	RPT_TRIG_ENVCL // Envelope Closure
	RPT_TRIG_MACAR // MAC Addresses Reporting
	RPT_TRIG_EVETH // Event Threshold
	RPT_TRIG_EVEQU // Event Quota
	RPT_TRIG_IPMJL // IP Multicast Join/Leave
)

type ReportingTrigger struct {
	Flags uint16
}

func newRptTrig(rpgTrig ReportingTrigger) *ie.IE {
	if rpgTrig.Flags == 0 {
		return nil
	}
	return ie.NewReportingTriggers(rpgTrig.Flags)
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
