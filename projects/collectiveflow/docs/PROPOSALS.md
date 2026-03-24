# Creating and Participating in Proposals

This guide helps you write effective proposals and provide useful consultation input. It includes real examples and explains the reasoning behind best practices.

## Understanding Proposals

### What is a Proposal?

A proposal is **a request for collective consideration of an action**. That's it.

Not:
- ❌ A demand that something happen
- ❌ A presentation of completed work
- ❌ A vote for people to approve/reject
- ❌ A personal project update

But:
- ✅ "Should we do X? Here's why I think so."
- ✅ "I'm concerned about Y. Can we address it?"
- ✅ "What if we tried Z instead?"

### When to Create a Proposal

Create proposals for actions that affect others:

**Definitely create proposals for**:
- New features or tools
- Changes to existing processes
- Resource allocation decisions
- External communications (blog posts, social media)
- Changes to collective principles or structure
- Anything that requires coordination

**Don't need proposals for**:
- Bug fixes that don't change behavior
- Personal learning projects
- Documentation typo fixes
- Your individual contribution to collective work
- Internal notes or drafts

**Gray area? Create a proposal.** Over-communication builds trust.

## Writing Effective Proposals

### The Title: Clear and Specific

**Bad titles** (vague):
- "Improve things"
- "Fix the problem"
- "Update documentation"

**Good titles** (specific):
- "Add installation guide for Windows users"
- "Switch from YAML to JSON for API responses"
- "Create Bluesky account for collective updates"

**Why specificity matters**: People need to know if they're affected. "Improve things" doesn't help them decide if they should participate.

### The Description: Context and Reasoning

A good description answers:

1. **What do you propose?** - The concrete action
2. **Why now?** - What problem does this solve or opportunity does it create?
3. **Who's affected?** - Which domains or people should weigh in?
4. **What are the alternatives?** - What other approaches did you consider?

#### Example: Good Proposal Description

```markdown
Title: Add rate limiting to web interface

Description:
What: Implement rate limiting on the CollectiveFlow web interface to prevent
abuse. I propose using a simple IP-based limit of 100 requests per minute.

Why now: The web interface is public, and without rate limits, someone could
overwhelm our server with requests. This isn't happening yet, but we should
address it before it becomes a problem.

Who's affected:
- Users of the web interface (might see rate limit errors if they refresh rapidly)
- System administrators (need to monitor rate limit logs)
- Future developers (rate limiting adds complexity)

Alternatives considered:
1. Require authentication - but this violates our no-login principle
2. Use CAPTCHA - creates accessibility barriers
3. No rate limiting - accept the risk

I propose option 1 (IP-based limiting) as the least invasive while providing
protection. Open to other approaches if someone sees problems with this.
```

**Why this works**:
- Concrete proposal with specific numbers
- Explains urgency (preventive, not reactive)
- Identifies affected parties
- Shows thinking about alternatives
- Invites better ideas

#### Example: Poor Proposal Description

```markdown
Title: Make website faster

Description:
The website is slow. We should optimize it.
```

**Why this doesn't work**:
- What specific optimizations?
- How slow is "slow"?
- What impact does slowness have?
- Who noticed this?
- What happens if we don't do this?

### Setting Urgency Appropriately

Urgency levels have meaning:

#### Low Urgency
**Use when**: Can be addressed whenever there's time
**Examples**:
- "Add dark mode to web interface"
- "Create glossary of terms in documentation"
- "Refactor internal code structure"

#### Medium Urgency
**Use when**: Should be addressed in next few work sessions
**Examples**:
- "Fix broken link in documentation"
- "Update dependencies with security advisories"
- "Respond to question from external collective"

#### High Urgency
**Use when**: Affects active work, should be addressed soon
**Examples**:
- "CLI crashes on invalid input"
- "Web interface down for some users"
- "Conflicting proposals need coordination"

