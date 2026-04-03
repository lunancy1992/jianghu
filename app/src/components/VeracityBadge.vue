<script setup lang="ts">
defineProps<{
  veracity: number // 0=待核实, 1=已证实, 2=已辟谣
}>()

const labels: Record<number, string> = {
  0: '待核实',
  1: '已证实',
  2: '已辟谣',
}
</script>

<template>
  <view
    class="veracity-badge"
    :class="{
      'badge-pending': veracity === 0,
      'badge-verified': veracity === 1,
      'badge-debunked': veracity === 2,
    }"
  >
    <text class="badge-text">{{ labels[veracity] || '未知' }}</text>
  </view>
</template>

<style lang="scss" scoped>
.veracity-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2rpx 12rpx;
  border-radius: 4rpx;
  flex-shrink: 0;
}

.badge-pending {
  border: 1rpx solid $color-text;
  background-color: transparent;

  .badge-text {
    color: $color-text;
    font-size: $font-size-xs;
    font-weight: 500;
  }
}

.badge-verified {
  background-color: $color-text;
  border: 1rpx solid $color-text;

  .badge-text {
    color: $color-bg;
    font-size: $font-size-xs;
    font-weight: 500;
  }
}

.badge-debunked {
  background-color: $color-danger;
  border: 1rpx solid $color-danger;

  .badge-text {
    color: $color-bg;
    font-size: $font-size-xs;
    font-weight: 500;
  }
}

@media (prefers-color-scheme: dark) {
  .badge-pending {
    border-color: $color-dark-text;

    .badge-text {
      color: $color-dark-text;
    }
  }

  .badge-verified {
    background-color: $color-dark-text;
    border-color: $color-dark-text;

    .badge-text {
      color: $color-dark-bg;
    }
  }
}
</style>
