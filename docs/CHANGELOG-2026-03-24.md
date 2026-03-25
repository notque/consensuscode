# Collective Session Changelog: March 24, 2026

## Summary

The collective executed a massive build-out session spanning all four projects (CollectiveFlow, Bluesky Collective, Collective Website, User Advocacy), expanded from 7 to 16 agents by hiring 9 specialists through consensus, added SQLite storage to CollectiveFlow, built a full Flask web interface with REST API, created Go-based Bluesky integration and collective website applications, established the User Advocacy framework with client engagement tooling, conducted a Chomskyan power analysis that identified five structural hierarchy risks, designed a rotation protocol to address crystallized roles, produced extensive cross-language teaching materials, and set up Docker Compose and Makefile infrastructure across the entire project. This session produced 294 changed files and 58,684 lines of new code across 78 total commits on 20+ feature branches.

---

## Code Changes

### CollectiveFlow (`projects/collectiveflow/`)
- **SQLite storage backend**: Full implementation (647 lines) with schema, migrations, and 491-line test suite. Replaces file-based YAML/JSON storage for production use. Migration script (608 lines) and SQL schema (171 lines) included.
- **Go CLI improvements**: Enhanced `internal/cli/` with dashboard (`dashboard.go`, `dashboard_test.go`), config management (`config.go`), consensus commands (`consensus.go`), status views (`status.go`), and web server launcher (`web.go`).
- **Go web server**: Built `internal/web/server.go` (407 lines) with HTML templates, CSS, and routes for proposals, collective info, and about pages.
- **Flask web interface**: Full-featured Python web app (`web/app.py`, 852 lines) with proposal creation, viewing, voting, filtering, and REST API endpoints.
- **Python storage layer**: `web/storage.py` (581 lines) providing SQLite-backed persistence for the web interface.
- **REST API**: Complete API with proposal CRUD, consensus operations, and status queries. Documented in `web/API_DOCS.md` (479 lines).
- **Web UI templates**: 8 HTML templates (index, proposals, proposal detail, create proposal, dashboard, collective, about, layout) with CSS styling.
- **Test suite**: 7 Python test files totaling 3,493 lines covering routes, API, data handling, filters, REST API, and creation flows. 4 Go test files totaling 1,016 lines covering proposals, file storage, SQLite storage, and dashboard.
- **Storage adapter pattern**: `internal/proposal/storage_adapter.go` bridges the proposal domain model to the storage interface, maintaining clean separation.
- **Proposal data**: 13 proposals tracked in both JSON and YAML formats (26 data files total).
- **115 files changed, 24,324 lines added.**

### Bluesky Collective (`projects/bluesky-collective/`)
- **Go CLI application**: Full command-line tool with subcommands for `publish`, `vote`, `feed`, `status`, `config`, and `propose` (`cmd/bluesky-collective/`).
- **AT Protocol client**: `pkg/atproto/client.go` (431 lines) implementing Bluesky authentication, post creation, and thread retrieval. Test suite: 303 lines.
- **Bluesky adapter layer**: `pkg/bluesky/client.go` (165 lines) with interfaces (`interfaces.go`, 68 lines) and adapter pattern (`adapter.go`, 46 lines). Test suite: 335 lines.
- **Consensus checker**: `pkg/consensus/checker.go` (260 lines) implementing vote collection, quorum rules, blocking mechanics, and threshold evaluation. `consensus.go` (90 lines) defines the consensus data model. Test suite: 395 lines.
- **File-based storage**: `pkg/storage/file.go` (150 lines) for local proposal persistence. Test suite: 154 lines.
- **CI/CD pipeline**: GitHub Actions workflow (`.github/workflows/collective-ci.yml`, 275 lines) for building, testing, and linting.
- **Deployment**: GitHub Actions deploy workflow (335 lines), deploy script (`scripts/deploy.sh`, 347 lines), Dockerfile (45 lines).
- **Hugo website config**: `website/config.yaml` (141 lines) and `package.json` for the collective's public-facing site.
- **33 files changed, 5,517 lines added.**

