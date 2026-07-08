import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import VueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig({
  plugins: [
    vue(),
    VueDevTools({
      launchEditor: resolve(__dirname, 'antigravity-launcher.bat')
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 5173,
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          vue: ['vue', 'vue-router', 'pinia'],
          codemirror: ['codemirror', '@codemirror/lang-sql', '@codemirror/view', '@codemirror/state'],
          tabulator: ['tabulator-tables'],
        },
      },
    },
  },
})
