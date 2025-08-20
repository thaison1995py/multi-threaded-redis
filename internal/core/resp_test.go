package core

import (
	"reflect"
	"testing"
)

func Test_readSimpleString(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple string",
			args: args{
				data: []byte("+OK\r\n"),
			},
			want:    "OK",
			want1:   5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readSimpleString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSimpleString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSimpleString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readSimpleString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readInt64(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "int64",
			args: args{
				data: []byte(":12345\r\n"),
			},
			want:  12345,
			want1: 8,
		},
		{
			name: "negative int64",
			args: args{
				data: []byte(":-12345\r\n"),
			},
			want:  -12345,
			want1: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readInt64(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readInt64() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readInt64() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readBulkString(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "bulk string",
			args: args{
				data: []byte("$6\r\nfoobar\r\n"),
			},
			want:    "foobar",
			want1:   12,
			wantErr: false,
		},
		{
			name: "bulk string with empty value",
			args: args{
				data: []byte("$0\r\n\r\n"),
			},
			want:    "",
			want1:   6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readBulkString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBulkString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readBulkString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readBulkString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readArray(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "array of strings",
			args: args{
				data: []byte("*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n"),
			},
			want:    []interface{}{"foo", "bar", "baz"},
			want1:   31,
			wantErr: false,
		},
		{
			name: "array with mixed types",
			args: args{
				data: []byte("*3\r\n$3\r\nfoo\r\n:42\r\n$4\r\nquxx\r\n"),
			},
			want:    []interface{}{"foo", int64(42), "quxx"},
			want1:   28,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readArray(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readArray() got = %v, want %v", got, tt.want)
				return
			}

			if got1 != tt.want1 {
				t.Errorf("readArray() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
