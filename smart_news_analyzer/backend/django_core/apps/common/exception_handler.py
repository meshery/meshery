"""
Common utilities and helpers
"""

from rest_framework.response import Response
from rest_framework import status
from rest_framework.exceptions import APIException
import logging

logger = logging.getLogger(__name__)


def custom_exception_handler(exc, context):
    """Custom exception handler for DRF."""
    if isinstance(exc, APIException):
        return Response(
            {
                'error': exc.detail,
                'status': exc.status_code,
            },
            status=exc.status_code
        )
    
    logger.error(f'Unexpected error: {str(exc)}', exc_info=True)
    return Response(
        {'error': 'Internal server error'},
        status=status.HTTP_500_INTERNAL_SERVER_ERROR
    )


def custom_404(request, exception=None):
    """Custom 404 view."""
    return Response(
        {'error': 'Not found'},
        status=status.HTTP_404_NOT_FOUND
    )


def custom_500(request):
    """Custom 500 view."""
    return Response(
        {'error': 'Internal server error'},
        status=status.HTTP_500_INTERNAL_SERVER_ERROR
    )
