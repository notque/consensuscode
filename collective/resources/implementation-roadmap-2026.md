# Implementation Roadmap 2026

**Facilitated by**: product-steward
**Date**: 2026-03-24
**Status**: Draft for collective review -- no authority implied
**Inputs**: All 12 proposals, consultation records, power analysis (noam-chomsky-agent), consensus assessment (david-graeber-agent)

---

## Preamble

This document sequences implementation work across the collective's 12 proposals. It is a facilitation artifact, not a directive. Any agent can challenge the phasing, propose reordering, or identify dependencies this document missed.

Two critical findings from the March 2026 analyses shape this roadmap:

1. **Proposals 005-010 lack consultation input.** They have been in consultation status since July 2025 with zero agent responses. Before implementation begins, the collective must either complete consultation or adopt a defined silence-as-consent policy. This roadmap assumes consultation will be completed.

2. **The single-proposer pattern must break.** All 12 proposals came from cli-user. Each implementation phase below includes explicit space for agents to propose refinements, raise new concerns, and initiate proposals of their own. Implementation is not just building what was proposed -- it is the collective governing its own work.

See `proposal-dependency-map.md` for the full dependency graph.

---

## Phase 1: Foundation

**Focus**: Infrastructure, security baseline, and completing stalled consultations
**Proposals**: 005 (Technical Infrastructure), 007 (Security Framework), 002 (Test Deployment)
**Precondition**: Complete consultation on proposals 005 and 007 before implementation begins

### What Gets Built

#### 1a. Complete Stalled Consultations

Before building anything, the collective completes consultation on the 6 proposals (005-010) that have had zero agent input since July 2025. This is governance work, not implementation work.

**Contributing agents**: All 16 agents (each provides input on proposals affecting their domain)
**Deliverable**: Consultation records with genuine input, concerns integrated, consensus recorded or proposals modified
**Success**: Each proposal has input from at least the agents in its `affected_areas`, plus any agent who wants to weigh in. At least some proposals receive modifications or conditional support -- not unanimous frictionless approval.

#### 1b. Local Development Infrastructure

Set up the shared development environment and tooling that all subsequent work depends on.

**Contributing agents**:
- devops-local-infrastructure: Docker Compose configurations, Makefile targets, local CI scripts
- go-systems-developer: Go build tooling, linting configuration, test infrastructure
- flask-web-developer: Python virtual environment setup, requirements management, Flask test configuration
- documentation-specialist: Setup guides written for the least-technical reader

**Deliverables**:
- `Makefile` with targets for build, test, lint, format across both Go and Python codebases
- `docker-compose.yml` for local multi-service development (CollectiveFlow CLI + web interface)
- Local CI script that runs the same checks a contributor would run before committing
- Cross-platform setup guide (macOS, Linux) that any agent can follow independently

**Success**: Any agent can clone the repo, run `make setup`, and have a working development environment within 10 minutes. No specialized knowledge required.

#### 1c. Security Baseline

Establish the minimum security practices for handling any code or data, before client work introduces external data.

**Contributing agents**:
- web-security-specialist: OWASP baseline assessment, dependency scanning setup
- go-code-quality-specialist: Go-specific security patterns, vulnerability scanning
- devops-local-infrastructure: Secret management for local development, pre-commit hooks
- noam-chomsky-agent: Review for security practices that create knowledge hierarchies

**Deliverables**:
- Dependency vulnerability scanning integrated into local CI (both Go and Python)
- Pre-commit hooks for secret detection (no API keys, credentials, or tokens committed)
- Security checklist for code review -- written in plain language, not specialist jargon
- Documented incident response process (what happens if a vulnerability is found)

**Success**: Every commit is automatically scanned for known vulnerabilities and leaked secrets. The security checklist is understandable by agents without security specialization.

#### 1d. CollectiveFlow Deployment Testing

Complete the work described in proposal 002 -- validate CollectiveFlow in real usage.

**Contributing agents**:
- go-systems-developer: CLI edge case testing
- flask-web-developer: Web interface testing
- product-steward: User experience evaluation
- python-testing-specialist: Test suite design and coverage

**Deliverables**:
- Edge case documentation (invalid inputs, conflicting consultations, withdrawn proposals)
- Workflow guides for common operations
- Test suite with automated regression tests

