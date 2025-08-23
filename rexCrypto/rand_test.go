package rexCrypto

import (
	"testing"
)

func Test_defaultRand_RandLowerString(t *testing.T) {
	type args struct {
		stringLen int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{1}, "1"},
		{"test1", args{2}, "12"},
		{"test1", args{3}, "123"},
		{"test1", args{4}, "1234"},
		{"test1", args{5}, "12345"},
		{"test1", args{6}, "123456"},
		{"test1", args{7}, "1234567"},
		{"test1", args{8}, "12345678"},
		{"test1", args{9}, "123456789"},
		{"test1", args{10}, "1234567890"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defaultRand{}
			if got := de.RandLowerString(tt.args.stringLen); len(got) != len(tt.want) {
				t.Errorf("RandLowerString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultRand_RandBytes(t *testing.T) {
	type args struct {
		btLen BitLen
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test16", args{Bits16Len}, make([]byte, Bits16Len), false},
		{"test32", args{Bits32Len}, make([]byte, Bits32Len), false},
		{"test64", args{Bits64Len}, make([]byte, Bits64Len), false},
		{"test128", args{Bits128Len}, make([]byte, Bits128Len), false},
		{"test256", args{Bits256Len}, make([]byte, Bits256Len), false},
		{"test512", args{Bits512Len}, make([]byte, Bits512Len), false},
		{"test1024", args{Bits1024Len}, make([]byte, Bits1024Len), false},
		{"test2048", args{Bits2048Len}, make([]byte, Bits2048Len), false},
		{"test4096", args{Bits4096Len}, make([]byte, Bits4096Len), false},
		{"test8192", args{Bits8192Len}, make([]byte, Bits8192Len), false},
		{"test16384", args{Bits16384Len}, make([]byte, Bits16384Len), false},
		{"testAny32", args{NewRand().GetAnyBtLen(32)}, make([]byte, Bits32Len), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &defaultRand{}
			got, err := r.RandBytes(tt.args.btLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("RandBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
