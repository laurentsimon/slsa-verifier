package verifiers

import (
	"context"
	"fmt"

	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/slsa-framework/slsa-verifier/v2/options"
	"github.com/slsa-framework/slsa-verifier/v2/register"
	_ "github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/builders/gcb"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/internal/builders/gha"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/utils"
)

func getProvenanceVerifier(builderOpts *options.BuilderOpts) (register.ProvenanceVerifier, error) {
	// By default, use the GHA builders
	verifier := register.ProvenanceVerifiers[gha.BuilderName]

	// If user provids a builderID, find the right verifier based on its ID.
	if builderOpts.ExpectedID != nil &&
		*builderOpts.ExpectedID != "" {
		name, _, err := utils.ParseBuilderID(*builderOpts.ExpectedID, false)
		if err != nil {
			return nil, err
		}
		for _, v := range register.ProvenanceVerifiers {
			if v.IsAuthoritativeFor(name) {
				return v, nil
			}
		}
		// No builder found.
		return nil, fmt.Errorf("%w: %s", serrors.ErrorBuilderVerifierNotSupported, *builderOpts.ExpectedID)
	}

	return verifier, nil
}

func VerifyImageProvenance(ctx context.Context, artifactImage string,
	provenance []byte,
	provenanceOpts *options.ProvenanceOpts,
	builderOpts *options.BuilderOpts,
) ([]byte, *utils.TrustedBuilderID, error) {
	verifier, err := getProvenanceVerifier(builderOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyImage(ctx, provenance, artifactImage, provenanceOpts, builderOpts)
}

func VerifyArtifactProvenance(ctx context.Context,
	provenance []byte, artifactHash string,
	provenanceOpts *options.ProvenanceOpts,
	builderOpts *options.BuilderOpts,
) ([]byte, *utils.TrustedBuilderID, error) {
	verifier, err := getProvenanceVerifier(builderOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyArtifact(ctx, provenance, artifactHash,
		provenanceOpts, builderOpts)
}

func VerifyNpmPackageProvenance(ctx context.Context,
	attestations []byte, tarballHash string,
	provenanceOpts *options.ProvenanceOpts,
	builderOpts *options.BuilderOpts,
) ([]byte, *utils.TrustedBuilderID, error) {
	verifier, err := getProvenanceVerifier(builderOpts)
	if err != nil {
		return nil, nil, err
	}

	return verifier.VerifyNpmPackage(ctx, attestations, tarballHash,
		provenanceOpts, builderOpts)
}
