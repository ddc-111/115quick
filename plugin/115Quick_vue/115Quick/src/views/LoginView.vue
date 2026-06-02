<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { checkIsFirst, login } from '@/api/user'
import { useAuthStore } from '@/stores/auth' // 引入 store
const router = useRouter()
const loading = ref(false)
const authStore = useAuthStore()
const loginForm = reactive({
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

onMounted(async () => {
  try {
    const { data } = await checkIsFirst()
    if (data.isFirst) {
      router.push('/first-set-auth')
    }
  } catch (error) {
    ElMessage.error('系统检查失败，请刷新重试')
  }
})

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        loading.value = true
        const { data } = await login(loginForm)
        console.log(data)
        authStore.setToken(data.accessToken)
        ElMessage.success('登录成功')
        router.push('/main')
      } catch (error) {
        ElMessage.error('登录失败，请检查用户名和密码')
      } finally {
        loading.value = false
      }
    } else {
      ElMessage.error('请填写正确的登录信息')
      return false
    }
  })
}
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <h2>系统登录</h2>
      <el-form ref="formRef" :model="loginForm" :rules="rules" label-width="0" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="用户名">
            <template #prefix>
              <el-icon>
                <User />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码" show-password>
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1a237e 0%, #0d47a1 100%);
}

.login-box {
  width: 480px;
  padding: 50px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
}

.login-box h2 {
  text-align: center;
  margin-bottom: 40px;
  color: #1a237e;
  font-size: 28px;
  font-weight: 600;
}

.login-form {
  margin-top: 30px;
}

.login-form :deep(.el-input__wrapper) {
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.login-form :deep(.el-input__inner) {
  height: 42px;
  font-size: 16px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  margin-top: 10px;
  background: linear-gradient(135deg, #1a237e 0%, #0d47a1 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(26, 35, 126, 0.2);
}

.login-form :deep(.el-form-item) {
  margin-bottom: 25px;
}

.login-form :deep(.el-icon) {
  font-size: 18px;
  color: #1a237e;
}
</style>
