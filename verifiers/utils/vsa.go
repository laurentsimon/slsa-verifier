package utils

import (
	"fmt"
	"strings"

	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
)

type TrustedVerifierID struct {
	name, version string
}

// TrustedVerifierIDNew creates a new VerifierID structure.
func TrustedVerifierIDNew(VerifierID string) (*TrustedVerifierID, error) {
	name, version, err := ParseVerifierID(VerifierID, true)
	if err != nil {
		return nil, err
	}

	return &TrustedVerifierID{
		name:    name,
		version: version,
	}, nil
}

// Matches matches the VerifierID string against the reference VerifierID.
// If the VerifierID contains a semver, the full VerifierID must match.
// Otherwise, only the name needs to match.
// `allowRef: true` indicates that the matching need not be an eaxct
// match. In this case, if the VerifierID version is a GitHub ref
// `refs/tags/name`, we will consider it equal to user-provided
// VerifierID `name`.
func (b *TrustedVerifierID) Matches(VerifierID string, allowRef bool) error {
	name, version, err := ParseVerifierID(VerifierID, false)
	if err != nil {
		return err
	}

	if name != b.name {
		return fmt.Errorf("%w: expected name '%s', got '%s'", serrors.ErrorMismatchVerifierID,
			name, b.name)
	}

	if version != "" && version != b.version {
		// If allowRef is true, try the long version `refs/tags/<name>` match.
		if allowRef &&
			"refs/tags/"+version == b.version {
			return nil
		}
		return fmt.Errorf("%w: expected version '%s', got '%s'", serrors.ErrorMismatchVerifierID,
			version, b.version)
	}

	return nil
}

func (b *TrustedVerifierID) Name() string {
	return b.name
}

func (b *TrustedVerifierID) Version() string {
	return b.version
}

func (b *TrustedVerifierID) String() string {
	// TODO: Update to using `@`
	return fmt.Sprintf("%s/%s", b.name, b.version)
}

func ParseVerifierID(id string, needVersion bool) (string, string, error) {
	// WARNING: temporary implementation for Google VSA.
	parts := strings.Split(id, "/")
	if len(parts) < 1 {
		return "", "", fmt.Errorf("%w: verifierID: '%s'",
			serrors.ErrorInvalidFormat, id)
	}
	version := parts[len(parts)-1]
	// This happens if the user-provided ID is of the form `<name>/`.
	if needVersion {
		if version == "" {
			return "", "", fmt.Errorf("%w: verifierID: '%s'",
				serrors.ErrorInvalidFormat, id)
		}
		if !strings.HasPrefix(version, "v") {
			return "", "", fmt.Errorf("%w: verifierID: '%s'",
				serrors.ErrorInvalidFormat, id)
		}
	}

	// If the latest part of the ID is not a version,
	// use all the parts to construct the ID with an
	// empty version.
	if !strings.HasPrefix(version, "v") {
		return strings.Join(parts, "/"), "", nil
	}

	return strings.Join(parts[:len(parts)-1], "/"), version, nil
	/*
		parts := strings.Split(id, "@")
		if len(parts) == 2 {
			if parts[1] == "" {
				return "", "", fmt.Errorf("%w: verifierID: '%s'",
					serrors.ErrorInvalidFormat, id)
			}
			return parts[0], parts[1], nil
		}

		if len(parts) == 1 && !needVersion {
			return parts[0], "", nil
		}

		return "", "", fmt.Errorf("%w: verifierID: '%s'",
			serrors.ErrorInvalidFormat, id)
	*/
}
