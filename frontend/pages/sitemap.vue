<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Header -->
    <div class="text-center mb-10 space-y-3">
      <span class="mimo-badge mimo-badge-indigo">Navegação Geral</span>
      <h1 class="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
        Mapa do Site
      </h1>
      <p class="text-sm text-slate-400">
        Acesse rapidamente todas as rotas e seções disponíveis na plataforma ADHP
      </p>
    </div>

    <!-- Routes Card -->
    <div class="mimo-card p-6 sm:p-8 border-[#363e66] bg-[#1c2038] shadow-2xl">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <NuxtLink
          v-for="route in routes"
          :key="route"
          :to="route"
          class="p-4 rounded-xl bg-[#141728] border border-[#2d3458] hover:border-[#7c4ed4] hover:bg-[#1a1e35] transition-all duration-200 flex items-center justify-between group"
        >
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-lg bg-purple-500/10 border border-purple-500/20 text-purple-300 flex items-center justify-center text-xs group-hover:bg-purple-500/20 transition">
              <i class="pi pi-link"></i>
            </span>
            <span class="text-sm font-semibold text-slate-200 group-hover:text-white transition font-mono">
              {{ route }}
            </span>
          </div>

          <i class="pi pi-arrow-right text-xs text-slate-500 group-hover:text-purple-300 group-hover:translate-x-1 transition-all"></i>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const routes = computed(() => {
  return router.getRoutes()
    .map(r => r.path)
    .filter(p => !p.includes(':') && p !== '/sitemap')
    .sort();
});

useSeoMeta({
  title: 'Mapa do Site • ADHP',
  description: 'Navegação e rotas do Ambiente de Desenvolvimento de Habilidades em Programação.'
});
</script>