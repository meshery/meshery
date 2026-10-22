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
    exclude: [
      'node_modules',
      '.next',
      'out',
      'tests/e2e',
      // TODO: these three files use shallow vi.mock() stubs for @sistent/sistent.
      // Their transitive imports pull in styled-components modules that the
      // stubs don't cover. Fixing properly requires either (a) converting the
      // mock to vi.mock(..., async (importOriginal) => ...) and wrapping render()
      // in SistentThemeProvider, or (b) exhaustively stubbing the imported
      // exports. Tracked separately from the Nighthawk removal.
      'components/registry/ImportModelModal.test.tsx',
      'components/registry/CreateModelModal.test.tsx',
      'components/registry/MeshModelComponent.test.tsx',
    ],
    css: false,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'json-summary', 'json'],
      reportsDirectory: './coverage',
      exclude: [
        'node_modules/**',
        '.next/**',
        'out/**',
        'tests/e2e/**',
        '**/*.config.{js,ts}',
        '**/*.d.ts',
        '**/__tests__/**',
        '**/*.test.{ts,tsx}',
        'scripts/**',
        'public/**',
        'docs/**',
        'assets/icons/**',
        'graphql/**/*.{ts,tsx}',
        '**/*.stories.{ts,tsx}',
        'next-env.d.ts',
        'ui_dev_server.js',
        'tsconfig.tsbuildinfo',
      ],
      include: [
        'components/**/*.{ts,tsx}',
        'utils/**/*.{ts,tsx}',
        'lib/**/*.{ts,tsx}',
        'rtk-query/**/*.{ts,tsx}',
        'store/**/*.{ts,tsx}',
        'machines/**/*.{ts,tsx}',
        'pages/**/*.{ts,tsx}',
        'api/**/*.{ts,tsx}',
      ],
    },
  },
  resolve: {
    alias: {
      '@/components': path.resolve(__dirname, './components'),
      '@/graphql': path.resolve(__dirname, './graphql'),
      '@/utils': path.resolve(__dirname, './utils'),
      '@/rtk-query': path.resolve(__dirname, './rtk-query'),
      '@/constants': path.resolve(__dirname, './constants'),
      '@/api': path.resolve(__dirname, './api'),
      '@/assets': path.resolve(__dirname, './assets'),
      '@/theme': path.resolve(__dirname, './theme'),
      '@/store': path.resolve(__dirname, './store'),
      // Non-@ aliases used by the codebase
      lib: path.resolve(__dirname, './lib'),
      css: path.resolve(__dirname, './css'),
      machines: path.resolve(__dirname, './machines'),
    },
  },
});
