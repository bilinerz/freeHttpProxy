package httpProxy

import "testing"

func TestRandKey(t *testing.T) {
	k := RandKey()
	t.Log("随机密钥：", k)
}

func TestParseBaseKey(t *testing.T) {
	k := RandKey()

	t.Log("base64加密后的key：", k)

	key, err := ParseBaseKey(k)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("base64解密后的Key：", key)
}
