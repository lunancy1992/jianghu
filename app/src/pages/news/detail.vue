<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import VeracityBadge from '@/components/VeracityBadge.vue'
import CommentList from '@/components/CommentList.vue'
import type { Comment } from '@/components/CommentList.vue'
import { get } from '@/services/api'
import { formatTime } from '@/utils/format'
import { useUserStore } from '@/stores/user'

interface NewsDetail {
  id: string
  title: string
  summary: string
  source: string
  source_url: string
  published_at: string
  veracity: number
  comment_text: string
  comment_count: number
}

const userStore = useUserStore()
const newsDetail = ref<NewsDetail | null>(null)
const comments = ref<Comment[]>([])
const pendingCommentNotice = ref('')
const isLoading = ref(true)
const newsId = ref('')

onLoad((query) => {
  if (query?.id) {
    newsId.value = query.id
    loadDetail(query.id)
    loadComments(query.id)
  }
})

async function loadDetail(id: string) {
  isLoading.value = true
  try {
    newsDetail.value = await get<NewsDetail>(`/news/${id}`)
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    isLoading.value = false
  }
}

async function loadComments(id: string) {
  try {
    const result = await get<{ list: Comment[] }>(`/news/${id}/comments`)
    comments.value = result.list || []
  } catch {
    // Silently fail
  }
}

function handleOpenSource() {
  if (!newsDetail.value?.source_url) return
  // #ifdef H5
  window.open(newsDetail.value.source_url, '_blank')
  // #endif
  // #ifndef H5
  uni.setClipboardData({
    data: newsDetail.value.source_url,
    success: () => {
      uni.showToast({ title: '链接已复制', icon: 'none' })
    },
  })
  // #endif
}

function handleCommentPosted(payload: { status: number }) {
  loadComments(newsId.value)
  userStore.updateCoinBalance(-1)
  pendingCommentNotice.value = payload.status === 1 ? '' : '评论已提交审核，通过后会显示在评论区'
}
</script>

<template>
  <view class="detail-page">
    <view v-if="isLoading" class="loading-text">
      <text>加载中...</text>
    </view>

    <view v-if="newsDetail" class="detail-content">
      <!-- Header -->
      <view class="detail-header">
        <view class="detail-meta">
          <text class="meta-source">{{ newsDetail.source }}</text>
          <text class="meta-time">{{ formatTime(newsDetail.published_at) }}</text>
          <VeracityBadge :veracity="newsDetail.veracity" />
        </view>
        <text class="detail-title">{{ newsDetail.title }}</text>
      </view>

      <!-- Summary -->
      <view class="detail-summary">
        <text class="summary-text">{{ newsDetail.summary }}</text>
      </view>

      <!-- Editor comment -->
      <view v-if="newsDetail.comment_text" class="editor-comment">
        <view class="comment-label">
          <text class="label-text">编辑点评</text>
        </view>
        <text class="comment-body">{{ newsDetail.comment_text }}</text>
      </view>

      <!-- Source link -->
      <view class="source-action" @tap="handleOpenSource">
        <text class="source-btn-text">查看原文</text>
        <text class="source-arrow">&#8594;</text>
      </view>

      <!-- Divider -->
      <view class="divider-thick" />

      <!-- Coin hint -->
      <view v-if="userStore.isLoggedIn" class="coin-hint">
        <text class="coin-text">积分余额: {{ userStore.coinBalance }} · 评论消耗1积分</text>
      </view>
      <view v-if="pendingCommentNotice" class="comment-notice">
        <text class="comment-notice-text">{{ pendingCommentNotice }}</text>
      </view>

      <!-- Comments -->
      <CommentList
        :comments="comments"
        :news-id="newsId"
        @comment-posted="handleCommentPosted"
      />
    </view>
  </view>
</template>

<style lang="scss" scoped>
.detail-page {
  min-height: 100vh;
  background-color: $color-bg;
}

.detail-content {
  padding-bottom: 140rpx;
}

.detail-header {
  padding: $spacing-lg;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;
}

.meta-source {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.meta-time {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
}

.detail-title {
  font-size: $font-size-xxl;
  font-weight: 700;
  color: $color-text;
  line-height: $line-height-tight;
}

.detail-summary {
  padding: 0 $spacing-lg $spacing-lg;
}

.summary-text {
  font-size: $font-size-base;
  color: $color-text;
  line-height: $line-height-relaxed;
}

.editor-comment {
  margin: 0 $spacing-lg $spacing-lg;
  padding: $spacing-md $spacing-lg;
  background-color: $color-bg-secondary;
  border-radius: $radius-sm;
  border-left: 6rpx solid $color-text;
}

.comment-label {
  margin-bottom: $spacing-xs;
}

.label-text {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 2rpx;
}

.comment-body {
  font-size: $font-size-base;
  color: $color-text-secondary;
  line-height: $line-height-relaxed;
}

.source-action {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 $spacing-lg $spacing-lg;
  padding: $spacing-md;
  border: 1rpx solid $color-border;
  border-radius: $radius-sm;
}

.source-btn-text {
  font-size: $font-size-base;
  color: $color-text-secondary;
  margin-right: $spacing-xs;
}

.source-arrow {
  font-size: $font-size-base;
  color: $color-text-tertiary;
}

.coin-hint {
  padding: $spacing-sm $spacing-lg;
  background-color: $color-bg-secondary;
}

.coin-text {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.comment-notice {
  padding: $spacing-sm $spacing-lg;
  background-color: $color-bg-secondary;
}

.comment-notice-text {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

@media (prefers-color-scheme: dark) {
  .detail-page {
    background-color: $color-dark-bg;
  }

  .detail-title {
    color: $color-dark-text;
  }

  .summary-text {
    color: $color-dark-text;
  }

  .editor-comment {
    background-color: $color-dark-bg-secondary;
    border-left-color: $color-dark-text;
  }

  .label-text {
    color: $color-dark-text-secondary;
  }

  .comment-body {
    color: $color-dark-text-secondary;
  }

  .source-action {
    border-color: $color-dark-border;
  }

  .source-btn-text {
    color: $color-dark-text-secondary;
  }

  .coin-hint {
    background-color: $color-dark-bg-secondary;
  }

  .comment-notice {
    background-color: $color-dark-bg-secondary;
  }

  .comment-notice-text {
    color: $color-dark-text-secondary;
  }
}
</style>
