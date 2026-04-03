<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { get, post } from '@/services/api'

import TabBar from '@/components/TabBar.vue'
import { useUserStore } from '@/stores/user'

import { computed } from 'vue'

import VeracityBadge from '@/components/VeracityBadge.vue'

// Types
interface Comment {
  id: number
  news_id: number
  user_id: number
  content: string
  stance: string
  like_count: number
  status: number
  created_at: string
  nickname: string
  avatar: string
  news_title?: string
}
interface News {
  id: number
  title: string
  source: string
  source_url: string
  published_at: string
  veracity: number
  comment_text: string
  is_read: boolean
}
interface Headline {
  news_id: number
  rank: number
  title: string
  expire_at: string
}
const userStore = useUserStore()
const isAdmin = computed(() => userStore.userInfo?.role === 'admin')
const activeTab = ref<'audit' | 'headlines'>('audit')
const auditList = ref<Comment[]>([])
const newsList = ref<News[]>([])
const selectedNews = ref<Set<number>>(new Set())
const headlineItems = ref<Headline[]>([])
const isLoading = ref(false)
onMounted(() => {
  if (!isAdmin.value) {
    uni.showToast({ title: '无权限访问', icon: 'none' })
    uni.switchTab({ url: '/pages/profile/index' })
    return
  }
  loadAuditList()
  loadNewsList()
})
async function loadAuditList() {
  isLoading.value = true
  try {
    const result = await get<{ list: Comment[], total: number }>('/admin/audit/queue')
    auditList.value = result.list || []
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    isLoading.value = false
  }
}
async function loadNewsList() {
  try {
    const result = await get<{ list: News[] }>('/news?size=50')
    newsList.value = result.list || []
  } catch {
    uni.showToast({ title: '加载新闻失败', icon: 'none' })
  }
}
async function approveComment(id: number) {
  try {
    await post(`/admin/audit/${id}/approve`)
    uni.showToast({ title: '已通过', icon: 'success' })
    loadAuditList()
  } catch {
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}
async function rejectComment(id: number) {
  try {
    await post(`/admin/audit/${id}/reject`, { reason: '不合规' })
    uni.showToast({ title: '已拒绝', icon: 'none' })
    loadAuditList()
  } catch {
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}
function toggleNewsSelection(newsId: number) {
  if (selectedNews.value.has(newsId)) {
    selectedNews.value.delete(newsId)
  } else {
    selectedNews.value.add(newsId)
  }
}
async function setHeadlines() {
  if (selectedNews.value.size === 0) {
    uni.showToast({ title: '请选择新闻', icon: 'none' })
    return
  }
  const headlines: Headline[] = []
  let rank = 1
  for (const newsId of selectedNews.value) {
    const news = newsList.value.find(n => n.id === newsId)
    if (news) {
      headlines.push({
        news_id: newsId,
        rank: rank++,
        title: news.title,
        expire_at: ''
      })
    }
  }
  try {
    await post('/admin/headlines', { headlines })
    uni.showToast({ title: '头条已设置', icon: 'success' })
    selectedNews.value.clear()
  } catch {
    uni.showToast({ title: '设置失败', icon: 'none' })
  }
}
</script>
<template>
  <view class="admin-page">
    <!-- Header -->
    <view class="admin-header">
      <text class="header-title">管理后台</text>
    </view>
    <!-- Tabs -->
    <view class="tabs-bar">
      <view class="tab-item" :class="{ 'tab-item--active': activeTab === 'audit' }" @tap="activeTab = 'audit'">
        <text class="tab-text">评论审核</text>
      </view>
      <view class="tab-item" :class="{ 'tab-item--active': activeTab === 'headlines' }" @tap="activeTab = 'headlines'">
        <text class="tab-text">头条设置</text>
      </view>
    </view>
    <!-- Audit Tab -->
    <view v-if="activeTab === 'audit'" class="tab-content">
      <view v-if="isLoading" class="loading-text">
        <text>加载中...</text>
      </view>
      <view v-else-if="auditList.length === 0" class="empty-state">
        <text>暂无待审核评论</text>
      </view>
      <view v-else class="audit-list">
        <view v-for="item in auditList" :key="item.id" class="audit-item">
          <view class="audit-header">
            <text class="audit-user">{{ item.nickname }}</text>
            <text class="audit-time">{{ item.created_at }}</text>
          </view>
          <text class="audit-content">{{ item.content }}</text>
          <view class="audit-actions">
            <view class="action-btn action-btn--approve" @tap="approveComment(item.id)">
              <text>通过</text>
            </view>
            <view class="action-btn action-btn--reject" @tap="rejectComment(item.id)">
              <text>拒绝</text>
            </view>
          </view>
        </view>
      </view>
    </view>
    <!-- Headlines Tab -->
    <view v-if="activeTab === 'headlines'" class="tab-content">
      <view class="headlines-toolbar">
        <text class="toolbar-hint">选择新闻设为头条（已选 {{ selectedNews.size }} 条）</text>
        <view class="toolbar-btn" @tap="setHeadlines">
          <text>确认设置</text>
        </view>
      </view>
      <view class="news-select-list">
        <view v-for="news in newsList" :key="news.id" class="news-select-item" :class="{ 'news-select-item--selected': selectedNews.has(news.id) }" @tap="toggleNewsSelection(news.id)">
          <view class="news-checkbox">
            <text v-if="selectedNews.has(news.id)" class="check-mark">✓</text>
          </view>
          <view class="news-info">
            <text class="news-select-title">{{ news.title }}</text>
            <text class="news-select-meta">{{ news.source }} · {{ news.published_at }}</text>
          </view>
        </view>
      </view>
    </view>
    <!-- TabBar -->
    <TabBar :current="3" />
  </view>
</template>
<style lang="scss" scoped>
.admin-page {
  min-height: 100vh;
  background-color: $color-bg-secondary;
  padding-bottom: 120rpx;
}
.admin-header {
  background-color: $color-text;
  padding: $spacing-lg;
}
.header-title {
  font-size: $font-size-xl;
  font-weight: 700;
  color: $color-bg;
}
.tabs-bar {
  display: flex;
  background-color: $color-bg;
  border-bottom: 1rpx solid $color-border;
}
.tab-item {
  flex: 1;
  padding: $spacing-md;
  text-align: center;
}
.tab-item--active {
  border-bottom: 2rpx solid $color-text;
  .tab-text {
    font-weight: 600;
    color: $color-text;
  }
}
.tab-text {
  font-size: $font-size-base;
  color: $color-text-secondary;
}
.tab-content {
  padding: $spacing-md;
}
.audit-list {
  background-color: $color-bg;
}
.audit-item {
  padding: $spacing-lg;
  border-bottom: 1rpx solid $color-border;
}
.audit-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}
.audit-user {
  font-size: $font-size-sm;
  font-weight: 600;
  color: $color-text;
}
.audit-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}
.audit-content {
  font-size: $font-size-base;
  color: $color-text;
  line-height: $line-height-base;
}
.audit-actions {
  display: flex;
  gap: $spacing-md;
  margin-top: $spacing-md;
}
.action-btn {
  padding: $spacing-sm $spacing-lg;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
}
.action-btn--approve {
  background-color: $color-text;
  color: $color-bg;
}
.action-btn--reject {
  background-color: transparent;
  border: 1rpx solid $color-border;
  color: $color-text-secondary;
}
.headlines-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-md;
  background-color: $color-bg;
  margin-bottom: $spacing-sm;
}
.toolbar-hint {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}
.toolbar-btn {
  background-color: $color-text;
  color: $color-bg;
  padding: $spacing-sm $spacing-lg;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
}
.news-select-list {
  background-color: $color-bg;
}
.news-select-item {
  display: flex;
  align-items: flex-start;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1rpx solid $color-border;
}
.news-select-item--selected {
  background-color: rgba(0, 0, 0, 0.02);
}
.news-checkbox {
  width: 40rpx;
  height: 40rpx;
  border: 2rpx solid $color-border;
  border-radius: 8rpx;
  margin-right: $spacing-md;
  display: flex;
  align-items: center;
  justify-content: center;
}
.check-mark{
  color: $color-text;
  font-size: $font-size-lg;
}
.news-info {
  flex: 1;
  min-width: 0;
}
.news-select-title{
  font-size: $font-size-base;
  color: $color-text;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}
