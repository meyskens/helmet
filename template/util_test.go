package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeValues(t *testing.T) {
	type args struct {
		old map[string]interface{}
		new map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "test 2 new keys",
			args: args{
				old: map[string]interface{}{
					"test1": "ok",
				},
				new: map[string]interface{}{
					"test2": "ok",
				},
			},
			want: map[string]interface{}{
				"test1": "ok",
				"test2": "ok",
			},
		},
		{
			name: "test override keys",
			args: args{
				old: map[string]interface{}{
					"test1": "ok",
				},
				new: map[string]interface{}{
					"test1": "evenmoreok",
				},
			},
			want: map[string]interface{}{
				"test1": "evenmoreok",
			},
		},
		{
			name: "test nested keys",
			args: args{
				old: map[string]interface{}{
					"test1": map[string]interface{}{
						"test2": "ok",
					},
				},
				new: map[string]interface{}{
					"test1": map[string]interface{}{
						"test2": "whoreadsthis",
					},
				},
			},
			want: map[string]interface{}{
				"test1": map[string]interface{}{
					"test2": "whoreadsthis",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeValues(tt.args.old, tt.args.new)
			assert.Equal(t, tt.args.old, tt.want)
		})
	}
}
