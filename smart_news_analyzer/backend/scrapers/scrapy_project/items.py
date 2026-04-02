"""
Define the fields for scraped articles
"""

import scrapy


class ArticleItem(scrapy.Item):
    """Article item"""
    title = scrapy.Field()
    url = scrapy.Field()
    source = scrapy.Field()
    author = scrapy.Field()
    published_date = scrapy.Field()
    content = scrapy.Field()
    summary = scrapy.Field()
    image_url = scrapy.Field()
    keywords = scrapy.Field()
    tags = scrapy.Field()


class NewsItem(scrapy.Item):
    """General news item"""
    title = scrapy.Field()
    url = scrapy.Field()
    content = scrapy.Field()
    source = scrapy.Field()
    author = scrapy.Field()
    date = scrapy.Field()
    categories = scrapy.Field()
    sentiment = scrapy.Field()
    stocks = scrapy.Field()


class YouTubeVideoItem(scrapy.Item):
    """YouTube video item"""
    video_id = scrapy.Field()
    title = scrapy.Field()
    channel = scrapy.Field()
    description = scrapy.Field()
    published_at = scrapy.Field()
    url = scrapy.Field()
    thumbnail_url = scrapy.Field()
    duration = scrapy.Field()


class TelegramMessageItem(scrapy.Item):
    """Telegram message item"""
    chat_id = scrapy.Field()
    chat_title = scrapy.Field()
    message_id = scrapy.Field()
    text = scrapy.Field()
    date = scrapy.Field()
    url = scrapy.Field()
    media_url = scrapy.Field()


class TwitterTweetItem(scrapy.Item):
    """Twitter/X tweet item"""
    tweet_id = scrapy.Field()
    author = scrapy.Field()
    text = scrapy.Field()
    created_at = scrapy.Field()
    url = scrapy.Field()
    retweets = scrapy.Field()
    likes = scrapy.Field()
    media_urls = scrapy.Field()
