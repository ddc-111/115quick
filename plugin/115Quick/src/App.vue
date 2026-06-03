<template>
  <div class="app-container dark">
    <el-container v-if="serverStore.isConnected">
      <el-aside width="200px">
        <AppLayout />
      </el-aside>
      <el-main>
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
    <ServerSetup v-else />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useServerStore } from '@/stores/server'
import AppLayout from '@/components/AppLayout.vue'
import ServerSetup from '@/views/ServerSetup.vue'

const serverStore = useServerStore()
const router = useRouter()

// 处理来自background的消息
function handleMessage(message: any) {
  if (message.type === 'NAVIGATE_TO' && message.path) {
    router.push(message.path);
  }
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  chrome.runtime.onMessage.addListener(handleMessage);
})

onUnmounted(() => {
  chrome.runtime.onMessage.removeListener(handleMessage);
})
</script>

<style lang="scss">
.app-container {
  min-height: 100vh;
  background-color: var(--el-bg-color);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
