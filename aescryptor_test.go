package boltsec

import (
	"bytes"
	"testing"
)

func TestEnc(t *testing.T) {
	data := []struct {
		content []byte
		secret  []byte
	}{
		{[]byte("Foo"), []byte("Boo")},
		{[]byte("Foo"), []byte("Car")},
		{[]byte("Foo"), []byte("")},
		{[]byte(""), []byte("Car")},
		{[]byte("Long input with more than 16 characters"), []byte("Car")},
		{[]byte(`"{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"`), []byte("input with more than 16 characters")},
		{[]byte(`"{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"
            "{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}{"info":{"name":"gXeMfp.zip"
            ,"type":"","size":79448, "comment":"test"}}{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"
            }}{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}{"info":{"name":"gXeMfp.zip","type":"",
            "size":79448, "comment":"test"}}{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}{"info":
            {"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"`), []byte("input with more than 16 characters")},
	}

	for _, iter := range data {
		//Logger.Printf("TestEnc for key: %v, content:%v", iter.secret, iter.content)
		ac, err := newAESCryptor(iter.secret)
		if err != nil {
			t.Errorf("newAESCryptor with key '%v' return err: %s", iter.secret, err)
			continue
		}
		enc, err := ac.encrypt(iter.content)
		if err != nil {
			t.Errorf("Unable to encrypt '%v' with key '%v': %v", iter.content, iter.secret, err)
			continue
		}
		dec, err := ac.decrypt(enc)
		if err != nil {
			t.Errorf("Unable to decrypt '%v' with key '%v': %v", enc, iter.secret, err)
			continue
		}
		if !bytes.Equal(dec, iter.content) {
			t.Errorf("Decrypt Key %v\n  Input: %v\n  Expect: %v\n  Actual: %v", iter.secret, enc, iter.content, dec)
		}
	}
}

func BenchmarkInitLongSecret(b *testing.B) {
	secret := `"{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"`
	for n := 0; n < b.N; n++ {
		_, err := newAESCryptor([]byte(secret))
		if err != nil {
			b.Errorf("newAESCryptor with key '%v' return err: %s", secret, err)
			continue
		}
	}
}

func BenchmarkInitShortSecret(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := newAESCryptor([]byte("secret"))
		if err != nil {
			b.Errorf("newAESCryptor with key '%v' return err: %s", []byte("secret"), err)
			continue
		}
	}
}

func BenchmarkReuseCryptor(b *testing.B) {
	content := []byte(`"{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"`)
	secret := []byte("secret")

	ac, err := newAESCryptor(secret)
	if err != nil {
		b.Errorf("newAESCryptor with key '%v' return err: %s", secret, err)
	}

	for n := 0; n < b.N; n++ {

		enc, err := ac.encrypt(content)
		if err != nil {
			b.Errorf("Unable to encrypt '%v' with key '%v': %v", content, secret, err)
		}
		dec, err := ac.decrypt(enc)
		if err != nil {
			b.Errorf("Unable to decrypt '%v' with key '%v': %v", enc, secret, err)
		}
		if !bytes.Equal(dec, content) {
			b.Errorf("Decrypt Key %v\n  Input: %v\n  Expect: %v\n  Actual: %v", secret, enc, content, dec)
		}
	}
}

func BenchmarkNewCryptor(b *testing.B) {
	content := []byte(`"{"info":{"name":"gXeMfp.zip","type":"","size":79448, "comment":"test"}}"`)
	secret := []byte("secret")

	for n := 0; n < b.N; n++ {
		ac, err := newAESCryptor(secret)
		if err != nil {
			b.Errorf("newAESCryptor with key '%v' return err: %s", secret, err)
		}

		enc, err := ac.encrypt(content)
		if err != nil {
			b.Errorf("Unable to encrypt '%v' with key '%v': %v", content, secret, err)
		}
		dec, err := ac.decrypt(enc)
		if err != nil {
			b.Errorf("Unable to decrypt '%v' with key '%v': %v", enc, secret, err)
		}
		if !bytes.Equal(dec, content) {
			b.Errorf("Decrypt Key %v\n  Input: %v\n  Expect: %v\n  Actual: %v", secret, enc, content, dec)
		}
	}
}
