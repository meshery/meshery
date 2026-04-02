"""
URLs for Predictions App
"""

from django.urls import path, include
from rest_framework.routers import DefaultRouter
from apps.predictions.views import StockPredictionViewSet

app_name = 'predictions'

router = DefaultRouter()
router.register(r'', StockPredictionViewSet, basename='prediction')

urlpatterns = [
    path('', include(router.urls)),
]
