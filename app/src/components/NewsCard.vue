<script setup lang="ts">
import { formatTime } from '@/utils/format'
import VeracityBadge from './VeracityBadge.vue'

export interface NewsCardData {
  id: string
  title: string
  source: string
  published_at: string
  source_url: string
  veracity: number
  comment_text: string
  is_read?: boolean
}

const props = defineProps<{
  news: NewsCardData
}>()

const emit = defineEmits<{
  (e: 'read', id: string): void
  (e: 'open-source', url: string): void
}>()

function handleTap() {
  emit('read', props.news.id)
  uni.navigateTo({
    url: `/pages/news/detail?id=${props.news.id}`,
  })
}

function handleOpenSource() {
  emit('open-source', props.news.source_url)
  // #ifdef H5
  window.open(props.news.source_url, '_blank')
  // #endif
  // #ifndef H5
  uni.setClipboardData({
    data: props.news.source_url,
    success: () => {
      uni.showToast({ title: '链接已复制', icon: 'none' })
    },
  })
  // #endif
}
</script>

<template>
  <view class="news-card" @tap="handleTap">
    <!-- Meta row: time + source + veracity -->
    <view class="card-meta">
      <view class="meta-left">
        <text class="meta-time">{{ formatTime(news.published_at) }}</text>
        <text class="meta-divider">|</text>
        <text class="meta-source">{{ news.source }}</text>
      </view>
      <VeracityBadge :veracity="news.veracity" />
    </view>

    <!-- Title -->
    <view class="card-title">
      <text class="title-text">{{ news.title }}</text>
      <text v-if="news.is_read" class="read-tag">已阅</text>
    </view>

    <!-- Editor comment -->
    <view v-if="news.comment_text" class="card-comment">
      <text class="comment-text">{{ news.comment_text }}</text>
    </view>

    <!-- Source link -->
    <view class="card-footer">
      <text class="source-link" @tap.stop="handleOpenSource">查看原文</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.news-card {
  background-color: $color-bg;
  padding: $spacing-lg $spacing-lg $spacing-md;
  border-bottom: 1rpx solid $color-border;

  &:active {
    background-color: $color-bg-secondary;
  }
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}

.meta-left {
  display: flex;
  align-items: center;
  flex: 1;
  overflow: hidden;
}

.meta-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
  flex-shrink: 0;
}

.meta-divider {
  font-size: $font-size-xs;
  color: $color-border;
  margin: 0 $spacing-xs;
  flex-shrink: 0;
}

.meta-source {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-title {
  margin-bottom: $spacing-sm;
}

.read-tag {
  display: inline-block;
  margin-top: $spacing-xs;
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.title-text {
  font-size: $font-size-lg;
  font-weight: 600;
  color: $color-text;
  line-height: $line-height-tight;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-comment {
  margin-bottom: $spacing-sm;
  padding: $spacing-sm $spacing-md;
  background-color: $color-bg-secondary;
  border-radius: $radius-sm;
}

.comment-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-height-relaxed;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-footer {
  display: flex;
  align-items: center;
}

.source-link {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
  text-decoration: underline;
}

@media (prefers-color-scheme: dark) {
  .news-card {
    background-color: $color-dark-bg-secondary;
    border-bottom-color: $color-dark-border;

    &:active {
      background-color: $color-dark-bg;
    }
  }

  .title-text {
    color: $color-dark-text;
  }

  .card-comment {
    background-color: $color-dark-bg;
  }

  .comment-text {
    color: $color-dark-text-secondary;
  }
}
</style>
