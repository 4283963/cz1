const API_BASE = '/api'

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  })
  const data = await res.json()
  if (!res.ok) {
    throw new Error(data.error || `请求失败 (${res.status})`)
  }
  return data
}

export async function getTanks() {
  return request(`${API_BASE}/tanks`)
}

export async function createTank(data) {
  return request(`${API_BASE}/tanks`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

export async function getWaterRecords(tankId, limit = 50) {
  const url = tankId
    ? `${API_BASE}/water?tank_id=${tankId}&limit=${limit}`
    : `${API_BASE}/water?limit=${limit}`
  return request(url)
}

export async function getLatestWater(tankId) {
  return request(`${API_BASE}/water/tank/${tankId}/latest`)
}

export async function createWaterRecord(data) {
  return request(`${API_BASE}/water`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

export async function getBreedingRecords(tankId) {
  const url = tankId
    ? `${API_BASE}/breeding?tank_id=${tankId}`
    : `${API_BASE}/breeding`
  return request(url)
}

export async function createBreedingRecord(data) {
  return request(`${API_BASE}/breeding`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

export async function updateBreedingRecord(id, data) {
  return request(`${API_BASE}/breeding/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

export async function deleteBreedingRecord(id) {
  return request(`${API_BASE}/breeding/${id}`, {
    method: 'DELETE'
  })
}

export async function getAlertConfig() {
  return request(`${API_BASE}/alert/config`)
}

export async function updateAlertConfig(data) {
  return request(`${API_BASE}/alert/config`, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

export async function getAlertStatus() {
  return request(`${API_BASE}/alert/status`)
}

export async function triggerAlertNotify() {
  return request(`${API_BASE}/alert/notify`, {
    method: 'POST'
  })
}
