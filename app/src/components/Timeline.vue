<script setup lang="ts">
import { formatTime } from '@/utils/format'
import VeracityBadge from './VeracityBadge.vue'

export interface TimelineNode {
  id: string
  title: string
  content: string
  node_time: string
  source: string
  veracity: number
}

defineProps<{
  nodes: TimelineNode[]
}>()
</script>

<template>
  <view class="timeline">
    <view
      v-for="(node, index) in nodes"
      :key="node.id"
      class="timeline-item"
      :class="{ 'timeline-item--last': index === nodes.length - 1 }"
    >
      <!-- Left line and dot -->
      <view class="timeline-left">
        <view class="timeline-dot" :class="{ 'dot-first': index === 0 }" />
        <view v-if="index < nodes.length - 1" class="timeline-line" />
      </view>

      <!-- Content -->
      <view class="timeline-content">
        <view class="timeline-header">
          <text class="timeline-time">{{ formatTime(node.node_time) }}</text>
          <VeracityBadge :veracity="node.veracity" />
        </view>
        <text class="timeline-text">{{ node.title || node.content }}</text>
        <text class="timeline-source">来源: {{ node.source }}</text>
      </view>
    </view>

    <view v-if="!nodes || nodes.length === 0" class="empty-state">
      <text>暂无时间线数据</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.timeline {
  padding: $spacing-md $spacing-lg;
}

.timeline-item {
  display: flex;
  padding-bottom: $spacing-lg;

  &--last {
    padding-bottom: 0;
  }
}

.timeline-left {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 32rpx;
  flex-shrink: 0;
  margin-right: $spacing-md;
}

.timeline-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: $color-text-tertiary;
  border: 2rpx solid $color-text-tertiary;
  flex-shrink: 0;
  margin-top: 8rpx;
}

.dot-first {
  background-color: $color-text;
  border-color: $color-text;
  width: 20rpx;
  height: 20rpx;
}

.timeline-line {
  width: 2rpx;
  flex: 1;
  background-color: $color-border;
  margin-top: $spacing-xs;
}

.timeline-content {
  flex: 1;
  padding-bottom: $spacing-sm;
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-xs;
}

.timeline-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.timeline-text {
  font-size: $font-size-base;
  color: $color-text;
  line-height: $line-height-relaxed;
  margin-bottom: $spacing-xs;
}

.timeline-source {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

@media (prefers-color-scheme: dark) {
  .timeline-dot {
    background-color: $color-dark-text-secondary;
    border-color: $color-dark-text-secondary;
  }

  .dot-first {
    background-color: $color-dark-text;
    border-color: $color-dark-text;
  }

  .timeline-line {
    background-color: $color-dark-border;
  }

  .timeline-text {
    color: $color-dark-text;
  }
}
</style>