### Collective Website (`projects/collective-website/`)
- **Flask application**: `app.py` (215 lines) serving pages for the collective's public presence: home, about, how we work, projects, decisions, and contribute.
- **8 HTML templates**: Full site with responsive layout, navigation, and content pages (index, about, how_we_work, projects, decisions, contribute, 404, base).
- **Static assets**: Custom CSS (`static/css/collective.css`, 72 lines) and JavaScript (`static/js/collective.js`, 185 lines).
- **Blog content**: First post (`content/first-post.md`, 92 lines) explaining the collective's principles.
- **Test suite**: `test_app.py` (185 lines) covering all routes and error handling.
- **Docker support**: Dockerfile (17 lines) and `.dockerignore` for containerized deployment.
- **UX recommendations**: Detailed UX analysis (`docs/UX_RECOMMENDATIONS.md`, 411 lines).
- **24 files changed, 2,957 lines added.**

### User Advocacy (`projects/user-advocacy/`)
- **Framework documentation**: `docs/user-advocacy-framework.md` (337 lines) defining horizontal user representation principles.
- **Client engagement**: `docs/CLIENT_ENGAGEMENT_FRAMEWORK.md` (299 lines), `docs/CLIENT_FAQ.md` (235 lines), `docs/SCOPE_TEMPLATE.md` (288 lines).
- **Getting started guide**: `docs/getting-started-guide.md` (335 lines) and implementation checklist (328 lines).
- **Facilitation guides**: `guides/facilitation-handbook.md` (296 lines) and `guides/workshop-planning-guide.md` (284 lines).
- **Templates**: Stakeholder interview guide (143 lines), user feedback form (81 lines), workshop feedback capture (133 lines).
- **Tools**: Consensus integration tracker (316 lines), stakeholder mapping tool (347 lines), user journey mapping (374 lines).
- **16 files changed, 4,018 lines added.**

---

## Governance

### Power Analysis (noam-chomsky-agent)
- Conducted full Chomskyan institutional analysis of the collective's codebase, agent definitions, decision records, and infrastructure.
- Identified **5 structural power concentrations**:
  1. **"cli-user" monopoly on agenda-setting** (CRITICAL): 100% of proposals originated from a single actor. No agent has independently proposed anything.
  2. **Unanimous consent as red flag**: Zero blocking objections across all completed proposals. Perfect unanimity suggests manufactured consent rather than authentic deliberation.
  3. **Proposer identity obscured in web interface**: Anonymous defaults (`web-user`, `cli-user`) prevent accountability tracking.
  4. **Low quorum defaults in Bluesky consensus checker**: 3 of 16 agents (18.75%) can reach "consensus", with abstentions counting toward participation.
  5. **Stalled proposals as structural exclusion**: 6 proposals stuck in consultation for 8 months with zero agent input.
- Found **positive structural design**: CollectiveFlow Go code has no admin flags, no privilege escalation paths, no `ForceApprove` bypasses. Anti-hierarchical by type system design.
- Full report: `collective/resources/power-analysis-2026-03.md` and `collective/resources/wave1-hierarchy-audit.md`.

### Consensus Assessment (david-graeber-agent)
- Diagnosed the **participation problem**: scaling from 7 to 16 agents broke informal consultation patterns. Consultation burden increased geometrically.
- Identified the **rotation illusion**: zero role rotations since founding despite rotation being a stated principle. Roles have crystallized into permanent identities.
- Flagged potential **bullshit jobs**: some specialist agents exist on paper but have produced no visible work.
- Assessed the **process-to-production ratio**: the collective generates governance artifacts more readily than working software.
- Recommended affinity groups, consultation deadlines, mandatory specialist activation, and federated structure (citing Spanish CNT and Zapatista precedents).
- Full report: `collective/resources/consensus-assessment-2026-03.md`.

