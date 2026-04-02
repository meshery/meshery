"""
Serializers for Articles App
"""

from rest_framework import serializers
from apps.articles.models import Article, ArticleSource, ArticleAnalysis, ArticleCollection


class ArticleSourceSerializer(serializers.ModelSerializer):
    class Meta:
        model = ArticleSource
        fields = ['name', 'display_name', 'url', 'active']


class ArticleAnalysisSerializer(serializers.ModelSerializer):
    class Meta:
        model = ArticleAnalysis
        fields = ['summary_method', 'sentiment_model', 'nlp_response', 'confidence_scores', 'entities', 'processing_time']


class ArticleListSerializer(serializers.ModelSerializer):
    """Lightweight serializer for list views"""
    source_name = serializers.CharField(source='source.display_name', read_only=True)
    sentiment_label = serializers.CharField(read_only=True)
    
    class Meta:
        model = Article
        fields = [
            'id', 'title', 'slug', 'summary', 'url_hash',
            'source', 'source_name', 'published_date',
            'sentiment', 'sentiment_label', 'sentiment_score',
            'stock_symbols', 'impact_level', 'views_count',
            'created_at'
        ]
        read_only_fields = ['id', 'url_hash', 'created_at']


class ArticleDetailSerializer(serializers.ModelSerializer):
    """Detailed serializer for article detail views"""
    source = ArticleSourceSerializer(read_only=True)
    analysis = ArticleAnalysisSerializer(read_only=True)
    tags = serializers.StringRelatedField(many=True, read_only=True)
    sentiment_label = serializers.CharField(read_only=True)
    is_recent = serializers.BooleanField(read_only=True)
    
    class Meta:
        model = Article
        fields = [
            'id', 'title', 'slug', 'content', 'summary',
            'original_url', 'url_hash', 'source', 'author',
            'published_date', 'sentiment', 'sentiment_label',
            'sentiment_score', 'detected_stocks', 'keywords',
            'stock_symbols', 'impact_level', 'tags', 'views_count',
            'relevance_score', 'is_analyzed', 'analysis',
            'created_at', 'updated_at', 'is_recent'
        ]
        read_only_fields = ['id', 'url_hash', 'created_at', 'updated_at']


class ArticleCreateSerializer(serializers.ModelSerializer):
    """Serializer for creating articles"""
    source_name = serializers.CharField(write_only=True)
    
    class Meta:
        model = Article
        fields = [
            'title', 'content', 'original_url', 'source_name',
            'author', 'published_date', 'stock_symbols'
        ]
    
    def create(self, validated_data):
        source_name = validated_data.pop('source_name')
        source, _ = ArticleSource.objects.get_or_create(name=source_name)
        validated_data['source'] = source
        return super().create(validated_data)


class ArticleCollectionSerializer(serializers.ModelSerializer):
    articles = ArticleListSerializer(many=True, read_only=True)
    
    class Meta:
        model = ArticleCollection
        fields = ['id', 'name', 'slug', 'description', 'articles', 'visibility', 'created_at']


class ArticleWebhookSerializer(serializers.Serializer):
    """Serializer for webhook data from scrapers"""
    title = serializers.CharField(max_length=500)
    content = serializers.CharField()
    url = serializers.URLField()
    source = serializers.CharField(max_length=100)
    author = serializers.CharField(max_length=255, required=False)
    published_date = serializers.DateTimeField()
    thumbnail_url = serializers.URLField(required=False)
    
    def create(self, validated_data):
        article, created = Article.objects.get_or_create(
            original_url=validated_data['url'],
            defaults={
                'title': validated_data['title'],
                'content': validated_data['content'],
                'author': validated_data.get('author', ''),
                'published_date': validated_data['published_date'],
                'source_id': ArticleSource.objects.get_or_create(
                    name=validated_data['source']
                )[0].id,
            }
        )
        return article
