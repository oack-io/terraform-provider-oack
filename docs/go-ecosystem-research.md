# Go Ecosystem Research: Badges, Linters & Code Quality

Research based on analysis of 30 top Go repositories on GitHub (kubernetes, terraform, vault, prometheus, grafana, etcd, moby, minio, gin, cobra, etc.).

## 1. README Badges — Adoption Ranking

| Rank | Badge | Adoption | URL Pattern |
|------|-------|----------|-------------|
| 1 | **Go Report Card** | ~17/30 | `goreportcard.com/badge/github.com/...` |
| 2 | **GitHub Actions CI** | ~14/30 | `github.com/.../actions/workflows/...badge.svg` |
| 3 | **License** | ~12/30 | `img.shields.io/badge/license-...` |
| 4 | **GoDoc / pkg.go.dev** | ~11/30 | `pkg.go.dev/badge/...` |
| 5 | **Release/Version** | ~8/30 | `img.shields.io/github/v/release/...` |
| 6 | **Codecov** | ~6/30 | `codecov.io/gh/.../badge.svg` |
| 7 | **Docker Pulls** | ~5/30 | `img.shields.io/docker/pulls/...` |
| 8 | **OpenSSF Scorecard** | 3/30 | `api.scorecard.dev/projects/...` |
| 9 | **CII Best Practices** | 3/30 | `bestpractices.dev/...` |
| 10 | **Sourcegraph** | 3/30 | `sourcegraph.com/github.com/...` |

### Terraform provider-specific badges (from hashicorp/* repos)

| Badge | Example |
|-------|---------|
| Terraform Registry | `img.shields.io/badge/terraform-registry-blueviolet` |
| Terraform version compat | `img.shields.io/badge/terraform-%3E%3D1.0-blue` |

## 2. Linters — Adoption in golangci-lint Configs

Based on 19 repos with `.golangci.yml` files.

| Rank | Linter | Count | In Our Repo? | Recommendation |
|------|--------|-------|--------------|----------------|
| 1 | govet | 17/19 | Yes | — |
| 2 | staticcheck | 16/19 | Yes | — |
| 3 | ineffassign | 15/19 | Yes | — |
| 4 | revive | 14/19 | Yes | — |
| 5 | unconvert | 13/19 | Yes | — |
| 6 | unused | 12/19 | Yes | — |
| 7 | gocritic | 12/19 | Yes | — |
| 8 | misspell | 11/19 | Yes | — |
| 9 | errcheck | 10/19 | Yes | — |
| 10 | **depguard** | 10/19 | **No** | Add — prevents unwanted imports |
| 11 | **bodyclose** | 8/19 | **No** | **Add** — critical for HTTP clients |
| 12 | **errorlint** | 8/19 | Yes | — |
| 13 | **gosec** | 7/19 | **No** | **Add** — security scanner |
| 14 | **usestdlibvars** | 7/19 | **No** | Add — catches missed stdlib constants |
| 15 | **nakedret** | 6/19 | **No** | Add — catches naked returns |
| 16 | **nolintlint** | 6/19 | **No** | **Add** — enforces nolint discipline |
| 17 | **whitespace** | 6/19 | **No** | Add — blank line consistency |
| 18 | copyloopvar | 6/19 | Yes | — |
| 19 | **unparam** | 5/19 | **No** | Add — unused function params |
| 20 | **durationcheck** | 5/19 | **No** | Add — duration multiplication bugs |
| 21 | **testifylint** | 5/19 | **No** | Add if using testify |
| 22 | exhaustive | — | Yes | — |
| 23 | intrange | — | Yes | — |

### Formatters

| Formatter | Adoption | In Our Repo? |
|-----------|----------|--------------|
| gofmt | 5/19 | Yes |
| goimports | 6/19 | **No** — recommended |
| gofumpt | 4/19 | No |

## 3. Security & Quality CI Tools

| Tool | Adoption | In Our Repo? | Recommendation |
|------|----------|--------------|----------------|
| golangci-lint | 19/30 | Yes | — |
| **CodeQL** | 6/30 | **No** | **Add** — free, catches vuln patterns |
| **govulncheck** | 2/30 | **No** | **Add** — Go vulnerability database |
| **Dependabot** | 5/30 | **No** | **Add** — auto-update dependencies |
| gosec (via lint) | 7/19 | **No** | Add to golangci-lint |
| OpenSSF Scorecard | 3/30 | No | Optional |
| Trivy | 1/30 | No | Optional |
| gocyclo (complexity) | 3/19 | No | Optional |

## 4. Gap Analysis — What to Apply

### Priority 1 — Must have (industry standard)

| Item | Action |
|------|--------|
| README badges | Add: Go Report Card, CI status, pkg.go.dev, License, Release, Terraform Registry |
| gosec linter | Add to `.golangci.yml` — security scanning |
| bodyclose linter | Add — critical for HTTP client code |
| nolintlint linter | Add — enforces proper `//nolint` annotations |
| Dependabot | Add `.github/dependabot.yml` |

### Priority 2 — Recommended (high adoption)

| Item | Action |
|------|--------|
| CodeQL workflow | Add `.github/workflows/codeql.yml` |
| govulncheck workflow | Add to CI — weekly vulnerability check |
| usestdlibvars linter | Add — catches `http.StatusOK` vs `200` etc. |
| nakedret linter | Add — catches `return` without values |
| whitespace linter | Add — blank line consistency |
| durationcheck linter | Add — duration multiplication bugs |
| goimports formatter | Replace gofmt with goimports |

### Priority 3 — Nice to have

| Item | Action |
|------|--------|
| unparam linter | Add — finds unused function parameters |
| depguard linter | Add — restrict imports |
| testifylint linter | Add if using testify |
| OpenSSF Scorecard badge | Add after above items are in place |
| Codecov integration | Add if test coverage tracking desired |
