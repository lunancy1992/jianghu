import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get } from '@/services/api'

export interface NewsItem {
  id: string
  news_id?: string
  title: string
  summary?: string
  source: string
  source_url: string
  published_at: string
  veracity: number
  comment_text: string
  category?: string
  cover_image?: string
  comment_count?: number
  created_at?: string
  is_read?: boolean
}

interface PaginatedResult<T> {
  list: T[]
  total: number
  page: number
  size: number
}

export const useNewsStore = defineStore('news', () => {
  const headlines = ref<NewsItem[]>([])
  const newsList = ref<NewsItem[]>([])
  const currentCategory = ref<string>('全部')
  const newsPage = ref(1)
  const newsTotal = ref(0)
  const isLoadingHeadlines = ref(false)
  const isLoadingNews = ref(false)
  const hasMoreNews = ref(true)

  async function fetchHeadlines(date?: string) {
    if (isLoadingHeadlines.value) return
    isLoadingHeadlines.value = true
    try {
      const result = date
        ? await get<NewsItem[]>('/headlines/history', { date })
        : await get<NewsItem[]>('/headlines')
      headlines.value = (result || []).map((item) => ({
        ...item,
        is_read: item.is_read ?? false,
      }))
    } catch {
      headlines.value = []
    } finally {
      isLoadingHeadlines.value = false
    }
  }

  async function fetchNews(refresh = false) {
    if (isLoadingNews.value) return
    if (!refresh && !hasMoreNews.value) return

    if (refresh) {
      newsPage.value = 1
      hasMoreNews.value = true
    }

    isLoadingNews.value = true
    try {
      const params: Record<string, any> = {
        page: newsPage.value,
        size: 20,
      }
      if (currentCategory.value !== '全部') {
        params.category = currentCategory.value
      }

      const result = await get<PaginatedResult<NewsItem>>('/news', params)
      const list = result.list || []
      if (refresh) {
        newsList.value = list
      } else {
        newsList.value.push(...list)
      }
      newsTotal.value = result.total || 0
      newsPage.value++
      hasMoreNews.value = list.length >= 20
    } catch {
      // Silently fail
    } finally {
      isLoadingNews.value = false
    }
  }

  async function searchNews(keyword: string): Promise<NewsItem[]> {
    try {
      const result = await get<PaginatedResult<NewsItem>>('/news/search', {
        q: keyword,
        page: 1,
        size: 20,
      })
      return result.list
    } catch {
      return []
    }
  }

  function markHeadlineRead(newsId: string) {
    headlines.value = headlines.value.map((item) => ({
      ...item,
      is_read: item.news_id === newsId || item.id === newsId ? true : item.is_read,
    }))
  }

  function markAllHeadlinesRead() {
    headlines.value = headlines.value.map((item) => ({
      ...item,
      is_read: true,
    }))
  }

  function setCategory(category: string) {
    currentCategory.value = category
    fetchNews(true)
  }

  return {
    headlines,
    newsList,
    currentCategory,
    newsPage,
    newsTotal,
    isLoadingHeadlines,
    isLoadingNews,
    hasMoreNews,
    fetchHeadlines,
    fetchNews,
    searchNews,
    markHeadlineRead,
    markAllHeadlinesRead,
    setCategory,
  }
})
