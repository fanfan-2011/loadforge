<template>
  <div class="test-list">
    <div class="stats-overview">
      <div class="stat-card"><span class="num">{{ tests.length }}</span><span class="label">总测试数</span></div>
      <div class="stat-card"><span class="num">{{ totalReqs }}</span><span class="label">总请求数</span></div>
      <div class="stat-card"><span class="num">{{ avgQps.toFixed(1) }}</span><span class="label">平均 QPS</span></div>
    </div>

    <div v-if="tests.length === 0" class="empty">
      <p>暂无测试记录</p>
      <p class="hint">运行 <code>loadforge bench -n 1000 -c 10 https://example.com --report</code> 开始测试</p>
    </div>

    <div v-else class="test-table">
      <div class="table-header">
        <span class="col-time">测试时间</span>
        <span class="col-url">目标 URL</span>
        <span class="col-req">请求数</span>
        <span class="col-qps">QPS</span>
        <span class="col-status">状态</span>
        <span class="col-action">操作</span>
      </div>
      <div v-for="t in sortedTests" :key="t.id" class="table-row">
        <span class="col-time">{{ formatTime(t.id) }}</span>
        <span class="col-url" :title="t.url">{{ t.url }}</span>
        <span class="col-req">{{ t.total_reqs || t.total_requests }}</span>
        <span class="col-qps">{{ (t.qps || 0).toFixed(1) }}</span>
        <span class="col-status"><span class="badge done">{{ t.status }}</span></span>
        <span class="col-action"><button @click="$emit('select-test', t.id)" class="view-btn">查看报告</button></span>
      </div>
    </div>
  </div>
</template>

<script>
const API = ''

export default {
  data() {
    return { tests: [] }
  },
  computed: {
    sortedTests() { return [...this.tests].reverse() },
    totalReqs() { return this.tests.reduce((s, t) => s + (t.total_reqs || t.total_requests || 0), 0) },
    avgQps() {
      const valid = this.tests.filter(t => t.qps > 0)
      return valid.length ? valid.reduce((s, t) => s + t.qps, 0) / valid.length : 0
    }
  },
  mounted() { this.fetchTests() },
  methods: {
    async fetchTests() {
      try {
        const r = await fetch(`${API}/api/tests`)
        this.tests = await r.json()
      } catch(e) { console.error('获取测试列表失败:', e) }
    },
    formatTime(ts) {
      const d = new Date(parseInt(ts))
      return isNaN(d.getTime()) ? ts : d.toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.stats-overview { display: flex; gap: 16px; margin-bottom: 24px; }
.stat-card {
  background: #1a2a3a; border: 1px solid #2a3a4a; border-radius: 12px;
  padding: 20px 32px; flex: 1; text-align: center;
}
.stat-card .num { display: block; font-size: 32px; font-weight: 700; color: #00d4ff; }
.stat-card .label { display: block; font-size: 13px; color: #8899aa; margin-top: 4px; }
.empty { text-align: center; padding: 60px 20px; color: #667788; }
.empty code { background: #1a2a3a; padding: 2px 6px; border-radius: 4px; font-size: 13px; }
.hint { margin-top: 12px; font-size: 14px; }
.test-table { background: #1a2a3a; border: 1px solid #2a3a4a; border-radius: 12px; overflow: hidden; }
.table-header, .table-row { display: flex; padding: 12px 16px; align-items: center; }
.table-header { background: #0f1923; font-size: 12px; color: #667788; text-transform: uppercase; }
.table-row { border-top: 1px solid #2a3a4a; font-size: 14px; }
.table-row:hover { background: #1e3040; }
.col-time { width: 180px; }
.col-url { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.col-req { width: 100px; text-align: right; }
.col-qps { width: 100px; text-align: right; }
.col-status { width: 80px; text-align: center; }
.col-action { width: 100px; text-align: center; }
.badge { display: inline-block; padding: 2px 10px; border-radius: 10px; font-size: 12px; }
.badge.done { background: #00d4ff22; color: #00d4ff; }
.view-btn {
  background: #00d4ff; color: #000; border: none; padding: 6px 14px;
  border-radius: 6px; cursor: pointer; font-size: 12px; font-weight: 600;
}
.view-btn:hover { background: #00e5ff; }
</style>
