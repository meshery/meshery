'use client';

import React, { useEffect, useState } from 'react';
import useSWR from 'swr';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement } from 'chart.js';
import { Doughnut, Line } from 'react-chartjs-2';
import axios from 'axios';

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement);

const fetcher = (url: string) =>  axios.get(url).then(res => res.data);

export default function Dashboard() {
  const [sentiment, setSentiment] = useState({buy: 0, hold: 0, sell: 0});
  const [loading, setLoading] = useState(true);

  const { data: predictions } = useSWR(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/predictions/market-sentiment/`,
    fetcher,
    { revalidateOnFocus: false }
  );

  useEffect(() => {
    if (predictions) {
      setSentiment(predictions.details || { buy: 0, hold: 0, sell: 0 });
      setLoading(false);
    }
  }, [predictions]);

  const chartData = {
    labels: ['BUY', 'HOLD', 'SELL'],
    datasets: [{
      data: [sentiment.buy, sentiment.hold, sentiment.sell],
      backgroundColor: ['#10B981', '#F59E0B', '#EF4444'],
      borderColor: ['#047857', '#D97706', '#DC2626'],
      borderWidth: 2,
    }],
  };

  return (
    <div className="space-y-8">
      <h1 className="text-4xl font-bold text-gray-900">Market Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold mb-4">Sentiment Distribution</h3>
          {!loading && <Doughnut data={chartData} />}
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold mb-4">Key Metrics</h3>
          <div className="space-y-3">
            <p className="text-gray-700">Total Articles: <span className="font-bold">0</span></p>
            <p className="text-gray-700">Active Predictions: <span className="font-bold">0</span></p>
            <p className="text-gray-700">Avg Confidence: <span className="font-bold">0%</span></p>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold mb-4">Market Sentiment</h3>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span>Bullish</span>
              <span className="text-green-600 font-bold">{sentiment.buy}%</span>
            </div>
            <div className="flex justify-between">
              <span>Bearish</span>
              <span className="text-red-600 font-bold">{sentiment.sell}%</span>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-4">Top Trending News</h2>
        <div className="space-y-4">
          {/* Placeholder for trending articles */}
          <p className="text-gray-500">Loading trending articles...</p>
        </div>
      </div>
    </div>
  );
}
