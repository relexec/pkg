package inspect_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/relexec/pkg/inspect"
	"github.com/stretchr/testify/require"
)

func TestIsUUID(t *testing.T) {
	uuidStr := uuid.NewString()
	uuidObj, _ := uuid.Parse(uuidStr)

	// last segment only 11...
	wrongLenUUIDDashes := "47415846-6dfa-4cb1-8b96-c3e91bf6fcc"
	wrongLenUUID := "474158466dfa4cb18b96c3e91bf6fcc"
	badCharUUID := "47415846-6dfa-4cb1-8b96-c3e91bf6fccX"
	dashesWrongUUID := "4741584-66dfa-4cb1-8b96-c3e91bf6fcca"

	cases := []struct {
		name    string
		subject any
		exp     bool
	}{
		{
			"google UUID is a UUID",
			uuidObj,
			true,
		},
		{
			"valid UUID string is a UUID",
			uuidStr,
			true,
		},
		{
			"valid UUID []byte is a UUID",
			[]byte(uuidStr),
			true,
		},
		{
			"nil isn't a UUID",
			nil,
			false,
		},
		{
			"int isn't a UUID",
			1,
			false,
		},
		{
			"wrong length with dashes",
			wrongLenUUIDDashes,
			false,
		},
		{
			"wrong length no dashes",
			wrongLenUUID,
			false,
		},
		{
			"bad char",
			badCharUUID,
			false,
		},
		{
			"dashes in wrong place",
			dashesWrongUUID,
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			got := inspect.IsUUID(c.subject)
			require.Equal(c.exp, got)
		})
	}
}
