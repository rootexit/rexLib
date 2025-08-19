package rexCustomAwsSign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/rootexit/rexLib/rexHeaders"
	"net/http"
	"sort"
	"strings"
	"time"
)

/*
	note: 这个签名部分主要是借鉴aws的V4签名，相关文档可以访问，我愿称之为《伟大的杰作》
	https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
*/

type (
	CustomSigner interface {
		WithMaxSkew(maxSkew time.Duration)
		WithIgnoredHeaders(IgnoredHeaders map[string]string)
		WithNeedSignHeaders(NeedSignHeaders map[string]string)
		GetDeriveKeyPrefix() string
		GetTimeFormat() string
		GetAuthHeaderPrefix() string
		GetShortTimeFormat() string
		GetVersionRequest() string
		GetEmptyStringSHA256() string
		GetMaxSkew() time.Duration
		GetDoubleSpace() string
		GetAuthHeaderSignatureElem() string
		GetHeaderDate() string
		GetHeaderContentSha256() string
		GetIgnoredHeaders() map[string]string
		GetNeedSignHeaders() map[string]string
		SignAuth(accessKeyID, credentialString, signedHeaders, signature string) string
		BuildSignature(Region, ServiceName, SecretAccessKey, stringToSign string, Time time.Time) string
		DeriveSigningKey(region, service, secretKey string, dt time.Time) []byte
		HmacSHA256(key []byte, data []byte) []byte
		FormatDate(now time.Time) string
		Sha256Content(bodyBytes []byte) (contentSha256 string)
		HashSHA256(data []byte) []byte
		BuildCanonicalHeaders(r *http.Request) (canonicalHeaders string, signedHeaderStr string)
		StripExcessSpaces(vals []string)
		BuildCanonicalString(r *http.Request, canonicalHeaders, signedHeaders, bodyDigest string) string
		BuildCredentialString(region, service string, dt time.Time) string
		FormatShortTime(dt time.Time) string
		BuildSigningScope(region, service string, dt time.Time) string
		BuildStringToSign(reqTime, credentialString, canonicalString string) string
	}
	customSigner struct {
		DeriveKeyPrefix         string
		TimeFormat              string
		AuthHeaderPrefix        string
		ShortTimeFormat         string
		VersionRequest          string
		EmptyStringSHA256       string
		MaxSkew                 time.Duration
		DoubleSpace             string
		AuthHeaderSignatureElem string
		HeaderDate              string
		HeaderContentSha256     string
		IgnoredHeaders          map[string]string
		NeedSignHeaders         map[string]string
	}
)

func NewCustomSigner(shortName string, version uint) CustomSigner {
	if shortName == "AWS" {
		return &customSigner{
			DeriveKeyPrefix:         "AWS4",
			TimeFormat:              "20060102T150405Z",
			AuthHeaderPrefix:        "AWS4-HMAC-SHA256",
			ShortTimeFormat:         "20060102",
			VersionRequest:          "aws4_request",
			EmptyStringSHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			MaxSkew:                 5 * time.Minute,
			DoubleSpace:             "  ",
			AuthHeaderSignatureElem: "Signature=",
			HeaderDate:              "X-Amz-Date",
			HeaderContentSha256:     "X-Amz-Content-Sha256",
			IgnoredHeaders: map[string]string{
				rexHeaders.HeaderAuthorization:  "",
				rexHeaders.HeaderUserAgent:      "",
				"X-Amzn-Trace-Id":               "",
				rexHeaders.HeaderAcceptEncoding: "",
				rexHeaders.HeaderConnection:     "",
				rexHeaders.HeaderContentLength:  "",
				rexHeaders.HeaderAccept:         "",
			},
			NeedSignHeaders: map[string]string{
				"X-Amz-Date":           "",
				"X-Amz-Content-Sha256": "",
			},
		}
	} else {
		tmpHeaderDate := fmt.Sprintf("X-%s-Date", shortName)
		tmpHeaderContentSha256 := fmt.Sprintf("X-%s-Content-Sha256", shortName)
		return &customSigner{
			DeriveKeyPrefix:         strings.ToUpper(fmt.Sprintf("%s%d", shortName, version)),
			TimeFormat:              "20060102T150405Z",
			AuthHeaderPrefix:        fmt.Sprintf("%s%d-HMAC-SHA25", shortName, version),
			ShortTimeFormat:         "20060102",
			VersionRequest:          strings.ToLower(fmt.Sprintf("%s%d_request", shortName, version)),
			EmptyStringSHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			MaxSkew:                 5 * time.Minute,
			DoubleSpace:             "  ",
			AuthHeaderSignatureElem: "Signature=",
			HeaderDate:              tmpHeaderDate,
			HeaderContentSha256:     tmpHeaderContentSha256,
			IgnoredHeaders: map[string]string{
				rexHeaders.HeaderAuthorization:  "",
				rexHeaders.HeaderUserAgent:      "",
				"X-Amzn-Trace-Id":               "",
				rexHeaders.HeaderAcceptEncoding: "",
				rexHeaders.HeaderConnection:     "",
				rexHeaders.HeaderContentLength:  "",
				rexHeaders.HeaderAccept:         "",
			},
			NeedSignHeaders: map[string]string{
				rexHeaders.HeaderXRequestIDFor: "",
				rexHeaders.HeaderContentType:   "",
				tmpHeaderDate:                  "",
				tmpHeaderDate:                  "",
			},
		}
	}
}

