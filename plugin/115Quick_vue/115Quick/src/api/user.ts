import request from '@/utils/request'
import axios from 'axios'

export interface LoginForm {
  username: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  accessExpire: number
  refreshAfter: number
}

export const checkIsFirst = () => {
  return request.get<{ isFirst: boolean }>('v1/User/api/isFirst')
}

export const firstSetAuth = (data: LoginForm) => {
  return request.post('v1/User/api/firstSetAuth', data)
}

export const login = (data: LoginForm) => {
  return axios.post<LoginResponse>('http://192.168.31.81:8888/v1/User/api/login', data, {
    headers: {
      'Content-Type': 'application/json'
    },
    timeout: 5000
  })
}
