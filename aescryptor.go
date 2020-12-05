package boltsec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// The AESCryptor struct to keep the key and block value which can
// be reused by all the encrypt and decrypt actions as those values
// are static and won't be changed as long as the secret is not changed.
type aesCryptor struct {
	rawkey []byte
	key    []byte
	block  cipher.Block
}

// The newAESCryptor return a pointer to the AESCryptor struct,
// error is returned if any error happened to the aes.NewCipher
func newAESCryptor(secret []byte) (result *aesCryptor, err error) {
	result = new(aesCryptor)
	result.rawkey = secret
	data := sha256.Sum256(secret)
	result.key = data[0:]

	result.block, err = aes.NewCipher(result.key)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// The encrypt function encrypt the data with the cipher value that is calculated
// when the AESCryptor is initialized
func (ac *aesCryptor) encrypt(data []byte) ([]byte, error) {
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(ac.block, iv)

	stream.XORKeyStream(encrypted, data)
	return output, nil
}

// The decrypt function decrypt the data with the cipher value that is calculated
// when the AESCryptor is initialized. Be cautious that the decrypt directly update
// the decrypted value in the data field, thus make sure the data field is modifiable,
// otherwise copy the original encrypted content to a new []byte before calling this function
func (ac *aesCryptor) decrypt(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return []byte(""), errors.New("cipherText too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(ac.block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}
