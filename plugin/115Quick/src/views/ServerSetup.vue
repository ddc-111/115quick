<template>
  <div class="server-setup">
    <el-card class="setup-card">
      <template #header>
        <div class="card-header">
          <h2>连接服务器</h2>
          <p>请输入115Quick服务器地址</p>
        </div>
      </template>
      <el-form :model="form" label-position="top">
        <el-form-item label="服务器地址">
          <el-input
            v-model="form.serverUrl"
            placeholder="http://192.168.1.100:8889"
            size="large"
          >
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleConnect"
            style="width: 100%"
          >
            连接
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import { useServerStore } from '@/stores/server'

const serverStore = useServerStore()
const loading = ref(false)
const form = ref({
  serverUrl: 'http://localhost:8889'
})

async function handleConnect() {
  if (!form.value.serverUrl) {
    ElMessage.warning('请输入服务器地址')
    return
  }

  loading.value = true
  try {
    const success = await serverStore.testConnection(form.value.serverUrl)
    if (success) {
      serverStore.setServerUrl(form.value.serverUrl)
      ElMessage.success('连接成功')
    } else {
      ElMessage.error('连接失败，请检查服务器地址')
    }
  } catch {
    ElMessage.error('连接失败，请检查网络')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.server-setup {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--el-bg-color);
}

.setup-card {
  width: 400px;
  background-color: var(--el-bg-color-overlay);
  border: 1px solid rgba(255, 255, 255, 0.1);

  .card-header {
    text-align: center;

    h2 {
      font-size: 24px;
      margin-bottom: 8px;
    }

    p {
      color: var(--el-text-color-regular);
    }
  }
}
</style>
