#!/bin/bash

# Development Setup for Collective Voice Website

echo "Setting up Collective Voice Website for development..."

# Check if Python 3 is available
if ! command -v python3 &> /dev/null; then
    echo "Python 3 is required but not installed."
    exit 1
fi

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "Installing dependencies..."
pip install -r requirements.txt

# Run tests to verify setup
echo "Running tests..."
python test_app.py

echo ""
echo "Setup complete! To start development:"
echo "1. source venv/bin/activate"
echo "2. python run.py"
echo ""
echo "The website will be available at http://127.0.0.1:5000"