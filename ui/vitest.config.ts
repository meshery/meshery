import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./vitest.setup.ts'],
    include: ['**/__tests__/**/*.test.{ts,tsx}', '**/*.test.{ts,tsx}'],
    exclude: ['node_modules', '.next', 'out', 'tests/e2e'],
    css: false,
  },
  resolve: {
    alias: {
      '@/components': path.resolve(__dirname, './components'),
      '@/utils': path.resolve(__dirname, './utils'),
      '@/rtk-query': path.resolve(__dirname, './rtk-query'),
      '@/constants': path.resolve(__dirname, './constants'),
      '@/api': path.resolve(__dirname, './api'),
      '@/assets': path.resolve(__dirname, './assets'),
      '@/themes': path.resolve(__dirname, './themes'),
      '@/theme': path.resolve(__dirname, './theme'),
      '@/store': path.resolve(__dirname, './store'),
      // Non-@ aliases used by the codebase
      lib: path.resolve(__dirname, './lib'),
      css: path.resolve(__dirname, './css'),
      machines: path.resolve(__dirname, './machines'),
    },
  },
});
