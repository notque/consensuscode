# Deployment Guide

This guide helps you deploy CollectiveFlow in your collective's infrastructure. It covers both the CLI tool and web interface, with emphasis on approaches that align with horizontal principles.

## Deployment Philosophy

CollectiveFlow's deployment approach reflects our values:

- **Local-first**: Works on personal machines without cloud dependencies
- **Simple infrastructure**: Avoids complex orchestration that creates knowledge hierarchies
- **Transparent operation**: All data visible, no hidden state
- **Cost-effective**: Free to run, no cloud provider lock-in
- **Collective ownership**: Anyone in the collective can manage it

## Quick Start: Single-User Deployment

For individual use or initial testing:

### CLI Tool

```bash
# Build the binary
go build -o collectiveflow ./cmd/collectiveflow

# Move to your PATH (optional but convenient)
sudo mv collectiveflow /usr/local/bin/

# Or keep it local
mkdir -p ~/.local/bin
mv collectiveflow ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"  # Add to .bashrc or .zshrc

# Test it works
collectiveflow --help
```

**Data location**: By default, proposals are stored in `./data/proposals/` relative to where you run the command.

**Configure data location**:
```bash
# Set persistent location
export COLLECTIVEFLOW_DATA="$HOME/.collectiveflow/data"
mkdir -p "$COLLECTIVEFLOW_DATA/proposals"

# Add to shell config to make permanent
echo 'export COLLECTIVEFLOW_DATA="$HOME/.collectiveflow/data"' >> ~/.bashrc
```

## Shared Collective Deployment

For a collective sharing proposals across multiple people:

### Option 1: Shared Directory (Local Network)

**Best for**: Small collectives with shared file server or NFS mount

```bash
# On shared storage
mkdir -p /shared/collective/collectiveflow/data/proposals

# Each person configures their environment
export COLLECTIVEFLOW_DATA="/shared/collective/collectiveflow/data"

# Everyone uses the same data directory
collectiveflow proposal list  # Shows collective's proposals
```

**Pros**:
- Simple setup
- No server to maintain
- Direct file access (transparency)

**Cons**:
- Requires shared filesystem
- Concurrent writes need coordination
- No remote access without VPN

**Concurrency note**: YAML files handle concurrent reads fine. Concurrent writes are rare (proposals are created infrequently). For high-concurrency collectives, see database deployment.

### Option 2: Git-Based Sync

**Best for**: Distributed collectives without shared filesystem

```bash
# Initialize git repo for proposals
cd data
git init
git add proposals/
git commit -m "initial proposals"

# Push to shared repo (GitHub, GitLab, self-hosted)
git remote add origin git@github.com:your-collective/proposals.git
git push -u origin main

# Each person clones
git clone git@github.com:your-collective/proposals.git ~/.collectiveflow/data

# Workflow
collectiveflow proposal create "My proposal"
cd ~/.collectiveflow/data
git add .
git commit -m "add proposal"
git push

# Others pull to get updates
git pull
```

**Pros**:
- Works for remote collectives
- Git provides version history automatically
- No central server needed (can use any git host)
- Merge conflicts are rare (proposals have unique IDs)

**Cons**:
- Manual sync step (git pull/push)
- Requires git knowledge
- Not real-time

**Automation tip**: Use cron or systemd timer to auto-sync:

```bash
# Crontab entry to sync every 5 minutes
*/5 * * * * cd ~/.collectiveflow/data && git pull --rebase && git push
```

### Option 3: Centralized File Server

**Best for**: Collectives wanting centralized storage with rsync/scp access

```bash
# Server setup (one-time)
mkdir -p /var/collectiveflow/data/proposals
chown collectiveflow:collectiveflow /var/collectiveflow
chmod 770 /var/collectiveflow

# Client setup
export COLLECTIVEFLOW_DATA="/mnt/collectiveflow"

# Mount remote directory (SSHFS example)
mkdir -p /mnt/collectiveflow
sshfs collective@server.org:/var/collectiveflow /mnt/collectiveflow
```

**Pros**:
- Central location
- Real-time for all users
- Simple backup (one location)

**Cons**:
- Requires server
- Network dependency
- Single point of failure

## Web Interface Deployment

### Development/Testing

For trying out the web interface:

```bash
cd web

# Set up Python environment
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Run Flask development server
export COLLECTIVEFLOW_DATA="../data"
python app.py

# Access at http://localhost:5000
```

