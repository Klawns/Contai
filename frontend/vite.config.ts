import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

const manualVendorChunks = [
  {
    name: 'vendor-react',
    packages: ['react', 'react-dom', 'react-router-dom'],
  },
  {
    name: 'vendor-query',
    packages: ['@tanstack/react-query', 'axios'],
  },
  {
    name: 'vendor-charts',
    packages: ['recharts'],
  },
  {
    name: 'vendor-motion',
    packages: ['motion'],
  },
  {
    name: 'vendor-forms',
    packages: ['react-hook-form', '@hookform/resolvers', 'zod', '@daypicker/react'],
  },
  {
    name: 'vendor-ui',
    packages: ['lucide-react'],
  },
]

function isPackage(id: string, packageName: string) {
  const normalizedId = id.replaceAll('\\', '/')
  const packagePath = `/node_modules/${packageName}/`

  return normalizedId.includes(packagePath)
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@edusites/bancos-brasil': fileURLToPath(
        new URL('./node_modules/@edusites/bancos-brasil/src/core.js', import.meta.url),
      ),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rolldownOptions: {
      output: {
        manualChunks(id) {
          for (const chunk of manualVendorChunks) {
            if (chunk.packages.some((packageName) => isPackage(id, packageName))) {
              return chunk.name
            }
          }

          return undefined
        },
      },
    },
  },
})
