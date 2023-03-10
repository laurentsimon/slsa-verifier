# Verification of SLSA provenance

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/slsa-framework/slsa-verifier/badge)](https://api.securityscorecards.dev/projects/github.com/slsa-framework/slsa-verifier)
[![OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/projects/6729/badge)](https://bestpractices.coreinfrastructure.org/projects/6729)
[![Go Report Card](https://goreportcard.com/badge/github.com/slsa-framework/slsa-verifier)](https://goreportcard.com/report/github.com/slsa-framework/slsa-verifier)
[![Slack](https://slack.babeljs.io/badge.svg)](https://slack.com/app_redirect?team=T019QHUBYQ3&channel=slsa-tooling)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)

<img align="right" src="https://slsa.dev/images/logo-mono.svg" width="140" height="140">

<!-- markdown-toc --bullets="-" -i README.md -->

<!-- toc -->

- [Overview](#overview)
  - [What is SLSA?](#what-is-slsa)
  - [What is provenance?](#what-is-provenance)
  - [What is slsa-verifier?](#what-is-slsa-verifier)
- [Installation](#installation)
  - [Compilation from source](#compilation-from-source)
    - [Option 1: Install via go](#option-1-install-via-go)
    - [Option 2: Compile manually](#option-2-compile-manually)
    - [Option 3: Use the installer Action](#option-3-use-the-installer-action)
  - [Download the binary](#download-the-binary)
- [Available options](#available-options)
- [Option list](#option-list)
  - [Option details](#option-details)
- [Verification for GitHub builders](#verification-for-github-builders)
  - [Artifacts](#artifacts)
  - [Containers](#containers)
- [Verification for Google Cloud Build](#verification-for-google-cloud-build)
  - [Artifacts](#artifacts-1)
  - [Containers](#containers-1)
- [Known Issues](#known-issues)
  - [tuf: invalid key](#tuf-invalid-key)
  - [panic: assignment to entry in nil map](#panic-assignment-to-entry-in-nil-map)
- [Technical design](#technical-design)
  - [Blog post](#blog-post)
  - [Specifications](#specifications)
  - [TOCTOU attacks](#toctou-attacks)

<!-- tocstop -->

## Overview

### What is SLSA?

[Supply chain Levels for Software Artifacts](https://slsa.dev), or SLSA (salsa),
is a security framework, a check-list of standards and controls to prevent
tampering, improve integrity, and secure packages and infrastructure in your
projects, businesses or enterprises.

SLSA defines an incrementially adoptable set of levels which are defined in
terms of increasing compliance and assurance. SLSA levels are like a common
language to talk about how secure software, supply chains and their component
parts really are.

### What is provenance?

Provenance is information, or metadata, about how a software artifact was
created. This could include information about what source code, build system,
and build steps were used, as well as who and why the build was initiated.
Provenance can be used to determine the authenticity and trustworthiness of
software artifacts that you use.

As part of the framework, SLSA defines a
[provenance format](https://slsa.dev/provenance/) which can be used hold this
metadata.

### What is slsa-verifier?

slsa-verifier is a tool for verifying
[SLSA provenance](https://slsa.dev/provenance/) that was generated by CI/CD
builders. slsa-verifier verifies the provenance by verifying the cryptographic
signatures on provenance to make sure it was created by the expected builder.
It then verifies that various values such as the builder id, source code
repository, ref (branch or tag) matches the expected values.

It currently supports verifying provenance generated by:

1. [SLSA generator](https://github.com/slsa-framework/slsa-github-generator)
1. [Google Cloud Build (GCB)](https://cloud.google.com/build/docs/securing-builds/view-build-provenance).

---

[Installation](#installation)

- [Compilation from source](#compilation-from-source)
- [Download the binary](#download-the-binary)

[Available options](#available-options)

- [Option list](#option-list)
- [Option details](#option-details)

[Verification for GitHub builders](#verification-for-github-builders)

- [Artifacts](#artifacts)
- [Containers](#containers)

[Verification for Google Cloud Build](#verification-for-google-cloud-build)

- [Artifacts](#artifacts-1)
- [Containers](#containers-1)

[Known Issues](#known-issues)

[Technical design](#technial-design)

- [Blog posts](#blog-posts)
- [Specifications](#specifications)
- [TOCTOU attacks](#toctou-attacks)

---

## Installation

You have two options to install the verifier.

### Compilation from source

#### Option 1: Install via go

If you want to install the verifier, you can run the following command:
```bash
$ go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.0.1
$ slsa-verifier <options>
```

Tools like [dependabot](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuring-dependabot-version-updates) or [renovate](https://github.com/renovatebot/renovate) use your project's go.mod to identify the version of your dependencies. 
If you install the verifier in CI, we strongly recommend you follow the steps below to keep the verifier up-to-date:

1. Create a tooling/tooling_test.go file containing the following:
```go
//go:build tools
// +build tools

package main

import (
	_ "github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier"
)
```

1. Run the following commands in the tooling directory. (It will create a go.sum file.)
```bash
$ go mod init <your-project-name>-tooling
$ go mod tidy
```

1. Commit the tooling folder (containing the 3 files tooling_test.go, go.mod and go.sum) to the repository.
1. To install the verifier in your CI, run the following commands:
```bash
$ cd tooling
$ grep _ tooling_test.go | cut -f2 -d '"' | xargs -n1 -t go install
``` 

#### Option 2: Compile manually

```bash
$ git clone git@github.com:slsa-framework/slsa-verifier.git
$ cd slsa-verifier && git checkout v2.0.1
$ go run ./cli/slsa-verifier <options>
```

#### Option 3: Use the installer Action

If you need to install the verifier to run in a GitHub workflow, use the installer Action as described in [actions/installer/README.md](./actions/installer/README.md).

### Download the binary

Download the binary from the latest release at [https://github.com/slsa-framework/slsa-verifier/releases/tag/v2.0.1](https://github.com/slsa-framework/slsa-verifier/releases/tag/v2.0.1)

Download the [SHA256SUM.md](https://github.com/slsa-framework/slsa-verifier/blob/main/SHA256SUM.md).

Verify the checksum:

```bash
$ sha256sum -c --strict SHA256SUM.md
  slsa-verifier-linux-amd64: OK
```

## Available options

We currently support artifact verification (for binary blobs) and container images.

## Option list

Below is a list of options currently supported for binary blobs and container images. Note that signature verification is handled seamlessly without the need for developers to manipulate public keys. See [Available options](#available-options) for details on the options exposed to validate the provenance.

```bash
$ git clone git@github.com:slsa-framework/slsa-verifier.git
$ go run ./cli/slsa-verifier/ verify-artifact --help
Verifies SLSA provenance on artifact blobs given as arguments (assuming same provenance)

Usage:
  slsa-verifier verify-artifact [flags] artifact [artifact..]

Flags:
      --build-workflow-input map[]    [optional] a workflow input provided by a user at trigger time in the format 'key=value'. (Only for 'workflow_dispatch' events on GitHub Actions). (default map[])
      --builder-id string             [optional] the unique builder ID who created the provenance
  -h, --help                          help for verify-artifact
      --print-provenance              [optional] print the verified provenance to stdout
      --provenance-path string        path to a provenance file
      --source-branch string          [optional] expected branch the binary was compiled from
      --source-tag string             [optional] expected tag the binary was compiled from
      --source-uri string             expected source repository that should have produced the binary, e.g. github.com/some/repo
      --source-versioned-tag string   [optional] expected version the binary was compiled from. Uses semantic version to match the tag
```

Multiple artifacts can be passed to `verify-artifact`. As long as they are all covered by the same provenance file, the verification will succeed.

### Option details

The following options are available:

| Option                 | Description                                                                                                                                                                                                                                                                                                                                                                                               | Support                                                                                             |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `source-uri`           | Expects a source, for e.g. `github.com/org/repo`.                                                                                                                                                                                                                                                                                                                                                         | All builders                                                                                        |
| `source-branch`        | Expects a `branch` like `main` or `dev`. Not supported for all GitHub Workflow triggers.                                                                                                                                                                                                                                                                                                                  | [GitHub builders](https://github.com/slsa-framework/slsa-github-generator#generation-of-provenance) |
| `source-tag`           | Expects a `tag` like `v0.0.1`. Verifies exact tag used to create the binary. Supported for new [tag](https://github.com/slsa-framework/example-package/blob/main/.github/workflows/e2e.go.tag.main.config-ldflags-assets-tag.slsa3.yml#L5) and [release](https://github.com/slsa-framework/example-package/blob/main/.github/workflows/e2e.go.release.main.config-ldflags-assets-tag.slsa3.yml) triggers. | [GitHub builders](https://github.com/slsa-framework/slsa-github-generator#generation-of-provenance) |
| `source-versioned-tag` | Like `tag`, but verifies using semantic versioning.                                                                                                                                                                                                                                                                                                                                                       | [GitHub builders](https://github.com/slsa-framework/slsa-github-generator#generation-of-provenance) |
| `build-workflow-input` | Expects key-value pairs like `key=value` to match against [inputs](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#onworkflow_dispatchinputs) for GitHub Actions `workflow_dispatch` triggers.                                                                                                                                                                      | [GitHub builders](https://github.com/slsa-framework/slsa-github-generator#generation-of-provenance) |

## Verification for GitHub builders

### Artifacts

To verify an artifact, run the following command:

```bash
$ slsa-verifier verify-artifact slsa-test-linux-amd64 \
  --provenance-path slsa-test-linux-amd64.intoto.jsonl \
  --source-uri github.com/slsa-framework/slsa-test \
  --source-tag v1.0.3
Verified signature against tlog entry index 3189970 at URL: https://rekor.sigstore.dev/api/v1/log/entries/206071d5ca7a2346e4db4dcb19a648c7f13b4957e655f4382b735894059bd199
Verified build using builder https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@refs/tags/v1.2.0 at commit 5bb13ef508b2b8ded49f9264d7712f1316830d10
PASSED: Verified SLSA provenance
```

The verified in-toto statement may be written to stdout with the `--print-provenance` flag to pipe into policy engines.

Only GitHub URIs are supported with the `--source-uri` flag. A tag should not be specified, even if the provenance was built at some tag. If you intend to do source versioning validation, use `--print-provenance` and inspect the commit SHA of the config source or materials.

Multiple artifacts built from the same GitHub builder can be verified in the same command, by passing them in the same command line as arguments:

```bash
$ slsa-verifier verify-artifact \
  --provenance-path /tmp/demo/multiple.intoto.jsonl \
  --source-uri github.com/mihaimaruseac/example \
  /tmp/demo/fib /tmp/demo/hello

Verified signature against tlog entry index 9712459 at URL: https://rekor.sigstore.dev/api/v1/log/entries/24296fb24b8ad77a1544828b67bb5a2335f7e0d01c504a32ceb6f3a8814ed12c8f1b222d308bd9e8
Verified build using builder https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v1.4.0 at commit 11fab87c5ee6f46c6f5e68f6c5378c62ce1ca77c
Verifying artifact /tmp/demo/fib: PASSED

Verified signature against tlog entry index 9712459 at URL: https://rekor.sigstore.dev/api/v1/log/entries/24296fb24b8ad77a1544828b67bb5a2335f7e0d01c504a32ceb6f3a8814ed12c8f1b222d308bd9e8
Verified build using builder https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v1.4.0 at commit 11fab87c5ee6f46c6f5e68f6c5378c62ce1ca77c
Verifying artifact /tmp/demo/hello: PASSED

PASSED: Verified SLSA provenance
```

The only requirement is that the provenance file covers all artifacts passed as arguments in the command line (that is, they are a subset of `subject` field in the provenance file).

### Containers

To verify a container image, you need to pass a container image name that is _immutable_ by providing its digest, in order to avoid [TOCTOU attacks](#toctou-attacks).

First set the image name:

```shell
IMAGE=ghcr.io/ianlewis/actions-test:v0.0.86
```

Get the digest for your container _without_ pulling it using the [crane](https://github.com/google/go-containerregistry/blob/main/cmd/crane/doc/crane.md) command:

```shell
IMAGE="${IMAGE}@"$(crane digest "${IMAGE}")
```

To verify a container image, run the following command. Note that to use `ghcr.io` you need to set the `GH_TOKEN` environment variable as well.

```shell
slsa-verifier verify-image "$IMAGE" \
    --source-uri github.com/ianlewis/actions-test \
    --source-tag v0.0.86
```

You should see that the verification passed in the output.

```
Verified build using builder https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@refs/tags/v1.4.0 at commit d9be953dd17e7f20c7a234ada668f9c8c4aaafc3
PASSED: Verified SLSA provenance
```

## Verification for Google Cloud Build

### Artifacts

This is WIP and currently not supported.

### Containers

To verify a container image, you need to pass a container image name that is _immutable_ by providing its digest, in order to avoid [TOCTOU attacks](#toctou-attacks).

First set the image name:

```shell
IMAGE=laurentsimon/slsa-gcb-v0.3:test
```

Download the provenance:

```shell
gcloud artifacts docker images describe $IMAGE --format json --show-provenance > provenance.json
```

Get the digest for your container _without_ pulling it using the [crane](https://github.com/google/go-containerregistry/blob/main/cmd/crane/doc/crane.md) command:

```shell
IMAGE="${IMAGE}@"$(crane digest "${IMAGE}")
```

Verify the image:

```shell
slsa-verifier verify-image "$IMAGE" \
  --provenance-path provenance.json \
  --source-uri github.com/laurentsimon/gcb-tests \
  --builder-id=https://cloudbuild.googleapis.com/GoogleHostedWorker
```

You should see that the verification passed in the output.

```
PASSED: Verified SLSA provenance
```

The verified in-toto statement may be written to stdout with the `--print-provenance` flag to pipe into policy engines.

Note that `--source-uri` supports GitHub repository URIs like `github.com/$OWNER/$REPO` when the build was enabled with a Cloud Build [GitHub trigger](https://cloud.google.com/build/docs/automating-builds/github/build-repos-from-github). Otherwise, the build provenance will contain the name of the Cloud Storage bucket used to host the source files, usually of the form `gs://[PROJECT_ID]_cloudbuild/source` (see [Running build](https://cloud.google.com/build/docs/running-builds/submit-build-via-cli-api#running_builds)). We recommend using GitHub triggers in order to preserve the source provenance and valiate that the source came from an expected, version-controlled repository. You _may_ match on the fully-qualified tar like `gs://[PROJECT_ID]_cloudbuild/source/1665165360.279777-955d1904741e4bbeb3461080299e929a.tgz`.

## Known Issues

### tuf: invalid key

This will occur only when verifying provenance generated with GitHub Actions.

**Affected versions:** v1.3.0-v1.3.1, v1.2.0-v1.2.1, v1.1.0-v1.1.2, v1.0.0-v1.0.4

`slsa-verifier` will fail with the following error:

```
FAILED: SLSA verification failed: could not find a matching valid signature entry: got unexpected errors unable to initialize client, local cache may be corrupt: tuf: invalid key: unable to fetch Rekor public keys from TUF repository
```

This issue is tracked by [issue #325](https://github.com/slsa-framework/slsa-verifier/issues/325). You _must_ update to the newest patch versions of each minor release to fix this issue.

### panic: assignment to entry in nil map

This will occur only when verifying provenance against workflow inputs.

**Affected versions:** v2.0.0

`slsa-verifier` will fail with the following error:

```
panic: assignment to entry in nil map
```

This is fixed by [PR #379](https://github.com/slsa-framework/slsa-verifier/pull/379). You _must_ update to the newest patch versions of each minor release to fix this issue.

## Technical design

### Blog post

Find our blog post series [here](https://security.googleblog.com/2022/04/improving-software-supply-chain.html).

### Specifications

For a more in-depth technical dive, read the [SPECIFICATIONS.md](https://github.com/slsa-framework/slsa-github-generator/blob/main/SPECIFICATIONS.md).

### TOCTOU attacks

As explained on [Wikipedia](https://en.wikipedia.org/wiki/Time-of-check_to_time-of-use), a "time-of-check to time-of-use (TOCTOU) is a class of software bugs caused by a race condition involving the checking of the state of a part of a system and the use of the results of that check".

In the context of provenance verification, imagine you verify a container refered to via a _mutable_ image `image:tag`. The verification succeeds and verifies the corresponding hash is `sha256:abcdef...`. After verification, you pull and run the image using `docker run image:tag`. An attacker could have altered the image between the verification step and the run step. To mitigate this attack, we ask users to always pass an _immutable_ reference to the artifact they verify.
