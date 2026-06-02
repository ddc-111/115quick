import axios from 'axios'
import { useServerStore } from '@/stores/server'
import { ElMessage } from 'element-plus'

const request = axios.create({
  timeout: 30000
})

request.interceptors.request.use(
  (config) => {
    const serverStore = useServerStore()
    if (!serverStore.serverUrl) {
      return Promise.reject(new Error('未配置服务器地址'))
    }
    config.baseURL = serverStore.getBaseUrl()
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code === 200) {
      return data.data
    }
    ElMessage.error(data.Msg || '请求失败')
    return Promise.reject(new Error(data.Msg || '请求失败'))
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request
