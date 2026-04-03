<script setup lang="ts">
import { onLoad, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { ref } from 'vue'
import SearchBar from '@/components/SearchBar.vue'
import NewsCard from '@/components/NewsCard.vue'
import TabBar from '@/components/TabBar.vue'
import { useNewsStore } from '@/stores/news'

const newsStore = useNewsStore()

const categories = ['全部', '国际', '国内', '科技', '财经', '社会']
const searchResults = ref<any[]>([])
const isSearching = ref(false)

onLoad(() => {
  newsStore.fetchNews(true)
})

onPullDownRefresh(() => {
  if (isSearching.value) {
    uni.stopPullDownRefresh()
    return
  }
  newsStore.fetchNews(true).finally(() => {
    uni.stopPullDownRefresh()
  })
})

onReachBottom(() => {
  if (isSearching.value) return
  newsStore.fetchNews()
})

function handleCategoryChange(cat: string) {
  isSearching.value = false
  searchResults.value = []
  newsStore.setCategory(cat)
}

async function handleSearch(keyword: string) {
  if (!keyword) {
    isSearching.value = false
    searchResults.value = []
    return
  }
  isSearching.value = true
  searchResults.value = await newsStore.searchNews(keyword)
}

function handleRead(id: string) {
  console.log('Read news:', id)
}

function handleOpenSource(url: string) {
  console.log('Open source:', url)
}
</script>

<template>
  <view class="news-page">
    <!-- Search -->
    <SearchBar @search="handleSearch" />

    <!-- Category tabs -->
    <view v-if="!isSearching" class="category-bar">
      <scroll-view scroll-x class="category-scroll">
        <view class="category-list">
          <view
            v-for="cat in categories"
            :key="cat"
            class="category-item"
            :class="{ 'category-item--active': newsStore.currentCategory === cat }"
            @tap="handleCategoryChange(cat)"
          >
            <text class="category-text">{{ cat }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Search results -->
    <view v-if="isSearching" class="news-list">
      <view v-if="searchResults.length === 0" class="empty-state">
        <text>未找到相关结果</text>
      </view>
      <NewsCard
        v-for="item in searchResults"
        :key="item.id"
        :news="{
          id: item.id,
          title: item.title,
          source: item.source,
          published_at: item.published_at,
          source_url: item.source_url,
          veracity: item.veracity,
          comment_text: item.comment_text,
        }"
        @read="handleRead"
        @open-source="handleOpenSource"
      />
    </view>

    <!-- News list -->
    <view v-else class="news-list">
      <view v-if="newsStore.isLoadingNews && newsStore.newsList.length === 0" class="loading-text">
        <text>加载中...</text>
      </view>

      <NewsCard
        v-for="item in newsStore.newsList"
        :key="item.id"
        :news="{
          id: item.id,
          title: item.title,
          source: item.source,
          published_at: item.published_at,
          source_url: item.source_url,
          veracity: item.veracity,
          comment_text: item.comment_text,
        }"
        @read="handleRead"
        @open-source="handleOpenSource"
      />

      <view v-if="newsStore.isLoadingNews && newsStore.newsList.length > 0" class="loading-text">
        <text>加载更多...</text>
      </view>

      <view v-if="!newsStore.hasMoreNews && newsStore.newsList.length > 0" class="loading-text">
        <text>没有更多了</text>
      </view>

      <view v-if="!newsStore.isLoadingNews && newsStore.newsList.length === 0" class="empty-state">
        <text>暂无资讯</text>
      </view>
    </view>

    <!-- TabBar -->
    <TabBar :current="1" />
  </view>
</template>

<style lang="scss" scoped>
.news-page {
  min-height: 100vh;
  background-color: $color-bg;
  padding-bottom: 120rpx;
}

.category-bar {
  background-color: $color-bg;
  border-bottom: 1rpx solid $color-border;
}

.category-scroll {
  white-space: nowrap;
}

.category-list {
  display: flex;
  padding: 0 $spacing-md;
}

.category-item {
  padding: $spacing-md $spacing-lg;
  flex-shrink: 0;
  position: relative;

  &--active {
    .category-text {
      color: $color-text;
      font-weight: 600;
    }

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 40rpx;
      height: 4rpx;
      background-color: $color-text;
      border-radius: 2rpx;
    }
  }
}

.category-text {
  font-size: $font-size-base;
  color: $color-text-secondary;
}

.news-list {
  background-color: $color-bg;
}

@media (prefers-color-scheme: dark) {
  .news-page {
    background-color: $color-dark-bg;
  }

  .category-bar {
    background-color: $color-dark-bg;
    border-bottom-color: $color-dark-border;
  }

  .category-item--active {
    .category-text {
      color: $color-dark-text;
    }

    &::after {
      background-color: $color-dark-text;
    }
  }

  .category-text {
    color: $color-dark-text-secondary;
  }

  .news-list {
    background-color: $color-dark-bg;
  }
}
</style>
