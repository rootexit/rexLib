package rexShortId

import "testing"

func TestShortId(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				raw: "01KJVZV1WQ7T5THC264HC8W8KK",
			},
			want: "10KT3ElpSCRwmjKLH",
		},
		// TODO: Add test cases.
		{
			name: "case2",
			args: args{
				raw: "01KJW03BG159GSYDF5XHVBXXMM",
			},
			want: "1tNXCSkJGtwRPkO1K",
		},
		// TODO: Add test cases.
		{
			name: "case3",
			args: args{
				raw: "01KJW03D7EPBSGDPXA61S19TXY",
			},
			want: "1rgtRVOz4jzn16prD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortId(tt.args.raw); got != tt.want {
				t.Errorf("ShortId() = %v, want %v", got, tt.want)
			}
		})
	}
}
