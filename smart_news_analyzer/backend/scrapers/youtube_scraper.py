"""
YouTube news scraper using YouTube API
"""

import logging
import requests
from datetime import datetime, timedelta
from decouple import config
from scrapy_project.items import YouTubeVideoItem

logger = logging.getLogger(__name__)


class YouTubeScraper:
    """Scrape financial news from YouTube channels"""
    
    # Financial news channels to scrape
    CHANNELS = {
        'CNBC': 'UCrxmJ33LIsJXV-dJCSR8B3Q',
        'Bloomberg Quicktake': 'UCrxmJ33LIsJXV-dJCSR8B3Q',
        'Yahoo Finance': 'UBA9TbCSt_oENMi23aOMfsw',
        'Investing.com': 'UC6_syHT5oWcT-nT9gVUSM6Q',
        'The Motley Fool': 'UCz2c6H9LeMVPbZj2c9E3Igg',
    }
    
    def __init__(self):
        self.api_key = config('YOUTUBE_API_KEY', '')
        self.base_url = 'https://www.googleapis.com/youtube/v3'
        self.session = requests.Session()
    
    def scrape(self):
        """Scrape all configured channels"""
        videos = []
        
        for channel_name, channel_id in self.CHANNELS.items():
            try:
                channel_videos = self.get_channel_videos(channel_id, channel_name)
                videos.extend(channel_videos)
                logger.info(f'Scraped {len(channel_videos)} videos from {channel_name}')
            except Exception as e:
                logger.error(f'Error scraping {channel_name}: {e}')
        
        return videos
    
    def get_channel_videos(self, channel_id, channel_name):
        """Get latest videos from a channel"""
        videos = []
        
        try:
            # Get uploads playlist
            response = self.session.get(
                f'{self.base_url}/channels',
                params={
                    'part': 'contentDetails',
                    'id': channel_id,
                    'key': self.api_key
                }
            )
            
            if response.status_code != 200:
                logger.warning(f'Failed to get channel {channel_id}: {response.status_code}')
                return videos
            
            uploads_playlist = response.json()['items'][0]['contentDetails']['relatedPlaylists']['uploads']
            
            # Get videos from uploads playlist
            video_response = self.session.get(
                f'{self.base_url}/playlistItems',
                params={
                    'part': 'snippet',
                    'playlistId': uploads_playlist,
                    'maxResults': 20,
                    'key': self.api_key
                }
            )
            
            if video_response.status_code != 200:
                logger.warning(f'Failed to get videos from {channel_id}')
                return videos
            
            for item in video_response.json().get('items', []):
                snippet = item['snippet']
                
                video_item = YouTubeVideoItem()
                video_item['video_id'] = snippet['resourceId']['videoId']
                video_item['title'] = snippet['title']
                video_item['channel'] = channel_name
                video_item['description'] = snippet['description']
                video_item['published_at'] = snippet['publishedAt']
                video_item['url'] = f"https://www.youtube.com/watch?v={video_item['video_id']}"
                video_item['thumbnail_url'] = snippet['thumbnails']['default']['url']
                
                videos.append(video_item)
        
        except Exception as e:
            logger.error(f'Error processing channel {channel_id}: {e}')
        
        return videos
    
    def process_items(self):
        """Process and send videos to webhook"""
        import requests as req
        
        webhook_url = config('SCRAPERS_WEBHOOK_URL', 'http://backend:8000/api/articles/webhook/')
        
        for video in self.scrape():
            payload = {
                'title': video['title'],
                'content': video['description'],
                'url': video['url'],
                'source': 'youtube',
                'author': video['channel'],
                'published_date': video['published_at'],
                'thumbnail_url': video['thumbnail_url'],
            }
            
            try:
                response = req.post(webhook_url, json=payload, timeout=10)
                if response.status_code in [200, 201]:
                    logger.info(f'Sent YouTube video: {video["title"][:50]}')
                else:
                    logger.warning(f'Webhook failed: {response.status_code}')
            except Exception as e:
                logger.error(f'Error sending to webhook: {e}')


def main():
    scraper = YouTubeScraper()
    scraper.process_items()


if __name__ == '__main__':
    main()
