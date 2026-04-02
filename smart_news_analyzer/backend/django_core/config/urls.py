"""
URL Configuration for Smart News Analyzer
"""

from django.contrib import admin
from django.conf import settings
from django.conf.urls.static import static
from django.urls import path, include, re_path
from rest_framework import permissions
from drf_spectacular.views import SpectacularAPIView, SpectacularSwaggerView, SpectacularRedocView
from rest_framework_simplejwt.views import (
    TokenObtainPairView,
    TokenRefreshView,
)

urlpatterns = [
    # Admin
    path('admin/', admin.site.urls),

    # Health Check
    path('health/', lambda request: __import__('django.http', fromlist=['JsonResponse']).JsonResponse({'status': 'healthy'})),

    # API Documentation
    path('api/schema/', SpectacularAPIView.as_view(), name='schema'),
    path('api/docs/', SpectacularSwaggerView.as_view(url_name='schema'), name='swagger-ui'),
    path('api/redoc/', SpectacularRedocView.as_view(url_name='schema'), name='redoc'),

    # Authentication
    path('api/auth/token/', TokenObtainPairView.as_view(), name='token_obtain_pair'),
    path('api/auth/token/refresh/', TokenRefreshView.as_view(), name='token_refresh'),

    # API v1
    path('api/v1/articles/', include('apps.articles.urls', namespace='articles')),
    path('api/v1/predictions/', include('apps.predictions.urls', namespace='predictions')),
    path('api/v1/tasks/', include('apps.tasks.urls', namespace='tasks')),
    path('api/v1/scrapers/', include('apps.scrapers.urls', namespace='scrapers')),
    path('api/v1/users/', include('apps.users.urls', namespace='users')),

    # API v1 (without versioning in URL for compatibility)
    path('api/articles/', include('apps.articles.urls')),
    path('api/predictions/', include('apps.predictions.urls')),
    path('api/tasks/', include('apps.tasks.urls')),
    path('api/scrapers/', include('apps.scrapers.urls')),
    path('api/users/', include('apps.users.urls')),

    # Django Debug Toolbar (development only)
]

if settings.DEBUG:
    urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
    urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)
    
    import debug_toolbar
    urlpatterns = [
        path('__debug__/', include(debug_toolbar.urls)),
    ] + urlpatterns

# Custom 404 and 500 handlers
handler404 = 'apps.common.views.custom_404'
handler500 = 'apps.common.views.custom_500'
