package rpc

import (
	"reflect"
	"testing"
)

func Test_suitableMethods(t *testing.T) {
	type args struct {
		typ       reflect.Type
		reportErr bool
	}
	tests := []struct {
		name string
		args args
		want map[string]*methodType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := suitableMethods(tt.args.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("suitableMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExportedOrBuiltinType(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExportedOrBuiltinType(tt.args.t); got != tt.want {
				t.Errorf("isExportedOrBuiltinType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_call(t *testing.T) {
	type fields struct {
		name   string
		rcvr   reflect.Value
		typ    reflect.Type
		method map[string]*methodType
	}
	type args struct {
		mtype  *methodType
		argv   reflect.Value
		replyv reflect.Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				name:   tt.fields.name,
				rcvr:   tt.fields.rcvr,
				typ:    tt.fields.typ,
				method: tt.fields.method,
			}
			if err := s.call(tt.args.mtype, tt.args.argv, tt.args.replyv); (err != nil) != tt.wantErr {
				t.Errorf("service.call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
