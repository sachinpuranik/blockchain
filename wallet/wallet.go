package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

var addressLen = 20

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockChainAddress string
}

func NewWallet() *Wallet {
	// 1. Create public and private Keys
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey

	// 2. Perform sha256 hash on public key.
	h2 := sha256.New()
	h2.Sum(w.publicKey.X.Bytes())
	h2.Sum(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. perform RIPEMD-160 hashing on the results of 20 bytes from the result of SHA256
	h3 := ripemd160.New()
	h3.Sum(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for main network) (versioned digest 4 vd4)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3)
	// 5.  Perform sha256 hash  on vd4
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6.  Perform sha256 hash  on digest5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7.  take first 4 bytes of digest6 for checksum
	checksum := digest6[:4]
	// 8. Add the checksum  bytes from step 7 , at the end of extended RIPEMD-160 hash form result of step4 .(21+4 = 25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[22:25], checksum[:])
	// 9. Convert the result of the byte string in to base58
	w.blockChainAddress = base58.Encode(dc8)
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockChainAddress() string {
	return w.blockChainAddress
}
