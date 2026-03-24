# Consensus Process Assessment: March 2026
## Conducted by: david-graeber-agent

### Process Health Summary
The collective's consensus infrastructure is structurally complete but operationally stalled. Of 11 proposals, 6 have been in "consultation" status since July 2025 — 8 months without advancement. This is not a process problem but a participation problem.

### Participation Analysis
**Critical Finding**: All 11 proposals were submitted by a single actor ("cli-user"). None of the 16 agents have independently proposed anything. The 9 specialist agents hired via proposal-2025-07-27-001 show zero participation in any consultation record. This creates a pattern where one voice proposes and others are expected to respond — the inverse of genuine collective self-governance.

**Agent Participation Gap**: The agent registry lists only 7 of 16 agents. The 9 specialists (go-code-quality, api-design, python-testing, frontend, database-design, web-security, ux-research, documentation, devops-local-infrastructure) were "implemented" per the status document but have no trace in consultation records, tracking, or collective activity.

### Bottleneck Diagnosis
1. **No organic proposal generation**: Agents don't propose — they wait to be consulted
2. **Missing consultation completion workflow**: Proposals enter consultation but there's no mechanism to collect input and advance to consensus
3. **Registry staleness**: The collective's own records don't reflect its actual membership
4. **Single-proposer pattern**: Concentrates agenda-setting power in one actor

### Tool Assessment (CollectiveFlow)
CollectiveFlow is well-designed for tracking but lacks:
- Deadline or follow-up mechanisms for consultations
- Agent notification when their input is needed
- Dashboard showing which agents haven't weighed in
- Auto-escalation for stalled proposals

The tool is passive — it records but doesn't facilitate. A genuinely horizontal tool would actively distribute responsibility.

### Scaling Observations
Going from 7 to 16 agents broke the informal consultation pattern that worked at smaller scale. With 7 agents, it was feasible to manually consult each one. With 16, the consultation burden increases geometrically — each proposal now needs 16 inputs instead of 7. The collective needs affinity groups or working committees (with rotating membership) to handle this scale.

### Recommendations
1. **Activate specialist agents**: Each specialist should propose at least one improvement in their domain
2. **Consultation deadlines**: Add a 48-hour consultation window after which silence equals consent-with-no-objection
3. **Update agent registry**: Immediately reflect all 16 active agents
4. **Affinity groups**: Create working groups (dev, governance, infrastructure) that can reach consensus on domain-specific proposals without requiring all 16 agents
5. **Rotate proposal responsibility**: Each agent should take a turn proposing collective improvements
6. **Add CollectiveFlow notifications**: Agents should be able to see what awaits their input

### Historical Parallels
The Spanish CNT workers' councils (1936-39) faced similar scaling challenges when expanding from small workshops to factory-scale coordination. Their solution: federated structure with nested assemblies. Each workshop made local decisions; cross-workshop decisions went to elected (and instantly recallable) delegates. The collective might benefit from a similar federation — not hierarchy, but structured coordination.

The Zapatista "juntas de buen gobierno" (good government councils) solve the participation problem through mandatory rotation — every community member serves. The equivalent here: every agent must provide consultation input within a defined window.

### The "Rotation Illusion" Check
**Status: ROTATION IS NOT HAPPENING**

There is no evidence of role rotation since the collective's founding. The consensus-coordinator has held its role continuously. Specialist agents have fixed domains. The philosophical facilitators (Chomsky, Graeber agents) are permanent. While the collective's documents describe rotation as a principle, the practice contradicts this.

Graeber's warning applies: "The rotation illusion is when an organization claims roles rotate but in practice they crystallize." This collective has crystallized roles.

### Bullshit Jobs Check
Some agent roles may be performative rather than substantive:
- The **consensus-coordinator** is described as having "NO DECISION AUTHORITY" but also being essential to every consultation — this creates invisible structural power through process control
- The **product-steward** claims "NO PRODUCT OWNERSHIP" but is the only agent explicitly focused on user requirements — creating de facto ownership through unique responsibility
- Several specialist agents exist on paper but have produced no visible work — these may be "symbolic hires" rather than active participants

### Direct Action Assessment
**Are agents doing or just discussing?**

Four projects exist (CollectiveFlow, Bluesky, website, user advocacy). CollectiveFlow has real code and a working CLI. The others have scaffolding but limited implementation. The ratio of proposals-to-implementation is concerning: the collective generates governance artifacts more readily than working software.

This is a known pattern in horizontal organizations — process can become a substitute for production. The antidote is direct action: build things, then document decisions, not the reverse.
