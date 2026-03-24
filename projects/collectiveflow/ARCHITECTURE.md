# CollectiveFlow Architecture

## Overview

CollectiveFlow uses a **simple, file-based architecture** that maintains horizontal principles and avoids complex infrastructure.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CollectiveFlow System                     │
└─────────────────────────────────────────────────────────────┘

┌──────────────┐          ┌──────────────┐          ┌──────────────┐
│   CLI Tool   │          │ Web Interface│          │   Proposals  │
│   (Go)       │◄────────►│   (Flask)    │◄────────►│   (YAML)     │
│              │          │              │          │              │
│  - Create    │          │  - Browse    │          │  - Storage   │
│  - View      │          │  - Create    │          │  - Backup    │
│  - Consensus │          │  - Status    │          │  - History   │
└──────────────┘          └──────────────┘          └──────────────┘
       │                         │                          │
       │                         │                          │
       └─────────────────────────┴──────────────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │   File System   │
                        │  data/proposals/│
                        └─────────────────┘
```

## Components

### 1. CLI Tool (Go)
**Location:** `./cmd/collectiveflow/`

**Purpose:** Command-line interface for proposal and consensus management

**Features:**
- Create proposals
- Start consensus processes
- Add agent input
- View status and history
- Complete consensus decisions

**Technology:**
- Language: Go 1.21+
- Framework: Cobra CLI
- Storage: File-based via interface
- Configuration: Viper

**No hierarchy:** CLI has same privileges as web interface. Both are equal ways to interact with proposals.

### 2. Web Interface (Flask)
**Location:** `./web/`

**Purpose:** Visual interface for browsing and creating proposals

**Features:**
- Browse all proposals by status
- View detailed proposal information
- Create new proposals via form
- See consensus history
- Collective activity dashboard
- RESTful API endpoints

**Technology:**
- Language: Python 3.13+
- Framework: Flask
- Templates: Jinja2
- Styling: Tailwind CSS (via CDN)
- API: JSON responses

**No authentication:** By design, no user system or roles. Everyone is equal.

### 3. Data Storage (YAML)
**Location:** `./data/proposals/`

**Purpose:** Human-readable, Git-friendly proposal storage

**Format:**
```yaml
id: proposal-2025-11-05-abc123
title: "Proposal Title"
description: "Detailed description"
proposer: agent-name
date: "2025-11-05T10:30:00"
status: consultation
urgency: medium
affected_areas:
  - area1
  - area2
consensus_status: "Awaiting input from 3 agents"
consensus_history:
  - timestamp: "2025-11-05T10:30:00"
    event: proposal_created
    actor: agent-name
    details: "Created with urgency: medium"
consultations:
  - agent: agent1
    position: support
    comment: "I agree with this approach"
    timestamp: "2025-11-05T11:00:00"
```

**Why YAML?**
- Human-readable (anyone can understand it)
- Git-friendly (easy to diff and review)
- No database complexity
- Direct file editing possible
- Language-agnostic (Go and Python both read it)

### 4. Optional: Docker Containers
**Location:** `./docker-compose.yml`, `./web/Dockerfile`

**Purpose:** Optional containerization for isolation and portability

**Services:**
- `web`: Flask application in container
- `filebrowser`: Optional tool for browsing proposal files

**Philosophy:** Docker is a convenience, not a requirement. Local development is equally valid.

## Data Flow

### Creating a Proposal (CLI)
```
User Command
     │
     ▼
CLI (collectiveflow)
     │
     ▼
Validation & ID Generation
     │
     ▼
YAML File Creation
     │
     ▼
data/proposals/proposal-2025-11-05-abc123.yaml
```

### Creating a Proposal (Web)
```
User Browser
     │
     ▼
Flask Form Handler
     │
     ▼
Validation & ID Generation
     │
     ▼
YAML File Creation
     │
     ▼
data/proposals/proposal-2025-11-05-abc123.yaml
```

### Viewing Proposals (Web)
```
User Browser
     │
     ▼
Flask Route Handler
     │
     ▼
Read YAML Files
     │
     ▼
Template Rendering
     │
     ▼
HTML Response
```

## Deployment Topologies

### Topology 1: Single User (Simplest)
```
┌─────────────────────┐
│   Your Laptop       │
│                     │
│  ┌──────────────┐   │
│  │ CLI + Web    │   │
│  │ (no Docker)  │   │
│  └──────────────┘   │
│         │           │
│         ▼           │
│  ┌──────────────┐   │
│  │ data/        │   │
│  │ proposals/   │   │
│  └──────────────┘   │
└─────────────────────┘
```

### Topology 2: Docker Local Development
```
┌─────────────────────────────┐
│   Your Laptop               │
│                             │
│  ┌──────────────────────┐   │
│  │  Docker Container    │   │
│  │  ┌────────────────┐  │   │
│  │  │ Flask Web      │  │   │
│  │  └────────────────┘  │   │
│  └──────────────────────┘   │
│         │ (mounted)         │
│         ▼                   │
│  ┌──────────────┐           │
│  │ data/        │           │
│  │ proposals/   │           │
│  └──────────────┘           │
└─────────────────────────────┘
```

### Topology 3: Shared Server
```
┌───────────────────────────────────┐
│   Server (any Linux box)          │
│                                   │
│  ┌────────────────┐               │
│  │ nginx/caddy    │               │
│  │ (HTTPS proxy)  │               │
│  └────────────────┘               │
│         │                         │
│         ▼                         │
│  ┌────────────────┐               │
│  │ Flask Web      │               │
│  │ (systemd/Docker)│              │
│  └────────────────┘               │
│         │                         │
│         ▼                         │
│  ┌────────────────┐               │
│  │ data/proposals/│               │
│  │ (Git backup)   │               │
│  └────────────────┘               │
└───────────────────────────────────┘
          │
          ▼ (access from anywhere)
