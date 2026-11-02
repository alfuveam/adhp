<template>
  <div class="space-y-20 lg:space-y-28 py-6 sm:py-10">
    
    <!-- ========================================== -->
    <!-- HERO SECTION (Rich Interactive EdTech)    -->
    <!-- ========================================== -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-12 items-center">
        
        <!-- Left Hero Content -->
        <div class="lg:col-span-6 space-y-6 text-center lg:text-left">
          
          <!-- Academic Origin Badge -->
          <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-[#1c203b] text-purple-300 text-xs font-bold border border-[#323963]">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>UNESC • Ciência da Computação</span>
          </div>

          <!-- Main Value Headline -->
          <h1 class="text-3xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight text-white leading-[1.15]">
            Aprenda lógica de programação com <span class="text-[#c084fc]">testes automáticos</span> e repetição espaçada.
          </h1>

          <!-- Concise, Human Subhead -->
          <p class="text-base sm:text-lg text-slate-300 max-w-xl mx-auto lg:mx-0 leading-relaxed font-normal">
            Escreva código em Python e Go, receba validação instantânea no sandbox de testes e fixe o aprendizado de longo prazo com o método científico da curva de Ebbinghaus.
          </p>

          <!-- Primary Actions (Tactile Mimo 3D Buttons) -->
          <div class="pt-2 flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
            <template v-if="!isAuthenticated">
              <NuxtLink 
                to="/account/createaccount" 
                class="w-full sm:w-auto mimo-btn-primary px-7 py-3.5 text-base"
              >
                <span>Começar Gratuitamente</span>
                <i class="pi pi-arrow-right text-xs"></i>
              </NuxtLink>

              <NuxtLink 
                to="/account/login" 
                class="w-full sm:w-auto mimo-btn-secondary px-6 py-3.5 text-base"
              >
                <span>Acessar Conta</span>
              </NuxtLink>
            </template>

            <template v-else>
              <button 
                @click="redirectToDashboard" 
                class="w-full sm:w-auto mimo-btn-primary px-7 py-3.5 text-base"
              >
                <i class="pi pi-th-large text-sm"></i>
                <span>Acessar Meu Painel</span>
              </button>
            </template>
          </div>

          <!-- Proof / Feature Highlights -->
          <div class="pt-3 flex flex-wrap items-center justify-center lg:justify-start gap-3 text-xs font-semibold text-slate-300">
            <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-[#171a30] border border-[#2b3258]">
              <i class="pi pi-bolt text-emerald-400 text-sm"></i>
              <span>Testes no Sandbox em &lt; 50ms</span>
            </div>
            <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-[#171a30] border border-[#2b3258]">
              <i class="pi pi-history text-purple-400 text-sm"></i>
              <span>Curva de Retenção Ebbinghaus</span>
            </div>
            <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-[#171a30] border border-[#2b3258]">
              <i class="pi pi-code text-cyan-400 text-sm"></i>
              <span>Ambientes Go &amp; Python</span>
            </div>
          </div>

        </div>

        <!-- Right Hero Mockup: Interactive Code & Test Workstation -->
        <div class="lg:col-span-6">
          <div class="mimo-card p-0 overflow-hidden shadow-2xl border-[#323963] bg-[#121424]">
            
            <!-- Window Titlebar with Language Switcher -->
            <div class="flex items-center justify-between px-4 py-3 bg-[#191d36] border-b border-[#282f54]">
              <!-- Window Control Dots -->
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full bg-rose-500/80"></span>
                <span class="w-3 h-3 rounded-full bg-amber-500/80"></span>
                <span class="w-3 h-3 rounded-full bg-emerald-500/80"></span>
              </div>

              <!-- Interactive Language Tabs -->
              <div class="flex items-center p-0.5 rounded-lg bg-[#111322] border border-[#272d50]">
                <button 
                  id="tab-btn-python"
                  type="button"
                  @click="activeCodeTab = 'python'"
                  class="px-3 py-1 rounded-md text-xs font-bold transition flex items-center gap-1.5 cursor-pointer"
                  :class="activeCodeTab === 'python' ? 'bg-[#7c4ed4] text-white shadow-sm' : 'text-slate-400 hover:text-white'"
                >
                  <span class="w-2 h-2 rounded-full bg-amber-400"></span>
                  <span>Python</span>
                </button>
                <button 
                  id="tab-btn-golang"
                  type="button"
                  @click="activeCodeTab = 'golang'"
                  class="px-3 py-1 rounded-md text-xs font-bold transition flex items-center gap-1.5 cursor-pointer"
                  :class="activeCodeTab === 'golang' ? 'bg-[#0284c7] text-white shadow-sm' : 'text-slate-400 hover:text-white'"
                >
                  <span class="w-2 h-2 rounded-full bg-cyan-400"></span>
                  <span>Golang</span>
                </button>
              </div>

              <!-- Status Tag -->
              <span class="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider bg-emerald-500/15 text-emerald-300 border border-emerald-500/30 flex items-center gap-1">
                <i class="pi pi-check text-[9px]"></i>
                <span>{{ activeCodeTab === 'python' ? 'Python 3.12' : 'Go 1.22' }}</span>
              </span>
            </div>

            <!-- Code Editor Body (Python Tab) -->
            <div v-if="activeCodeTab === 'python'" class="p-5 font-mono text-xs leading-relaxed space-y-1 bg-[#0f111f] text-slate-200">
              <div class="text-slate-500 text-[11px] pb-1"># Desafio: Soma de Elementos Pares (Lista 01)</div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">1</span>
                <div><span class="text-purple-400 font-semibold">def</span> <span class="text-yellow-300 font-semibold">soma_pares</span>(valores: list[int]) -> int:</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">2</span>
                <div class="pl-4">total = <span class="text-amber-300">0</span></div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">3</span>
                <div class="pl-4"><span class="text-purple-400 font-semibold">for</span> num <span class="text-purple-400 font-semibold">in</span> valores:</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">4</span>
                <div class="pl-8"><span class="text-purple-400 font-semibold">if</span> num % <span class="text-amber-300">2</span> == <span class="text-amber-300">0</span>:</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">5</span>
                <div class="pl-12">total += num</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">6</span>
                <div class="pl-4"><span class="text-purple-400 font-semibold">return</span> total</div>
              </div>
            </div>

            <!-- Code Editor Body (Golang Tab) -->
            <div v-else class="p-5 font-mono text-xs leading-relaxed space-y-1 bg-[#0f111f] text-slate-200">
              <div class="text-slate-500 text-[11px] pb-1">// Desafio: Soma de Elementos Pares (Lista 01)</div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">1</span>
                <div><span class="text-cyan-400 font-semibold">package</span> main</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">2</span>
                <div><span class="text-cyan-400 font-semibold">func</span> <span class="text-yellow-300 font-semibold">SomaPares</span>(valores []<span class="text-cyan-300">int</span>) <span class="text-cyan-300">int</span> {</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">3</span>
                <div class="pl-4">total := <span class="text-amber-300">0</span></div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">4</span>
                <div class="pl-4"><span class="text-purple-400 font-semibold">for</span> _, num := <span class="text-purple-400 font-semibold">range</span> valores {</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">5</span>
                <div class="pl-8"><span class="text-purple-400 font-semibold">if</span> num%<span class="text-amber-300">2</span> == <span class="text-amber-300">0</span> { total += num }</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">6</span>
                <div class="pl-4">}</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">7</span>
                <div class="pl-4"><span class="text-purple-400 font-semibold">return</span> total</div>
              </div>
              <div class="flex items-start gap-3">
                <span class="text-slate-600 select-none w-4 text-right">8</span>
                <div>}</div>
              </div>
            </div>

            <!-- Test Results & Spaced Repetition Panel -->
            <div class="border-t border-[#262c4e] bg-[#15182e] p-5 space-y-4">
              
              <!-- Result Header -->
              <div class="flex items-center justify-between pb-2 border-b border-[#24294a]">
                <div class="flex items-center gap-2">
                  <span class="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center text-xs">
                    <i class="pi pi-check"></i>
                  </span>
                  <span class="text-xs sm:text-sm font-bold text-white">3 de 3 Testes Unitários Aprovados</span>
                </div>
                <span class="text-[11px] font-mono text-slate-400">Sandbox: {{ activeCodeTab === 'python' ? '0.038s' : '0.019s' }}</span>
              </div>

              <!-- Test Chips -->
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 text-[11px] font-mono">
                <div class="p-2 rounded-xl bg-[#0f1120] border border-[#252b4c] flex items-center justify-between">
                  <span class="text-slate-300">{{ activeCodeTab === 'python' ? 'test_pares' : 'TestSomaPares' }}</span>
                  <span class="text-emerald-400 font-bold flex items-center gap-1">
                    <i class="pi pi-check text-[10px]"></i> OK
                  </span>
                </div>
                <div class="p-2 rounded-xl bg-[#0f1120] border border-[#252b4c] flex items-center justify-between">
                  <span class="text-slate-300">{{ activeCodeTab === 'python' ? 'test_misto' : 'TestMistos' }}</span>
                  <span class="text-emerald-400 font-bold flex items-center gap-1">
                    <i class="pi pi-check text-[10px]"></i> OK
                  </span>
                </div>
                <div class="p-2 rounded-xl bg-[#0f1120] border border-[#252b4c] flex items-center justify-between">
                  <span class="text-slate-300">{{ activeCodeTab === 'python' ? 'test_vazia' : 'TestSliceVazio' }}</span>
                  <span class="text-emerald-400 font-bold flex items-center gap-1">
                    <i class="pi pi-check text-[10px]"></i> OK
                  </span>
                </div>
              </div>

              <!-- Visual Ebbinghaus Timeline (Spaced Repetition Progression) -->
              <div class="p-3.5 rounded-2xl bg-[#1b203c] border border-[#323963] space-y-2.5">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <i class="pi pi-calendar-plus text-purple-400 text-xs"></i>
                    <span class="text-xs font-bold text-white">Curva de Repetição Espaçada</span>
                  </div>
                  <span class="text-[10px] font-bold px-2 py-0.5 rounded-md bg-purple-500/20 text-purple-300 border border-purple-500/30">
                    Retenção Estimada: 94%
                  </span>
                </div>

                <!-- 5 Steps of Hermann Ebbinghaus -->
                <div class="grid grid-cols-5 gap-1.5 pt-1 text-center font-mono">
                  <div class="p-1.5 rounded-lg bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 text-[10px] font-bold">
                    <div>1h</div>
                    <div class="text-[9px] opacity-80">Feito ✓</div>
                  </div>
                  <div class="p-1.5 rounded-lg bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 text-[10px] font-bold">
                    <div>9h</div>
                    <div class="text-[9px] opacity-80">Feito ✓</div>
                  </div>
                  <div class="p-1.5 rounded-lg bg-purple-600 text-white font-bold text-[10px] shadow-sm animate-pulse">
                    <div>1d</div>
                    <div class="text-[9px]">Hoje ⏳</div>
                  </div>
                  <div class="p-1.5 rounded-lg bg-[#111322] border border-[#272d50] text-slate-500 text-[10px]">
                    <div>6d</div>
                    <div class="text-[9px]">Fila</div>
                  </div>
                  <div class="p-1.5 rounded-lg bg-[#111322] border border-[#272d50] text-slate-500 text-[10px]">
                    <div>31d</div>
                    <div class="text-[9px]">Fila</div>
                  </div>
                </div>
              </div>

            </div>

          </div>
        </div>

      </div>
    </section>

    <!-- ========================================== -->
    <!-- SECTION 2: COMO FUNCIONA (VISUAL WORKFLOW) -->
    <!-- ========================================== -->
    <section id="metodologia" class="py-14 bg-[#14172c]/70 border-y border-[#262c4e]">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-12">
        
        <div class="max-w-2xl space-y-2 text-center sm:text-left">
          <span class="mimo-badge mimo-badge-purple uppercase tracking-wider text-xs font-bold">
            Metodologia Ativa
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
            Como o ADHP transforma o estudo de algoritmos
          </h2>
          <p class="text-sm text-slate-300">
            A combinação de desafios estruturados, feedback imediato em sandbox e ciência da memória humana.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          
          <!-- Step 1: Trilhas -->
          <div class="mimo-card-interactive p-7 space-y-5 flex flex-col justify-between">
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="w-10 h-10 rounded-2xl bg-purple-500/20 border border-purple-500/30 text-purple-300 flex items-center justify-center font-bold text-sm">
                  01
                </span>
                <span class="text-xs font-bold uppercase tracking-wider text-purple-300 font-mono">Estrutura</span>
              </div>
              <h3 class="text-lg font-bold text-white">Desafios Progressivos</h3>
              <p class="text-sm text-slate-300 leading-relaxed">
                Listas organizadas pelo professor para guiar seu raciocínio do primeiro conceito até estruturas de dados avançadas, com requisitos claros.
              </p>
            </div>

            <!-- Mini Graphic: Roadmap Node Tree -->
            <div class="p-3.5 rounded-xl bg-[#0f1120] border border-[#252b4c] space-y-2 text-xs">
              <div class="flex items-center justify-between text-slate-300">
                <span class="flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
                  <span>1. Fundamentos &amp; Sintaxe</span>
                </span>
                <span class="text-emerald-400 font-bold">✓</span>
              </div>
              <div class="flex items-center justify-between text-white font-semibold">
                <span class="flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-purple-400"></span>
                  <span>2. Condicionais &amp; Laços</span>
                </span>
                <span class="text-purple-300 text-[11px]">Atual</span>
              </div>
              <div class="flex items-center justify-between text-slate-500">
                <span class="flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-slate-600"></span>
                  <span>3. Slices &amp; Ponteiros</span>
                </span>
                <span class="text-slate-600 text-[11px]">Bloqueado</span>
              </div>
            </div>
          </div>

          <!-- Step 2: Feedback Instantâneo -->
          <div class="mimo-card-interactive p-7 space-y-5 flex flex-col justify-between">
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="w-10 h-10 rounded-2xl bg-cyan-500/20 border border-cyan-500/30 text-cyan-300 flex items-center justify-center font-bold text-sm">
                  02
                </span>
                <span class="text-xs font-bold uppercase tracking-wider text-cyan-300 font-mono">Execução</span>
              </div>
              <h3 class="text-lg font-bold text-white">Feedback em Sandbox</h3>
              <p class="text-sm text-slate-300 leading-relaxed">
                Ao submeter o código, a plataforma executa sua solução em contêineres isolados, testando casos normais e extremos em milissegundos.
              </p>
            </div>

            <!-- Mini Graphic: Terminal Pass Box -->
            <div class="p-3.5 rounded-xl bg-[#0f1120] border border-[#252b4c] font-mono text-[11px] space-y-1.5 text-slate-300">
              <div class="text-slate-500">$ go test -v ./...</div>
              <div class="text-emerald-400 flex items-center gap-1.5">
                <i class="pi pi-check text-[10px]"></i> PASS: TestSomaPares (0.012s)
              </div>
              <div class="text-emerald-400 flex items-center gap-1.5">
                <i class="pi pi-check text-[10px]"></i> PASS: TestListaVazia (0.008s)
              </div>
              <div class="text-cyan-300 font-bold pt-1 border-t border-[#222746]">
                Resultado: 100% de cobertura
              </div>
            </div>
          </div>

          <!-- Step 3: Repetição Espaçada -->
          <div class="mimo-card-interactive p-7 space-y-5 flex flex-col justify-between">
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="w-10 h-10 rounded-2xl bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 flex items-center justify-center font-bold text-sm">
                  03
                </span>
                <span class="text-xs font-bold uppercase tracking-wider text-emerald-300 font-mono">Memória</span>
              </div>
              <h3 class="text-lg font-bold text-white">Retenção de Longo Prazo</h3>
              <p class="text-sm text-slate-300 leading-relaxed">
                Baseado na curva de Hermann Ebbinghaus. Exercícios resolvidos retornam em intervalos calculados (1h, 9h, 1d, 6d, 31d) para combater o esquecimento.
              </p>
            </div>

            <!-- Mini Graphic: Spaced Repetition vs Natural Decay -->
            <div class="p-3.5 rounded-xl bg-[#0f1120] border border-[#252b4c] space-y-2 text-xs">
              <div class="flex items-center justify-between text-slate-400">
                <span>Sem repetição:</span>
                <span class="text-rose-400 font-semibold font-mono">~20% em 30d</span>
              </div>
              <div class="w-full bg-[#1b203c] h-1.5 rounded-full overflow-hidden">
                <div class="bg-rose-500 h-full w-[20%]"></div>
              </div>
              <div class="flex items-center justify-between text-white font-semibold pt-1">
                <span class="text-emerald-300 flex items-center gap-1">
                  <i class="pi pi-check-circle text-xs"></i> Com ADHP:
                </span>
                <span class="text-emerald-400 font-bold font-mono">92% fixado</span>
              </div>
              <div class="w-full bg-[#1b203c] h-1.5 rounded-full overflow-hidden">
                <div class="bg-emerald-400 h-full w-[92%]"></div>
              </div>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- ========================================== -->
    <!-- SECTION 3: TRILHAS DE CÓDIGO (COURSES)     -->
    <!-- ========================================== -->
    <section id="trilhas" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-10">
      
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div class="space-y-2">
          <span class="mimo-badge mimo-badge-blue uppercase tracking-wider text-xs font-bold">
            Catálogo de Trilhas
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
            Linguagens Suportadas
          </h2>
          <p class="text-sm text-slate-300">
            Ambientes prontos com compilador e interpretador integrados para prática imediata no navegador.
          </p>
        </div>

        <NuxtLink 
          to="/account/createaccount" 
          class="text-sm font-bold text-purple-300 hover:text-purple-200 transition flex items-center gap-1.5"
        >
          <span>Criar conta de estudante</span>
          <i class="pi pi-arrow-right text-xs"></i>
        </NuxtLink>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <!-- Track 1: Golang -->
        <div class="mimo-card-interactive p-8 space-y-6 border-[#2d355c]">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-4">
              <span class="w-14 h-14 rounded-2xl bg-cyan-500/20 border border-cyan-500/30 text-cyan-300 flex items-center justify-center font-extrabold text-xl shadow-inner">
                Go
              </span>
              <div>
                <h3 class="text-xl font-extrabold text-white">Trilha Golang</h3>
                <span class="text-xs text-cyan-300 font-semibold">Linguagem Compilada e Concorrente</span>
              </div>
            </div>
            <span class="mimo-badge mimo-badge-blue text-[11px]">Básico ao Avançado</span>
          </div>

          <p class="text-sm text-slate-300 leading-relaxed">
            Estruture seu raciocínio com a sintaxe robusta de Go. Pratique tipos estáticos, ponteiros, slices, structs, interfaces e algoritmos concorrentes.
          </p>

          <div class="p-4 rounded-xl bg-[#0f1120] border border-[#252b4c] space-y-2.5 text-xs text-slate-300">
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-cyan-400 text-xs"></i>
              <span>Avaliação automatizada com a biblioteca nativa <code class="text-cyan-300">testing</code></span>
            </div>
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-cyan-400 text-xs"></i>
              <span>Diagnóstico detalhado de erros de compilação</span>
            </div>
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-cyan-400 text-xs"></i>
              <span>Execução em contêiner isolado no backend Go</span>
            </div>
          </div>

          <div class="pt-2">
            <NuxtLink 
              to="/account/createaccount" 
              class="w-full mimo-btn-blue py-3 text-sm"
            >
              <span>Explorar Trilha de Go</span>
              <i class="pi pi-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>

        <!-- Track 2: Python -->
        <div class="mimo-card-interactive p-8 space-y-6 border-[#2d355c]">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-4">
              <span class="w-14 h-14 rounded-2xl bg-amber-400/20 border border-amber-400/30 text-amber-300 flex items-center justify-center font-extrabold text-xl shadow-inner">
                Py
              </span>
              <div>
                <h3 class="text-xl font-extrabold text-white">Trilha Python</h3>
                <span class="text-xs text-amber-300 font-semibold">Lógica e Estruturas de Dados</span>
              </div>
            </div>
            <span class="mimo-badge mimo-badge-amber text-[11px]">Iniciante Amigável</span>
          </div>

          <p class="text-sm text-slate-300 leading-relaxed">
            Foque diretamente na resolução do algoritmo. Pratique estruturas de repetição, listas, dicionários, funções puras e manipulação de strings.
          </p>

          <div class="p-4 rounded-xl bg-[#0f1120] border border-[#252b4c] space-y-2.5 text-xs text-slate-300">
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-amber-400 text-xs"></i>
              <span>Bateria de testes unitários com asserções detalhadas</span>
            </div>
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-amber-400 text-xs"></i>
              <span>Retorno claro de stdout e stderr para depuração</span>
            </div>
            <div class="flex items-center gap-2">
              <i class="pi pi-check text-amber-400 text-xs"></i>
              <span>Raciocínio lógico acelerado com foco no problema</span>
            </div>
          </div>

          <div class="pt-2">
            <NuxtLink 
              to="/account/createaccount" 
              class="w-full mimo-btn-primary py-3 text-sm"
            >
              <span>Explorar Trilha de Python</span>
              <i class="pi pi-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>

      </div>

    </section>

    <!-- ========================================== -->
    <!-- SECTION 4: IMPACTO ACADÊMICO (DUAL IMPACT) -->
    <!-- ========================================== -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-10">
      
      <div class="max-w-2xl space-y-2">
        <span class="mimo-badge mimo-badge-emerald uppercase tracking-wider text-xs font-bold">
          Benefícios Reais
        </span>
        <h2 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
          Projetado para as duas pontas da sala de aula
        </h2>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <!-- For Students -->
        <div class="mimo-card p-8 space-y-5 border-[#2c345a]">
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 flex items-center justify-center text-xl">
            <i class="pi pi-graduation-cap"></i>
          </div>
          <h3 class="text-xl font-bold text-white">Para o Estudante (Discente)</h3>
          <ul class="space-y-3.5 text-sm text-slate-300 leading-relaxed">
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-emerald-400 mt-1 text-xs"></i>
              <span><strong>Fim da incerteza:</strong> receba feedback imediato se a lógica está correta sem precisar esperar o próximo dia de aula.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-emerald-400 mt-1 text-xs"></i>
              <span><strong>Dicas pedagógicas:</strong> orientações cadastradas pelo professor para destravar sem entregar a resposta pronta.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-emerald-400 mt-1 text-xs"></i>
              <span><strong>Memória permanente:</strong> as revisões espaçadas impedem que você esqueça a sintaxe dias antes da avaliação final.</span>
            </li>
          </ul>
        </div>

        <!-- For Teachers -->
        <div class="mimo-card p-8 space-y-5 border-[#2c345a]">
          <div class="w-12 h-12 rounded-2xl bg-purple-500/20 border border-purple-500/30 text-purple-300 flex items-center justify-center text-xl">
            <i class="pi pi-user-edit"></i>
          </div>
          <h3 class="text-xl font-bold text-white">Para o Docente (Professor)</h3>
          <ul class="space-y-3.5 text-sm text-slate-300 leading-relaxed">
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-purple-400 mt-1 text-xs"></i>
              <span><strong>Zero correção manual:</strong> a suíte de testes unitários automatizada substitui horas de compilação manual de arquivos.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-purple-400 mt-1 text-xs"></i>
              <span><strong>Gestão completa de listas:</strong> crie trilhas, organize sequências pedagógicas e configure múltiplos casos de teste.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <i class="pi pi-check text-purple-400 mt-1 text-xs"></i>
              <span><strong>Métricas de retenção:</strong> acompanhe a taxa de acerto na primeira tentativa vs nas repetições espaçadas da turma.</span>
            </li>
          </ul>
        </div>

      </div>

    </section>

    <!-- ========================================== -->
    <!-- BOTTOM CALL TO ACTION BANNER               -->
    <!-- ========================================== -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="mimo-card p-8 sm:p-12 border-[#363e69] bg-[#161a33] flex flex-col md:flex-row md:items-center justify-between gap-8">
        <div class="space-y-2 max-w-xl">
          <h3 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
            Pronto para colocar sua lógica à prova?
          </h3>
          <p class="text-sm sm:text-base text-slate-300 leading-relaxed">
            Crie sua conta gratuitamente para resolver os desafios das trilhas de Golang e Python.
          </p>
        </div>

        <div class="flex items-center gap-4 shrink-0">
          <NuxtLink 
            v-if="!isAuthenticated" 
            to="/account/createaccount" 
            class="mimo-btn-primary px-7 py-3.5 text-base"
          >
            <span>Cadastrar Como Aluno</span>
            <i class="pi pi-arrow-right text-xs"></i>
          </NuxtLink>
          <button 
            v-else 
            @click="redirectToDashboard" 
            class="mimo-btn-primary px-7 py-3.5 text-base"
          >
            <span>Ir Para o Meu Painel</span>
            <i class="pi pi-arrow-right text-xs"></i>
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

const route = useRoute();
const isAuthenticated = ref(false);
const activeCodeTab = ref<'python' | 'golang'>('python');
const { $authService } = useNuxtApp();

const redirectToDashboard = () => {
  if ($authService.isDocente()) {
    navigateTo("/dashboard_docente");
  } else if ($authService.isDiscente()) {
    navigateTo("/dashboard_discente");
  }
};

const checkIsAuthenticated = () => {
  if ($authService.isAuthenticated()) {
    try {
      const jwtPayload = $authService.parseJwt($authService.getToken() || '');
      if (jwtPayload && jwtPayload.exp && jwtPayload.exp < Date.now() / 1000) {
        $authService.logout();
      }
    } catch {
      // ignore
    }
  }
  isAuthenticated.value = $authService.isAuthenticated();
};

watch(route, () => {
  checkIsAuthenticated();
});

onMounted(() => {
  checkIsAuthenticated();
});
</script>