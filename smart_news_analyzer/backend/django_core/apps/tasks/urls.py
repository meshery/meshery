"""
URLs for Tasks, Scrapers, Users
"""

from django.urls import path

# Tasks URLs
tasks_patterns = [
    path('', lambda request: __import__('django.http', fromlist=['JsonResponse']).JsonResponse({'status': 'ok'})),
]

# Scrapers URLs
scrapers_patterns = [
    path('', lambda request: __import__('django.http', fromlist=['JsonResponse']).JsonResponse({'status': 'ok'})),
]

# Users URLs  
users_patterns = [
    path('', lambda request: __import__('django.http', fromlist=['JsonResponse']).JsonResponse({'status': 'ok'})),
]
