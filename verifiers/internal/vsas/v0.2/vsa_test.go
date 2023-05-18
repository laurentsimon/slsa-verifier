package v02

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options/vsa"
)

func Test_FromStatement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		statement intoto.Statement
		predicate VsaPredicate
		err       error
	}{
		{
			name: "invalid predicate",
			err:  serrors.ErrorInvalidPredicate,
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa + "A",
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			predicate: VsaPredicate{
				VerificationResult: "PASSED",
				PolicyLevel:        "SLSA_BUILD_LEVEL_2",
			},
		},
		{
			name: "short level",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					// Google bug.
					PolicyLevel: "SLSA_L2",
				},
			},
			predicate: VsaPredicate{
				VerificationResult: "PASSED",
				PolicyLevel:        "SLSA_BUILD_LEVEL_2",
			},
		},
		{
			name: "long level",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			predicate: VsaPredicate{
				VerificationResult: "PASSED",
				PolicyLevel:        "SLSA_BUILD_LEVEL_2",
			},
		},
		{
			name: "invalid level v1.0",
			err:  serrors.ErrorInvalidVsaLevel,
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					// This is v1.0 specs.
					PolicyLevel: "SLSA_BUILD_LEVEL_2",
				},
			},
		},
		{
			name: "invalid level",
			err:  serrors.ErrorInvalidVsaLevel,
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        "INVALID",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			predicate, err := FromStatement(&tt.statement)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(err, tt.err, cmpopts.EquateErrors()))
			}
			if err != nil {
				return
			}
			if !cmp.Equal(*predicate, tt.predicate) {
				t.Errorf(cmp.Diff(*predicate, tt.predicate))
			}
		})
	}
}

func Test_VerifyLevels(t *testing.T) {
	for predicateLevel := vsa.BuildLevel0; predicateLevel <= vsa.BuildLevel3; predicateLevel++ {
		for userLevel := vsa.BuildLevel0; userLevel <= predicateLevel; userLevel++ {
			// Lower levels must pass verification.
			statement := intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        fmt.Sprintf("SLSA_LEVEL_%d", predicateLevel),
				},
			}
			predicate, err := FromStatement(&statement)
			if err != nil {
				t.Errorf("FromStatement: %v", err)
			}
			err = predicate.VerifyLevels([]vsa.Level{userLevel.New()})
			if err != nil {
				t.Errorf(fmt.Sprintf("statement: %v, userLevel: %v, %v", statement, userLevel, cmp.Diff(err, nil, cmpopts.EquateErrors())))
			}

			// Higher levels must fail verification.
			if predicateLevel == userLevel {
				continue
			}

			statement = intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        fmt.Sprintf("SLSA_LEVEL_%d", userLevel),
				},
			}
			predicate, err = FromStatement(&statement)
			if err != nil {
				t.Errorf("FromStatement: %v", err)
			}
			err = predicate.VerifyLevels([]vsa.Level{predicateLevel.New()})
			if !cmp.Equal(err, serrors.ErrorMismatchVsaLevel, cmpopts.EquateErrors()) {
				t.Errorf(fmt.Sprintf("statement: %v, predicateLevel: %v, %v", statement, predicateLevel,
					cmp.Diff(err, serrors.ErrorMismatchVsaLevel, cmpopts.EquateErrors())))
			}
		}
	}
}

func Test_VerifyLevels_invalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		predicateLevel vsa.Level
		levels         []vsa.Level
		err            error
	}{
		{
			name:           "too many different track levels",
			predicateLevel: vsa.BuildLevel1.New(),
			levels:         []vsa.Level{vsa.BuildLevel1.New(), vsa.SourceLevel1.New()},
			err:            serrors.ErrorInvalidVsaLevel,
		},
		{
			name:           "too many same track levels",
			predicateLevel: vsa.BuildLevel1.New(),
			levels:         []vsa.Level{vsa.BuildLevel1.New(), vsa.BuildLevel3.New()},
			err:            serrors.ErrorInvalidVsaLevel,
		},
		{
			name:           "different tracks",
			predicateLevel: vsa.BuildLevel1.New(),
			levels:         []vsa.Level{vsa.SourceLevel1.New()},
			err:            serrors.ErrorMismatchVsaLevel,
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			statement := intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: PredicateVsa,
				},
				Predicate: VsaPredicate{
					VerificationResult: "PASSED",
					PolicyLevel:        fmt.Sprintf("SLSA_LEVEL_%d", tt.predicateLevel.ToInt()),
				},
			}

			predicate, err := FromStatement(&statement)
			if err != nil {
				t.Errorf("FromStatement: %v", err)
			}
			err = predicate.VerifyLevels(tt.levels)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(err, tt.err, cmpopts.EquateErrors()))
			}
		})
	}
}