**Do not use Flask's development server in production** - it's not secure or performant.

### Production Deployment

For production web interface:

#### Option 1: Gunicorn + Nginx (Recommended)

**Best for**: Serious deployments, multiple users

**Setup**:

1. **Install dependencies**:
```bash
# On Ubuntu/Debian
sudo apt install nginx python3-venv python3-pip

cd /opt/collectiveflow/web
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pip install gunicorn  # Production WSGI server
```

2. **Create systemd service** (`/etc/systemd/system/collectiveflow-web.service`):
```ini
[Unit]
Description=CollectiveFlow Web Interface
After=network.target

[Service]
Type=notify
User=collectiveflow
Group=collectiveflow
WorkingDirectory=/opt/collectiveflow/web
Environment="COLLECTIVEFLOW_DATA=/var/collectiveflow/data"
Environment="PATH=/opt/collectiveflow/web/venv/bin"
ExecStart=/opt/collectiveflow/web/venv/bin/gunicorn \
    --bind 127.0.0.1:8000 \
    --workers 4 \
    --access-logfile /var/log/collectiveflow/access.log \
    --error-logfile /var/log/collectiveflow/error.log \
    app:app

[Install]
WantedBy=multi-user.target
```

3. **Configure Nginx** (`/etc/nginx/sites-available/collectiveflow`):
```nginx
server {
    listen 80;
    server_name collective.example.org;

    # Redirect to HTTPS (after setting up SSL)
    # return 301 https://$server_name$request_uri;

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /static {
        alias /opt/collectiveflow/web/static;
        expires 30d;
    }
}
```

4. **Enable and start**:
```bash
# Enable Nginx site
sudo ln -s /etc/nginx/sites-available/collectiveflow /etc/nginx/sites-enabled/
sudo nginx -t  # Test configuration
sudo systemctl reload nginx

# Start CollectiveFlow web service
sudo systemctl enable collectiveflow-web
sudo systemctl start collectiveflow-web
sudo systemctl status collectiveflow-web
```

5. **Set up SSL with Let's Encrypt** (recommended):
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d collective.example.org
```

**Pros**:
- Production-ready
- Good performance
- Automatic HTTPS
- Professional setup

**Cons**:
- More complex than development server
- Requires root access to set up
- Need to understand Nginx/systemd

#### Option 2: Docker Compose (Simplified)

**Best for**: Collectives familiar with Docker

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  web:
    build:
      context: .
      dockerfile: web/Dockerfile
    ports:
      - "5000:5000"
    volumes:
      - ./data:/app/data:ro  # Read-only data mount
    environment:
      - COLLECTIVEFLOW_DATA=/app/data
      - FLASK_ENV=production
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - web
    restart: unless-stopped
```

Create `web/Dockerfile`:

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY web/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt gunicorn

COPY web/ .

EXPOSE 5000

CMD ["gunicorn", "--bind", "0.0.0.0:5000", "--workers", "4", "app:app"]
```

**Deploy**:
```bash
docker-compose up -d
docker-compose logs -f  # View logs
```

**Pros**:
- Isolated environment
- Easy to update (rebuild container)
- Consistent across machines

**Cons**:
- Requires Docker knowledge
- More moving parts
- Container adds complexity

#### Option 3: Simple Systemd Service (Minimal)

**Best for**: Very small collectives, local access only

Just run Gunicorn behind systemd, no Nginx:

```bash
# Install
cd /opt/collectiveflow/web
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt gunicorn

# Systemd service (similar to above, but bind to 0.0.0.0:5000)
# Access directly at http://server:5000

# For HTTPS, set up SSL certificate and use Gunicorn SSL options
gunicorn --certfile=cert.pem --keyfile=key.pem --bind 0.0.0.0:5443 app:app
```

**Pros**:
- Simpler than Nginx setup
- Fewer components

**Cons**:
- No static file optimization
- No HTTPS termination (must use Gunicorn SSL)
- Less production-ready

## Database Backend Deployment (Future)

CollectiveFlow currently uses file-based storage. When database backend is implemented:

### PostgreSQL Option

```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Create database and user
sudo -u postgres psql
CREATE DATABASE collectiveflow;
CREATE USER collectiveflow WITH PASSWORD 'secure-password';
GRANT ALL PRIVILEGES ON DATABASE collectiveflow TO collectiveflow;

# Configure CollectiveFlow
export COLLECTIVEFLOW_STORAGE="postgres"
export COLLECTIVEFLOW_DB_URL="postgresql://collectiveflow:secure-password@localhost/collectiveflow"

