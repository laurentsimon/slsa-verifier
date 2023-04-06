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

package verify

import (
	"fmt"
	"strconv"
	"strings"

	serrors "github.com/slsa-framework/slsa-verifier/v2/errors"
	"github.com/spf13/cobra"
)

type Interface interface {
	// AddFlags adds this options' flags to the cobra command.
	AddFlags(cmd *cobra.Command)
}

// VerifyOptions is the top-level options for all `provenance` commands.
type VerifyOptions struct {
	/* Source requirements */
	SourceURI        string
	SourceBranch     string
	SourceTag        string
	SourceVersionTag string
	/* Builder Requirements */
	BuildWorkflowInputs workflowInputs
	BuilderID           string
	/* Other */
	ProvenancePath  string
	PrintProvenance bool
}

var _ Interface = (*VerifyOptions)(nil)

// AddFlags implements Interface.
func (o *VerifyOptions) AddFlags(cmd *cobra.Command) {
	/* Builder options */
	cmd.Flags().Var(&o.BuildWorkflowInputs, "build-workflow-input",
		"[optional] a workflow input provided by a user at trigger time in the format 'key=value'. (Only for 'workflow_dispatch' events on GitHub Actions).")

	cmd.Flags().StringVar(&o.BuilderID, "builder-id", "", "[optional] the unique builder ID who created the provenance")

	/* Source options */
	cmd.Flags().StringVar(&o.SourceURI, "source-uri", "",
		"expected source repository that should have produced the binary, e.g. github.com/some/repo")

	cmd.Flags().StringVar(&o.SourceBranch, "source-branch", "", "[optional] expected branch the binary was compiled from")

	cmd.Flags().StringVar(&o.SourceTag, "source-tag", "", "[optional] expected tag the binary was compiled from")

	cmd.Flags().StringVar(&o.SourceVersionTag, "source-versioned-tag", "",
		"[optional] expected version the binary was compiled from. Uses semantic version to match the tag")

	/* Other options */
	cmd.Flags().StringVar(&o.ProvenancePath, "provenance-path", "",
		"path to a provenance file")

	cmd.Flags().BoolVar(&o.PrintProvenance, "print-provenance", false,
		"[optional] print the verified provenance to stdout")

	cmd.MarkFlagRequired("source-uri")
	cmd.MarkFlagsMutuallyExclusive("source-versioned-tag", "source-tag")
}

type VSABuildLevel int

const (
	VSABuildLevel0 VSABuildLevel = iota
	VSABuildLevel1
	VSABuildLevel2
	VSABuildLevel3
)

type VSASourceLevel int

const (
	VSASourceLevel0 VSASourceLevel = iota
	VSASourceLevel1
	VSASourceLevel2
	VSASourceLevel3
)

type VSALevels []interface{}

// VerifyVSAOptions is the top-level options for all `vsa` commands.
type VerifyVSAOptions struct {
	VSAPath        string
	VerifierID     string
	VerifiedLevels VSALevels
	ResourceURI    string
}

var _ Interface = (*VerifyVSAOptions)(nil)

// AddFlags implements Interface.
func (o *VerifyVSAOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.VSAPath, "vsa-path", "", "the path to a VSA")
	cmd.Flags().StringVar(&o.VerifierID, "verifier-id", "", "the unique verifier ID who created the VSA")
	cmd.Flags().StringVar(&o.ResourceURI, "resource-uri", "", "the resource URI the VSA is for")
	cmd.Flags().Var(&o.VerifiedLevels, "verified-levels", "verified levels of the artifacts")
}

func (ls *VSALevels) Type() string {
	return "SLSA_BUILD_LEVEL_1, SLSA_SOURCE_LEVEL_1, ..."
}

func (ls *VSALevels) String() string {
	return fmt.Sprintf("%v", *ls)
}

