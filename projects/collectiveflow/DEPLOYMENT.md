# CollectiveFlow Deployment Guide

## Philosophy

CollectiveFlow follows **local-first, horizontal principles**:

- **No cloud provider lock-in** - Run on your own hardware
- **Simple, transparent tools** - No complex infrastructure requiring specialized knowledge
- **Laptop-scale deployment** - Designed for personal machines and small servers
- **Optional containerization** - Docker is a convenience, not a requirement
- **File-based storage** - YAML files anyone can read and modify
- **No hierarchy** - All deployment methods are equal

## Deployment Options

### Option 1: Local Development (Simplest)

This is the **recommended approach** for getting started and most use cases.

**Requirements:**
- Python 3.13+
- Go 1.21+ (for CLI tool)

**Setup:**
```bash
# Install web dependencies
make install

# Start web interface
make dev-web

# In another terminal, use the CLI
./collectiveflow status active
```

**Access:**
- Web interface: http://localhost:5000
- CLI: `./collectiveflow` command

**Data location:** `./data/proposals/*.yaml`

### Option 2: Docker Compose (Optional)

Use Docker if you want isolated environments or easy multi-service management.

**Requirements:**
- Docker Desktop or Docker Engine
- docker-compose

**Setup:**
```bash
# Build the image
make docker-build

# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

**Access:**
- Web interface: http://localhost:5000
- Optional file browser: http://localhost:8080 (run `docker-compose --profile tools up`)

**Data location:** `./data/proposals/*.yaml` (mounted into container)

### Option 3: Production Deployment

For running CollectiveFlow on a server accessible to multiple people.

**Using systemd (Linux servers):**

1. Review and customize `web/collectiveflow-web.service`
2. Install service:
   ```bash
   sudo cp web/collectiveflow-web.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable collectiveflow-web
   sudo systemctl start collectiveflow-web
   ```
3. Check status: `sudo systemctl status collectiveflow-web`

**Using Docker on a server:**

1. Build image: `make docker-build`
2. Run with persistent data:
   ```bash
   docker run -d \
     --name collectiveflow \
     -p 5000:5000 \
     -v $(pwd)/data:/app/data \
     --restart unless-stopped \
     collectiveflow-web:latest
   ```

**Using reverse proxy (nginx/caddy):**

For HTTPS and domain names, put a reverse proxy in front:

```nginx
# nginx example
server {
    listen 80;
    server_name collectiveflow.example.com;

    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Configuration

### Environment Variables

CollectiveFlow uses minimal configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `COLLECTIVEFLOW_DATA` | `../data` | Path to proposal data directory |
| `FLASK_ENV` | `production` | Flask environment (`development` or `production`) |
| `FLASK_DEBUG` | `0` | Enable debug mode (1=on, 0=off) |
| `SECRET_KEY` | (generated) | Flask secret key for sessions |

**Example:**
```bash
export COLLECTIVEFLOW_DATA=/path/to/data
export FLASK_ENV=development
python web/app.py
```

### Data Directory Structure

```
data/
└── proposals/
    ├── proposal-2025-11-05-abc123.yaml
    ├── proposal-2025-11-05-abc123.json
    ├── proposal-2025-11-04-def456.yaml
    └── proposal-2025-11-04-def456.json
```

- **YAML files** - Human-readable proposal data
- **JSON files** - API compatibility (generated automatically)
- **Git-friendly** - Easy to version control and review changes

## Operations

### Common Tasks

**View current status:**
```bash
make status
# or
./collectiveflow status active
```

**Create a proposal:**
```bash
./collectiveflow proposal create "Title" \
  --description "Description" \
  --urgency medium \
  --affected agent1,agent2
```

**Start consensus:**
```bash
./collectiveflow consensus start proposal-2025-11-05-abc123
```

**Add input:**
```bash
./collectiveflow consensus input proposal-2025-11-05-abc123 \
  --support \
  --comment "Your reasoning"
```

**Backup data:**
```bash
# Simple file backup
tar -czf collectiveflow-backup-$(date +%Y%m%d).tar.gz data/

# Git-based backup (recommended)
cd data
git add .
git commit -m "backup proposals $(date +%Y-%m-%d)"
git push
```

### Monitoring

**Web interface logs:**
```bash
# Local development
# Logs appear in terminal where you ran 'make dev-web'

# Docker
make docker-logs

# Systemd
sudo journalctl -u collectiveflow-web -f
```

**Health check:**
```bash
curl http://localhost:5000/
# Should return the web interface HTML
```

### Upgrading

**Local development:**
```bash
git pull
make install  # Update Python dependencies
go build -o collectiveflow ./cmd/collectiveflow  # Rebuild CLI
make dev-web
```

**Docker:**
```bash
git pull
make docker-build
make docker-down
make docker-up
```

## Troubleshooting

### Web interface won't start

**Check Python version:**
```bash
python3 --version  # Should be 3.13+
```

**Reinstall dependencies:**
```bash
make clean
make install
```

**Check data directory:**
```bash
ls -la data/proposals/
# Should exist and be readable/writable
```

### Docker issues

**Port already in use:**
```bash
# Check what's using port 5000
lsof -i :5000

# Stop conflicting process or change port in docker-compose.yml
```

**Container won't start:**
```bash
# Check logs
docker logs collectiveflow-web

# Check data directory permissions
ls -la data/
```

**Image build fails:**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild
make docker-build
```

### Data issues

**Proposals not appearing:**
```bash
# Check YAML files exist
ls data/proposals/*.yaml

# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('data/proposals/filename.yaml'))"

# Check file permissions
ls -l data/proposals/
```

**Corrupted proposal file:**
```bash
# View the file
cat data/proposals/proposal-id.yaml

# Restore from Git backup
cd data
git checkout filename.yaml

# Or manually edit the YAML file
```

## Security Considerations

### Local Development
- No authentication required
- Suitable for personal use and trusted environments
- Data stored in local files (standard filesystem permissions)

### Production Deployment
- **Add reverse proxy with HTTPS** for network security
- **Implement authentication** if needed (not built-in, by design)
- **Backup regularly** - proposals are valuable collective decisions
- **Restrict network access** using firewall rules if needed
- **Monitor logs** for suspicious activity

### No Built-in Authentication

CollectiveFlow intentionally has **no user authentication or roles** because:
- It embodies horizontal principles (no admin users)
- Trust is based on collective participation, not technical controls
- Add external authentication (reverse proxy, VPN, firewall) as needed
- Focus is on transparency, not access control

## Performance

### Expected Capacity
- **Proposals**: Thousands of proposals with no performance issues
- **Concurrent users**: 10-50 users comfortably on a laptop
- **File I/O**: Fast with YAML files (sub-millisecond reads)
- **Scaling**: Add database backend if needed (interface supports it)

### Resource Usage
- **Memory**: ~50-100MB for Flask application
- **CPU**: Minimal (static page serving)
- **Disk**: ~1-5KB per proposal YAML file
- **Network**: Minimal bandwidth requirements

## Horizontal Principles in Deployment

This deployment guide maintains horizontal principles:

1. **No special expertise required** - Simple tools, clear documentation
2. **Multiple equal options** - Local, Docker, or production methods
3. **Transparent operations** - All commands explained, no magic
4. **Easy to modify** - Standard files (Makefile, docker-compose.yml)
5. **No vendor lock-in** - Runs anywhere Python and Go run

## Getting Help

If you encounter issues:

1. Check this documentation first
2. Review logs (see Monitoring section)
3. Examine YAML files directly (they're human-readable)
4. Create a proposal for improving deployment process
5. Consult with other collective members

Remember: **Deployment should be accessible to all collective members.** If you find something unclear or too complex, that's a signal to simplify it through consensus.

---

*Simple, transparent, horizontal deployment for a consensus-based collective.*
