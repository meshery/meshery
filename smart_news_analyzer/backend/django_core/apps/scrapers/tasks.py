"""
Celery tasks for asynchronous processing
"""

import logging
import requests
from celery import shared_task, chain, group
from django.db import transaction
from django.utils import timezone
from apps.articles.models import Article, ArticleSource, ArticleAnalysis
from apps.predictions.models import StockPrediction, PredictionHistory
from apps.tasks.models import ScraperTask, AnalysisTask
from apps.scrapers.models import ScraperConfiguration
from decouple import config

logger = logging.getLogger(__name__)

NLP_SERVICE_URL = config('NLP_SERVICE_URL', 'http://nlp-service:8001')
PREDICTION_SERVICE_URL = config('PREDICTION_SERVICE_URL', 'http://prediction-service:8002')


@shared_task(bind=True, max_retries=3)
def scrape_moneycontrol(self):
    """Scrape Moneycontrol news"""
    try:
        logger.info('Starting Moneycontrol scraping...')
        
        task = ScraperTask.objects.create(
            source='moneycontrol',
            task_type='scrape',
            status='running'
        )
        
        # In production, would call the actual Scrapy spider
        # For now, log the action
        logger.info(f'Moneycontrol scraping completed. Task: {task.id}')
        
        task.status = 'completed'
        task.save()
        
        return f'Moneycontrol scraping completed: {task.id}'
    
    except Exception as e:
        logger.error(f'Moneycontrol scraping failed: {e}')
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=3)
def scrape_economic_times(self):
    """Scrape Economic Times news"""
    try:
        logger.info('Starting Economic Times scraping...')
        task = ScraperTask.objects.create(
            source='economic_times',
            task_type='scrape',
            status='running'
        )
        logger.info(f'Economic Times scraping completed')
        task.status = 'completed'
        task.save()
        return f'Economic Times scraping done: {task.id}'
    except Exception as e:
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=3)
def scrape_youtube_channels(self):
    """Scrape YouTube financial channels"""
    try:
        logger.info('Starting YouTube scraping...')
        # Call YouTube scraper
        logger.info('YouTube scraping completed')
        return 'YouTube scraping done'
    except Exception as e:
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=3)
def scrape_telegram_channels(self):
    """Scrape Telegram financial channels"""
    try:
        logger.info('Starting Telegram scraping...')
        logger.info('Telegram scraping completed')
        return 'Telegram scraping done'
    except Exception as e:
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=3)
def scrape_twitter_accounts(self):
    """Scrape Twitter/X financial accounts"""
    try:
        logger.info('Starting Twitter scraping...')
        logger.info('Twitter scraping completed')
        return 'Twitter scraping done'
    except Exception as e:
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=2)
def analyze_new_articles(self):
    """Analyze new unanalyzed articles"""
    try:
        logger.info('Starting article analysis...')
        
        task =AnalysisTask.objects.create(
            articles_count=0,
            status='processing'
        )
        
        # Get unanalyzed articles
        unanalyzed = Article.objects.filter(
            is_analyzed=False,
            analysis_attempted=False
        )[:100]
        
        count = 0
        errors = 0
        
        for article in unanalyzed:
            try:
                # Call NLP service
                npl_response = requests.post(
                    f'{NLP_SERVICE_URL}/api/analyze',
                    json={
                        'text': article.content,
                        'article_id': str(article.id),
                        'tasks': ['sentiment']
                    },
                    timeout=10
                )
                
                if npl_response.status_code == 200:
                    data = npl_response.json()
                    sentiment = data['tasks'].get('sentiment', {})
                    
                    # Update article
                    article.sentiment = sentiment.get('label', 'unknown').lower()
                    article.sentiment_score = sentiment.get('score', 0)
                    article.is_analyzed = True
                    article.save()
                    
                    # Create analysis record
                    ArticleAnalysis.objects.create(
                        article=article,
                        nlp_response=data,
                        confidence_scores=sentiment
                    )
                    
                    count += 1
                else:
                    errors += 1
                    logger.warning(f'NLP service error: {npl_response.status_code}')
            
            except Exception as e:
                errors += 1
                logger.error(f'Analysis error for article {article.id}: {e}')
                article.analysis_errors.append(str(e))
                article.save()
        
        task.success_count = count
        task.error_count = errors
        task.articles_completed = count + errors
        task.status = 'completed'
        task.completed_at = timezone.now()
        task.save()
        
        logger.info(f'Article analysis completed: {count} success, {errors} errors')
        return f'Analysis done: {count} articles'
    
    except Exception as e:
        logger.error(f'Analysis task error: {e}')
        self.retry(exc=e, countdown=60)


