package jose

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
	"github.com/lestrrat-go/jwx/v4/jwt"

	"github.com/relexec/pkg/http/request"
)

// Authorizer authorizes HTTP requests using JSON Object Signing and Encryption
// (JOSE) auth infrastructure.
type Authorizer struct {
	alg        jwa.SignatureAlgorithm
	signKeys   jwk.Set
	verifyKeys jwk.Set
	// validateClaims contains any claims that will be validated against every
	// parsed token.
	validateClaims []jwt.ValidateOption
}

// Authorize performs authorization on the supplied HTTP request and returns an
// Auth with the validated identity. Authorize returns an error if the supplied
// HTTP request cannot be verified for the given conditions.
func (a Authorizer) Authorize(
	ctx huma.Context,
) (*request.Auth, error) {
	r, _ := humachi.Unwrap(ctx)
	token, err := jwt.ParseRequest(r, jwt.WithKey(a.alg, a.verifyKeys))
	if err != nil {
		return nil, fmt.Errorf("failed parsing token: %w", err)
	}
	rctx := ctx.Context()
	vopts := []jwt.ValidateOption{
		jwt.WithContext(rctx),
	}
	vopts = append(vopts, a.validateClaims...)
	err = jwt.Validate(token, vopts...)
	if err != nil {
		return nil, fmt.Errorf("failed validating token: %w", err)
	}
	return &request.Auth{
		Identity: "TODO",
	}, nil
}

var _ request.Authorizer = (*Authorizer)(nil)

// New returns an initialized Authorizer that uses JSON Object Signing and
// Encryption (JOSE) authentication and authorization infrastructure.
func New(
	ctx context.Context,
	cfg Config,
) (*Authorizer, error) {
	sigAlg := cfg.SignatureAlgorithm
	if sigAlg == "" {
		sigAlg = DefaultJWASignatureAlgorithm
	}
	ka, err := jwa.KeyAlgorithmFrom(string(sigAlg))
	if err != nil {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: %w", err,
		)
	}
	alg, ok := ka.(jwa.SignatureAlgorithm)
	if !ok {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: "+
				"%s is not a signature algorithm",
			sigAlg,
		)
	}

	signKeyPath := cfg.SignKeyPath
	if signKeyPath == "" {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: " +
				"sign_key_path is required",
		)
	}
	_, err = os.Stat(signKeyPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: "+
				"sign_key_path %s does not exist",
			signKeyPath,
		)
	} else if err != nil {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: "+
				"failed opening sign_key_path %s: %w",
			signKeyPath, err,
		)
	}

	signKeys, err := jwk.ParseFS(
		os.DirFS(filepath.Dir(signKeyPath)),
		filepath.Base(signKeyPath),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed initializing JOSE auth handler: "+
				"failed parsing sign key: %w",
			err,
		)
	}

	var verifyKeys jwk.Set

	if signatureAlgorithmAsymmetrical(sigAlg) {
		verifyKeyPath := cfg.VerifyKeyPath
		if verifyKeyPath == "" {
			return nil, fmt.Errorf(
				"failed initializing JOSE auth handler: " +
					"verify_key_path is required",
			)
		}
		_, err = os.Stat(verifyKeyPath)
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf(
				"failed initializing JOSE auth handler: "+
					"verify_key_path %s does not exist",
				verifyKeyPath,
			)
		} else if err != nil {
			return nil, fmt.Errorf(
				"failed initializing JOSE auth handler: "+
					"failed opening verify_key_path %s: %w",
				verifyKeyPath, err,
			)
		}

		verifyKeys, err = jwk.ParseFS(
			os.DirFS(filepath.Dir(verifyKeyPath)),
			filepath.Base(verifyKeyPath),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed initializing JOSE auth handler: "+
					"failed parsing verify key: %w",
				err,
			)
		}
	} else {
		verifyKeys = signKeys
	}

	var claims []jwt.ValidateOption
	if cfg.ValidateIssuer != "" {
		claims = append(claims, jwt.WithIssuer(cfg.ValidateIssuer))
	}
	if cfg.ValidateAudience != "" {
		claims = append(claims, jwt.WithAudience(cfg.ValidateAudience))
	}

	return &Authorizer{
		alg:            alg,
		signKeys:       signKeys,
		verifyKeys:     verifyKeys,
		validateClaims: claims,
	}, nil
}
