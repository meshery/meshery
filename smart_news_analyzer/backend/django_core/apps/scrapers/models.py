"""
Models for Scrapers App
"""

from django.db import models
import uuid


class ScraperConfiguration(models.Model):
    """Configuration for each scraper"""
    SOURCES = [
        ('moneycontrol', 'Moneycontrol'),
        ('economic_times', 'Economic Times'),
        ('bloomberg', 'Bloomberg'),
        ('reuters', 'Reuters'),
        ('cnbc', 'CNBC'),
        ('youtube', 'YouTube'),
        ('telegram', 'Telegram'),
        ('twitter', 'Twitter/X'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    
    source = models.CharField(max_length=50, choices=SOURCES, unique=True)
    display_name = models.CharField(max_length=255)
    
    # Configuration
    is_enabled = models.BooleanField(default=True)
    api_key = models.CharField(max_length=500, null=True, blank=True)
    api_secret = models.CharField(max_length=500, null=True, blank=True)
    
    # Behavior
    scrape_interval_hours = models.IntegerField(default=4)
    timeout_seconds = models.IntegerField(default=30)
    max_retries = models.IntegerField(default=3)
    concurrent_requests = models.IntegerField(default=5)
    
    # Filters
    keywords_filter = models.JSONField(default=list, blank=True)  # Target keywords
    exclude_keywords = models.JSONField(default=list, blank=True)  # Exclude keywords
    
    # Metadata
    last_scraped_at = models.DateTimeField(null=True, blank=True)
    total_articles_scraped = models.PositiveIntegerField(default=0)
    total_errors = models.PositiveIntegerField(default=0)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "Scraper Configurations"
    
    def __str__(self):
        return self.display_name


class ScraperRun(models.Model):
    """Individual scraper execution"""
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('running', 'Running'),
        ('completed', 'Completed'),
        ('failed', 'Failed'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    
    config = models.ForeignKey(ScraperConfiguration, on_delete=models.CASCADE, related_name='runs')
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    
    # Statistics
    articles_found = models.IntegerField(default=0)
    articles_new = models.IntegerField(default=0)
    articles_duplicates = models.IntegerField(default=0)
    articles_failed = models.IntegerField(default=0)
    
    # Performance
    started_at = models.DateTimeField(auto_now_add=True)
    completed_at = models.DateTimeField(null=True, blank=True)
    
    # Error tracking
    error_message = models.TextField(blank=True)
    stack_trace = models.TextField(blank=True)
    
    # Additional logs
    logs = models.TextField(blank=True)
    
    class Meta:
        verbose_name_plural = "Scraper Runs"
        ordering = ['-started_at']
        indexes = [
            models.Index(fields=['config', '-started_at']),
            models.Index(fields=['status']),
        ]
    
    def __str__(self):
        return f"{self.config.display_name} - {self.get_status_display()}"
    
    @property
    def duration_seconds(self):
        if self.completed_at:
            return (self.completed_at - self.started_at).total_seconds()
        return None


class ScrapedURL(models.Model):
    """Track scraped URLs to avoid duplicates"""
    url = models.URLField(unique=True, db_index=True)
    url_hash = models.CharField(max_length=64, unique=True)
    
    source = models.CharField(max_length=50)
    
    last_scraped_at = models.DateTimeField(null=True, blank=True)
    scrape_count = models.IntegerField(default=0)
    
    status_code = models.IntegerField(null=True, blank=True)
    is_active = models.BooleanField(default=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "Scraped URLs"
        indexes = [
            models.Index(fields=['source', 'is_active']),
        ]
    
    def __str__(self):
        return self.url


class ScraperError(models.Model):
    """Track scraper errors"""
    ERROR_TYPES = [
        ('connection', 'Connection Error'),
        ('timeout', 'Timeout'),
        ('parsing', 'Parsing Error'),
        ('authentication', 'Authentication Error'),
        ('rate_limit', 'Rate Limit'),
        ('unknown', 'Unknown Error'),
    ]
    
    scraper_run = models.ForeignKey(ScraperRun, on_delete=models.CASCADE, related_name='errors', null=True, blank=True)
    
    source = models.CharField(max_length=50)
    error_type = models.CharField(max_length=50, choices=ERROR_TYPES)
    error_message = models.TextField()
    
    url = models.URLField(null=True, blank=True)
    stack_trace = models.TextField(blank=True)
    
    retry_count = models.IntegerField(default=0)
    
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        verbose_name_plural = "Scraper Errors"
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.source} - {self.error_type}"
