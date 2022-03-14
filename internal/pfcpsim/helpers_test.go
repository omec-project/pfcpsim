// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wmnsk/go-pfcp/ie"
)

func Test_parseAppFilter(t *testing.T) {
	type args struct {
		filterString string
	}

	type want struct {
		SDFFilter  string
		gateStatus uint8
	}

	tests := []struct {
		name    string
		args    *args
		want    *want
		wantErr bool
	}{
		{name: "Correct app filter",
			args: &args{
				filterString: "udp:10.0.0.0/8:80-80:allow",
			},
			want: &want{
				SDFFilter:  "permit out udp from 10.0.0.0/8 to assigned 80-80",
				gateStatus: ie.GateStatusOpen,
			},
		},
		{name: "Correct app filter with deny",
			args: &args{
				filterString: "udp:10.0.0.0/8:80-80:deny",
			},
			want: &want{
				SDFFilter:  "permit out udp from 10.0.0.0/8 to assigned 80-80",
				gateStatus: ie.GateStatusClosed,
			},
		},
		{name: "Correct app filter with deny-all policy",
			args: &args{
				filterString: "ip:0.0.0.0/0:any:deny",
			},
			want: &want{
				SDFFilter:  "permit out ip from 0.0.0.0/0 to assigned",
				gateStatus: ie.GateStatusClosed,
			},
		},
		{name: "Correct app filter with deny-all policy 2",
			args: &args{
				filterString: "ip:any:any:deny",
			},
			want: &want{
				SDFFilter:  "permit out ip from any to assigned",
				gateStatus: ie.GateStatusClosed,
			},
		},
		{name: "Correct app filter with allow-all policy",
			args: &args{
				filterString: "ip:any:any:allow",
			},
			want: &want{
				SDFFilter:  "permit out ip from any to assigned",
				gateStatus: ie.GateStatusOpen,
			},
		},
		{name: "Correct app filter with allow-all policy 2",
			args: &args{
				filterString: "ip:0.0.0.0/0:any:allow",
			},
			want: &want{
				SDFFilter:  "permit out ip from 0.0.0.0/0 to assigned",
				gateStatus: ie.GateStatusOpen,
			},
		},
		{name: "incorrect app filter bad protocol",
			args: &args{
				filterString: "test:10.0.0.0/8:80-80:allow",
			},
			want:    &want{},
			wantErr: true,
		},
		{name: "incorrect app filter bad IP format",
			args: &args{
				filterString: "ip:10/8:80-80:allow",
			},
			want:    &want{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				filter, gateStatus, err := parseAppFilter(tt.args.filterString)
				if tt.wantErr {
					require.Error(t, err)
					return
				}

				require.Equal(t, tt.want.SDFFilter, filter)
				require.Equal(t, tt.want.gateStatus, gateStatus)
			},
		)
	}
}