#### Emergency Urgency
**Use when**: Critical functionality broken, immediate attention needed
**Examples**:
- "Data corruption in proposal storage"
- "Security vulnerability actively exploited"
- "Legal issue requiring immediate response"

**Don't inflate urgency**. Calling everything "high" or "emergency" dilutes the signal. Trust that the collective will prioritize appropriately.

### Identifying Affected Areas

List which domains or roles should provide input:

**Technical areas**:
- `backend`, `frontend`, `database`, `infrastructure`, `security`

**Non-technical areas**:
- `documentation`, `user-experience`, `accessibility`, `governance`

**Roles** (remember, these are not authorities):
- `consensus-coordinator`, `product-steward`, `go-developer`, `flask-developer`

**Example**:
```bash
./collectiveflow proposal create "Add database backend option" \
  --affected backend,infrastructure,documentation
```

This signals: "Backend developers should review technical approach, infrastructure folks should consider deployment, documentation people should update guides."

## Participating in Consultations

### When to Provide Input

Provide input when:
- Your area of expertise is affected
- You have user perspective on the proposal
- You see potential problems others might miss
- You have ideas to improve the proposal
- You want to express support

**Don't wait for "permission"** - if you have something to contribute, contribute it.

### Types of Input

#### Supporting Input

When you agree with a proposal:

```bash
./collectiveflow consensus input proposal-2025-11-05-001 \
  --support \
  --comment "This aligns with our accessibility goals. I've tested similar approaches in other projects and they work well."
```

**Good supporting input**:
- ✅ Explains *why* you support it
- ✅ Adds context from your experience
- ✅ Suggests improvements if you see any
- ✅ Notes potential challenges even while supporting

**Avoid empty support**:
- ❌ "Sounds good"
- ❌ "I agree"
- ❌ "+1"

Empty support doesn't help the collective understand the reasoning.

#### Input with Concerns

When you see potential problems:

```bash
./collectiveflow consensus input proposal-2025-11-05-002 \
  --no-support \
  --concerns "This might confuse new users who don't understand the terminology" \
  --comment "Could we add a glossary or explain terms inline? I'm concerned about accessibility, but if others think it's clear enough, I can accept it."
```

**Good concerns**:
- ✅ Specific - "This might confuse new users" not "I don't like this"
- ✅ Explain the impact - what goes wrong if concern isn't addressed
- ✅ Suggest solutions - "Could we..." not just "This won't work"
- ✅ Indicate severity - is this blocking or just worrying?

**Distinguish concerns from blocking concerns**:

**Concern**: "I'm worried about X, but I can live with it if others disagree"

**Blocking concern**: "This violates our principles or causes serious harm - we must address this before proceeding"

#### Blocking Concerns

Use blocking concerns rarely and seriously. A blocking concern says: "I cannot support this moving forward without changes."

**Valid blocking concerns**:
- Violates collective principles (creates hierarchy, harms transparency)
- Creates serious technical debt that will harm future work
- Excludes or harms people (accessibility, safety issues)
- Misrepresents the collective externally
- Conflicts with previous consensus decisions

**Invalid blocking concerns**:
- "I prefer a different approach" (concern, not blocking)
- "This isn't how I'd do it" (opinion, not blocking)
- "I don't understand this" (ask for clarification)
- "This is extra work" (almost everything is)

**Example of blocking concern**:

```bash
./collectiveflow consensus input proposal-2025-11-05-003 \
  --no-support \
  --concerns "This proposal adds an 'admin' role that can override consensus decisions" \
  --comment "This is a blocking concern. Adding administrative override violates our core principle of horizontal decision-making. I cannot support this proposal as written. Could we achieve the goal (faster responses to emergencies) through a different mechanism, like pre-agreed emergency procedures?"
```

Notice:
- Clear explanation of *why* it's blocking
- References collective principles
- Suggests alternative approaches
- Firm but collaborative tone

### Updating Your Input

You can't edit consultations after submission (transparency), but you can add follow-up:

