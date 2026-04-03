import { get } from './api'

export interface News {
  id: number
  title: string
  source: string
  published_at: string
}

export interface Comment {
  id: number
  news_id: number
  user_id: number
  nickname: string
  content: string
  created_at: string
}

export interface AuditQueueResponse {
  list: Comment[]
  total: number
}

export function getAuditQueue(page = number = 20): Promise<AuditQueueResponse> {
  return get<AuditQueueResponse>('/admin/audit/queue', { page })
}

export function approveComment(id: number): Promise<void> {
  return post<void>(`/admin/audit/${id}/approve`)
}
export function rejectComment(id: number, reason: string): Promise<void> {
  return post<void>(`/admin/audit/${id}/reject`, { reason })
}
export function getNewsList(): Promise<News[]> {
  return get<News[]>('/news')
}
export function isAdmin(): boolean {
  const userInfo = useUserStore()
  return userInfo.value?.role === 'admin'
}
