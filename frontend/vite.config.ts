import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import VueDevTools from 'vite-plugin-vue-devtools'
import net from 'net'

function getFreePort(startPort: number = 5173): Promise<number> {
  return new Promise((resolve) => {
    const server = net.createServer()
    server.on('error', () => {
      resolve(getFreePort(startPort + 1))
    })
    server.listen(startPort, () => {
      server.close(() => {
        resolve(startPort)
      })
    })
  })
}

export default defineConfig(async () => {
  const port = await getFreePort(5173)

  return {
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
      port: port,
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
  }
})
