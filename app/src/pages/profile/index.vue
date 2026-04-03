<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import TabBar from '@/components/TabBar.vue'

const userStore = useUserStore()

// 检查是否是管理员
const isAdmin = computed(() => userStore.userInfo?.role === 'admin')

interface MenuItem {
  label: string
  icon: string
  path: string
}

const menuItems: MenuItem[] = [
  { label: '评论记录', icon: '评', path: '' },
  { label: '举证记录', icon: '证', path: '' },
  { label: '身份认证', icon: '认', path: '' },
  { label: '设置', icon: '设', path: '' },
]

onShow(() => {
  if (userStore.isLoggedIn) {
    userStore.fetchProfile()
  }
})

function goToLogin() {
  uni.navigateTo({ url: '/pages/auth/login' })
}

function goToCoinHistory() {
  if (!userStore.isLoggedIn) {
    goToLogin()
    return
  }
  uni.navigateTo({ url: '/pages/profile/coin-history' })
}

function goToAdmin() {
  uni.navigateTo({ url: '/pages/admin/index' })
}

function handleMenuTap(item: MenuItem) {
  if (!userStore.isLoggedIn) {
    goToLogin()
    return
  }
  if (item.path) {
    uni.navigateTo({ url: item.path })
  } else {
    uni.showToast({ title: '暂未开放', icon: 'none' })
  }
}

function handleLogout() {
  uni.showModal({
    title: '确认退出',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.showToast({ title: '已退出登录', icon: 'none' })
      }
    },
  })
}
</script>

<template>
  <view class="profile-page">
    <!-- User info section -->
    <view class="user-section">
      <view v-if="userStore.isLoggedIn" class="user-info">
        <view class="avatar-wrapper">
          <image
            class="avatar"
            :src="userStore.avatar || '/static/default-avatar.png'"
            mode="aspectFill"
          />
        </view>
        <view class="user-detail">
          <text class="user-nickname">{{ userStore.nickname }}</text>
          <text class="user-phone" v-if="userStore.userInfo?.phone">
            {{ userStore.userInfo.phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') }}
          </text>
        </view>
      </view>
      <view v-else class="login-prompt" @tap="goToLogin">
        <view class="avatar-wrapper">
          <view class="avatar-placeholder">
            <text class="avatar-placeholder-text">?</text>
          </view>
        </view>
        <view class="user-detail">
          <text class="login-text">点击登录</text>
          <text class="login-hint">登录后解锁更多功能</text>
        </view>
      </view>
    </view>

    <!-- Coin section -->
    <view class="coin-section" @tap="goToCoinHistory">
      <view class="coin-info">
        <text class="coin-label">积分余额</text>
        <text class="coin-balance">{{ userStore.isLoggedIn ? userStore.coinBalance : '--' }}</text>
      </view>
      <text class="coin-arrow">&#8250;</text>
    </view>

    <!-- Menu section -->
    <view class="menu-section">
      <view
        v-for="item in menuItems"
        :key="item.label"
        class="menu-item"
        @tap="handleMenuTap(item)"
      >
        <view class="menu-left">
          <view class="menu-icon-box">
            <text class="menu-icon">{{ item.icon }}</text>
          </view>
          <text class="menu-label">{{ item.label }}</text>
        </view>
        <text class="menu-arrow">&#8250;</text>
      </view>
    </view>

    <!-- Admin section -->
    <view v-if="isAdmin" class="admin-section">
      <view class="admin-item" @tap="goToAdmin">
        <view class="menu-left">
          <view class="menu-icon-box admin-icon">
            <text class="menu-icon">管</text>
          </view>
          <text class="menu-label">管理后台</text>
        </view>
        <text class="menu-arrow">&#8250;</text>
      </view>
    </view>

    <!-- Logout -->
    <view v-if="userStore.isLoggedIn" class="logout-section">
      <view class="logout-btn" @tap="handleLogout">
        <text class="logout-text">退出登录</text>
      </view>
    </view>

    <!-- TabBar -->
    <TabBar :current="3" />
  </view>
</template>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background-color: $color-bg-secondary;
  padding-bottom: 120rpx;
}

.user-section {
  background-color: $color-text;
  padding: $spacing-xl $spacing-lg $spacing-lg;
}

.user-info,
.login-prompt {
  display: flex;
  align-items: center;
}

.avatar-wrapper {
  margin-right: $spacing-lg;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background-color: $color-bg-secondary;
}

.avatar-placeholder {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-placeholder-text {
  font-size: $font-size-xxl;
  color: rgba(255, 255, 255, 0.5);
}

.user-detail {
  display: flex;
  flex-direction: column;
}

.user-nickname {
  font-size: $font-size-xl;
  font-weight: 600;
  color: $color-bg;
  margin-bottom: $spacing-xs;
}

.user-phone {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
}

.login-text {
  font-size: $font-size-xl;
  font-weight: 600;
  color: $color-bg;
  margin-bottom: $spacing-xs;
}

.login-hint {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
}

.coin-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: $color-bg;
  padding: $spacing-lg;
  margin-bottom: 2rpx;
}

.coin-info {
  display: flex;
  flex-direction: column;
}

.coin-label {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  margin-bottom: $spacing-xs;
}

.coin-balance {
  font-size: 56rpx;
  font-weight: 700;
  color: $color-text;
  line-height: 1;
}

.coin-arrow {
  font-size: $font-size-xl;
  color: $color-text-tertiary;
}

.menu-section {
  background-color: $color-bg;
  margin-top: $spacing-sm;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1rpx solid $color-border;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background-color: $color-bg-secondary;
  }
}

.menu-left {
  display: flex;
  align-items: center;
}

.menu-icon-box {
  width: 48rpx;
  height: 48rpx;
  border-radius: $radius-sm;
  background-color: $color-bg-secondary;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: $spacing-md;
}

.menu-icon {
  font-size: $font-size-sm;
  color: $color-text;
  font-weight: 600;
}

.menu-label {
  font-size: $font-size-base;
  color: $color-text;
}

.menu-arrow {
  font-size: $font-size-xl;
  color: $color-text-tertiary;
}

.logout-section {
  margin-top: $spacing-xl;
  padding: 0 $spacing-lg;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80rpx;
  background-color: $color-bg;
  border-radius: $radius-sm;
}

.logout-text {
  font-size: $font-size-base;
  color: $color-text-secondary;
}

.admin-section {
  margin-top: $spacing-sm;
  background-color: $color-bg;
}

.admin-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  background-color: $color-bg;
}

.admin-icon {
  background-color: $color-text;
  .menu-icon {
    color: $color-bg;
  }
}

@media (prefers-color-scheme: dark) {
  .profile-page {
    background-color: $color-dark-bg;
  }

  .coin-section {
    background-color: $color-dark-bg-secondary;
  }

  .coin-balance {
    color: $color-dark-text;
  }

  .menu-section {
    background-color: $color-dark-bg-secondary;
  }

  .menu-item {
    border-bottom-color: $color-dark-border;

    &:active {
      background-color: $color-dark-bg;
    }
  }

  .menu-icon-box {
    background-color: $color-dark-bg;
  }

  .menu-icon {
    color: $color-dark-text;
  }

  .menu-label {
    color: $color-dark-text;
  }

  .logout-btn {
    background-color: $color-dark-bg-secondary;
  }

  .logout-text {
    color: $color-dark-text-secondary;
  }
}
</style>