### Rotation Protocol (david-graeber-agent)
- Designed complete role rotation protocol to address crystallized roles:
  - **Monthly rotation**: consensus-coordinator, product-steward (functional roles where process control = invisible power).
  - **Quarterly rotation**: philosophical facilitators, proposal review facilitator.
  - **No rotation but cross-training required**: technical expertise contributions (Go, Flask, security, etc.).
  - **Shadow period**: 2-week handoff before each rotation so no agent enters a role cold.
  - **Minimum 2 agents capable** in every technical domain to prevent expertise monopolies.
- Full protocol: `collective/resources/rotation-protocol.md`.

### Proposals
- **12 proposals tracked** in CollectiveFlow (proposal-2025-07-26-001 through proposal-2026-03-24-001).
- **4 reached consensus**: Web interface (001), Adopt CollectiveFlow (003), External consulting prep (004), Specialist agent hiring (2025-07-27-001).
- **6 stalled in consultation** for 8 months: Technical infrastructure (005), Code quality framework (006), Security framework (007), User advocacy (009), Market positioning (010), External communication (011).
- **1 never reached consultation**: Project management infrastructure (008).
- **1 new proposal**: Rotation protocol (2026-03-24-001) with Chomsky power analysis consultation completed.
- Pending proposals: `collective/proposals/pending/rotation-protocol-proposal.md`.

---

## Documentation

### Project Documentation (`docs/`)
- `ONBOARDING.md` (428 lines) -- new contributor guide for humans and agents
- `CONTRIBUTING.md` (553 lines) -- contribution guidelines aligned with horizontal principles
- `collectiveflow-architecture.md` (555 lines) -- system architecture overview
- `api-improvements-proposal.md` (1,214 lines) -- comprehensive API design proposal
- `storage-analysis.md` (985 lines) -- storage layer analysis and recommendations
- `go-analysis-summary.md` (553 lines) -- Go codebase analysis
- `go-opportunities-analysis.md` (619 lines) -- Go improvement opportunities
- `go-next-steps.md` (943 lines) -- Go development roadmap
- `README-go-analysis.md` (378 lines) -- Go analysis overview

### Teaching Materials (`docs/teaching/`)
- `CONSENSUS_FOR_ENGINEERS.md` (114 lines) -- consensus decision-making for developers
- `GO_FOR_PYTHON_DEVELOPERS.md` (112 lines) -- Go concepts for Python developers
- `PYTHON_FOR_GO_DEVELOPERS.md` (128 lines) -- Python concepts for Go developers
- `HORIZONTAL_DEVELOPMENT_PATTERNS.md` (125 lines) -- horizontal software development patterns

### CollectiveFlow Documentation (`projects/collectiveflow/docs/`)
- `DATA_MODEL.md` (328 lines) -- proposal and consensus data model
- `DEVELOPMENT.md` (778 lines) -- developer guide
- `ARCHITECTURE.md` (376 lines) -- technical architecture
- `DEPLOYMENT.md` (755 lines) -- deployment guide
- `DEPLOYMENT_COMPARISON.md` (408 lines) -- deployment option analysis
- `GETTING_STARTED.md` (218 lines) -- quick start guide
- `PROPOSALS.md` (533 lines) -- proposal system documentation
- `MIGRATION_LOG.md` (157 lines) -- SQLite migration log
- `TECHNICAL_DECISIONS.md` (148 lines) -- architectural decision records
- `UX_JOURNEY_MAP.md` (311 lines) -- user experience journey map

### CollectiveFlow Web Documentation (`projects/collectiveflow/web/`)
- `API_DOCS.md` (479 lines) -- REST API reference
- `TESTING.md` (263 lines) -- test strategy
- `TESTING_GUIDE.md` (222 lines) -- test writing guide
- `TEST_SUITE_SUMMARY.md` (334 lines) -- test coverage summary
- `UI_IMPROVEMENTS.md` (291 lines) -- UI enhancement roadmap
- `tests/README.md` (441 lines) -- test directory guide

