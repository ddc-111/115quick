<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Download, Link } from '@element-plus/icons-vue'
import axios from 'axios'
import { Buffer } from 'buffer'
import CryptoJS from 'crypto-js'
import request from '@/utils/request'

const downloadPath = ref('')
const qrCodeUrl = ref('')
const isScanning = ref(false)

const formData = reactive({
  downloadPath: '',
  bindCode: ''
})

const rules = {
  downloadPath: [
    { required: true, message: '请输入下载路径', trigger: 'blur' }
  ],
  bindCode: [
    { required: true, message: '请输入绑定码', trigger: 'blur' }
  ]
}

const formRef = ref()

const handleSavePath = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid: boolean) => {
    if (valid) {
      // TODO: 保存下载路径到本地存储或后端
      ElMessage.success('下载路径保存成功')
    } else {
      ElMessage.error('请填写正确的下载路径')
      return false
    }
  })
}

const startScan = async () => {
  try {
    isScanning.value = true
    const response = await request.get('/v1/Config/api/get115ScanQrCode')

    if (response.data) {
      qrCodeUrl.value ="https://minico.qq.com/qrcode/get?type=2&r=1&size=300&b=miniapp&text="+response.data.qrcodeUrl

    } else {
      ElMessage.error('获取二维码失败')
      stopScan()
    }
  } catch (error) {
    console.error('获取二维码错误:', error)
    ElMessage.error('获取二维码失败')
    stopScan()
  }
}

const stopScan = () => {
  isScanning.value = false
  qrCodeUrl.value = ''
}
</script>

<template>
  <div class="main-container">
    <el-card class="path-card">
      <template #header>
        <div class="card-header">
          <el-icon>
            <Download />
          </el-icon>
          <span>下载路径设置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="下载路径" prop="downloadPath">
          <el-input v-model="formData.downloadPath" placeholder="请输入下载文件保存路径">
            <template #append>
              <el-button @click="handleSavePath">保存</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="bind-card">
      <template #header>
        <div class="card-header">
          <el-icon>
            <Link />
          </el-icon>
          <span>115账号绑定</span>
        </div>
      </template>

      <div class="bind-container">
        <div v-if="!isScanning" class="start-scan">
          <el-button type="primary" @click="startScan">开始扫码绑定</el-button>
        </div>
        <div v-else class="qr-code-container">
          <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="115绑定二维码" class="qr-code" />
          <div class="scan-tips">请使用115手机客户端扫描二维码完成绑定</div>
          <el-button @click="stopScan">取消绑定</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.main-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.path-card,
.bind-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: bold;
}

.bind-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.qr-code-container {
  text-align: center;
}

.qr-code {
  width: 200px;
  height: 200px;
  margin-bottom: 16px;
}

.scan-tips {
  color: #666;
  margin-bottom: 16px;
}

.start-scan {
  padding: 40px 0;
}
</style>
