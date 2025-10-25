import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  test: {
    environment: 'jsdom',
    include: ['src/**/__tests__/**/*.{test,spec}.{js,ts}'],
    setupFiles: ['./vitest.setup.ts'],
    globals: true,
    deps: {
      inline: ['@inertiajs/svelte'],
    },
    // Enable type checking in tests
    typecheck: {
      enabled: true,
    }
  },
});
