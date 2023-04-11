package verifiers

import (
	"context"
	"fmt"

	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options"
	"github.com/slsa-framework/slsa-verifier/v2/register"
	_ "github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/vsas/static"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/utils"
)

func getVsaVerifier(verifierOpts *options.VerifierOpts) (register.VsaVerifier, error) {
	// If user provids a verifierID, find the right verifier based on its ID.
	if verifierOpts.ExpectedID != "" {
		name, _, err := utils.ParseVerifierID(verifierOpts.ExpectedID, false)
		if err != nil {
			return nil, err
		}
		for _, v := range register.VsaVerifiers {
			if v.IsAuthoritativeFor(name) {
				return v, nil
			}
		}
	}

	// No verifier found.
	return nil, fmt.Errorf("%w: %s", serrors.ErrorVsaVerifierNotSupported, verifierOpts.ExpectedID)
}

func VerifyImageVsa(ctx context.Context,
	artifactImage string, vsa []byte,
	vsaOpts *options.VsaOpts,
	verifierOpts *options.VerifierOpts,
) ([]byte, *utils.TrustedVerifierID, error) {
	verifier, err := getVsaVerifier(verifierOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyImage(ctx, artifactImage, vsa, vsaOpts, verifierOpts)
}

func VerifyArtifactVsa(ctx context.Context,
	vsa []byte,
	vsaOpts *options.VsaOpts,
	verifierOpts *options.VerifierOpts,
) ([]byte, *utils.TrustedVerifierID, error) {
	verifier, err := getVsaVerifier(verifierOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyArtifact(ctx, vsa,
		vsaOpts, verifierOpts)
}

func VerifyNpmPackageVsa(ctx context.Context,
	attestations []byte,
	vsaOpts *options.VsaOpts,
	verifierOpts *options.VerifierOpts,
) ([]byte, *utils.TrustedVerifierID, error) {
	verifier, err := getVsaVerifier(verifierOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyNpmPackage(ctx, attestations,
		vsaOpts, verifierOpts)
}
