package auth

import (
	"github.com/cenkalti/backoff"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/zenoss/zenkit/logging"
)

// GetKeysFromFS creates a slice of jwt.Key from keys in files
func GetKeysFromFS(logger logging.ErrorLogger, files []string) ([]jwt.Key, error) {
	var parsedKeys []jwt.Key
	for _, keyFile := range files {
		keyBytes, err := ReadKeyFromFS(logger, keyFile)
		if err != nil {
			return []jwt.Key{}, errors.Wrap(err, "Unable to read key file from fs")
		}
		key := ConvertToKey(keyBytes)
		parsedKeys = append(parsedKeys, key)
	}
	return parsedKeys, nil
}

// ConvertToKey converts a byte slice to a jwt.Key
func ConvertToKey(key []byte) jwt.Key {
	pubkey, err := jwtgo.ParseRSAPublicKeyFromPEM(key)
	if err == nil {
		return pubkey
	}

	return key
}

func ReadKeyFromFS(logger logging.ErrorLogger, filename string) ([]byte, error) {
	// Get the secret key
	var key []byte
	readKey := func() error {
		data, err := afero.ReadFile(FS, filename)
		if err != nil {
			logger.LogError("Unable to load auth key. Retrying.", "keyfile", filename, "err", err)
			return errors.Wrap(err, "unable to load auth key")
		}
		key = data
		return nil
	}
	// Docker sometimes doesn't mount the secret right away, so we'll do a short retry
	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = KeyFileTimeout
	if err := backoff.Retry(readKey, boff); err != nil {
		return nil, errors.Wrap(err, "unable to load auth key within the timeout")
	}
	return key, nil
}
