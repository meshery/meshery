'use client';

import Link from 'next/link';

export default function Home() {
  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-green-600 to-blue-600 text-white py-20 px-8 rounded-lg">
        <h1 className="text-5xl font-bold mb-4">Smart News Analyzer</h1>
        <p className="text-xl mb-8">AI-Powered Financial News Intelligence Platform</p>
        <div className="flex gap-4">
          <Link href="/dashboard" className="bg-white text-green-600 px-6 py-3 rounded font-semibold hover:bg-gray-100">
            Go to Dashboard
          </Link>
          <Link href="/news-feed" className="border-2 border-white px-6 py-3 rounded font-semibold hover:bg-white hover:text-green-600">
            Browse News
          </Link>
        </div>
      </div>

      {/* Features */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="bg-white p-8 rounded-lg shadow">
          <h3 className="text-2xl font-bold mb-4">📰 News Aggregation</h3>
          <p className="text-gray-700">Automatically scrapes financial news from multiple sources including websites, YouTube, Telegram, and Twitter.</p>
        </div>

        <div className="bg-white p-8 rounded-lg shadow">
          <h3 className="text-2xl font-bold mb-4">🤖 AI Analysis</h3>
          <p className="text-gray-700">Uses advanced NLP models (BART, FinBERT) to summarize articles and analyze sentiment.</p>
        </div>

        <div className="bg-white p-8 rounded-lg shadow">
          <h3 className="text-2xl font-bold mb-4">📈 Market Predictions</h3>
          <p className="text-gray-700">Generates BUY/HOLD/SELL signals based on sentiment analysis and market trends.</p>
        </div>
      </div>

      {/* How It Works */}
      <div className="bg-white p-8 rounded-lg shadow">
        <h2 className="text-3xl font-bold mb-8">How It Works</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="text-center">
            <div className="bg-green-100 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">1</div>
            <h4 className="font-semibold mb-2">Scrape</h4>
            <p className="text-sm text-gray-600">Collect news from multiple sources</p>
          </div>
          <div className="text-center">
            <div className="bg-blue-100 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">2</div>
            <h4 className="font-semibold mb-2">Analyze</h4>
            <p className="text-sm text-gray-600">Process with NLP models</p>
          </div>
          <div className="text-center">
            <div className="bg-yellow-100 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">3</div>
            <h4 className="font-semibold mb-2">Predict</h4>
            <p className="text-sm text-gray-600">Generate market signals</p>
          </div>
          <div className="text-center">
            <div className="bg-purple-100 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">4</div>
            <h4 className="font-semibold mb-2">Display</h4>
            <p className="text-sm text-gray-600">Show on dashboard</p>
          </div>
        </div>
      </div>
    </div>
  );
}
