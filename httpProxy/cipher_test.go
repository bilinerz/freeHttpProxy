package httpProxy

import (
	"reflect"
	"testing"
)

func TestCipher(t *testing.T) {
	k := RandKey()
	p, _ := ParseBaseKey(k)

	//创建一个解码器
	cipher := NewCipher(p)

	// 假设原数据是 [0～255]
	org := make([]byte, keyLength)
	for i := 0; i < keyLength; i++ {
		org[i] = byte(i)
	}

	// 复制一份原数据到 tmp
	tmp := make([]byte, keyLength)
	copy(tmp, org)
	t.Log("原数据", tmp) //[0～255]

	// 加密 tmp
	cipher.encode(tmp)
	t.Log("原数据替换为随机生成的密码：", tmp) //将原数据转为随机生成的 encodeKey

	// 解密 tmp
	cipher.decode(tmp)
	t.Log("将随机密码还原为原数据：", tmp) //将随机生成的密码转为原密码

	if !reflect.DeepEqual(org, tmp) {
		t.Error("解码编码数据后无法还原数据，数据不对应")
	}
}
