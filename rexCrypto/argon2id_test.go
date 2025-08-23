package rexCrypto

import (
	"crypto/subtle"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

var (
	secret        = NewRand().RandBytesNoErr(Bits256Len)
	pepperCurrent = NewRand().RandBytesNoErr(Bits256Len)
	pepperOld     = NewRand().RandBytesNoErr(Bits256Len)
	configCurrent = DefaultArgon2Config()
	configOld     = Argon2Config{
		Time:      2,
		MemoryKiB: 64, // 65536 = 64 * 1024 KiB = 64 MiB
		Threads:   2,
		KeyLen:    32,
		SaltLen:   32,
	}
)

func TestArgon2id(t *testing.T) {
	// note: 先用oldConfig+pepperOld生成PHC字符串
	argonOld := NewArgon2Tool(&configOld)
	argonNew := NewArgon2Tool(configCurrent)

	phcOldPepper, err := argonOld.HashToPHC(secret, pepperOld)
	if err != nil {
		t.Fatalf("HashToPHC with old config failed: %v", err)
	}
	logx.Infof("HashToPHC with old config: %v", phcOldPepper)

	// 验证PHC字符串, 先测试旧配置下，验证通过，同时提示需要用新pepper进行重新hash
	verified, rehashWith, err := argonOld.VerifyWithPepperSet(phcOldPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with old config failed")
	}
	if rehashWith == nil {
		// note: 如果不为nil, 说明需要重新哈希，这里应该提示需要用新的pepper做hash
		// note: 同时还有可能一种情况，就是2次哈希的配置是不同的，需要升级
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// 验证PHC字符串, 先测试新配置下，验证通过，同时提示需要用新配置进行重新hash
	verified, rehashWith, err = argonNew.VerifyWithPepperSet(phcOldPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet with new config failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with new config failed")
	}
	if rehashWith == nil {
		// rehashWith 应该为新的pepperCurrent
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// note: 再用新的pepper和旧配置生成PHC字符串
	phcNewPepper, err := argonOld.HashToPHC(secret, pepperCurrent)
	if err != nil {
		t.Fatalf("HashToPHC with old config failed: %v", err)
	}
	logx.Infof("HashToPHC with old config: %v", phcNewPepper)

	// 验证PHC字符串, 先测试旧配置下，验证通过，同时提示需要用新pepper进行重新hash
	verified, rehashWith, err = argonOld.VerifyWithPepperSet(phcNewPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with old config failed")
	}
	if rehashWith == nil {
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// 验证PHC字符串, 先测试新配置下，验证通过，同时提示需要用新配置进行重新hash
	verified, rehashWith, err = argonNew.VerifyWithPepperSet(phcNewPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet with new config failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with new config failed")
	}
	if rehashWith == nil {
		// rehashWith 应该为新的pepperCurrent
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// note: 再用旧的pepper和新配置生成PHC字符串
	phcNewConfigOldPepper, err := argonNew.HashToPHC(secret, pepperOld)
	if err != nil {
		t.Fatalf("HashToPHC with old config failed: %v", err)
	}
	logx.Infof("HashToPHC with new config: %v", phcNewConfigOldPepper)

	// 验证PHC字符串, 先测试旧配置下，验证通过，同时提示需要用新pepper进行重新hash
	verified, rehashWith, err = argonOld.VerifyWithPepperSet(phcNewConfigOldPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with old config failed")
	}
	if rehashWith == nil {
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// 验证PHC字符串, 先测试新配置下，验证通过，同时提示需要用新配置进行重新hash
	verified, rehashWith, err = argonNew.VerifyWithPepperSet(phcNewConfigOldPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet with new config failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with new config failed")
	}
	if rehashWith == nil {
		// rehashWith 应该为新的pepperCurrent
		if BytesEqual(rehashWith, pepperCurrent) {
			t.Fatalf("Expected rehashWith to be not nil, got: %v", rehashWith)
		}
	}

	// note: 再用新的pepper和新配置生成PHC字符串
	phcNewConfigNewPepper, err := argonNew.HashToPHC(secret, pepperCurrent)
	if err != nil {
		t.Fatalf("HashToPHC with old config failed: %v", err)
	}
	logx.Infof("HashToPHC with new config: %v", phcNewConfigNewPepper)

	// 验证PHC字符串, 先测试旧配置下，验证通过，同时提示需要用新pepper进行重新hash, 实际上没有这种情况可能发生
	//verified, rehashWith, err = argonOld.VerifyWithPepperSet(phcNewConfigNewPepper, secret, pepperCurrent, pepperOld)
	//if err != nil {
	//	t.Fatalf("VerifyWithPepperSet failed: %v", err)
	//}
	//if !verified {
	//	t.Fatalf("Verification with old config failed")
	//}
	//if len(rehashWith) == 0 {
	//	t.Fatalf("Expected rehashWith to be be nil, got: %v", rehashWith)
	//}

	// 验证PHC字符串, 先测试新配置下，验证通过，同时提示需要用新配置进行重新hash
	verified, rehashWith, err = argonNew.VerifyWithPepperSet(phcNewConfigNewPepper, secret, pepperCurrent, pepperOld)
	if err != nil {
		t.Fatalf("VerifyWithPepperSet with new config failed: %v", err)
	}
	if !verified {
		t.Fatalf("Verification with new config failed")
	}
	if rehashWith != nil {
		t.Fatalf("Expected rehashWith to be be nil, got: %v", rehashWith)
	}

	//argonNew := NewArgon2Tool(&configCurrent)
}

func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}
