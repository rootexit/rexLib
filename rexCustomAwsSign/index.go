package rexCustomAwsSign

import (
	"crypto/sha256"
	"encoding/hex"
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

const (
	timeFormat = "20060102T150405Z"
	//authHeaderPrefix = "AWS4-HMAC-SHA256"
	authHeaderPrefix = "REX1-HMAC-SHA256"
	shortTimeFormat  = "20060102"
	//awsV4Request     = "aws4_request"
	rexV1Request = "rex1_request"

	// emptyStringSHA256 is a SHA256 of an empty string
	EmptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
	maxSkew           = 5 * time.Minute
	doubleSpace       = "  "
)

var ignoredHeaders = map[string]string{
	rexHeaders.HeaderAuthorization: "",
	rexHeaders.HeaderUserAgent:     "",
	"X-Amzn-Trace-Id":              "",
	rexHeaders.HeaderXRequestIDFor: "",
}

func FormatDate(now time.Time) string {
	return now.Format(timeFormat)
}

func Sha256Content(bodyBytes []byte) (contentSha256 string) {
	if len(bodyBytes) != 0 {
		calculatedSha256 := hex.EncodeToString(hashSHA256(bodyBytes))
		return calculatedSha256
	} else {
		return EmptyStringSHA256
	}
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func BuildCanonicalHeaders(r *http.Request) (canonicalHeaders string, signedHeaderStr string) {
	var headers []string
	var signedHeaders []string
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
		for ignoredKey := range ignoredHeaders {
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
	stripExcessSpaces(headerItems)
	signedHeaderStr = strings.Join(headers, ";")
	canonicalHeaders = strings.Join(headerItems, "\n")
	return canonicalHeaders, signedHeaderStr
}

func stripExcessSpaces(vals []string) {
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
		j = strings.Index(str, doubleSpace)
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

func BuildCanonicalString(r *http.Request, canonicalHeaders, signedHeaders, bodyDigest string) string {
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

func BuildCredentialString(region, service string, dt time.Time) string {
	credentialString := buildSigningScope(region, service, dt)
	return credentialString
}

func formatShortTime(dt time.Time) string {
	return dt.UTC().Format(shortTimeFormat)
}

func buildSigningScope(region, service string, dt time.Time) string {
	return strings.Join([]string{
		formatShortTime(dt),
		region,
		service,
		rexV1Request,
	}, "/")
}
func BuildStringToSign(reqTime, credentialString, canonicalString string) string {
	stringToSign := strings.Join([]string{
		authHeaderPrefix,
		reqTime,
		credentialString,
		hex.EncodeToString(hashSHA256([]byte(canonicalString))),
	}, "\n")
	return stringToSign
}
