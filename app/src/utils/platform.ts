/**
 * Platform detection helpers for uni-app.
 */

export function isH5(): boolean {
  // #ifdef H5
  return true
  // #endif
  // #ifndef H5
  return false
  // #endif
}

export function isWeixin(): boolean {
  // #ifdef MP-WEIXIN
  return true
  // #endif
  // #ifndef MP-WEIXIN
  return false
  // #endif
}

export function isToutiao(): boolean {
  // #ifdef MP-TOUTIAO
  return true
  // #endif
  // #ifndef MP-TOUTIAO
  return false
  // #endif
}

export function isApp(): boolean {
  // #ifdef APP-PLUS
  return true
  // #endif
  // #ifndef APP-PLUS
  return false
  // #endif
}

export function getStatusBarHeight(): number {
  const systemInfo = uni.getSystemInfoSync()
  return systemInfo.statusBarHeight || 0
}

/**
 * Open a URL. Uses webview on mini-programs, window.open on H5.
 */
export function openUrl(url: string): void {
  // #ifdef H5
  window.open(url, '_blank')
  // #endif
  // #ifndef H5
  uni.navigateTo({
    url: `/pages/webview/index?url=${encodeURIComponent(url)}`,
    fail: () => {
      uni.setClipboardData({
        data: url,
        success: () => {
          uni.showToast({ title: '链接已复制', icon: 'none' })
        },
      })
    },
  })
  // #endif
}
