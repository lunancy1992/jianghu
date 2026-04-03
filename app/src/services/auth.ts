import { post, get } from './api'

interface LoginResult {
  token: string
  user: UserInfo
}

interface UserInfo {
  id: string
  phone: string
  nickname: string
  avatar: string
  coin_balance: number
  is_verified: boolean
  role: string
}

export function getToken(): string {
  return uni.getStorageSync('token') || ''
}

export function isLoggedIn(): boolean {
  return !!getToken()
}

export async function sendSMS(phone: string): Promise<void> {
  await post('/auth/sms/send', { phone })
}

export async function loginSMS(phone: string, code: string): Promise<LoginResult> {
  const result = await post<LoginResult>('/auth/sms/login', { phone, code })
  uni.setStorageSync('token', result.token)
  uni.setStorageSync('user_info', JSON.stringify(result.user))
  return result
}

export async function fetchProfile(): Promise<UserInfo> {
  const user = await get<UserInfo>('/user/profile')
  uni.setStorageSync('user_info', JSON.stringify(user))
  return user
}

export function logout(): void {
  uni.removeStorageSync('token')
  uni.removeStorageSync('user_info')
}

export function getCachedUser(): UserInfo | null {
  const raw = uni.getStorageSync('user_info')
  if (!raw) return null
  try {
    return JSON.parse(raw) as UserInfo
  } catch {
    return null
  }
}
