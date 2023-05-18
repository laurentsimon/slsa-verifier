package register

import (
	"context"

	"github.com/slsa-framework/slsa-verifier/v2/options"
	"github.com/slsa-framework/slsa-verifier/v2/verifiers/utils"
)

var ProvenanceVerifiers = make(map[string]ProvenanceVerifier)

type ProvenanceVerifier interface {
	// IsAuthoritativeFor checks whether a verifier can
	// verify provenance for a given builder identified by its
	// `BuilderID`.
	IsAuthoritativeFor(builderIDName string) bool

	// VerifyArtifact verifies a provenance for a supplied artifact.
	VerifyArtifact(ctx context.Context,
		provenance []byte, artifactHash string,
		provenanceOpts *options.ProvenanceOpts,
		builderOpts *options.BuilderOpts,
	) ([]byte, *utils.TrustedBuilderID, error)

	// VerifyImage verifies a provenance for a supplied OCI image.
	VerifyImage(ctx context.Context,
		provenance []byte, artifactImage string,
		provenanceOpts *options.ProvenanceOpts,
		builderOpts *options.BuilderOpts,
	) ([]byte, *utils.TrustedBuilderID, error)

	VerifyNpmPackage(ctx context.Context,
		attestations []byte, tarballHash string,
		provenanceOpts *options.ProvenanceOpts,
		builderOpts *options.BuilderOpts,
	) ([]byte, *utils.TrustedBuilderID, error)
}

func RegisterProvenanceVerifier(name string, verifier ProvenanceVerifier) {
	ProvenanceVerifiers[name] = verifier
}

var VsaVerifiers = make(map[string]VsaVerifier)

type VsaVerifier interface {
	// IsAuthoritativeFor checks whether a verifier can
	// verify VSA for a given verifier identified by its
	// `VerifierID`.
	IsAuthoritativeFor(verifierID string) bool

	// VerifyArtifact verifies a VSA for a supplied artifact.
	VerifyArtifact(ctx context.Context,
		vsa []byte,
		vsaOpts *options.VsaOpts,
		verifierOpts *options.VerifierOpts,
	) ([]byte, *utils.TrustedVerifierID, error)

	// VerifyImage verifies a VSA for a supplied OCI image.
	VerifyImage(ctx context.Context,
		artifactImage string, vsa []byte,
		vsaOpts *options.VsaOpts,
		verifierOpts *options.VerifierOpts,
	) ([]byte, *utils.TrustedVerifierID, error)

	VerifyNpmPackage(ctx context.Context,
		vsa []byte,
		vsaOpts *options.VsaOpts,
		verifierOpts *options.VerifierOpts,
	) ([]byte, *utils.TrustedVerifierID, error)
}

func RegisterVsaVerifier(name string, verifier VsaVerifier) {
	VsaVerifiers[name] = verifier
}
