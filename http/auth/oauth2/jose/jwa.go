package jose

// JWASignatureAlgorithm is a string designating the JSON Web Algorithm (JWA)
// signature algorithm.
//
// Ref: https://datatracker.ietf.org/doc/html/rfc7518
type JWASignatureAlgorithm string

const (
	// HMAC signature algorithm using SHA-256
	JWASignatureAlgorithmHS256 JWASignatureAlgorithm = "HS256"
	// HMAC signature algorithm using SHA-384
	JWASignatureAlgorithmHS384 JWASignatureAlgorithm = "HS384"
	// HMAC signature algorithm using SHA-512
	JWASignatureAlgorithmHS512 JWASignatureAlgorithm = "HS512"
	// RSASSA-PKCD-v1.5 signature algorithm using SHA-256
	JWASignatureAlgorithmRS256 JWASignatureAlgorithm = "RS256"
	// RSASSA-PKCD-v1.5 signature algorithm using SHA-384
	JWASignatureAlgorithmRS384 JWASignatureAlgorithm = "RS384"
	// RSASSA-PKCD-v1.5 signature algorithm using SHA-512
	JWASignatureAlgorithmRS512 JWASignatureAlgorithm = "RS512"
	// ECDSA signature algorithm using P-256 curve SHA-256
	JWASignatureAlgorithmES256 JWASignatureAlgorithm = "ES256"
	// ECDSA signature algorithm using P-384 curve SHA-384
	JWASignatureAlgorithmES384 JWASignatureAlgorithm = "ES384"
	// ECDSA signature algorithm using P-512 curve SHA-512
	JWASignatureAlgorithmES512 JWASignatureAlgorithm = "ES512"
)

const (
	DefaultJWASignatureAlgorithm = JWASignatureAlgorithmRS256
)

// signatureAlgorithmAsymmetrical returns true if the supplied JWA signature
// algorithm is asymmetrical.
func signatureAlgorithmAsymmetrical(alg JWASignatureAlgorithm) bool {
	switch alg {
	case
		JWASignatureAlgorithmHS256,
		JWASignatureAlgorithmHS384,
		JWASignatureAlgorithmHS512:
		return false
	default:
		return true
	}
}
