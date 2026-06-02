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
import { onMounted } from 'vue'
import { useServerStore } from '@/stores/server'
import AppLayout from '@/components/AppLayout.vue'
import ServerSetup from '@/views/ServerSetup.vue'

const serverStore = useServerStore()

onMounted(() => {
  document.documentElement.classList.add('dark')
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
