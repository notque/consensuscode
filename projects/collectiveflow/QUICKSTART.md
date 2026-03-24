# CollectiveFlow Quick Start

Get CollectiveFlow running in under 5 minutes.

## Choose Your Path

### Path 1: Local Development (Recommended)

**One command setup:**
```bash
make dev-setup
```

**Then start the web interface:**
```bash
make dev-web
```

Open http://localhost:5000 in your browser.

### Path 2: Docker (Optional)

**Two commands:**
```bash
make docker-build
make docker-up
```

Open http://localhost:5000 in your browser.

## What You Get

- **Web interface** at http://localhost:5000
  - View all proposals
  - Create new proposals
  - Track consensus status
  - See collective activity

- **CLI tool** via `./collectiveflow`
  - `./collectiveflow status active` - Show active proposals
  - `./collectiveflow proposal create "Title" --description "Text"` - Create proposal
  - `./collectiveflow consensus start [id]` - Begin consensus process

## Your First Proposal

### Via Web Interface:
1. Go to http://localhost:5000
2. Click "New Proposal"
3. Fill in title and description
4. Submit

### Via CLI:
```bash
./collectiveflow proposal create "My First Proposal" \
  --description "Learning how to use CollectiveFlow" \
  --urgency low
```

## Common Commands

```bash
make help           # See all available commands
make status         # Check current proposals
make dev-web        # Start web interface
make docker-logs    # View Docker logs
make clean          # Remove build artifacts
```

## Data Location

All proposals are stored as YAML files in:
```
./data/proposals/
```

You can read, edit, and backup these files directly. They're human-readable.

## Next Steps

1. Read `DEPLOYMENT.md` for production deployment
2. Check `CLAUDE.md` for development principles
3. Create your first proposal
4. Share with your collective

## Troubleshooting

**Python version error?**
```bash
python3 --version  # Need 3.13+
```

**Port 5000 already in use?**
```bash
# Find what's using it
lsof -i :5000
# Stop that process or edit docker-compose.yml to use different port
```

**Docker not starting?**
```bash
make docker-logs  # Check what went wrong
```

## Philosophy

CollectiveFlow is:
- **Simple** - Files, not databases
- **Transparent** - YAML you can read
- **Horizontal** - No admin users
- **Local-first** - Your machine, your data
- **Optional containers** - Docker is convenience, not requirement

---

**Need help?** Check `DEPLOYMENT.md` or create a proposal for improving documentation.

*Built by consensus, for consensus, through consensus.*
