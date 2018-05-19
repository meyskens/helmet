package template

import (
	"io/ioutil"
	"path"
	"reflect"
	"testing"
)

const emptyNOTES = `1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace test -l "app=empty,release=unit" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl port-forward $POD_NAME 8080:80
`

func loadChartFromPath(p string) (*ChartInfo, map[string]interface{}, map[string][]byte) {
	chartFile, _ := ioutil.ReadFile(path.Join(p, "Chart.yaml"))
	valuesFile, _ := ioutil.ReadFile(path.Join(p, "values.yaml"))
	chart, _ := ParseChartFile(chartFile)
	values, _ := ParseValuesFile(valuesFile)

	files, _ := ioutil.ReadDir(path.Join(p, "templates"))
	filesMap := map[string][]byte{}
	for _, f := range files {
		content, _ := ioutil.ReadFile(path.Join(p, "templates", f.Name()))
		filesMap[f.Name()] = content
	}

	return chart, values, filesMap
}

func loadOutputFromPath(p string) map[string][]byte {
	files, _ := ioutil.ReadDir(p)
	filesMap := map[string][]byte{}
	for _, f := range files {
		content, _ := ioutil.ReadFile(path.Join(p, f.Name()))
		filesMap[f.Name()] = content
	}

	return filesMap
}

func TestChart_CreateManifests(t *testing.T) {
	emptyChartInfo, emptyValues, emptyFiles := loadChartFromPath("../testfiles/empty/")
	emptynnChartInfo, emptynnValues, emptynnFiles := loadChartFromPath("../testfiles/empty-nonotes/")
	emptyOut := loadOutputFromPath("../testfiles/empty-output/")
	type fields struct {
		templateFiles map[string][]byte
		values        map[string]interface{}
		chartInfo     *ChartInfo
		release       Release
	}
	type args struct {
		release Release
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string][]byte
		want1   []byte
		wantErr bool
	}{
		{
			name: "Test empty chart file",
			fields: fields{
				templateFiles: emptyFiles,
				chartInfo:     emptyChartInfo,
				values:        emptyValues,
			},
			args: args{
				release: NewRelease("unit", "test"),
			},
			want:  emptyOut,
			want1: []byte(emptyNOTES),
		},
		{
			name: "Test empty chart file without NOTES.txt",
			fields: fields{
				templateFiles: emptynnFiles,
				chartInfo:     emptynnChartInfo,
				values:        emptynnValues,
			},
			args: args{
				release: NewRelease("unit", "test"),
			},
			want:  emptyOut,
			want1: []byte{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chart{
				templateFiles: tt.fields.templateFiles,
				values:        tt.fields.values,
				chartInfo:     tt.fields.chartInfo,
				release:       tt.fields.release,
			}
			got, got1, err := c.CreateManifests(tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chart.CreateManifests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chart.CreateManifests() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Chart.CreateManifests() got1 = '%v', want '%v'", string(got1), string(tt.want1))
			}
		})
	}
}
