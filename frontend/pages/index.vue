<template>
  <div class="home-page">
    <section class="hero">
      <div>
        <div class="eyebrow">Para quem está começando</div>
        <h1>Seu primeiro programa começa com <em>uma tentativa.</em></h1>
        <p class="intro">
          Aprenda lógica de programação colocando suas ideias em prática. Resolva exercícios em
          Python e Go, confira o resultado e tente de novo quando precisar.
        </p>
        <div class="actions">
          <NuxtLink v-if="!isAuthenticated" class="button" to="/account/createaccount"
            >Começar a praticar
            <span aria-hidden="true" style="margin-left: 15px">→</span></NuxtLink
          ><button v-else class="button" @click="redirectToDashboard">Ir para meus estudos →</button
          ><a class="textlink" href="#metodologia">Entenda como funciona</a>
        </div>
        <p class="note">Não precisa acertar de primeira. Faz parte do aprendizado.</p>
      </div>
      <div>
        <div class="example">
          <div class="example-label">Uma primeira experiência <span>Python</span></div>
          <div class="exercise">
            <h3>Vamos dizer olá?</h3>
            <p>Este programa mostra uma mensagem com seu nome.</p>
          </div>
          <div class="code">
            <span class="line">1</span>nome = <span class="str">"Ana"</span><br /><span class="line"
              >2</span
            ><span class="func">print</span>(<span class="str">"Olá, "</span> + nome)
          </div>
          <div class="run">
            <span>2 linhas. Um começo.</span
            ><button class="button" type="button" @click="exampleExecuted = true">
              {{ exampleExecuted ? 'Ver resultado novamente ↻' : 'Ver resultado ▷' }}
            </button>
          </div>
          <div class="output" v-if="exampleExecuted" id="result" role="status">
            <small>Saída do programa</small><code>Olá, Ana</code>
            <div style="margin-top: 8px; color: #67795c; font-size: 12px">
              O programa juntou a saudação com o nome escolhido.
            </div>
          </div>
        </div>
        <div class="caption">Pequenos exercícios para entender uma ideia de cada vez.</div>
      </div>
    </section>
    <section class="how" id="metodologia">
      <div class="section-top">
        <h2>Um passo de cada vez.</h2>
        <p>Da primeira tentativa à próxima descoberta.</p>
      </div>
      <div class="steps">
        <article>
          <span class="stepnum">01 / ESCOLHER</span>
          <h3>Encontre seu ponto de partida</h3>
          <p>
            Siga as listas organizadas pelo professor e pratique os conceitos que está aprendendo.
          </p>
        </article>
        <article>
          <span class="stepnum">02 / EXPERIMENTAR</span>
          <h3>Escreva, teste, entenda</h3>
          <p>
            Veja o resultado da sua solução. Se algo não funcionar, use a dica e experimente outro
            caminho.
          </p>
        </article>
        <article>
          <span class="stepnum">03 / REVISITAR</span>
          <h3>Volte ao que aprendeu</h3>
          <p>
            Retome os exercícios nas revisões para praticar novamente e fortalecer seu aprendizado.
          </p>
        </article>
      </div>
    </section>
    <section class="resources" id="recursos" aria-labelledby="resources-title">
      <div class="section-top">
        <h2 id="resources-title">Além de escrever código.</h2>
        <p>Ferramentas para acompanhar cada tentativa.</p>
      </div>
      <div class="resources-grid">
        <dl class="feature-list">
          <div>
            <dt>Um editor no próprio navegador</dt>
            <dd>
              Abra o exercício e escreva sua solução ali mesmo. O editor destaca a sintaxe de Python
              e Go para ajudar na leitura do código.
            </dd>
          </div>
          <div>
            <dt>Verificações para entender sua solução</dt>
            <dd>
              Teste o programa e confira sua saída. Quando algo não funcionar, veja uma orientação
              inicial e abra os detalhes técnicos para investigar.
            </dd>
          </div>
          <div>
            <dt>Uma dica quando você precisar</dt>
            <dd>
              Consulte as dicas cadastradas pelo professor para encontrar um próximo caminho, sem
              sair do exercício.
            </dd>
          </div>
          <div>
            <dt>Seu progresso como ponto de partida</dt>
            <dd>
              Distinga o que já concluiu do que está disponível e continue pelo próximo exercício.
              As atividades seguintes são liberadas conforme você avança.
            </dd>
          </div>
        </dl>
        <aside class="learning-preview" aria-label="Exemplos dos recursos do ambiente">
          <p class="preview-label">Durante a prática</p>
          <h3>Uma tentativa também mostra o que estudar.</h3>
          <p>
            O resultado ajuda você a decidir o próximo ajuste. As mensagens técnicas continuam
            disponíveis para consulta.
          </p>
          <div class="guidance-example">
            <span>Exemplo de orientação</span>
            <strong>Vamos ajustar sua solução</strong>
            <p>
              Confira os espaços no início das linhas. Linhas do mesmo bloco precisam ter o mesmo
              recuo.
            </p>
          </div>
          <div class="progress-example">
            <span>Na sua lista de exercícios</span>
            <ul>
              <li>
                <i class="pi pi-check-circle" aria-hidden="true"></i
                ><span>Concluído<small>Solução aprovada nas verificações</small></span>
              </li>
              <li>
                <i class="pi pi-play-circle" aria-hidden="true"></i
                ><span>Disponível<small>Pronto para você tentar</small></span>
              </li>
              <li>
                <i class="pi pi-lock" aria-hidden="true"></i
                ><span>Bloqueado<small>Conclua os anteriores para avançar</small></span>
              </li>
            </ul>
          </div>
        </aside>
      </div>
    </section>
    <section class="practice" id="trilhas">
      <div class="practice-grid">
        <div>
          <h2>O que você vai praticar</h2>
          <NuxtLink
            class="row"
            :to="
              isAuthenticated
                ? $authService.isDocente()
                  ? '/dashboard_docente'
                  : '/dashboard_discente'
                : '/account/createaccount'
            "
            ><span class="lang">Py</span
            ><span
              ><strong>Lógica de programação com Python</strong
              ><small>Variáveis, condições e repetições</small></span
            ><span class="arrow">↗</span></NuxtLink
          ><NuxtLink
            class="row"
            :to="
              isAuthenticated
                ? $authService.isDocente()
                  ? '/dashboard_docente'
                  : '/dashboard_discente'
                : '/account/createaccount'
            "
            ><span class="lang" style="background: #e5edf0; color: #3e6a7d">Go</span
            ><span
              ><strong>Lógica de programação com Go</strong
              ><small>Uma outra linguagem para explorar os fundamentos</small></span
            ><span class="arrow">↗</span></NuxtLink
          >
        </div>
        <aside class="aside">
          <h3>Aprender também é retomar.</h3>
          <p>
            As revisões espaçadas ajudam você a voltar aos exercícios ao longo do tempo. Seu painel
            reúne as atividades e as próximas revisões.
          </p>
        </aside>
      </div>
    </section>
    <section class="teacher-section" aria-labelledby="teacher-title">
      <div class="teacher-intro">
        <p class="eyebrow">Para quem ensina</p>
        <h2 id="teacher-title">Prepare o caminho de quem está aprendendo.</h2>
        <p>
          Organize a prática da sua turma com trilhas, listas e exercícios de Python e Go. Você
          define as atividades e as verificações de cada solução.
        </p>
        <NuxtLink
          class="textlink"
          :to="
            isAuthenticated
              ? $authService.isDocente()
                ? '/dashboard_docente'
                : '/dashboard_discente'
              : '/account/login'
          "
          >{{ isAuthenticated ? 'Acessar meu painel' : 'Já tenho acesso de professor' }}
          <span aria-hidden="true">→</span></NuxtLink
        >
      </div>
      <dl class="teacher-features">
        <div>
          <dt>Organize a sequência de estudo</dt>
          <dd>Crie e edite trilhas, reúna exercícios em listas e ajuste a ordem das atividades.</dd>
        </div>
        <div>
          <dt>Defina o ponto de partida e os testes</dt>
          <dd>
            Prepare o código inicial e os testes que vão verificar as soluções enviadas pelos
            estudantes.
          </dd>
        </div>
        <div>
          <dt>Cadastre dicas para as tentativas</dt>
          <dd>
            Adicione orientações aos exercícios para apoiar quem precisar de ajuda durante a
            prática.
          </dd>
        </div>
      </dl>
    </section>
    <section class="questions" aria-labelledby="questions-title">
      <h2 id="questions-title">Antes de começar</h2>
      <details>
        <summary>Preciso instalar um editor no computador?</summary>
        <p>
          Para resolver as atividades do ADHP, você usa o editor da própria página. O código é
          enviado para execução e verificação pelo ambiente.
        </p>
      </details>
      <details>
        <summary>O que acontece se eu errar um exercício?</summary>
        <p>
          Você pode ler o retorno da tentativa, consultar uma dica disponível, ajustar o código e
          testar novamente. O exercício é marcado como concluído quando sua solução passa pelas
          verificações.
        </p>
      </details>
      <details>
        <summary>Como funcionam as revisões?</summary>
        <p>
          Depois de concluir exercícios, você pode retomá-los nas revisões espaçadas. O painel
          mostra as atividades agendadas e quando cada uma fica disponível para revisar.
        </p>
      </details>
      <details>
        <summary>Posso começar sem conhecer Python ou Go?</summary>
        <p>
          O ambiente tem foco em quem está iniciando na programação. As atividades são organizadas
          pelo professor; siga a sequência da sua trilha e use as dicas disponíveis ao longo do
          caminho.
        </p>
      </details>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