# Run CLI (will use database)
collectiveflow proposal list
```

**When to use database**:
- Hundreds of active proposals
- Complex querying needed
- Performance becomes an issue

**When not to use database**:
- Small collectives (files are simpler)
- Transparency is more important than performance
- No one in collective comfortable with database admin

## Backup Strategy

**Critical**: Back up proposal data regularly. Losing consensus history damages collective trust.

### File-Based Backup

```bash
# Simple: Daily backup to another location
#!/bin/bash
# /etc/cron.daily/collectiveflow-backup

BACKUP_DIR="/backup/collectiveflow"
DATA_DIR="/var/collectiveflow/data"
DATE=$(date +%Y%m%d)

mkdir -p "$BACKUP_DIR"
tar -czf "$BACKUP_DIR/collectiveflow-$DATE.tar.gz" -C "$DATA_DIR" .

# Keep last 30 days
find "$BACKUP_DIR" -name "collectiveflow-*.tar.gz" -mtime +30 -delete
```

### Git-Based Backup

If using git sync (Option 2 above), git history is your backup:

```bash
# Verify you have history
git log

# Restore to specific point
git checkout <commit-hash>

# Or revert specific change
git revert <commit-hash>
```

### Remote Backup

```bash
# Rsync to remote server
rsync -avz /var/collectiveflow/data/ backup-server:/backups/collectiveflow/

# Or use rclone for cloud storage (if needed)
rclone sync /var/collectiveflow/data/ remote:collectiveflow-backup/
```

**Automation**: Add to cron for automatic backup:

```bash
# Crontab: Daily at 2 AM
0 2 * * * /usr/local/bin/collectiveflow-backup.sh
```

## Monitoring

### Log Files

**CLI**: Logs to stdout/stderr by default

**Web interface**: Configure logging in systemd service

```bash
# View web logs
sudo journalctl -u collectiveflow-web -f

# Or if using file logging
tail -f /var/log/collectiveflow/error.log
```

### Health Checks

Simple health check for web interface:

```bash
# Check if web interface is responding
curl http://localhost:5000/

# More sophisticated check
#!/bin/bash
if curl -f http://localhost:5000/ > /dev/null 2>&1; then
    echo "CollectiveFlow is healthy"
else
    echo "CollectiveFlow is down!"
    # Send alert (email, SMS, etc.)
fi
```

Add to cron for monitoring:
```bash
*/5 * * * * /usr/local/bin/collectiveflow-healthcheck.sh
```

### Metrics (Optional)

If your collective wants metrics:

```bash
# Simple: Log access patterns
tail -f /var/log/nginx/access.log | grep collectiveflow

# Advanced: Use prometheus + grafana (requires more setup)
# Only add if collective decides monitoring benefits outweigh complexity
```

## Security Considerations

### No Built-In Authentication

CollectiveFlow has no authentication by design (horizontal principles). For security:

**Option 1: Network isolation**
- Deploy on VPN-only network
- Use firewall rules to limit access
- SSH tunneling for remote access

**Option 2: Reverse proxy authentication**
- Use Nginx basic auth or OAuth proxy
- Authentication is external to CollectiveFlow
- Everyone who authenticates has full access (no roles)

Example Nginx basic auth:

```nginx
server {
    location / {
        auth_basic "Collective Access";
        auth_basic_user_file /etc/nginx/.htpasswd;

        proxy_pass http://127.0.0.1:8000;
        # ... other proxy settings
    }
}
```

Create password file:
```bash
sudo htpasswd -c /etc/nginx/.htpasswd collective-member
```

**Important**: Use this for access control, not to create hierarchy. Everyone in the collective should have access.

### HTTPS

Always use HTTPS for production deployments:

```bash
# Let's Encrypt (free, automated)
sudo certbot --nginx -d collective.example.org

