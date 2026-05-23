package internal

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceJWKS() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceJWKSRead,

		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PEM-encoded RSA or ECDSA public key",
			},
			"jwks": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes-compatible JSON Web Key Set for the public key",
			},
		},
	}
}

func datasourceJWKSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	jwks, err := jwksFromPublicKeyPEM(d.Get("public_key").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("jwks", jwks); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(Hashcode(jwks))
	return diags
}

type jwksDocument struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Use       string `json:"use"`
	KeyType   string `json:"kty"`
	KeyID     string `json:"kid"`
	Curve     string `json:"crv,omitempty"`
	Algorithm string `json:"alg"`
	Modulus   string `json:"n,omitempty"`
	Exponent  string `json:"e,omitempty"`
	X         string `json:"x,omitempty"`
	Y         string `json:"y,omitempty"`
}

// jwksFromPublicKeyPEM converts one PEM-encoded public key into a JWKS document
// compatible with Kubernetes service account issuer discovery.
func jwksFromPublicKeyPEM(content string) (string, error) {
	publicKey, err := parsePublicKeyPEM(content)
	if err != nil {
		return "", err
	}

	key, err := publicKeyToJWK(publicKey)
	if err != nil {
		return "", err
	}
	jwks := jwksDocument{
		Keys: []jwk{key},
	}

	data, err := json.Marshal(jwks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWKS: %w", err)
	}
	return string(data), nil
}

// parsePublicKeyPEM decodes PEM-formatted content containing a PKIX public key
// or a PKCS1 RSA public key.
func parsePublicKeyPEM(content string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(content))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM public key")
	}

	switch block.Type {
	case "PUBLIC KEY":
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKIX public key: %w", err)
		}
		return publicKey, nil
	case "RSA PUBLIC KEY":
		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS1 RSA public key: %w", err)
		}
		return publicKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type %q: expected PUBLIC KEY or RSA PUBLIC KEY", block.Type)
	}
}

// publicKeyToJWK converts supported public keys to the same public JWK shape
// Kubernetes publishes from /openid/v1/jwks.
func publicKeyToJWK(publicKey crypto.PublicKey) (jwk, error) {
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return rsaPublicKeyToJWK(publicKey)
	case *ecdsa.PublicKey:
		return ecdsaPublicKeyToJWK(publicKey)
	default:
		return jwk{}, fmt.Errorf("unsupported public key type %T: only RSA and ECDSA public keys are supported", publicKey)
	}
}

func rsaPublicKeyToJWK(publicKey *rsa.PublicKey) (jwk, error) {
	keyID, err := kubernetesKeyIDFromPublicKey(publicKey)
	if err != nil {
		return jwk{}, err
	}

	return jwk{
		Use:       "sig",
		KeyType:   "RSA",
		KeyID:     keyID,
		Algorithm: "RS256",
		Modulus:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		Exponent:  base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}, nil
}

func ecdsaPublicKeyToJWK(publicKey *ecdsa.PublicKey) (jwk, error) {
	keyID, err := kubernetesKeyIDFromPublicKey(publicKey)
	if err != nil {
		return jwk{}, err
	}

	var curve, algorithm string
	switch publicKey.Curve {
	case elliptic.P256():
		curve = "P-256"
		algorithm = "ES256"
	case elliptic.P384():
		curve = "P-384"
		algorithm = "ES384"
	case elliptic.P521():
		curve = "P-521"
		algorithm = "ES512"
	default:
		return jwk{}, fmt.Errorf("unknown ECDSA public key curve: must be P-256, P-384, or P-521")
	}

	size := (publicKey.Curve.Params().BitSize + 7) / 8
	return jwk{
		Use:       "sig",
		KeyType:   "EC",
		KeyID:     keyID,
		Curve:     curve,
		Algorithm: algorithm,
		X:         base64.RawURLEncoding.EncodeToString(fixedLengthBytes(publicKey.X, size)),
		Y:         base64.RawURLEncoding.EncodeToString(fixedLengthBytes(publicKey.Y, size)),
	}, nil
}

func fixedLengthBytes(value *big.Int, size int) []byte {
	data := value.Bytes()
	if len(data) >= size {
		return data
	}

	padded := make([]byte, size)
	copy(padded[size-len(data):], data)
	return padded
}

// kubernetesKeyIDFromPublicKey derives a Kubernetes-compatible service account
// key ID (`kid`) from a public key.
func kubernetesKeyIDFromPublicKey(publicKey crypto.PublicKey) (string, error) {
	// This must match kube-apiserver's service account token implementation.
	// Kubernetes puts the same Key ID (`kid`) in JWT headers and in the JWKS
	// entry served from /openid/v1/jwks, so relying parties can choose the
	// right verification key.
	//
	// There are several plausible but incompatible ways to identify a key:
	// hashing the PEM text, hashing only the RSA modulus, using a certificate
	// thumbprint, or using the RFC 7638 JWK thumbprint. Kubernetes does none
	// of those for its built-in service account signer. Instead, it marshals
	// the public key as DER SubjectPublicKeyInfo with x509.MarshalPKIXPublicKey,
	// hashes those bytes with SHA-256, then base64url-encodes the digest without
	// padding.
	//
	// Keep this aligned with upstream Kubernetes defaults:
	// https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/jwt.go#L90-L110
	//   keyIDFromPublicKey derives the kid from DER SubjectPublicKeyInfo.
	// https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/jwt.go#L273-L291
	//   StaticPublicKeysGetter attaches that kid to each public key.
	// https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/openidmetadata.go#L237-L252
	//   openIDKeysetJSON renders the JWKS document.
	// https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/openidmetadata.go#L291-L313
	//   jwkFromPublicKey preserves the kid and adds use=sig and the matching
	//   signing algorithm for RSA or ECDSA keys.
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to serialize public key to DER format: %w", err)
	}

	hash := sha256.Sum256(publicKeyDERBytes)
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}
