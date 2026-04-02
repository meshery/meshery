"""
Models for Users App
"""

from django.db import models
from django.contrib.auth.models import User
import uuid


class UserProfile(models.Model):
    """Extended user profile"""
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='profile')
    
    # Profile information
    bio = models.TextField(blank=True)
    avatar = models.ImageField(upload_to='avatars/', null=True, blank=True)
    
    # Preferences
    preferred_stocks = models.JSONField(default=list)  # List of stock symbols
    notification_email = models.EmailField(blank=True)
    
    # Settings
    enable_email_alerts = models.BooleanField(default=True)
    enable_push_notifications = models.BooleanField(default=False)
    notification_frequency = models.CharField(
        max_length=20,
        choices=[
            ('immediate', 'Immediate'),
            ('daily', 'Daily'),
            ('weekly', 'Weekly'),
        ],
        default='daily'
    )
    
    # Statistics
    total_views = models.PositiveIntegerField(default=0)
    total_favorites = models.PositiveIntegerField(default=0)
    
    # Timestamps
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "User Profiles"
    
    def __str__(self):
        return f"{self.user.username}'s Profile"


class UserFavorite(models.Model):
    """User's favorite articles"""
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='favorites')
    article = models.ForeignKey('articles.Article', on_delete=models.CASCADE, related_name='favorited_by')
    
    added_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        verbose_name_plural = "User Favorites"
        unique_together = ['user', 'article']
    
    def __str__(self):
        return f"{self.user.username} favorited {self.article.title[:50]}"


class UserReadingHistory(models.Model):
    """Track user reading history"""
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='reading_history')
    article = models.ForeignKey('articles.Article', on_delete=models.CASCADE, related_name='read_by')
    
    read_at = models.DateTimeField(auto_now_add=True)
    time_spent_seconds = models.IntegerField(default=0)
    is_completed = models.BooleanField(default=False)
    
    class Meta:
        verbose_name_plural = "User Reading Histories"
        unique_together = ['user', 'article']
    
    def __str__(self):
        return f"{self.user.username} read {self.article.title[:50]}"


class UserAlert(models.Model):
    """User-created alerts for specific stocks"""
    ALERT_TYPES = [
        ('price_target', 'Price Target'),
        ('sentiment_shift', 'Sentiment Shift'),
        ('article_mention', 'Article Mention'),
        ('prediction', 'Prediction Change'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='alerts')
    
    stock_symbol = models.CharField(max_length=10)
    alert_type = models.CharField(max_length=50, choices=ALERT_TYPES)
    
    # Alert parameters
    condition = models.JSONField(default=dict)  # e.g., {'type': 'above', 'value': 100}
    
    # Status
    is_active = models.BooleanField(default=True)
    triggered_count = models.IntegerField(default=0)
    last_triggered_at = models.DateTimeField(null=True, blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "User Alerts"
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.user.username} - {self.stock_symbol} ({self.alert_type})"


class UserNotification(models.Model):
    """Notifications for users"""
    NOTIFICATION_TYPES = [
        ('article', 'New Article'),
        ('prediction', 'New Prediction'),
        ('alert', 'Alert Triggered'),
        ('system', 'System'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='notifications')
    
    notification_type = models.CharField(max_length=50, choices=NOTIFICATION_TYPES)
    title = models.CharField(max_length=255)
    message = models.TextField()
    
    # Content references
    article = models.ForeignKey('articles.Article', on_delete=models.SET_NULL, null=True, blank=True)
    prediction = models.ForeignKey('predictions.StockPrediction', on_delete=models.SET_NULL, null=True, blank=True)
    
    # Status
    is_read = models.BooleanField(default=False)
    is_sent = models.BooleanField(default=False)
    
    created_at = models.DateTimeField(auto_now_add=True)
    sent_at = models.DateTimeField(null=True, blank=True)
    read_at = models.DateTimeField(null=True, blank=True)
    
    class Meta:
        verbose_name_plural = "User Notifications"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['user', '-created_at']),
            models.Index(fields=['is_read']),
        ]
    
    def __str__(self):
        return f"{self.user.username} - {self.title}"


class UserActivityLog(models.Model):
    """Log user activities for analytics"""
    ACTIVITY_TYPES = [
        ('view_article', 'View Article'),
        ('search', 'Search'),
        ('filter', 'Filter'),
        ('export', 'Export'),
        ('login', 'Login'),
        ('logout', 'Logout'),
    ]
    
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='activity_logs')
    activity_type = models.CharField(max_length=50, choices=ACTIVITY_TYPES)
    
    details = models.JSONField(default=dict, blank=True)
    ip_address = models.GenericIPAddressField(null=True, blank=True)
    user_agent = models.TextField(blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        verbose_name_plural = "User Activity Logs"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['user', '-created_at']),
            models.Index(fields=['activity_type']),
        ]
    
    def __str__(self):
        return f"{self.user.username} - {self.activity_type}"
