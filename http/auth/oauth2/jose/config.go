package jose

// Config contains configuration for the Javascript Object Signing and
// Encryption (JOSE) auth infrastructure.
//
// Ref: https://www.redhat.com/en/blog/jose-json-object-signing-and-encryption
type Config struct {
	// SignatureAlgorithm contains the JSON Web Algorithms (JWA) signature
	// algorithm. Defaults to "RS256".
	SignatureAlgorithm JWASignatureAlgorithm `json:"signature_algorithm,omitempty"`
	// SignKeyPath points at the file containing the contents of the signing
	// key. When using key pairs (PKE), this is the private key file.
	SignKeyPath string `json:"sign_key_path,omitempty"`
	// VerifyKeyPath points at the file containing the contents of the
	// verifying key. When using key pairs (PKE), this is the public key file.
	VerifyKeyPath string `json:"verify_key_path,omitempty"`
	// ValidateIssuer is the exact issuer that should match the "iss" claim in
	// the supplied JWT token. Leave empty to not validate the issuer (not
	// recommended).
	ValidateIssuer string `json:"validate_issuer,omitempty"`
	// ValidateAudience is the exact audience that should match the "aud" claim
	// in the supplied JWT token. Leave empty to not validate the audience (not
	// recommended).
	ValidateAudience string `json:"validate_audience,omitempty"`
	// JWK contains configuration for the JSON Web Key (JWK) parts of the JOSE
	// auth infrastructure.
	JWK *JWKConfig `json:"jwk,omitempty"`
}

// JWKConfig contains configuration for the JSON Web Key (JWK) parts of the
// JOSE auth infrastructure.
//
// Ref: https://datatracker.ietf.org/doc/html/rfc7517
type JWKConfig struct {
	// URL is the URL to find the JWKs used in authenticating JWTs.
	URL string `json:"url"`
}