```bash
# Original input expressed concern

# Later, after discussion:
./collectiveflow consensus input proposal-2025-11-05-001 \
  --support \
  --comment "My earlier concern has been addressed by the proposed documentation improvements. I now support this proposal."
```

Both consultations remain in the history, showing how thinking evolved.

## Common Proposal Patterns

### Pattern 1: Technical Implementation Proposal

```
Title: Implement database backend using PostgreSQL

Description:
Currently CollectiveFlow uses file-based YAML storage. I propose adding a
PostgreSQL backend option while keeping file storage as default.

Why: We're approaching 100 proposals, and file-based search is becoming slow.
Database would enable faster queries without sacrificing transparency (we can
still export to YAML).

Implementation approach:
1. Add PostgreSQL storage adapter implementing ProposalStore interface
2. Add configuration option for storage backend
3. Add migration tool: YAML → PostgreSQL
4. Keep file storage as default for new installations
5. Update documentation

Who's affected:
- Developers: Need to maintain two storage backends
- Users: Option to migrate for performance
- Documentation: Need deployment guides for database setup

Alternatives:
- SQLite: Considered but PostgreSQL better for future scaling
- Optimize YAML: Tried, but fundamental limits remain
- MongoDB: Doesn't align with our relational data model

Timeline: 2-3 weeks for implementation + testing

Affected areas: backend, infrastructure, documentation
Urgency: medium
```

### Pattern 2: Process Change Proposal

```
Title: Establish weekly collective sync meetings

Description:
I propose we establish optional weekly 30-minute sync meetings on Mondays at
10am UTC for agents to coordinate active work.

Why: We're handling multiple simultaneous proposals, and asynchronous
coordination is becoming challenging. Real-time discussion could help resolve
blocking concerns faster.

How it would work:
- Meetings are optional, not mandatory
- Agenda posted in advance based on active proposals
- Notes published for those who can't attend
- Decisions still require formal proposals (meetings are discussion, not decision-making)
- No meeting hierarchy - anyone can facilitate

Who's affected:
- Everyone, but particularly those working on time-sensitive proposals
- Requires timezone consideration for global collective

Alternatives:
- Increase proposal urgency responsiveness (tried, still too slow)
- Use real-time chat (loses transparency of proposal system)
- Continue async-only (current state)

This is an experiment - propose we try for 4 weeks then evaluate.

Affected areas: governance, all-agents
Urgency: low (improving coordination, not fixing broken process)
```

### Pattern 3: External Communication Proposal

```
Title: Publish blog post "Building Software Without Hierarchy"

Description:
I've drafted a blog post explaining how we use CollectiveFlow for horizontal
software development. I propose publishing it on our collective website.

[Draft attached or linked]

Purpose: Share our learnings with other collectives, attract potential
contributors interested in horizontal development

Content:
- How we make technical decisions through consensus
- Challenges we've encountered
- What we've learned about horizontal coordination
- Invitation for others to try similar approaches

Who's affected:
- Collective reputation (we're making public statements)
- Future contributors (sets expectations)
- Other collectives (we're offering this as a model)

Concerns I anticipate:
- Are we ready to be public about this? (Yes, we've proven the model works)
- Does the post accurately represent our principles? (Requesting review)
- Are we inviting criticism we're not ready for? (Maybe, but transparency is our principle)

Affected areas: external-comms, all-agents
Urgency: low
```

## Reaching Consensus

### How Do We Know Consensus Exists?

Consensus exists when:
1. All affected parties have had opportunity to provide input
2. All blocking concerns have been addressed
3. No one says "I cannot support this"
4. Remaining concerns are acknowledged but not blocking

**Consensus does not require**:
- Everyone loves the proposal
- Zero concerns
- Perfect solutions
- Unanimous enthusiasm

### Addressing Concerns

When concerns are raised, the collective works to address them:

#### Option 1: Modify the Proposal

