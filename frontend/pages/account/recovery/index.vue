<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center p-4 sm:p-6">
    <Toast />

    <div class="w-full max-w-md study-card p-8 border-[#dce1d9] bg-[#ffffff] shadow-none space-y-6">
      <div class="space-y-2 text-center">
        <div
          class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-sky-900/5 border border-slate-200 text-sky-800 mb-1"
        >
          <i class="pi pi-envelope text-lg"></i>
        </div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Recuperar senha</h1>
        <p class="text-sm text-slate-600">
          {{
            step === 1
              ? 'Informe o e-mail da sua conta e enviaremos um código para criar uma nova senha.'
              : 'Digite o código que enviamos para o seu e-mail e escolha uma nova senha.'
          }}
        </p>
      </div>

      <div
        role="note"
        class="flex gap-3 p-3 rounded-lg border text-sm leading-relaxed"
        :class="
          isMicrosoftEmail
            ? 'bg-amber-100 border-amber-400 text-amber-900'
            : 'bg-amber-50 border-amber-200 text-amber-900'
        "
      >
        <i class="pi pi-exclamation-triangle mt-0.5"></i>
        <p>
          Estamos com problemas para enviar e-mails para endereços da Microsoft
          (<strong>Hotmail, Outlook, Live e MSN</strong>). Se o seu e-mail for de um desses, o código
          pode demorar, cair no lixo eletrônico ou não chegar. Se não receber, procure o professor
          responsável para redefinir sua senha.
        </p>
      </div>

      <form v-if="step === 1" @submit.prevent="requestCode" class="flex flex-col gap-4">
        <div class="flex flex-col gap-1.5">
          <label for="recovery-email" class="text-sm font-semibold text-slate-600">E-mail</label>
          <InputText
            id="recovery-email"
            v-model="email"
            type="email"
            autocomplete="email"
            placeholder="seu.email@exemplo.com"
            required
            class="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
        </div>

        <TurnstileWidget ref="turnstileRef" v-model="turnstileToken" />

        <button
          type="submit"
          :disabled="submitting"
          :aria-busy="submitting"
          class="study-btn-primary w-full mt-2 py-3 rounded-lg text-sm font-semibold"
        >
          <span>{{ submitting ? 'Enviando…' : 'Enviar código por e-mail' }}</span>
          <i class="pi pi-send text-sm"></i>
        </button>
      </form>

      <form v-else @submit.prevent="resetPassword" class="flex flex-col gap-4">
        <p class="text-sm text-slate-600">
          Código enviado para <strong class="text-slate-800">{{ email }}</strong>. Ele vale por 15
          minutos. Confira também a caixa de spam.
        </p>

        <div class="flex flex-col gap-1.5">
          <label for="recovery-code" class="text-sm font-semibold text-slate-600">Código</label>
          <InputText
            id="recovery-code"
            v-model="code"
            inputmode="numeric"
            autocomplete="one-time-code"
            maxlength="6"
            placeholder="000000"
            required
            class="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-center text-lg tracking-[0.5em] font-mono"
            fluid
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="recovery-password" class="text-sm font-semibold text-slate-600">Nova senha</label>
          <Password
            inputId="recovery-password"
            v-model="newPassword"
            autocomplete="new-password"
            placeholder="Mínimo de 9 caracteres"
            :feedback="false"
            toggleMask
            class="w-full"
            inputClass="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="recovery-password-again" class="text-sm font-semibold text-slate-600"
            >Repita a nova senha</label
          >
          <Password
            inputId="recovery-password-again"
            v-model="newPasswordAgain"
            autocomplete="new-password"
            placeholder="Digite a senha de novo"
            :feedback="false"
            toggleMask
            class="w-full"
            inputClass="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
        </div>

        <button
          type="submit"
          :disabled="submitting"
          :aria-busy="submitting"
          class="study-btn-primary w-full mt-2 py-3 rounded-lg text-sm font-semibold"
        >
          <span>{{ submitting ? 'Salvando…' : 'Salvar nova senha' }}</span>
          <i class="pi pi-check text-sm"></i>
        </button>

        <button
          type="button"
          class="text-sm font-semibold text-sky-800 hover:text-sky-900"
          @click="backToEmail"
        >
          Não recebeu? Enviar outro código
        </button>
      </form>

      <div class="pt-4 border-t border-[#dce1d9] text-center">
        <p class="text-sm text-slate-600">
          Lembrou a senha?
          <NuxtLink to="/account/login" class="font-semibold text-sky-800 hover:text-sky-900 ml-1">
            Voltar para o login
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const config = useRuntimeConfig();

const step = ref<1 | 2>(1);
const submitting = ref(false);
const email = ref('');
const code = ref('');
const newPassword = ref('');
const newPasswordAgain = ref('');
const turnstileToken = ref('');
const turnstileRef = ref<{ reset: () => void } | null>(null);

const MICROSOFT_DOMAINS = /@(hotmail|outlook|live|msn|windowslive)\./i;
const isMicrosoftEmail = computed(() => MICROSOFT_DOMAINS.test(email.value.trim()));

async function post(path: string, body: Record<string, unknown>) {
  const response = await fetch(config.public.BACKEND_API_URL + path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  const data = await response.json().catch(() => ({}));
  return { ok: response.ok, data };
}

function connectionError(error: unknown) {
  console.error(error);
  toast.add({
    severity: 'error',
    summary: 'Erro de conexão',
    detail: 'Não foi possível conectar ao servidor.',
    life: 3000,
  });
}

async function requestCode() {
  if (submitting.value) return;
  if (!turnstileToken.value) {
    toast.add({ severity: 'warn', summary: 'Confirme a verificação de segurança', life: 3000 });
    return;
  }
  submitting.value = true;
  try {
    const { ok, data } = await post('/v1/recovery/request', {
      email: email.value.trim(),
      turnstile_token: turnstileToken.value,
    });
    if (ok) {
      toast.add({ severity: 'success', summary: 'Código enviado', detail: data.message, life: 4000 });
      step.value = 2;
    } else {
      toast.add({
        severity: 'error',
        summary: 'Não foi possível enviar o código',
        detail: data.error || 'Tente novamente.',
        life: 5000,
      });
    }
  } catch (error) {
    connectionError(error);
  } finally {
    submitting.value = false;
    turnstileRef.value?.reset();
  }
}

async function resetPassword() {
  if (submitting.value) return;
  if (newPassword.value.length < 9 || newPassword.value.length > 100) {
    toast.add({ severity: 'warn', summary: 'A nova senha deve ter entre 9 e 100 caracteres.', life: 3000 });
    return;
  }
  if (newPassword.value !== newPasswordAgain.value) {
    toast.add({ severity: 'warn', summary: 'As senhas não coincidem.', life: 3000 });
    return;
  }
  submitting.value = true;
  try {
    const { ok, data } = await post('/v1/recovery/reset', {
      email: email.value.trim(),
      code: code.value.trim(),
      new_password: newPassword.value,
    });
    if (ok) {
      toast.add({ severity: 'success', summary: 'Senha alterada', detail: data.message, life: 3000 });
      setTimeout(() => navigateTo('/account/login'), 1500);
    } else {
      toast.add({
        severity: 'error',
        summary: 'Não foi possível alterar a senha',
        detail: data.error || 'Tente novamente.',
        life: 4000,
      });
    }
  } catch (error) {
    connectionError(error);
  } finally {
    submitting.value = false;
  }
}

function backToEmail() {
  code.value = '';
  newPassword.value = '';
  newPasswordAgain.value = '';
  step.value = 1;
}
</script>
