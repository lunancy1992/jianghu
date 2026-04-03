import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get } from '@/services/api'

export interface CoinTransaction {
  id: string
  type: string
  note: string
  amount: number
  balance_after?: number
  created_at: string
}

interface PaginatedResult<T> {
  list: T[]
  total: number
  page: number
  size: number
}

export const useCoinStore = defineStore('coin', () => {
  const balance = ref(0)
  const transactions = ref<CoinTransaction[]>([])
  const txPage = ref(1)
  const hasMoreTx = ref(true)
  const isLoading = ref(false)

  async function fetchBalance() {
    try {
      const result = await get<{ balance: number }>('/coin/balance')
      balance.value = result.balance
    } catch {
      // Silently fail
    }
  }

  async function fetchTransactions(refresh = false) {
    if (isLoading.value) return
    if (!refresh && !hasMoreTx.value) return

    if (refresh) {
      txPage.value = 1
      hasMoreTx.value = true
    }

    isLoading.value = true
    try {
      const result = await get<PaginatedResult<CoinTransaction>>('/coin/transactions', {
        page: txPage.value,
        size: 20,
      })
      if (refresh) {
        transactions.value = result.list
      } else {
        transactions.value.push(...result.list)
      }
      txPage.value++
      hasMoreTx.value = result.list.length >= 20
    } catch {
      // Silently fail
    } finally {
      isLoading.value = false
    }
  }

  return {
    balance,
    transactions,
    isLoading,
    hasMoreTx,
    fetchBalance,
    fetchTransactions,
  }
})
