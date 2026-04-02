"""
FastAPI NLP Service for text summarization and sentiment analysis
"""

from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import List, Optional
import logging
import torch
from transformers import pipeline, BartForConditionalGeneration, BartTokenizer
from decouple import config
import time

# Setup
logger = logging.getLogger(__name__)
app = FastAPI(title="Smart News NLP Service", version="1.0.0")

# Configuration
SUMMARIZER_MODEL = config('SUMMARIZER_MODEL', 'facebook/bart-large-cnn')
SENTIMENT_MODEL = config('SENTIMENT_MODEL', 'ProsusAI/finbert')
GPU_ENABLED = config('NLP_GPU_ENABLED', False, cast=bool)
BATCH_SIZE = config('NLP_BATCH_SIZE', 32, cast=int)

# Device
device = 'cuda' if (GPU_ENABLED and torch.cuda.is_available()) else 'cpu'
logger.info(f'Using device: {device}')

# Load models on startup
logger.info('Loading NLP models...')
try:
    summarizer = pipeline('summarization', model=SUMMARIZER_MODEL, device=0 if device == 'cuda' else -1)
    logger.info(f'Loaded summarizer: {SUMMARIZER_MODEL}')
except Exception as e:
    logger.error(f'Failed to load summarizer: {e}')
    summarizer = None

try:
    sentiment_analyzer = pipeline('sentiment-analysis', model=SENTIMENT_MODEL, device=0 if device == 'cuda' else -1)
    logger.info(f'Loaded sentiment model: {SENTIMENT_MODEL}')
except Exception as e:
    logger.error(f'Failed to load sentiment model: {e}')
    sentiment_analyzer = None


# Pydantic Models
class TextInput(BaseModel):
    text: str = Field(..., min_length=10, max_length=50000)
    article_id: Optional[str] = None


class SummarizationRequest(BaseModel):
    text: str = Field(..., min_length=100)
    max_length: int = 150
    min_length: int = 50


class SentimentRequest(BaseModel):
    text: str = Field(..., min_length=10)


class BatchAnalysisRequest(BaseModel):
    texts: List[TextInput]
    tasks: List[str] = Field(default=['summarize', 'sentiment'])


class SummarizationResponse(BaseModel):
    summary: str
    original_length: int
    summary_length: int
    compression_ratio: float
    processing_time: float


class SentimentResponse(BaseModel):
    label: str
    score: float
    article_id: Optional[str] = None


class BatchAnalysisResponse(BaseModel):
    results: List[dict]
    total_processing_time: float
    success_count: int
    error_count: int


# Health Check
@app.get('/health')
async def health():
    """Health check endpoint"""
    return {
        'status': 'healthy',
        'summarizer_loaded': summarizer is not None,
        'sentiment_analyzer_loaded': sentiment_analyzer is not None,
        'device': device,
    }


# Summarization Endpoints
@app.post('/api/summarize', response_model=SummarizationResponse)
async def summarize(request: SummarizationRequest):
    """Summarize text using BART"""
    if not summarizer:
        raise HTTPException(status_code=503, detail='Summarizer not loaded')
    
    start_time = time.time()
    
    try:
        # BART requires minimum length
        if len(request.text.split()) < 50:
            raise HTTPException(
                status_code=400,
                detail='Text too short for summarization (min 50 words)'
            )
        
        # Generate summary
        summary_result = summarizer(
            request.text,
            max_length=request.max_length,
            min_length=request.min_length,
            do_sample=False
        )
        
        summary_text = summary_result[0]['summary_text']
        processing_time = time.time() - start_time
        
        return SummarizationResponse(
            summary=summary_text,
            original_length=len(request.text),
            summary_length=len(summary_text),
            compression_ratio=len(summary_text) / len(request.text),
            processing_time=processing_time,
        )
    
    except Exception as e:
        logger.error(f'Summarization error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


# Sentiment Analysis Endpoints
@app.post('/api/sentiment', response_model=SentimentResponse)
async def analyze_sentiment(request: SentimentRequest):
    """Analyze sentiment of text using FinBERT"""
    if not sentiment_analyzer:
        raise HTTPException(status_code=503, detail='Sentiment analyzer not loaded')
    
    try:
        result = sentiment_analyzer(request.text[:512])  # Limit to 512 tokens
        
        return SentimentResponse(
            label=result[0]['label'],
            score=result[0]['score'],
        )
    
    except Exception as e:
        logger.error(f'Sentiment analysis error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


@app.post('/api/sentiment-batch', response_model=List[SentimentResponse])
async def analyze_sentiment_batch(texts: List[SentimentRequest]):
    """Analyze sentiment of multiple texts"""
    if not sentiment_analyzer:
        raise HTTPException(status_code=503, detail='Sentiment analyzer not loaded')
    
    try:
        results = []
        for request in texts:
            result = sentiment_analyzer(request.text[:512])
            results.append(SentimentResponse(
                label=result[0]['label'],
                score=result[0]['score'],
            ))
        
        return results
    
    except Exception as e:
        logger.error(f'Batch sentiment analysis error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


# Combined Analysis
@app.post('/api/analyze', response_model=dict)
async def analyze_article(
    text: str = Field(..., min_length=10),
    article_id: Optional[str] = None,
    tasks: List[str] = ['sentiment']
):
    """Analyze article with multiple NLP tasks"""
    results = {
        'article_id': article_id,
        'tasks': {}
    }
    
    try:
        # Sentiment analysis
        if 'sentiment' in tasks and sentiment_analyzer:
            sentiment_result = sentiment_analyzer(text[:512])
            results['tasks']['sentiment'] = {
                'label': sentiment_result[0]['label'],
                'score': float(sentiment_result[0]['score']),
            }
        
        # Summarization (if text is long enough)
        if 'summarize' in tasks and summarizer and len(text.split()) >= 50:
            summary_result = summarizer(text, max_length=150, min_length=50, do_sample=False)
            results['tasks']['summarize'] = {
                'summary': summary_result[0]['summary_text'],
                'compression_ratio': len(summary_result[0]['summary_text']) / len(text),
            }
        
        return results
    
    except Exception as e:
        logger.error(f'Analysis error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


# Batch Processing
@app.post('/api/batch-analyze', response_model=BatchAnalysisResponse)
async def batch_analyze(request: BatchAnalysisRequest, background_tasks: BackgroundTasks):
    """Batch analyze multiple texts"""
    start_time = time.time()
    success_count = 0
    error_count = 0
    results = []
    
    try:
        for item in request.texts:
            try:
                result = await analyze_article(
                    text=item.text,
                    article_id=item.article_id,
                    tasks=request.tasks
                )
                results.append(result)
                success_count += 1
            except Exception as e:
                logger.error(f'Error analyzing text {item.article_id}: {e}')
                error_count += 1
                results.append({
                    'article_id': item.article_id,
                    'error': str(e)
                })
        
        processing_time = time.time() - start_time
        
        return BatchAnalysisResponse(
            results=results,
            total_processing_time=processing_time,
            success_count=success_count,
            error_count=error_count,
        )
    
    except Exception as e:
        logger.error(f'Batch analysis error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


# Root endpoint
@app.get('/')
async def root():
    """Root endpoint"""
    return {
        'service': 'Smart News NLP Service',
        'version': '1.0.0',
        'endpoints': [
            '/health',
            '/api/summarize',
            '/api/sentiment',
            '/api/sentiment-batch',
            '/api/analyze',
            '/api/batch-analyze',
        ]
    }


if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='0.0.0.0', port=8001, workers=1)
