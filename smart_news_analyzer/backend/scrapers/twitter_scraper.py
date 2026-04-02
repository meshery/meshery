"""
Twitter/X news scraper using snscrape
"""

import logging
from datetime import datetime, timedelta
from decouple import config
import requests

logger = logging.getLogger(__name__)

try:
    import snscrape.modules.twitter as sntwitter
except ImportError:
    logger.warning('snscrape not installed. Twitter scraping will not work.')
    sntwitter = None


class TwitterScraper:
    """Scrape financial news from Twitter/X accounts"""
    
    ACCOUNTS = [
        '@ReutersBiz',
        '@BloombergMarkets',
        '@CNBC',
        '@WSJ',
        '@FT',
        '@YahooFinance',
        '@MarketWatch',
    ]
    
    def __init__(self):
        self.bearer_token = config('TWITTER_BEARER_TOKEN', '')
        self.session = requests.Session()
    
    def scrape_account(self, account_name):
        """Scrape tweets from an account"""
        if not sntwitter:
            logger.error('snscrape not installed')
            return []
        
        tweets = []
        query = f'from:{account_name} lang:en -is:reply -is:retweet'
        
        try:
            # Scrape last 50 tweets
            for i, tweet in enumerate(sntwitter.TwitterSearchScraper(query).get_items()):
                if i >= 50:
                    break
                
                tweet_data = {
                    'title': tweet.content[:100].replace('\n', ' '),
                    'content': tweet.content,
                    'url': f'https://twitter.com/{tweet.author.username}/status/{tweet.id}',
                    'source': 'twitter',
                    'author': tweet.author.username,
                    'published_date': tweet.date.isoformat(),
                }
                
                tweets.append(tweet_data)
                logger.debug(f'Scraped tweet from {account_name}: {tweet_data["title"][:50]}')
        
        except Exception as e:
            logger.error(f'Error scraping {account_name}: {e}')
        
        return tweets
    
    def scrape_all(self):
        """Scrape all configured accounts"""
        all_tweets = []
        
        for account in self.ACCOUNTS:
            try:
                tweets = self.scrape_account(account)
                all_tweets.extend(tweets)
                logger.info(f'Scraped {len(tweets)} tweets from {account}')
            except Exception as e:
                logger.error(f'Error processing {account}: {e}')
        
        return all_tweets
    
    def process_items(self):
        """Process tweets and send to webhook"""
        webhook_url = config('SCRAPERS_WEBHOOK_URL', 'http://backend:8000/api/articles/webhook/')
        
        tweets = self.scrape_all()
        
        for tweet in tweets:
            try:
                response = requests.post(webhook_url, json=tweet, timeout=10)
                if response.status_code in [200, 201]:
                    logger.info(f'Sent Tweet: {tweet["title"][:50]}')
                else:
                    logger.warning(f'Webhook failed: {response.status_code}')
            except Exception as e:
                logger.error(f'Error sending to webhook: {e}')


def main():
    scraper = TwitterScraper()
    scraper.process_items()


if __name__ == '__main__':
    main()
