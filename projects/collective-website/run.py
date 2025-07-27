#!/usr/bin/env python3
"""
Development server runner for the Collective Voice website
"""

import os
from app import app

if __name__ == '__main__':
    # Set development environment
    os.environ['FLASK_ENV'] = 'development'
    app.run(debug=True, host='127.0.0.1', port=5000)