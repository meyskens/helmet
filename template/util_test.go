package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeValues(t *testing.T) {
	type args struct {
		old map[interface{}]interface{}
		new map[interface{}]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[interface{}]interface{}
	}{
		{
			name: "test 2 new keys",
			args: args{
				old: map[interface{}]interface{}{
					"test1": "ok",
				},
				new: map[interface{}]interface{}{
					"test2": "ok",
				},
			},
			want: map[interface{}]interface{}{
				"test1": "ok",
				"test2": "ok",
			},
		},
		{
			name: "test override keys",
			args: args{
				old: map[interface{}]interface{}{
					"test1": "ok",
				},
				new: map[interface{}]interface{}{
					"test1": "evenmoreok",
				},
			},
			want: map[interface{}]interface{}{
				"test1": "evenmoreok",
			},
		},
		{
			name: "test nested keys",
			args: args{
				old: map[interface{}]interface{}{
					"test1": map[interface{}]interface{}{
						"test2": "ok",
					},
				},
				new: map[interface{}]interface{}{
					"test1": map[interface{}]interface{}{
						"test2": "whoreadsthis",
					},
				},
			},
			want: map[interface{}]interface{}{
				"test1": map[interface{}]interface{}{
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
