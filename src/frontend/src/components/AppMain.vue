<template>
  <main class="flex-1 py-8 px-4 sm:px-6 lg:px-8 bg-gray-50 dark:bg-gray-950 transition-colors duration-300">
    <div class="w-full max-w-[960px] mx-auto">
      <div class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-4 gap-3 mb-6">
        <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center transition-colors duration-300">
          <div class="text-3xl font-bold text-blue-600 dark:text-blue-400" :class="{ 'animate-pulse': isLoading }">{{ nodes.length }}</div>
          <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">注册实例</div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden transition-colors duration-300">
        <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">节点列表 <span class="text-sm font-normal text-gray-500 dark:text-gray-400">{{ instanceCount }}</span></h2>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 dark:bg-gray-900/50">
              <tr>
                <th class="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-300 border-b border-gray-200 dark:border-gray-700">域名</th>
                <th class="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-300 border-b border-gray-200 dark:border-gray-700">节点URL</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="isLoading">
                <td colspan="2" class="text-center px-4 py-10 text-gray-500 dark:text-gray-400">加载中...</td>
              </tr>
              <tr v-else-if="nodes.length === 0">
                <td colspan="2" class="text-center px-4 py-10 text-gray-500 dark:text-gray-400">暂无节点数据</td>
              </tr>
              <tr v-for="node in nodes" :key="node.name" class="border-b border-gray-100 dark:border-gray-700 last:border-b-0 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <td class="px-4 py-3 text-gray-900 dark:text-gray-100">{{ node.name }}</td>
                <td class="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">{{ node.url }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const isLoading = ref(false)
const nodes = ref([])
let timer = null
let countdown = 30

const instanceCount = computed(() => nodes.value.length > 0 ? `(${nodes.value.length})` : '')

const loadData = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/nodes/public')
    const data = await res.json()
    if (data.nodes) {
      nodes.value = data.nodes
    }
  } catch (e) {
    console.error('Failed to load data:', e)
  } finally {
    isLoading.value = false
  }
}

timer = setInterval(() => {
  countdown--
  if (countdown <= 0) {
    countdown = 30
    loadData()
  }
}, 1000)

onMounted(() => {
  loadData()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

defineExpose({ loadData })
</script>
