<script setup lang="ts">
import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app'
import { ref } from 'vue'
import NewsCard from '@/components/NewsCard.vue'
import TabBar from '@/components/TabBar.vue'
import { useNewsStore } from '@/stores/news'
import { formatHeadlineDate } from '@/utils/format'
import { post } from '@/services/api'
import { useUserStore } from '@/stores/user'

const newsStore = useNewsStore()
const userStore = useUserStore()
const todayDate = ref(formatHeadlineDate())
const hasRead = ref(false)
const currentHistoryDate = ref('')

async function loadHeadlines() {
  await newsStore.fetchHeadlines(currentHistoryDate.value || undefined)
}

onLoad(() => {
  loadHeadlines()
})

onPullDownRefresh(() => {
  hasRead.value = false
  currentHistoryDate.value = ''
  todayDate.value = formatHeadlineDate()
  loadHeadlines().finally(() => {
    uni.stopPullDownRefresh()
  })
})

async function handleRead(id: string) {
  newsStore.markHeadlineRead(id)
  if (!userStore.isLoggedIn) return

  try {
    await post(`/news/${id}/read`)
  } catch {
    // ignore read mark failure
  }
}

function handleOpenSource(_url: string) {}

async function markAllRead() {
  if (!hasRead.value) {
    newsStore.markAllHeadlinesRead()
    hasRead.value = true
    currentHistoryDate.value = new Date().toISOString().slice(0, 10)
    todayDate.value = `${currentHistoryDate.value} 历史头条`
    await loadHeadlines()
    uni.showToast({ title: '已切换历史头条', icon: 'none' })
    return
  }

  await loadHeadlines()
}
</script>

<template>
  <view class="headlines-page">
    <!-- Header -->
    <view class="headlines-header">
      <view class="header-content">
        <text class="header-title">江湖小报</text>
        <text class="header-subtitle">每日要闻 · 事实为先</text>
      </view>
      <text class="header-date">{{ todayDate }}</text>
    </view>

    <!-- Loading state -->
    <view v-if="newsStore.isLoadingHeadlines && newsStore.headlines.length === 0" class="loading-text">
      <text>加载中...</text>
    </view>

    <!-- Headlines list -->
    <view class="headlines-list">
      <view v-for="(item, index) in newsStore.headlines" :key="item.id" class="headline-item">
        <view class="headline-index">
          <text class="index-text" :class="{ 'index-top': index < 3 }">{{ index + 1 }}</text>
        </view>
        <view class="headline-content">
          <NewsCard
            :news="{
              id: item.news_id || item.id,
              title: item.title,
              source: item.source,
              published_at: item.published_at,
              source_url: item.source_url,
              veracity: item.veracity,
              comment_text: item.comment_text,
              is_read: item.is_read,
            }"
            @read="handleRead"
            @open-source="handleOpenSource"
          />
        </view>
      </view>
    </view>

    <!-- Empty state -->
    <view v-if="!newsStore.isLoadingHeadlines && newsStore.headlines.length === 0" class="empty-state">
      <text>{{ hasRead ? '当天暂无历史头条' : '暂无头条' }}</text>
    </view>

    <!-- Bottom action -->
    <view v-if="newsStore.headlines.length > 0" class="headlines-footer">
      <view class="footer-btn" @tap="markAllRead">
        <text class="footer-btn-text">{{ hasRead ? '已阅 · 查看历史' : '标记已阅' }}</text>
      </view>
    </view>

    <!-- TabBar -->
    <TabBar :current="0" />
  </view>
</template>

<style lang="scss" scoped>
.headlines-page {
  min-height: 100vh;
  background-color: $color-bg;
  padding-bottom: 160rpx;
}

.headlines-header {
  background-color: $color-text;
  padding: $spacing-xl $spacing-lg $spacing-lg;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.header-content {
  display: flex;
  flex-direction: column;
}

.header-title {
  font-size: $font-size-title;
  font-weight: 700;
  color: $color-bg;
  letter-spacing: 4rpx;
}

.header-subtitle {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
  margin-top: $spacing-xs;
  letter-spacing: 2rpx;
}

.header-date {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
}

.headlines-list {
  background-color: $color-bg;
}

.headline-item {
  display: flex;
}

.headline-index {
  width: 60rpx;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: $spacing-lg;
  flex-shrink: 0;
}

.index-text {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
  font-weight: 600;
  font-style: italic;
}

.index-top {
  color: $color-text;
  font-size: $font-size-md;
  font-weight: 700;
}

.headline-content {
  flex: 1;
  min-width: 0;
}

.headlines-footer {
  padding: $spacing-xl $spacing-lg;
  display: flex;
  justify-content: center;
}

.footer-btn {
  padding: $spacing-sm $spacing-xl;
  border: 1rpx solid $color-border;
  border-radius: $radius-sm;
}

.footer-btn-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

@media (prefers-color-scheme: dark) {
  .headlines-page {
    background-color: $color-dark-bg;
  }

  .headlines-list {
    background-color: $color-dark-bg;
  }

  .index-top {
    color: $color-dark-text;
  }

  .footer-btn {
    border-color: $color-dark-border;
  }

  .footer-btn-text {
    color: $color-dark-text-secondary;
  }
}
</style>
