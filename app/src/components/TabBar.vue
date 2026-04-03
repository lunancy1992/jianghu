<script setup lang="ts">
import { ref } from 'vue'

interface TabItem {
  pagePath: string
  text: string
  iconText: string
}

const tabs: TabItem[] = [
  { pagePath: '/pages/headlines/index', text: '头条', iconText: '报' },
  { pagePath: '/pages/news/index', text: '资讯', iconText: '讯' },
  { pagePath: '/pages/events/index', text: '追踪', iconText: '踪' },
  { pagePath: '/pages/profile/index', text: '我的', iconText: '我' },
]

const props = defineProps<{
  current: number
}>()

function switchTab(index: number) {
  if (index === props.current) return
  uni.switchTab({
    url: tabs[index].pagePath,
  })
}
</script>

<template>
  <view class="custom-tabbar safe-bottom">
    <view
      v-for="(tab, index) in tabs"
      :key="tab.pagePath"
      class="tabbar-item"
      :class="{ 'tabbar-item--active': current === index }"
      @tap="switchTab(index)"
    >
      <text class="tabbar-icon">{{ tab.iconText }}</text>
      <text class="tabbar-label">{{ tab.text }}</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.custom-tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  height: 100rpx;
  background-color: $color-bg;
  border-top: 1rpx solid $color-border;
  z-index: 999;
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xs 0;

  &--active {
    .tabbar-icon {
      color: $color-text;
      font-weight: 700;
    }

    .tabbar-label {
      color: $color-text;
      font-weight: 600;
    }
  }
}

.tabbar-icon {
  font-size: $font-size-lg;
  color: $color-text-tertiary;
  font-weight: 500;
  margin-bottom: 2rpx;
}

.tabbar-label {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

@media (prefers-color-scheme: dark) {
  .custom-tabbar {
    background-color: $color-dark-bg;
    border-top-color: $color-dark-border;
  }

  .tabbar-item--active {
    .tabbar-icon {
      color: $color-dark-text;
    }

    .tabbar-label {
      color: $color-dark-text;
    }
  }
}
</style>
