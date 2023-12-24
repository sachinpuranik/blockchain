package wallet

import (
	// "crypto/ed25519"
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestGenertePrivateKey(t *testing.T) {
// 	pk, err := GeneratePrivateKey()
// 	require.NoError(t, err)
// 	require.Equal(t, len(pk.Bytes()), ed25519.PrivateKeySize)
// 	require.Equal(t, len(pk.Public().Bytes()), ed25519.PublicKeySize)
// }

// func TestSignMessage(t *testing.T) {
// 	pk, _ := GeneratePrivateKey()
// 	msg := "hello world"
// 	sig := pk.Sign(msg)
// 	assert.True(t, sig.Verify(pk.Public(), msg))
// 	assert.False(t, sig.Verify(pk.Public(), "abc"))

// 	differentKey, _ := GeneratePrivateKey()
// 	assert.False(t, sig.Verify(differentKey.Public(), msg))
// }

func TestWallet(t *testing.T) {
	a := NewWallet()
	b := NewWallet()
	require.NotEqual(t, a.privateKey, b.privateKey)
	require.NotEqual(t, a.publicKey, b.publicKey)
	require.NotEqual(t, a.PublicKeyStr(), b.PublicKeyStr())
}