.news-select-meta{
  font-size: $font-size-xs;
  color: $color-text-tertiary;
  margin-top: 4rpx;
}
.loading-text, .empty-state {
  text-align: center;
  color: $color-text-tertiary;
  padding: $spacing-xl;
}
@media (prefers-color-scheme: dark) {
  .admin-page {
    background-color: $color-dark-bg;
  }
  .tabs-bar {
    background-color: $color-dark-bg;
    border-bottom-color: $color-dark-border;
  }
  .tab-item--active {
    border-bottom-color: $color-dark-text;
    .tab-text {
      color: $color-dark-text;
    }
  }
  .audit-list, .audit-item {
    background-color: $color-dark-bg-secondary;
    border-bottom-color: $color-dark-border;
  }
  .audit-user {
    color: $color-dark-text;
  }
  .audit-content {
    color: $color-dark-text;
  }
  .action-btn--approve {
    background-color: $color-dark-text;
    color: $color-dark-bg;
  }
  .action-btn--reject {
    border-color: $color-dark-border;
    color: $color-dark-text-secondary;
  }
  .headlines-toolbar, {
    background-color: $color-dark-bg-secondary;
  }
  .toolbar-btn {
    background-color: $color-dark-text;
    color: $color-dark-bg;
  }
  .news-select-list, .news-select-item {
    background-color: $color-dark-bg-secondary;
    border-bottom-color: $color-dark-border;
  }
  .news-select-item--selected {
    background-color: rgba(255, 255, 255, 0.05);
  }
  .news-checkbox{
    border-color: $color-dark-border;
  }
  .check-mark{
    color: $color-dark-text;
  }
  .news-select-title{
    color: $color-dark-text;
  }
}
</style>
