package rexCrypto

import (
	"crypto/elliptic"
	"crypto/x509/pkix"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestECDSAGenerateECCCertificate(t *testing.T) {
	type args struct {
		curve   elliptic.Curve
		subject pkix.Name
	}
	tests := []struct {
		name        string
		args        args
		wantCertPem []byte
		wantKeyPem  []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "TestECDSAGenerateECCCertificate",
			args: args{
				curve: elliptic.P384(),
				subject: pkix.Name{
					Country:            []string{"CN"},
					Organization:       []string{"Technology-99"},
					OrganizationalUnit: []string{"Technology-99"},
					Locality:           []string{"Beijing"},
					Province:           []string{"Beijing"},
					StreetAddress:      []string{"Chaoyang"},
					PostalCode:         []string{"100000"},
					CommonName:         "Technology-99",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCertPem, gotKeyPem, err := ECDSAGenerateECCCertificate(tt.args.curve, tt.args.subject)
			logx.Infof("gotCertPem: %s", gotCertPem)
			logx.Infof("gotKeyPem: %s", gotKeyPem)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECDSAGenerateECCCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCertPem, tt.wantCertPem) {
				t.Errorf("ECDSAGenerateECCCertificate() gotCertPem = %v, want %v", gotCertPem, tt.wantCertPem)
			}
			if !reflect.DeepEqual(gotKeyPem, tt.wantKeyPem) {
				t.Errorf("ECDSAGenerateECCCertificate() gotKeyPem = %v, want %v", gotKeyPem, tt.wantKeyPem)
			}
		})
	}
}
