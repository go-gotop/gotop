package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// Config 用于设定加密参数
type Config struct {
	// Passphrase 用户提供的主密码，用于派生加密密钥
	Passphrase string
	// Salt 固定或事先生成的随机 salt，用于派生 key，长度建议固定（例如 16 字节）
	Salt []byte
	// Iteration PBKDF2 迭代次数，如 100000
	Iteration int
	// KeyLen AES256 的 key 长度为 32 字节
	KeyLen int
}

// DeriveKey 从 Passphrase 和 Salt 中派生出 AES 密钥
func DeriveKey(cfg *Config) ([]byte, error) {
	if cfg.Iteration <= 0 {
		cfg.Iteration = 100000
	}
	if cfg.KeyLen <= 0 {
		cfg.KeyLen = 32 // AES-256
	}
	if len(cfg.Salt) == 0 {
		return nil, errors.New("salt is required to derive key")
	}
	return pbkdf2.Key([]byte(cfg.Passphrase), cfg.Salt, cfg.Iteration, cfg.KeyLen, sha256.New), nil
}

// Encrypt 使用 AES-GCM 对明文进行加密
// 返回值: (nonce | cipherText), nonce 与 cipherText 拼接在一起的字节数组
func Encrypt(cfg *Config, plaintext []byte) ([]byte, error) {
	key, err := DeriveKey(cfg)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, plaintext, nil)
	// 将 nonce 和 cipherText 拼接返回：nonce + cipherText
	result := append(nonce, cipherText...)

	return result, nil
}

// Decrypt 使用 AES-GCM 对密文进行解密
// 参数 data 格式为: nonce + cipherText
func Decrypt(cfg *Config, data []byte) ([]byte, error) {
	key, err := DeriveKey(cfg)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}