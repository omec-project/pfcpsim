// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-present Ian Chen <ychen.cs10@nycu.edu.tw>

package session

import (
	"log"
	"time"

	"github.com/wmnsk/go-pfcp/ie"
)

type measurementMethodParams struct {
	event int
	volum int
	durat int
}

type urrBuilder struct {
	urrID             uint32
	method            IEMethod
	measurementMethod *measurementMethodParams
	measurementPeriod time.Duration
	measurementInfo   uint8
	rptTrig           ReportingTrigger
	volumThreshold    *volumThreshold
	volumeQuota       *volumeQuota
}

const UrrNoFuzz = 0
const UrrWithMeasurementInfo = 1
const UrrWithMeasurementMethod = 2
const UrrWithMeasurementPeriod = 3
const UrrMax = 4

// NewURRBuilder returns a urrBuilder.
func NewURRBuilder() *urrBuilder {
	return &urrBuilder{}
}

func (b *urrBuilder) FuzzIE(ieType int, arg uint) *urrBuilder {
	switch ieType {
	case UrrWithMeasurementInfo:
		log.Println("Fuzz: UrrWithMeasurementInfo")
		return b.WithMeasurementInfo(uint8(arg))
	case UrrWithMeasurementMethod:
		log.Println("Fuzz: UrrWithMeasurementMethod")

		args := []int{0, 1, 0}
		i := arg % 3
		args[i] = int(arg)

		return b.WithMeasurementMethod(args[0], args[1], args[2])
	case UrrWithMeasurementPeriod:
		log.Println("Fuzz: UrrWithMeasurementPeriod")
		return b.WithMeasurementPeriod(time.Duration(arg))
	default:
	}

	return b
}

func (b *urrBuilder) WithID(id uint32) *urrBuilder {
	b.urrID = id
	return b
}

func (b *urrBuilder) WithMethod(method IEMethod) *urrBuilder {
	b.method = method
	return b
}

func (b *urrBuilder) WithMeasurementMethod(event, volum, durat int) *urrBuilder {
	b.measurementMethod = &measurementMethodParams{
		event: event,
		volum: volum,
		durat: durat,
	}

	return b
}

func (b *urrBuilder) WithMeasurementPeriod(period time.Duration) *urrBuilder {
	b.measurementPeriod = period
	return b
}

func (b *urrBuilder) WithMeasurementInfo(info uint8) *urrBuilder {
	b.measurementInfo = info
	return b
}

func (b *urrBuilder) WithVolumeThreshold(flags uint8, tvol, uvol, dvol uint64) *urrBuilder {
	b.volumThreshold = &volumThreshold{
		flags: flags,
		tvol:  tvol,
		uvol:  uvol,
		dvol:  dvol,
	}

	return b
}

func (b *urrBuilder) WithVolumeQuota(flags uint8, tvol, uvol, dvol uint64) *urrBuilder {
	b.volumeQuota = &volumeQuota{
		flags: flags,
		tvol:  tvol,
		uvol:  uvol,
		dvol:  dvol,
	}

	return b
}

// TS 29.244 5.2.2.2
// When provisioning a URR, the CP function shall provide the reporting trigger(s) in the Reporting Triggers IE of the
// URR which shall cause the UP function to generate and send a Usage Report for this URR to the CP function.
func (b *urrBuilder) WithReportingTrigger(rptTrig ReportingTrigger) *urrBuilder {
	b.rptTrig = rptTrig
	return b
}

func (b *urrBuilder) validate() {
	if b.urrID == 0 {
		panic("URR ID is not set")
	}

	if b.measurementInfo > 0 && b.measurementInfo > MNOP {
		panic("Measurement Information is not valid")
	}
}

func (b *urrBuilder) Build() *ie.IE {
	if doCheck {
		b.validate()
	}

	createFunc := ie.NewCreateURR
	if b.method == Update {
		createFunc = ie.NewUpdateURR
	}

	urr := createFunc(ie.NewURRID(b.urrID),
		newMeasurementMethod(b.measurementMethod),
		ie.NewMeasurementPeriod(b.measurementPeriod),
		NewRptTrig(b.rptTrig),
		newVolumeThreshold(b.volumThreshold),
		newVolumeQuota(b.volumeQuota),
		newMeasurementInfo(b.measurementInfo),
	)

	if b.method == Delete {
		return ie.NewRemoveURR(urr)
	}

	return urr
}
