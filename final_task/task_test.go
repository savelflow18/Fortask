package main

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	type args struct {
		cxt context.Context
		s   string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.cxt, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_info_size(t *testing.T) {
	type args struct {
		cxt  context.Context
		file *os.File
	}
	tests := []struct {
		name string
		args args
		want int
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := info_size(tt.args.cxt, tt.args.file); got != tt.want {
				t.Errorf("info_size() = %v, want %v", got, tt.want)
			}
		})
	}
}