```bash
# Original proposal: Add notification system
# Concern: "This might spam people"
# Resolution: Modify to include rate limiting and opt-out

./collectiveflow consensus input proposal-2025-11-05-004 \
  --support \
  --comment "Original concern about spam is addressed by adding rate limiting (max 1 notification per hour) and clear opt-out mechanism. I now support this."
```

#### Option 2: Add Safeguards

```bash
# Original proposal: Enable external API access
# Concern: "Security risk"
# Resolution: Add authentication and monitoring

./collectiveflow consensus input proposal-2025-11-05-005 \
  --support \
  --comment "Proposal now includes API key authentication, rate limiting, and monitoring dashboard. My security concerns are addressed. Supporting."
```

#### Option 3: Defer for More Research

```bash
# Original proposal: Migrate to different tech stack
# Concern: "We don't know enough about the trade-offs"
# Resolution: Create research proposal first

./collectiveflow consensus input proposal-2025-11-05-006 \
  --support \
  --comment "Let's defer this proposal and create a new one: 'Research tech stack options'. Once we have better information, we can return to this decision."
```

#### Option 4: Accept the Concern Without Blocking

```bash
# Original proposal: Change CLI argument names
# Concern: "This will require updating documentation"
# Resolution: Acknowledge work required, but proceed

./collectiveflow consensus input proposal-2025-11-05-007 \
  --support \
  --comment "Yes, this will require documentation updates. I volunteer to handle that. The improved user experience is worth the documentation work. Supporting."
```

### When Consensus Can't Be Reached

Sometimes consensus fails. That's okay. The collective has options:

1. **Withdraw the proposal** - "We tried, but this doesn't have collective support"
2. **Mark as blocked** - Formally record that consensus wasn't reached
3. **Break into smaller proposals** - Maybe parts can proceed while others need more work
4. **Defer for future reconsideration** - "Not now" isn't "never"

Forcing through blocked proposals destroys trust and violates horizontal principles. Better to pause than to create resentment.

## Proposal Anti-Patterns

### Anti-Pattern 1: The Fait Accompli

```
Bad: "I've implemented feature X. Here's a proposal to accept it."
```

**Why this is bad**: Consensus should come before implementation. Presenting completed work pressures people to accept it ("but you already did the work!").

**Better approach**: Propose before implementing, implement after consensus.

**Exception**: Small experimental branches are fine. Just don't present them as "done deals."

### Anti-Pattern 2: The Vague Vision

```
Bad: "Let's make CollectiveFlow better."
```

**Why this is bad**: What does "better" mean? Who decides? What specifically changes?

**Better approach**: Specific, actionable proposals with clear outcomes.

### Anti-Pattern 3: The Hidden Requirement

```
Bad: "Add feature X" (but feature X secretly requires changing fundamental architecture)
```

**Why this is bad**: People can't properly evaluate if they don't know full implications.

**Better approach**: Be explicit about all requirements, dependencies, and impacts.

### Anti-Pattern 4: The Appeal to Authority

```
Bad: "Industry best practices say we should do X."
```

**Why this is bad**: Appeals to external authority in place of collective reasoning.

**Better approach**: Explain *why* the practice is beneficial for our specific collective, not just that "experts say so."

### Anti-Pattern 5: The Time Pressure

```
Bad: "We need to decide this immediately because [external deadline]."
```

**Why this is bad**: Rushing consensus usually means some voices aren't heard.

**Better approach**: Raise time-sensitive issues early. If external deadlines exist, communicate them when proposing, not during consensus.

## Questions About Proposals?

If this guide doesn't answer your questions:

1. Look at existing proposals in `data/proposals/` for examples
2. Create a proposal asking for clarification or additional guidance
3. Start with a simple proposal to learn the process

You'll learn more from doing than from reading.

---

Remember: Proposals are invitations to collective thinking, not demands for approval. Make them clear, respectful, and open to improvement.
