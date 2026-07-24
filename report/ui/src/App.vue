<template>
  <div class="app">
    <header class="header">
      <div class="header-content">
        <h1 @click="currentView = 'list'" style="cursor:pointer">⚡ LoadForge 测试报告</h1>
        <nav>
          <button :class="{ active: currentView === 'list' }" @click="currentView = 'list'">测试列表</button>
        </nav>
      </div>
    </header>
    <main class="main">
      <TestList v-if="currentView === 'list'" @select-test="showTest" />
      <TestDetail v-else :test-id="selectedTestId" @back="currentView = 'list'" />
    </main>
  </div>
</template>

<script>
import TestList from './components/TestList.vue'
import TestDetail from './components/TestDetail.vue'

export default {
  components: { TestList, TestDetail },
  data() {
    return { currentView: 'list', selectedTestId: null }
  },
  methods: {
    showTest(id) { this.selectedTestId = id; this.currentView = 'detail' }
  }
}
</script>

<style>
.app { min-height: 100vh; }
.header {
  background: linear-gradient(135deg, #1a2a3a 0%, #0f1923 100%);
  border-bottom: 1px solid #2a3a4a;
  padding: 16px 24px;
}
.header-content {
  max-width: 1200px; margin: 0 auto;
  display: flex; justify-content: space-between; align-items: center;
}
.header h1 { font-size: 20px; color: #00d4ff; }
.header nav button {
  background: #1a2a3a; color: #a0b0c0; border: 1px solid #2a3a4a;
  padding: 8px 16px; border-radius: 6px; cursor: pointer;
}
.header nav button.active { background: #00d4ff33; color: #00d4ff; border-color: #00d4ff; }
.main { max-width: 1200px; margin: 0 auto; padding: 24px; }
</style>
