package vsas

import (
	"context"
	"fmt"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	dsselib "github.com/secure-systems-lab/go-securesystemslib/dsse"
	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options"
	vsa02 "github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/vsas/v0.2"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/utils"
)

type VsaVerifier struct {
	verifier *utils.DsseVerifier
	envelope *dsselib.Envelope
}

func VsaVerifierNew(vsa []byte, verifier *utils.DsseVerifier) (*VsaVerifier, error) {
	env, err := utils.EnvelopeFromBytes(vsa)
	if err != nil {
		return nil, err
	}

	return &VsaVerifier{
		verifier: verifier,
		envelope: env,
	}, nil
}

func (v *VsaVerifier) Verify(ctx context.Context, vsaOpts *options.VsaOpts,
	verifierOpts *options.VerifierOpts, encoding utils.SignatureEncoding,
) (*utils.TrustedVerifierID, error) {
	// Verify the signature.
	if err := v.verifySignature(ctx, encoding); err != nil {
		return nil, err
	}

	// Extract the statement.
	statement, err := utils.StatementFromEnvelope(v.envelope)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Verify the payload.
	return v.verifyPayload(statement, vsaOpts, verifierOpts)
}

func (v *VsaVerifier) verifySignature(ctx context.Context, encoding utils.SignatureEncoding) error {
	return v.verifier.Verify(ctx, v.envelope, &encoding)
}

func (v *VsaVerifier) verifyPayload(statement *intoto.Statement, vsaOpts *options.VsaOpts,
	verifierOpts *options.VerifierOpts,
) (*utils.TrustedVerifierID, error) {
	// Extract the VSA predicate.
	vsaPredicate, err := vsa02.FromStatement(statement)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	trustedVerifier, err := utils.TrustedVerifierIDNew(vsaPredicate.Verifier.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := trustedVerifier.Matches(verifierOpts.ExpectedID, false); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Verify the hash.
	if err := v.verifySubject(statement, vsaOpts.ExpectedDigest); err != nil {
		return nil, err
	}

	// Verify verification result.
	if vsaPredicate.VerificationResult != "PASSED" {
		return nil, fmt.Errorf("%w: verification result is '%s'", serrors.ErrorVsaResultFailure,
			vsaPredicate.VerificationResult)
	}

	// Verify policy levels.
	if err := vsaPredicate.VerifyLevels(vsaOpts.ExpectedLevels); err != nil {
		return nil, err
	}

	if vsaOpts.ExpectedResourceURI != nil &&
		*vsaOpts.ExpectedResourceURI != vsaPredicate.ResourceURI {
		return nil, fmt.Errorf("%w: expected '%s'. Got '%s'",
			serrors.ErrorMismatchVsaResourceURI, *vsaOpts.ExpectedResourceURI,
			vsaPredicate.ResourceURI)
	}

	return trustedVerifier, nil
}

func (v *VsaVerifier) verifySubject(statement *intoto.Statement, expectedHash string) error {
	for _, subject := range statement.StatementHeader.Subject {
		digestSet := subject.Digest
		hash, exists := digestSet["sha256"]
		if !exists {
			return fmt.Errorf("%w: %s", serrors.ErrorInvalidDssePayload, "no sha256 subject digest")
		}

		if hash == expectedHash {
			return nil
		}
	}

	return fmt.Errorf("expected hash '%s' not found: %w", expectedHash, serrors.ErrorMismatchHash)
}
