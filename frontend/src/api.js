const API_BASE = '/api'

export async function getTanks() {
  const res = await fetch(`${API_BASE}/tanks`)
  return res.json()
}

export async function createTank(data) {
  const res = await fetch(`${API_BASE}/tanks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function getWaterRecords(tankId, limit = 50) {
  const url = tankId
    ? `${API_BASE}/water?tank_id=${tankId}&limit=${limit}`
    : `${API_BASE}/water?limit=${limit}`
  const res = await fetch(url)
  return res.json()
}

export async function getLatestWater(tankId) {
  const res = await fetch(`${API_BASE}/water/tank/${tankId}/latest`)
  return res.json()
}

export async function createWaterRecord(data) {
  const res = await fetch(`${API_BASE}/water`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function getBreedingRecords(tankId) {
  const url = tankId
    ? `${API_BASE}/breeding?tank_id=${tankId}`
    : `${API_BASE}/breeding`
  const res = await fetch(url)
  return res.json()
}

export async function createBreedingRecord(data) {
  const res = await fetch(`${API_BASE}/breeding`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function updateBreedingRecord(id, data) {
  const res = await fetch(`${API_BASE}/breeding/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function deleteBreedingRecord(id) {
  const res = await fetch(`${API_BASE}/breeding/${id}`, {
    method: 'DELETE'
  })
  return res.json()
}
