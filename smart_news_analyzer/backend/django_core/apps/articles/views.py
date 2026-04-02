"""
Views/ViewSets for Articles App
"""

from rest_framework import viewsets, status, permissions
from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework.filters import SearchFilter, OrderingFilter
from django_filters.rest_framework import DjangoFilterBackend
from django.db.models import Q
from apps.articles.models import Article, ArticleSource, ArticleCollection
from apps.articles.serializers import (
    ArticleListSerializer, ArticleDetailSerializer, ArticleCreateSerializer,
    ArticleCollectionSerializer, ArticleWebhookSerializer
)


class ArticleViewSet(viewsets.ModelViewSet):
    """ViewSet for Article API"""
    queryset = Article.objects.all()
    permission_classes = [permissions.IsAuthenticated]
    filter_backends = [DjangoFilterBackend, SearchFilter, OrderingFilter]
    filterset_fields = ['sentiment', 'source', 'is_analyzed']
    search_fields = ['title', 'content', 'keywords', 'stock_symbols']
    ordering_fields = ['published_date', 'created_at', 'views_count', 'sentiment_score']
    ordering = ['-published_date']
    
    def get_serializer_class(self):
        """Return appropriate serializer based on action"""
        if self.action == 'create' or self.action == 'webhook':
            return ArticleCreateSerializer
        elif self.action == 'retrieve':
            return ArticleDetailSerializer
        return ArticleListSerializer
    
    def retrieve(self, request, *args, **kwargs):
        """Override retrieve to increment view count"""
        instance = self.get_object()
        instance.views_count += 1
        instance.save(update_fields=['views_count'])
        
        serializer = self.get_serializer(instance)
        return Response(serializer.data)
    
    @action(detail=False, methods=['POST'], permission_classes=[permissions.AllowAny])
    def webhook(self, request):
        """Webhook endpoint for scrapers to submit articles"""
        serializer = ArticleWebhookSerializer(data=request.data)
        if serializer.is_valid():
            article = serializer.save()
            return Response(
                ArticleDetailSerializer(article).data,
                status=status.HTTP_201_CREATED
            )
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
    
    @action(detail=False, methods=['GET'])
    def trending(self, request):
        """Get trending articles"""
        articles = self.get_queryset().filter(
            is_analyzed=True
        ).order_by('-views_count', '-sentiment_score')[:10]
        
        serializer = self.get_serializer(articles, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def by_sentiment(self, request):
        """Get articles by sentiment"""
        sentiment = request.query_params.get('sentiment', None)
        if not sentiment:
            return Response(
                {'error': 'sentiment parameter required'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        articles = self.get_queryset().filter(sentiment=sentiment)[:20]
        serializer = self.get_serializer(articles, many=True)
        return Response(serializer.data)
    
    @action(detail=False, methods=['GET'])
    def search(self, request):
        """Search articles"""
        query = request.query_params.get('q', '')
        if not query or len(query) < 2:
            return Response(
                {'error': 'Query parameter q required (min 2 chars)'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        articles = self.get_queryset().filter(
            Q(title__icontains=query) |
            Q(content__icontains=query) |
            Q(keywords__contains=[query])
        )[:50]
        
        serializer = self.get_serializer(articles, many=True)
        return Response(serializer.data)
    
    @action(detail=True, methods=['POST'])
    def favorite(self, request, pk=None):
        """Add article to favorites"""
        article = self.get_object()
        from apps.users.models import UserFavorite
        
        favorite, created = UserFavorite.objects.get_or_create(
            user=request.user,
            article=article
        )
        
        message = "Added to favorites" if created else "Already in favorites"
        return Response({'message': message}, status=status.HTTP_200_OK)
    
    @action(detail=True, methods=['DELETE'])
    def unfavorite(self, request, pk=None):
        """Remove article from favorites"""
        article = self.get_object()
        from apps.users.models import UserFavorite
        
        UserFavorite.objects.filter(
            user=request.user,
            article=article
        ).delete()
        
        return Response({'message': 'Removed from favorites'}, status=status.HTTP_200_OK)
    
    @action(detail=False, methods=['GET'])
    def my_favorites(self, request):
        """Get user's favorite articles"""
        from apps.users.models import UserFavorite
        
        favorites = UserFavorite.objects.filter(
            user=request.user
        ).values_list('article_id', flat=True)
        
        articles = Article.objects.filter(id__in=favorites)
        serializer = self.get_serializer(articles, many=True)
        return Response(serializer.data)


class ArticleSourceViewSet(viewsets.ReadOnlyModelViewSet):
    """ViewSet for Article Sources"""
    queryset = ArticleSource.objects.filter(active=True)
    serializer_class = ArticleSourceSerializer
    permission_classes = [permissions.IsAuthenticated]


class ArticleCollectionViewSet(viewsets.ModelViewSet):
    """ViewSet for Article Collections"""
    queryset = ArticleCollection.objects.all()
    serializer_class = ArticleCollectionSerializer
    permission_classes = [permissions.IsAuthenticated]
    
    def perform_create(self, serializer):
        serializer.save()
