export const API_BASE_URL = '/api/v1'

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: Record<string, any>
  header?: Record<string, string>
  showError?: boolean
}

function getToken(): string {
  return uni.getStorageSync('token') || ''
}

function setToken(token: string): void {
  uni.setStorageSync('token', token)
}

async function refreshToken(): Promise<boolean> {
  const currentToken = getToken()
  if (!currentToken) return false

  try {
    const res: any = await new Promise((resolve, reject) => {
      uni.request({
        url: `${API_BASE_URL}/auth/refresh`,
        method: 'POST',
        data: { token: currentToken },
        success: resolve,
        fail: reject,
      })
    })

    if (res.statusCode === 200 && res.data?.code === 0) {
      setToken(res.data.data.token)
      return true
    }
    return false
  } catch {
    return false
  }
}

let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

function onTokenRefreshed(token: string) {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

function addRefreshSubscriber(cb: (token: string) => void) {
  refreshSubscribers.push(cb)
}

export async function request<T = any>(options: RequestOptions): Promise<T> {
  const { url, method = 'GET', data, header = {}, showError = true } = options

  const token = getToken()
  if (token) {
    header['Authorization'] = `Bearer ${token}`
  }
  header['Content-Type'] = header['Content-Type'] || 'application/json'

  return new Promise((resolve, reject) => {
    uni.request({
      url: url.startsWith('http') ? url : `${API_BASE_URL}${url}`,
      method,
      data,
      header,
      success: async (res: any) => {
        if (res.statusCode === 401) {
          if (!isRefreshing) {
            isRefreshing = true
            const refreshed = await refreshToken()
            isRefreshing = false

            if (refreshed) {
              const newToken = getToken()
              onTokenRefreshed(newToken)
              // Retry original request
              header['Authorization'] = `Bearer ${newToken}`
              uni.request({
                url: url.startsWith('http') ? url : `${API_BASE_URL}${url}`,
                method,
                data,
                header,
                success: (retryRes: any) => {
                  if (retryRes.data?.code === 0) {
                    resolve(retryRes.data.data)
                  } else {
                    reject(retryRes.data)
                  }
                },
                fail: reject,
              })
            } else {
              uni.removeStorageSync('token')
              uni.removeStorageSync('user_info')
              uni.navigateTo({ url: '/pages/auth/login' })
              reject(new Error('登录已过期'))
            }
          } else {
            // Wait for token refresh
            return new Promise<void>((retryResolve) => {
              addRefreshSubscriber((newToken: string) => {
                header['Authorization'] = `Bearer ${newToken}`
                uni.request({
                  url: url.startsWith('http') ? url : `${API_BASE_URL}${url}`,
                  method,
                  data,
                  header,
                  success: (retryRes: any) => {
                    if (retryRes.data?.code === 0) {
                      resolve(retryRes.data.data)
                    } else {
                      reject(retryRes.data)
                    }
                    retryResolve()
                  },
                  fail: (err) => {
                    reject(err)
                    retryResolve()
                  },
                })
              })
            })
          }
          return
        }

        const body = res.data as ApiResponse<T>
        if (body.code === 0) {
          resolve(body.data)
        } else {
          if (showError) {
            uni.showToast({ title: body.message || '请求失败', icon: 'none' })
          }
          reject(body)
        }
      },
      fail: (err) => {
        if (showError) {
          uni.showToast({ title: '网络错误', icon: 'none' })
        }
        reject(err)
      },
    })
  })
}

export function get<T = any>(url: string, data?: Record<string, any>) {
  return request<T>({ url, method: 'GET', data })
}

export function post<T = any>(url: string, data?: Record<string, any>) {
  return request<T>({ url, method: 'POST', data })
}

export function put<T = any>(url: string, data?: Record<string, any>) {
  return request<T>({ url, method: 'PUT', data })
}

export function del<T = any>(url: string, data?: Record<string, any>) {
  return request<T>({ url, method: 'DELETE', data })
}