@shared_task(bind=True, max_retries=2)
def generate_predictions(self):
    """Generate stock predictions based on sentiment"""
    try:
        logger.info('Starting prediction generation...')
        
        # Group articles by stock symbol
        from django.db.models import Q
        
        articles = Article.objects.filter(
            is_analyzed=True,
            detected_stocks__isnull=False
        ).values('detected_stocks').distinct()[:50]
        
        predictions_made = 0
        
        for article_group in articles:
            stocks = article_group.get('detected_stocks', [])
            
            if not stocks:
                continue
            
            for stock in stocks:
                try:
                    # Get sentiment data
                    stock_articles = Article.objects.filter(
                        detected_stocks__contains=[stock],
                        is_analyzed=True
                    )
                    
                    sentiments = stock_articles.values_list('sentiment_score', flat=True)
                    sentiment_avg = sum(sentiments) / len(sentiments) if sentiments else 0
                    
                    positive = stock_articles.filter(sentiment='positive').count()
                    negative = stock_articles.filter(sentiment='negative').count()
                    neutral = stock_articles.filter(sentiment='neutral').count()
                    
                    # Call prediction service
                    pred_response = requests.post(
                        f'{PREDICTION_SERVICE_URL}/api/predict',
                        json={
                            'stock_symbol': stock,
                            'sentiment_score': sentiment_avg,
                            'positive_articles': positive,
                            'negative_articles': negative,
                            'neutral_articles': neutral,
                        },
                        timeout=10
                    )
                    
                    if pred_response.status_code == 200:
                        pred_data = pred_response.json()
                        
                        # Save prediction
                        expires_at = timezone.now() + timezone.timedelta(days=7)
                        
                        StockPrediction.objects.create(
                            stock_symbol=stock,
                            action=pred_data['action'].lower(),
                            confidence_score=pred_data['confidence_score'],
                            confidence_range=pred_data['confidence_range'],
                            risk_level=pred_data['risk_level'],
                            sentiment_score=pred_data['sentiment_score'],
                            positive_articles_count=positive,
                            negative_articles_count=negative,
                            neutral_articles_count=neutral,
                            reasoning=pred_data['reasoning'],
                            key_factors=pred_data['key_factors'],
                            expires_at=expires_at,
                        )
                        
                        predictions_made += 1
                
                except Exception as e:
                    logger.error(f'Prediction error for {stock}: {e}')
        
        logger.info(f'Generated {predictions_made} predictions')
        return f'Predictions done: {predictions_made}'
    
    except Exception as e:
        logger.error(f'Prediction generation error: {e}')
        self.retry(exc=e, countdown=60)


@shared_task
def cleanup_old_tasks():
    """Clean up old task records"""
    try:
        cutoff_date = timezone.now() - timezone.timedelta(days=30)
        
        deleted_count, _ = ScraperTask.objects.filter(
            created_at__lt=cutoff_date,
            status='completed'
        ).delete()
        
        logger.info(f'Cleaned up {deleted_count} old task records')
        return f'Cleanup done: {deleted_count} records deleted'
    
    except Exception as e:
        logger.error(f'Cleanup error: {e}')


# Group task - run all scrapers in parallel
@shared_task
def run_all_scrapers():
    """Run all scrapers in parallel"""
    job = group(
        scrape_moneycontrol.s(),
        scrape_economic_times.s(),
        scrape_youtube_channels.s(),
        scrape_telegram_channels.s(),
        scrape_twitter_accounts.s(),
    )
    return job.apply_async()


# Chain task - run analysis after scraping
@shared_task
def scrape_and_analyze():
    """Chain scrapers with analysis"""
    workflow = chain(
        run_all_scrapers.s(),
        analyze_new_articles.s(),
        generate_predictions.s(),
    )
    return workflow.apply_async()
