import { defineConfig } from 'vite'
    import vue from '@vitejs/plugin-vue'

    // https://vite.dev/config/
    export default defineConfig({
      plugins: [vue()],
      base: './',
      server: {
        proxy: {
          '/api': {
            target: 'http://localhost:3000',
            changeOrigin: true,
            // rewrite: (path) => path.replace(/^\/api/, '')
          },
          '/ws': {
            target: 'ws://localhost:3000',
            changeOrigin: true,
            ws: true
          }
        }
      }
    })