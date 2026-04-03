import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authService from '@/services/auth'

interface UserInfo {
  id: string
  phone: string
  nickname: string
  avatar: string
  coin_balance: number
  is_verified: boolean
  role: string
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(authService.getCachedUser())
  const token = ref<string>(authService.getToken())

  const isLoggedIn = computed(() => !!token.value)
  const nickname = computed(() => userInfo.value?.nickname || '未登录')
  const avatar = computed(() => userInfo.value?.avatar || '')
  const coinBalance = computed(() => userInfo.value?.coin_balance || 0)

  async function login(phone: string, code: string) {
    const result = await authService.loginSMS(phone, code)
    token.value = result.token
    userInfo.value = result.user
  }

  function logout() {
    authService.logout()
    token.value = ''
    userInfo.value = null
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const user = await authService.fetchProfile()
      userInfo.value = user
    } catch {
      // Token may be invalid
    }
  }

  function updateCoinBalance(amount: number) {
    if (userInfo.value) {
      userInfo.value.coin_balance += amount
    }
  }

  return {
    userInfo,
    token,
    isLoggedIn,
    nickname,
    avatar,
    coinBalance,
    login,
    logout,
    fetchProfile,
    updateCoinBalance,
  }
})
