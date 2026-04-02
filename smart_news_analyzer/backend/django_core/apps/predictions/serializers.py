"""
Serializers for Predictions App
"""

from rest_framework import serializers
from apps.predictions.models import StockPrediction, PredictionHistory, PredictionAlert
from apps.articles.serializers import ArticleListSerializer


class PredictionListSerializer(serializers.ModelSerializer):
    """Lightweight serializer for prediction lists"""
    class Meta:
        model = StockPrediction
        fields = [
            'id', 'stock_symbol', 'stock_name', 'action', 'confidence_score',
            'confidence_range', 'risk_level', 'sentiment_score', 'reasoning',
            'is_active', 'is_trending', 'created_at', 'expires_at'
        ]
        read_only_fields = ['id', 'created_at']


class PredictionDetailSerializer(serializers.ModelSerializer):
    """Detailed serializer for prediction details"""
    related_articles = ArticleListSerializer(many=True, read_only=True)
    
    class Meta:
        model = StockPrediction
        fields = [
            'id', 'stock_symbol', 'stock_name', 'action', 'confidence_score',
            'confidence_range', 'risk_level', 'sentiment_score',
            'positive_articles_count', 'negative_articles_count',
            'neutral_articles_count', 'related_articles', 'reasoning',
            'key_factors', 'current_price', 'target_price', 'stop_loss',
            'prediction_period', 'actual_direction', 'accuracy_feedback',
            'is_accurate', 'is_active', 'is_trending', 'created_at',
            'updated_at', 'expires_at', 'days_until_expiry'
        ]
        read_only_fields = ['id', 'created_at', 'updated_at']


class PredictionHistorySerializer(serializers.ModelSerializer):
    class Meta:
        model = PredictionHistory
        fields = [
            'id', 'stock_symbol', 'action', 'confidence_score',
            'actual_direction', 'was_accurate', 'sentiment_at_prediction',
            'created_at', 'resolved_at'
        ]


class PredictionAlertSerializer(serializers.ModelSerializer):
    class Meta:
        model = PredictionAlert
        fields = [
            'id', 'prediction', 'message', 'alert_type',
            'is_sent', 'sent_at', 'created_at'
        ]
