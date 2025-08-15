package rexCommon

import "testing"

func TestMaskEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				email: "a@example.com",
			},
			want: "*@example.com",
		},
		{
			name: "test2",
			args: args{
				email: "ab@example.com",
			},
			want: "a*@example.com",
		},
		{
			name: "test3",
			args: args{
				email: "alice@example.com",
			},
			want: "a***e@example.com",
		},
		{
			name: "test4",
			args: args{
				email: "张三@例子.公司",
			},
			want: "张*@例子.公司",
		},
		{
			name: "test5",
			args: args{
				email: "no-at-symbol",
			},
			want: "no-at-symbol",
		},
		{
			name: "test6",
			args: args{
				email: "@example.com",
			},
			want: "@example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskEmail(tt.args.email); got != tt.want {
				t.Errorf("MaskEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
