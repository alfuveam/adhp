// https://nuxt.com/docs/api/configuration/nuxt-config
import { defineNuxtConfig } from 'nuxt/config';
import Aura from "@primevue/themes/aura";

export default defineNuxtConfig({
  app: {
    head: {
      title: 'ADHP: AMBIENTE PARA O DESENVOLVIMENTO DE HABILIDADES EM PROGRAMAÇÃO', // Default title for all pages
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
        // { name: 'description', content: t('footer.left_text_one') },
        { name: 'google-site-verification', content: '' },
        { property: 'og:title', content: 'ADHP: AMBIENTE PARA O DESENVOLVIMENTO DE HABILIDADES EM PROGRAMAÇÃO' },
        { property: 'og:url', content: import.meta.env.NUXT_PUBLIC_SITE_URL },
        { property: 'og:image', content: import.meta.env.NUXT_PUBLIC_SITE_URL + '/assets/adhp_256x256.webp?v=2' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'ADHP: AMBIENTE PARA O DESENVOLVIMENTO DE HABILIDADES EM PROGRAMAÇÃO' },
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:site', content: import.meta.env.NUXT_PUBLIC_SITE_URL },
        { name: 'twitter:title', content: 'ADHP: AMBIENTE PARA O DESENVOLVIMENTO DE HABILIDADES EM PROGRAMAÇÃO' },
        // { name: 'twitter:description', content: t('footer.left_text_one') },
        { name: 'twitter:image', content: import.meta.env.NUXT_PUBLIC_SITE_URL + '/assets/adhp_256x256.webp?v=2' },
        { name: 'keywords', content: 'adhp, golang, python' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/icon.svg' },
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico?v=2' },
        { rel: 'icon', type: 'image/png', href: '/favicon.png?v=2' },
        { rel: 'apple-touch-icon', href: '/favicon.png?v=2' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&display=swap' },
      ],
    },
  },

  compatibilityDate: '2024-04-03',
  devtools: { enabled: true },
  css: [
    'primeicons/primeicons.css',
  ],
  modules: [
    '@nuxtjs/tailwindcss',
    '@primevue/nuxt-module',
    'nuxt-monaco-editor'
  ],
  runtimeConfig: {
    public: {
      BACKEND_API_URL: process.env.NUXT_PUBLIC_NUXT_API_URL,
      USER_TYPE_DISCENTE: process.env.NUXT_PUBLIC_DISCENTE,
      USER_TYPE_DOCENTE: process.env.NUXT_PUBLIC_DOCENTE,
      USE_MONACO_EDITOR: process.env.NUXT_USE_MONACO_EDITOR,
      METRICAS: {
        INICIO: 1,
        TENTOU: 2,
        RODOU: 3
      },
      METRICAS_FEEDBACK: {
        EXERCICIO: 1,
        REPETICAO: 2,
      }
    },
  },

  vite: {
    optimizeDeps: {
      // Like Vue itself, PrimeVue ships ESM. Keep its components on the same
      // runtime graph instead of rebundling them as lazy pages are discovered.
      exclude: ['primevue', '@primevue/forms', '@primevue/core', '@primevue/icons'],
      include: ['zod', 'chart.js/auto', 'quill'],
    },
    server: {
      // hostnames derivados das origens do CORS_ORIGIN ("*" ou vazio libera tudo)
      allowedHosts: process.env.CORS_ORIGIN && process.env.CORS_ORIGIN.trim() !== '*'
        ? process.env.CORS_ORIGIN.split(',').map((o) => {
            o = o.trim();
            try { return new URL(o).hostname; } catch { return o; }
          })
        : true,
    },
  },

  primevue: {
    options: {
      theme: {
        preset: Aura,
        options: { darkModeSelector: false },
      },
      ripple: true,
    },
    autoImport: true,
  },

  plugins: [
    '~/plugins/authService.ts',
    '~/plugins/localCookies.ts',
  ]
})
