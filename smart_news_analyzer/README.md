# 🚀 SMART NEWS SUMMARIZER + MARKET ADVISOR

## Enterprise-Grade AI-Powered Financial News Intelligence Platform

A complete production-ready microservices architecture for:
- **Scraping** financial news from multiple sources (websites, YouTube, Telegram, Twitter/X)
- **Analyzing** sentiment and summarizing news using advanced NLP models
- **Predicting** stock market movements with AI-driven insights
- **Visualizing** market intelligence through a modern web dashboard

---

## 📋 Table of Contents

- [System Architecture](#-system-architecture)
- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [Technologies](#-technologies)
- [Configuration](#️-configuration)
- [API Documentation](#-api-documentation)
- [Deployment](#-deployment)
- [Development](#-development)
- [Contributing](#-contributing)

---

## 🏗️ System Architecture

### Microservices Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     NGINX Reverse Proxy                      │
└──────────┬──────────────────────────────────────────────────┘
           │
    ┌──────┴──────┐
    │             │
┌───▼──────┐  ┌──▼──────────┐
│ Frontend │  │ Django Core  │
│(Next.js) │  │ (DRF API)    │
└──────────┘  └─────┬────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
   ┌────▼────┐ ┌───▼────┐ ┌──▼──────┐
   │ Scrapers │ │NLP Svc │ │Pred Svc │
   └────┬────┘ └───┬────┘ └──┬──────┘
        │          │         │
        └──────┬───┴────┬────┘
               │        │
          ┌────▼─┐  ┌──▼──────┐
          │Redis │  │PostgreSQL│
          └──────┘  └──────────┘

      ┌──────────────────────┐
      │  Celery Workers      │
      │  (Async Tasks)       │
      └──────────────────────┘
```

### Services

| Service | Port | Technology | Purpose |
|---------|------|-----------|---------|
| **Frontend** | 3000 | Next.js 14, React, TypeScript | Web Dashboard UI |
| **Django Core** | 8000 | Django, DRF, PostgreSQL | API Gateway & Business Logic |
| **NLP Service** | 8001 | FastAPI, Transformers, BART | Text Summarization & Sentiment Analysis |
| **Prediction Service** | 8002 | FastAPI, Scikit-learn | Stock Market Prediction Engine |
| **Scrapers** | N/A | Scrapy, YouTube API, Telethon | Data Collection |
| **Celery Worker** | N/A | Python Celery | Async Task Processing |
| **Celery Beat** | N/A | Django Celery Beat | Scheduled Tasks |
| **PostgreSQL** | 5432 | PostgreSQL 15 | Primary Database |
| **Redis** | 6379 | Redis 7 | Cache & Message Broker |
| **Nginx** | 80/443 | Nginx Alpine | Reverse Proxy & Load Balancing |

---

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose (v2.0+)
- Git
- Python 3.11+ (for local development)
- Node.js 18+ (for local frontend development)

### Installation

1. **Clone Repository**
```bash
git clone https://github.com/yourusername/smart-news-analyzer.git
cd smart-news-analyzer
```

2. **Create Environment Configuration**
```bash
cp .env.example .env
# Edit .env with your API keys and credentials
nano .env
```

3. **Start Services**
```bash
bash scripts/start.sh
```

This will:
- Build all Docker images
- Start all containers
- Run database migrations
- Create default admin user
- Load initial fixtures

4. **Access Applications**

| Application | URL |
|-------------|-----|
| Frontend Dashboard | http://localhost:3000 |
| API Endpoints | http://localhost:8000/api |
| Admin Panel | http://localhost:8000/admin |
| API Docs | http://localhost:8000/api/docs |

### Default Credentials

- **Admin User**: `admin` / `admin`
- **Redis**: No password (default)
- **PostgreSQL**: `postgres` / `postgres-password-change-me`

⚠️ **Change all default credentials in production!**

---

## 📁 Project Structure

```
smart-news-analyzer/
├── backend/
│   ├── django_core/                 # Main Django Application
│   │   ├── config/                  # Django Settings
│   │   ├── apps/
│   │   │   ├── articles/            # Article CRUD
│   │   │   ├── tasks/               # Task Management
│   │   │   ├── users/               # User Management
│   │   │   ├── scrapers/            # Scraper Management
│   │   │   └── predictions/         # Predictions API
│   │   ├── fixtures/                # Initial Data
│   │   ├── static/                  # Static Files
│   │   ├── media/                   # User Uploads
│   │   ├── manage.py
│   │   └── requirements.txt
│   │
│   ├── scrapers/                    # Scrapy Project
│   │   ├── scrapy_project/
│   │   │   ├── spiders/             # Web Scrapers
│   │   │   ├── pipelines/           # Data Processing
│   │   │   ├── items.py
│   │   │   └── settings.py
│   │   ├── youtube_scraper.py       # YouTube API Handler
│   │   ├── telegram_scraper.py      # Telegram API Handler
│   │   ├── twitter_scraper.py       # Twitter/X Handler
│   │   └── requirements.txt
│   │
│   ├── nlp_service/                 # FastAPI NLP Service
│   │   ├── main.py
│   │   ├── models/                  # ML Models
│   │   ├── services/
│   │   │   ├── summarizer.py
│   │   │   └── sentiment.py
│   │   ├── schemas/
│   │   └── requirements.txt
│   │
│   └── prediction_service/          # FastAPI Prediction Service
│       ├── main.py
│       ├── models/
│       ├── services/
│       └── requirements.txt
│
├── frontend/                        # Next.js Application
│   ├── app/                         # App Router
│   ├── components/                  # React Components
│   ├── pages/
│   │   ├── dashboard/
│   │   ├── news-feed/
│   │   ├── predictions/
│   │   ├── articles/
│   │   └── admin/
│   ├── public/                      # Static Assets
│   ├── styles/                      # TailwindCSS
│   ├── utils/
│   ├── hooks/                       # Custom React Hooks
│   ├── package.json
│   └── next.config.js
│
├── docker/                          # Docker Configuration
│   └── [Additional configs]
│
├── nginx/                           # Nginx Configuration
│   ├── nginx.conf
│   └── conf.d/
│
├── scripts/                         # Utility Scripts
│   ├── start.sh
│   ├── stop.sh
│   ├── init_db.sql
│   └── backup.sh
│
├── docker-compose.yml               # Service Orchestration
├── .env.example                     # Environment Template
├── README.md                        # This File
└── .gitignore
```

---

## 🛠️ Technologies

### Backend
- **Python 3.11**
- **Django 4.2** - Web Framework
- **Django REST Framework** - API Development
- **FastAPI** - High-performance API framework
- **Celery** - Asynchronous Task Queue
- **PostgreSQL 15** - Database
- **Redis 7** - Caching & Message Broker

### ML/AI
- **Transformers (Hugging Face)** - NLP Models
- **BART** - Abstractive Summarization
- **FinBERT** - Financial Sentiment Analysis
- **spaCy** - Named Entity Recognition
- **Scikit-learn** - ML Utilities

### Data Collection
- **Scrapy** - Web Scraping Framework
- **Telethon** - Telegram API Client
- **snscrape** - Twitter/X Scraper
- **YouTube Data API v3** - YouTube Integration

### Frontend
- **Next.js 14** - React Framework
- **TypeScript** - Type Safety
- **TailwindCSS** - Styling
- **ShadCN UI** - Component Library
- **Chart.js** - Data Visualization
- **Axios** - HTTP Client

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Orchestration
- **Nginx** - Reverse Proxy
- **Gunicorn** - WSGI Server

---

## ⚙️ Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Django
DJANGO_SECRET_KEY=your-secret-key
DJANGO_DEBUG=False
DJANGO_ALLOWED_HOSTS=localhost,api.example.com

# Database
DB_ENGINE=django.db.backends.postgresql
DB_NAME=smart_news_db
DB_USER=postgres
DB_PASSWORD=secure-password
DB_HOST=postgres
DB_PORT=5432

# Redis
REDIS_URL=redis://redis:6379/0

# API Keys
YOUTUBE_API_KEY=your-youtube-api-key
TWITTER_BEARER_TOKEN=your-twitter-token
TELEGRAM_API_ID=your-telegram-api-id
TELEGRAM_API_HASH=your-telegram-api-hash

# External Services
NLP_SERVICE_URL=http://nlp-service:8001
PREDICTION_SERVICE_URL=http://prediction-service:8002

# Models
SUMMARIZER_MODEL=facebook/bart-large-cnn
SENTIMENT_MODEL=ProsusAI/finbert
```

See `.env.example` for complete configuration options.

---

## 📚 API Documentation

### API Base URLs

- **REST API**: `http://localhost:8000/api/`
- **API Documentation**: `http://localhost:8000/api/docs/`
- **Admin Dashboard**: `http://localhost:8000/admin/`

### Core Endpoints

#### Articles

```
GET    /api/articles/                 # List articles
POST   /api/articles/                 # Create article
GET    /api/articles/{id}/            # Get article details
PUT    /api/articles/{id}/            # Update article
DELETE /api/articles/{id}/            # Delete article

GET    /api/articles/search/          # Search articles
GET    /api/articles/by-sentiment/    # Filter by sentiment
GET    /api/articles/trending/        # Get trending articles
```

#### Predictions

```
GET    /api/predictions/              # List predictions
GET    /api/predictions/{symbol}/     # Get stock prediction
POST   /api/predictions/generate/     # Generate new prediction
GET    /api/predictions/history/      # Prediction history
```

#### Tasks

```
GET    /api/tasks/                    # List tasks
GET    /api/tasks/{id}/               # Task status
POST   /api/tasks/trigger-scraping/   # Start scraping job
GET    /api/tasks/celery-status/      # Celery worker status
```

#### User Management

```
POST   /api/auth/register/            # Register user
POST   /api/auth/login/               # Login
POST   /api/auth/logout/              # Logout
GET    /api/users/profile/            # Get user profile
```

---

## 🚀 Deployment

### Production Deployment

#### 1. **Prepare Server**

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo bash get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. **Clone Repository**

```bash
git clone your-repo.git /opt/smart-news-analyzer
cd /opt/smart-news-analyzer
```

#### 3. **Configure Environment**

```bash
cp .env.example .env
nano .env  # Edit with production values
```

**Critical Production Settings:**

```bash
DJANGO_DEBUG=False
DJANGO_SECRET_KEY=<strong-random-key>
DJANGO_ALLOWED_HOSTS=yourdomain.com,www.yourdomain.com
DB_PASSWORD=<strong-database-password>
```

#### 4. **Setup SSL/HTTPS**

```bash
# Using Let's Encrypt with Certbot
sudo apt install certbot python3-certbot-nginx
sudo certbot certonly --nginx -d yourdomain.com
```

Enable in nginx.conf:
```nginx
listen 443 ssl http2;
ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
```

#### 5. **Start Services**

```bash
bash scripts/start.sh

# Or manually:
docker-compose -f docker-compose.yml up -d
```

#### 6. **Monitoring & Logs**

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f backend

# Check service status
docker-compose ps
```

### Scaling Considerations

- **Multiple Celery Workers**: `docker-compose up -d --scale celery-worker=4`
- **Load Balancing**: Nginx upstream configuration
- **Database Replication**: PostgreSQL replication for HA
- **Redis Cluster**: For distributed caching

---

## 🔧 Development

### Local Development Setup

#### Backend Development

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r backend/django_core/requirements.txt

# Run migrations
cd backend/django_core
python manage.py migrate

# Start development server
python manage.py runserver
```

#### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Access at http://localhost:3000
```

#### Scraper Testing

```bash
cd backend/scrapers

# Run specific spider
scrapy crawl moneycontrol_spider

# Test pipeline
scrapy crawl news_pipeline -a limit=5
```

### Testing

```bash
# Backend tests
cd backend/django_core
python manage.py test

# Run with coverage
coverage run --source='.' manage.py test
coverage report
```

### Database Management

```bash
# Create migrations
docker-compose exec backend python manage.py makemigrations

# Apply migrations
docker-compose exec backend python manage.py migrate

# Create superuser
docker-compose exec backend python manage.py createsuperuser

# Backup database
docker-compose exec postgres pg_dump -U postgres smart_news_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres smart_news_db < backup.sql
```

---

## 📊 Monitoring

### Health Checks

All services have built-in health checks:

```bash
curl http://localhost:8000/health/          # Backend
curl http://localhost:8001/health            # NLP Service
curl http://localhost:8002/health            # Prediction Service
curl http://localhost/health                 # Nginx
```

### Celery Monitoring

```bash
# View active tasks
docker-compose exec celery-worker celery -A config inspect active

# View registered tasks
docker-compose exec celery-worker celery -A config inspect registered

# Real-time monitoring
docker-compose exec celery-worker celery -A config events
```

### Database Monitoring

```bash
# Access PostgreSQL
docker-compose exec postgres psql -U postgres -d smart_news_db

# Common queries
SELECT COUNT(*) FROM articles;
SELECT COUNT(*) FROM predictions;
SELECT COUNT(*) FROM users;
```

---

## 🤝 Contributing

Contributions are welcome! Please follow our guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

---

## 🆘 Support

For issues, questions, or suggestions:
- Create an Issue on GitHub
- Check existing documentation
- Review API docs at http://localhost:8000/api/docs/

---

## 🎯 Roadmap

- [ ] Mobile app (React Native)
- [ ] Advanced ML prediction models
- [ ] Real-time WebSocket updates
- [ ] Multi-language support
- [ ] Advanced analytics dashboard
- [ ] Social sharing features
- [ ] Custom alerts & notifications
- [ ] Portfolio tracking
- [ ] Backtesting framework

---

**Built with ❤️ for financial intelligence**
