<template>
  <div class="test-detail">
    <button class="back-btn" @click="$emit('back')">← 返回列表</button>

    <div v-if="loading" class="loading">加载中...</div>

    <template v-if="data">
      <div class="overview-cards">
        <div class="card"><span class="num">{{ data.result?.total_requests }}</span><span class="lbl">总请求</span></div>
        <div class="card"><span class="num">{{ (data.result?.qps || 0).toFixed(1) }}</span><span class="lbl">QPS</span></div>
        <div class="card"><span class="num success">{{ data.result?.success_count }}</span><span class="lbl">成功</span></div>
        <div class="card"><span class="num" :class="(data.result?.fail_rate || 0) > 0.05 ? 'danger' : ''">{{ data.result?.fail_count }}</span><span class="lbl">失败</span></div>
        <div class="card"><span class="num">{{ (data.result?.latency?.avg || 0).toFixed(1) }}ms</span><span class="lbl">平均延迟</span></div>
        <div class="card"><span class="num">{{ (data.result?.throughput_mb_s || 0).toFixed(2) }}</span><span class="lbl">吞吐量 MB/s</span></div>
      </div>

      <div class="section">
        <h2>总览</h2>
        <div class="overview-grid">
          <div class="info-item"><span class="key">目标 URL</span><span class="val">{{ data.config?.url }}</span></div>
          <div class="info-item"><span class="key">测试时间</span><span class="val">{{ formatTime(data.result?.timestamp) }}</span></div>
          <div class="info-item"><span class="key">持续时间</span><span class="val">{{ (data.result?.duration_seconds || 0).toFixed(2) }}s</span></div>
          <div class="info-item"><span class="key">请求配置</span><span class="val">{{ data.config?.method || 'GET' }} | {{ data.config?.concurrency || 0 }} 并发 | {{ data.config?.num_requests || 0 }} 请求</span></div>
        </div>
      </div>

      <div class="charts-row">
        <div class="chart-box"><h2>QPS 曲线</h2><div ref="qpsChart" style="height:300px"></div></div>
        <div class="chart-box"><h2>延迟曲线</h2><div ref="latChart" style="height:300px"></div></div>
      </div>

      <div class="charts-row">
        <div class="chart-box"><h2>状态码分布</h2><div ref="statusChart" style="height:300px"></div></div>
        <div class="chart-box"><h2>错误统计</h2><div ref="errorChart" style="height:300px"></div></div>
      </div>

      <div class="section">
        <h2>延迟分析</h2>
        <div class="latency-grid">
          <div class="lat-item"><span class="lbl">P50</span><span class="val">{{ latVal('p50_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">P75</span><span class="val">{{ latVal('p75_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">P90</span><span class="val">{{ latVal('p90_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">P95</span><span class="val">{{ latVal('p95_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">P99</span><span class="val">{{ latVal('p99_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">P999</span><span class="val">{{ latVal('p999_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">Min</span><span class="val">{{ latVal('min_ms') }}</span></div>
          <div class="lat-item"><span class="lbl">Max</span><span class="val">{{ latVal('max_ms') }}</span></div>
        </div>
      </div>

      <div class="section">
        <h2>详细信息</h2>
        <div class="detail-grid">
          <div class="info-item"><span class="key">HTTP方法</span><span class="val">{{ data.config?.method || 'GET' }}</span></div>
          <div class="info-item"><span class="key">请求Header</span><span class="val">{{ headerStr }}</span></div>
          <div class="info-item"><span class="key">Body大小</span><span class="val">{{ (data.config?.body?.length || 0) }} bytes</span></div>
          <div class="info-item"><span class="key">平均响应大小</span><span class="val">{{ formatBytes(data.result?.avg_response_size_bytes) }}</span></div>
          <div class="info-item"><span class="key">上传流量</span><span class="val">{{ formatBytes(data.result?.bytes_sent) }}</span></div>
          <div class="info-item"><span class="key">下载流量</span><span class="val">{{ formatBytes(data.result?.bytes_received) }}</span></div>
        </div>
      </div>

      <div v-if="tips.length" class="section tips">
        <h2>💡 性能分析建议</h2>
        <div v-for="(tip, i) in tips" :key="i" class="tip">{{ tip }}</div>
        <div v-if="data.result?.latency?.p99_ms > data.result?.latency?.p95_ms * 1.5" class="tip">⚠️ 发现 P99 延迟突增，服务端可能存在瓶颈</div>
        <div class="tip-sub">可能原因: 服务端压力过高 / 数据库响应慢 / 连接池不足</div>
      </div>
    </template>
  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  props: { testId: String },
  emits: ['back'],
  data() {
    return { data: null, loading: true }
  },
  computed: {
    tips() { return this.data?.result?.performance_tips || [] },
    headerStr() {
      const h = this.data?.config?.headers
      return h && Object.keys(h).length ? Object.entries(h).map(([k,v]) => `${k}:${v}`).join(', ') : '无'
    }
  },
  async mounted() {
    await this.loadData()
    this.$nextTick(() => this.renderCharts())
  },
  methods: {
    async loadData() {
      try {
        const r = await fetch(`/api/tests/${this.testId}`)
        this.data = await r.json()
        this.loading = false
      } catch(e) { console.error(e); this.loading = false }
    },
    latVal(key) {
      const v = this.data?.result?.latency?.[key]
      return v != null ? v.toFixed(2) + 'ms' : '-'
    },
    formatTime(ts) { return ts ? new Date(parseInt(ts)).toLocaleString('zh-CN') : '-' },
    formatBytes(b) {
      if (!b) return '0 B'
      const u = ['B','KB','MB','GB']
      let i = 0
      let v = b
      while (v >= 1024 && i < u.length - 1) { v /= 1024; i++ }
      return v.toFixed(2) + ' ' + u[i]
    },
    renderCharts() {
      if (!this.data?.timeline?.length) return

      const times = this.data.timeline.map((_, i) => i + 1 + 's')
      const qpsData = this.data.timeline.map(t => +(t.qps || 0).toFixed(1))
      const latData = this.data.timeline.map(t => +(t.latency || 0).toFixed(1))

      this.initChart(this.$refs.qpsChart, 'QPS (req/s)', times, qpsData, '#00d4ff')
      this.initChart(this.$refs.latChart, '延迟 (ms)', times, latData, '#ff6b6b')

      // 状态码分布
      if (this.data.result?.status_codes) {
        const sc = this.data.result.status_codes
        this.initPieChart(this.$refs.statusChart, '状态码', Object.keys(sc), Object.values(sc))
      }

      // 错误统计
      if (this.data.result?.errors && Object.keys(this.data.result.errors).length) {
        const err = this.data.result.errors
        this.initPieChart(this.$refs.errorChart, '错误类型', Object.keys(err).map(k => k.substring(0,20)), Object.values(err))
      }
    },
    initChart(el, name, x, y, color) {
      if (!el) return
      const chart = echarts.init(el)
      chart.setOption({
        tooltip: { trigger: 'axis' },
        grid: { left: 60, right: 20, top: 30, bottom: 30 },
        xAxis: { type: 'category', data: x, axisLabel: { color: '#8899aa' }, axisLine: { lineStyle: { color: '#2a3a4a' } } },
        yAxis: { type: 'value', axisLabel: { color: '#8899aa' }, splitLine: { lineStyle: { color: '#1a2a3a' } } },
        series: [{ type: 'line', data: y, smooth: true, symbol: 'none', lineStyle: { color, width: 2 }, areaStyle: { color: new echarts.graphic.LinearGradient(0,0,0,1,[{offset:0,color:color+'44'},{offset:1,color:color+'00'}]) } }]
      })
    },
    initPieChart(el, name, keys, vals) {
      if (!el) return
      const chart = echarts.init(el)
      chart.setOption({
        tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
        series: [{
          type: 'pie', radius: ['35%', '60%'], center: ['50%', '50%'],
          data: keys.map((k, i) => ({ name: k, value: vals[i] })),
          label: { color: '#a0b0c0' },
          labelLine: { lineStyle: { color: '#2a3a4a' } },
          itemStyle: { borderRadius: 4, borderWidth: 2, borderColor: '#0f1923' }
        }]
      })
    }
  }
}
</script>

<style scoped>
.back-btn { background: #1a2a3a; color: #00d4ff; border: 1px solid #2a3a4a; padding: 8px 16px; border-radius: 6px; cursor: pointer; margin-bottom: 20px; }
.loading { text-align: center; padding: 60px; color: #667788; }
.overview-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 12px; margin-bottom: 24px; }
.card { background: #1a2a3a; border: 1px solid #2a3a4a; border-radius: 10px; padding: 16px; text-align: center; }
.card .num { display: block; font-size: 24px; font-weight: 700; color: #00d4ff; }
.card .num.success { color: #4caf50; }
.card .num.danger { color: #ff5252; }
.card .lbl { display: block; font-size: 12px; color: #667788; margin-top: 4px; }
.section { background: #1a2a3a; border: 1px solid #2a3a4a; border-radius: 10px; padding: 20px; margin-bottom: 20px; }
.section h2 { font-size: 16px; color: #00d4ff; margin-bottom: 16px; }
.overview-grid, .detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.info-item { display: flex; justify-content: space-between; padding: 8px 12px; background: #0f1923; border-radius: 6px; }
.info-item .key { color: #667788; font-size: 13px; }
.info-item .val { color: #e0e0e0; font-size: 13px; word-break: break-all; }
.charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
.chart-box { background: #1a2a3a; border: 1px solid #2a3a4a; border-radius: 10px; padding: 16px; }
.chart-box h2 { font-size: 14px; color: #8899aa; margin-bottom: 8px; }
.latency-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.lat-item { background: #0f1923; border-radius: 8px; padding: 12px; text-align: center; }
.lat-item .lbl { display: block; color: #667788; font-size: 12px; margin-bottom: 4px; }
.lat-item .val { font-size: 20px; font-weight: 700; color: #ffd93d; }
.tips .tip { background: #ffd93d11; border-left: 3px solid #ffd93d; padding: 10px 14px; margin-bottom: 8px; border-radius: 0 6px 6px 0; color: #ffd93d; }
.tip-sub { margin-top: 8px; color: #667788; font-size: 13px; }
</style>
