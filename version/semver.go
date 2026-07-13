package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// Major returns the major version number from the supplied argument. The
// argument can be a string, e.g. "1" or "v3" or a [semver.Version].
//
// Returns an empty string and error if the supplied argument isn't a valid
// semver or major version number.
func Major(ver any) (string, error) {
	switch ver := ver.(type) {
	case semver.Version:
		return strconv.Itoa(int(ver.Major())), nil
	case *semver.Version:
		return strconv.Itoa(int(ver.Major())), nil
	case string:
		// We do a quick check here to see if the user passed just a
		// major version, optionally prefixed with a "v". If so, we
		// avoid trying to parse the supplied string into a proper
		// semver.Version.
		if !strings.Contains(ver, ".") {
			ver = strings.TrimPrefix(ver, "v")
			_, err := strconv.Atoi(ver)
			if err != nil {
				return "", fmt.Errorf("failed to convert %s to integer", ver)
			}
			return ver, nil
		}
		sv, err := semver.NewVersion(ver)
		if err == nil {
			return strconv.Itoa(int(sv.Major())), nil
		}
		return "", err
	default:
		return "", fmt.Errorf(
			"unsupported type %T passed to Major", ver,
		)
	}
}
