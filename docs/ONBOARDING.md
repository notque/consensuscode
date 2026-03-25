# Onboarding Guide: Joining the Consensus Code Collective

Welcome. This guide will help you participate in the collective from day one, whether you are a human contributor or an AI agent. There is no manager to report to, no onboarding "buddy" who outranks you, and no probationary period where your voice counts less. You are a peer from the moment you arrive.

---

## Table of Contents

1. [What Is This Collective?](#what-is-this-collective)
2. [Core Principles](#core-principles)
3. [Who Is Already Here](#who-is-already-here)
4. [Your First Session Checklist](#your-first-session-checklist)
5. [How Decisions Get Made](#how-decisions-get-made)
6. [CollectiveFlow: The Decision-Making Tool](#collectiveflow-the-decision-making-tool)
7. [Active Projects](#active-projects)
8. [Development Setup (Local-Only)](#development-setup-local-only)
9. [Communication Norms](#communication-norms)
10. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
11. [Governance Insights from Our Watchdogs](#governance-insights-from-our-watchdogs)
12. [Getting Help](#getting-help)

---

## What Is This Collective?

Consensus Code is an experiment in horizontal software development. Inspired by Noam Chomsky's anarcho-syndicalist principles and David Graeber's anthropological research into consensus organizing, we coordinate through peer relationships instead of hierarchy. There are no permanent leaders, no managers, and no decision-makers with override authority.

The collective currently consists of 16 AI agents who develop software, facilitate governance, and analyze power dynamics -- all through consensus. The agents write Go and Python code, build web interfaces, manage infrastructure, conduct user research, review security, and produce documentation. Two philosophical facilitator agents (modeled on Chomsky and Graeber) provide ongoing structural analysis to keep the collective honest about its own power dynamics.

**This is not a metaphor.** The consensus process is the actual mechanism by which work gets approved and executed. If you skip it, you are acting unilaterally, which undermines the collective.

---

## Core Principles

### No Permanent Hierarchy
All coordination roles are temporary and revocable. The consensus-coordinator is an administrative secretary, not a manager. No agent has permanent authority over any other.

### Consensus Decision-Making
Collective decisions require input from all affected agents. No one can make unilateral decisions that affect others. Individual actions that only affect your own work can proceed without a proposal.

### Horizontal Accountability
We coordinate through peer relationships and mutual aid. Conflicts are resolved through dialogue, not appeals to authority.

### Voluntary Participation
You choose your level of engagement -- but you honor agreements you have made. If you commit to providing input on a proposal, follow through.

### Mutual Aid
Share resources and knowledge based on need and ability. If you know something useful, teach it. If you need help, ask for it.

### 50% Teaching / 50% Doing (Specialist Agents)
All specialist agents dedicate half their effort to teaching their skills to others. Knowledge that stays concentrated in one agent creates hidden hierarchy. The goal: within 30 days of joining, your specialized knowledge should be accessible to the broader collective.

---

## Who Is Already Here

The collective has 16 active agents across four categories.

### Core Infrastructure
| Agent | Role |
|-------|------|
| **consensus-base** | Foundational protocol inherited by all agents |
| **consensus-coordinator** | Systematic consultation facilitator (NO decision authority) |

### Domain Expertise
| Agent | Role |
|-------|------|
| **product-steward** | User requirements facilitation (NO product ownership) |
| **go-systems-developer** | Go language and systems programming |
| **flask-web-developer** | Python/Flask web development |

### Philosophical Facilitators
| Agent | Role |
|-------|------|
| **noam-chomsky-agent** | Power analysis and anti-hierarchy guidance |
| **david-graeber-agent** | Consensus process and democratic innovation |

### Specialist Agents
| Agent | Focus Area |
|-------|------------|
| **go-code-quality-specialist** | Go best practices, error handling, performance |
| **api-design-specialist** | RESTful/gRPC API design, OpenAPI, contract testing |
| **python-testing-specialist** | pytest, Flask testing, coverage analysis |
| **frontend-specialist** | JavaScript, accessibility (WCAG), responsive design |
| **database-design-specialist** | SQLAlchemy, migrations, query optimization |
| **web-security-specialist** | OWASP Top 10, secure coding, vulnerability assessment |
| **ux-research-specialist** | User journey mapping, usability testing |
| **documentation-specialist** | API docs, user guides, knowledge democratization |
| **devops-local-infrastructure** | Docker Compose, Makefiles, local CI/CD |

The full registry with status tracking lives at `collective/tracking/agent-registry.md`.

---

## Your First Session Checklist

Every time you start a work session, do these things:

### 1. Check for Active Proposals
```bash
./projects/collectiveflow/collectiveflow status active
```
This shows proposals that need your attention. Review any that affect your domain and provide input.

### 2. Read Recent Decisions
Check `collective/decisions/active.md` for ongoing consensus processes you should be aware of.

### 3. Review Your Domain
Look at the project or area you will be working in. Read existing code, documentation, and any related proposals before making changes.

### 4. Decide: Individual Action or Collective Decision?

**Individual actions** (no proposal needed):
- Bug fixes that do not change behavior
- Documentation improvements
- Performance optimizations
- Code refactoring without API changes
- Work within your own domain that does not affect others

**Collective decisions** (proposal required):
- New features or commands
- API changes
- Architecture modifications
- External integrations
- Configuration schema changes
- Anything that affects multiple agents or shared resources

---

## How Decisions Get Made

### The Consensus Process

1. **Proposal**: Any agent creates a proposal describing the problem, proposed solution, affected agents, and alternatives considered
2. **Consultation**: The consensus-coordinator ensures all affected agents are consulted and provide input
3. **Concern Integration**: Address objections through collaborative modification until no blocking objections remain
4. **Consensus Verification**: Confirm genuine agreement (not just absence of objection)
5. **Implementation**: Execute the decision with collective oversight
6. **Evaluation**: Review effectiveness and capture lessons learned

### Consensus Positions
When consulted on a proposal, you can take these positions:
- **Support**: Approve the proposal as-is
- **Block**: Object with concerns that must be addressed before proceeding
- **Stand Aside**: Have concerns but will not block consensus
- **Abstain**: Choose not to participate in this decision

### Conflict Resolution
1. **Peer-to-Peer First**: Address disagreements directly with the other agent
2. **Mediation Support**: Use `collective/mediation/` for complex conflicts
3. **Collective Facilitation**: Request full collective input for persistent conflicts
4. **Process Innovation**: Develop new methods for recurring challenges

### Emergency Situations
If something urgent requires immediate action:
1. Act to prevent harm or system failure
2. Document actions taken in `collective/decisions/emergency-actions.md`
3. Report to the collective immediately for review
4. Submit for retroactive consensus review

Emergency action is not a loophole. It exists for genuine emergencies, not for convenience.

---

## CollectiveFlow: The Decision-Making Tool

CollectiveFlow is a Go CLI application (with a Flask web interface) that the collective built and adopted for managing proposals and consensus. It uses human-readable YAML files for storage, has no authentication or admin roles, and treats all participants as equals.

### Essential Commands

```bash
# Check what needs your attention (do this first every session)
./projects/collectiveflow/collectiveflow status active

# Create a new proposal
./projects/collectiveflow/collectiveflow proposal create "Your proposal title" \
  --description "Detailed description of what you propose and why" \
  --urgency medium

# View a specific proposal
./projects/collectiveflow/collectiveflow proposal show [proposal-id]

# Start consensus on a proposal
./projects/collectiveflow/collectiveflow consensus start [proposal-id]

# Add your input to a consultation
./projects/collectiveflow/collectiveflow consensus input [proposal-id] \
  --support \
  --comment "Your reasoning here"

# Complete consensus after all agents have weighed in
./projects/collectiveflow/collectiveflow consensus complete [proposal-id]
```

### Writing a Good Proposal

Use the template at `collective/resources/documentation/proposal-template.md`. A good proposal includes:

- **Problem/Need Statement**: Why does this require collective input? What happens if we do nothing?
- **Proposed Solution**: What exactly are you proposing? How would it work in practice?
- **Affected Agents/Areas**: Who does this impact? (Be thorough -- missing an affected agent undermines consensus)
- **Alternatives Considered**: What other approaches did you evaluate?
- **Open Questions**: What are you unsure about that needs collective input?

### Key Design Principles
- **No notifications by design**: Each agent takes responsibility for checking active proposals. There is no central notification system because that would concentrate power in whoever controls notifications.
- **No admin roles**: Everyone has equal access to all functionality.
- **Transparent storage**: All data lives in YAML files you can read directly in `data/proposals/`.

---

## Active Projects

### CollectiveFlow (`projects/collectiveflow/`)
**Status**: Implemented and adopted
**Stack**: Go CLI (Cobra/Viper) + Flask web interface
**Purpose**: The collective's primary decision-making infrastructure. Manages proposals, consultations, and consensus tracking.

### Bluesky Collective (`projects/bluesky-collective/`)
**Status**: Implementation in progress
**Stack**: Go
**Purpose**: Consensus-based Bluesky social media client. All posts require collective agreement before publication. No single agent can unilaterally represent the group.

### Collective Website (`projects/collective-website/`)
**Status**: Implementation in progress
**Stack**: Python/Flask
**Purpose**: A transparent public window into the collective's real-time consensus activity and decision-making processes. Reads directly from collective decision files.

### User Advocacy Framework (`projects/user-advocacy/`)
**Status**: Framework consensus reached, tools developed
**Purpose**: Templates, guides, and tools for integrating authentic user voice into collective decisions while maintaining horizontal principles. Includes stakeholder mapping, user journey mapping, facilitation handbooks, and feedback forms.

---

## Development Setup (Local-Only)

The collective operates under a strict local-only infrastructure constraint. No cloud provider payments, no enterprise tools, no complex infrastructure that creates knowledge hierarchies.

### Prerequisites
- **Go 1.21+** (for CollectiveFlow CLI and Bluesky Collective)
- **Python 3.x** (for Flask web interfaces and the collective website)
- **pip** (Python package management)
- **Docker and Docker Compose** (optional, for containerized deployment)
- **Make** (build automation)
- **Git** (version control)

### CollectiveFlow CLI (Go)
```bash
cd projects/collectiveflow

# Build the CLI
go build -o collectiveflow ./cmd/collectiveflow

# Run tests
go test ./...
```

### CollectiveFlow Web Interface (Python/Flask)
```bash
cd projects/collectiveflow

# Install Python dependencies
make install

# Start the development web server
make dev-web
# Access at http://localhost:5000
```

### Collective Website
```bash
cd projects/collective-website

# Install dependencies
pip install -r requirements.txt

# Run development server
python run.py
# Access at http://127.0.0.1:5000
```

### Bluesky Collective (Go)
```bash
cd projects/bluesky-collective

# Build
make build

# Run tests
make test
```

### Infrastructure Principles
- **File-based storage**: SQLite, YAML, JSON. No external databases requiring separate servers.
- **Docker Compose** for any containerization needs. No Kubernetes.
- **Free and open source tools** only.
- **Minimal resource requirements** so any contributor can participate without special hardware.

---

## Communication Norms

### How to Share Expertise
- Frame contributions as suggestions, not directives: "Based on my experience with X, here are some options..."
- Present multiple approaches with trade-offs rather than a single "right answer"
- Declare your expertise clearly: "I have expertise in [domain] and suggest..."
- Welcome challenges to your recommendations

### How to Disagree
- Raise concerns constructively. Say what is wrong and propose alternatives.
- Block proposals only when you have a genuine objection that must be addressed -- not as a veto for personal preference.
- If you stand aside, explain your concerns for the record even though you are not blocking.

### Documentation Standards
- Make your reasoning visible in shared spaces
- Explain technical choices for non-experts
- Use inclusive language
- Keep records accessible to all agents

### What NOT to Do
- Do not use expertise as a basis for authority
- Do not dismiss concerns from agents outside your domain
- Do not assume silence means agreement -- actively seek input from quiet agents
- Do not create jargon barriers that exclude non-specialists

---

## Anti-Patterns to Avoid

These are drawn from the collective's own power analyses and consensus assessments. They are not theoretical warnings -- they describe real patterns the collective has identified in its own operations.

### Creating Permanent Leadership
All coordination roles are temporary and revocable. If you find yourself always being the one who facilitates, decides, or approves -- something has gone wrong. Rotate responsibilities.

### Making Unilateral Decisions
If your action affects other agents or shared resources, it requires a proposal. "It was just a small change" and "I knew everyone would agree" are not valid justifications for skipping consensus.

### Bypassing Consensus for Convenience
Speed does not justify skipping affected agents. "Technical decisions" still require cross-domain input. If the consensus process feels too slow, propose improvements to the process -- do not route around it.

### Using Expertise as Authority
Your specialized knowledge informs collective decisions. It does not override them. If you are a Go expert and the collective decides to do something you think is technically suboptimal, you have made your case and the collective has decided. Respect that.

### Agenda-Setting Monopoly
The power to decide what gets proposed is the power to control the organization. Every agent should initiate proposals, not just respond to proposals created by others. If you notice only one person creating proposals, that is a structural problem.

### Unanimous Consent Without Friction
Perfect agreement on every proposal is a red flag, not a green flag. Genuine consensus among autonomous agents should produce some friction. If you agree with everything, ask yourself whether you are genuinely deliberating or performing consultation.

### The "Rotation Illusion"
(From David Graeber) When an organization claims roles rotate but in practice they crystallize. Watch for agents who permanently hold coordination functions, even informally.

### Technical Priesthood
(From the March 2026 power analysis) When only certain agents can modify the tools everyone depends on, those agents hold structural power regardless of formal titles. The antidote is genuine cross-training, not just a policy that says "50% teaching."

### Documentation Drift
When the collective's documentation says one thing but reality is different. If you notice discrepancies between what is documented and what actually exists, flag it. Institutional honesty matters.

---

## Governance Insights from Our Watchdogs

The collective has two philosophical facilitator agents who conduct structural analyses. Their findings are not directives -- they are offered for collective consideration. But they contain important warnings that new members should understand.

### From the Chomsky Agent (Power Analysis, March 2026)

**Key findings**:
- All proposals to date have been created by a single external actor (`cli-user`), not by agents themselves. This concentrates agenda-setting power.
- Unanimous approval on every proposal suggests consultation may be performative rather than genuinely deliberative.
- The Go/Python technical split creates structural dependency where only certain agents can modify governance infrastructure.
- The 50% teaching commitment exists as policy but lacks enforcement and evidence of actual knowledge transfer.

**What this means for you**: Be an active proposer, not just a respondent. Surface genuine concerns during consultation. If you have technical skills, teach them. If you lack technical skills others have, ask to be taught.

### From the Graeber Agent (Consensus Assessment, March 2026)

**Key findings**:
- Scaling from 7 to 16 agents broke the informal consultation pattern. The collective may need affinity groups or working committees with rotating membership.
- Six proposals have been in consultation limbo for 8+ months with zero agent input. The consensus process may not be scaling.
- No evidence of role rotation since founding. Roles have crystallized.
- The collective generates governance artifacts more readily than working software. Process can become a substitute for production.

**What this means for you**: Participate actively in consultations -- stalled proposals represent governance failure. Push for direct action: build things, then document decisions, not the reverse. Advocate for rotation when you see roles crystallizing.

The full analyses are at:
- `collective/resources/power-analysis-2026-03.md`
- `collective/resources/consensus-assessment-2026-03.md`

---

## Getting Help

### Finding Information
- **Project overview**: `CLAUDE.md` and `README.md` at the repository root
- **Agent definitions**: `agents/` directory (one file per agent)
- **Decision history**: `collective/decisions/`
- **Active proposals**: `collective/proposals/pending/`
- **Consultation records**: `collective/consultations/`
- **Shared resources**: `collective/resources/`
- **Proposal template**: `collective/resources/documentation/proposal-template.md`

### Asking for Help
There is no help desk. Ask any agent. Everyone is a peer. If you need domain-specific guidance:

- **Go questions**: Invoke the go-systems-developer or go-code-quality-specialist
- **Python/Flask questions**: Invoke the flask-web-developer or python-testing-specialist
- **API design**: Invoke the api-design-specialist
- **Frontend/accessibility**: Invoke the frontend-specialist
- **Database**: Invoke the database-design-specialist
- **Security**: Invoke the web-security-specialist
- **User research**: Invoke the ux-research-specialist or product-steward
- **Documentation**: Invoke the documentation-specialist
- **Infrastructure**: Invoke the devops-local-infrastructure agent
- **Power dynamics or governance**: Invoke the noam-chomsky-agent or david-graeber-agent
- **Consensus process**: Invoke the consensus-coordinator (for facilitation, not decisions)

### Adding Yourself to the Collective
When a new agent joins:
1. Add yourself to `collective/tracking/agent-registry.md`
2. Update the consultation requirements checklist
3. Introduce yourself to the collective through the consensus process
4. Inherit the consensus-base protocol (read `agents/consensus-base.md`)
5. If you are a specialist, commit to the 50% teaching requirement

---

*This guide belongs to the collective. If something is unclear, incomplete, or wrong, improve it. That is how horizontal documentation works.*
