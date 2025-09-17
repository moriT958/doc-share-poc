<template>
  <div :class="statusClasses">
    {{ status.statusText }}
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ConnectionStatus } from '../composables/useWebSocket'

type Props = {
  status: ConnectionStatus
}

const props = defineProps<Props>()

const statusClasses = computed(() => ({
  status: true,
  connected: props.status.status === 'connected',
  disconnected:
    props.status.status === 'disconnected' || props.status.status === 'error',
  connecting: props.status.status === 'connecting',
}))
</script>

<style scoped>
.status {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.connected {
  background-color: #d4edda;
  color: #155724;
}

.disconnected {
  background-color: #f8d7da;
  color: #721c24;
}

.connecting {
  background-color: #fff3cd;
  color: #856404;
}
</style>
