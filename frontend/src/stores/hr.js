import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

const DEFAULT_TIMEZONE = 'WIB (+7)'

function normalizeDivision(raw = {}) {
  return {
    id: raw.id,
    name: raw.name || '',
    description: raw.description || '',
    recruitmentStatus: raw.recruitmentStatus || 'Stabil',
    headEmployeeId: raw.headEmployeeId ?? null,
    headPositionId: raw.headPositionId ?? null,
    head: raw.head || '',
    headTitle: raw.headTitle || ''
  }
}

function normalizePosition(raw = {}) {
  return {
    id: raw.id,
    title: raw.title || '',
    code: raw.code || '',
    divisionId: raw.divisionId ?? raw.division_id ?? null,
    divisionName: raw.divisionName || raw.division_name || raw.division?.name || '',
    parentId: raw.parentId ?? raw.parent_id ?? null,
    notes: raw.notes || '',
    grade: raw.grade || '',
    salaryRange: raw.salaryRange || raw.salary_range || ''
  }
}

function normalizeEmployee(raw = {}) {
  const divisionName = raw.divisionName || raw.division_name || raw.division?.name || ''
  const positionTitle = raw.positionTitle || raw.position_title || raw.position?.title || ''

  return {
    id: raw.id,
    employeeCode: raw.employeeCode || raw.employee_code || raw.nik || '',
    nik: raw.nik || raw.employeeCode || raw.employee_code || '',
    fullName: raw.fullName || raw.full_name || '',
    birthPlace: raw.birthPlace || raw.birth_place || '',
    birthDate: raw.birthDate || raw.birth_date || null,
    gender: raw.gender || '',
    bloodType: raw.bloodType || raw.blood_type || '',
    maritalStatus: raw.maritalStatus || raw.marital_status || '',
    religion: raw.religion || '',
    identityType: raw.identityType || raw.identity_type || '',
    identityNumber: raw.identityNumber || raw.identity_number || '',
    addressKtp: raw.addressKtp || raw.address_ktp || '',
    addressDomicile: raw.addressDomicile || raw.address_domicile || raw.address || '',
    address: raw.address || raw.addressDomicile || raw.address_domicile || '',
    phone: raw.phone || '',
    email: raw.email || '',
    timezone: raw.timezone || DEFAULT_TIMEZONE,
    divisionId: raw.divisionId ?? raw.division_id ?? null,
    divisionName,
    positionId: raw.positionId ?? raw.position_id ?? null,
    positionTitle,
    employmentType: raw.employmentType || raw.employment_type || '',
    status: raw.status || 'Aktif',
    joinDate: raw.joinDate || raw.join_date || null,
    photo: raw.photo || ''
  }
}

export const useHRStore = defineStore('hr', () => {
  const divisions = ref([])
  const positions = ref([])
  const employees = ref([])
  const hydrated = ref(false)
  const loading = ref(false)

  async function ensureHydrated() {
    if (hydrated.value) return
    await refreshAll()
  }

  async function refreshAll() {
    loading.value = true
    try {
      const [divisionResponse, positionResponse, employeeResponse] = await Promise.all([
        api.get('/employees/divisions'),
        api.get('/employees/positions'),
        api.get('/employees')
      ])
      divisions.value = (divisionResponse.data?.data || []).map(normalizeDivision)
      positions.value = (positionResponse.data?.data || []).map(normalizePosition)
      employees.value = (employeeResponse.data?.data || []).map(normalizeEmployee)
      hydrated.value = true
    } finally {
      loading.value = false
    }
  }

  async function createDivision(payload) {
    const response = await api.post('/employees/divisions', payload)
    const division = normalizeDivision(response.data?.data || response.data)
    divisions.value = [...divisions.value, division]
    return division
  }

  async function updateDivision(id, payload) {
    const response = await api.put(`/employees/divisions/${id}`, payload)
    const updated = normalizeDivision(response.data?.data || response.data)
    divisions.value = divisions.value.map((division) => (division.id === updated.id ? updated : division))
    return updated
  }

  async function deleteDivision(id) {
    await api.delete(`/employees/divisions/${id}`)
    divisions.value = divisions.value.filter((division) => division.id !== id)
    positions.value = positions.value.map((position) =>
      position.divisionId === id ? { ...position, divisionId: null, divisionName: '' } : position
    )
    employees.value = employees.value.map((employee) =>
      employee.divisionId === id ? { ...employee, divisionId: null, divisionName: '' } : employee
    )
  }

  async function createPosition(payload) {
    const response = await api.post('/employees/positions', payload)
    const position = normalizePosition(response.data?.data || response.data)
    positions.value = [...positions.value, position]
    return position
  }

  async function updatePosition(id, payload) {
    const response = await api.put(`/employees/positions/${id}`, payload)
    const updated = normalizePosition(response.data?.data || response.data)
    positions.value = positions.value.map((position) => (position.id === updated.id ? updated : position))
    return updated
  }

  async function deletePosition(id) {
    await api.delete(`/employees/positions/${id}`)
    positions.value = positions.value
      .map((position) => (position.parentId === id ? { ...position, parentId: null } : position))
      .filter((position) => position.id !== id)
    employees.value = employees.value.map((employee) =>
      employee.positionId === id ? { ...employee, positionId: null, positionTitle: '' } : employee
    )
  }

  async function createEmployee(payload) {
    const response = await api.post('/employees', payload)
    const employee = normalizeEmployee(response.data?.data || response.data)
    if (employee.divisionId) {
      const division = findDivision(employee.divisionId)
      if (division) employee.divisionName = division.name
    }
    if (employee.positionId) {
      const position = findPosition(employee.positionId)
      if (position) employee.positionTitle = position.title
    }
    employees.value = [...employees.value, employee]
    return employee
  }

  async function updateEmployee(id, payload) {
    const response = await api.put(`/employees/${id}`, payload)
    const updated = normalizeEmployee(response.data?.data || response.data)
    if (updated.divisionId) {
      const division = findDivision(updated.divisionId)
      if (division) updated.divisionName = division.name
    }
    if (updated.positionId) {
      const position = findPosition(updated.positionId)
      if (position) updated.positionTitle = position.title
    }
    employees.value = employees.value.map((employee) => (employee.id === updated.id ? updated : employee))
    return updated
  }

  async function deleteEmployee(id) {
    await api.delete(`/employees/${id}`)
    employees.value = employees.value.filter((employee) => employee.id !== id)
  }

  async function deleteEmployeesBatch(ids) {
    const payload = Array.isArray(ids) ? ids.map((id) => Number(id)).filter((id) => !Number.isNaN(id)) : []
    if (!payload.length) return
    await api.delete('/employees', { data: { ids: payload } })
    const toRemove = new Set(payload)
    employees.value = employees.value.filter((employee) => !toRemove.has(Number(employee.id)))
  }

  function employeeCountByDivision(divisionId) {
    return employees.value.filter((employee) => employee.divisionId === divisionId).length
  }

  function findDivision(id) {
    return divisions.value.find((division) => division.id === id) || null
  }

  function findPosition(id) {
    return positions.value.find((position) => position.id === id) || null
  }

  function findEmployee(id) {
    return employees.value.find((employee) => employee.id === id) || null
  }

  return {
    divisions,
    positions,
    employees,
    loading,
    hydrated,
    ensureHydrated,
    refreshAll,
    createDivision,
    updateDivision,
    deleteDivision,
    createPosition,
    updatePosition,
    deletePosition,
    createEmployee,
    updateEmployee,
    deleteEmployee,
    deleteEmployeesBatch,
    employeeCountByDivision,
    findDivision,
    findPosition,
    findEmployee
  }
})
