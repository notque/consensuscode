"""
Configuration settings for the Collective Voice website
"""

import os
from pathlib import Path

class Config:
    """Base configuration class"""
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'collective-development-key-change-in-production'
    
    # Collective data paths
    COLLECTIVE_ROOT = Path(__file__).parent.parent.parent / 'projects' / 'collectiveflow'
    
    # Flask settings
    DEBUG = False
    TESTING = False
    
    # Cache settings
    CACHE_TIMEOUT = 60  # seconds
    
    # API settings
    API_VERSION = '1.0'
    MAX_VOICES_DISPLAY = 10
    
class DevelopmentConfig(Config):
    """Development configuration"""
    DEBUG = True
    ENV = 'development'

class ProductionConfig(Config):
    """Production configuration"""
    DEBUG = False
    ENV = 'production'
    
    # Override with environment variables in production
    SECRET_KEY = os.environ.get('SECRET_KEY')
    if not SECRET_KEY:
        raise ValueError("SECRET_KEY environment variable must be set in production")

class TestingConfig(Config):
    """Testing configuration"""
    TESTING = True
    DEBUG = True
    WTF_CSRF_ENABLED = False

# Configuration mapping
config = {
    'development': DevelopmentConfig,
    'production': ProductionConfig,
    'testing': TestingConfig,
    'default': DevelopmentConfig
}