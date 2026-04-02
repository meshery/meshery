"""
Models for Predictions App
"""

from django.db import models
from django.utils import timezone
from apps.articles.models import Article
import uuid


class StockPrediction(models.Model):
    """Stock market predictions based on news sentiment"""
    ACTION_CHOICES = [
        ('buy', 'BUY'),
        ('hold', 'HOLD'),
        ('sell', 'SELL'),
    ]
    
    CONFIDENCE_RANGES = [
        ('very_high', 'Very High (90-100%)'),
        ('high', 'High (70-89%)'),
        ('medium', 'Medium (50-69%)'),
        ('low', 'Low (30-49%)'),
        ('very_low', 'Very Low (0-29%)'),
    ]
    
    RISK_LEVELS = [
        ('very_low', 'Very Low'),
        ('low', 'Low'),
        ('medium', 'Medium'),
        ('high', 'High'),
        ('very_high', 'Very High'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    
    # Stock identification
    stock_symbol = models.CharField(max_length=10, db_index=True)
    stock_name = models.CharField(max_length=255, blank=True)
    
    # Prediction
    action = models.CharField(max_length=10, choices=ACTION_CHOICES)
    confidence_score = models.FloatField()  # 0-1
    confidence_range = models.CharField(max_length=20, choices=CONFIDENCE_RANGES)
    risk_level = models.CharField(max_length=20, choices=RISK_LEVELS, default='medium')
    
    # Analysis basis
    sentiment_score = models.FloatField()  # -1 to 1
    positive_articles_count = models.PositiveIntegerField(default=0)
    negative_articles_count = models.PositiveIntegerField(default=0)
    neutral_articles_count = models.PositiveIntegerField(default=0)
    
    # Related articles
    related_articles = models.ManyToManyField(Article, blank=True, related_name='predictions')
    
    # Reasoning
    reasoning = models.TextField()
    key_factors = models.JSONField(default=list)  # List of contributing factors
    
    # Price targets (if available)
    current_price = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)
    target_price = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)
    stop_loss = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)
    
    # Timeframe
    prediction_period = models.CharField(
        max_length=20,
        choices=[
            ('1day', '1 Day'),
            ('1week', '1 Week'),
            ('1month', '1 Month'),
            ('3months', '3 Months'),
        ],
        default='1week'
    )
    
    # Performance tracking
    actual_direction = models.CharField(
        max_length=10,
        choices=[('up', 'Up'), ('down', 'Down'), ('neutral', 'Neutral'), ('unknown', 'Unknown')],
        default='unknown'
    )
    accuracy_feedback = models.FloatField(null=True, blank=True)  # 0-1
    is_accurate = models.BooleanField(null=True, blank=True)
    
    # Status
    is_active = models.BooleanField(default=True)
    is_trending = models.BooleanField(default=False)
    
    # Timestamps
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    expires_at = models.DateTimeField()
    
    class Meta:
        verbose_name_plural = "Stock Predictions"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['stock_symbol', '-created_at']),
            models.Index(fields=['action', '-created_at']),
            models.Index(fields=['is_active', 'is_trending']),
        ]
        unique_together = [('stock_symbol', 'prediction_period', 'created_at')]
    
    def __str__(self):
        return f"{self.stock_symbol} - {self.action} ({self.confidence_score*100:.1f}%)"
    
    def save(self, *args, **kwargs):
        """Set expiry based on prediction period"""
        if not self.expires_at:
            periods = {
                '1day': timezone.timedelta(days=1),
                '1week': timezone.timedelta(days=7),
                '1month': timezone.timedelta(days=30),
                '3months': timezone.timedelta(days=90),
            }
            self.expires_at = self.created_at + periods.get(self.prediction_period, timezone.timedelta(days=7))
        
        super().save(*args, **kwargs)
    
    @property
    def is_expired(self):
        """Check if prediction has expired"""
        return timezone.now() > self.expires_at
    
    @property
    def days_until_expiry(self):
        """Get days until prediction expires"""
        delta = self.expires_at - timezone.now()
        return max(0, delta.days)


class PredictionHistory(models.Model):
    """Historical record of predictions for comparison"""
    stock_symbol = models.CharField(max_length=10, db_index=True)
    action = models.CharField(max_length=10)
    confidence_score = models.FloatField()
    actual_direction = models.CharField(
        max_length=10,
        choices=[('up', 'Up'), ('down', 'Down'), ('neutral', 'Neutral')],
        default='neutral'
    )
    
    was_accurate = models.BooleanField(null=True)
    sentiment_at_prediction = models.FloatField()
    
    created_at = models.DateTimeField(auto_now_add=True)
    resolved_at = models.DateTimeField(null=True, blank=True)
    
    class Meta:
        verbose_name_plural = "Prediction Histories"
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.stock_symbol} - {self.action}"


class PredictionAlert(models.Model):
    """Alerts for significant predictions"""
    prediction = models.ForeignKey(StockPrediction, on_delete=models.CASCADE, related_name='alerts')
    
    message = models.TextField()
    alert_type = models.CharField(
        max_length=20,
        choices=[
            ('strong_buy', 'Strong Buy Signal'),
            ('strong_sell', 'Strong Sell Signal'),
            ('reversal', 'Trend Reversal'),
            ('milestone', 'Price Milestone'),
        ]
    )
    
    is_sent = models.BooleanField(default=False)
    sent_at = models.DateTimeField(null=True, blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    
    def __str__(self):
        return f"{self.alert_type} - {self.prediction.stock_symbol}"
