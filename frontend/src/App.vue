<template>
  <div class="app">
    <div class="header">
      <h1>🐠 孔雀鱼繁育监控系统</h1>
      <div class="nav">
        <button :class="{ active: currentView === 'dashboard' }" @click="currentView = 'dashboard'">水质看板</button>
        <button :class="{ active: currentView === 'breeding' }" @click="currentView = 'breeding'">繁育记录</button>
      </div>
    </div>
    <div class="content">
      <Dashboard v-if="currentView === 'dashboard'" />
      <Breeding v-else />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Dashboard from './views/Dashboard.vue'
import Breeding from './views/Breeding.vue'

const currentView = ref('dashboard')
</script>
