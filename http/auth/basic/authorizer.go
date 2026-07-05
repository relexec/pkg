package basic

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/relexec/pkg/http/request"
)

type cred struct {
	username      string
	password      string
	permissions   []string
	organizations []string
	domains       []string
	namespaces    []string
}

// Authorizer that examines HTTP Authorization header and checks against known
// credentials.
//
// NOTE(jaypipes): This Authorizer should never be used in production. It is
// only useful for local developer testing.
type Authorizer struct {
	// creds is a map of credentials, keyed by username.
	creds map[string]cred
}

// Authorize performs authorization on the supplied HTTP request and returns an
// Auth with the validated identity. Authorize returns an error if the supplied
// HTTP request cannot be verified for the given conditions.
func (a Authorizer) Authorize(
	ctx huma.Context,
) (*request.Auth, error) {
	authHeader := ctx.Header("Authorization")

	if authHeader == "" {
		return nil, fmt.Errorf("Authorization header required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return nil, fmt.Errorf("invalid Authorization header")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid Authorization header")
	}
	cred, err := a.credentialsFor(string(decoded))
	if err != nil {
		return nil, err
	}
	return &request.Auth{Identity: cred.username}, nil
}

func (a Authorizer) credentialsFor(header string) (*cred, error) {
	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Authorization header")
	}
	u := parts[0]
	pw := parts[1]
	c, ok := a.creds[u]
	if !ok {
		return nil, fmt.Errorf("unauthorized user")
	}
	if c.password != pw {
		return nil, fmt.Errorf("unauthorized user")
	}
	return &c, nil
}

var _ request.Authorizer = (*Authorizer)(nil)

// New returns an initialized Authorizer only useful for development and
// testing.
func New(
	ctx context.Context,
	cfg Config,
) (*Authorizer, error) {
	creds := map[string]cred{}
	scanner := bufio.NewScanner(strings.NewReader(testCredentials))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 6 {
			return nil, fmt.Errorf(
				"expected line %q to contain 6 parts but got %d",
				line, len(parts),
			)
		}
		username := parts[0]
		password := parts[1]
		perms := strings.Split(parts[2], ",")
		orgs := strings.Split(parts[3], ",")
		doms := strings.Split(parts[4], ",")
		nss := strings.Split(parts[5], ",")
		creds[username] = cred{
			username:      username,
			password:      password,
			permissions:   perms,
			organizations: orgs,
			domains:       doms,
			namespaces:    nss,
		}
	}
	return &Authorizer{
		creds: creds,
	}, nil
}
