<script setup lang="ts">
import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { get } from '@/services/api'
import { formatTime } from '@/utils/format'
import VeracityBadge from '@/components/VeracityBadge.vue'
import TabBar from '@/components/TabBar.vue'

interface EventItem {
  id: string
  title: string
  summary: string
  status: string // 'ongoing' | 'resolved'
  veracity: number
  latest_update: string
  updated_at: string
  node_count: number
}

const events = ref<EventItem[]>([])
const isLoading = ref(false)

onLoad(() => {
  loadEvents()
})

onPullDownRefresh(() => {
  loadEvents().finally(() => {
    uni.stopPullDownRefresh()
  })
})

async function loadEvents() {
  isLoading.value = true
  try {
    const result = await get<{ list: EventItem[] }>('/events')
    events.value = result.list || []
  } catch {
    events.value = []
  } finally {
    isLoading.value = false
  }
}

function goToDetail(id: string) {
  uni.navigateTo({
    url: `/pages/events/detail?id=${id}`,
  })
}
</script>

<template>
  <view class="events-page">
    <view v-if="isLoading && events.length === 0" class="loading-text">
      <text>加载中...</text>
    </view>

    <view class="events-list">
      <view
        v-for="event in events"
        :key="event.id"
        class="event-card"
        @tap="goToDetail(event.id)"
      >
        <view class="event-header">
          <view class="event-status" :class="event.status === 'ongoing' ? 'status-ongoing' : 'status-resolved'">
            <text class="status-text">{{ event.status === 'ongoing' ? '进行中' : '已解决' }}</text>
          </view>
          <VeracityBadge :veracity="event.veracity" />
        </view>

        <text class="event-title">{{ event.title }}</text>

        <view class="event-preview">
          <text class="preview-text">{{ event.latest_update || event.summary || '暂无最新进展' }}</text>
        </view>

        <view class="event-footer">
          <text class="footer-time">{{ formatTime(event.updated_at) }} 更新</text>
          <text class="footer-nodes">{{ event.node_count }} 条动态</text>
        </view>
      </view>
    </view>

    <view v-if="!isLoading && events.length === 0" class="empty-state">
      <text>暂无追踪事件</text>
    </view>

    <!-- TabBar -->
    <TabBar :current="2" />
  </view>
</template>

<style lang="scss" scoped>
.events-page {
  min-height: 100vh;
  background-color: $color-bg-secondary;
  padding: $spacing-sm 0;
  padding-bottom: 120rpx;
}

.events-list {
  padding: 0;
}

.event-card {
  background-color: $color-bg;
  padding: $spacing-lg;
  margin-bottom: 2rpx;

  &:active {
    background-color: $color-bg-secondary;
  }
}

.event-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}

.event-status {
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

.event-title {
  font-size: $font-size-lg;
  font-weight: 600;
  color: $color-text;
  line-height: $line-height-tight;
  margin-bottom: $spacing-sm;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.event-preview {
  margin-bottom: $spacing-sm;
}

.preview-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-height-base;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.event-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.footer-nodes {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

@media (prefers-color-scheme: dark) {
  .events-page {
    background-color: $color-dark-bg;
  }

  .event-card {
    background-color: $color-dark-bg-secondary;

    &:active {
      background-color: $color-dark-bg;
    }
  }

  .status-ongoing {
    background-color: $color-dark-text;

    .status-text {
      color: $color-dark-bg;
    }
  }

  .status-resolved {
    background-color: $color-dark-border;
  }

  .event-title {
    color: $color-dark-text;
  }

  .preview-text {
    color: $color-dark-text-secondary;
  }
}
</style>