func (ls *VSALevels) Set(value string) error {
	parts := strings.Split(value, ",")
	levels := make([]interface{}, len(parts))
	for i := range parts {
		levelStr := strings.TrimSpace(parts[i])
		// Build level.
		blf := parseLevel[VSABuildLevel]
		bl, berr := blf(levelStr, "SLSA_BUILD_LEVEL_")
		if berr == nil {
			levels[i] = bl
			continue
		}

		// Source level.
		slf := parseLevel[VSABuildLevel]
		sl, serr := slf(levelStr, "SLSA_SOURCE_LEVEL_")
		if serr == nil {
			levels[i] = sl
			continue
		}
		return fmt.Errorf("parse: [build track:%v, source track:%v]", serr, berr)

	}
	*ls = levels
	return nil
}

type levelType interface {
	VSABuildLevel | VSASourceLevel
}

func parseLevel[T levelType](levelStr, prefix string) (*T, error) {
	if !strings.HasPrefix(levelStr, prefix) ||
		len(levelStr) != len(prefix)+1 {
		return nil, fmt.Errorf("invalid level: %s", levelStr)
	}
	level, err := strconv.Atoi(levelStr[len(levelStr)-1:])
	if err != nil {
		return nil, fmt.Errorf("%w: not a number: %q", err, levelStr)
	}
	if level < 0 || level > 4 {
		return nil, fmt.Errorf("invalid level: %q", levelStr)
	}
	ret := T(level)
	return &ret, nil
}

// VerifyNpmOptions is the top-level options for the `verifyNpmPackage` command.
type VerifyNpmOptions struct {
	VerifyOptions
	/* Other */
	AttestationsPath string
	PackageName      string
	PackageVersion   string
}

var _ Interface = (*VerifyNpmOptions)(nil)

// AddFlags implements Interface.
func (o *VerifyNpmOptions) AddFlags(cmd *cobra.Command) {
	/* Builder options */
	cmd.Flags().Var(&o.BuildWorkflowInputs, "build-workflow-input",
		"[optional] a workflow input provided by a user at trigger time in the format 'key=value'. (Only for 'workflow_dispatch' events on GitHub Actions).")

	cmd.Flags().StringVar(&o.BuilderID, "builder-id", "", "[optional] the unique builder ID who created the provenance")

	/* Source options */
	cmd.Flags().StringVar(&o.SourceURI, "source-uri", "",
		"expected source repository that should have produced the binary, e.g. github.com/some/repo")

	cmd.Flags().StringVar(&o.SourceBranch, "source-branch", "", "[optional] expected branch the binary was compiled from")

	cmd.Flags().StringVar(&o.SourceTag, "source-tag", "", "[optional] expected tag the binary was compiled from")

	cmd.Flags().StringVar(&o.SourceVersionTag, "source-versioned-tag", "",
		"[optional] expected version the binary was compiled from. Uses semantic version to match the tag")

	cmd.Flags().StringVar(&o.AttestationsPath, "attestations-path", "",
		"path to a file containing the attestations")

	cmd.Flags().StringVar(&o.PackageName, "package-name", "",
		"the package name")

	cmd.Flags().StringVar(&o.PackageVersion, "package-version", "",
		"the package version")

	cmd.Flags().BoolVar(&o.PrintProvenance, "print-provenance", false,
		"[optional] print the verified provenance to stdout")

	cmd.MarkFlagRequired("source-uri")
	cmd.MarkFlagRequired("builder-id")
	cmd.MarkFlagRequired("package-name")
	cmd.MarkFlagRequired("package-version")
	cmd.MarkFlagsMutuallyExclusive("source-versioned-tag", "source-tag")
}

type workflowInputs struct {
	kv map[string]string
}

func (i *workflowInputs) Type() string {
	return fmt.Sprintf("%v", i.kv)
}

func (i *workflowInputs) String() string {
	return fmt.Sprintf("%v", i.kv)
}

func (i *workflowInputs) Set(value string) error {
	l := strings.Split(value, "=")
	if len(l) != 2 {
		return fmt.Errorf("%w: expected 'key=value' format, got '%s'", serrors.ErrorInvalidFormat, value)
	}
	if i.kv == nil {
		i.kv = make(map[string]string)
	}
	i.kv[l[0]] = l[1]
	return nil
}

func (i *workflowInputs) AsMap() map[string]string {
	return i.kv
}
