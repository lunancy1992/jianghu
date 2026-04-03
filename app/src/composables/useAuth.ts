import { useUserStore } from '@/stores/user'

/**
 * Auth composable for checking login status and redirecting.
 */
export function useAuth() {
  const userStore = useUserStore()

  function requireLogin(callback?: () => void): boolean {
    if (!userStore.isLoggedIn) {
      uni.navigateTo({
        url: '/pages/auth/login',
      })
      return false
    }
    callback?.()
    return true
  }

  function checkLogin(): boolean {
    return userStore.isLoggedIn
  }

  return {
    isLoggedIn: userStore.isLoggedIn,
    userInfo: userStore.userInfo,
    requireLogin,
    checkLogin,
  }
}