# Certbot auto-renews, but verify
sudo certbot renew --dry-run
```

### File Permissions

Protect data directory:

```bash
# Only collectiveflow user can write
sudo chown -R collectiveflow:collectiveflow /var/collectiveflow
sudo chmod 750 /var/collectiveflow
sudo chmod 640 /var/collectiveflow/data/proposals/*.yaml
```

### Updates

Keep system and dependencies updated:

```bash
# System packages
sudo apt update && sudo apt upgrade

# Go updates (rebuild CLI)
go version  # Check for updates at golang.org
go build -o collectiveflow ./cmd/collectiveflow

# Python dependencies (web interface)
cd web
source venv/bin/activate
pip install --upgrade -r requirements.txt
```

## Troubleshooting

### CLI Issues

**Problem**: Command not found
```bash
# Solution: Check PATH
which collectiveflow
echo $PATH

# Add to PATH if needed
export PATH="/path/to/collectiveflow:$PATH"
```

**Problem**: Can't find data directory
```bash
# Solution: Set COLLECTIVEFLOW_DATA
export COLLECTIVEFLOW_DATA="/correct/path/to/data"

# Verify
collectiveflow proposal list
```

**Problem**: Permission denied reading/writing proposals
```bash
# Solution: Fix permissions
sudo chown -R $USER /path/to/collectiveflow/data
chmod u+rw /path/to/collectiveflow/data/proposals/*.yaml
```

### Web Interface Issues

**Problem**: Web interface won't start
```bash
# Check logs
sudo journalctl -u collectiveflow-web -n 50

# Common causes:
# - Port 8000 already in use (change port)
# - Python environment issues (recreate venv)
# - COLLECTIVEFLOW_DATA not set (check service file)
```

**Problem**: Nginx 502 Bad Gateway
```bash
# Gunicorn isn't running
sudo systemctl status collectiveflow-web
sudo systemctl restart collectiveflow-web

# Check Gunicorn logs
sudo journalctl -u collectiveflow-web -f
```

**Problem**: Can't see proposals in web interface
```bash
# Check data directory
echo $COLLECTIVEFLOW_DATA
ls -la $COLLECTIVEFLOW_DATA/proposals/

# Check service has correct COLLECTIVEFLOW_DATA
sudo systemctl cat collectiveflow-web | grep COLLECTIVEFLOW_DATA
```

### Performance Issues

**Problem**: CLI slow with many proposals
```bash
# Current limitation of file-based storage
# Solutions:
# 1. Clean up old implemented proposals (archive)
# 2. Wait for database backend implementation
# 3. Optimize queries (if you can contribute code)
```

**Problem**: Web interface slow
```bash
# Add more Gunicorn workers
# Edit systemd service: --workers 8  # (increase number)
sudo systemctl daemon-reload
sudo systemctl restart collectiveflow-web

# Enable caching in Nginx
# Add to nginx config:
location /static {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

## Scaling Considerations

### When to Scale

For most collectives, single-server deployment is sufficient. Consider scaling when:

- More than 50 concurrent users
- More than 1000 active proposals
- Response times consistently > 1 second
- Geographic distribution requires edge servers

### Horizontal Scaling

If you need to scale:

1. **Database backend**: Migrate from files to PostgreSQL
2. **Load balancer**: Multiple web servers behind HAProxy/Nginx
3. **Read replicas**: For query-heavy workloads

**Warning**: Scaling adds complexity. Only scale when genuinely needed. Complexity can create knowledge hierarchies.

## Migration Between Deployments

### From Local to Shared

```bash
# Copy your local proposals to shared location
cp -r ~/collectiveflow/data/proposals/* /shared/collectiveflow/data/proposals/

# Update configuration
export COLLECTIVEFLOW_DATA="/shared/collectiveflow/data"
```

### From File to Database (Future)

When database backend is available:

```bash
# Export from files
collectiveflow export --format yaml --output proposals-export.yaml

# Import to database
export COLLECTIVEFLOW_STORAGE="postgres"
collectiveflow import proposals-export.yaml
```

## Deployment Checklist

Before deploying to production:

- [ ] Data directory configured and accessible
- [ ] Backup strategy implemented and tested
- [ ] Monitoring/health checks in place
- [ ] HTTPS configured (if web interface)
- [ ] Access control appropriate for collective
- [ ] Documentation shared with collective
- [ ] Tested restore from backup
- [ ] Collective consensus on deployment approach

## Getting Help

If deployment issues arise:

1. Check logs (systemd journal, nginx logs, application logs)
2. Verify configuration (environment variables, file permissions)
3. Test components individually (CLI works? Web works? Nginx works?)
4. Create proposal asking for deployment help from collective
5. Include error messages and what you've tried

Remember: Deployment complexity should serve the collective, not create barriers to participation. If this guide makes deployment seem too complex, that's feedback for the collective to simplify the architecture.

---

Your deployment serves collective coordination. Keep it as simple as possible while meeting your collective's needs.
