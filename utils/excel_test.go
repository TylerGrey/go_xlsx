package utils

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"testing"
)

func TestExcelize(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *excelize.File
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Excelize(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Excelize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Excelize() got = %v, want %v", got, tt.want)
			}
		})
	}
}