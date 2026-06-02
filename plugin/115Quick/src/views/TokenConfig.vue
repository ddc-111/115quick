<template>
  <div class="token-config page-container">
    <div class="page-header">
      <h1 class="page-title">Token配置</h1>
    </div>

    <!-- 当前状态 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">当前状态</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="配置状态">
          <el-tag :type="tokenInfo.configured ? 'success' : 'info'">
            {{ tokenInfo.configured ? '已配置' : '未配置' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="有效性">
          <el-tag :type="tokenInfo.valid ? 'success' : 'danger'">
            {{ tokenInfo.valid ? '有效' : '无效' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="过期时间" :span="2">
          {{ tokenInfo.expiresAt || '未知' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态消息" :span="2">
          {{ tokenInfo.message || '未知' }}
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- Token配置 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">配置Token</h3>
      <el-form :model="form" label-position="top">
        <el-form-item label="Access Token">
          <el-input
            v-model="form.accessToken"
            type="textarea"
            :rows="3"
            placeholder="请输入115 Access Token"
          />
        </el-form-item>
        <el-form-item label="Refresh Token">
          <el-input
            v-model="form.refreshToken"
            type="textarea"
            :rows="3"
            placeholder="请输入115 Refresh Token"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="saving"
            @click="handleSave"
            style="width: 100%"
          >
            保存Token
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 获取方法说明 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">获取Token方法</h3>
      <el-steps direction="vertical" :active="3">
        <el-step title="登录115网页版" description="访问 115.com 并登录您的账号" />
        <el-step title="打开开发者工具" description="按 F12 打开浏览器开发者工具" />
        <el-step title="获取Token" description="在开发者工具中找到 Access Token 和 Refresh Token" />
      </el-steps>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { setToken, getTokenStatus } from '@/api/token'

const saving = ref(false)
const tokenInfo = ref<any>({})
const form = ref({
  accessToken: '',
  refreshToken: ''
})

async function loadTokenStatus() {
  try {
    tokenInfo.value = await getTokenStatus()
  } catch (error) {
    console.error('获取Token状态失败:', error)
  }
}

async function handleSave() {
  if (!form.value.accessToken || !form.value.refreshToken) {
    ElMessage.warning('请填写完整的Token信息')
    return
  }

  saving.value = true
  try {
    await setToken(form.value.accessToken, form.value.refreshToken)
    ElMessage.success('Token保存成功')
    loadTokenStatus()
    form.value.accessToken = ''
    form.value.refreshToken = ''
  } catch (error) {
    ElMessage.error('Token保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadTokenStatus()
})
</script>
