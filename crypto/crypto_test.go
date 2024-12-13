package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantErr   bool
		wantKeyLen int
	}{
		{
			name: "valid_config",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				Iteration: 100000,
				KeyLen:    32,
			},
			wantErr:    false,
			wantKeyLen: 32,
		},
		{
			name: "empty_salt",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte{},
				Iteration: 100000,
				KeyLen:    32,
			},
			wantErr: true,
		},
		{
			name: "default_iteration",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				KeyLen:    32,
			},
			wantErr:    false,
			wantKeyLen: 32,
		},
		{
			name: "default_keylen",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				Iteration: 100000,
			},
			wantErr:    false,
			wantKeyLen: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := DeriveKey(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantKeyLen, len(key))
			
			// 验证相同配置生成相同密钥
			key2, err := DeriveKey(tt.cfg)
			require.NoError(t, err)
			assert.Equal(t, key, key2)
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		plaintext []byte
		wantErr   bool
	}{
		{
			name: "normal_text",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				Iteration: 100000,
				KeyLen:    32,
			},
			plaintext: []byte("Hello, World!"),
			wantErr:   false,
		},
		{
			name: "empty_text",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				Iteration: 100000,
				KeyLen:    32,
			},
			plaintext: []byte{},
			wantErr:   false,
		},
		{
			name: "long_text",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte("testsalt12345678"),
				Iteration: 100000,
				KeyLen:    32,
			},
			plaintext: bytes.Repeat([]byte("A"), 1000),
			wantErr:   false,
		},
		{
			name: "invalid_config",
			cfg: &Config{
				Passphrase: "test123",
				Salt:      []byte{},
				Iteration: 100000,
				KeyLen:    32,
			},
			plaintext: []byte("Hello, World!"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			ciphertext, err := Encrypt(tt.cfg, tt.plaintext)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, ciphertext)

			// 解密
			decrypted, err := Decrypt(tt.cfg, ciphertext)
			require.NoError(t, err)
			
			if !bytes.Equal(tt.plaintext, decrypted) {
				t.Errorf("Decrypt() got = %v, want %v", decrypted, tt.plaintext)
			}

			// 验证不同的明文产生不同的密文
			if len(tt.plaintext) > 0 {
				ciphertext2, err := Encrypt(tt.cfg, tt.plaintext)
				require.NoError(t, err)
				assert.NotEqual(t, ciphertext, ciphertext2, "两次加密应产生不同的密文")
			}
		})
	}
}

func TestDecrypt_InvalidInput(t *testing.T) {
	cfg := &Config{
		Passphrase: "test123",
		Salt:      []byte("testsalt12345678"),
		Iteration: 100000,
		KeyLen:    32,
	}

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "empty_input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "too_short_input",
			input:   []byte("too short"),
			wantErr: true,
		},
		{
			name:    "invalid_ciphertext",
			input:   bytes.Repeat([]byte("A"), 32),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(cfg, tt.input)
			assert.Error(t, err)
		})
	}
}

func TestEncryptDecrypt_DifferentConfigs(t *testing.T) {
	plaintext := []byte("Hello, World!")
	cfg1 := &Config{
		Passphrase: "test123",
		Salt:      []byte("testsalt12345678"),
		Iteration: 100000,
		KeyLen:    32,
	}
	cfg2 := &Config{
		Passphrase: "test456",
		Salt:      []byte("testsalt12345678"),
		Iteration: 100000,
		KeyLen:    32,
	}

	// 使用 cfg1 加密
	ciphertext, err := Encrypt(cfg1, plaintext)
	require.NoError(t, err)

	// 使用 cfg2 解密应该失败
	_, err = Decrypt(cfg2, ciphertext)
	assert.Error(t, err)
} 