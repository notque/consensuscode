# Getting Started with CollectiveFlow

Welcome to CollectiveFlow - a tool for horizontal, consensus-based decision-making. This guide will help you understand what the tool does, why it exists, and how to start participating.

## What is CollectiveFlow?

CollectiveFlow is a decision-making tool built on **libertarian socialist principles**. It helps groups coordinate without creating hierarchy or administrative power.

### Why Not Just Use Email or Chat?

Good question. Here's what makes CollectiveFlow different:

- **Structured consensus tracking**: Everyone's input is recorded, not lost in chat scrollback
- **No implicit hierarchy**: Email threads often privilege whoever "started" the conversation
- **Transparent history**: Complete record of how decisions were made and why
- **Equal participation**: No one has special "admin" powers to override decisions
- **Systematic consultation**: Tools to ensure all affected people are consulted

### Key Principle: No One is in Charge

This is not a project management tool. There are no "project managers" or "administrators" in CollectiveFlow. Everyone participates as equals, and decisions require genuine consensus - not majority votes or manager approval.

## Understanding the Workflow

CollectiveFlow supports a simple process:

```
1. Someone proposes an action
   ↓
2. Everyone affected provides input
   ↓
3. Concerns are addressed collectively
   ↓
4. Decision is made when no blocking concerns remain
   ↓
5. Implementation happens with continued transparency
```

### What Makes This "Consensus"?

**Consensus is not voting**. In voting, 51% can force the other 49% to comply. In consensus:

- Everyone affected by a decision participates
- All concerns are heard and addressed
- The goal is finding solutions everyone can live with
- One blocking concern stops progress until it's resolved
- No one can override concerns through authority

## Your First Time Using CollectiveFlow

### Step 1: Check What's Happening

Before you do anything else, see what the collective is working on:

```bash
./collectiveflow status active
```

This shows you:
- Proposals currently under consideration
- Which ones need your input
- What stage each proposal is in

**Why this matters**: You can't participate if you don't know what's happening. Checking status regularly is how everyone stays informed in a horizontal system.

### Step 2: Create Your First Proposal

Let's say you want to suggest a documentation improvement. Here's how:

```bash
./collectiveflow proposal create "Improve installation documentation" \
  --description "The current installation docs assume people know what a virtual environment is. Let's add more context for newcomers." \
  --urgency low \
  --affected documentation,newcomers
```

**What just happened?**

- You created a proposal with a unique ID (like `proposal-2025-11-05-001`)
- The proposal is now in "proposed" status
- It's stored in a YAML file anyone can read at `data/proposals/`
- The collective can now consider it

**Understanding urgency levels**:
- `low`: Can be addressed whenever there's time
- `medium`: Should be considered in the next few work sessions
- `high`: Needs attention soon, affects active work
- `emergency`: Blocking critical functionality, address immediately

### Step 3: Start the Consensus Process

Once you've created a proposal, start gathering input:

```bash
./collectiveflow consensus start proposal-2025-11-05-001
```

This changes the proposal status to "consultation" - signaling that input is being actively gathered.

### Step 4: Provide Your Input on Someone Else's Proposal

When you see a proposal that affects your area, add your input:

```bash
# If you support it
./collectiveflow consensus input proposal-2025-11-05-002 \
  --support \
  --comment "This aligns with our accessibility principles"

# If you have concerns
./collectiveflow consensus input proposal-2025-11-05-002 \
  --no-support \
  --concerns "This might create a learning curve for CLI users" \
  --comment "Can we provide both options?"
```

**Understanding concerns vs blocking concerns**:

- **Concern**: "I'm worried about X, but if others think it's okay, I can live with it"
- **Blocking concern**: "This violates our principles or causes serious problems - we need to address this before proceeding"

The tool records both, but the collective determines which concerns block consensus.

### Step 5: Check Consensus Status

See where a proposal stands:

```bash
./collectiveflow consensus status proposal-2025-11-05-001
```

This shows:
- Who has provided input
- What concerns were raised
- Whether consensus has been reached

### Step 6: Complete Consensus

When all concerns are addressed and everyone supports the proposal:

```bash
./collectiveflow consensus complete proposal-2025-11-05-001
```

The proposal moves to "consensus" status, meaning the collective has agreed to proceed.

## Common Questions from New Users

### "Who decides if we have consensus?"

Everyone who participates. The tool tracks input, but humans decide together whether concerns are addressed and consensus exists.

### "What if someone blocks everything?"

Consensus requires good faith participation. If someone consistently blocks without engagement, that's a collective relationship problem, not a tool problem. The solution is conversation, not giving anyone power to override blocking concerns.

### "How do I know when to provide input?"

Check `./collectiveflow status active` regularly. If you're affected by a proposal or have relevant expertise, your input helps the collective make better decisions.

### "What if I disagree with a decision?"

Consensus doesn't mean you love every decision. It means:
- You were consulted
- Your concerns were heard
- The final solution addresses serious problems
- You can live with the outcome

If you have a blocking concern, express it clearly and work with others to find a solution.

### "Can someone delete or change decisions?"

No. CollectiveFlow maintains a complete history. Proposals can be withdrawn or modified, but the history of what happened and why is permanent.

### "What if I make a mistake in a proposal?"

You can't edit proposals after creation, but you can:
- Withdraw the proposal and create a new one
- Add clarifications through the consultation process
- Ask others to help refine the proposal

This permanence might feel restrictive, but it serves transparency - everyone can see how proposals evolved.

## Where to Learn More

Now that you understand the basics:

1. **[Architecture Overview](ARCHITECTURE.md)**: Learn why CollectiveFlow is designed the way it is
2. **[Proposal Guide](PROPOSALS.md)**: Deep dive into creating effective proposals
3. **[Development Guide](DEVELOPMENT.md)**: Set up your development environment
4. **[Deployment Guide](DEPLOYMENT.md)**: Run CollectiveFlow in your collective

## Getting Help

CollectiveFlow is maintained by a horizontal collective. If you have questions:

1. Check the documentation in `docs/`
2. Look at existing proposals in `data/proposals/` to see examples
3. Create a proposal asking for clarification or documentation improvements
4. The collective will respond through the consensus process

Remember: There's no "help desk" or "support team" because there's no hierarchy. Help comes from the collective, and you're part of that collective.

## A Note on Jargon

We've tried to avoid jargon, but some terms are important:

- **Consensus**: Agreement reached through addressing everyone's concerns
- **Horizontal**: Coordination without hierarchy or authority
- **Consultation**: Gathering input from affected people
- **Blocking concern**: A concern serious enough to stop progress until addressed
- **Libertarian socialism**: Political philosophy emphasizing freedom and equality without imposed authority

If you encounter terms that aren't clear, that's a documentation problem we should fix through consensus.

---

Welcome to the collective. Your participation makes collective decision-making work.
