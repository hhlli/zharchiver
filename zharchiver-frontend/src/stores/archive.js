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
  const currentGroup = ref(null)
  const itemToDelete = ref(null)
  const itemToDeleteType = ref('answer')
  const itemToDeleteCallback = ref(null)

  // 记录展开状态的组 key 集合
  const expandedGroups = ref(new Set())

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
        currentGroup.value = null // 进入详情时清除组
      }
    } catch (err) {
      console.error('获取详情失败:', err)
    }
  }

  const selectGroup = (group) => {
    currentGroup.value = group
    currentAnswer.value = null
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
        // 如果正展示该组详情，删除后刻刷新组内容
        if (currentGroup.value) {
          const updatedAnswers = currentGroup.value.answers.filter(a => a.answer_id !== id)
          if (updatedAnswers.length === 0) {
            currentGroup.value = null
          } else {
            currentGroup.value = { ...currentGroup.value, answers: updatedAnswers, count: updatedAnswers.length }
          }
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

  const toggleGroup = (groupKey) => {
    const s = new Set(expandedGroups.value)
    if (s.has(groupKey)) {
      s.delete(groupKey)
    } else {
      s.add(groupKey)
    }
    expandedGroups.value = s
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

  // 分组后的列表：相同 title（非空）归并，空 title 独立展示
  const groupedAnswers = computed(() => {
    const list = filteredAnswers.value
    const groups = []
    const questionIdMap = new Map() // question_id -> group index
    const titleMap = new Map()      // title -> group index

    for (const item of list) {
      const hasTitle = item.title && item.title.trim() !== ''
      const hasQid = item.question_id && item.question_id !== '0' && item.question_id !== 0

      // 无标题：独立展示，不参与任何组
      if (!hasTitle) {
        groups.push({
          groupKey: `solo_${item.answer_id}`,
          title: item.title,
          questionId: item.question_id,
          tag: item.tag,
          tagColor: item.tag_color,
          answers: [item],
          isSolo: true
        })
        continue
      }

      // 尝试用 question_id 找已有组
      let idx = hasQid ? questionIdMap.get(item.question_id) : undefined

      // 再尝试用 title 找已有组
      if (idx === undefined) {
        idx = titleMap.get(item.title)
      }

      if (idx !== undefined) {
        groups[idx].answers.push(item)
        if (hasQid && !questionIdMap.has(item.question_id)) {
          questionIdMap.set(item.question_id, idx)
        }
      } else {
        const newIdx = groups.length
        const groupKey = hasQid ? `qid_${item.question_id}` : `title_${item.title}`
        groups.push({
          groupKey,
          title: item.title,
          questionId: item.question_id,
          tag: item.tag,
          tagColor: item.tag_color,
          answers: [item],
          isSolo: false
        })
        titleMap.set(item.title, newIdx)
        if (hasQid) {
          questionIdMap.set(item.question_id, newIdx)
        }
      }
    }

    return groups.map(g => ({
      ...g,
      count: g.answers.length,
      isExpanded: expandedGroups.value.has(g.groupKey)
    }))
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
    currentGroup,
    itemToDelete,
    itemToDeleteType,
    itemToDeleteCallback,
    expandedGroups,
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
    selectGroup,
    onArchiveSuccess,
    confirmDelete,
    toggleGroup,

    tags,
    filteredAnswers,
    groupedAnswers
  }
})
