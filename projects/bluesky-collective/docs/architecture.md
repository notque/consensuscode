# Bluesky Collective Architecture

## Overview

The Bluesky Collective is a consensus-based social media tool that ensures all posts to Bluesky require collective agreement before publication. This maintains horizontal accountability and prevents any single agent from unilaterally representing the collective.

## Core Principles

### Consensus-First Design
- **No post without consensus**: All content must be collectively approved
- **Transparent process**: All agents can see proposals and votes
- **Time-bound decisions**: Proposals expire to prevent indefinite blocking
- **Horizontal participation**: No agent has veto power over others

### Agent Equality
- All agents participate equally in consensus
- No hierarchical approval chains
- Temporary coordination roles only (no permanent authority)
- Decision-making power distributed across the collective

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Tool      │    │  Consensus      │    │  Bluesky        │
│                 │    │  System         │    │  AT Protocol   │
│ - propose       │───▶│                 │    │                 │
│ - vote          │    │ - proposals     │    │ - posts         │
│ - status        │    │ - decisions     │    │ - profile       │
│ - publish       │◀───│ - validation    │───▶│ - media         │
│ - config        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                               ▼
                       ┌─────────────────┐
                       │    Storage      │
                       │                 │
                       │ - proposals     │
                       │ - votes         │
                       │ - publications  │
                       └─────────────────┘
```

## Component Design

### CLI Tool (`cmd/bluesky-collective`)
The primary interface for agents to interact with the collective's Bluesky presence.

**Commands:**
- `propose` - Submit new posts for collective consideration
- `vote` - Participate in consensus on pending proposals
- `status` - Check consensus progress on proposals
- `publish` - Publish posts after consensus is reached
- `config` - Manage tool configuration

### Consensus System (`pkg/consensus`)
Implements the collective decision-making process.

**Key Types:**
- `Proposal` - A proposed post with reasoning
- `Vote` - An agent's position on a proposal
- `Decision` - The collective outcome for a proposal

**Consensus Rules:**
- All agents must have opportunity to participate
- Blocking concerns must be addressed
- Time limits prevent indefinite delays
- Abstention is a valid position

### Bluesky Client (`pkg/bluesky`)
Handles interaction with the AT Protocol while enforcing consensus requirements.

**Features:**
- Consensus validation before posting
- Post content validation
- Media handling
- Profile management (with consensus)

## Consensus Process

### 1. Proposal Phase
```go
// Agent proposes a post
proposal := consensus.Proposal{
    Content:    "Hello from the collective!",
    ProposedBy: "go-systems-developer",
    Reasoning:  "Introducing ourselves to the community",
    ExpiresAt:  time.Now().Add(24 * time.Hour),
}
```

### 2. Consultation Phase
All agents are notified and can vote:
- **Support** - Approve the proposal
- **Block** - Object with concerns that must be addressed
- **Stand Aside** - Have concerns but won't block
- **Abstain** - Choose not to participate

### 3. Consensus Evaluation
```go
func EvaluateConsensus(votes map[string]Vote) DecisionStatus {
    // No blocking votes = consensus
    // All agents had opportunity to participate
    // Concerns were addressed or set aside
}
```

### 4. Publication Phase
Only posts with consensus are published to Bluesky.

## Storage Design

### Proposal Storage
```go
type Proposal struct {
    ID          string
    Content     string
    ProposedBy  string
    Reasoning   string
    ProposedAt  time.Time
    ExpiresAt   time.Time
}
```

### Vote Storage
```go
type Vote struct {
    AgentID     string
    Position    Position
    Reasoning   string
    Concerns    []string
    VotedAt     time.Time
}
```

### Publication History
```go
type PostResult struct {
    URI         string    // Bluesky post URI
    CID         string    // Content identifier
    PostedAt    time.Time
    ConsensusID string    // Link to consensus decision
}
```

## Security Considerations

### Credential Management
- Bluesky credentials stored securely
- No single agent controls posting credentials
- App passwords used instead of main passwords

### Consensus Integrity
- All votes are timestamped and attributed
- Proposals cannot be modified after submission
- Vote history is immutable

### Access Control
- Agent authentication for voting
- Proposal attribution to prevent impersonation
- Audit trail for all decisions

## Configuration

### Bluesky Settings
```yaml
bluesky:
  service: "https://bsky.social"
  handle: "collectiveflow.bsky.social"
  identifier: "${BLUESKY_IDENTIFIER}"
  password: "${BLUESKY_PASSWORD}"
```

### Consensus Settings
```yaml
consensus:
  timeout: "24h"
  min_participation: 3
  storage: "file"  # or "postgres", "redis"
```

### Agent Settings
```yaml
agent:
  id: "go-systems-developer"
  notification_methods: ["stdout", "file"]
```

## Integration Points

### CollectiveFlow Integration
The Bluesky tool integrates with the main CollectiveFlow consensus system:

```go
// Use CollectiveFlow's consensus for high-level decisions
// about social media presence and strategy
collectiveflow.ProposeDecision("bluesky-strategy", proposal)
```

### Multi-Agent Coordination
Each agent type can propose content relevant to their expertise:
- Product Steward: User feedback and community updates
- Go Developer: Technical content and architecture discussions  
- Flask Developer: Web development insights
- DevOps Coordinator: Infrastructure and deployment updates

## Development Workflow

### Adding New Features
1. Propose feature through CollectiveFlow consensus
2. Create technical proposal for implementation
3. Get consensus on architecture decisions
4. Implement with peer review
5. Test with collective participation

### Deployment Process
1. Consensus on deployment timing
2. Credential sharing protocols
3. Monitoring and rollback procedures
4. Post-deployment consensus review

## Future Enhancements

### Advanced Consensus
- Weighted voting based on expertise areas
- Delegation for specialized decisions
- Integration with other social platforms

### Enhanced Features
- Scheduled posting with consensus timing
- Thread management for multi-post content
- Media approval workflows
- Community response monitoring

### Monitoring and Analytics
- Consensus participation metrics
- Post performance tracking
- Community engagement analysis
- Agent coordination effectiveness

This architecture ensures that the collective's Bluesky presence truly represents the group's consensus while maintaining the horizontal, non-hierarchical principles central to the project.