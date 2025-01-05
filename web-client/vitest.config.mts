import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
    plugins: [tsconfigPaths(), react()],
    test: {
      environment: 'jsdom',
      globals: true,
      setupFiles: './vitest.setup.ts',
      exclude: [
        '**/e2e/**',          // Exclure tous les fichiers dans le dossier e2e
        '**/tests-examples/**', // Exclure tous les fichiers dans tests-examples
        '**/node_modules/**',   // Assurez-vous que node_modules est exclu
      ],
    },
  });