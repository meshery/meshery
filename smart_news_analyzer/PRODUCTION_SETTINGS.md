# Production Settings & Security Hardening

## Django Settings for Production

### 1. Essential Security Settings

```python
# backend/django_core/config/settings.py (production overrides)

DEBUG = False
ALLOWED_HOSTS = [
    'yourdomain.com',
    'www.yourdomain.com',
    'api.yourdomain.com'
]

# Security
SECRET_KEY = os.environ.get('DJANGO_SECRET_KEY')  # Must be unique per environment
SECURE_SSL_REDIRECT = True
SESSION_COOKIE_SECURE = True
CSRF_COOKIE_SECURE = True
SECURE_BROWSER_XSS_FILTER = True
SECURE_CONTENT_SECURITY_POLICY = {
    'default-src': ("'self'",),
    'script-src': ("'self'", "'unsafe-inline'"),
    'style-src': ("'self'", "'unsafe-inline'"),
}

# HSTS (HTTP Strict Transport Security)
SECURE_HSTS_SECONDS = 31536000
SECURE_HSTS_INCLUDE_SUBDOMAINS = True
SECURE_HSTS_PRELOAD = True

# CORS - Restrict to frontend domain only
CORS_ALLOWED_ORIGINS = [
    'https://yourdomain.com',
    'https://www.yourdomain.com',
]
CORS_ALLOW_CREDENTIALS = True
```

### 2. Database Connection Pooling

```python
# backend/django_core/config/settings.py

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': os.environ.get('DB_NAME'),
        'USER': os.environ.get('DB_USER'),
        'PASSWORD': os.environ.get('DB_PASSWORD'),
        'HOST': os.environ.get('DB_HOST'),
        'PORT': os.environ.get('DB_PORT', '5432'),
        'CONN_MAX_AGE': 600,  # Connection pooling
        'OPTIONS': {
            'connect_timeout': 10,
        }
    }
}

# Connection pooling with PgBouncer
if os.environ.get('USE_PGBOUNCER'):
    DATABASES['default']['HOST'] = 'pgbouncer'
    DATABASES['default']['PORT'] = '6432'
```

### 3. Caching Configuration

```python
CACHES = {
    'default': {
        'BACKEND': 'django_redis.cache.RedisCache',
        'LOCATION': 'redis://{}:{}/{}'.format(
            os.environ.get('REDIS_HOST', 'redis'),
            os.environ.get('REDIS_PORT', '6379'),
            os.environ.get('REDIS_DB_CACHE', '1')
        ),
        'OPTIONS': {
            'CLIENT_CLASS': 'django_redis.client.DefaultClient',
            'CONNECTION_POOL_KWARGS': {'max_connections': 50},
            'SOCKET_CONNECT_TIMEOUT': 5,
            'SOCKET_TIMEOUT': 5,
            'IGNORE_EXCEPTIONS': True,
        },
        'KEY_PREFIX': 'smart-news-prod',
        'TIMEOUT': 3600,
    }
}

# Session backend
SESSION_ENGINE = 'django.contrib.sessions.backends.cache'
SESSION_CACHE_ALIAS = 'default'
```

### 4. Logging for Production

```python
LOGGING = {
    'version': 1,
    'disable_existing_loggers': False,
    'formatters': {
        'verbose': {
            'format': '{levelname} {asctime} {module} {process:d} {thread:d} {message}',
            'style': '{',
        },
        'json': {
            '()': 'pythonjsonlogger.jsonlogger.JsonFormatter',
            'format': '%(asctime)s %(name)s %(levelname)s %(message)s'
        }
    },
    'handlers': {
        'file': {
            'level': 'ERROR',
            'class': 'logging.handlers.RotatingFileHandler',
            'filename': '/var/log/django/error.log',
            'maxBytes': 1024 * 1024 * 10,  # 10MB
            'backupCount': 5,
            'formatter': 'json',
        },
        'sentry': {
            'level': 'ERROR',
            'class': 'sentry_sdk.integrations.logging.EventHandler',
        }
    },
    'root': {
        'handlers': ['file', 'sentry'],
        'level': 'INFO',
    },
}

# Sentry integration
import sentry_sdk
sentry_sdk.init(
    dsn=os.environ.get('SENTRY_DSN'),
    traces_sample_rate=0.1,
    environment=os.environ.get('ENVIRONMENT', 'production')
)
```

