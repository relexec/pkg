package version

import (
	"sort"

	"github.com/Masterminds/semver/v3"
)

// Set wraps a map of major version to [semver.Version] for versions in that
// major version. The versions in each major version bucket are sorted earliest
// to latest.
type Set map[string]semver.Collection

// Add adds one or more supplied [semver.Version] to the Set.
func (s Set) Add(versions ...semver.Version) {
	for _, v := range versions {
		major, _ := Major(v)
		c, ok := s[major]
		if !ok {
			c = semver.Collection{}
		}
		c = append(c, &v)
		s[major] = c
	}
	for major, c := range s {
		sort.Sort(c)
		s[major] = c
	}
}

// Contains returns true if the supplied [semver.Version] is contained within
// the Set.
func (s Set) Contains(v semver.Version) bool {
	for _, c := range s {
		for _, cv := range c {
			if cv.Equal(&v) {
				return true
			}
		}
	}
	return false
}
