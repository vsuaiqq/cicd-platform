package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GenerateSSHKeyPair() (privatePEM []byte, publicKey string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}

	privPEM, err := ssh.MarshalPrivateKey(priv, "deploy key")
	if err != nil {
		return nil, "", err
	}
	privatePEM = pem.EncodeToMemory(privPEM)

	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, "", err
	}
	publicKey = strings.TrimRight(string(ssh.MarshalAuthorizedKey(sshPub)), "\n")

	return privatePEM, publicKey, nil
}

func GenerateWebhookSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
