import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [
    vue({
      script: {
        defineModel: true,
        propsDestructure: true
      }
    })
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8088',
        changeOrigin: true,
      },
      '/d/': {
        target: 'http://localhost:8088',
        changeOrigin: true,
      },
      '/s/': {
        target: 'http://localhost:8088',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: '../assets/dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'axios'],
          ui: ['lucide-vue-next']
        }
      }
    }
  },
  esbuild: {
    target: 'es2020'
  }
})