#!/bin/bash

# Collective Voice Website Deployment Script

set -e

echo "Deploying Collective Voice Website..."

# Set production environment
export FLASK_ENV=production

# Install dependencies
echo "Installing dependencies..."
pip install -r requirements.txt

# Run tests
echo "Running tests..."
python -m pytest test_app.py -v

# Create necessary directories
mkdir -p logs

# Set secure permissions
chmod 600 config.py

# Start application with gunicorn
echo "Starting application server..."
gunicorn -w 4 -b 0.0.0.0:8000 app:app \
    --access-logfile logs/access.log \
    --error-logfile logs/error.log \
    --log-level info \
    --daemon

echo "Collective Voice website deployed successfully!"
echo "Access at: http://localhost:8000"