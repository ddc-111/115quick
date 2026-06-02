<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { firstSetAuth } from '@/api/user'

const router = useRouter()

const authForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const formRef = ref()
const loading = ref(false)

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        loading.value = true
        await firstSetAuth(authForm)
        ElMessage.success('账户设置成功，请登录')
        router.push('/login')
      } catch (error) {
        ElMessage.error('设置失败，请重试')
      } finally {
        loading.value = false
      }
    } else {
      ElMessage.error('请填写正确的信息')
      return false
    }
  })
}
</script>

<template>
  <div class="first-auth-container">
    <div class="first-auth-box">
      <h2>首次使用设置</h2>
      <p class="subtitle">请设置您的管理员账户</p>
      <el-form ref="formRef" :model="authForm" :rules="rules" label-width="0" class="auth-form">
        <el-form-item prop="username">
          <el-input v-model="authForm.username" placeholder="用户名">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="authForm.password" type="password" placeholder="密码" show-password>
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="submit-button" @click="handleSubmit">
            确认设置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.first-auth-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1a237e 0%, #0d47a1 100%);
}

.first-auth-box {
  width: 480px;
  padding: 50px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
}

.first-auth-box h2 {
  text-align: center;
  margin-bottom: 10px;
  color: #1a237e;
  font-size: 28px;
  font-weight: 600;
}

.subtitle {
  text-align: center;
  color: #666;
  margin-bottom: 30px;
}

.auth-form {
  margin-top: 30px;
}

.auth-form :deep(.el-input__wrapper) {
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.auth-form :deep(.el-input__inner) {
  height: 42px;
  font-size: 16px;
}

.submit-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  margin-top: 10px;
  background: linear-gradient(135deg, #1a237e 0%, #0d47a1 100%);
  border: none;
  transition: all 0.3s ease;
}

.submit-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(26, 35, 126, 0.2);
}

.auth-form :deep(.el-form-item) {
  margin-bottom: 25px;
}

.auth-form :deep(.el-icon) {
  font-size: 18px;
  color: #1a237e;
}
</style>
