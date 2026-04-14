import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

export default defineConfig({
  plugins: [
    react(),
    // Replaces CRA's webpack ProvidePlugin for Buffer, process, etc.
    nodePolyfills({
      globals: {
        Buffer: true,
        process: true,
      },
      // Replaces webpack fallback config
      protocolImports: true,
    }),
  ],

  // Handles @docker/docker-mui-theme ESM/CJS interop
  // (was previously handled by webpack config-overrides.js)
  optimizeDeps: {
    include: ['@docker/docker-mui-theme'],
  },

  build: {
    outDir: 'build',        // Keep same output dir as CRA (Dockerfile expects this)
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: undefined,
      },
    },
  },

  server: {
    port: 3000,             // Keep same port as CRA dev server
    open: false,
  },
})
