package template

import (
	"reflect"
	"testing"
)

func TestNewRelease(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
		want Release
	}{
		{
			name: "test release",
			args: args{
				name:      "test",
				namespace: "default",
			},
			want: Release{
				Name:      "test",
				Service:   "helmet",
				Namespace: "default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRelease(tt.args.name, tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}
