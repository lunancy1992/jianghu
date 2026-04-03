<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'search', keyword: string): void
}>()

const keyword = ref('')

function handleConfirm() {
  const val = keyword.value.trim()
  if (val) {
    emit('search', val)
  }
}

function handleClear() {
  keyword.value = ''
  emit('search', '')
}
</script>

<template>
  <view class="search-bar">
    <view class="search-inner">
      <text class="search-icon">&#128269;</text>
      <input
        v-model="keyword"
        class="search-input"
        placeholder="搜索新闻、事件..."
        placeholder-class="search-placeholder"
        confirm-type="search"
        @confirm="handleConfirm"
      />
      <text v-if="keyword" class="search-clear" @tap="handleClear">&#10005;</text>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.search-bar {
  padding: $spacing-sm $spacing-lg;
  background-color: $color-bg;
}

.search-inner {
  display: flex;
  align-items: center;
  height: 72rpx;
  padding: 0 $spacing-md;
  background-color: $color-bg-secondary;
  border-radius: $radius-sm;
}

.search-icon {
  font-size: $font-size-md;
  color: $color-text-tertiary;
  margin-right: $spacing-sm;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  height: 72rpx;
  font-size: $font-size-base;
  color: $color-text;
}

.search-placeholder {
  color: $color-text-tertiary;
  font-size: $font-size-base;
}

.search-clear {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
  padding: $spacing-xs;
  margin-left: $spacing-xs;
  flex-shrink: 0;
}

@media (prefers-color-scheme: dark) {
  .search-bar {
    background-color: $color-dark-bg;
  }

  .search-inner {
    background-color: $color-dark-bg-secondary;
  }

  .search-input {
    color: $color-dark-text;
  }
}
</style>
