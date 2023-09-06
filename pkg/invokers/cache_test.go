package invokers

import (
	"testing"
)

func TestCacheInvokerGenerateKey(t *testing.T) {
	key, err := generateKey(&MockCommand{
		ExecuteFunc: nil,
		MarshalJSONFunc: func() ([]byte, error) {
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
	})
	expected := "*invokers.MockCommand: {\"ExecuteFunc\":null}"
	if err != nil {
		t.Error(err)
	}
	if key != expected {
		t.Errorf("Expected key to be %s', got '%s'", expected, key)
	}
}
