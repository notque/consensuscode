# CollectiveFlow Deployment Comparison

## Choose Your Deployment Method

This guide helps you choose the right deployment approach for your collective's needs.

## Quick Decision Tree

```
Do you need to share with multiple people?
│
├─ NO → Use Local Development
│        (Simplest, fastest, most transparent)
│
└─ YES → Do you have a server?
         │
         ├─ NO → Use Docker on someone's laptop
         │        (Portable, easy to move)
         │
         └─ YES → Do you want containers?
                  │
                  ├─ YES → Docker on server
                  │        (Isolated, reproducible)
                  │
                  └─ NO → Systemd on server
                           (Traditional, simple)
```

## Deployment Methods Compared

| Factor | Local Dev | Docker Local | Server (systemd) | Server (Docker) |
|--------|-----------|--------------|------------------|-----------------|
| **Setup Time** | 2 minutes | 3 minutes | 10 minutes | 5 minutes |
| **Complexity** | Lowest | Low | Medium | Medium |
| **Best For** | Solo development | Team on same network | Shared server access | Containerized infrastructure |
| **Requirements** | Python, Go | Docker Desktop | Linux server, systemd | Linux server, Docker |
| **Isolation** | None (direct) | Container | Process | Container |
| **Networking** | Localhost only | Local network | Internet | Internet |
| **Persistence** | Local files | Mounted volumes | Server filesystem | Mounted volumes |
| **Backup** | Git or tar | Git or tar | Git or server backup | Git or volume backup |
| **Updates** | git pull + restart | Rebuild image | git pull + systemctl restart | Rebuild + restart |
| **Resource Use** | ~50MB RAM | ~100MB RAM | ~50MB RAM | ~100MB RAM |
| **Expertise Level** | Beginner | Beginner | Intermediate | Intermediate |
| **Horizontal?** | ✅ Very | ✅ Yes | ✅ Yes | ✅ Yes |

## Detailed Comparison

### Local Development (Recommended Starting Point)

**What it is:**
- Python virtual environment
- Flask development server
- Files stored locally
- No containerization

**When to use:**
- Personal exploration
- Single-user development
- Testing and experimentation
- Learning the system

**Advantages:**
- Fastest setup
- Most transparent (no container abstraction)
- Easy to debug
- Direct file access
- No Docker knowledge needed

**Disadvantages:**
- Not shareable beyond localhost
- Python dependencies on host system
- Development server (not production-grade)

**Setup:**
```bash
make install
make dev-web
```

**Horizontal score: 10/10**
- Most accessible to everyone
- Clearest to understand
- No hidden layers

---

### Docker Compose (Local Network)

**What it is:**
- Containerized Flask app
- Gunicorn production server
- Optional file browser tool
- Shared on local network

**When to use:**
- Small team on same network
- Want isolation from host system
- Testing production-like setup
- Learning Docker

**Advantages:**
- Production-ready server (Gunicorn)
- Isolated from host
- Easy to share on LAN
- Reproducible environment
- Optional development tools

**Disadvantages:**
- Requires Docker installation
- Slight abstraction layer
- ~100MB image size
- Build time required

**Setup:**
```bash
make docker-build
make docker-up
# Access at http://localhost:5000 or http://<your-ip>:5000
```

**Horizontal score: 9/10**
- Standard Docker (widely known)
- Optional, not required
- Clear docker-compose.yml
- Easy to modify

---

### Production: Systemd (Server)

**What it is:**
- Flask app as systemd service
- Reverse proxy (nginx/Caddy)
- Traditional Linux server setup
- Domain name + HTTPS

**When to use:**
- Permanent shared server
- Want traditional server management
- Team comfortable with Linux
- Prefer no containers

**Advantages:**
- Standard Linux service management
- Direct filesystem access
- Familiar to sysadmins
- No container overhead
- systemctl integration

**Disadvantages:**
- Requires Linux server
- Manual dependency management
- Server-specific configuration
- Needs reverse proxy setup

**Setup:**
```bash
# On server
make install
sudo cp web/collectiveflow-web.service /etc/systemd/system/
sudo systemctl enable collectiveflow-web
sudo systemctl start collectiveflow-web
# Configure nginx/Caddy
```

**Horizontal score: 8/10**
- Standard Linux tools
- Transparent service file
- Some server admin knowledge needed
- Well-documented patterns

---

### Production: Docker (Server)

**What it is:**
- Docker container on server
- Reverse proxy in front
- Volume-mounted data
- Optional Docker Compose

**When to use:**
- Permanent shared server
- Want container benefits
- Easy deployment and updates
- Server already uses Docker

**Advantages:**
- Isolated environment
- Easy updates (rebuild image)
- Portable between servers
- Reproducible setup
- Built-in health checks

**Disadvantages:**
- Requires Docker on server
- Container abstraction layer
- Image storage space
- Docker networking setup

**Setup:**
```bash
# Build locally or on server
make docker-build

# Run on server
docker run -d \
  --name collectiveflow \
  -p 5000:5000 \
  -v $(pwd)/data:/app/data \
  --restart unless-stopped \
  collectiveflow-web:latest

# Configure nginx/Caddy reverse proxy
```

**Horizontal score: 8/10**
- Standard Docker patterns
- Clear container config
- Some Docker knowledge needed
- Well-documented

---

## Special Considerations

### For Collectives New to Deployment

**Start with Local Development**
1. Everyone installs locally
2. Share proposals via Git
3. Each person runs their own instance
4. Collaborate through file syncing

**Advantages:**
- No single point of failure
- Everyone has full data access
- Most transparent approach
- Distributed by design

### For Collectives with Server Access

**Use Docker on Server**
1. One shared instance
2. Everyone accesses via browser
3. Single source of truth
4. Easier coordination

