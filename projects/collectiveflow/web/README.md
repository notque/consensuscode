# CollectiveFlow Web Interface

A horizontal, non-hierarchical web interface for the CollectiveFlow consensus system.

## Principles

This web interface embodies the collective's core principles:
- **No authentication or login** - All information is equally accessible
- **No admin panels or special roles** - Every user has the same access
- **Transparent decision history** - All proposals and consultations are visible
- **Horizontal design** - No hierarchical navigation or special privileges

## Features

- View all proposals organized by status
- See detailed proposal information including consultations
- Track the collective's activity and statistics
- Access the same data as the CLI, but through a web browser
- Responsive design works on all devices

## Installation

1. Install Python dependencies:
```bash
cd web
pip install -r requirements.txt
```

2. Run the Flask development server:
```bash
python app.py
```

3. Open http://localhost:5000 in your browser

## Configuration

The web interface reads proposals from the same data directory as the CLI. By default, it looks for `../data/proposals/`. You can override this with the `COLLECTIVEFLOW_DATA` environment variable:

```bash
export COLLECTIVEFLOW_DATA=/path/to/collectiveflow/data
python app.py
```

## API Endpoints

The web interface also provides JSON API endpoints:
- `GET /api/proposals` - List all proposals
- `GET /api/proposal/<id>` - Get specific proposal details

## Development

This is a simple Flask application designed for accessibility and transparency. It uses:
- Flask for the web framework
- Tailwind CSS (via CDN) for styling
- No JavaScript required for core functionality
- No database - reads directly from YAML files

## Contributing

To propose changes to the web interface, use the CollectiveFlow CLI to create a proposal. All changes should maintain our horizontal principles - no features that create hierarchy or special privileges.