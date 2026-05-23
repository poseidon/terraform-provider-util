package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const jwksPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyutR5sZOMK+JwDafWm4i
5ncsyed9IwlY+bmsBwg89fYfggr/aO5Xno30njifscTDdbK2y0uon3OPXPjARipI
t6duEzlesdwKcBImBkyxe/I3GzS0ARpJtzpRMV+8mceNRbr3ce++SClEd/SvBWrR
/RJ1w862R86iwsS3k6LJbPTe+efT3ppyeAiLG/9ds7bQGwdRMT6Zr882eJl66ZwN
t5iZ8GebI8A1MLKSaJsCHFP8jRhKtjeme5/IcWGdwbaEKye5DXcgKHjCcewwHco2
pcmj4PQyBO/PI547oPlfSYQrcdN8VFiMBFxIscicLXbEH2GfQKj5cJ7/nbsC+gAh
wQIDAQAB
-----END PUBLIC KEY-----`

const jwksECDSAPublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEH6cuzP8XuD5wal6wf9M6xDljTOPL
X2i8uIp/C/ASqiIGUeeKQtX0/IR3qCXyThP/dbCiHrF3v1cuhBOHY8CLVg==
-----END PUBLIC KEY-----`

const jwksExample = `
data "util_jwks" "example" {
  public_key = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyutR5sZOMK+JwDafWm4i
5ncsyed9IwlY+bmsBwg89fYfggr/aO5Xno30njifscTDdbK2y0uon3OPXPjARipI
t6duEzlesdwKcBImBkyxe/I3GzS0ARpJtzpRMV+8mceNRbr3ce++SClEd/SvBWrR
/RJ1w862R86iwsS3k6LJbPTe+efT3ppyeAiLG/9ds7bQGwdRMT6Zr882eJl66ZwN
t5iZ8GebI8A1MLKSaJsCHFP8jRhKtjeme5/IcWGdwbaEKye5DXcgKHjCcewwHco2
pcmj4PQyBO/PI547oPlfSYQrcdN8VFiMBFxIscicLXbEH2GfQKj5cJ7/nbsC+gAh
wQIDAQAB
-----END PUBLIC KEY-----
EOT
}
`

const jwksExpected = `{"keys":[{"use":"sig","kty":"RSA","kid":"HvfxvoDusbEuKLtBer21gBVQO2m09IlZO_gcBCAhf0M","alg":"RS256","n":"yutR5sZOMK-JwDafWm4i5ncsyed9IwlY-bmsBwg89fYfggr_aO5Xno30njifscTDdbK2y0uon3OPXPjARipIt6duEzlesdwKcBImBkyxe_I3GzS0ARpJtzpRMV-8mceNRbr3ce--SClEd_SvBWrR_RJ1w862R86iwsS3k6LJbPTe-efT3ppyeAiLG_9ds7bQGwdRMT6Zr882eJl66ZwNt5iZ8GebI8A1MLKSaJsCHFP8jRhKtjeme5_IcWGdwbaEKye5DXcgKHjCcewwHco2pcmj4PQyBO_PI547oPlfSYQrcdN8VFiMBFxIscicLXbEH2GfQKj5cJ7_nbsC-gAhwQ","e":"AQAB"}]}`

const jwksECDSAExpected = `{"keys":[{"use":"sig","kty":"EC","kid":"SoABiieYuNx4UdqYvZRVeuC6SihxgLrhLy9peHMHpTc","crv":"P-256","alg":"ES256","x":"H6cuzP8XuD5wal6wf9M6xDljTOPLX2i8uIp_C_ASqiI","y":"BlHnikLV9PyEd6gl8k4T_3Wwoh6xd79XLoQTh2PAi1Y"}]}`

func TestJWKS_RSAPublicKey(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: jwksExample,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.util_jwks.example", "jwks", jwksExpected),
				),
			},
		},
	})

	got, err := jwksFromPublicKeyPEM(jwksPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if got != jwksExpected {
		t.Fatalf("got %q, want %q", got, jwksExpected)
	}
}

func TestJWKS_ECDSAPublicKey(t *testing.T) {
	got, err := jwksFromPublicKeyPEM(jwksECDSAPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if got != jwksECDSAExpected {
		t.Fatalf("got %q, want %q", got, jwksECDSAExpected)
	}
}

func TestJWKSInvalidPEM(t *testing.T) {
	_, err := jwksFromPublicKeyPEM("not a PEM public key")
	if err == nil || !strings.Contains(err.Error(), "failed to decode PEM public key") {
		t.Fatalf("got error %v, want failed PEM decode error", err)
	}
}

func TestJWKSWrongPEMBlockType(t *testing.T) {
	content := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("not a public key"),
	})

	_, err := jwksFromPublicKeyPEM(string(content))
	if err == nil || !strings.Contains(err.Error(), `unsupported PEM block type "CERTIFICATE"`) {
		t.Fatalf("got error %v, want unsupported PEM block type error", err)
	}
}

func TestJWKSUnsupportedPublicKey(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatal(err)
	}
	content := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	_, err = jwksFromPublicKeyPEM(string(content))
	if err == nil || !strings.Contains(err.Error(), "only RSA and ECDSA public keys are supported") {
		t.Fatalf("got error %v, want unsupported public key type error", err)
	}
}
