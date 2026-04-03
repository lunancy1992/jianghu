<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { sendSMS } from '@/services/auth'

const userStore = useUserStore()

const phone = ref('')
const smsCode = ref('')
const countdown = ref(0)
const isSending = ref(false)
const isLogging = ref(false)

let timer: ReturnType<typeof setInterval> | null = null

function isPhoneValid(): boolean {
  return /^1\d{10}$/.test(phone.value)
}

async function handleSendSMS() {
  if (!isPhoneValid()) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  if (countdown.value > 0 || isSending.value) return

  isSending.value = true
  try {
    await sendSMS(phone.value)
    uni.showToast({ title: '验证码已发送', icon: 'none' })
    startCountdown()
  } catch {
    uni.showToast({ title: '发送失败，请重试', icon: 'none' })
  } finally {
    isSending.value = false
  }
}

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (timer) clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function handleLogin() {
  if (!isPhoneValid()) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  if (!smsCode.value || smsCode.value.length < 4) {
    uni.showToast({ title: '请输入验证码', icon: 'none' })
    return
  }
  if (isLogging.value) return

  isLogging.value = true
  try {
    await userStore.login(phone.value, smsCode.value)
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => {
      uni.switchTab({ url: '/pages/profile/index' })
    }, 1000)
  } catch {
    uni.showToast({ title: '登录失败，请检查验证码', icon: 'none' })
  } finally {
    isLogging.value = false
  }
}
</script>

<template>
  <view class="login-page">
    <view class="login-header">
      <text class="login-title">江湖小报</text>
      <text class="login-subtitle">手机号登录</text>
    </view>

    <view class="login-form">
      <!-- Phone input -->
      <view class="form-item">
        <view class="input-wrapper">
          <text class="input-prefix">+86</text>
          <input
            v-model="phone"
            class="form-input"
            type="number"
            placeholder="请输入手机号"
            :maxlength="11"
          />
        </view>
      </view>

      <!-- SMS code input -->
      <view class="form-item">
        <view class="input-wrapper">
          <input
            v-model="smsCode"
            class="form-input sms-input"
            type="number"
            placeholder="请输入验证码"
            :maxlength="6"
          />
          <view
            class="sms-btn"
            :class="{ 'sms-btn--disabled': countdown > 0 || !isPhoneValid() }"
            @tap="handleSendSMS"
          >
            <text class="sms-btn-text">
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </text>
          </view>
        </view>
      </view>

      <!-- Login button -->
      <view
        class="login-btn"
        :class="{ 'login-btn--disabled': !isPhoneValid() || smsCode.length < 4 || isLogging }"
        @tap="handleLogin"
      >
        <text class="login-btn-text">{{ isLogging ? '登录中...' : '登录' }}</text>
      </view>

      <view class="login-agreement">
        <text class="agreement-text">
          登录即表示同意《用户协议》和《隐私政策》
        </text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background-color: $color-bg;
  padding: 0 $spacing-lg;
}

.login-header {
  padding-top: 120rpx;
  padding-bottom: $spacing-xl;
}

.login-title {
  font-size: $font-size-title;
  font-weight: 700;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-sm;
  letter-spacing: 4rpx;
}

.login-subtitle {
  font-size: $font-size-md;
  color: $color-text-secondary;
  display: block;
}

.login-form {
  padding-top: $spacing-lg;
}

.form-item {
  margin-bottom: $spacing-lg;
}

.input-wrapper {
  display: flex;
  align-items: center;
  height: 96rpx;
  border-bottom: 2rpx solid $color-border;
}

.input-prefix {
  font-size: $font-size-lg;
  color: $color-text;
  font-weight: 500;
  margin-right: $spacing-md;
  flex-shrink: 0;
}

.form-input {
  flex: 1;
  height: 96rpx;
  font-size: $font-size-lg;
  color: $color-text;
}

.sms-input {
  flex: 1;
}

.sms-btn {
  flex-shrink: 0;
  padding: $spacing-xs $spacing-md;
  border-left: 1rpx solid $color-border;
  margin-left: $spacing-md;

  &--disabled {
    opacity: 0.4;
    pointer-events: none;
  }
}

.sms-btn-text {
  font-size: $font-size-sm;
  color: $color-text;
  font-weight: 500;
  white-space: nowrap;
}

.login-btn {
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: $color-text;
  border-radius: $radius-sm;
  margin-top: $spacing-xl;

  &--disabled {
    opacity: 0.4;
    pointer-events: none;
  }
}

.login-btn-text {
  font-size: $font-size-lg;
  color: $color-bg;
  font-weight: 600;
}

.login-agreement {
  margin-top: $spacing-xl;
  display: flex;
  justify-content: center;
}

.agreement-text {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

@media (prefers-color-scheme: dark) {
  .login-page {
    background-color: $color-dark-bg;
  }

  .login-title {
    color: $color-dark-text;
  }

  .login-subtitle {
    color: $color-dark-text-secondary;
  }

  .input-wrapper {
    border-bottom-color: $color-dark-border;
  }

  .input-prefix {
    color: $color-dark-text;
  }

  .form-input {
    color: $color-dark-text;
  }

  .sms-btn {
    border-left-color: $color-dark-border;
  }

  .sms-btn-text {
    color: $color-dark-text;
  }

  .login-btn {
    background-color: $color-dark-text;
  }

  .login-btn-text {
    color: $color-dark-bg;
  }
}
</style>
