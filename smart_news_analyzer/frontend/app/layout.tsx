'use client';

import React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <title>Smart News Analyzer</title>
        <meta name="description" content="AI-powered financial news intelligence" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body className="bg-gray-50">
        <nav className="bg-white shadow-md">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <Link href="/" className="text-2xl font-bold text-green-600">
                📊 Smart News
              </Link>
              
              <div className="hidden md:flex space-x-8">
                <Link href="/dashboard" className="text-gray-700 hover:text-green-600">
                  Dashboard
                </Link>
                <Link href="/news-feed" className="text-gray-700 hover:text-green-600">
                  News
                </Link>
                <Link href="/predictions" className="text-gray-700 hover:text-green-600">
                  Predictions
                </Link>
                <Link href="/admin" className="text-gray-700 hover:text-red-600">
                  Admin
                </Link>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {children}
        </main>

        <footer className="bg-gray-100 py-8 mt-12">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-gray-600">
            <p>&copy; 2024 Smart News Analyzer. All rights reserved.</p>
          </div>
        </footer>
      </body>
    </html>
  );
}
