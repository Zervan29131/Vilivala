import axios from 'axios'

// 请求拦截器：统一携带token
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 登录接口
export const login = (data) => {
  return axios.post('/api/v1/user/login', data).then(res => res.data)
}

// 注册接口
export const register = (data) => {
  return axios.post('/api/v1/user/register', data).then(res => res.data)
}

// 获取用户信息
export const getUserInfo = () => {
  return axios.get('/api/v1/user/info').then(res => res.data)
}

// 修改密码
export const changePassword = (data) => {
  return axios.put('/api/v1/user/password', data).then(res => res.data)
}