"""
Models for Tasks App
"""

from django.db import models
from django.contrib.auth.models import User
from django_celery_results.models import TaskResult
import uuid


class ScraperTask(models.Model):
    """Track scraper execution tasks"""
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('running', 'Running'),
        ('completed', 'Completed'),
        ('failed', 'Failed'),
        ('stopped', 'Stopped'),
    ]
    
    SOURCES = [
        ('moneycontrol', 'Moneycontrol'),
        ('economic_times', 'Economic Times'),
        ('bloomberg', 'Bloomberg'),
        ('reuters', 'Reuters'),
        ('cnbc', 'CNBC'),
        ('youtube', 'YouTube'),
        ('telegram', 'Telegram'),
        ('twitter', 'Twitter/X'),
        ('all', 'All Sources'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    
    # Task identification
    celery_task_id = models.CharField(max_length=100, unique=True, null=True, blank=True)
    source = models.CharField(max_length=50, choices=SOURCES)
    task_type = models.CharField(
        max_length=50,
        choices=[
            ('scrape', 'Scraping'),
            ('analyze', 'Analysis'),
            ('predict', 'Prediction'),
        ]
    )
    
    # Status and progress
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    progress_percentage = models.IntegerField(default=0)
    
    # Results
    articles_processed = models.IntegerField(default=0)
    articles_created = models.IntegerField(default=0)
    articles_updated = models.IntegerField(default=0)
    
    # Errors and logs
    error_message = models.TextField(blank=True)
    logs = models.TextField(blank=True)
    
    # Timing
    started_at = models.DateTimeField(null=True, blank=True)
    completed_at = models.DateTimeField(null=True, blank=True)
    
    # Metadata
    executed_by = models.ForeignKey(User, null=True, blank=True, on_delete=models.SET_NULL)
    is_scheduled = models.BooleanField(default=False)
    retry_count = models.IntegerField(default=0)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "Scraper Tasks"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['status', '-created_at']),
            models.Index(fields=['celery_task_id']),
        ]
    
    def __str__(self):
        return f"{self.source} - {self.get_task_type_display()} ({self.get_status_display()})"
    
    @property
    def duration(self):
        """Get task duration in seconds"""
        if self.started_at and self.completed_at:
            delta = self.completed_at - self.started_at
            return delta.total_seconds()
        return None
    
    @property
    def is_running(self):
        return self.status == 'running'
    
    @property
    def is_completed(self):
        return self.status in ['completed', 'failed']


class AnalysisTask(models.Model):
    """Track NLP analysis tasks"""
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('processing', 'Processing'),
        ('completed ', 'Completed'),
        ('failed', 'Failed'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    
    celery_task_id = models.CharField(max_length=100, unique=True, null=True, blank=True)
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    
    # Target
    articles_count = models.IntegerField()
    articles_completed = models.IntegerField(default=0)
    
    # Results
    success_count = models.IntegerField(default=0)
    error_count = models.IntegerField(default=0)
    
    # Timing
    started_at = models.DateTimeField(null=True, blank=True)
    completed_at = models.DateTimeField(null=True, blank=True)
    
    error_message = models.TextField(blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        verbose_name_plural = "Analysis Tasks"
        ordering = ['-created_at']
    
    def __str__(self):
        return f"Analysis Task - {self.articles_count} articles ({self.get_status_display()})"


class TaskLog(models.Model):
    """Detailed logs for tasks"""
    task_id = models.CharField(max_length=100, db_index=True)
    task_type = models.CharField(max_length=50)
    
    level = models.CharField(
        max_length=10,
        choices=[
            ('debug', 'Debug'),
            ('info', 'Info'),
            ('warning', 'Warning'),
            ('error', 'Error'),
        ]
    )
    
    message = models.TextField()
    metadata = models.JSONField(default=dict, blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        verbose_name_plural = "Task Logs"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['task_id', '-created_at']),
        ]
    
    def __str__(self):
        return f"{self.task_type} - {self.level}: {self.message[:50]}"
