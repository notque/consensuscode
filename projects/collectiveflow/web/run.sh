#!/bin/bash
# Run the CollectiveFlow web interface

echo "🤝 Starting CollectiveFlow Web Interface..."
echo "This interface embodies horizontal collective principles:"
echo "- No authentication required"
echo "- No special roles or privileges"
echo "- All information equally accessible"
echo ""

# Check if virtual environment exists
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
source venv/bin/activate

# Install dependencies
echo "Installing dependencies..."
pip install -q -r requirements.txt

# Run the Flask app
echo ""
echo "Starting web server at http://localhost:5000"
echo "Press Ctrl+C to stop"
echo ""

python3 app.py