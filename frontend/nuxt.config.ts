// https://nuxt.com/docs/api/configuration/nuxt-config
import { defineNuxtConfig } from 'nuxt/config';
import Aura from "@primevue/themes/aura";

export default defineNuxtConfig({
  app: {
    head: {
      title: 'TCC Felipe 2025', // Default title for all pages
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
        // { name: 'description', content: t('footer.left_text_one') },
        { name: 'google-site-verification', content: '' },
        { property: 'og:title', content: 'TCC Felipe 2025' },
        { property: 'og:url', content: import.meta.env.NUXT_PUBLIC_SITE_URL },
        { property: 'og:image', content: import.meta.env.NUXT_PUBLIC_SITE_URL + '/assets/tcc_256x256.webp' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'TCC Felipe 2025' },
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:site', content: import.meta.env.NUXT_PUBLIC_SITE_URL },
        { name: 'twitter:title', content: 'TCC Felipe 2025' },
        // { name: 'twitter:description', content: t('footer.left_text_one') },
        { name: 'twitter:image', content: import.meta.env.NUXT_PUBLIC_SITE_URL + '/assets/tcc_256x256.webp' },
        { name: 'keywords', content: 'tcc, tcc 2025, golang, python'}
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      ],
    },
  },

  compatibilityDate: '2024-04-03',
  devtools: { enabled: true },
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
    server: {
      allowedHosts: [
        'localhost',
        '127.0.0.1',
        '10.0.0.156',
        'tcc.rest',
        'tcc.com.br',
      ],
    },
  },

  primevue: {
    options: {
        theme: {
            preset: Aura,
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
