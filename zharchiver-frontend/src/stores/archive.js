import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

const API_BASE = ''

export const useArchiveStore = defineStore('archive', () => {
  // === State ===
  const isDesktopSidebarOpen = ref(true)
  const showAddModal = ref(false)
  const showLoginModal = ref(false)
  const showTagManageModal = ref(false)

  const currentView = ref('home')
  const activeSetting = ref('auth')
  const viewMode = ref(localStorage.getItem('viewMode') || 'grid')

  const searchQuery = ref('')
  const activeCategory = ref('all') // 'all' or tag name

  const answers = ref([])
  const currentAnswer = ref(null)
  const itemToDelete = ref(null)
  const itemToDeleteType = ref('answer')
  const itemToDeleteCallback = ref(null)

  watch(itemToDelete, (newVal) => {
    if (!newVal) {
      itemToDeleteType.value = 'answer'
      itemToDeleteCallback.value = null
    }
  })

  const showGlobalAlert = ref(false)
  const globalAlertTitle = ref('')
  const globalAlertMessage = ref('')

  const toastVisible = ref(false)
  const toastMessage = ref('')
  const toastType = ref('success') // 'success', 'error', 'info'

  // === Actions ===
  const setViewMode = (mode) => {
    viewMode.value = mode
    localStorage.setItem('viewMode', mode)
  }

  const showAlert = (title, message) => {
    globalAlertTitle.value = title
    globalAlertMessage.value = message
    showGlobalAlert.value = true
  }

  let toastTimer = null
  const showToast = (message, type = 'success') => {
    toastMessage.value = message
    toastType.value = type
    toastVisible.value = true
    
    if (toastTimer) clearTimeout(toastTimer)
    toastTimer = setTimeout(() => {
      toastVisible.value = false
    }, 3000)
  }

  const apiFetch = async (url, options = {}) => {
    const token = localStorage.getItem('token')
    const headers = { ...options.headers }
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const res = await fetch(url, { ...options, headers })
    
    if (res.status === 401) {
      showLoginModal.value = true
      throw new Error('未授权，请重新登录')
    }
    
    return res
  }

  const fetchAnswersList = async () => {
    try {
      const res = await apiFetch(`${API_BASE}/api/answers`)
      if (res.ok) {
        answers.value = await res.json()
      }
    } catch (err) {
      console.error('获取列表失败:', err)
    }
  }

  const selectAnswer = async (id) => {
    try {
      const res = await apiFetch(`${API_BASE}/api/answers/${id}`)
      if (res.ok) {
        currentAnswer.value = await res.json()
      }
    } catch (err) {
      console.error('获取详情失败:', err)
    }
  }

  const onArchiveSuccess = async (id) => {
    await fetchAnswersList()
    if (id) {
      selectAnswer(id)
    }
  }

  const confirmDelete = async () => {
    if (!itemToDelete.value) return

    if (itemToDeleteType.value === 'comment') {
      if (itemToDeleteCallback.value) {
        await itemToDeleteCallback.value()
      }
      itemToDelete.value = null
      return
    }

    const id = itemToDelete.value.answer_id
    itemToDelete.value = null
    
    try {
      const res = await apiFetch(`${API_BASE}/api/answers/${id}`, { method: 'DELETE' })
      if (res.ok) {
        if (currentAnswer.value && currentAnswer.value.answer_id === id) {
          currentAnswer.value = null
        }
        await fetchAnswersList()
      } else {
        const errData = await res.json()
        showAlert('删除失败', errData.message)
      }
    } catch (e) {
      console.error(e)
      showAlert('错误', '网络请求失败')
    }
  }

  // === Getters ===
  const tags = computed(() => {
    const map = new Map()
    answers.value.forEach(a => {
      if (a.tag && !map.has(a.tag)) {
        map.set(a.tag, a.tag_color || 'blue')
      }
    })
    return Array.from(map.entries())
      .map(([name, color]) => ({ name, color }))
      .sort((a, b) => a.name.localeCompare(b.name))
  })

  const filteredAnswers = computed(() => {
    let list = answers.value

    if (activeCategory.value !== 'all') {
      list = list.filter(item => item.tag === activeCategory.value)
    }

    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      list = list.filter(item => 
        item.title.toLowerCase().includes(q) || 
        item.author_name.toLowerCase().includes(q)
      )
    }
    
    return list
  })

  return {
    API_BASE,
    isDesktopSidebarOpen,
    showAddModal,
    showLoginModal,
    showTagManageModal,
    currentView,
    activeSetting,
    viewMode,
    searchQuery,
    activeCategory,
    answers,
    currentAnswer,
    itemToDelete,
    itemToDeleteType,
    itemToDeleteCallback,
    showGlobalAlert,
    globalAlertTitle,
    globalAlertMessage,
    toastVisible,
    toastMessage,
    toastType,
    
    setViewMode,
    showAlert,
    showToast,
    apiFetch,
    fetchAnswersList,
    selectAnswer,
    onArchiveSuccess,
    confirmDelete,

    tags,
    filteredAnswers
  }
})