const route = useRoute();
const isAuthenticated = ref(false);
const exampleExecuted = ref(false);
const { $authService } = useNuxtApp();

const redirectToDashboard = () => {
  if ($authService.isDocente()) {
    navigateTo('/dashboard_docente');
  } else if ($authService.isDiscente()) {
    navigateTo('/dashboard_discente');
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
<style scoped>
.hero {
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 78px;
  align-items: center;
  padding: 82px 0 70px;
}
.eyebrow {
  color: #596f61;
  font-size: 12px;
  font-weight: 650;
  letter-spacing: 1.4px;
  text-transform: uppercase;
}
h1 {
  font-size: 49px;
  line-height: 1.14;
  font-weight: 650;
  letter-spacing: -2px;
  max-width: 530px;
  margin: 20px 0 23px;
}
h1 em {
  font-style: normal;
  color: #285b75;
}
.intro {
  color: #65706b;
  max-width: 465px;
  font-size: 17px;
}
.actions {
  display: flex;
  align-items: center;
  gap: 22px;
  margin-top: 28px;
}
.textlink {
  font-size: 14px;
  text-decoration: underline;
  text-underline-offset: 5px;
}
.note {
  font-size: 12px;
  color: #727b74;
  margin-top: 17px;
}
.example {
  background: white;
  border: 1px solid #dce0d8;
  border-radius: 10px;
  box-shadow: 0 12px 30px #303a3210;
  overflow: hidden;
}
.example-label {
  padding: 17px 23px;
  border-bottom: 1px solid #e8eae4;
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}
.example-label span {
  color: #798278;
  font-size: 12px;
}
.exercise {
  padding: 22px 24px 18px;
}
.exercise h3 {
  font-size: 19px;
  margin: 0 0 7px;
  font-weight: 600;
}
.exercise p {
  font-size: 14px;
  color: #6b746d;
  margin: 0;
}
.code {
  margin: 0 22px;
  background: #202d35;
  color: #e7ebec;
  border-radius: 6px;
  padding: 23px 17px;
  font: 14px/1.95 monospace;
}
.line {
  color: #7e8b92;
  display: inline-block;
  width: 24px;
  user-select: none;
}
.str {
  color: #ddc08c;
}
.func {
  color: #96c2d7;
}
.run {
  padding: 16px 23px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.run button {
  font-size: 12px;
  padding: 9px 15px;
}
.run span {
  font-size: 12px;
  color: #7b817c;
}
.output {
  margin: 0 23px 23px;
  background: #f1f5ee;
  border-left: 3px solid #7d9c6e;
  padding: 12px 15px;
  font-size: 13px;
}
.output small {
  display: block;
  font-size: 10px;
  color: #6e7b67;
  letter-spacing: 0.8px;
  text-transform: uppercase;
  margin-bottom: 5px;
}
.caption {
  text-align: center;
  font-size: 12px;
  color: #798175;
  margin-top: 15px;
}
.how {
  border-top: 1px solid #dce0d7;
  padding: 39px 0 48px;
}
.section-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 25px;
}
h2 {
  font-size: 25px;
  letter-spacing: -0.6px;
  font-weight: 600;
  margin: 0 0 26px;
}
.section-top p {
  font-size: 13px;
  color: #778075;
}
.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 46px;
}
.stepnum {
  font: 13px monospace;
  color: #7d8c76;
}
.steps h3 {
  font-size: 17px;
  font-weight: 600;
  margin: 11px 0 9px;
}
.steps p {
  font-size: 14px;
  color: #697369;
  margin: 0;
}
.practice {
  padding: 35px 0 52px;
}
.practice-grid {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 70px;
}
.row {
  border-top: 1px solid #dce0d7;
  padding: 17px 0;
  display: flex;
  align-items: center;
  gap: 18px;
}
.lang {
  width: 40px;
  height: 40px;
  background: #ecece2;
  border-radius: 6px;
  display: grid;
  place-items: center;
  font: 13px monospace;
  color: #625f40;
}
.row strong {
  font-size: 15px;
  font-weight: 600;
  display: block;
}
.row small {
  font-size: 12px;
  color: #7a8175;
}
.arrow {
  margin-left: auto;
  color: #285b75;
}
.aside {
  border-left: 2px solid #c9d4c0;
  padding: 0 0 0 24px;
  align-self: center;
}
.aside h3 {
  font-size: 18px;
  font-weight: 550;
  margin: 0 0 8px;
}
.aside p {
  font-size: 14px;
  color: #697369;
  margin: 0;
}
@media (max-width: 760px) {
  .hero {
    grid-template-columns: 1fr;
    gap: 35px;
    padding: 42px 0;
  }
  h1 {
    font-size: 38px;
  }
  .steps,
  .practice-grid {
    grid-template-columns: 1fr;
    gap: 25px;
  }
  .section-top {
    display: block;
  }
  .section-top p {
    margin-top: -15px;
    margin-bottom: 25px;
  }
  .hero .example {
    transform: none;
  }
  .practice-grid {
    gap: 32px;
  }
  .actions {
    gap: 17px;
    flex-wrap: wrap;
  }
  .code {
    font-size: 13px;
  }
  .run {
    gap: 12px;
    flex-wrap: wrap;
  }
  .aside {
    margin-bottom: 10px;
  }
}

.resources {
  border-top: 1px solid #dce0d7;
  padding: 42px 0 56px;
}
.resources-grid {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 76px;
  align-items: start;
}
.feature-list > div {
  padding: 20px 0;
  border-bottom: 1px solid #dce0d7;
}
.feature-list > div:first-child {
  padding-top: 0;
}
.feature-list dt,
.teacher-features dt {
  font-size: 17px;
  font-weight: 600;
  margin-bottom: 8px;
}
.feature-list dd,
.teacher-features dd {
  font-size: 15px;
  line-height: 1.75;
  color: #65706b;
  margin: 0;
}
.learning-preview {
  padding: 26px;
  background: #fff;
  border: 1px solid #dce0d7;
  border-radius: 10px;
}
.preview-label {
  font-size: 12px;
  color: #285b75;
  margin-bottom: 12px;
}
.learning-preview h3 {
  font-size: 21px;
  line-height: 1.4;
  font-weight: 600;
  letter-spacing: -0.4px;
  margin-bottom: 12px;
}
.learning-preview > p:not(.preview-label),
.guidance-example p {
  font-size: 14px;
  color: #65706b;
}
.guidance-example {
  border-left: 3px solid #87601d;
  background: #faf8f2;
  padding: 14px 16px;
  margin: 22px 0;
}
.guidance-example > span,
.progress-example > span {
  display: block;
  font-size: 12px;
  color: #65706b;
  margin-bottom: 9px;
}
.guidance-example strong {
  font-size: 14px;
  font-weight: 600;
}
.guidance-example p {
  margin-top: 8px;
}
.progress-example ul {
  display: grid;
  gap: 15px;
}
.progress-example li {
  display: flex;
  align-items: start;
  gap: 12px;
  font-size: 14px;
}
.progress-example i {
  margin-top: 4px;
  color: #285b75;
}
.progress-example li:first-child i {
  color: #326449;
}
.progress-example li:last-child i {
  color: #65706b;
}
.progress-example small {
  display: block;
  color: #65706b;
  font-size: 12px;
  margin-top: 2px;
}
.teacher-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 76px;
  background: #edf1eb;
  border-radius: 10px;
  padding: 40px;
  margin: 16px 0 58px;
}
.teacher-intro h2 {
  margin: 14px 0 18px;
  line-height: 1.3;
}
.teacher-intro > p:not(.eyebrow) {
  color: #65706b;
  font-size: 15px;
}
.teacher-intro .textlink {
  display: inline-block;
  margin-top: 24px;
  color: #285b75;
}
.teacher-features {
  display: grid;
  gap: 24px;
}
.teacher-features dt {
  font-size: 16px;
}
.questions {
  padding-bottom: 20px;
  max-width: 820px;
}
.questions details {
  border-top: 1px solid #dce0d7;
  padding: 18px 0;
}
.questions summary {
  cursor: pointer;
  font-weight: 500;
  font-size: 16px;
}
.questions details p {
  color: #65706b;
  font-size: 15px;
  margin: 14px 0 0 16px;
  line-height: 1.8;
}
@media (max-width: 760px) {
  .resources-grid,
  .teacher-section {
    grid-template-columns: 1fr;
    gap: 30px;
  }
  .teacher-section {
    padding: 25px;
    margin-bottom: 40px;
  }
  .learning-preview {
    padding: 22px;
  }
}
</style>
