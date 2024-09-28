package keygen

import (
	"os"
)

const KeySize = 256

func SaveKeys(private, public []byte, privateKeyPath, publicKeyPath string) error {
	if err := os.WriteFile(privateKeyPath, private, 0666); err != nil {
		return err
	}
	if err := os.WriteFile(publicKeyPath, public, 0666); err != nil {
		return err
	}
	return nil
}
