package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const valuesFile = `
# Default values for empty.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1

image:
  repository: nginx
  tag: stable
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  path: /
  hosts:
    - chart-example.local
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #  cpu: 100m
  #  memory: 128Mi
  # requests:
  #  cpu: 100m
  #  memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}
`

func TestParseValuesFile(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[interface{}]interface{}
		wantErr bool
	}{
		{
			name: "Parse correct values file",
			args: args{
				in: []byte(valuesFile),
			},
			want: map[interface{}]interface{}{
				"replicaCount": 1,
				"image": map[interface{}]interface{}{
					"repository": "nginx",
					"tag":        "stable",
					"pullPolicy": "IfNotPresent",
				},
				"service": map[interface{}]interface{}{
					"type": "ClusterIP",
					"port": 80,
				},
				"ingress": map[interface{}]interface{}{
					"enabled":     false,
					"annotations": map[interface{}]interface{}{},
					"path":        "/",
					"hosts": []interface{}{
						"chart-example.local",
					},
					"tls": []interface{}{},
				},
				"resources":    map[interface{}]interface{}{},
				"nodeSelector": map[interface{}]interface{}{},
				"affinity":     map[interface{}]interface{}{},
				"tolerations":  []interface{}{},
			},
			wantErr: false,
		},
		{
			name: "test invalid yaml",
			args: args{
				in: []byte("holoworlldddd\n\r\a"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValuesFile(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseValuesFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
