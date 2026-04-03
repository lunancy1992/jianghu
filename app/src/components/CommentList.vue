<script setup lang="ts">
import { ref } from 'vue'
import { formatTime } from '@/utils/format'
import { useAuth } from '@/composables/useAuth'
import { post, del } from '@/services/api'

export interface Comment {
  id: string
  user_id: string
  nickname: string
  avatar: string
  content: string
  created_at: string
  like_count: number
  is_liked: boolean
}

const props = defineProps<{
  comments: Comment[]
  newsId: string
}>()

const emit = defineEmits<{
  (e: 'comment-posted', payload: { status: number }): void
}>()

const { requireLogin } = useAuth()
const inputText = ref('')
const isSubmitting = ref(false)

function handleLike(comment: Comment) {
  requireLogin(async () => {
    try {
      if (comment.is_liked) {
        await del('/like', { comment_id: Number(comment.id) })
        comment.is_liked = false
        comment.like_count = Math.max(0, comment.like_count - 1)
      } else {
        await post('/like', { comment_id: Number(comment.id) })
        comment.is_liked = true
        comment.like_count++
      }
    } catch {
      // Silently fail
    }
  })
}

async function handleSubmit() {
  if (!inputText.value.trim()) return
  if (!requireLogin()) return
  if (isSubmitting.value) return

  isSubmitting.value = true
  try {
    const comment = await post<{ status: number }>(`/news/${props.newsId}/comments`, {
      content: inputText.value.trim(),
    })
    inputText.value = ''
    emit('comment-posted', { status: comment.status })
    uni.showToast({
      title: comment.status === 1 ? '评论成功，消耗1积分' : '评论已提交审核，消耗1积分',
      icon: 'none',
    })
  } catch {
    // Error handled by api layer
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <view class="comment-list">
    <view class="comment-header">
      <text class="comment-title">评论 ({{ comments.length }})</text>
    </view>

    <!-- Comment items -->
    <view v-for="comment in comments" :key="comment.id" class="comment-item">
      <view class="comment-avatar">
        <image
          class="avatar-img"
          :src="comment.avatar || '/static/default-avatar.png'"
          mode="aspectFill"
        />
      </view>
      <view class="comment-body">
        <text class="comment-nickname">{{ comment.nickname }}</text>
        <text class="comment-content">{{ comment.content }}</text>
        <view class="comment-footer">
          <text class="comment-time">{{ formatTime(comment.created_at) }}</text>
          <view class="comment-like" @tap="handleLike(comment)">
            <text class="like-icon">{{ comment.is_liked ? '&#9829;' : '&#9825;' }}</text>
            <text class="like-count">{{ comment.like_count || '' }}</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="comments.length === 0" class="empty-state">
      <text>暂无评论，快来抢沙发</text>
    </view>

    <!-- Input bar -->
    <view class="comment-input-bar safe-bottom">
      <input
        v-model="inputText"
        class="comment-input"
        placeholder="发表评论（消耗1积分）"
        :maxlength="500"
        confirm-type="send"
        @confirm="handleSubmit"
      />
      <view
        class="submit-btn"
        :class="{ 'submit-btn--active': inputText.trim() }"
        @tap="handleSubmit"
      >
        <text class="submit-text">发送</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.comment-list {
  background-color: $color-bg;
  padding-bottom: 120rpx;
}

.comment-header {
  padding: $spacing-lg $spacing-lg $spacing-md;
  border-bottom: 1rpx solid $color-border;
}

.comment-title {
  font-size: $font-size-md;
  font-weight: 600;
  color: $color-text;
}

.comment-item {
  display: flex;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1rpx solid $color-border;
}

.comment-avatar {
  width: 64rpx;
  height: 64rpx;
  flex-shrink: 0;
  margin-right: $spacing-md;
}

.avatar-img {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background-color: $color-bg-secondary;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-nickname {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  margin-bottom: $spacing-xs;
  display: block;
}

.comment-content {
  font-size: $font-size-base;
  color: $color-text;
  line-height: $line-height-relaxed;
  margin-bottom: $spacing-xs;
  display: block;
}

.comment-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.comment-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.comment-like {
  display: flex;
  align-items: center;
  padding: $spacing-xs;
}

.like-icon {
  font-size: $font-size-md;
  color: $color-text-tertiary;
  margin-right: 4rpx;
}

.like-count {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.comment-input-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: $spacing-sm $spacing-lg;
  background-color: $color-bg;
  border-top: 1rpx solid $color-border;
  z-index: 100;
}

.comment-input {
  flex: 1;
  height: 72rpx;
  padding: 0 $spacing-md;
  background-color: $color-bg-secondary;
  border-radius: $radius-sm;
  font-size: $font-size-base;
  color: $color-text;
}

.submit-btn {
  margin-left: $spacing-sm;
  padding: 0 $spacing-lg;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: $color-border;
  border-radius: $radius-sm;

  &--active {
    background-color: $color-text;
  }
}

.submit-text {
  font-size: $font-size-sm;
  color: $color-bg;
}

@media (prefers-color-scheme: dark) {
  .comment-list {
    background-color: $color-dark-bg;
  }

  .comment-header {
    border-bottom-color: $color-dark-border;
  }

  .comment-title {
    color: $color-dark-text;
  }

  .comment-item {
    border-bottom-color: $color-dark-border;
  }

  .comment-content {
    color: $color-dark-text;
  }

  .comment-input-bar {
    background-color: $color-dark-bg;
    border-top-color: $color-dark-border;
  }

  .comment-input {
    background-color: $color-dark-bg-secondary;
    color: $color-dark-text;
  }
}
</style>
