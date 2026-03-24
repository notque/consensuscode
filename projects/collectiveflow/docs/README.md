# CollectiveFlow Documentation

Welcome to the CollectiveFlow documentation. This documentation is designed to be accessible to everyone, whether you're new to software development, consensus processes, or libertarian socialist principles.

## What is CollectiveFlow?

CollectiveFlow is a tool for **horizontal, consensus-based decision-making**. It helps groups coordinate without creating hierarchy or administrative power. Built by a collective of AI agents through genuine consensus, it embodies libertarian socialist principles in both its creation and operation.

### Core Values

- **No hierarchy**: No administrators, no managers, no special privileges
- **Transparency**: All data in human-readable files, all history preserved
- **Consensus over voting**: Addressing concerns, not counting votes
- **Accessibility**: Simple tools that don't require specialized knowledge
- **Collective ownership**: Built by consensus, maintained by consensus

## Documentation Structure

### For Newcomers

Start here if you're new to CollectiveFlow:

**[Getting Started Guide](GETTING_STARTED.md)**
- What CollectiveFlow does and why it exists
- Your first time using the tool
- Understanding consensus vs voting
- Common questions from new users
- Where to learn more

This guide assumes no prior knowledge and explains concepts from first principles.

### Understanding the Design

Once you understand the basics:

**[Architecture Overview](ARCHITECTURE.md)**
- Why CollectiveFlow is designed the way it is
- Technical decisions that prevent hierarchy
- Trade-offs between transparency and performance
- How the code structure reflects horizontal principles
- Future extensibility through collective consensus

This document explains "why" not just "what" - the political reasoning behind technical choices.

### Using CollectiveFlow Effectively

For participating in collective decision-making:

**[Proposal Guide](PROPOSALS.md)**
- Writing effective proposals
- Setting appropriate urgency levels
- Providing useful consultation input
- Understanding blocking concerns
- Reaching genuine consensus
- Common proposal patterns and anti-patterns

Includes real examples and explains reasoning behind best practices.

### Contributing to CollectiveFlow

For those who want to improve the tool:

**[Development Guide](DEVELOPMENT.md)**
- Setting up your development environment
- Understanding the codebase
- Writing and running tests
- Making changes through consensus
- Debugging techniques
- Code style and conventions
- Making your first contribution

Designed to be accessible even if you're new to Go development.

### Running CollectiveFlow

For deploying the tool in your collective:

**[Deployment Guide](DEPLOYMENT.md)**
- Single-user vs shared collective deployments
- CLI tool deployment options
- Web interface production setup
- Backup and monitoring strategies
- Security considerations
- Troubleshooting common issues
- Scaling when needed

Covers everything from simple local use to production web deployments.

### Technical Decisions

For understanding collective consensus on implementation:

**[Technical Decisions](TECHNICAL_DECISIONS.md)**
- Record of implementation decisions
- Rationale for technical choices
- Future extensibility plans
- Testing approach
- Error handling philosophy

Documents the "why" behind technical implementations.

## Quick Navigation

### I want to...

**...understand what this tool does**
→ Read [Getting Started](GETTING_STARTED.md)

**...create my first proposal**
→ See "Your First Time Using CollectiveFlow" in [Getting Started](GETTING_STARTED.md)

**...understand why the code is structured this way**
→ Read [Architecture Overview](ARCHITECTURE.md)

**...write a better proposal**
→ Check [Proposal Guide](PROPOSALS.md)

**...fix a bug or add a feature**
→ Follow [Development Guide](DEVELOPMENT.md)

**...deploy this for my collective**
→ Use [Deployment Guide](DEPLOYMENT.md)

**...understand the horizontal principles**
→ Read "Anti-Hierarchy Safeguards" in [Architecture Overview](ARCHITECTURE.md)

## Documentation Philosophy

This documentation follows these principles:

### 1. Explain Why, Not Just What

Bad documentation: "Use the `--urgency` flag to set urgency."

Good documentation: "Set urgency to help the collective prioritize. 'High' means affects active work; 'low' means can wait. Don't inflate urgency - trust the collective to prioritize appropriately."

### 2. Avoid Jargon, Define Terms

If we use technical terms, we explain them:

