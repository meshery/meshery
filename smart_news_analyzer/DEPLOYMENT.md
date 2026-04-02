# 🚀 DEPLOYMENT GUIDE

## Smart News Summarizer + Market Advisor
### Production Deployment Instructions

---

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [AWS EC2 Deployment](#aws-ec2-deployment)
3. [DigitalOcean Deployment](#digitalocean-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Monitoring & Maintenance](#monitoring--maintenance)
6. [Troubleshooting](#troubleshooting)
7. [Security Checklist](#security-checklist)

---

## Prerequisites

### System Requirements
- **OS**: Ubuntu 22.04 LTS (recommended)
- **CPU**: Minimum 4 cores (8+ recommended)
- **RAM**: Minimum 8GB (16GB+ recommended)
- **Storage**: Minimum 50GB SSD
- **Bandwidth**: 10+ Mbps

### Required Software
- Docker (v24.0+)
- Docker Compose (v2.20+)
- Git
- Python 3.11+ (for local development)
- Node.js 18+ (for frontend development)

### Accounts & Credentials
- GitHub repository access
- AWS or DigitalOcean account (for hosting)
- API Keys:
  - YouTube API
  - Twitter/X Bearer Token
  - Telegram API credentials
  - Domain registrar account

---

## AWS EC2 Deployment

### 1. Launch EC2 Instance

```bash
# Create security group
aws ec2 create-security-group \
  --group-name smart-news-sg \
  --description "Security group for Smart News Analyzer"

# Add ingress rules
aws ec2 authorize-security-group-ingress \
  --group-name smart-news-sg \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
  --group-name smart-news-sg \
  --protocol tcp \
  --port 443 \
  --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
  --group-name smart-news-sg \
  --protocol tcp \
  --port 22 \
  --cidr YOUR_IP/32
```

### 2. Connect & Update System

```bash
# SSH into instance
ssh -i your-key.pem ubuntu@your-ec2-public-ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo bash get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

### 3. Clone Repository

```bash
cd ~
git clone -b main https://github.com/yourusername/smart-news-analyzer.git
cd smart-news-analyzer
```

### 4. Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit with production values
nano .env

# Key configurations for production:
DJANGO_DEBUG=False
DJANGO_SECRET_KEY=$(python -c 'from django.core.management.utils import get_random_secret_key; print(get_random_secret_key())')
DJANGO_ALLOWED_HOSTS=yourdomain.com,www.yourdomain.com
DB_PASSWORD=$(openssl rand -base64 32)
ENVIRONMENT=production
```

### 5. Deploy with Docker Compose

```bash
# Build and start services
docker-compose build
docker-compose up -d

# Verify services
docker-compose ps

# Check logs
docker-compose logs -f backend
```

### 6. Setup SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain certificate
sudo certbot certonly --standalone \
  -d yourdomain.com \
  -d www.yourdomain.com

# Update nginx.conf with certificate paths
# ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
# ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

# Restart nginx
docker-compose restart nginx

# Auto-renewal (crontab)
sudo certbot renew --quiet
```

---

## DigitalOcean Deployment

### 1. Create Droplet

```bash
# Using doctl CLI
doctl compute droplet create smart-news-analyzer \
  --region nyc3 \
  --image ubuntu-22-04-x64 \
  --size s-2vcpu-4gb \
  --enable-backups \
  --enable-ipv6
```

### 2. Initial Setup

```bash
# SSH into droplet
ssh root@YOUR_DROPLET_IP

# Create non-root user
useradd -m -s /bin/bash appuser
usermod -aG sudo appuser
su - appuser

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo bash get-docker.sh
sudo usermod -aG docker $USER
```

### 3. Deploy Application

```bash
# Clone repo
git clone https://github.com/yourusername/smart-news-analyzer.git
cd smart-news-analyzer

# Configure environment
cp .env.example .env
nano .env

# Start services
docker-compose up -d

# Setup Firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

---

## Kubernetes Deployment

### 1. Create Kubernetes Manifests

```yaml
# kubernetes/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: smart-news
```

```yaml
# kubernetes/postgres-deployment.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: smart-news
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: db-config
              key: db-name
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        ports:
        - containerPort: 5432
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 20Gi
```

### 2. Deploy to Kubernetes

```bash
# Create namespace
kubectl apply -f kubernetes/namespace.yaml

# Create secrets
kubectl create secret generic db-secret \
  --from-literal=username=postgres \
  --from-literal=password=$(openssl rand -base64 32) \
  -n smart-news

# Create ConfigMap
kubectl create configmap db-config \
  --from-literal=db-name=smart_news_db \
  -n smart-news

# Deploy PostgreSQL
kubectl apply -f kubernetes/postgres-deployment.yaml

# Deploy application services
kubectl apply -f kubernetes/
```

---

## Monitoring & Maintenance

### 1. Enable Health Checks

```bash
# Monitor container health
docker-compose ps

# Check service health
curl http://localhost:8000/health/
curl http://localhost:8001/health
curl http://localhost:8002/health
curl http://localhost/health
```

### 2. Setup Monitoring with Prometheus

```bash
# Create prometheus config
cat > monitoring/prometheus.yml << EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'django'
    static_configs:
      - targets: ['localhost:8000']

  - job_name: 'postgres'
    static_configs:
      - targets: ['localhost:5432']
EOF

# Run Prometheus
docker run -d \
  -p 9090:9090 \
  -v $(pwd)/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

### 3. Log Aggregation

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f backend
docker-compose logs -f nlp-service

# Save logs to file
docker-compose logs > deployment.log

# Rotate logs
docker-compose logs --tail 100 > logs/archive-$(date +%Y-%m-%d).log
```

### 4. Backup Strategy

```bash
#!/bin/bash
# backup.sh - Daily backup script

BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)

# Backup PostgreSQL
docker-compose exec -T postgres pg_dump \
  -U postgres smart_news_db > \
  $BACKUP_DIR/db_backup_$TIMESTAMP.sql

# Backup media files
tar -czf $BACKUP_DIR/media_backup_$TIMESTAMP.tar.gz \
  backend/django_core/media/

# Upload to S3
aws s3 sync $BACKUP_DIR s3://your-backup-bucket/

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -mtime +30 -delete
```

### 5. Update Application

```bash
# Pull latest code
git pull origin main

# Build new images
docker-compose build

# Apply migrations
docker-compose exec backend python manage.py migrate

# Restart services
docker-compose up -d --no-deps backend celery-worker celery-beat
```

---

## Troubleshooting

### Services Won't Start

```bash
# Check logs
docker-compose logs backend

# Check dependencies
docker-compose ps

# Rebuild images
docker-compose build --no-cache

# Full restart
docker-compose down -v
docker-compose up -d
```

### Database Connection Issues

```bash
# Check PostgreSQL
docker-compose exec postgres psql -U postgres -c "\l"

# Test connection
docker-compose exec postgres pg_isready

# Check migrations
docker-compose exec backend python manage.py showmigrations
```

### Memory Issues

```bash
# Check resource usage
docker stats

# Limit memory
docker-compose down
# Edit docker-compose.yml to add memory limits

# Restart
docker-compose up -d
```

### NLP Service Slow

```bash
# Check GPU availability
docker-compose exec nlp-service python -c "import torch; print(torch.cuda.is_available())"

# Reduce batch size in .env
NLP_BATCH_SIZE=16

# Restart service
docker-compose restart nlp-service
```

---

## Security Checklist

### ✅ Before Production

- [ ] Change all default passwords
- [ ] Update Django SECRET_KEY
- [ ] Enable HTTPS/SSL
- [ ] Setup firewall rules
- [ ] Enable database backups
- [ ] Configure email notifications
- [ ] Setup monitoring and alerts
- [ ] Implement rate limiting
- [ ] Enable CORS only for trusted origins
- [ ] Rotate credentials regularly

### ✅ Database Security

```bash
# Create database backup user
docker-compose exec postgres psql -U postgres << EOF
CREATE USER backup WITH ENCRYPTED PASSWORD 'strong-password';
GRANT CONNECT ON DATABASE smart_news_db TO backup;
GRANT USAGE ON SCHEMA public TO backup;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup;
EOF

# Enable SSL for database connections
# Edit postgresql.conf and restart
```

### ✅ Application Security

```bash
# Enable security headers in nginx.conf
add_header X-Frame-Options "SAMEORIGIN";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
add_header Referrer-Policy "strict-origin-when-cross-origin";
add_header Content-Security-Policy "default-src 'self'";

# Restart nginx
docker-compose restart nginx
```

---

## Performance Optimization

### Database Optimization

```sql
-- Create indexes
CREATE INDEX idx_articles_published ON articles(published_date);
CREATE INDEX idx_articles_sentiment ON articles(sentiment);
CREATE INDEX idx_predictions_symbol ON predictions(stock_symbol);

-- Run VACUUM
VACUUM ANALYZE;
```

### Cache Optimization

```bash
# Monitor Redis
docker-compose exec redis redis-cli INFO stats

# Clear old caches
docker-compose exec redis redis-cli FLUSHDB

# Set max memory policy
docker-compose exec redis redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

### Load Balancing

```nginx
# nginx/conf.d/load-balance.conf
upstream backend_pool {
    server backend:8000 weight=1;
    server backend:8000 weight=1;
}

server {
    listen 80;
    location /api/ {
        proxy_pass http://backend_pool;
    }
}
```

---

## Recovery Procedures

### Restore from Backup

```bash
# Download backup
aws s3 cp s3://your-backup-bucket/db_backup_YYYY-MM-DD.sql ./

# Drop existing database
docker-compose exec postgres dropdb -U postgres smart_news_db

# Create new database
docker-compose exec postgres createdb -U postgres smart_news_db

# Restore backup
docker-compose exec -T postgres psql -U postgres smart_news_db < db_backup_YYYY-MM-DD.sql

# Restart services
docker-compose restart
```

---

## Support & Documentation

- 📚 [Main README](../README.md)
- 🐛 [GitHub Issues](https://github.com/yourusername/smart-news-analyzer/issues)
- 💬 [Discussions](https://github.com/yourusername/smart-news-analyzer/discussions)

---

**Last Updated**: 2024-04-02
**Version**: 1.0.0
