package rexUserAgent

import (
	"encoding/hex"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"github.com/ua-parser/uap-go/uaparser"
	"testing"
)

func TestUserAgentUtilsFromIpAndUa(t *testing.T) {
	type args struct {
		ip       string
		ua       string
		region   *ip2region.Ip2Region
		uaparser *uaparser.Parser
	}
	parser, _ := uaparser.New("regexes.yaml")
	region, _ := ip2region.New("ip2region.db")
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				ip:       "223.76.226.47",
				ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
				region:   region,
				uaparser: parser,
			},
			want:    Client{},
			wantErr: false,
		},
		{
			name: "case2",
			args: args{
				ip:       "223.76.226.47:7920",
				ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
				region:   region,
				uaparser: parser,
			},
			want:    Client{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserAgentUtilsFromIpAndUa(tt.args.ip, tt.args.ua, tt.args.region, tt.args.uaparser)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAgentUtilsFromIpAndUa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ip = %v", got.IP)
			t.Logf("Port = %v", got.Port)
			t.Logf("IPHash = %v", hex.EncodeToString(got.IPHash))
			t.Logf("City = %v", got.City)
			t.Logf("OsFamily = %v", got.OsFamily)
		})
	}
}