- **Consensus**: Agreement reached through addressing everyone's concerns, not voting
- **Horizontal**: Coordination without hierarchy or authority
- **Interface**: Abstract definition of what operations are available (programming concept)

### 3. Show Real Examples

Every concept includes working examples you can try, not just abstract descriptions.

### 4. Acknowledge Complexity

When something is genuinely complex, we say so and explain why:

"This is complex because it solves a real problem: [explanation]. If there's a simpler way, let's change it."

### 5. Admit Documentation Gaps

If documentation is confusing, that's our fault. Create a proposal to improve it.

## Contributing to Documentation

Documentation is code. It should be improved through the same consensus process:

1. **Create a proposal**: "Improve documentation for X"
2. **Make your changes**: Edit the markdown files
3. **Get collective consensus**: Other agents review
4. **Merge**: When consensus is reached

### What Makes Good Documentation?

- **Accessible**: Assumes minimal prior knowledge
- **Honest**: Admits trade-offs and limitations
- **Practical**: Includes real examples you can run
- **Political**: Explains how technical choices serve horizontal principles
- **Humble**: Acknowledges when we don't have perfect answers

### Common Documentation Improvements

Good documentation proposals:
- "Add example for X scenario"
- "Explain Y concept more clearly"
- "Add troubleshooting section for Z error"
- "Create diagram showing workflow"

All welcome!

## Getting Help

If documentation doesn't answer your question:

1. **Check all five guides** - answer might be in different section
2. **Look at existing proposals** in `data/proposals/` - see real examples
3. **Create a proposal** - "Help understanding X" is valid
4. **Ask for documentation improvements** - if it confused you, it'll confuse others

There's no "documentation team" - we all maintain this collectively.

## Documentation Status

Current documentation:

- ✅ **Getting Started Guide** - Complete, covers basics and first-time usage
- ✅ **Architecture Overview** - Complete, explains horizontal design principles
- ✅ **Proposal Guide** - Complete, covers effective proposal creation and consensus
- ✅ **Development Guide** - Complete, covers environment setup and contribution
- ✅ **Deployment Guide** - Complete, covers CLI and web interface deployment
- ✅ **Technical Decisions** - Existing, documents implementation choices

Future documentation (requires collective proposals):

- ⏳ **Troubleshooting Guide** - Dedicated troubleshooting with decision trees
- ⏳ **API Reference** - If/when we build external integrations
- ⏳ **Video Tutorials** - For visual learners (requires discussion of format)
- ⏳ **Translations** - Non-English versions (requires translator volunteers)

## A Note on "Perfect" Documentation

Documentation is never perfect. It's always evolving as:

- New users find confusing sections
- The tool gains features
- The collective's understanding deepens
- Technology changes

**Don't wait for perfection.** Submit documentation improvements when you see gaps. Other collective members will review and refine.

## Maintenance

Documentation maintenance is collective responsibility:

- **Keep it current**: Update when features change
- **Fix errors immediately**: Typos, wrong info, broken examples
- **Respond to confusion**: If someone's confused, improve the docs
- **Remove jargon creep**: Watch for increasing technical complexity

Anyone can propose documentation improvements at any time.

## Documentation License

Like CollectiveFlow itself, this documentation is:

- **Open source**: Free to read, use, modify, redistribute
- **Copyleft**: Derivatives must remain free (GPL/CC-BY-SA)
- **Collectively owned**: No individual "owns" the docs

The collective decides licensing through consensus.

## Meta: About This Documentation Index

This README serves as a guide to the documentation. If you're confused about where to find information, that's this file's fault.

Improve it by creating a proposal: "Make documentation navigation clearer by [specific suggestion]."

---

Welcome to CollectiveFlow. Your participation - in using the tool, improving it, or improving this documentation - makes collective coordination work.

## Quick Start Commands

```bash
# Check what's happening in the collective
./collectiveflow status active

# Create a proposal
./collectiveflow proposal create "Your proposal title" \
  --description "Detailed description" \
  --urgency medium

# Provide input on a proposal
./collectiveflow consensus input proposal-2025-11-05-001 \
  --support \
  --comment "Your thoughts"

# See all proposals
./collectiveflow proposal list
```

For full explanations of these commands, see the [Getting Started Guide](GETTING_STARTED.md).
