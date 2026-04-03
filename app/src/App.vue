<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'

onLaunch(() => {
  console.log('江湖小报 App Launch')
})

onShow(() => {
  console.log('App Show')
})
</script>

<style lang="scss">
@import './styles/global.scss';
</style>
