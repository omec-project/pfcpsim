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
		{name: "incorrect app filter",
			args: &args{
				filterString: "test:10.0.0.0/8:80-80:allow",
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
