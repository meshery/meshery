FROM python:3.11-slim

WORKDIR /app

# Install system dependencies including Chromium for web scraping
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create necessary directories
RUN mkdir -p /app/logs /app/data

# Set up scrapy settings
ENV SCRAPY_SETTINGS_MODULE=scrapy_project.settings

# Create logs directory
RUN mkdir -p /app/logs

# Default command to run scrapy
CMD ["python", "-m", "scrapy", "crawl", "news_pipeline"]
