"""
FastAPI Prediction Service for stock market predictions
"""

from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field
from typing import List, Optional
import logging
from datetime import datetime, timedelta
import numpy as np
from decouple import config

logger = logging.getLogger(__name__)
app = FastAPI(title="Stock Prediction Service", version="1.0.0")


# Models
class SentimentData(BaseModel):
    stock_symbol: str
    sentiment_scores: List[float]
    positive_count: int
    negative_count: int
    neutral_count: int
    average_sentiment: float


class PredictionRequest(BaseModel):
    stock_symbol: str
    sentiment_score: float = Field(-1, ge=-1, le=1)
    positive_articles: int = 0
    negative_articles: int = 0
    neutral_articles: int = 0
    keywords: Optional[List[str]] = None
    timeframe: str = '1week'


class PredictionResponse(BaseModel):
    stock_symbol: str
    action: str  # BUY, HOLD, SELL
    confidence_score: float
    confidence_range: str
    risk_level: str
    reasoning: str
    key_factors: List[str]
    sentiment_score: float
    timestamp: str


class MarketSentimentResponse(BaseModel):
    overall_sentiment: str
    bullish_count: int
    bearish_count: int
    neutral_count: int
    sentiment_score: float


# Prediction engine
class PredictionEngine:
    """ML-based prediction engine for stock movements"""
    
    def __init__(self):
        self.confidence_threshold = 0.6
        self.risk_multiplier = 1.0
    
    def predict(self, request: PredictionRequest) -> PredictionResponse:
        """Generate prediction based on sentiment and article data"""
        
        # Calculate sentiment-based signals
        total_articles = (request.positive_articles + 
                         request.negative_articles + 
                         request.neutral_articles)
        
        if total_articles == 0:
            sentiment_ratio = request.sentiment_score
        else:
            positive_ratio = request.positive_articles / total_articles
            negative_ratio = request.negative_articles / total_articles
            sentiment_ratio = (positive_ratio - negative_ratio)
        
        # Combine real sentiment with article data
        combined_sentiment = (request.sentiment_score * 0.6 + sentiment_ratio * 0.4)
        
        # Generate prediction
        if combined_sentiment > 0.3:
            action = 'BUY'
            confidence = min(0.95, 0.5 + abs(combined_sentiment) * 0.7)
            risk_level = self._calculate_risk(confidence, 'BUY')
        elif combined_sentiment < -0.3:
            action = 'SELL'
            confidence = min(0.95, 0.5 + abs(combined_sentiment) * 0.7)
            risk_level = self._calculate_risk(confidence, 'SELL')
        else:
            action = 'HOLD'
            confidence = 0.5 + abs(combined_sentiment) * 0.3
            risk_level = 'low'
        
        # Determine confidence range
        if confidence >= 0.9:
            confidence_range = 'very_high'
        elif confidence >= 0.7:
            confidence_range = 'high'
        elif confidence >= 0.5:
            confidence_range = 'medium'
        elif confidence >= 0.3:
            confidence_range = 'low'
        else:
            confidence_range = 'very_low'
        
        # Generate reasoning
        reasoning, factors = self._generate_reasoning(request, combined_sentiment, action)
        
        return PredictionResponse(
            stock_symbol=request.stock_symbol,
            action=action,
            confidence_score=confidence,
            confidence_range=confidence_range,
            risk_level=risk_level,
            reasoning=reasoning,
            key_factors=factors,
            sentiment_score=combined_sentiment,
            timestamp=datetime.now().isoformat()
        )
    
    def _calculate_risk(self, confidence: float, action: str) -> str:
        """Calculate risk level based on confidence and action"""
        if confidence >= 0.85:
            base_risk = 'very_high' if action == 'BUY' else 'very_low'
        elif confidence >= 0.75:
            base_risk = 'high' if action == 'BUY' else 'low'
        elif confidence >= 0.65:
            base_risk = 'medium'
        else:
            base_risk = 'low' if action == 'BUY' else 'high'
        
        return base_risk
    
    def _generate_reasoning(self, request: PredictionRequest, 
                           sentiment: float, action: str) -> tuple:
        """Generate reasoning and key factors"""
        factors = []
        
        # Article sentiment factors
        if request.positive_articles > request.negative_articles:
            factors.append('Increased positive sentiment in news')
        elif request.negative_articles > request.positive_articles:
            factors.append('Increased negative sentiment in news')
        
        if request.neutral_articles > (request.positive_articles + request.negative_articles):
            factors.append('Mixed market sentiment')
        
        # Keyword factors
        if request.keywords:
            bullish_keywords = ['gain', 'surge', 'grow', 'profit', 'bull', 'bullish']
            bearish_keywords = ['loss', 'drop', 'fall', 'bear', 'bearish', 'crash']
            
            bullish_count = sum(1 for k in request.keywords if k.lower() in bullish_keywords)
            bearish_count = sum(1 for k in request.keywords if k.lower() in bearish_keywords)
            
            if bullish_count > bearish_count:
                factors.append('Bullish keywords detected')
            elif bearish_count > bullish_count:
                factors.append('Bearish keywords detected')
        
        # Sentiment trend
        if sentiment > 0.5:
            factors.append('Strong positive sentiment trend')
        elif sentiment < -0.5:
            factors.append('Strong negative sentiment trend')
        
        # Generate reasoning text
        if action == 'BUY':
            reasoning = f'Positive sentiment indicators ({request.positive_articles} positive articles) suggest upward momentum.'
        elif action == 'SELL':
            reasoning = f'Negative sentiment indicators ({request.negative_articles} negative articles) suggest downward pressure.'
        else:
            reasoning = 'Mixed sentiment signals suggest holding current position.'
        
        return reasoning, factors[:3]  # Limit to top 3 factors


