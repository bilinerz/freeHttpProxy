package httpProxy

type Cipher struct {
	// 加密密钥
	encodeKey *Key
	// 解密密钥
	decodeKey *Key
}

// 加密原数据
func (cipher *Cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodeKey[v]
	}
}

// 解码加密后的数据到原数据
func (cipher *Cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodeKey[v]
	}
}

// NewCipher 新建一个编码解码器
func NewCipher(encodeKey *Key) *Cipher {
	decodeKey := &Key{}
	for i, v := range encodeKey {
		encodeKey[i] = v
		decodeKey[v] = byte(i)
	}
	return &Cipher{
		encodeKey: encodeKey,
		decodeKey: decodeKey,
	}
}
