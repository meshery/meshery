"""
URLs for Articles App
"""

from django.urls import path, include
from rest_framework.routers import DefaultRouter
from apps.articles.views import ArticleViewSet, ArticleSourceViewSet, ArticleCollectionViewSet

app_name = 'articles'

router = DefaultRouter()
router.register(r'', ArticleViewSet, basename='article')
router.register(r'sources', ArticleSourceViewSet, basename='source')
router.register(r'collections', ArticleCollectionViewSet, basename='collection')

urlpatterns = [
    path('', include(router.urls)),
]
