package v02

import (
	"encoding/json"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options/vsa"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/vsas/common"
)

const (
	// PredicateVsa represents a VSA for an artifact.
	PredicateVsa = "https://slsa.dev/verification_summary/v0.2"
)

// VsaPredicate is the provenance predicate definition.
type VsaPredicate struct {
	Verifier           common.VsaVerifier `json:"verifier"`
	VerificationResult string             `json:"verification_result"`
	PolicyLevel        string             `json:"policy_level"`
	ResourceURI        string             `json:"resource_uri"`
	// Other fields are not needed for verification.
}

func FromStatement(statement *intoto.Statement) (*VsaPredicate, error) {
	if statement.PredicateType != PredicateVsa {
		return nil, fmt.Errorf("%w: %v", serrors.ErrorInvalidPredicate, statement.PredicateType)
	}

	// Marshal the predicate.
	predicateBytes, err := json.Marshal(statement.Predicate)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	var vsa VsaPredicate
	if err := json.Unmarshal(predicateBytes, &vsa); err != nil {
		return nil, fmt.Errorf("%w: %s", serrors.ErrorInvalidDssePayload, err.Error())
	}

	// Normalize the level.
	switch {
	// SLSA_LX. Bug in Google's VSA.
	case strings.HasPrefix(vsa.PolicyLevel, "SLSA_L") && len(vsa.PolicyLevel) == len("SLSA_L")+1:
		vsa.PolicyLevel = strings.ReplaceAll(vsa.PolicyLevel, "SLSA_L", "SLSA_BUILD_LEVEL_")
	// v0.2 specs.
	case strings.HasPrefix(vsa.PolicyLevel, "SLSA_LEVEL_") && len(vsa.PolicyLevel) == len("SLSA_LEVEL_")+1:
		vsa.PolicyLevel = strings.ReplaceAll(vsa.PolicyLevel, "SLSA_LEVEL_", "SLSA_BUILD_LEVEL_")
	default:
		return nil, fmt.Errorf("%w: %v", serrors.ErrorInvalidVsaLevel, vsa.PolicyLevel)
	}

	return &vsa, nil
}

func (p *VsaPredicate) VerifyLevels(expectedLevels []vsa.Level) error {
	predicateLevel, err := vsa.LevelFromString(p.PolicyLevel)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if len(expectedLevels) != 1 {
		return fmt.Errorf("%w: a single track level is supported, %d provided",
			serrors.ErrorInvalidVsaLevel, len(expectedLevels))
	}
	expectedLevel := expectedLevels[0]
	if expectedLevel.Track() != predicateLevel.Track() {
		return fmt.Errorf("%w: expected '%s' track. Got '%s' track", serrors.ErrorMismatchVsaLevel,
			expectedLevel.Track(), predicateLevel.Track())
	}

	if predicateLevel.LowerThan(expectedLevel) {
		return fmt.Errorf("%w: expected level %s. Got %s", serrors.ErrorMismatchVsaLevel,
			expectedLevel.ToString(), p.PolicyLevel)
	}

	return nil
}