### 5. Celery Production Configuration

```python
# backend/django_core/config/settings.py

CELERY_BROKER_URL = 'redis://{}:{}/{}'.format(
    os.environ.get('REDIS_HOST', 'redis'),
    os.environ.get('REDIS_PORT', '6379'),
    os.environ.get('REDIS_DB_CELERY', '0')
)

CELERY_RESULT_BACKEND = 'redis://{}:{}/{}'.format(
    os.environ.get('REDIS_HOST', 'redis'),
    os.environ.get('REDIS_PORT', '6379'),
    os.environ.get('REDIS_DB_RESULTS', '2')
)

# Task configuration
CELERY_TASK_SERIALIZER = 'json'
CELERY_RESULT_SERIALIZER = 'json'
CELERY_ACCEPT_CONTENT = ['json']
CELERY_TIMEZONE = 'UTC'

# Retry configuration
CELERY_TASK_RETRY_KWARGS = {'max_retries': 3}
CELERY_TASK_RETRY_DELAY = 60  # 1 minute

# Task time limits
CELERY_TASK_TIME_LIMIT = 30 * 60  # 30 minutes hard limit
CELERY_TASK_SOFT_TIME_LIMIT = 25 * 60  # 25 minutes soft limit

# Beat schedule adjustments
CELERY_BEAT_SCHEDULE = {
    'scrape-moneycontrol': {
        'task': 'apps.scrapers.tasks.scrape_moneycontrol',
        'schedule': crontab(minute=0, hour='*/4'),  # Every 4 hours
    },
    'analyze-articles': {
        'task': 'apps.scrapers.tasks.analyze_new_articles',
        'schedule': crontab(minute='*/15'),  # Every 15 minutes
    },
    'generate-predictions': {
        'task': 'apps.scrapers.tasks.generate_predictions',
        'schedule': crontab(minute=0, hour='0,6,12,18'),  # 4x daily
    },
}
```

---

## Environment Variables for Production

```bash
# .env (Production Configuration)

# Django
DJANGO_DEBUG=False
DJANGO_SECRET_KEY=your-super-secret-key-here-change-this
ENVIRONMENT=production
ALLOWED_HOSTS=yourdomain.com,www.yourdomain.com,api.yourdomain.com

# Database
DB_ENGINE=postgresql
DB_NAME=smart_news_db
DB_USER=postgres
DB_PASSWORD=strong-password-here
DB_HOST=postgres
DB_PORT=5432
DB_BACKUP_RETENTION_DAYS=30

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_DB_CACHE=1
REDIS_DB_CELERY=0
REDIS_DB_RESULTS=2

# NLP Service
NLP_SERVICE_URL=http://nlp-service:8001
NLP_MODEL_BART=facebook/bart-large-cnn
NLP_MODEL_FINBERT=ProsusAI/finbert
NLP_BATCH_SIZE=32
NLP_MAX_LENGTH=1024

# Prediction Service
PREDICTION_SERVICE_URL=http://prediction-service:8002
PREDICTION_CONFIDENCE_THRESHOLD=0.6
PREDICTION_EXPIRY_HOURS=24

# API Keys (Use AWS Secrets Manager or HashiCorp Vault)
YOUTUBE_API_KEY=your-youtube-api-key
TWITTER_BEARER_TOKEN=your-twitter-bearer-token
TELEGRAM_API_ID=your-telegram-api-id
TELEGRAM_API_HASH=your-telegram-api-hash

# Email Configuration
EMAIL_BACKEND=django.core.mail.backends.smtp.EmailBackend
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_HOST_USER=apikey
EMAIL_HOST_PASSWORD=SG.your-sendgrid-key
DEFAULT_FROM_EMAIL=noreply@yourdomain.com

# Security
SECURE_SSL_REDIRECT=True
SESSION_COOKIE_SECURE=True
CSRF_COOKIE_SECURE=True
SECURE_HSTS_SECONDS=31536000

# Sentry Error Tracking
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project

# AWS S3 (for media storage)
AWS_STORAGE_BUCKET_NAME=your-bucket-name
AWS_S3_REGION_NAME=us-east-1
AWS_S3_CUSTOM_DOMAIN=cdn.yourdomain.com
AWS_S3_USE_SSL=True

# Monitoring
DATADOG_API_KEY=your-datadog-api-key
LOG_LEVEL=INFO
```

