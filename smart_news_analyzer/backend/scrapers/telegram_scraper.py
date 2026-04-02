"""
Telegram news scraper using Telethon
"""

import asyncio
import logging
from datetime import datetime
from decouple import config
import requests

logger = logging.getLogger(__name__)

try:
    from telethon import TelegramClient
except ImportError:
    logger.warning('Telethon not installed. Telegram scraping will not work.')
    TelegramClient = None


class TelegramScraper:
    """Scrape financial news from Telegram channels"""
    
    CHANNELS = [
        '@business',
        '@financialtimes',
        '@Bloomberg',
        '@CNBC',
        '@Reuters',
    ]
    
    def __init__(self):
        self.api_id = config('TELEGRAM_API_ID', '')
        self.api_hash = config('TELEGRAM_API_HASH', '')
        self.session_name = config('TELEGRAM_SESSION_NAME', 'telegram_session')
        self.client = None
    
    async def connect(self):
        """Connect to Telegram"""
        if not TelegramClient:
            logger.error('Telethon not installed')
            return False
        
        try:
            self.client = TelegramClient(self.session_name, self.api_id, self.api_hash)
            await self.client.start()
            logger.info('Connected to Telegram')
            return True
        except Exception as e:
            logger.error(f'Failed to connect to Telegram: {e}')
            return False
    
    async def scrape_channel(self, channel_name):
        """Scrape messages from a channel"""
        messages = []
        
        try:
            channel = await self.client.get_entity(channel_name)
            
            async for message in self.client.iter_messages(channel, limit=50):
                if not message.text:
                    continue
                
                msg_data = {
                    'title': message.text[:100].replace('\n', ' '),
                    'content': message.text,
                    'url': f'https://t.me/{channel.username}/{message.id}' if hasattr(channel, 'username') else f'https://t.me/c/{channel.id}/{message.id}',
                    'source': 'telegram',
                    'author': channel_name,
                    'published_date': message.date.isoformat(),
                }
                
                messages.append(msg_data)
                logger.debug(f'Scraped message from {channel_name}: {msg_data["title"][:50]}')
        
        except Exception as e:
            logger.error(f'Error scraping {channel_name}: {e}')
        
        return messages
    
    async def scrape_all(self):
        """Scrape all configured channels"""
        all_messages = []
        
        if not await self.connect():
            return all_messages
        
        try:
            for channel in self.CHANNELS:
                messages = await self.scrape_channel(channel)
                all_messages.extend(messages)
                logger.info(f'Scraped {len(messages)} messages from {channel}')
        
        finally:
            if self.client:
                await self.client.disconnect()
        
        return all_messages
    
    def process_items(self):
        """Process messages and send to webhook"""
        webhook_url = config('SCRAPERS_WEBHOOK_URL', 'http://backend:8000/api/articles/webhook/')
        
        # Run async code
        messages = asyncio.run(self.scrape_all())
        
        for message in messages:
            try:
                response = requests.post(webhook_url, json=message, timeout=10)
                if response.status_code in [200, 201]:
                    logger.info(f'Sent Telegram message: {message["title"][:50]}')
                else:
                    logger.warning(f'Webhook failed: {response.status_code}')
            except Exception as e:
                logger.error(f'Error sending to webhook: {e}')


def main():
    scraper = TelegramScraper()
    scraper.process_items()


if __name__ == '__main__':
    main()
