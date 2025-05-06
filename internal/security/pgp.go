package security

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

var pubKeyRing openpgp.EntityList

func init() {
	keyData, err := os.ReadFile("internal/security/keys/pubkey.asc")
	if err != nil {
		panic("cannot read PGP public key: " + err.Error())
	}
	pubKeyRing, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(keyData))
	if err != nil {
		panic("cannot parse PGP public key: " + err.Error())
	}
}

func EncryptPGP(plainText string) ([]byte, error) {
	var buf bytes.Buffer

	armorWriter, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return nil, err
	}

	encryptWriter, err := openpgp.Encrypt(armorWriter, pubKeyRing, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	_, err = encryptWriter.Write([]byte(plainText))
	if err != nil {
		return nil, err
	}
	encryptWriter.Close()
	armorWriter.Close()

	return buf.Bytes(), nil
}

func DecryptPGP(cipher []byte) (string, error) {
	keyData, err := os.ReadFile("internal/security/keys/privkey.asc")
	if err != nil {
		return "", err
	}

	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(keyData))
	if err != nil {
		return "", err
	}

	pass := []byte(os.Getenv("PGP_PASSPHRASE"))

	for _, entity := range entityList {
		if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
			if err := entity.PrivateKey.Decrypt(pass); err != nil {
				return "", fmt.Errorf("cannot decrypt private key: %w", err)
			}
		}
		for _, sub := range entity.Subkeys {
			if sub.PrivateKey != nil && sub.PrivateKey.Encrypted {
				if err := sub.PrivateKey.Decrypt(pass); err != nil {
					return "", fmt.Errorf("failed to decrypt subkey: %w", err)
				}
			}
		}
	}

	block, err := armor.Decode(bytes.NewReader(cipher))
	if err != nil {
		return "", err
	}

	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