func (s *customSigner) WithMaxSkew(maxSkew time.Duration) {
	s.MaxSkew = maxSkew
}

func (s *customSigner) WithIgnoredHeaders(IgnoredHeaders map[string]string) {
	s.IgnoredHeaders = IgnoredHeaders
}

func (s *customSigner) WithNeedSignHeaders(NeedSignHeaders map[string]string) {
	s.NeedSignHeaders = NeedSignHeaders
}

func (s *customSigner) GetDeriveKeyPrefix() string {
	return s.DeriveKeyPrefix
}

func (s *customSigner) GetTimeFormat() string {
	return s.TimeFormat
}

func (s *customSigner) GetAuthHeaderPrefix() string {
	return s.AuthHeaderPrefix
}

func (s *customSigner) GetShortTimeFormat() string {
	return s.ShortTimeFormat
}

func (s *customSigner) GetVersionRequest() string {
	return s.VersionRequest
}

func (s *customSigner) GetEmptyStringSHA256() string {
	return s.EmptyStringSHA256
}

func (s *customSigner) GetMaxSkew() time.Duration {
	return s.MaxSkew
}

func (s *customSigner) GetDoubleSpace() string {
	return s.DoubleSpace
}

func (s *customSigner) GetAuthHeaderSignatureElem() string {
	return s.AuthHeaderSignatureElem
}

func (s *customSigner) GetHeaderDate() string {
	return s.HeaderDate
}

func (s *customSigner) GetHeaderContentSha256() string {
	return s.HeaderContentSha256
}

func (s *customSigner) GetIgnoredHeaders() map[string]string {
	return s.IgnoredHeaders
}

func (s *customSigner) GetNeedSignHeaders() map[string]string {
	return s.NeedSignHeaders
}

func (s *customSigner) SignAuth(accessKeyID, credentialString, signedHeaders, signature string) string {
	parts := []string{
		s.AuthHeaderPrefix + " Credential=" + accessKeyID + "/" + credentialString,
		"SignedHeaders=" + signedHeaders,
		s.AuthHeaderSignatureElem + signature,
	}
	return strings.Join(parts, ", ")
}

func (s *customSigner) BuildSignature(Region, ServiceName, SecretAccessKey, stringToSign string, Time time.Time) string {
	creds := s.DeriveSigningKey(Region, ServiceName, SecretAccessKey, Time)
	signature := s.HmacSHA256(creds, []byte(stringToSign))
	signatureHex := hex.EncodeToString(signature)
	return signatureHex
}

func (s *customSigner) DeriveSigningKey(region, service, secretKey string, dt time.Time) []byte {
	kDate := s.HmacSHA256([]byte(s.DeriveKeyPrefix+secretKey), []byte(s.FormatShortTime(dt)))
	kRegion := s.HmacSHA256(kDate, []byte(region))
	kService := s.HmacSHA256(kRegion, []byte(service))
	signingKey := s.HmacSHA256(kService, []byte(s.VersionRequest))
	return signingKey
}

func (s *customSigner) HmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func (s *customSigner) FormatDate(now time.Time) string {
	return now.UTC().Format(s.TimeFormat)
}

