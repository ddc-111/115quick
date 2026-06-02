import { inject } from 'vue';
import type { AxiosInstance } from 'axios';
import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

//const globalConfig = inject("globalConfig") as { apiUrl?: string } | undefined;


// 创建 axios 实例
const instance: AxiosInstance = axios.create({
  baseURL: "http://192.168.31.81:8888", // 避免 undefined
  timeout: 5000,
});

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    const token = authStore.token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
instance.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error("API 请求错误:", error);
    return Promise.reject(error);
  }
);

export default instance;
