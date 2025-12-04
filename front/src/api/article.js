import service from './index'

// 获取文章列表（分页/筛选）
export const getArticleList = (params) => {
  return service({
    url: '/v1/article/list',
    method: 'get',
    params
  })
}

// 获取文章详情
export const getArticleDetail = (id) => {
  return service({
    url: `/v1/article/${id}`,
    method: 'get'
  })
}

// 创建文章
export const createArticle = (data) => {
  return service({
    url: '/v1/article',
    method: 'post',
    data
  })
}

// 更新文章
export const updateArticle = (id, data) => {
  return service({
    url: `/v1/article/${id}`,
    method: 'put',
    data
  })
}

// 删除文章
export const deleteArticle = (id) => {
  return service({
    url: `/v1/article/${id}`,
    method: 'delete'
  })
}