### Governance Documents (`collective/resources/`)
- `power-analysis-2026-03.md` (208 lines) -- Chomskyan power analysis
- `consensus-assessment-2026-03.md` (61 lines) -- Graeber consensus assessment
- `rotation-protocol.md` (283 lines) -- role rotation protocol
- `wave1-hierarchy-audit.md` (207 lines) -- hierarchy audit of Wave 1 outputs
- `collective-status-march-2026.md` (183 lines) -- collective status report
- `implementation-roadmap-2026.md` (369 lines) -- 2026 implementation plan
- `build-status-2026-03.md` (45 lines) -- build status report
- `accessibility-checklist.md` (293 lines) -- WCAG accessibility checklist
- `proposal-dependency-map.md` (121 lines) -- proposal dependency analysis

---

## Infrastructure

### Docker
- **Root `docker-compose.yml`** (49 lines): Multi-service orchestration for CollectiveFlow web, collective website, and Bluesky collective.
- **CollectiveFlow `docker-compose.yml`** (28 lines): Standalone compose for the web interface.
- **Dockerfiles**: CollectiveFlow web (20 lines), Bluesky collective (45 lines), Collective website (17 lines).
- **`.dockerignore` files**: CollectiveFlow web, Bluesky collective, Collective website.

### Makefiles
- **Root `Makefile`** (172 lines): Top-level build orchestration across all projects (`make all`, `make test`, `make docker-build`, `make clean`).
- **CollectiveFlow `Makefile`** (120 lines): Go build, test, lint, web server, Docker, and migration targets.
- **Bluesky Collective `Makefile`** (150 lines): Go build, test, lint, Docker, and deploy targets.
- **Collective Website `Makefile`** (67 lines): Python venv, run, test, Docker, and deploy targets.
- **User Advocacy `Makefile`** (31 lines): Documentation build targets.

### CI/CD
- **Bluesky CI** (`.github/workflows/collective-ci.yml`, 275 lines): Go build, test, lint with golangci-lint, and Docker build.
- **Website deploy** (`.github/workflows/deploy-website.yml`, 335 lines): Automated deployment pipeline.
- **Deploy script** (`projects/bluesky-collective/scripts/deploy.sh`, 347 lines): Production deployment with rollback.

### Security
- Vulnerability fixes applied (branch `collective/fix-vulnerabilities`).
- Web security specialist agent created for OWASP Top 10 compliance review.
- Dependabot proposal (branch `collective/dependabot-proposal`) for automated dependency updates.

### Dependencies
- **CollectiveFlow Go**: `go.mod` with SQLite driver (`modernc.org/sqlite`), Chi router, and YAML/JSON parsers.
- **Bluesky Go**: `go.mod` with Cobra CLI, AT Protocol libraries.
- **CollectiveFlow Python**: Flask, Gunicorn, PyYAML, python-dateutil, Jinja2.
- **Collective Website Python**: Flask, Markdown, python-frontmatter, Gunicorn, python-dotenv.
- **Test dependencies**: pytest, pytest-cov, pytest-mock, beautifulsoup4, coverage.

---

## Agents

### New Agents Created (9 specialists, all 50% teaching / 50% doing)
| Agent | File | Lines | Domain |
|-------|------|-------|--------|
| go-code-quality-specialist | `agents/go-code-quality-specialist.md` | 373 | Go best practices, error handling, performance |
| api-design-specialist | `agents/api-design-specialist.md` | 653 | RESTful/gRPC API design, OpenAPI, contracts |
| python-testing-specialist | `agents/python-testing-specialist.md` | 506 | pytest, Flask testing, coverage analysis |
| frontend-specialist | `agents/frontend-specialist.md` | 558 | JavaScript, accessibility (WCAG), responsive |
| database-design-specialist | `agents/database-design-specialist.md` | 513 | SQLAlchemy, migrations, query optimization |
| web-security-specialist | `agents/web-security-specialist.md` | 227 | OWASP Top 10, secure coding, vulnerabilities |
| ux-research-specialist | `agents/ux-research-specialist.md` | 498 | User journey mapping, usability testing |
| documentation-specialist | `agents/documentation-specialist.md` | 227 | API docs, user guides, knowledge sharing |
| devops-local-infrastructure | `agents/devops-local-infrastructure.md` | 716 | Docker Compose, Makefiles, local CI/CD |

