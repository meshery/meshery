"""
Views for Predictions App
"""

from rest_framework import viewsets, status, permissions
from rest_framework.decorators import action
from rest_framework.response import Response
from django.utils import timezone
from django.db.models import Q
from apps.predictions.models import StockPrediction, PredictionHistory
from apps.predictions.serializers import (
    PredictionListSerializer, PredictionDetailSerializer,
    PredictionHistorySerializer
)


class StockPredictionViewSet(viewsets.ReadOnlyModelViewSet):
    """ViewSet for Stock Predictions"""
    queryset = StockPrediction.objects.filter(is_active=True)
    permission_classes = [permissions.IsAuthenticated]
    filterset_fields = ['stock_symbol', 'action', 'risk_level']
    search_fields = ['stock_symbol', 'stock_name']
    ordering_fields = ['confidence_score', 'created_at', '-created_at']
    ordering = ['-created_at']
    
    def get_serializer_class(self):
        if self.action == 'retrieve':
            return PredictionDetailSerializer
        return PredictionListSerializer
    
    @action(detail=False, methods=['GET'])
    def trending(self, request):
        """Get trending predictions"""
        predictions = self.get_queryset().filter(
            is_trending=True,
            expires_at__gt=timezone.now()
        )[:10]
        serializer = self.get_serializer(predictions, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def by_symbol(self, request):
        """Get predictions for a specific stock"""
        symbol = request.query_params.get('symbol', None)
        if not symbol:
            return Response(
                {'error': 'symbol parameter required'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        predictions = self.get_queryset().filter(
            stock_symbol__iexact=symbol,
            expires_at__gt=timezone.now()
        )
        
        serializer = self.get_serializer(predictions, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def buy_signals(self, request):
        """Get BUY recommendations"""
        predictions = self.get_queryset().filter(
            action='buy',
            confidence_score__gte=0.6,
            expires_at__gt=timezone.now()
        )[:20]
        
        serializer = self.get_serializer(predictions, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def sell_signals(self, request):
        """Get SELL recommendations"""
        predictions = self.get_queryset().filter(
            action='sell',
            confidence_score__gte=0.6,
            expires_at__gt=timezone.now()
        )[:20]
        
        serializer = self.get_serializer(predictions, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def hold_signals(self, request):
        """Get HOLD recommendations"""
        predictions = self.get_queryset().filter(
            action='hold',
            expires_at__gt=timezone.now()
        )[:20]
        
        serializer = self.get_serializer(predictions, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def sentiment_distribution(self, request):
        """Get distribution of prediction actions"""
        queryset = self.get_queryset().filter(expires_at__gt=timezone.now())
        
        distribution = {
            'buy': queryset.filter(action='buy').count(),
            'sell': queryset.filter(action='sell').count(),
            'hold': queryset.filter(action='hold').count(),
        }
        
        return Response(distribution)
    
    @action(detail=False, methods=['GET'])
    def market_sentiment(self, request):
        """Get overall market sentiment based on predictions"""
        queryset = self.get_queryset().filter(expires_at__gt=timezone.now())
        
        buy_count = queryset.filter(action='buy').count()
        sell_count = queryset.filter(action='sell').count()
        hold_count = queryset.filter(action='hold').count()
        
        total = buy_count + sell_count + hold_count
        
        if total == 0:
            return Response({'sentiment': 'neutral', 'details': {}})
        
        sentiment_score = (buy_count - sell_count) / total
        
        if sentiment_score > 0.3:
            sentiment = 'bullish'
        elif sentiment_score < -0.3:
            sentiment = 'bearish'
        else:
            sentiment = 'neutral'
        
        return Response({
            'sentiment': sentiment,
            'score': sentiment_score,
            'details': {
                'buy': buy_count,
                'sell': sell_count,
                'hold': hold_count,
            }
        })
    
    @action(detail=False, methods=['GET'])
    def history(self, request):
        """Get prediction history"""
        history = PredictionHistory.objects.all().order_by('-created_at')[:100]
        serializer = PredictionHistorySerializer(history, many=True)
        return Response(serializer.data)
