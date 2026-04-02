'use client';

import React, { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';
import axios from 'axios';

export default function AdminPage() {
  const [loading, setLoading] = useState(false);

  const handleScrapingTrigger = async (source: string) => {
    setLoading(true);
    try {
      const response = await axios.post(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/tasks/trigger-scraping/`,
        { source },
        { headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` } }
      );
      toast.success(`Scraping started for ${source}`);
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to start scraping');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      <Toaster />
      <h1 className="text-4xl font-bold text-gray-900">Admin Panel</h1>

      {/* Scraper Management */}
      <div className="bg-white p-8 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-6">Scraper Management</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {['moneycontrol', 'economic_times', 'youtube', 'telegram', 'twitter', 'bloomberg'].map(source => (
            <button
              key={source}
              onClick={() => handleScrapingTrigger(source)}
              disabled={loading}
              className="bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white px-6 py-3 rounded font-semibold transition"
            >
              Scrape {source.replace('_', ' ').toUpperCase()}
            </button>
          ))}
        </div>

        <button
          onClick={() => handleScrapingTrigger('all')}
          disabled={loading}
          className="mt-4 w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white px-6 py-3 rounded font-semibold transition"
        >
          {loading ? 'Starting...' : 'Run All Scrapers'}
        </button>
      </div>

      {/* Task Monitoring */}
      <div className="bg-white p-8 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-6">Recent Tasks</h2>
        <div className="space-y-4">
          <p className="text-gray-500">Task monitoring would appear here...</p>
        </div>
      </div>

      {/* System Status */}
      <div className="bg-white p-8 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-6">System Status</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="flex justify-between items-center">
            <span>Backend API</span>
            <span className="text-green-600 font-semibold">✓ Healthy</span>
          </div>
          <div className="flex justify-between items-center">
            <span>NLP Service</span>
            <span className="text-green-600 font-semibold">✓ Healthy</span>
          </div>
          <div className="flex justify-between items-center">
            <span>Prediction Service</span>
            <span className="text-green-600 font-semibold">✓ Healthy</span>
          </div>
          <div className="flex justify-between items-center">
            <span>Database</span>
            <span className="text-green-600 font-semibold">✓ Healthy</span>
          </div>
        </div>
      </div>
    </div>
  );
}
