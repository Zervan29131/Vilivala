import service from './index'

// 用户注册
export const userRegister = (data) => {
  return service({
    url: '/v1/user/register',
    method: 'post',
    data
  })
}

// 用户登录
export const userLogin = (data) => {
  return service({
    url: '/v1/user/login',
    method: 'post',
    data
  })
}

// 获取当前用户信息
export const getUserInfo = () => {
  return service({
    url: '/v1/user/info',
    method: 'get'
  })
}

// 修改密码
export const changePassword = (data) => {
  return service({
    url: '/v1/user/password',
    method: 'put',
    data
  })
}