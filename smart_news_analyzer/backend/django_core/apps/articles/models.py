"""
Models for Articles App
"""

from django.db import models
from django.urls import reverse
from django.utils.text import slugify
from taggit.managers import TaggableManager
from django.contrib.postgres.search import SearchVectorField
import hashlib
import uuid


class ArticleSource(models.Model):
    """Article source (website, platform, channel)"""
    SOURCES = [
        ('moneycontrol', 'Moneycontrol'),
        ('economic_times', 'Economic Times'),
        ('bloomberg', 'Bloomberg'),
        ('reuters', 'Reuters'),
        ('cnbc', 'CNBC'),
        ('youtube', 'YouTube'),
        ('telegram', 'Telegram'),
        ('twitter', 'Twitter/X'),
        ('other', 'Other'),
    ]
    
    name = models.CharField(max_length=100, choices=SOURCES, unique=True)
    display_name = models.CharField(max_length=255)
    url = models.URLField(null=True, blank=True)
    api_key = models.CharField(max_length=500, null=True, blank=True)
    active = models.BooleanField(default=True)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['display_name']
        verbose_name_plural = "Article Sources"
    
    def __str__(self):
        return self.display_name


class Article(models.Model):
    """Article model"""
    SENTIMENT_CHOICES = [
        ('positive', 'Positive'),
        ('negative', 'Negative'),
        ('neutral', 'Neutral'),
        ('mixed', 'Mixed'),
        ('unknown', 'Unknown'),
    ]
    
    # Core fields
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    title = models.CharField(max_length=500)
    slug = models.SlugField(unique=True, db_index=True)
    content = models.TextField()
    summary = models.TextField(null=True, blank=True)
    original_url = models.URLField(unique=True)
    url_hash = models.CharField(max_length=64, unique=True, db_index=True)
    
    # Source information
    source = models.ForeignKey(ArticleSource, on_delete=models.CASCADE, related_name='articles')
    author = models.CharField(max_length=255, null=True, blank=True)
    published_date = models.DateTimeField(db_index=True)
    
    # NLP Analysis
    sentiment = models.CharField(max_length=20, choices=SENTIMENT_CHOICES, default='unknown', db_index=True)
    sentiment_score = models.FloatField(null=True, blank=True)  # -1 to 1
    detected_stocks = models.JSONField(default=list, blank=True)  # List of stock symbols
    keywords = models.JSONField(default=list, blank=True)  # Extracted keywords
    
    # Stock impact
    stock_symbols = models.CharField(max_length=500, null=True, blank=True)  # Comma-separated
    impact_level = models.CharField(
        max_length=10,
        choices=[('high', 'High'), ('medium', 'Medium'), ('low', 'Low')],
        default='low'
    )
    
    # Metadata
    tags = TaggableManager(blank=True)
    views_count = models.PositiveIntegerField(default=0)
    relevance_score = models.FloatField(default=0.0)  # 0-1
    
    # Status
    is_analyzed = models.BooleanField(default=False)
    analysis_attempted = models.BooleanField(default=False)
    analysis_errors = models.JSONField(default=list, blank=True)
    
    # Search optimization
    search_vector = SearchVectorField(null=True, blank=True)
    
    # Timestamps
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    crawled_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        verbose_name_plural = "Articles"
        ordering = ['-published_date']
        indexes = [
            models.Index(fields=['-published_date']),
            models.Index(fields=['sentiment', '-published_date']),
            models.Index(fields=['source', '-published_date']),
            models.Index(fields=['is_analyzed']),
        ]
    
    def __str__(self):
        return self.title[:100]
    
    def save(self, *args, **kwargs):
        """Override save to generate slug and url_hash"""
        if not self.slug:
            self.slug = slugify(self.title)
        
        if not self.url_hash:
            self.url_hash = hashlib.sha256(self.original_url.encode()).hexdigest()
        
        super().save(*args, **kwargs)
    
    def get_absolute_url(self):
        return reverse('articles:detail', kwargs={'pk': self.pk})
    
    @property
    def sentiment_label(self):
        """Get human-readable sentiment label"""
        return dict(self.SENTIMENT_CHOICES).get(self.sentiment, 'Unknown')
    
    @property
    def is_recent(self):
        """Check if article is from last 24 hours"""
        from django.utils import timezone
        from datetime import timedelta
        return (timezone.now() - self.published_date) < timedelta(hours=24)


class ArticleAnalysis(models.Model):
    """Store detailed analysis results"""
    article = models.OneToOneField(Article, on_delete=models.CASCADE, related_name='analysis')
    
    summary_method = models.CharField(max_length=50, default='bart')
    sentiment_model = models.CharField(max_length=50, default='finbert')
    
    # Raw analysis results
    nlp_response = models.JSONField(default=dict)
    confidence_scores = models.JSONField(default=dict)
    
    # NER results
    entities = models.JSONField(default=list)
    
    # Analysis metadata
    processing_time = models.FloatField(null=True, blank=True)  # in seconds
    error_message = models.TextField(null=True, blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"Analysis of {self.article.title[:50]}"


class ArticleCollection(models.Model):
    """Collections/Lists for organizing articles"""
    name = models.CharField(max_length=255)
    slug = models.SlugField(unique=True)
    description = models.TextField(blank=True)
    articles = models.ManyToManyField(Article, related_name='collections')
    
    is_default = models.BooleanField(default=False)
    visibility = models.CharField(
        max_length=10,
        choices=[('public', 'Public'), ('private', 'Private')],
        default='private'
    )
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return self.name