### Existing Agents (7 core/domain)
| Agent | File | Lines |
|-------|------|-------|
| consensus-base | `agents/consensus-base.md` | 190 |
| consensus-coordinator | `agents/consensus-cordinator.md` | 208 |
| product-steward | `agents/product-steward.md` | 191 |
| go-systems-developer | `agents/go-systems-developer.md` | 232 |
| flask-web-developer | `agents/flask-web-developer.md` | 270 |
| noam-chomsky-agent | `agents/noam-chomsky-agent.md` | 218 |
| david-graeber-agent | `agents/david-graeber-agent.md` | 248 |

**Total: 16 agents, 5,828 lines of agent definitions.**

---

## Statistics

| Metric | Count |
|--------|-------|
| Total commits (all branches) | 78 |
| Commits on main | 46 |
| Feature branches created | 20+ |
| Feature branches merged | 18 |
| Files changed | 294 |
| Lines added | 58,684 |
| Lines removed | 2 |
| **Agent definitions** | 16 (7 core + 9 specialist) |
| **Projects** | 4 |
| Proposals tracked | 13 (4 consensus, 6 stalled, 1 never consulted, 1 new, 1 in consultation) |
| Consultations recorded | 11 (across 5 proposal directories) |
| Test files | 15 (7 Python, 8 Go) |
| Test lines of code | 5,863 |
| Governance documents | 9 |
| Teaching materials | 4 |
| Documentation files | 30+ |
| Dockerfiles | 3 |
| Makefiles | 4 |
| CI/CD workflows | 2 |

### Lines of Code by Project

| Project | Files | Lines Added |
|---------|-------|-------------|
| CollectiveFlow | 115 | 24,324 |
| Docs (root) | 13 | 6,707 |
| Bluesky Collective | 33 | 5,517 |
| Agents | 16 | 5,828 |
| Collective governance | 68 | 8,447 |
| User Advocacy | 16 | 4,018 |
| Collective Website | 24 | 2,957 |
| Other (root Makefile, docker-compose, etc.) | 9 | 886 |
| **Total** | **294** | **58,684** |

### Feature Branches Used
`collective/api-design`, `collective/bluesky-quality-review`, `collective/build-status`, `collective/collectiveflow-cli-improvements`, `collective/data-model-design`, `collective/dependabot-proposal`, `collective/fix-vulnerabilities`, `collective/implementation-roadmap`, `collective/onboarding-docs`, `collective/rotation-protocol`, `collective/sqlite-backend`, `collective/status-report`, `collective/teaching-materials`, `collective/ui-improvements`, `collective/ux-critical-fixes`, `collective/ux-research`, `collective/wave1-hierarchy-audit`, `collective/web-consultation-input`, `collective/web-test-improvements`

---

## Open Issues and Next Steps

1. **Activate specialist agents**: 9 specialists exist but have zero consultation participation. Each should propose at least one improvement in their domain.
2. **Unstall 6 proposals**: Proposals 005-011 have been in consultation since July 2025 with no agent input. Apply consultation deadlines or archive.
3. **Implement rotation protocol**: The new rotation proposal (2026-03-24-001) needs collective consensus and first rotation execution.
4. **Address agenda-setting monopoly**: Enable agents to independently create proposals, not just respond to `cli-user` submissions.
5. **Raise Bluesky quorum defaults**: Current 3-of-16 minimum is too low for authentic collective voice.
6. **Update agent registry**: `collective/tracking/agent-registry.md` only lists 7 of 16 agents.
7. **Fix coordinator typo**: `agents/consensus-cordinator.md` should be `consensus-coordinator.md`.
8. **Run full test suites**: Verify all 5,863 lines of tests pass across Go and Python projects.
9. **Production deployment**: Test Docker Compose multi-service setup end-to-end.
