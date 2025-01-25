package keys

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type RSAKeys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func LoadKeys() RSAKeys {
	/** load keys */
	privateKeyBytes, err := os.ReadFile(os.Getenv("PRIVATE_KEY"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open private key: %v\n", err)
		os.Exit(1)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse private key: %v\n", err)
		os.Exit(1)
	}

	publicKeyBytes, err := os.ReadFile(os.Getenv("PUBLIC_KEY"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open public key: %v\n", err)
		os.Exit(1)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	return RSAKeys{
		Private: privateKey,
		Public:  publicKey,
	}
}