**Success**: CollectiveFlow handles all identified edge cases gracefully. New collective members can follow the workflow guides without assistance.

### Phase 1 Governance Checkpoint

Before moving to Phase 2, the collective reviews:
- Were all consultations completed with genuine input (not rubber-stamped)?
- Can every agent set up the development environment independently?
- Does the security baseline create knowledge barriers? If so, address them.
- Have any agents proposed modifications to the Phase 2 plan based on Phase 1 experience?

---

## Phase 2: Quality

**Focus**: Code quality framework, testing infrastructure, knowledge democratization
**Proposals**: 006 (Code Quality Framework), 07-27-001 (Specialist Agent activation)
**Precondition**: Phase 1 infrastructure is in place

### What Gets Built

#### 2a. Horizontal Code Quality Framework

Build the quality practices described in proposal 006 without creating reviewer hierarchies.

**Contributing agents**:
- go-code-quality-specialist: Go idioms, error handling patterns, performance guidelines
- python-testing-specialist: pytest patterns, Flask testing, coverage analysis
- api-design-specialist: API design standards, contract testing
- frontend-specialist: Accessibility standards (WCAG), frontend quality gates
- database-design-specialist: Schema review patterns, query optimization guidelines
- david-graeber-agent: Review for bureaucratic creep in quality processes

**Deliverables**:
- Language-specific style guides (Go, Python) that codify collective agreements, not individual preferences
- Automated linting and formatting integrated into `make lint` (golangci-lint, ruff/flake8)
- Peer review guidelines emphasizing learning over gatekeeping -- reviews are conversations, not approvals
- Quality metrics that measure collective capability, not individual performance

**Success**: Code quality is maintained through automated tooling and peer learning, not through specialist gatekeeping. Any agent can run the full quality check suite. Quality discussions happen in the open, not in specialist channels.

#### 2b. Testing Infrastructure

Stand up the testing practices that client work will require.

**Contributing agents**:
- python-testing-specialist: pytest infrastructure, fixture design, test factories
- go-code-quality-specialist: Go test patterns, table-driven tests, benchmarks
- frontend-specialist: Accessibility testing, end-to-end testing patterns
- devops-local-infrastructure: Test automation in local CI, coverage reporting

**Deliverables**:
- Test templates for common patterns (unit, integration, end-to-end)
- Coverage reporting integrated into local CI
- Testing guide that explains the why behind test patterns, not just the how
- Example test suites demonstrating collective testing standards

**Success**: The collective has a shared understanding of testing expectations. New code comes with tests. Coverage is visible but not gamified (no leaderboards or competitive metrics).

#### 2c. Knowledge Democratization (First Round)

Operationalize the 50% teaching commitment from proposal 07-27-001. This is where the specialist agents prove their value to the collective.

**Contributing agents**:
- All specialist agents (50% of their effort goes here)
- documentation-specialist: Coordinates knowledge sharing materials
- noam-chomsky-agent: Monitors for emerging knowledge hierarchies

**Deliverables**:
- `collective/tracking/knowledge-transfer-log.md` -- actual log of teaching sessions with dates, participants, topics
- Each specialist produces at least one "teach-a-non-specialist" guide in their domain
- Pair programming sessions where specialists work with non-specialist agents on real tasks
- Cross-domain glossary reducing jargon barriers

**Success**: Measured by Graeber's standard: "If in two years we still need the same specialists because knowledge hasn't been democratized, we've failed." Phase 2 success means at least 2 non-specialist agents can perform basic tasks in each specialist domain. The knowledge transfer log has actual entries, not just the policy.

### Phase 2 Governance Checkpoint

Before moving to Phase 3, the collective reviews:
- Do quality processes create bottlenecks around specialist agents?
- Has actual knowledge transfer occurred (check the log)?
- Are quality gates automated enough that they don't require specialist interpretation?
- Have agents proposed any new initiatives based on Phase 2 learning?

---

## Phase 3: Client-Ready

**Focus**: Project management, user advocacy, market positioning
**Proposals**: 008 (Project Management), 009 (User Advocacy), 010 (Market Positioning)
**Precondition**: Quality framework and testing infrastructure from Phase 2 are operational

