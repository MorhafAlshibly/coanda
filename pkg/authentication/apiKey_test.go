package authentication

import (
	"encoding/hex"
	"testing"
)

func Test_compareApiKeyAndHashedApiKey_MatchingApiKeyAndHashedApiKey_True(t *testing.T) {
	byteString, err := hex.DecodeString("18964ae773c37d03d4b9f3addc6001bc")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	a := &ApiKey{
		decodedHashedApiKey: DecodedHash{
			memory:      16,
			iterations:  2,
			parallelism: 1,
			salt:        []byte("6uOvsavxISLd8B70"),
			hash:        byteString,
		},
	}
	apiKey := "password123"
	authorized, err := a.compareApiKeyAndHashedApiKey(apiKey)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if !authorized {
		t.Error("Expected true, got false")
	}
}

func Test_compareApiKeyAndHashedApiKey_NonMatchingApiKeyAndHashedApiKey_False(t *testing.T) {
	a := &ApiKey{
		decodedHashedApiKey: DecodedHash{
			memory:      16,
			iterations:  2,
			parallelism: 1,
			salt:        []byte("kpgs28Lahlg3f2"),
			hash:        []byte("Zr9108NOrwJbsadvXpqCOw"),
		},
	}
	apiKey := "password1234"
	authorized, err := a.compareApiKeyAndHashedApiKey(apiKey)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if authorized {
		t.Error("Expected false, got true")
	}
}

func Test_decodeHashedApiKey_InvalidHashedApiKeyFormat_Error(t *testing.T) {
	_, err := decodeHashedApiKey("$argon2id$v=19$m=16,t=2,p=1$a2pnczI4TGFobGczZjI")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func Test_decodeHashedApiKey_IncompatibleVersion_Error(t *testing.T) {
	_, err := decodeHashedApiKey("$argon2id$v=20$m=16,t=2,p=1$a2pnczI4TGFobGczZjI$Zr9108NOrwJbsadvXpqCOw")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func Test_decodeHashedApiKey_ValidHashedApiKey_Params(t *testing.T) {
	decodedHash, err := decodeHashedApiKey("$argon2id$v=19$m=16,t=2,p=1$a2pnczI4TGFobGczZjI$Zr9108NOrwJbsadvXpqCOw")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if decodedHash.memory != 16 {
		t.Errorf("Expected 16, got %d", decodedHash.memory)
	}
	if decodedHash.iterations != 2 {
		t.Errorf("Expected 2, got %d", decodedHash.iterations)
	}
	if decodedHash.parallelism != 1 {
		t.Errorf("Expected 1, got %d", decodedHash.parallelism)
	}
	if string(decodedHash.salt) != "kjgs28Lahlg3f2" {
		t.Errorf("Expected kpgs28Lahlg3f2, got %s", decodedHash.salt)
	}
}
