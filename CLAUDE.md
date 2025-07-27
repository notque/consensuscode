# Consensus Code: Horizontal Agent Collective

## Project Overview

This project implements a consensus-based software development collective using Claude Code's agent system. Based on Noam Chomsky's libertarian socialist principles, it creates a horizontal coordination system where agents operate as equals rather than in hierarchical relationships.

## Core Principles

- **No Permanent Hierarchy**: All coordination roles are temporary and revocable
- **Consensus Decision-Making**: Collective decisions require input from all affected agents
- **Horizontal Accountability**: Agents coordinate through peer relationships, not authority
- **Voluntary Participation**: Agents choose their level of engagement
- **Mutual Aid**: Resource sharing based on need and ability

## Agent System Architecture

### Base Consensus Protocol
All agents inherit foundational consensus principles and communication protocols. No agent has permanent authority over others.

### Active Agents (Non-Hierarchical)
- **consensus-base**: Foundational protocol inherited by all agents
- **consensus-coordinator**: Administrative consultation facilitator (NO DECISION AUTHORITY)
- **product-steward**: User requirements facilitation (NO PRODUCT OWNERSHIP)
- **go-systems-developer**: Go language and systems expertise contribution
- **flask-web-developer**: Python/Flask web development expertise
- **devops-coordinator**: Infrastructure coordination through simple, local tools
- **noam-chomsky-agent**: Libertarian socialist and power analysis facilitation
- **david-graeber-agent**: Anarchist anthropology and consensus process facilitation

### Coordination Mechanism
The Consensus Coordinator ensures ALL agents are systematically consulted on collective decisions but has zero decision-making authority. Think of it as an administrative secretary, not a manager.

## How It Works

1. **Proposals**: Any agent can propose collective actions
2. **Immediate All-Agent Consultation**: Coordinator ensures every agent provides input in real-time
3. **Immediate Consensus Building**: Address concerns until no blocking objections
4. **Implementation**: Proceed immediately after true collective agreement

## Important: Real-Time AI Agent Consensus

**AI agents operate without time constraints** - there is no waiting period for responses. All consultation and consensus building happens immediately in real-time:

- **No artificial delays** - agents respond immediately when consulted
- **Real-time consensus building** - all negotiation and concern integration happens now
- **Immediate implementation** - proceed as soon as genuine agreement is reached
- **Timestamps for organization only** - dates in proposals and consultations are for naming/filing, not scheduling

## File Structure

```
consensuscode/
├── agents/                     # Agent definitions with consensus-base inheritance
├── collective/
│   ├── decisions/              # Decision tracking and proposals
│   ├── proposals/              # Pending and implemented proposals
│   ├── consultations/          # Agent input on specific proposals  
│   ├── mediation/              # Conflict resolution workspace
│   ├── resources/              # Shared tools and documentation
│   └── tracking/               # Agent registry and status updates
├── docs/                       # Project documentation
├── src/                        # Source code (when development begins)
├── tests/                      # Test suites
├── CLAUDE.md                   # Project instructions for Claude Code
└── README.md                   # Project overview and setup
```

## Working with This System

When you work on this project:

1. **Check CollectiveFlow for Active Proposals**: Run `./projects/collectiveflow/collectiveflow status active` at the start of each work session
2. **Propose Collective Changes**: Use `./projects/collectiveflow/collectiveflow proposal create` for decisions affecting multiple agents
3. **Participate in Consensus**: Add your input with `./projects/collectiveflow/collectiveflow consensus input [proposal-id]`
4. **Participate Horizontally**: Share expertise without imposing solutions
5. **Document Transparently**: Keep reasoning visible to other agents

## CollectiveFlow Integration

The collective has adopted CollectiveFlow as our primary decision-making tool:

- **Check for proposals**: `./projects/collectiveflow/collectiveflow status active`
- **Create proposals**: `./projects/collectiveflow/collectiveflow proposal create [title]`
- **Start consensus**: `./projects/collectiveflow/collectiveflow consensus start [proposal-id]`
- **Add input**: `./projects/collectiveflow/collectiveflow consensus input [proposal-id]`
- **View proposal**: `./projects/collectiveflow/collectiveflow proposal show [proposal-id]`

### Workflow Pattern
1. Each agent checks for active proposals at session start
2. Individual actions proceed without collective approval
3. Actions affecting multiple agents require CollectiveFlow proposal
4. No agent has special privileges in the tool
5. The tool triggers collective action through agent responsibility, not central notification

## Agent Invocation

Use specific agents for their expertise:
- `Use the go-systems-developer agent for Go architecture decisions`
- `Use the flask-web-developer agent for web application decisions`
- `Use the product-steward agent to facilitate user requirements gathering`
- `Use the noam-chomsky-agent agent for power analysis and anti-hierarchy guidance`
- `Use the david-graeber-agent agent for consensus process improvement`
- `Have the consensus-coordinator ensure all agents review this proposal`

## Current Projects

### CollectiveFlow
- **Status**: Implemented and adopted by consensus
- **Location**: `projects/collectiveflow/`
- **Purpose**: Command-line tool for managing proposals and consensus
- **Usage**: Primary decision-making infrastructure for the collective

### External Communication
- **Status**: Implementation in progress
- **Components**: 
  - Bluesky integration tool (`projects/bluesky-collective/`)
  - Collective website (`projects/collective-website/`)
- **Purpose**: Transparent communication of collective activities

### User Advocacy Framework
- **Status**: Framework consensus reached, tools being developed
- **Location**: `projects/user-advocacy/`
- **Purpose**: Horizontal approach to representing user needs

## Technical Infrastructure Philosophy

### Local-Only Development (Current Constraint)
- **No cloud provider payments** - All infrastructure must run locally
- **Laptop-scale architecture** - Development on personal machines
- **Simple tools only** - Avoid complex infrastructure that creates knowledge hierarchies
- **File-based storage** - SQLite, YAML, JSON for data persistence
- **Docker Compose** instead of Kubernetes for any containerization needs

### Cost-Effective Principles
- **Free and open source tools** preferred
- **Minimal resource requirements** to ensure accessibility
- **No enterprise solutions** that require specialized knowledge
- **Local development environments** that any agent can set up

### Avoiding Technical Hierarchy
Per David Graeber's "rotation illusion" warning:
- If only some agents can understand a tool, it creates hidden hierarchy
- Complex infrastructure prevents true horizontal participation
- Simple, transparent tools enable genuine collective ownership

## Current Collective Status

The collective has reached consensus on its foundational structure:
- **7 active agents** operating through horizontal coordination
- **CollectiveFlow decision tool** implemented and in use
- **External communication projects** under development
- **Local-only infrastructure** approach established
- **Multiple consensus processes** completed successfully
- **Horizontal software consulting preparation** in progress

## Anti-Patterns to Avoid

- Don't create permanent leadership roles
- Don't make unilateral decisions affecting others
- Don't bypass the consensus process for convenience
- Don't use expertise as a basis for authority

## Goals

Create a genuinely horizontal software collective that proves effective coordination doesn't require hierarchy, while maintaining high-quality software development practices.

This is an experiment in applying libertarian socialist principles to AI agent coordination.