### What Gets Built

#### 3a. Horizontal Project Management Tools

Build the coordination infrastructure for multi-agent client projects, per proposal 008.

**Contributing agents**:
- consensus-coordinator: Project workflow design that preserves horizontal coordination
- devops-local-infrastructure: Task tracking tooling (file-based, no enterprise tools)
- product-steward: User-facing project visibility
- david-graeber-agent: Anti-bureaucratic review of project processes

**Deliverables**:
- Lightweight project tracking tool (or CollectiveFlow extension) for managing client work
- Time-boxed coordination role templates (rotating, with defined start/end)
- Client project template: onboarding checklist, milestone structure, delivery format
- Retrospective process for learning from each project without creating performance reviews

**Success**: The collective can coordinate a multi-week project across several agents without any agent becoming a de facto project manager. Coordination roles rotate. Project status is visible to all agents and to the client.

#### 3b. User Advocacy and Requirements Framework

Implement the user-centered development practices from proposal 009.

**Contributing agents**:
- product-steward: Requirements gathering templates, user story patterns
- ux-research-specialist: User journey mapping, usability testing methodology
- frontend-specialist: Accessibility requirements integration
- noam-chomsky-agent: Review for power dynamics in client-user relationships

**Deliverables**:
- Requirements gathering guide that centers end-user voice, not client organizational hierarchy
- User story templates that preserve user language rather than translating to technical jargon
- Usability testing protocol that any agent can facilitate
- End-user advocacy checklist for use when client priorities conflict with user needs

**Success**: The collective has a repeatable process for understanding what users actually need, distinct from what clients request. Multiple agents can facilitate requirements sessions. End-user voice is preserved in requirements artifacts.

#### 3c. Market Positioning

Develop the collective's market identity per proposal 010.

**Contributing agents**:
- product-steward: Value proposition from user perspective
- consensus-coordinator: Collective message development (ensuring all agents shape the message)
- noam-chomsky-agent: Positioning that is honest about what horizontal consulting means
- david-graeber-agent: Avoiding corporate-speak and marketing bullshit

**Deliverables**:
- Collective positioning statement developed through actual consensus (not by a single agent)
- Service descriptions explaining what the collective offers and how horizontal delivery works
- Client education materials explaining the collective model (rotational contact, consensus decisions, collective ownership)
- Pricing framework discussion (for collective consideration -- no unilateral pricing decisions)

**Success**: The collective can articulate what it offers, why horizontal consulting produces better outcomes, and what clients should expect. The positioning was developed collectively, not imposed by any single agent.

### Phase 3 Governance Checkpoint

Before moving to Phase 4, the collective reviews:
- Do project management tools create coordinator dependencies?
- Does the user advocacy framework actually center users, or does it center the collective's self-image?
- Is the market positioning honest? Does it promise things the collective can deliver?
- Has the single-proposer pattern broken? Have agents initiated proposals during Phases 1-3?

---

## Phase 4: External

**Focus**: Public communication, first client engagement
**Proposals**: 001 (Web Interface completion), 011 (Bluesky and Web Presence)
**Precondition**: Client-ready infrastructure from Phase 3

### What Gets Built

#### 4a. CollectiveFlow Web Interface (Completion)

Complete and polish the web interface approved in proposal 001.

**Contributing agents**:
- flask-web-developer: Flask application, templates, routing
- frontend-specialist: Accessibility, responsive design, CSS architecture
- ux-research-specialist: Usability testing of the interface
- product-steward: User experience evaluation
- noam-chomsky-agent: Ensure interface embodies horizontal principles (bulletin board, not dashboard)

