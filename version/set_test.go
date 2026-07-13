package version_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/relexec/pkg/version"
	"github.com/stretchr/testify/require"
)

var (
	v0_0_0 = semver.MustParse("0.0.0")
	v0_0_1 = semver.MustParse("0.0.1")
	v0_1_0 = semver.MustParse("0.1.0")
	v0_1_1 = semver.MustParse("0.1.1")
	v1_0_0 = semver.MustParse("1.0.0")
	v1_0_1 = semver.MustParse("1.0.1")
	v1_1_0 = semver.MustParse("1.1.0")
	v1_1_1 = semver.MustParse("1.1.1")
	v2_0_0 = semver.MustParse("2.0.0")
	v2_0_1 = semver.MustParse("2.0.1")
)

func TestSet(t *testing.T) {

	cases := []struct {
		name string
		add  []semver.Version
		exp  version.Set
	}{
		{
			"empty set",
			[]semver.Version{},
			version.Set{},
		},
		{
			"one in zero major",
			[]semver.Version{
				*v0_0_0,
			},
			version.Set{
				"0": semver.Collection{
					v0_0_0,
				},
			},
		},
		{
			"sorted zero major",
			[]semver.Version{
				*v0_0_1,
				*v0_0_0,
				*v0_1_0,
				*v0_1_1,
			},
			version.Set{
				"0": semver.Collection{
					v0_0_0,
					v0_0_1,
					v0_1_0,
					v0_1_1,
				},
			},
		},
		{
			"sorted multiple majors",
			[]semver.Version{
				*v0_1_1,
				*v0_0_1,
				*v1_1_0,
				*v1_0_0,
			},
			version.Set{
				"0": semver.Collection{
					v0_0_1,
					v0_1_1,
				},
				"1": semver.Collection{
					v1_0_0,
					v1_1_0,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			vs := version.Set{}
			vs.Add(c.add...)
			require.Equal(c.exp, vs)
		})
	}
}
