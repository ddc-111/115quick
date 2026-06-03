import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { cpSync, existsSync } from 'fs'

function copyManifest() {
  return {
    name: 'copy-manifest',
    closeBundle() {
      const src = resolve(__dirname, 'manifest.json')
      const dest = resolve(__dirname, 'dist/manifest.json')
      if (existsSync(src)) {
        cpSync(src, dest)
        console.log('✓ manifest.json copied')
      }

      const iconsSrc = resolve(__dirname, 'public/icons')
      const iconsDest = resolve(__dirname, 'dist/icons')
      if (existsSync(iconsSrc)) {
        cpSync(iconsSrc, iconsDest, { recursive: true })
        console.log('✓ icons copied')
      }
    }
  }
}

export default defineConfig({
  plugins: [vue(), copyManifest()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: 'dist',
    rollupOptions: {
      input: {
        main: 'index.html',
        content: resolve(__dirname, 'src/content/index.ts'),
        background: resolve(__dirname, 'src/background/index.ts')
      },
      output: {
        entryFileNames: (chunkInfo) => {
          if (chunkInfo.name === 'content') return 'content.js'
          if (chunkInfo.name === 'background') return 'background.js'
          return 'assets/[name]-[hash].js'
        }
      }
    }
  }
})
