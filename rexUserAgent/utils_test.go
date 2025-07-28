package rexUserAgent

import (
	"context"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestUserAgentUtils(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserAgentUtils(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserAgentUtils() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAgentUtilsWithFunc(t *testing.T) {
	type args struct {
		ctx      context.Context
		callback func(client Client)
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		// TODO: Add test cases.
		{
			name: "aa",
			args: args{
				ctx: context.WithValue(context.Background(), rexCtx.CtxClientIp{}, "223.76.226.47"),
				callback: func(client Client) {
					logx.Infof("打印中间过程的ip: %s", client.IPHash)
				},
			},
			want: Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserAgentUtilsWithFunc(tt.args.ctx, tt.args.callback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserAgentUtilsWithFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
