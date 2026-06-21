<template>
  <div>
    <h2 class="section-title">鱼缸水质看板</h2>
    <div class="tank-grid">
      <div
        v-for="tank in tanks"
        :key="tank.id"
        class="tank-card"
        @click="openTankDetail(tank)"
      >
        <h3>{{ tank.name }}</h3>
        <div class="water-metrics">
          <div class="metric temp">
            <div class="label">温度</div>
            <div class="value">
              {{ getLatestData(tank.id)?.temperature?.toFixed(1) || '--' }}
              <span class="unit">°C</span>
            </div>
          </div>
          <div class="metric ph">
            <div class="label">pH值</div>
            <div class="value">
              {{ getLatestData(tank.id)?.ph?.toFixed(2) || '--' }}
            </div>
          </div>
        </div>
        <div class="record-time">
          {{ getLatestData(tank.id)?.recorded_at || '暂无数据' }}
        </div>
      </div>
    </div>

    <div v-if="selectedTank" class="modal-overlay" @click.self="selectedTank = null">
      <div class="modal">
        <div class="tank-detail-header">
          <h2>{{ selectedTank.name }}</h2>
          <button class="btn btn-secondary" @click="selectedTank = null">关闭</button>
        </div>

        <div class="water-metrics" style="margin-top: 20px;">
          <div class="metric temp">
            <div class="label">当前温度</div>
            <div class="value">
              {{ latestWater?.temperature?.toFixed(1) || '--' }}
              <span class="unit">°C</span>
            </div>
          </div>
          <div class="metric ph">
            <div class="label">当前pH</div>
            <div class="value">
              {{ latestWater?.ph?.toFixed(2) || '--' }}
            </div>
          </div>
        </div>

        <div class="water-history">
          <h3>历史记录</h3>
          <div class="history-list" v-if="waterHistory.length > 0">
            <div v-for="record in waterHistory" :key="record.id" class="history-item">
              <span class="history-time">{{ record.recorded_at }}</span>
              <div class="history-values">
                <span>{{ record.temperature.toFixed(1) }}°C</span>
                <span>pH {{ record.ph.toFixed(2) }}</span>
              </div>
            </div>
          </div>
          <div v-else class="no-data">暂无历史数据</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getTanks, getLatestWater, getWaterRecords } from '../api.js'

const tanks = ref([])
const latestDataMap = ref({})
const selectedTank = ref(null)
const latestWater = ref(null)
const waterHistory = ref([])

const getLatestData = (tankId) => {
  return latestDataMap.value[tankId] || null
}

const loadTanks = async () => {
  const data = await getTanks()
  tanks.value = data
  data.forEach(async (tank) => {
    const latest = await getLatestWater(tank.id)
    if (latest && latest.id) {
      latestDataMap.value[tank.id] = latest
    }
  })
}

const openTankDetail = async (tank) => {
  selectedTank.value = tank
  latestWater.value = null
  waterHistory.value = []

  const latest = await getLatestWater(tank.id)
  latestWater.value = latest

  const history = await getWaterRecords(tank.id, 20)
  waterHistory.value = history
}

onMounted(() => {
  loadTanks()
})
</script>
