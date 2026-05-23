# jwks Data Source

`jwks` converts a PEM-encoded RSA or ECDSA public key into a JSON Web Key Set (JWKS) document. The output matches Kubernetes service account issuer discovery defaults, so it can be used when serving `/openid/v1/jwks` from static OIDC discovery metadata.

## Usage

```hcl
data "util_jwks" "cluster" {
  public_key = module.cluster.service_account_public_key
}

module "oidc-discovery" {
  source       = "some/module/for/oidc-discovery"
  cluster_name = "cluster"
  jwks         = data.util_jwks.cluster.jwks
}
```

## Kubernetes Compatibility

Kubernetes derives the `kid` for service account signing keys by marshaling the public key as DER SubjectPublicKeyInfo, hashing those bytes with SHA-256, and base64url-encoding the digest without padding.

This is intentionally different from hashing the PEM text, hashing only the RSA modulus, using a certificate thumbprint, or using an RFC 7638 JWK thumbprint. Matching Kubernetes is important because service account JWT headers and JWKS entries must use the same `kid`.

RSA keys are published with `alg = "RS256"`. ECDSA keys use the same curve mapping as Kubernetes: P-256 uses `ES256`, P-384 uses `ES384`, and P-521 uses `ES512`.

See upstream Kubernetes:

* [`keyIDFromPublicKey`](https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/jwt.go#L90-L110)
* [`StaticPublicKeysGetter`](https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/jwt.go#L273-L291)
* [`openIDKeysetJSON`](https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/openidmetadata.go#L237-L252)
* [`jwkFromPublicKey`](https://github.com/kubernetes/kubernetes/blob/v1.36.1/pkg/serviceaccount/openidmetadata.go#L291-L313)

## Argument Reference

* `public_key` - PEM-encoded RSA or ECDSA public key

## Argument Attributes

* `jwks` - Kubernetes-compatible JSON Web Key Set for the public key
