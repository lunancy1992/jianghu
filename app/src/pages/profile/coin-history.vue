<script setup lang="ts">
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import { useCoinStore } from '@/stores/coin'
import { formatTime, formatCoin } from '@/utils/format'

const coinStore = useCoinStore()

const typeLabels: Record<string, string> = {
  weekly_grant: '每周活跃奖励',
  comment_cost: '发表评论',
  comment_reward: '热评奖励',
  evidence_reward: '举证奖励',
  admin_adjust: '管理员调整',
}

const typeIcons: Record<string, string> = {
  weekly_grant: '周',
  comment_cost: '评',
  comment_reward: '赏',
  evidence_reward: '证',
  admin_adjust: '调',
}

onLoad(() => {
  coinStore.fetchBalance()
  coinStore.fetchTransactions(true)
})

onReachBottom(() => {
  coinStore.fetchTransactions()
})
</script>

<template>
  <view class="coin-history-page">
    <!-- Balance header -->
    <view class="balance-header">
      <text class="balance-label">当前积分</text>
      <text class="balance-value">{{ coinStore.balance }}</text>
    </view>

    <!-- Transaction list -->
    <view class="tx-list">
      <view v-if="coinStore.isLoading && coinStore.transactions.length === 0" class="loading-text">
        <text>加载中...</text>
      </view>

      <view
        v-for="tx in coinStore.transactions"
        :key="tx.id"
        class="tx-item"
      >
        <view class="tx-icon-box">
          <text class="tx-icon">{{ typeIcons[tx.type] || '?' }}</text>
        </view>
        <view class="tx-content">
          <text class="tx-desc">{{ tx.note || typeLabels[tx.type] || tx.type }}</text>
          <text class="tx-time">{{ formatTime(tx.created_at) }}</text>
        </view>
        <text class="tx-amount" :class="tx.amount > 0 ? 'amount-plus' : 'amount-minus'">
          {{ formatCoin(tx.amount) }}
        </text>
      </view>

      <view v-if="coinStore.isLoading && coinStore.transactions.length > 0" class="loading-text">
        <text>加载更多...</text>
      </view>

      <view v-if="!coinStore.hasMoreTx && coinStore.transactions.length > 0" class="loading-text">
        <text>没有更多了</text>
      </view>

      <view v-if="!coinStore.isLoading && coinStore.transactions.length === 0" class="empty-state">
        <text>暂无积分记录</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.coin-history-page {
  min-height: 100vh;
  background-color: $color-bg-secondary;
}

.balance-header {
  background-color: $color-text;
  padding: $spacing-xl $spacing-lg;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.balance-label {
  font-size: $font-size-sm;
  color: $color-text-tertiary;
  margin-bottom: $spacing-sm;
}

.balance-value {
  font-size: 80rpx;
  font-weight: 700;
  color: $color-bg;
  line-height: 1;
}

.tx-list {
  background-color: $color-bg;
  margin-top: $spacing-sm;
}

.tx-item {
  display: flex;
  align-items: center;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1rpx solid $color-border;
}

.tx-icon-box {
  width: 64rpx;
  height: 64rpx;
  border-radius: $radius-sm;
  background-color: $color-bg-secondary;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: $spacing-md;
  flex-shrink: 0;
}

.tx-icon {
  font-size: $font-size-md;
  color: $color-text;
  font-weight: 600;
}

.tx-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.tx-desc {
  font-size: $font-size-base;
  color: $color-text;
  margin-bottom: 4rpx;
}

.tx-time {
  font-size: $font-size-xs;
  color: $color-text-tertiary;
}

.tx-amount {
  font-size: $font-size-lg;
  font-weight: 600;
  flex-shrink: 0;
  margin-left: $spacing-md;
}

.amount-plus {
  color: $color-text;
}

.amount-minus {
  color: $color-text-secondary;
}

@media (prefers-color-scheme: dark) {
  .coin-history-page {
    background-color: $color-dark-bg;
  }

  .tx-list {
    background-color: $color-dark-bg-secondary;
  }

  .tx-item {
    border-bottom-color: $color-dark-border;
  }

  .tx-icon-box {
    background-color: $color-dark-bg;
  }

  .tx-icon {
    color: $color-dark-text;
  }

  .tx-desc {
    color: $color-dark-text;
  }

  .amount-plus {
    color: $color-dark-text;
  }

  .amount-minus {
    color: $color-dark-text-secondary;
  }
}
</style>