---

## Gunicorn Configuration for Production

```python
# backend/django_core/gunicorn_config.py

import multiprocessing
import os

bind = "0.0.0.0:8000"
workers = multiprocessing.cpu_count() * 2 + 1
worker_class = "sync"
worker_connections = 1000
max_requests = 1000
max_requests_jitter = 50
timeout = 30
keep_alive = 5

# Logging
accesslog = "/var/log/gunicorn/access.log"
errorlog = "/var/log/gunicorn/error.log"
loglevel = "info"

# Process naming
proc_name = "smart-news-api"

# Server mechanics
daemon = False
pidfile = None
umask = 0
user = None
group = None

# SSL (if not using reverse proxy)
# keyfile = "/etc/ssl/private/key.pem"
# certfile = "/etc/ssl/certs/cert.pem"
```

### Docker Integration

```dockerfile
# backend/Dockerfile (production)

FROM python:3.11-slim as builder
WORKDIR /build
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

FROM python:3.11-slim
WORKDIR /app
COPY --from=builder /root/.local /root/.local
ENV PATH=/root/.local/bin:$PATH
COPY . .

# Security: Create non-root user
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8000
CMD ["gunicorn", "--config", "gunicorn_config.py", "config.wsgi:application"]
```

---

## Nginx Security Headers

```nginx
# nginx/conf.d/security.conf

# Security Headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;

# Content Security Policy
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self' api.yourdomain.com; frame-ancestors 'none';" always;

# Rate Limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=login_limit:10m rate=5r/m;

server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;
    
    # SSL Certificates
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Gzip compression
    gzip on;
    gzip_types text/plain text/css text/xml text/javascript 
               application/x-javascript application/json application/xml+rss;
    gzip_min_length 1000;
    
    # Rate limiting
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;
        proxy_pass http://backend:8000;
    }
    
    location /api/token/ {
        limit_req zone=login_limit burst=5 nodelay;
        proxy_pass http://backend:8000;
    }
}
```

---

## Database Maintenance

```bash
#!/bin/bash
# scripts/maintenance.sh

# Daily backup
backup_database() {
    BACKUP_DIR="/backups/database"
    TIMESTAMP=$(date +%Y-%m-%d-%H-%M-%S)
    
    mkdir -p $BACKUP_DIR
    
    docker-compose exec -T postgres pg_dump \
        -U postgres smart_news_db | \
        gzip > $BACKUP_DIR/backup_$TIMESTAMP.sql.gz
    
    # Upload to S3
    aws s3 cp $BACKUP_DIR/backup_$TIMESTAMP.sql.gz \
        s3://your-backup-bucket/
    
    # Cleanup old backups (30 days)
    find $BACKUP_DIR -mtime +30 -delete
}

# Weekly optimization
optimize_database() {
    docker-compose exec postgres psql -U postgres smart_news_db << EOF
VACUUM ANALYZE;
REINDEX DATABASE smart_news_db;
EOF
}

# Run maintenance
backup_database
optimize_database
```

---

## Monitoring Checklist

- [ ] Setup CloudWatch/Datadog dashboards
- [ ] Configure alerts for CPU > 80%
- [ ] Configure alerts for Memory > 90%
- [ ] Configure alerts for Disk space < 10GB
- [ ] Monitor API response times (target: < 200ms)
- [ ] Monitor failed API requests (target: < 1%)
- [ ] Setup uptime monitoring (Pingdom/UptimeRobot)
- [ ] Monitor Celery task failures
- [ ] Monitor database connection pool usage
- [ ] Setup error tracking (Sentry)

---

## Scaling Recommendations

### Horizontal Scaling

```bash
# Scale Celery workers
docker-compose up -d --scale celery-worker=3

# Load balance with Nginx
upstream backend_cluster {
    server backend-1:8000;
    server backend-2:8000;
    server backend-3:8000;
}
```

### Vertical Scaling

- Increase CPU/memory allocation per service
- Increase Gunicorn workers
- Increase database connection pool
- Increase Redis max memory

---

**Last Updated**: 2024-04-02  
**Recommended Review**: Quarterly
