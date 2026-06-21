<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center;">
      <h2 class="section-title">鱼缸水质看板</h2>
      <div style="display: flex; gap: 10px;">
        <button class="btn btn-secondary" @click="loadAllData">🔄 刷新</button>
        <button
          class="btn"
          :class="hasAlert ? 'btn-danger' : 'btn-secondary'"
          @click="sendNotify"
          :disabled="!hasAlert"
          :title="!hasAlert ? '暂无报警' : '一键发送通知'"
        >
          📱 一键发通知
        </button>
        <button class="btn btn-primary" @click="openAlertConfig">⚙️ 报警设置</button>
      </div>
    </div>

    <div v-if="hasAlert" class="alert-banner">
      🚨 检测到温度异常！共 {{ alertingTanks.length }} 个鱼缸报警
    </div>

    <div class="tank-grid">
      <div
        v-for="tank in tanks"
        :key="tank.id"
        class="tank-card"
        :class="isAlerting(tank.id) ? 'tank-alert alert-' + getAlertType(tank.id) : ''"
        @click="openTankDetail(tank)"
      >
        <div class="alert-light" :class="{ 'on': isAlerting(tank.id) }">
          <span v-if="isAlerting(tank.id)">🔴</span>
          <span v-else>🟢</span>
        </div>
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
        <div v-if="isAlerting(tank.id)" class="alert-text">
          ⚠️ {{ getAlertMessage(tank.id) }}
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

        <div style="margin-top: 16px; font-size: 13px; color: #666;">
          正常温度范围：{{ alertConfig.temp_min?.toFixed(1) || '24.0' }}°C ~ {{ alertConfig.temp_max?.toFixed(1) || '28.0' }}°C
        </div>

        <div class="water-history">
          <h3>历史记录</h3>
          <div class="history-list" v-if="waterHistory.length > 0">
            <div v-for="record in waterHistory" :key="record.id" class="history-item">
              <span class="history-time">{{ record.recorded_at }}</span>
              <div class="history-values">
                <span :class="isTempAbnormal(record.temperature) ? 'temp-alert' : ''">
                  {{ record.temperature.toFixed(1) }}°C
                </span>
                <span>pH {{ record.ph.toFixed(2) }}</span>
              </div>
            </div>
          </div>
          <div v-else class="no-data">暂无历史数据</div>
        </div>
      </div>
    </div>

    <div v-if="showConfig" class="modal-overlay" @click.self="showConfig = false">
      <div class="modal">
        <h2>温度报警设置</h2>

        <div v-if="configError" class="error-msg">{{ configError }}</div>
        <div v-if="configSuccess" class="success-msg">{{ configSuccess }}</div>

        <div class="form-group">
          <label>温度下限（°C）</label>
          <input type="number" step="0.1" v-model.number="configForm.temp_min" />
        </div>

        <div class="form-group">
          <label>温度上限（°C）</label>
          <input type="number" step="0.1" v-model.number="configForm.temp_max" />
        </div>

        <div class="form-group">
          <label>连续几次超标才报警</label>
          <input type="number" min="1" max="10" v-model.number="configForm.consecutive_count" />
          <span style="font-size: 12px; color: #888;">推荐 3 次</span>
        </div>

        <div class="form-group">
          <label style="display: flex; align-items: center; gap: 8px;">
            <input
              type="checkbox"
              v-model.number="configForm.notify_enabled"
              :true-value="1"
              :false-value="0"
            />
            开启聊天软件通知（Webhook）
          </label>
        </div>

        <div v-if="configForm.notify_enabled === 1" class="form-group">
          <label>Webhook 地址</label>
          <input
            type="text"
            v-model="configForm.webhook_url"
            placeholder="粘贴企业微信/钉钉/飞书机器人的 Webhook 地址"
          />
          <span style="font-size: 12px; color: #888;">
            消息会以普通文本格式 POST 到这个地址
          </span>
        </div>

        <div class="btn-group">
          <button class="btn btn-secondary" @click="showConfig = false">取消</button>
          <button class="btn btn-primary" @click="saveConfig">保存设置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  getTanks,
  getLatestWater,
  getWaterRecords,
  getAlertConfig,
  updateAlertConfig,
  getAlertStatus,
  triggerAlertNotify
} from '../api.js'

const tanks = ref([])
const latestDataMap = ref({})
const selectedTank = ref(null)
const latestWater = ref(null)
const waterHistory = ref([])

const alertConfig = ref({ temp_min: 24, temp_max: 28, consecutive_count: 3, notify_enabled: 0, webhook_url: '' })
const alertData = ref([])

const showConfig = ref(false)
const configForm = ref({ temp_min: 24, temp_max: 28, consecutive_count: 3, notify_enabled: 0, webhook_url: '' })
const configError = ref('')
const configSuccess = ref('')

let refreshTimer = null

const getLatestData = (tankId) => {
  return latestDataMap.value[tankId] || null
}

const hasAlert = computed(() => alertingTanks.value.length > 0)
const alertingTanks = computed(() => alertData.value.filter(x => x.status?.is_alerting))

const isAlerting = (tankId) => {
  const item = alertData.value.find(x => x.tank_id === tankId)
  return item?.status?.is_alerting
}

const getAlertType = (tankId) => {
  const item = alertData.value.find(x => x.tank_id === tankId)
  return item?.status?.alert_type || ''
}

const getAlertMessage = (tankId) => {
  const item = alertData.value.find(x => x.tank_id === tankId)
  return item?.status?.message || ''
}

const isTempAbnormal = (temp) => {
  return temp < alertConfig.value.temp_min || temp > alertConfig.value.temp_max
}

const loadAllData = async () => {
  try {
    const [tankList, config, status] = await Promise.all([
      getTanks(),
      getAlertConfig(),
      getAlertStatus()
    ])
    tanks.value = tankList
    alertConfig.value = config
    alertData.value = status.alerts || []

    tankList.forEach(async (tank) => {
      const latest = await getLatestWater(tank.id)
      if (latest && latest.id) {
        latestDataMap.value[tank.id] = latest
      }
    })
  } catch (err) {
    console.error('加载数据失败:', err)
  }
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

const openAlertConfig = async () => {
  configError.value = ''
  configSuccess.value = ''
  configForm.value = { ...alertConfig.value }
  showConfig.value = true
}

const saveConfig = async () => {
  configError.value = ''
  configSuccess.value = ''

  if (configForm.value.temp_min >= configForm.value.temp_max) {
    configError.value = '温度下限必须小于温度上限'
    return
  }
  if (configForm.value.consecutive_count < 1 || configForm.value.consecutive_count > 10) {
    configError.value = '连续次数必须在 1 到 10 之间'
    return
  }
  if (configForm.value.notify_enabled === 1 && !configForm.value.webhook_url) {
    configError.value = '开启通知后必须填写 Webhook 地址'
    return
  }

  try {
    await updateAlertConfig(configForm.value)
    configSuccess.value = '设置保存成功！'
    setTimeout(() => {
      showConfig.value = false
      loadAllData()
    }, 800)
  } catch (err) {
    configError.value = err.message || '保存失败'
  }
}

const sendNotify = async () => {
  try {
    const res = await triggerAlertNotify()
    if (res.sent) {
      alert('通知已发送到您的聊天软件！')
    } else {
      alert(res.message || '当前没有需要发送的报警')
    }
  } catch (err) {
    alert('发送失败：' + (err.message || '未知错误'))
  }
}

onMounted(() => {
  loadAllData()
  refreshTimer = setInterval(() => {
    loadAllData()
  }, 15000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>
