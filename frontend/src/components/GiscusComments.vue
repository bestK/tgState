<template>
  <div v-if="config.enabled" class="giscus-comments">
    <div ref="giscusContainer" class="w-full min-h-[200px]"></div>
    <div v-if="config.repo.includes('your-username')" class="text-center text-sm text-gray-500 mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
      <p>💡 请在 <code>frontend/src/config/giscus.ts</code> 中配置你的 GitHub 仓库</p>
      <p class="mt-2 text-xs">访问 <a href="https://giscus.app/zh-CN" target="_blank" class="text-blue-600 hover:underline">giscus.app</a> 获取配置参数</p>
    </div>
  </div>
  <div v-else class="text-center text-gray-500 py-8">
    <p>评论系统已禁用</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { giscusConfig } from '@/config/giscus'

const giscusContainer = ref<HTMLElement>()
const config = giscusConfig

onMounted(() => {
  if (giscusContainer.value && config.enabled) {
    const script = document.createElement('script')
    script.src = 'https://giscus.app/client.js'
    script.setAttribute('data-repo', config.repo)
    script.setAttribute('data-repo-id', config.repoId)
    script.setAttribute('data-category', config.category)
    script.setAttribute('data-category-id', config.categoryId)
    script.setAttribute('data-mapping', config.mapping)
    script.setAttribute('data-strict', config.strict ? '1' : '0')
    script.setAttribute('data-reactions-enabled', config.reactionsEnabled ? '1' : '0')
    script.setAttribute('data-emit-metadata', config.emitMetadata ? '1' : '0')
    script.setAttribute('data-input-position', config.inputPosition)
    script.setAttribute('data-theme', config.theme)
    script.setAttribute('data-lang', config.lang)
    script.setAttribute('crossorigin', 'anonymous')
    script.async = true
    
    giscusContainer.value.appendChild(script)
  }
})
</script>

<style scoped>
.giscus-comments {
  max-width: 100%;
}
</style>