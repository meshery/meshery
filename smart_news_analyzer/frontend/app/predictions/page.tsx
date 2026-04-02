'use client';

import React from 'react';
import useSWR from 'swr';
import axios from 'axios';

const fetcher = (url: string) => axios.get(url).then(res => res.data);

export default function Predictions() {
  const { data: predictions, isLoading } = useSWR(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/predictions/`,
    fetcher,
    { revalidateOnFocus: false }
  );

  const actionColor = (action: string) => {
    switch(action?.toUpperCase()) {
      case 'BUY': return 'bg-green-100 text-green-800';
      case 'SELL': return 'bg-red-100 text-red-800';
      default: return 'bg-yellow-100 text-yellow-800';
    }
  };

  return (
    <div className="space-y-8">
      <h1 className="text-4xl font-bold text-gray-900">Stock Predictions</h1>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-green-50 p-6 rounded-lg border-l-4 border-green-500">
          <h3 className="text-lg font-semibold text-green-700 mb-2">Buy Signals</h3>
          <p className="text-3xl font-bold text-green-600">
            {predictions?.results?.filter((p: any) => p.action === 'BUY').length || 0}
          </p>
        </div>

        <div className="bg-yellow-50 p-6 rounded-lg border-l-4 border-yellow-500">
          <h3 className="text-lg font-semibold text-yellow-700 mb-2">Hold Signals</h3>
          <p className="text-3xl font-bold text-yellow-600">
            {predictions?.results?.filter((p: any) => p.action === 'HOLD').length || 0}
          </p>
        </div>

        <div className="bg-red-50 p-6 rounded-lg border-l-4 border-red-500">
          <h3 className="text-lg font-semibold text-red-700 mb-2">Sell Signals</h3>
          <p className="text-3xl font-bold text-red-600">
            {predictions?.results?.filter((p: any) => p.action === 'SELL').length || 0}
          </p>
        </div>
      </div>

      {/* Predictions Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full">
          <thead className="bg-gray-100 border-b">
            <tr>
              <th className="px-6 py-3 text-left font-semibold">Symbol</th>
              <th className="px-6 py-3 text-left font-semibold">Action</th>
              <th className="px-6 py-3 text-left font-semibold">Confidence</th>
              <th className="px-6 py-3 text-left font-semibold">Risk Level</th>
              <th className="px-6 py-3 text-left font-semibold">Details</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                  Loading predictions...
                </td>
              </tr>
            ) : predictions?.results?.length > 0 ? (
              predictions.results.map((pred: any) => (
                <tr key={pred.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 font-semibold">{pred.stock_symbol}</td>
                  <td className="px-6 py-4">
                    <span className={`px-3 py-1 rounded font-semibold ${actionColor(pred.action)}`}>
                      {pred.action?.toUpperCase()}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center">
                      <div className="w-24 bg-gray-200 rounded-full h-2 mr-2">
                        <div 
                          className="bg-blue-500 h-2 rounded-full"
                          style={{width: `${pred.confidence_score * 100}%`}}
                        ></div>
                      </div>
                      <span className="text-sm font-semibold">{(pred.confidence_score * 100).toFixed(0)}%</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 capitalize">{pred.risk_level}</td>
                  <td className="px-6 py-4 text-blue-600 hover:text-blue-800 cursor-pointer">
                    View →
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                  No predictions available yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
