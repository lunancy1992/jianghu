<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import Timeline from '@/components/Timeline.vue'
import type { TimelineNode } from '@/components/Timeline.vue'
import { API_BASE_URL, get, post } from '@/services/api'
import { useAuth } from '@/composables/useAuth'

interface EventInfo {
  id: string
  title: string
  description: string
  category: string
  status: string
  cover_image: string
  created_at: string
  updated_at: string
}

interface EventDetailResponse {
  event: EventInfo
  nodes: TimelineNode[]
  news: any[] | null
  evidences: any[] | null
}

interface EventDetail {
  event: EventInfo
  nodes: TimelineNode[]
}

const eventDetail = ref<EventDetail | null>(null)
const isLoading = ref(true)
const eventId = ref('')
const evidenceDescription = ref('')
const isSubmittingEvidence = ref(false)
const { requireLogin } = useAuth()

onLoad((query) => {
  if (query?.id) {
    eventId.value = query.id
    loadDetail(query.id)
  }
})

async function loadDetail(id: string) {
  isLoading.value = true
  try {
    const resp = await get<EventDetailResponse>(`/events/${id}`)
    eventDetail.value = {
      event: resp.event,
      nodes: resp.nodes || [],
    }
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    isLoading.value = false
  }
}

async function uploadEvidenceFile(filePath: string): Promise<string> {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${API_BASE_URL}/evidence`,
      filePath,
      name: 'file',
      header: {
        Authorization: `Bearer ${uni.getStorageSync('token') || ''}`,
      },
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (res.statusCode === 200 && data?.code === 0 && data?.data?.url) {
            resolve(data.data.url)
            return
          }
          reject(new Error(data?.message || '上传失败'))
        } catch (error) {
          reject(error)
        }
      },
      fail: reject,
    })
  })
}

function handleUploadEvidence() {
  requireLogin(() => {
    uni.chooseImage({
      count: 1,
      success: async (chooseRes) => {
        const filePath = chooseRes.tempFilePaths?.[0]
        if (!filePath || !eventId.value || isSubmittingEvidence.value) return

        isSubmittingEvidence.value = true
        try {
          const url = await uploadEvidenceFile(filePath)
          await post(`/events/${eventId.value}/evidence`, {
            type: 'image',
            url,
            description: evidenceDescription.value.trim(),
          })
          evidenceDescription.value = ''
          uni.showToast({ title: '举证已提交，待审核', icon: 'none' })
        } catch {
          uni.showToast({ title: '提交失败', icon: 'none' })
        } finally {
          isSubmittingEvidence.value = false
        }
      },
    })
  })
}
</script>

<template>
  <view class="event-detail-page">
    <view v-if="isLoading" class="loading-text">
      <text>加载中...</text>
    </view>

    <view v-if="eventDetail" class="detail-content">
      <!-- Header -->
      <view class="detail-header">
        <view class="header-status">
          <view class="status-badge" :class="eventDetail.event.status === 'ongoing' ? 'status-ongoing' : 'status-resolved'">
            <text class="status-text">{{ eventDetail.event.status === 'ongoing' ? '进行中' : '已解决' }}</text>
          </view>
        </view>
        <text class="detail-title">{{ eventDetail.event.title }}</text>
        <text class="detail-summary">{{ eventDetail.event.description }}</text>
      </view>

      <!-- Divider -->
      <view class="section-divider" />

      <!-- Timeline -->
      <view class="timeline-section">
        <view class="section-header">
          <text class="section-title">事件时间线</text>
          <text class="section-count">{{ eventDetail.nodes.length }} 条动态</text>
        </view>
        <Timeline :nodes="eventDetail.nodes" />
      </view>

      <!-- Upload evidence -->
      <view class="evidence-section">
        <textarea
          v-model="evidenceDescription"
          class="evidence-input"
          maxlength="200"
          placeholder="补充举证说明（可选）"
        />
        <view class="evidence-btn" @tap="handleUploadEvidence">
          <text class="evidence-btn-text">{{ isSubmittingEvidence ? '提交中...' : '上传证据' }}</text>
        </view>
        <text class="evidence-hint">登录后可上传图片举证，提交后将进入待审核状态</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.event-detail-page {
  min-height: 100vh;
  background-color: $color-bg;
}

.detail-content {
  padding-bottom: $spacing-xl;
}

.detail-header {
  padding: $spacing-lg;
}

.header-status {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;
}

.status-badge {
  padding: 4rpx 16rpx;
  border-radius: $radius-sm;
}

.status-ongoing {
  background-color: $color-text;

  .status-text {
    color: $color-bg;
    font-size: $font-size-xs;
    font-weight: 500;
  }
}

.status-resolved {
  background-color: $color-border;

  .status-text {
    color: $color-text-secondary;
    font-size: $font-size-xs;
    font-weight: 500;
  }
}

.detail-title {
  font-size: $font-size-xxl;
  font-weight: 700;
  color: $color-text;
  line-height: $line-height-tight;
  margin-bottom: $spacing-md;
  display: block;
}

.detail-summary {
  font-size: $font-size-base;
  color: $color-text-secondary;
  line-height: $line-height-relaxed;
  display: block;
}

.section-divider {
  height: 16rpx;
  background-color: $color-bg-secondary;
}

.timeline-section {
  padding-top: $spacing-md;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 $spacing-lg $spacing-sm;
}

.section-title {
  font-size: $font-size-md;
  font-weight: 600;
  color: $color-text;
}

.section-count {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.evidence-section {
  padding: $spacing-xl $spacing-lg;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.evidence-input {
  width: 100%;
  min-height: 160rpx;
  padding: $spacing-md;
  margin-bottom: $spacing-sm;
  background-color: $color-bg-secondary;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
  color: $color-text;
  box-sizing: border-box;
}

.evidence-btn {
  width: 100%;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1rpx solid $color-text;
  border-radius: $radius-sm;
  margin-bottom: $spacing-sm;
}

.evidence-btn-text {
  font-size: $font-size-md;
  color: $color-text;
  font-weight: 500;
}

.evidence-hint {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

@media (prefers-color-scheme: dark) {
  .event-detail-page {
    background-color: $color-dark-bg;
  }

  .status-ongoing {
    background-color: $color-dark-text;

    .status-text {
      color: $color-dark-bg;
    }
  }

  .detail-title {
    color: $color-dark-text;
  }

  .detail-summary {
    color: $color-dark-text-secondary;
  }

  .section-divider {
    background-color: $color-dark-bg-secondary;
  }

  .section-title {
    color: $color-dark-text;
  }

  .evidence-btn {
    border-color: $color-dark-text;
  }

  .evidence-btn-text {
    color: $color-dark-text;
  }
}
</style>
