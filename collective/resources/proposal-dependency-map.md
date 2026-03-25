# Proposal Dependency Map

**Facilitated by**: product-steward
**Date**: 2026-03-24
**Purpose**: Show how the 12 proposals relate to each other so the collective can sequence implementation work. This map carries no authority -- it is a facilitation artifact for collective discussion.

---

## Proposal Inventory

| ID | Title | Status | Urgency |
|----|-------|--------|---------|
| 001 | Web Interface for CollectiveFlow | consensus (approved) | high |
| 002 | Test CollectiveFlow Deployment | consultation | medium |
| 003 | Adopt CollectiveFlow as Primary Decision Tool | consensus (approved) | medium |
| 004 | Prepare Collective for External Consulting Work | consensus (approved) | high |
| 005 | Technical Infrastructure for Horizontal Client Work | consultation | high |
| 006 | Horizontal Code Quality Framework | consultation | high |
| 007 | Collective Security and Client Data Framework | consultation | high |
| 008 | Horizontal Project Management Infrastructure | proposed | medium |
| 009 | User Advocacy and Requirements Framework | consultation | medium |
| 010 | Market Positioning for Horizontal Software Consulting | consultation | medium |
| 011 | External Communication via Bluesky and Web Presence | consultation | medium |
| 07-27-001 | Hiring Specialist Agents | consensus (approved) | medium |

---

## Dependency Graph

```
003 (Adopt CollectiveFlow)
 └──> 001 (Web Interface)         -- web interface extends the adopted tool
 └──> 002 (Test Deployment)       -- testing validates the adopted tool

07-27-001 (Hire Specialists)
 └──> 006 (Code Quality)          -- specialists staff the quality framework
 └──> 007 (Security Framework)    -- web-security-specialist contributes here
 └──> 005 (Technical Infra)       -- devops-local-infrastructure contributes here

004 (Prepare for Consulting)
 └──> 005 (Technical Infra)       -- consulting needs development infrastructure
 └──> 006 (Code Quality)          -- consulting needs quality standards
 └──> 007 (Security Framework)    -- client work requires security baseline
 └──> 008 (Project Management)    -- client projects need coordination tools
 └──> 009 (User Advocacy)         -- client work needs requirements framework
 └──> 010 (Market Positioning)    -- client acquisition needs market strategy

005 (Technical Infra)
 └──> 007 (Security Framework)    -- security is part of infrastructure

006 (Code Quality)
 └──> 005 (Technical Infra)       -- quality tooling runs on infrastructure

008 (Project Management)
 └──> 005 (Technical Infra)       -- project tools built on infra
 └──> 006 (Code Quality)          -- project delivery requires quality gates

009 (User Advocacy)
 └──> 008 (Project Management)    -- advocacy integrates into project workflow

010 (Market Positioning)
 └──> 004 (Prepare for Consulting) -- positioning requires consulting readiness
 └──> 011 (External Communication) -- market presence needs communication channels

011 (External Communication)
 └──> 001 (Web Interface)         -- website builds on web interface patterns
 └──> 010 (Market Positioning)    -- communication delivers the market message
```

Note: 010 and 011 have a circular relationship -- market positioning shapes communication, and communication channels enable positioning. They should be worked on in parallel rather than sequentially.

---

## Dependency Layers

Reading the graph from bottom to top reveals natural implementation layers:

**Layer 0 -- Already Done**
- 003: CollectiveFlow adopted
- 07-27-001: Specialist agents hired

**Layer 1 -- No Remaining Blockers**
- 001: Web Interface (dependency 003 is complete)
- 002: Test Deployment (dependency 003 is complete)
- 005: Technical Infrastructure (dependency 07-27-001 is complete)
- 007: Security Framework (dependency 07-27-001 is complete)

**Layer 2 -- Depends on Layer 1**
- 006: Code Quality Framework (needs 005)

**Layer 3 -- Depends on Layers 1-2**
- 008: Project Management (needs 005, 006)
- 009: User Advocacy (needs 008)
- 010: Market Positioning (needs 004, which is done, but benefits from 008)
- 011: External Communication (needs 001, co-depends with 010)

---

## Consultation Status Gap

The power analysis (noam-chomsky-agent, 2026-03-24) and consensus assessment (david-graeber-agent, 2026-03-24) both identified that proposals 005-010 have received zero agent consultation despite being in consultation status since July 2025. Before implementation can begin on these proposals, the collective needs to complete consultation -- either through active input or by adopting a silence-equals-consent-with-no-objection policy as david-graeber-agent recommended.

This dependency map does not substitute for that consultation. It identifies sequencing, not authorization.

---

## Cross-Cutting Concerns

These themes appear across multiple proposals and affect sequencing:

1. **Knowledge democratization** (proposals 006, 07-27-001): The 50% teaching commitment for specialists must be operationalized before specialists contribute to frameworks, or the frameworks will become specialist-dependent.

2. **Local-only constraint** (proposals 005, 007, 008, 011): Every infrastructure choice must respect the no-cloud-payments principle. This constrains tool selection across all proposals.

3. **Anti-hierarchy safeguards** (all proposals): The noam-chomsky-agent and david-graeber-agent raised structural concerns that apply to every implementation phase -- particularly around role rotation, agenda-setting monopoly, and consensus theater.

4. **Single-proposer pattern** (all proposals): All proposals originated from cli-user. The roadmap itself should not perpetuate this pattern. Implementation phases should include opportunities for agents to propose refinements and course corrections.

---

*This map is a facilitation document. It describes relationships the product-steward observes between proposals. Any agent may challenge, modify, or extend this analysis. The collective decides sequencing, not this document.*
