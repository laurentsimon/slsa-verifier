package vsas

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	intotocommon "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options"
	"github.com/slsa-framework/slsa-verifier/v2/options/vsa"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/vsas/common"
	vsa02 "github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/vsas/v0.2"
)

func Test_verifySubject(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		statement intoto.Statement
		digest    string
		err       error
	}{
		{
			name:   "sha256 match",
			digest: "123",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
			},
		},
		{
			name:   "sha256 mismatch",
			digest: "1234",
			err:    serrors.ErrorMismatchHash,
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
			},
		},
		{
			name:   "no sha256 digest",
			digest: "1234",
			err:    serrors.ErrorInvalidDssePayload,
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha512": "123",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := VsaVerifier{}
			err := v.verifySubject(&tt.statement, tt.digest)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(err, tt.err, cmpopts.EquateErrors()))
			}
		})
	}
}

func Test_verifyPayload_v02(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		statement    intoto.Statement
		vsaOpts      options.VsaOpts
		verifierOpts options.VerifierOpts
		err          error
	}{
		{
			name: "valid full builder ID",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "123",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel2.New()},
			},
		},
		{
			name: "valid short builder ID",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "123",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel2.New()},
			},
		},
		{
			name: "valid level 1",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "123",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel1.New()},
			},
		},
		{
			name: "invalid level 3",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "123",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel3.New()},
			},
			err: serrors.ErrorMismatchVsaLevel,
		},
		{
			name: "digest mismatch",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "1234",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "1235",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel2.New()},
			},
			err: serrors.ErrorMismatchHash,
		},
		{
			name: "verification fail",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "FAILED",
					PolicyLevel:        "SLSA_LEVEL_2",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest: "123",
				ExpectedLevels: []vsa.Level{vsa.BuildLevel2.New()},
			},
			err: serrors.ErrorVsaResultFailure,
		},
		{
			name: "match resource uri",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
					ResourceURI:        "name://the-resource",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest:      "123",
				ExpectedLevels:      []vsa.Level{vsa.BuildLevel2.New()},
				ExpectedResourceURI: asPtr("name://the-resource"),
			},
		},
		{
			name: "mismatch resource uri",
			statement: intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        "SLSA_LEVEL_2",
					ResourceURI:        "name://the-resource-different",
				},
			},
			verifierOpts: options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			},
			vsaOpts: options.VsaOpts{
				ExpectedDigest:      "123",
				ExpectedLevels:      []vsa.Level{vsa.BuildLevel2.New()},
				ExpectedResourceURI: asPtr("name://the-resource"),
			},
			err: serrors.ErrorMismatchVsaResourceURI,
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := VsaVerifier{}
			verifierOutID, err := v.verifyPayload(&tt.statement, &tt.vsaOpts,
				&tt.verifierOpts)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(err, tt.err, cmpopts.EquateErrors()))
			}
			if err != nil {
				return
			}

			predicate, ok := tt.statement.Predicate.(vsa02.VsaPredicate)
			if !ok {
				t.Errorf("predicate not a v02 VSA predicate")
			}
			if verifierOutID.String() != predicate.Verifier.ID {
				t.Errorf(cmp.Diff(verifierOutID.String(), predicate.Verifier.ID, cmpopts.EquateErrors()))
			}
		})
	}
}

func asPtr(s string) *string {
	return &s
}

func Test_verifyPayload_v02_levels(t *testing.T) {
	for predicateLevel := vsa.BuildLevel0; predicateLevel < vsa.BuildLevel3; predicateLevel++ {
		for userLevel := vsa.BuildLevel0; userLevel <= predicateLevel; userLevel++ {
			// Lower levels must pass verification.
			statement := intoto.Statement{
				StatementHeader: intoto.StatementHeader{
					PredicateType: vsa02.PredicateVsa,
					Subject: []intoto.Subject{
						{
							Digest: intotocommon.DigestSet{
								"sha256": "123",
							},
						},
					},
				},
				Predicate: vsa02.VsaPredicate{
					Verifier: common.VsaVerifier{
						ID: "https://the.trusted.verifier/v0.1",
					},
					VerificationResult: "PASSED",
					PolicyLevel:        fmt.Sprintf("SLSA_LEVEL_%d", predicateLevel),
					ResourceURI:        "name://the-resource",
				},
			}
			verifierOpts := options.VerifierOpts{
				ExpectedID: "https://the.trusted.verifier/v0.1",
			}
			vsaOpts := options.VsaOpts{
				ExpectedDigest:      "123",
				ExpectedLevels:      []vsa.Level{userLevel.New()},
				ExpectedResourceURI: asPtr("name://the-resource"),
			}

			v := VsaVerifier{}
			verifierOutID, err := v.verifyPayload(&statement, &vsaOpts, &verifierOpts)
			if err != nil {
				t.Errorf(fmt.Sprintf("statement: %v, userLevel: %v, %v", statement,
					userLevel, cmp.Diff(err, nil, cmpopts.EquateErrors())))
			}

			predicate, err := vsa02.FromStatement(&statement)
			if err != nil {
				t.Errorf("FromStatement: %v", err)
			}
			if verifierOutID.String() != predicate.Verifier.ID {
				t.Errorf(cmp.Diff(verifierOutID.String(), predicate.Verifier.ID, cmpopts.EquateErrors()))
			}

			// Higher levels must fail verification.
			if predicateLevel == userLevel {
				continue
			}

			statement.Predicate = vsa02.VsaPredicate{
				Verifier: common.VsaVerifier{
					ID: "https://the.trusted.verifier/v0.1",
				},
				VerificationResult: "PASSED",
				PolicyLevel:        fmt.Sprintf("SLSA_LEVEL_%d", userLevel),
				ResourceURI:        "name://the-resource",
			}
			vsaOpts.ExpectedLevels = []vsa.Level{predicateLevel.New()}

			_, err = v.verifyPayload(&statement, &vsaOpts, &verifierOpts)
			if err == nil {
				t.Errorf(fmt.Sprintf("statement: %v, predicateLevel: %v, %v", statement, predicateLevel,
					cmp.Diff(err, serrors.ErrorMismatchVsaLevel, cmpopts.EquateErrors())))
			}
		}
	}
}
