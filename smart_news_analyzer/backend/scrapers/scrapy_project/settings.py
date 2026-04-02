"""
Scrapy settings for scrapy_project project.
"""

import os
from decouple import config

# Project name
BOT_NAME = 'smart_news_analyzer'

# Spider modules
SPIDER_MODULES = ['scrapy_project.spiders']
NEWSPIDER_MODULE = 'scrapy_project.spiders'

# Accept robots.txt
ROBOTSTXT_OBEY = False

# Configure maximum concurrent requests
CONCURRENT_REQUESTS = config('SCRAPY_CONCURRENT_REQUESTS', default=16, cast=int)
CONCURRENT_REQUESTS_PER_DOMAIN = config('SCRAPY_PIPELINE_CONCURRENCY', default=2, cast=int)

# Download delay
DOWNLOAD_DELAY = 3

# User-Agent
USER_AGENT = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'

# Middlewares
DOWNLOADER_MIDDLEWARES = {
    'scrapy.downloadermiddlewares.robotstxt.RobotsTxtMiddleware': 543,
    'scrapy.downloadermiddlewares.httpproxy.HttpProxyMiddleware': 543,
    'scrapy_rotating_proxy.middlewares.RotatingProxyMiddleware': 610,
    'scrapy_rotating_proxy.middlewares.BanDetectionMiddleware': 603,
    'scrapy_user_agent.middlewares.RandomUserAgentMiddleware': 560,
}

# Pipeline
ITEM_PIPELINES = {
    'scrapy_project.pipelines.DuplicateCheckPipeline': 100,
    'scrapy_project.pipelines.DataValidationPipeline': 200,
    'scrapy_project.pipelines.DatabasePipeline': 300,
    'scrapy_project.pipelines.DjangoWebhookPipeline': 400,
}

# Enable and configure logging
LOG_LEVEL = config('LOG_LEVEL', default='INFO')
LOG_FORMAT = '%(asctime)s [%(name)s] %(levelname)s: %(message)s'
LOG_FILE = config('LOG_FILE', default='scrapy_project.log')

# Twisted reactor
TWISTED_REACTOR = 'twisted.internet.asyncioreactor.AsyncioSelectorReactor'

# Configure external APIs
YOUTUBE_API_KEY = config('YOUTUBE_API_KEY', default='')
TELEGRAM_API_ID = config('TELEGRAM_API_ID', default='')
TELEGRAM_API_HASH = config('TELEGRAM_API_HASH', default='')
TWITTER_BEARER_TOKEN = config('TWITTER_BEARER_TOKEN', default='')

# Django settings
DJANGO_API_URL = config('DJANGO_API_URL', default='http://backend:8000/api')

# Timeout
DOWNLOAD_TIMEOUT = config('SCRAPER_TIMEOUT', default=30, cast=int)

# Retries
RETRY_TIMES = config('SCRAPER_RETRY_ATTEMPTS', default=3, cast=int)
RETRY_HTTP_CODES = [500, 502, 503, 504, 408]

# Memory usage
TELNETCONSOLE_ENABLED = True
TELNETCONSOLE_PORT = [6023, 6024, 6025]

# AutoThrottle
AUTOTHROTTLE_ENABLED = True
AUTOTHROTTLE_START_DELAY = 5
AUTOTHROTTLE_MAX_DELAY = 60
AUTOTHROTTLE_TARGET_CONCURRENCY = 8.0

# Compression
COMPRESSION_ENABLED = True
COMPRESSION_LEVEL = 6

# Cache DNS results
DNSCACHE_ENABLED = True
DNSCACHE_SIZE = 10000

# HTTP cache
HTTPCACHE_ENABLED = True
HTTPCACHE_EXPIRATION_SECS = 0
HTTPCACHE_DIR = 'httpcache'
HTTPCACHE_IGNORE_HTTP_CODES = [400, 403, 404, 408, 500, 502, 503, 504]
