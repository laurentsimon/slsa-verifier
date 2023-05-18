// Copyright 2022 SLSA Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier/verify"
	"github.com/spf13/cobra"
)

const (
	SUCCESS         = "PASSED: Verified SLSA provenance"
	FAILURE         = "FAILED: SLSA verification failed"
	vsaShort        = "Use Verification SLSA Summary (VSA)"
	provenanceShort = "Use SLSA provenance"
)

func verifyArtifactVsaCmd() *cobra.Command {
	o := &verify.VerifyVsaOptions{}
	cmd := &cobra.Command{
		Use:   "vsa",
		Short: vsaShort,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expects one artifact")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vsa hello")
			fmt.Println("levels:", o.VerifiedLevels)
			fmt.Println("verifier.id:", o.VerifierID)

			v := verify.VerifyArtifactVsaCommand{
				VsaPath:        o.VsaPath,
				VerifierID:     o.VerifierID,
				VerifiedLevels: o.VerifiedLevels,
				PrintVsa:       o.PrintVsa,
			}

			if cmd.Flags().Changed("resource-uri") {
				v.ResourceURI = &o.ResourceURI
			}

			if _, err := v.Exec(cmd.Context(), args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", FAILURE, err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", SUCCESS)
			}
		},
	}

	o.AddFlags(cmd)
	cmd.MarkFlagRequired("vsa-path")
	cmd.MarkFlagRequired("verifier-id")
	cmd.MarkFlagRequired("verified-levels")

	return cmd
}

func verifyArtifactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-artifact",
		Short: "Verifies provenance on artifact blobs",
	}

	cmd.AddCommand(verifyArtifactVsaCmd())
	cmd.AddCommand(verifyArtifactProvenanceCmd())
	return cmd
}

func verifyArtifactProvenanceCmd() *cobra.Command {
	o := &verify.VerifyProvenanceOptions{}

	cmd := &cobra.Command{
		Use:   "provenance",
		Short: provenanceShort,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("expects at least one artifact")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			v := verify.VerifyArtifactProvenanceCommand{
				ProvenancePath:      o.ProvenancePath,
				SourceURI:           o.SourceURI,
				PrintProvenance:     o.PrintProvenance,
				BuildWorkflowInputs: o.BuildWorkflowInputs.AsMap(),
			}
			if cmd.Flags().Changed("source-branch") {
				v.SourceBranch = &o.SourceBranch
			}
			if cmd.Flags().Changed("source-tag") {
				v.SourceTag = &o.SourceTag
			}
			if cmd.Flags().Changed("source-versioned-tag") {
				v.SourceVersionTag = &o.SourceVersionTag
			}
			if cmd.Flags().Changed("builder-id") {
				v.BuilderID = &o.BuilderID
			}

			if _, err := v.Exec(cmd.Context(), args); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", FAILURE, err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", SUCCESS)
			}
		},
	}

	o.AddFlags(cmd)
	// --provenance-path must be supplied when verifying an artifact.
	cmd.MarkFlagRequired("provenance-path")
	return cmd
}

func verifyImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-image",
		Short: "Verifies provenance on a container image",
	}

	cmd.AddCommand(verifyImageProvenanceCmd())
	return cmd
}

func verifyImageProvenanceCmd() *cobra.Command {
	o := &verify.VerifyProvenanceOptions{}

	cmd := &cobra.Command{
		Use:   "provenance",
		Short: provenanceShort,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expects a single path to an image")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			v := verify.VerifyImageProvenanceCommand{
				SourceURI:           o.SourceURI,
				PrintProvenance:     o.PrintProvenance,
				BuildWorkflowInputs: o.BuildWorkflowInputs.AsMap(),
			}
			if cmd.Flags().Changed("provenance-path") {
				v.ProvenancePath = &o.ProvenancePath
			}
			if cmd.Flags().Changed("source-branch") {
				v.SourceBranch = &o.SourceBranch
			}
			if cmd.Flags().Changed("source-tag") {
				v.SourceTag = &o.SourceTag
			}
			if cmd.Flags().Changed("source-versioned-tag") {
				v.SourceVersionTag = &o.SourceVersionTag
			}
			if cmd.Flags().Changed("builder-id") {
				v.BuilderID = &o.BuilderID
			}

			if _, err := v.Exec(cmd.Context(), args); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", FAILURE, err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", SUCCESS)
			}
		},
	}

	o.AddFlags(cmd)
	return cmd
}

func verifyNpmPackageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-npm-package",
		Short: "Verifies provenance on an npm package",
	}

	// NOTE: We can later add support for
	// VSA, publish, provenance.
	cmd.AddCommand(verifyNpmAttestationsCmd())
	return cmd
}

func verifyNpmAttestationsCmd() *cobra.Command {
	o := &verify.VerifyNpmAttestationsOptions{}

	cmd := &cobra.Command{
		Use: "attestations",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expects a single path to an image")
			}
			return nil
		},
		Short: "Uses attestations from the registry",
		Run: func(cmd *cobra.Command, args []string) {
			v := verify.VerifyNpmPackageProvenanceCommand{
				SourceURI:           o.SourceURI,
				PrintProvenance:     o.PrintProvenance,
				BuildWorkflowInputs: o.BuildWorkflowInputs.AsMap(),
			}
			if cmd.Flags().Changed("attestations-path") {
				v.AttestationsPath = o.AttestationsPath
			}
			if cmd.Flags().Changed("package-name") {
				v.PackageName = &o.PackageName
			}
			if cmd.Flags().Changed("package-version") {
				v.PackageVersion = &o.PackageVersion
			}
			if cmd.Flags().Changed("source-branch") {
				fmt.Fprintf(os.Stderr, "%s: --source-branch not supported\n", FAILURE)
				os.Exit(1)
			}
			if cmd.Flags().Changed("source-tag") {
				fmt.Fprintf(os.Stderr, "%s: --source-tag not supported\n", FAILURE)
				os.Exit(1)
			}
			if cmd.Flags().Changed("source-versioned-tag") {
				fmt.Fprintf(os.Stderr, "%s: --source-versioned-tag not supported\n", FAILURE)
				os.Exit(1)
			}
			if cmd.Flags().Changed("print-provenance") {
				fmt.Fprintf(os.Stderr, "%s: --print-provenance not supported\n", FAILURE)
				os.Exit(1)
			}
			if cmd.Flags().Changed("builder-id") {
				v.BuilderID = &o.BuilderID
			}

			if _, err := v.Exec(cmd.Context(), args); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", FAILURE, err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", SUCCESS)
			}
		},
	}

	o.AddFlags(cmd)
	return cmd
}
