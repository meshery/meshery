#!/bin/bash

# ============================================================
# Smart News Analyzer - Build and Start Script
# ============================================================

set -e

echo "🚀 Starting Smart News Analyzer..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from template..."
    cp .env.example .env
    echo "✅ .env created. Please update it with your credentials."
fi

# Build and start services
echo "🐳 Building Docker images..."
docker-compose build

echo "📦 Starting services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Run migrations
echo "🗄️  Running database migrations..."
docker-compose exec -T backend python manage.py migrate

# Create superuser if it doesn't exist
echo "👤 Creating superuser..."
docker-compose exec -T backend python manage.py shell << END
from django.contrib.auth.models import User
if not User.objects.filter(username='admin').exists():
    User.objects.create_superuser('admin', 'admin@example.com', 'admin')
    print("✅ Superuser created")
else:
    print("⚠️  Superuser already exists")
END

# Load initial data
echo "📊 Loading initial data..."
docker-compose exec -T backend python manage.py loaddata --format json /app/fixtures/*.json 2>/dev/null || true

echo ""
echo "✅ Smart News Analyzer is running!"
echo ""
echo "📍 Access points:"
echo "   - Frontend:      http://localhost:3000"
echo "   - API:           http://localhost:8000/api"
echo "   - Admin Panel:   http://localhost:8000/admin"
echo "   - Nginx Status:  http://localhost/health"
echo ""
echo "💡 Useful commands:"
echo "   - View logs:     docker-compose logs -f"
echo "   - Shell:         docker-compose exec backend python manage.py shell"
echo "   - Migrations:    docker-compose exec backend python manage.py migrate"
echo "   - Celery status: docker-compose exec celery-worker celery -A config inspect active"
echo ""