func (s *customSigner) Sha256Content(bodyBytes []byte) (contentSha256 string) {
	if len(bodyBytes) != 0 {
		calculatedSha256 := hex.EncodeToString(s.HashSHA256(bodyBytes))
		return calculatedSha256
	} else {
		return s.EmptyStringSHA256
	}
}

func (s *customSigner) HashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (s *customSigner) BuildCanonicalHeaders(r *http.Request) (canonicalHeaders string, signedHeaderStr string) {
	var headers []string
	var signedHeaders []string
	for k := range s.NeedSignHeaders {
		signedHeaders = append(signedHeaders, strings.ToLower(k))
	}
	//for k := range r.Header {
	//	signedHeaders = append(signedHeaders, strings.ToLower(k))
	//}
	headers = append(headers, strings.ToLower(rexHeaders.HeaderHost))
	signedHeaderVals := make(http.Header)
	for k, v := range r.Header {
		lowerCaseKey := strings.ToLower(k)
		isHasSignedHeader := false
		for _, header := range signedHeaders {
			if header == lowerCaseKey {
				isHasSignedHeader = true
				break
			}
		}
		if !isHasSignedHeader {
			continue
		}

		isHas := false
		for ignoredKey := range s.IgnoredHeaders {
			if strings.HasPrefix(k, ignoredKey) {
				isHas = true
				break
			}
		}
		if isHas {
			continue // ignored header
		}

		if _, ok := signedHeaderVals[lowerCaseKey]; ok {
			// include additional values
			signedHeaderVals[lowerCaseKey] = append(signedHeaderVals[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		signedHeaderVals[lowerCaseKey] = v
	}
	sort.Strings(headers)

	headerItems := make([]string, len(headers))
	for i, k := range headers {
		if k == "host" {
			if r.Host != "" {
				headerItems[i] = "host:" + r.Host
			} else {
				headerItems[i] = "host:" + r.URL.Host
			}
		} else {
			headerValues := make([]string, len(signedHeaderVals[k]))
			for i, v := range signedHeaderVals[k] {
				headerValues[i] = strings.TrimSpace(v)
			}
			headerItems[i] = k + ":" +
				strings.Join(headerValues, ",")
		}
	}
	s.StripExcessSpaces(headerItems)
	signedHeaderStr = strings.Join(headers, ";")
	canonicalHeaders = strings.Join(headerItems, "\n")
	return canonicalHeaders, signedHeaderStr
}

func (s *customSigner) StripExcessSpaces(vals []string) {
	var j, k, l, m, spaces int
	for i, str := range vals {
		// Trim trailing spaces
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- {
		}

		// Trim leading spaces
		for k = 0; k < j && str[k] == ' '; k++ {
		}
		str = str[k : j+1]

		// Strip multiple spaces.
		j = strings.Index(str, s.DoubleSpace)
		if j < 0 {
			vals[i] = str
			continue
		}

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ {
			if buf[k] == ' ' {
				if spaces == 0 {
					// First space.
					buf[m] = buf[k]
					m++
				}
				spaces++
			} else {
				// End of multiple spaces.
				spaces = 0
				buf[m] = buf[k]
				m++
			}
		}

		vals[i] = string(buf[:m])
	}
}

func (s *customSigner) BuildCanonicalString(r *http.Request, canonicalHeaders, signedHeaders, bodyDigest string) string {
	uri := r.URL.Path

	canonicalString := strings.Join([]string{
		r.Method,
		uri,
		r.URL.RawQuery,
		canonicalHeaders + "\n",
		signedHeaders,
		bodyDigest,
	}, "\n")
	return canonicalString
}

func (s *customSigner) BuildCredentialString(region, service string, dt time.Time) string {
	credentialString := s.BuildSigningScope(region, service, dt)
	return credentialString
}

func (s *customSigner) FormatShortTime(dt time.Time) string {
	return dt.UTC().Format(s.ShortTimeFormat)
}

func (s *customSigner) BuildSigningScope(region, service string, dt time.Time) string {
	return strings.Join([]string{
		s.FormatShortTime(dt),
		region,
		service,
		s.VersionRequest,
	}, "/")
}
func (s *customSigner) BuildStringToSign(reqTime, credentialString, canonicalString string) string {
	stringToSign := strings.Join([]string{
		s.AuthHeaderPrefix,
		reqTime,
		credentialString,
		hex.EncodeToString(s.HashSHA256([]byte(canonicalString))),
	}, "\n")
	return stringToSign
}