┌─────────────────┐  ┌─────────────────┐
│  Collective     │  │  Collective     │
│  Member 1       │  │  Member 2       │
└─────────────────┘  └─────────────────┘
```

## Storage Interface

The system uses a **storage interface** pattern that allows future backends without changing core logic:

```go
// Current: File-based storage
type FileStorage struct {
    DataDir string
}

// Future: Database storage (if consensus agrees)
type DatabaseStorage struct {
    DB *sql.DB
}

// Both implement the same interface
type Storage interface {
    SaveProposal(p Proposal) error
    GetProposal(id string) (Proposal, error)
    ListProposals() ([]Proposal, error)
}
```

**Why this matters:**
- Maintains flexibility without complexity
- File storage is sufficient for most use cases
- Database backend possible without rewriting everything
- Collective can decide through consensus when/if to change

## Security Model

### Current: Trust-Based
- No authentication or authorization
- All users have equal access
- Trust based on collective participation
- Security through external controls (firewall, VPN, reverse proxy)

### Philosophy
- Authentication creates hierarchy (admin vs user)
- Transparency more important than access control
- Add external security layers as needed
- Focus on collective decision-making, not technical gatekeeping

### Production Considerations
If deploying for broader access:
- Use reverse proxy with HTTPS (nginx, Caddy)
- Implement network-level controls (firewall)
- Consider VPN for access control
- Add authentication at proxy level (not in app)

## Scalability

### Expected Capacity
| Metric | Laptop | Small Server | Notes |
|--------|--------|--------------|-------|
| Proposals | 10,000+ | 100,000+ | YAML files scale well |
| Users | 1 | 50-100 | Concurrent web users |
| Response Time | <10ms | <50ms | File reads are fast |
| Disk Space | <10MB | <100MB | ~1KB per proposal |

### Scaling Strategies (If Needed)
1. **Add database backend** - SQLite → PostgreSQL
2. **Add caching** - Redis for frequently accessed proposals
3. **Add CDN** - For static assets
4. **Horizontal scaling** - Multiple Flask workers

**Philosophy:** Start simple, scale only when needed, maintain transparency.

## Development Workflow

### Local Development Loop
```
1. Edit code (app.py, templates, etc.)
2. Flask auto-reloads (development mode)
3. Refresh browser
4. Repeat
```

### Testing Workflow
```
1. Create test proposal (CLI or web)
2. Verify YAML file created
3. Check web interface displays it
4. Test consensus flow
5. Examine YAML updates
```

### Deployment Workflow
```
1. Test locally (make dev-web)
2. Test with Docker (make docker-build && make docker-up)
3. Review changes (git diff)
4. Commit (git commit)
5. Deploy to server (git pull + restart service)
```

## Anti-Patterns Avoided

### ❌ Complex Infrastructure
- No Kubernetes (unnecessarily complex for this scale)
- No microservices (monolith is simpler)
- No message queues (direct file I/O is sufficient)
- No distributed systems (single machine works fine)

### ❌ Hidden Knowledge
- No magic configuration (everything explicit)
- No complex build systems (simple Makefile)
- No proprietary formats (YAML and JSON)
- No specialized tools (standard Docker/Python/Go)

### ❌ Hierarchical Patterns
- No admin users or roles
- No privileged access levels
- No approval workflows with authority
- No central notification system (by design)

## Future Enhancements (Through Consensus)

Potential improvements that maintain horizontal principles:

1. **Real-time Updates** - WebSocket for live proposal updates
2. **Federation** - Connect multiple CollectiveFlow instances
3. **Mobile Apps** - Native iOS/Android interfaces
4. **API Webhooks** - External integrations
5. **Enhanced Search** - Full-text search in proposals
6. **Export Formats** - PDF, Markdown export
7. **Visualization** - Graphs of consensus patterns

All enhancements require collective consensus and must maintain:
- Simplicity and transparency
- No hierarchy introduction
- Accessibility to all skill levels
- Local-first principles

## Maintenance

### Regular Tasks
- **Backup data directory** (git push or tar archive)
- **Update dependencies** (pip/go modules)
- **Review logs** (check for errors)
- **Monitor disk space** (proposals grow slowly)

### Infrequent Tasks
- **Upgrade Flask/Go** (when security updates)
- **Review and archive old proposals** (if desired)
- **Performance tuning** (only if issues arise)

## Conclusion

CollectiveFlow's architecture embodies its principles:

- **Simple** - Files, not databases
- **Transparent** - YAML anyone can read
- **Horizontal** - No technical hierarchy
- **Accessible** - Standard tools and patterns
- **Flexible** - Can grow without complexity

The architecture proves that horizontal coordination doesn't require complex infrastructure. Simple, transparent systems enable genuine collective ownership.

---

*Architecture designed by consensus, for consensus, through consensus.*