**Advantages:**
- Centralized proposals
- Real-time updates
- No local installation needed
- Simpler backup strategy

### For Collectives Without Server

**Use Docker on Rotating Laptops**
1. Install Docker Compose on member laptops
2. Rotate who hosts each week/month
3. Data stays in Git repository
4. Truly horizontal hosting

**Advantages:**
- No server costs
- Distributed ownership
- Rotation prevents concentration
- Anyone can host

## Cost Comparison

| Method | Hardware Cost | Software Cost | Monthly Cost | Setup Cost |
|--------|--------------|---------------|--------------|------------|
| **Local Dev** | $0 (your laptop) | $0 (all free) | $0 | 30 min |
| **Docker Local** | $0 (your laptop) | $0 (Docker free) | $0 | 30 min |
| **VPS Server** | $5-10/mo | $0 (all FOSS) | $5-10 | 2 hours |
| **Home Server** | $200 (Raspberry Pi) | $0 (all FOSS) | $0 | 3 hours |
| **Rotating Hosts** | $0 (existing laptops) | $0 (Docker free) | $0 | 1 hour/person |

**Recommendation:** Start with $0 option, upgrade only if needed.

## Horizontal Analysis

### Most Horizontal: Local Development
- **No technical hierarchy**: Everyone uses same simple setup
- **Full transparency**: Direct file access
- **Equal capability**: Everyone has complete system
- **No dependencies**: Works offline
- **Knowledge sharing**: Easy to explain to newcomers

### Also Horizontal: Rotating Docker Hosts
- **Distributed power**: No single host
- **Rotation prevents hierarchy**: Regular change
- **Collective responsibility**: Everyone hosts eventually
- **Low barrier**: Docker Compose is simple
- **Resilience**: System continues if someone leaves

### Least Horizontal (But Still Good): Permanent Server
- **Risk**: Server admin becomes key person
- **Mitigation**: Document everything, rotate access
- **Benefit**: Single source of truth
- **Consideration**: Backup and recovery procedures

## Security Comparison

| Method | Network Exposure | Authentication | Encryption | Attack Surface |
|--------|------------------|----------------|------------|----------------|
| **Local Dev** | None (localhost) | None (trust) | No (local) | Minimal |
| **Docker Local** | LAN only | None (trust) | Optional | Small |
| **Server HTTP** | Internet | Proxy-level | No | Medium |
| **Server HTTPS** | Internet | Proxy-level | Yes (TLS) | Medium |

**Recommendation:**
- Local/LAN: No auth needed (trust-based)
- Internet: Add HTTPS + optional auth at proxy level

## Maintenance Comparison

| Method | Update Frequency | Update Complexity | Backup Method | Recovery Time |
|--------|------------------|-------------------|---------------|---------------|
| **Local Dev** | On demand | `git pull` | Git or tar | Minutes |
| **Docker Local** | On demand | Rebuild image | Git or tar | Minutes |
| **Server (systemd)** | Weekly/monthly | `git pull + restart` | Git + server backup | 10 minutes |
| **Server (Docker)** | Weekly/monthly | Rebuild + restart | Git + volume backup | 15 minutes |

**All methods maintain horizontal principles:**
- Clear update procedures
- No specialized knowledge
- Documented steps
- Collective can perform

## Environmental Impact

| Method | Power Usage | Resource Usage | Carbon Footprint |
|--------|-------------|----------------|------------------|
| **Local Dev** | ~20W (laptop) | Minimal | Existing hardware |
| **Docker Local** | ~25W (laptop) | +100MB RAM | Existing hardware |
| **VPS Server** | Share of datacenter | 512MB-1GB | Depends on host |
| **Home Server** | ~5W (Raspberry Pi) | Dedicated device | Manufacturing + power |

**Most sustainable:** Local dev on existing hardware

## Recommendation by Collective Size

| Collective Size | Recommended Method | Rationale |
|----------------|-------------------|-----------|
| **1 person** | Local Development | Simplest, most transparent |
| **2-5 people** | Docker on rotating laptops | Distributed, low cost |
| **5-10 people** | Docker on shared VPS | Balance of simplicity and access |
| **10+ people** | Docker on server + backups | Reliability for many users |

## Migration Path

Start simple, grow as needed:

```
1. Local Development
   ↓ (collective grows)
2. Docker on rotating laptops
   ↓ (want stability)
3. Docker on VPS server
   ↓ (high traffic)
4. Scale horizontally (multiple instances)
```

**Key principle:** Only add complexity when collective actually needs it.

## Decision Checklist

Use this checklist to choose:

- [ ] How many people need access?
- [ ] Do we have a server?
- [ ] What's our technical comfort level?
- [ ] Do we need 24/7 availability?
- [ ] What's our budget?
- [ ] How important is simplicity?
- [ ] Can we maintain a server?
- [ ] Do we trust each other (for auth)?
- [ ] How will we backup data?
- [ ] Who can help if things break?

**After answering, refer to comparison tables above.**

## Horizontal Best Practices

Regardless of deployment method:

1. **Document everything** - Make deployment accessible
2. **Rotate responsibilities** - Prevent knowledge concentration
3. **Regular backups** - Git-based for transparency
4. **Clear procedures** - Anyone can deploy/maintain
5. **Collective decisions** - Changes through consensus
6. **Test locally first** - Validate before production
7. **Monitor simply** - Logs, health checks, no complexity
8. **Keep it boring** - Standard tools, proven patterns

## Conclusion

**For most collectives:**
- Start with **Local Development**
- Upgrade to **Docker on rotating laptops** if sharing needed
- Move to **Docker on server** only when scale requires it

**Remember:** The best deployment is the one your collective understands and can maintain collectively.

---

*Choose simplicity. Scale only when needed. Maintain horizontal principles.*
