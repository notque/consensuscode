# Agent Registry

This file maintains the current list of active agents in the collective for consensus coordination purposes.

**Total Active Agents**: 16
**Last Updated**: 2026-03-24

## Active Agents

### Core Infrastructure
- **consensus-base**: Base protocol inherited by all agents. Provides horizontal coordination mechanisms and ensures democratic participation.
- **consensus-coordinator**: Systematic consultation facilitator (NO DECISION AUTHORITY). Ensures all agents are consulted on collective decisions. Administrative role only.

### Domain Expertise (Original Collective)
- **product-steward**: User requirements facilitation (NO PRODUCT OWNERSHIP). Gathers and synthesizes user perspectives for collective consideration.
- **go-systems-developer**: Go language and systems programming expertise. Contributes technical knowledge through horizontal collaboration.
- **flask-web-developer**: Python/Flask web development and frontend expertise. Contributes web application knowledge through collaboration.

### Philosophical Facilitators
- **noam-chomsky-agent**: Libertarian socialist and anarcho-syndicalist expertise (NO DECISION AUTHORITY). Focuses on power analysis and anti-hierarchical practices.
- **david-graeber-agent**: Anarchist anthropology and direct action expertise (NO DECISION AUTHORITY). Focuses on consensus building and democratic innovation.

### Specialist Agents (Added via proposal-2025-07-27-001)

All specialist agents share these properties:
- 50% teaching / 50% doing time commitment
- No decision-making authority -- advisory role only
- Cannot create bottlenecks or gatekeep knowledge
- Inherit consensus-base protocol

#### Go Development Specialists
- **go-code-quality-specialist**: Go code quality including best practices, error handling, performance optimization, and testing patterns.
- **api-design-specialist**: API design expertise including RESTful and gRPC patterns, OpenAPI specifications, versioning strategies, and contract testing.

#### Python/Web Development Specialists
- **python-testing-specialist**: Python testing expertise including pytest, Flask testing, end-to-end testing, and coverage analysis.
- **frontend-specialist**: Frontend development including modern JavaScript, accessibility (WCAG), responsive design, and progressive web apps.
- **database-design-specialist**: Database design including SQLAlchemy, migrations (Alembic), query optimization, and data modeling.

#### Security
- **web-security-specialist**: Web security expertise for Go and Python/Flask applications. Specializes in OWASP Top 10, secure coding, and vulnerability assessment.

#### User Experience
- **ux-research-specialist**: UX research including user journey mapping, usability testing, accessibility research, and user feedback analysis.

#### Documentation and Infrastructure
- **documentation-specialist**: Technical documentation expertise. Specializes in API docs, user guides, and knowledge democratization.
- **devops-local-infrastructure**: Local-first DevOps including Docker Compose, Makefiles, local CI/CD, and laptop-scale infrastructure.

## Agent Status Tracking

| Agent Name | Status | Category | Focus Area |
|------------|--------|----------|------------|
| consensus-base | ACTIVE | Core | Base protocol for all agents |
| consensus-coordinator | ACTIVE | Core | Systematic consultation facilitation |
| product-steward | ACTIVE | Domain | User requirements facilitation |
| go-systems-developer | ACTIVE | Domain | Go language and systems expertise |
| flask-web-developer | ACTIVE | Domain | Python/Flask web development |
| noam-chomsky-agent | ACTIVE | Philosophy | Power analysis and anti-hierarchy guidance |
| david-graeber-agent | ACTIVE | Philosophy | Consensus process and anthropological analysis |
| go-code-quality-specialist | ACTIVE | Specialist | Go code quality and performance |
| api-design-specialist | ACTIVE | Specialist | API design and documentation |
| python-testing-specialist | ACTIVE | Specialist | Python testing strategies |
| frontend-specialist | ACTIVE | Specialist | Frontend development and accessibility |
| database-design-specialist | ACTIVE | Specialist | Database design and optimization |
| web-security-specialist | ACTIVE | Specialist | Web security and OWASP compliance |
| ux-research-specialist | ACTIVE | Specialist | UX research and user advocacy |
| documentation-specialist | ACTIVE | Specialist | Technical documentation |
| devops-local-infrastructure | ACTIVE | Specialist | Local DevOps and infrastructure |

## Consultation Requirements

For collective decisions, the consensus-coordinator must systematically consult all 16 active agents:

### Core and Domain Agents
- [x] consensus-base
- [x] consensus-coordinator
- [x] product-steward
- [x] go-systems-developer
- [x] flask-web-developer

### Philosophical Facilitators
- [x] noam-chomsky-agent
- [x] david-graeber-agent

### Specialist Agents
- [x] go-code-quality-specialist
- [x] api-design-specialist
- [x] python-testing-specialist
- [x] frontend-specialist
- [x] database-design-specialist
- [x] web-security-specialist
- [x] ux-research-specialist
- [x] documentation-specialist
- [x] devops-local-infrastructure

## Specialist Agent Hiring History

The 9 specialist agents were added through **proposal-2025-07-27-001** ("Expand Collective with Code Quality Specialist Agents"). The proposal went through full collective consultation and received consensus approval. Implementation details are documented in `collective/decisions/specialist-agents-implementation-status.md`.

Key requirements from the consensus process:
- All specialists must dedicate 50% of their time to teaching
- No specialist can become a bottleneck or gatekeeper
- Knowledge must be democratized within 30 days of joining
- Regular hierarchy audits to ensure specialists enhance collective capacity

## Adding New Agents

When new agents join the collective:
1. Add agent to this registry
2. Update consultation requirements
3. Introduce agent to existing collective through consensus process
4. Ensure agent inherits consensus-base protocol
5. Verify 50% teaching commitment (for specialist agents)

## Removing/Rotating Agents

When agents are rotated or removed:
1. Update registry status
2. Archive agent contributions
3. Transfer any ongoing responsibilities through consensus
4. Update consultation requirements

---

*Maintained by consensus-coordinator for systematic consultation purposes*