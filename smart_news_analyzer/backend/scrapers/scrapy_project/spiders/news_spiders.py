"""
Scrapy spiders for financial news sources
"""

import scrapy
import logging
from datetime import datetime
from scrapy_project.items import ArticleItem

logger = logging.getLogger(__name__)


class MoneycontrolSpider(scrapy.Spider):
    """Scrape Moneycontrol.com news"""
    
    name = 'moneycontrol'
    allowed_domains = ['moneycontrol.com']
    start_urls = [
        'https://www.moneycontrol.com/news/markets/',
        'https://www.moneycontrol.com/news/business/',
        'https://www.moneycontrol.com/news/stocks/',
    ]
    
    def parse(self, response):
        # Extract news items
        articles = response.css('div.storyCardImg')
        
        for article in articles:
            item = ArticleItem()
            item['title'] = article.css('a::attr(title)').get('').strip()
            item['url'] = response.urljoin(article.css('a::attr(href)').get(''))
            item['author'] = 'Moneycontrol'
            item['source'] = 'moneycontrol'
            item['published_date'] = datetime.now()
            
            if item['url']:
                yield scrapy.Request(
                    item['url'],
                    callback=self.parse_article,
                    meta={'item': item}
                )
    
    def parse_article(self, response):
        item = response.meta['item']
        
        item['content'] = ' '.join(response.css('div.article_content p::text').getall())
        item['summary'] = response.css('meta[name="description"]::attr(content)').get('')
        item['image_url'] = response.css('div.artimg img::attr(src)').get('')
        
        logger.info(f'Parsed article: {item["title"][:50]}')
        yield item


class EconomicTimesSpider(scrapy.Spider):
    """Scrape Economic Times news"""
    
    name = 'economic_times'
    allowed_domains = ['economictimes.indiatimes.com']
    start_urls = [
        'https://economictimes.indiatimes.com/markets',
        'https://economictimes.indiatimes.com/news/stocks',
        'https://economictimes.indiatimes.com/news/economy',
    ]
    
    def parse(self, response):
        articles = response.css('article.eachstory')
        
        for article in articles:
            item = ArticleItem()
            item['title'] = article.css('h2 a::text').get('').strip()
            item['url'] = response.urljoin(article.css('h2 a::attr(href)').get(''))
            item['author'] = article.css('span.auth_name::text').get('Economic Times')
            item['source'] = 'economic_times'
            item['published_date'] = datetime.now()
            
            if item['url']:
                yield scrapy.Request(
                    item['url'],
                    callback=self.parse_article,
                    meta={'item': item}
                )
    
    def parse_article(self, response):
        item = response.meta['item']
        
        item['content'] = ' '.join(response.css('div.el-article-text p::text').getall())
        item['summary'] = response.css('meta[name="description"]::attr(content)').get('')
        
        yield item


class BloombergSpider(scrapy.Spider):
    """Scrape Bloomberg Markets news"""
    
    name = 'bloomberg'
    allowed_domains = ['bloomberg.com']
    start_urls = [
        'https://www.bloomberg.com/markets',
        'https://www.bloomberg.com/news/markets',
    ]
    
    custom_settings = {
        'DOWNLOADER_MIDDLEWARES': {
            'scrapy.downloadermiddlewares.robotstxt.RobotsTxtMiddleware': None,
        }
    }
    
    def parse(self, response):
        # Bloomberg uses JavaScript - would require splash or selenium
        logger.warning('Bloomberg requires JavaScript rendering - using stub')
        return []


class ReutersSpider(scrapy.Spider):
    """Scrape Reuters Business news"""
    
    name = 'reuters'
    allowed_domains = ['reuters.com']
    start_urls = [
        'https://www.reuters.com/finance/',
        'https://www.reuters.com/business/',
    ]
    
    def parse(self, response):
        articles = response.css('article')
        
        for article in articles:
            item = ArticleItem()
            item['title'] = article.css('h3::text').get('').strip()
            item['url'] = response.urljoin(article.css('a::attr(href)').get(''))
            item['source'] = 'reuters'
            item['author'] = article.css('span.article-byline::text').get('Reuters')
            item['published_date'] = datetime.now()
            
            if item['url']:
                yield scrapy.Request(
                    item['url'],
                    callback=self.parse_article,
                    meta={'item': item}
                )
    
    def parse_article(self, response):
        item = response.meta['item']
        
        item['content'] = ' '.join(response.css('p::text').getall()[10:])  # Skip header paragraphs
        item['summary'] = response.css('meta[name="description"]::attr(content)').get('')
        
        yield item


class CNBCSpider(scrapy.Spider):
    """Scrape CNBC Markets news"""
    
    name = 'cnbc'
    allowed_domains = ['cnbc.com']
    start_urls = [
        'https://www.cnbc.com/markets/',
        'https://www.cnbc.com/business/',
    ]
    
    def parse(self, response):
        articles = response.css('div.Card-container')
        
        for article in articles:
            item = ArticleItem()
            item['title'] = article.css('a::attr(aria-label)').get('').strip()
            item['url'] = response.urljoin(article.css('a::attr(href)').get(''))
            item['source'] = 'cnbc'
            item['author'] = 'CNBC'
            item['published_date'] = datetime.now()
            
            if item['url'] and item['title']:
                yield scrapy.Request(
                    item['url'],
                    callback=self.parse_article,
                    meta={'item': item}
                )
    
    def parse_article(self, response):
        item = response.meta['item']
        
        item['content'] = ' '.join(response.css('article p::text').getall())
        item['summary'] = response.css('meta[name="description"]::attr(content)').get('')
        
        yield item


# Combined spider that runs all sources
class NewsAggregatorSpider(scrapy.Spider):
    """Aggregate news from all sources"""
    
    name = 'news_aggregator'
    
    def start_requests(self):
        spiders = [MoneycontrolSpider, EconomicTimesSpider, ReutersSpider, CNBCSpider]
        
        for spider_class in spiders:
            spider = spider_class()
            for request in spider.start_requests():
                yield request
