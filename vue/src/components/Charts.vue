<template>
  <div ref="chartRef" style="width: 100%; height: 400px;"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps<{ data: any }>()

const chartRef = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const handleResize = () => {
  chartInstance?.resize()
}

// 初始化图表
const initChart = () => {
  if (chartRef.value) {
    chartInstance = echarts.init(chartRef.value)
    window.addEventListener('resize', handleResize)

    if (props.data) {
      updateChart(props.data)
    }
  }
}

// 更新图表数据
const updateChart = (data: any) => {
  if (!chartInstance) return

  const categories = data.map((item: any) => item.name || `选项${item.id}`)
  const values = data.map((item: any) => item.value || 0)

  const option = {
    title: {
      text: '投票结果统计',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: categories
    },
    yAxis: {
      type: 'value',
      name: '票数'
    },
    series: [
      {
        name: '票数',
        type: 'bar',
        data: values,
        barWidth: '40%',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}'
        }
      }
    ],
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    }
  }

  // 渲染图表
  chartInstance.setOption(option)
}

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  window.removeEventListener('resize', handleResize)
})

// 监听数据变化
watch(() => props.data, (newData) => {
  if (newData && chartInstance) {
    updateChart(newData)
  }
}, { deep: true })
</script>