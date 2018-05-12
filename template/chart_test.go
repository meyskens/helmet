package template

import (
	"reflect"
	"testing"
)

const chartFile = `
apiVersion: v1
appVersion: "1.0"
description: A Helm chart for Kubernetes
name: empty
version: 0.1.0
`

func TestParseChartFile(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ChartInfo
		wantErr bool
	}{
		{
			name: "test normal chart file",
			args: args{
				in: []byte(chartFile),
			},
			want: &ChartInfo{
				APPVersion:  "1.0",
				APIVersion:  "v1",
				Description: "A Helm chart for Kubernetes",
				Name:        "empty",
				Version:     "0.1.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChartFile(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChartFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseChartFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
