package main

import (
	"bytes"
	"crypto/sha512"
	"github.com/bwesterb/go-ristretto"
	"github.com/r2ishiguro/vrf/go/vrf_ed25519"
	"github.com/r2ishiguro/vrf/go/vrf_ed25519/edwards25519"
	"golang.org/x/crypto/sha3"
)

const (
	PublicKeySize    = 32
	SecretKeySize    = 64
	Size             = 32
	N2               = 32
	N                = N2 / 2
	intermediateSize = 32
	ProofSize        = 32 + 32 + intermediateSize
)

// GenerateKey creates a public/secret key pair.
func GenerateKey() (*[PublicKeySize]byte, *[SecretKeySize]byte) {
	var secretKey ristretto.Scalar
	var pk = new([PublicKeySize]byte)
	var sk = new([SecretKeySize]byte)
	var digest [64]byte

	secretKey.Rand() // Generate a new secret key
	copy(sk[:32], secretKey.Bytes())

	h := sha3.NewShake256()
	h.Write(sk[:32])
	h.Read(digest[:])

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A ristretto.Point
	var hBytes [32]byte
	copy(hBytes[:], digest[:32])

	var hBytesScalar ristretto.Scalar
	hBytesScalar.SetBytes(&hBytes)

	A.ScalarMultBase(&hBytesScalar) // compute public key
	A.BytesInto(pk)

	copy(sk[32:], pk[:])
	return pk, sk
}

func Evaluate(sk []byte, pk []byte, m []byte) (r []byte, pi []byte) {
	x := expandSecret(sk)
	h := vrf_ed25519.ECVRF_hash_to_curve(m, pk)
	r = vrf_ed25519.ECP2OS(vrf_ed25519.GeScalarMult(h, x))
	pi, err := vrf_ed25519.ECVRF_prove(pk, sk, m)
	if err != nil {
		return
	}
	return r, pi
}

func Verify(pk []byte, pi []byte, m []byte, r []byte) (bool, error) {
	r1 := decode_proof(pi)
	if !bytes.Equal(r, vrf_ed25519.ECP2OS(r1)) {
		return false, nil
	}
	return vrf_ed25519.ECVRF_verify(pk, pi, m)
}

func expandSecret(sk []byte) *[32]byte {
	// copied from golang.org/x/crypto/ed25519/ed25519.go -- has to be the same
	digest := sha512.Sum512(sk[:32])
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64
	h := new([32]byte)
	copy(h[:], digest[:])
	return h
}

func decode_proof(pi []byte) (r *edwards25519.ExtendedGroupElement) {
	i := 0
	sign := pi[i]
	i++
	if sign != 2 && sign != 3 {
		return nil
	}
	r = vrf_ed25519.OS2ECP(pi[i:i+N2], sign-2)
	i += N2
	if r == nil {
		return nil
	}
	return
}
