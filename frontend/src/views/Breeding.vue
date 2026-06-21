<template>
  <div>
    <h2 class="section-title">繁育记录管理</h2>

    <div class="breeding-list">
      <div v-for="record in breedingRecords" :key="record.id" class="breeding-item">
        <div class="breeding-info">
          <h4>{{ record.strain }}</h4>
          <p>鱼缸: {{ getTankName(record.tank_id) }}</p>
          <p>配对日期: {{ record.pair_date }}</p>
          <p v-if="record.expected_birth_date">预计产仔: {{ record.expected_birth_date }}</p>
          <p v-if="record.notes" style="color: #666;">{{ record.notes }}</p>
        </div>
        <div style="text-align: right;">
          <span :class="'status-badge status-' + record.status">
            {{ statusText(record.status) }}
          </span>
          <div style="margin-top: 10px; display: flex; gap: 8px;">
            <button class="btn btn-secondary" @click="editRecord(record)">编辑</button>
            <button class="btn btn-danger" @click="deleteRecord(record.id)">删除</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="breedingRecords.length === 0" class="empty-state">
      暂无繁育记录，点击右下角 + 添加
    </div>

    <button class="add-btn" @click="openAddModal">+</button>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h2>{{ editingId ? '编辑繁育记录' : '新增繁育登记' }}</h2>

        <div v-if="errorMsg" class="error-msg">
          {{ errorMsg }}
        </div>

        <div class="form-group">
          <label>鱼缸</label>
          <select v-model="form.tank_id">
            <option v-for="tank in tanks" :key="tank.id" :value="tank.id">
              {{ tank.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>品系</label>
          <input type="text" v-model="form.strain" placeholder="如：白子孔雀鱼" />
        </div>

        <div class="form-group">
          <label>配对日期</label>
          <input type="date" v-model="form.pair_date" />
        </div>

        <div class="form-group">
          <label>预计产仔日期</label>
          <input type="date" v-model="form.expected_birth_date" />
        </div>

        <div class="form-group">
          <label>状态</label>
          <select v-model="form.status">
            <option value="breeding">繁育中</option>
            <option value="born">已产仔</option>
            <option value="completed">已完成</option>
          </select>
        </div>

        <div class="form-group">
          <label>备注</label>
          <textarea v-model="form.notes" rows="3" placeholder="选填"></textarea>
        </div>

        <div class="btn-group">
          <button class="btn btn-secondary" @click="closeModal">取消</button>
          <button class="btn btn-primary" @click="submitForm">
            {{ editingId ? '保存修改' : '确认添加' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getTanks, getBreedingRecords, createBreedingRecord, updateBreedingRecord, deleteBreedingRecord } from '../api.js'

const tanks = ref([])
const breedingRecords = ref([])
const showModal = ref(false)
const editingId = ref(null)
const errorMsg = ref('')

const form = ref({
  tank_id: 1,
  strain: '',
  pair_date: '',
  expected_birth_date: '',
  status: 'breeding',
  notes: ''
})

const statusText = (status) => {
  const map = {
    breeding: '繁育中',
    born: '已产仔',
    completed: '已完成'
  }
  return map[status] || status
}

const getTankName = (tankId) => {
  const tank = tanks.value.find(t => t.id === tankId)
  return tank ? tank.name : '未知'
}

const loadTanks = async () => {
  tanks.value = await getTanks()
}

const loadBreedingRecords = async () => {
  breedingRecords.value = await getBreedingRecords()
}

const openAddModal = () => {
  editingId.value = null
  errorMsg.value = ''
  form.value = {
    tank_id: tanks.value[0]?.id || 1,
    strain: '',
    pair_date: new Date().toISOString().split('T')[0],
    expected_birth_date: '',
    status: 'breeding',
    notes: ''
  }
  showModal.value = true
}

const editRecord = (record) => {
  editingId.value = record.id
  errorMsg.value = ''
  form.value = { ...record }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingId.value = null
}

const submitForm = async () => {
  errorMsg.value = ''

  if (!form.value.strain || !form.value.pair_date) {
    errorMsg.value = '请填写品系和配对日期'
    return
  }

  if (typeof form.value.tank_id === 'string') {
    form.value.tank_id = parseInt(form.value.tank_id)
  }

  try {
    if (editingId.value) {
      await updateBreedingRecord(editingId.value, form.value)
    } else {
      await createBreedingRecord(form.value)
    }
    closeModal()
    loadBreedingRecords()
  } catch (err) {
    errorMsg.value = err.message || '提交失败，请重试'
  }
}

const deleteRecord = async (id) => {
  if (!confirm('确定要删除这条记录吗？')) return
  await deleteBreedingRecord(id)
  loadBreedingRecords()
}

onMounted(async () => {
  await loadTanks()
  await loadBreedingRecords()
})
</script>
