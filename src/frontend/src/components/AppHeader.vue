<template>
  <header class="bg-white dark:bg-gray-900 transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center space-x-2">
          <div
            class="w-10 h-10 rounded-lg flex items-center justify-center transition-colors duration-300 bg-white dark:bg-gray-800 shadow-sm">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
              class="w-8 h-8 transition-colors duration-300 fill-gray-900 dark:fill-white" aria-hidden="true">
              <path
                d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0 1 12 6.836a9.59 9.59 0 0 1 2.504.337c1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.202 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.161 22 16.416 22 12c0-5.523-4.477-10-10-10z">
              </path>
            </svg>
          </div>
          <a href="/" class="text-xl font-semibold hover:opacity-80 transition-colors duration-300 text-gray-900 dark:text-white">GHProxy Registry</a>
        </div>
        <div class="flex items-center space-x-4">
          <span class="text-xs text-gray-500 dark:text-gray-400">自动刷新: <span class="text-blue-600 dark:text-blue-400 font-medium">{{ countdown }}</span>s</span>
          <button @click="$emit('toggleDark')"
            class="w-6 h-6 flex items-center justify-center cursor-pointer hover:opacity-80 transition-colors duration-300"
            :title="isDark ? '切换到亮色模式' : '切换到暗黑模式'">
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
              fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
              class="text-yellow-400 transition-colors duration-300">
              <circle cx="12" cy="12" r="4"></circle>
              <path d="M12 2v2"></path>
              <path d="M12 20v2"></path>
              <path d="m4.93 4.93 1.41 1.41"></path>
              <path d="m17.66 17.66 1.41 1.41"></path>
              <path d="M2 12h2"></path>
              <path d="M20 12h2"></path>
              <path d="m6.34 17.66-1.41 1.41"></path>
              <path d="m19.07 4.93-1.41 1.41"></path>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
              class="text-gray-600 transition-colors duration-300">
              <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  isDark: {
    type: Boolean,
    required: true
  },
  countdown: {
    type: Number,
    default: 30
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggleDark', 'refresh'])
</script>