# Initialize prediction engine
engine = PredictionEngine()


# Health check
@app.get('/health')
async def health():
    """Health check endpoint"""
    return {'status': 'healthy', 'service': 'Stock Prediction Service'}


# Prediction endpoints
@app.post('/api/predict', response_model=PredictionResponse)
async def predict(request: PredictionRequest):
    """Generate stock prediction"""
    try:
        prediction = engine.predict(request)
        return prediction
    except Exception as e:
        logger.error(f'Prediction error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


@app.post('/api/predict-batch', response_model=List[PredictionResponse])
async def predict_batch(requests: List[PredictionRequest]):
    """Generate batch predictions"""
    try:
        predictions = [engine.predict(req) for req in requests]
        return predictions
    except Exception as e:
        logger.error(f'Batch prediction error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


@app.get('/api/market-sentiment')
async def market_sentiment(
    timeframe: str = Query('1day', regex='^(1hour|1day|1week|1month)$'),
    limit: int = Query(20, ge=1, le=100)
):
    """Get overall market sentiment"""
    try:
        # This would typically query a database of recent predictions
        return {
            'timeframe': timeframe,
            'overall_sentiment': 'neutral',
            'details': {
                'bullish': 0,
                'bearish': 0,
                'neutral': 0,
            }
        }
    except Exception as e:
        logger.error(f'Market sentiment error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


@app.get('/api/stock/{symbol}/prediction')
async def get_stock_prediction(symbol: str):
    """Get latest prediction for a stock"""
    try:
        # This would query a database for the latest prediction
        return {
            'stock_symbol': symbol,
            'message': 'No data available yet'
        }
    except Exception as e:
        logger.error(f'Prediction lookup error: {e}')
        raise HTTPException(status_code=500, detail=str(e))


# Root endpoint
@app.get('/')
async def root():
    """Root endpoint"""
    return {
        'service': 'Stock Prediction Service',
        'version': '1.0.0',
        'endpoints': [
            '/health',
            '/api/predict',
            '/api/predict-batch',
            '/api/market-sentiment',
            '/api/stock/{symbol}/prediction',
        ]
    }


if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='0.0.0.0', port=8002, workers=1)
