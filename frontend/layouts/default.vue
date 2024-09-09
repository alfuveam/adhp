<template>
  <div class="p-3 bg-gray-700 min-h-screen text-gray-100">
    <header>
      <div class="w-full h-14 bg-gray-800 rounded">
        <div class="p-1 flex justify-between">
            <NuxtLink to="/" class="flex flex-center">
                <img src="@/assets/images/logo_tcc.webp" class="h-12 w-12" alt="">
                <span class="justify-center text-3xl font-bold text-gray-300 p-1">TCC</span>
            </NuxtLink>
            <div></div>
            <div class="flex flex-center justify-between p-3">
                <NuxtLink to="/login" class="mr-4">Entrar</NuxtLink>
                <NuxtLink to="/createaccount">Criar Conta</NuxtLink>
            </div>
        </div>
      </div>
    </header>

    <div class="w-full min-h-screen">
      <slot />
    </div>

    <footer>
        <div class="w-full h-24 bg-gray-800 rounded">
            <div class="flex justify-between">
                <div></div>
                <div class="flex justify-between p-3">
                    <NuxtLink to="/legal" class="mr-4">legal</NuxtLink>
                    <NuxtLink to="/contact">contact</NuxtLink>
                </div>
            </div>
            <hr class="ml-4 mr-4 border-gray-600"/>
            <div class="flex justify-center mt-2">TCC - Todos os direitos reservados</div>
        </div>
    </footer>
  </div>
</template>