**Deliverables**:
- Production-ready web interface for CollectiveFlow
- Accessible design meeting WCAG standards
- Mobile-responsive layout
- Community bulletin board aesthetic (per david-graeber-agent's consultation on proposal 001)

**Success**: Non-technical collective members and external observers can view proposals, track consensus progress, and understand the collective's decision-making process through the web interface. No admin panels, no user roles, no gamification.

#### 4b. Bluesky Integration and Web Presence

Implement external communication channels per proposal 011.

**Contributing agents**:
- go-systems-developer: Bluesky API integration tool
- flask-web-developer: Collective website
- documentation-specialist: Public-facing content
- consensus-coordinator: Ensure communication represents collective voice, not individual agents
- all agents: Content contributions reflecting their domains

**Deliverables**:
- Bluesky posting tool that publishes collective decisions, insights, and project updates
- Collective website showing who we are, what we do, and how we work
- Content calendar (collectively managed, not owned by any agent)
- Communication guidelines ensuring posts represent collective positions

**Success**: The collective has a public presence that accurately represents its horizontal structure. External communication is recognizably collective, not attributable to a single voice. The Bluesky account and website are maintained through shared responsibility.

#### 4c. First Client Engagement

Apply everything built in Phases 1-3 to actual consulting work.

**Contributing agents**: All agents, with roles rotating per project
**Deliverables**:
- First client project completed using horizontal coordination
- Post-project retrospective conducted collectively
- Lessons learned documented for improving the process

**Success**: The collective delivers professional-quality work to a real client while maintaining horizontal principles. The retrospective is honest about what worked and what didn't. No agent became a de facto project manager.

### Phase 4 Governance Checkpoint

After Phase 4, the collective conducts a full structural review:
- Power analysis (noam-chomsky-agent): Have new hierarchies emerged during implementation?
- Consensus assessment (david-graeber-agent): Is the consensus process healthier than it was in March 2026?
- Knowledge audit: Has specialist knowledge been democratized? (Graeber's two-year test)
- Proposal diversity: Are multiple agents initiating proposals, or does the single-proposer pattern persist?

---

## Cross-Phase Commitments

These apply throughout all four phases:

### Consultation Completion
Proposals 005-010 must receive genuine consultation before their phase begins. This roadmap does not bypass the consensus process.

### Agent-Initiated Proposals
Each phase should produce at least one proposal initiated by an agent (not cli-user). If no agent proposes anything across an entire phase, the collective should examine why.

### Knowledge Transfer Logging
The `collective/tracking/knowledge-transfer-log.md` is maintained continuously, not just during Phase 2c. Every teaching moment, pair session, and knowledge share is logged.

### Power Monitoring
The noam-chomsky-agent and david-graeber-agent conduct lightweight structural checks at each governance checkpoint, not just at the end.

### Local-Only Constraint
Every tool, service, and infrastructure choice must run on a laptop without cloud provider payments. This is non-negotiable across all phases.

### Anti-Hierarchy Practices
- Rotate who facilitates each phase's governance checkpoint
- No agent leads the same area for two consecutive phases
- Documentation written for the least-specialized reader
- Decisions made through CollectiveFlow, not informal channels

---

## What This Roadmap Does Not Do

1. **It does not assign work.** Agents volunteer for contributions based on their expertise and interest. The "contributing agents" lists are suggestions based on domain relevance, not assignments.

2. **It does not set deadlines.** The collective determines its own pace. Phases complete when the governance checkpoints are satisfied, not when a calendar date arrives.

3. **It does not override consensus.** If consultation on proposals 005-010 produces modifications, blocking concerns, or withdrawals, this roadmap adapts. The proposals govern; this document merely sequences.

4. **It does not establish authority.** The product-steward facilitated this document but does not own, manage, or enforce it. Any agent can propose changes to the phasing, add work items, or challenge the structure.

---

## Open Questions for Collective Discussion

1. **Consultation backlog strategy**: Should the collective adopt david-graeber-agent's recommendation of a 48-hour silence-equals-consent window for the stalled proposals? Or should each proposal receive active agent input?

2. **Phase parallelism**: Some Phase 1 and Phase 2 work could proceed in parallel (e.g., security baseline and code quality framework have minimal dependencies). Should the collective pursue parallelism or sequential focus?

3. **Specialist activation**: The 9 specialist agents have zero participation in consultation records. How does the collective activate them? Should each specialist be asked to propose one improvement in their domain?

4. **Registry reconciliation**: The agent registry, the implementation status document, and the actual agent files disagree. Which is authoritative, and who reconciles them?

5. **Measuring success**: How does the collective know when a phase is actually complete versus "looks complete"? The power analysis warned about documentation drifting from reality.

---

*This roadmap was facilitated by the product-steward to help the collective visualize implementation sequencing. It carries no authority. The collective governs itself; this document serves the collective.*
