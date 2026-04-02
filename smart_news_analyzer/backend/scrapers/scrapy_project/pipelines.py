"""
Scrapy pipelines for processing items
"""

import hashlib
import logging
import requests
from datetime import datetime
import psycopg2
from decouple import config
import time

logger = logging.getLogger(__name__)


class DuplicateCheckPipeline:
    """Check for duplicate articles"""
    
    def __init__(self):
        self.urls_seen = set()
    
    def process_item(self, item, spider):
        url = item.get('url')
        if not url:
            raise DropItem(f'No URL in item')
        
        # Create URL hash
        url_hash = hashlib.sha256(url.encode()).hexdigest()
        
        if url_hash in self.urls_seen:
            raise DropItem(f'Duplicate article: {url}')
        
        self.urls_seen.add(url_hash)
        item['url_hash'] = url_hash
        
        return item


class DataValidationPipeline:
    """Validate scraped data"""
    
    def process_item(self, item, spider):
        # Check required fields
        required_fields = ['title', 'url', 'source']
        
        for field in required_fields:
            if not item.get(field):
                raise DropItem(f'Missing required field: {field}')
        
        # Validate URL
        if not item['url'].startswith(('http://', 'https://')):
            raise DropItem(f'Invalid URL: {item["url"]}')
        
        # Clean text fields
        for text_field in ['title', 'content', 'summary']:
            if text_field in item:
                item[text_field] = item[text_field].strip()
        
        # Parse date if string
        if 'published_date' in item and isinstance(item['published_date'], str):
            try:
                item['published_date'] = datetime.fromisoformat(item['published_date'])
            except:
                item['published_date'] = datetime.now()
        
        logger.info(f'Validated article: {item["title"][:50]}')
        return item


class DatabasePipeline:
    """Store articles in database"""
    
    def __init__(self):
        self.db_config = {
            'dbname': config('DB_NAME', 'smart_news_db'),
            'user': config('DB_USER', 'postgres'),
            'password': config('DB_PASSWORD', 'postgres'),
            'host': config('DB_HOST', 'postgres'),
            'port': config('DB_PORT', '5432'),
        }
        self.connection = None
        self.cursor = None
    
    def open_spider(self, spider):
        try:
            self.connection = psycopg2.connect(**self.db_config)
            self.cursor = self.connection.cursor()
            logger.info('Database connection opened')
        except Exception as e:
            logger.error(f'Failed to connect to database: {e}')
    
    def close_spider(self, spider):
        if self.cursor:
            self.cursor.close()
        if self.connection:
            self.connection.close()
            logger.info('Database connection closed')
    
    def process_item(self, item, spider):
        try:
            # For now, just log - actual database insertion handled by webhook
            logger.debug(f'Pipeline processing: {item["title"][:50]}')
            return item
        except Exception as e:
            logger.error(f'Database error: {e}')
            raise


class DjangoWebhookPipeline:
    """Send scraped articles to Django API via webhook"""
    
    def __init__(self):
        self.webhook_url = config('SCRAPERS_WEBHOOK_URL', 'http://backend:8000/api/articles/webhook/')
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'SmartNewsScraper/1.0'
        })
    
    def process_item(self, item, spider):
        try:
            payload = {
                'title': item.get('title'),
                'content': item.get('content', item.get('summary', '')),
                'url': item.get('url'),
                'source': item.get('source'),
                'author': item.get('author'),
                'published_date': item.get('published_date', datetime.now()).isoformat() if hasattr(item.get('published_date'), 'isoformat') else str(item.get('published_date')),
                'thumbnail_url': item.get('image_url'),
            }
            
            # Send webhook
            response = self.session.post(
                self.webhook_url,
                json=payload,
                timeout=10
            )
            
            if response.status_code in [200, 201]:
                logger.info(f'Article sent to webhook: {item["title"][:50]}')
                return item
            else:
                logger.warning(f'Webhook returned status {response.status_code}: {response.text[:100]}')
                return item
        
        except requests.exceptions.RequestException as e:
            logger.error(f'Webhook request failed: {e}')
            # Don't drop the item, it will be retried
            return item
        except Exception as e:
            logger.error(f'Webhook processing error: {e}')
            return item


class DropItem(Exception):
    """Drop an item from processing"""
    pass
