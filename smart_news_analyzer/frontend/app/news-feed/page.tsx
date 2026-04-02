'use client';

import React, { useState } from 'react';
import useSWR from 'swr';
import axios from 'axios';

const fetcher = (url: string) => axios.get(url).then(res => res.data);

export default function NewsFeed() {
  const [sentiment, setSentiment] = useState('all');
  const [source, setSource] = useState('all');

  const { data: articles, isLoading } = useSWR(
    sentiment !== 'all' 
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/articles/by-sentiment/?sentiment=${sentiment}`
      : `${process.env.NEXT_PUBLIC_API_BASE_URL}/articles/`,
    fetcher,
    { revalidateOnFocus: false }
  );

  return (
    <div className="space-y-8">
      <h1 className="text-4xl font-bold text-gray-900">Financial News Feed</h1>

      {/* Filters */}
      <div className="bg-white p-6 rounded-lg shadow flex gap-4">
        <div>
          <label className="block text-sm font-medium mb-2">Sentiment</label>
          <select 
            value={sentiment}
            onChange={(e) => setSentiment(e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="all">All</option>
            <option value="positive">Positive</option>
            <option value="negative">Negative</option>
            <option value="neutral">Neutral</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Source</label>
          <select 
            value={source}
            onChange={(e) => setSource(e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="all">All Sources</option>
            <option value="moneycontrol">Moneycontrol</option>
            <option value="economic_times">Economic Times</option>
            <option value="youtube">YouTube</option>
          </select>
        </div>
      </div>

      {/* Articles List */}
      <div className="space-y-4">
        {isLoading ? (
          <p className="text-gray-500">Loading articles...</p>
        ) : articles?.results?.length > 0 ? (
          articles.results.map((article: any) => (
            <div key={article.id} className="bg-white p-6 rounded-lg shadow hover:shadow-lg transition">
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-xl font-bold text-gray-900 flex-1">{article.title}</h3>
                <span className={`px-3 py-1 rounded text-sm font-semibold text-white ${
                  article.sentiment === 'positive' ? 'bg-green-500' :
                  article.sentiment === 'negative' ? 'bg-red-500' :
                  'bg-gray-500'
                }`}>
                  {article.sentiment?.toUpperCase()}
                </span>
              </div>

              <p className="text-gray-700 mb-4 line-clamp-2">{article.summary || article.content}</p>

              <div className="flex justify-between items-center text-sm text-gray-500">
                <span>{article.source_name}</span>
                <span>{new Date(article.published_date).toLocaleDateString()}</span>
              </div>

              <a 
                href={article.original_url} 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-green-600 hover:text-green-700 font-semibold mt-4 block"
              >
                Read Full Article →
              </a>
            </div>
          ))
        ) : (
          <p className="text-gray-500">No articles found.</p>
        )}
      </div>
    </div>
  );
}
