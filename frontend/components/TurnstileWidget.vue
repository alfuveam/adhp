<template>
  <div ref="container" class="flex justify-center min-h-[65px]"></div>
</template>

<script setup lang="ts">
// Widget do Cloudflare Turnstile (renderização explícita).
// v-model recebe o token; chame reset() pelo ref depois de cada envio.
import { onBeforeUnmount, onMounted, ref } from 'vue';

const token = defineModel<string>({ default: '' });

const config = useRuntimeConfig();
const container = ref<HTMLElement | null>(null);
let widgetId: string | null = null;

declare global {
  interface Window {
    turnstile?: any;
    __turnstileLoading?: Promise<void>;
  }
}

function loadScript(): Promise<void> {
  if (window.turnstile) return Promise.resolve();
  if (!window.__turnstileLoading) {
    window.__turnstileLoading = new Promise((resolve, reject) => {
      const s = document.createElement('script');
      s.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit';
      s.async = true;
      s.defer = true;
      s.onload = () => resolve();
      s.onerror = () => {
        window.__turnstileLoading = undefined;
        reject(new Error('Falha ao carregar o Turnstile'));
      };
      document.head.appendChild(s);
    });
  }
  return window.__turnstileLoading;
}

onMounted(async () => {
  try {
    await loadScript();
  } catch (e) {
    console.error(e);
    return;
  }
  if (!container.value) return;
  widgetId = window.turnstile.render(container.value, {
    sitekey: config.public.TURNSTILE_SITE_KEY || '1x00000000000000000000AA',
    language: 'pt-br',
    callback: (t: string) => (token.value = t),
    'expired-callback': () => (token.value = ''),
    'error-callback': () => (token.value = ''),
  });
});

onBeforeUnmount(() => {
  if (widgetId !== null) window.turnstile?.remove(widgetId);
});

function reset() {
  token.value = '';
  if (widgetId !== null) window.turnstile?.reset(widgetId);
}

defineExpose({ reset });
</script>
