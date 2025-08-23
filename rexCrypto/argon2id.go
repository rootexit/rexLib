package rexCrypto

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strconv"
	"strings"
)

const (
	DefaultArgon2Time      uint32 = 3
	DefaultArgon2MemoryKiB uint32 = 65536 // 64 MiB
	DefaultArgon2Threads   uint8  = 2
	DefaultArgon2KeyLen    uint32 = 32
	DefaultArgon2SaltLen          = Bits256Len // 128 bits
)

var (
	Argon2ErrEmptySecret   = errors.New("empty secret")
	Argon2ErrPhcFormatBad  = errors.New("bad phc format")
	Argon2ErrPhcParamsBad  = errors.New("bad phc params")
	Argon2ErrPhcVersionBad = errors.New("bad alg version")
)

type (
	Argon2Config struct {
		Time      uint32 // Time cose
		MemoryKiB uint32 // Memory cost
		Threads   uint8  // Parallelism degree
		KeyLen    uint32 // KeyLen
		SaltLen   BitLen // Salt
	}
	argon2Tool interface {
		HashToPHC(secret, pepper []byte) (SelfContained string, err error)
		VerifyPhc(phc string, secret, pepper []byte) (result bool, err error)
		VerifyWithPepperSet(phc string, secret []byte, pepperCurrent, pepperOld []byte) (ok bool, rehashWith []byte, err error)
		NeedsRehash(phc string) (bool, error)
		ParsePHC(phc string) (algo, ver string, m uint32, t uint32, p uint8, salt, hash []byte, err error)
	}
	defaultArgon2Tool struct {
		conf *Argon2Config
	}
)

func DefaultArgon2Config() *Argon2Config {
	return &Argon2Config{
		Time:      DefaultArgon2Time,
		MemoryKiB: DefaultArgon2MemoryKiB, // 65536 = 64 * 1024 KiB = 64 MiB
		Threads:   DefaultArgon2Threads,
		KeyLen:    DefaultArgon2KeyLen,
		SaltLen:   DefaultArgon2SaltLen,
	}
}

func NewArgon2Tool(conf *Argon2Config) argon2Tool {
	if conf == nil {
		conf = DefaultArgon2Config()
	}
	return &defaultArgon2Tool{
		conf: conf,
	}
}

// Hash 生成PHC字符串。pepper可选(服务端机密,从kms轮转)
func (d *defaultArgon2Tool) HashToPHC(secret, pepper []byte) (SelfContained string, err error) {
	if len(secret) == 0 {
		return "", Argon2ErrEmptySecret
	}
	salt := NewRand().RandBytesNoErr(d.conf.SaltLen)
	input := append(pepper, secret...) // pepper||secret
	sum := argon2.IDKey(input, salt, d.conf.Time, d.conf.MemoryKiB, d.conf.Threads, d.conf.KeyLen)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		d.conf.MemoryKiB, d.conf.Time, d.conf.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(sum)), nil
}

// Verify 校验PHC字符串
func (d *defaultArgon2Tool) VerifyPhc(phc string, secret, pepper []byte) (result bool, err error) {
	// 期望: $argon2id$v=19$m=...,t=...,p=...$<salt>$<hash>
	parts := strings.Split(phc, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, Argon2ErrPhcFormatBad
	}
	params := parts[3] // m=..,t=..,p=..
	var mem, time uint64
	var threads uint64
	for _, kv := range strings.Split(params, ",") {
		ab := strings.SplitN(kv, "=", 2)
		if len(ab) != 2 {
			return false, Argon2ErrPhcParamsBad
		}
		switch ab[0] {
		case "m":
			mem, _ = strconv.ParseUint(ab[1], 10, 32)
		case "t":
			time, _ = strconv.ParseUint(ab[1], 10, 32)
		case "p":
			threads, _ = strconv.ParseUint(ab[1], 10, 8)
		}
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, Argon2ErrPhcParamsBad
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, Argon2ErrPhcParamsBad
	}

	input := append(pepper, secret...)
	got := argon2.IDKey(input, salt, uint32(time), uint32(mem), uint8(threads), uint32(len(want)))

	// 常数时间比较
	ok := subtle.ConstantTimeCompare(got, want) == 1
	return ok, nil
}

// VerifyWithPepperSet：支持 pepper 轮转（current 优先，其次 old）。
// 通过则返回 ok=true；如果是 old 通过，rehashWith 指示用哪个 pepper 重算。
func (d *defaultArgon2Tool) VerifyWithPepperSet(phc string, secret []byte, pepperCurrent, pepperOld []byte) (ok bool, rehashWith []byte, err error) {
	ok, err = d.VerifyPhc(phc, secret, pepperCurrent)
	if err != nil {
		return false, nil, err
	}
	if ok {
		need, _ := d.NeedsRehash(phc)
		if need {
			return true, pepperCurrent, nil
		}
		return true, nil, nil
	}
	// 尝试旧 pepper
	if len(pepperOld) > 0 {
		ok2, err2 := d.VerifyPhc(phc, secret, pepperOld)
		if err2 != nil {
			return false, nil, err2
		}
		if ok2 {
			return true, pepperCurrent, nil // 用新 pepper 重算
		}
	}
	return false, nil, nil
}

// NeedsRehash：根据当前期望参数判断是否需要升级（登录成功后可触发重算）。
func (d *defaultArgon2Tool) NeedsRehash(phc string) (bool, error) {
	_, ver, m, t, p, _, want, err := d.ParsePHC(phc)
	if err != nil {
		return false, err
	}
	if ver != "19" {
		return true, nil
	}
	if m != d.conf.MemoryKiB || t != d.conf.Time || p != d.conf.Threads || uint32(len(want)) != d.conf.KeyLen {
		return true, nil
	}
	return false, nil
}

// 解析 PHC 串（最小够用版）
func (d *defaultArgon2Tool) ParsePHC(phc string) (algo, ver string, m uint32, t uint32, p uint8, salt, hash []byte, err error) {
	parts := strings.Split(phc, "$")
	// ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 {
		err = Argon2ErrPhcFormatBad
		return
	}
	algo = parts[1]
	if !strings.HasPrefix(parts[2], "v=") {
		err = Argon2ErrPhcVersionBad
		return
	}
	ver = strings.TrimPrefix(parts[2], "v=")

	var mem, tim uint64
	var thr uint64
	for _, kv := range strings.Split(parts[3], ",") {
		ab := strings.SplitN(kv, "=", 2)
		if len(ab) != 2 {
			err = Argon2ErrPhcParamsBad
			return
		}
		switch ab[0] {
		case "m":
			mem, _ = strconv.ParseUint(ab[1], 10, 32)
		case "t":
			tim, _ = strconv.ParseUint(ab[1], 10, 32)
		case "p":
			thr, _ = strconv.ParseUint(ab[1], 10, 8)
		}
	}
	m = uint32(mem)
	t = uint32(tim)
	p = uint8(thr)

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return
	}
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	return
}
