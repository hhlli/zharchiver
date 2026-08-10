import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

const API_BASE = ''
const PAGE_LIMIT = 50

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
  const activeCategory = ref('all')

  const answers = ref([])
  const currentAnswer = ref(null)
  const currentGroup = ref(null)
  const itemToDelete = ref(null)
  const itemToDeleteType = ref('answer')
  const itemToDeleteCallback = ref(null)
  const allTags = ref([]) // 所有标签，独立接口拉取，不受分页影响

  // 分页状态
  const currentPage = ref(1)
  const totalCount = ref(0)
  const isLoadingMore = ref(false)
  const hasMore = computed(() => answers.value.length < totalCount.value)
  const savedScrollY = ref(0) // 记录进入详情前的滚动位置

  // 展开状态
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
  const toastType = ref('success')

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
    toastTimer = setTimeout(() => { toastVisible.value = false }, 3000)
  }

  const apiFetch = async (url, options = {}) => {
    const token = localStorage.getItem('token')
    const headers = { ...options.headers }
    if (token) headers['Authorization'] = `Bearer ${token}`
    const res = await fetch(url, { ...options, headers })
    if (res.status === 401) {
      showLoginModal.value = true
      throw new Error('未授权，请重新登录')
    }
    return res
  }

  // 构建带参数的 URL
  const buildListURL = (page) => {
    const params = new URLSearchParams()
    params.set('page', page)
    params.set('limit', PAGE_LIMIT)
    if (activeCategory.value !== 'all') params.set('tag', activeCategory.value)
    if (searchQuery.value.trim()) params.set('search', searchQuery.value.trim())
    return `${API_BASE}/api/answers?${params.toString()}`
  }

  // 独立获取全部标签（不受分页限制）
  const fetchAllTags = async () => {
    try {
      const res = await apiFetch(`${API_BASE}/api/tags`)
      if (res.ok) {
        allTags.value = await res.json()
      }
    } catch (err) {
      console.error('获取标签失败:', err)
    }
  }

  // 重置并拉取第 1 页（切换 tag/search 时调用）
  const fetchAnswersList = async () => {
    currentPage.value = 1
    answers.value = []
    totalCount.value = 0
    currentGroup.value = null
    currentAnswer.value = null
    try {
      const res = await apiFetch(buildListURL(1))
      if (res.ok) {
        const data = await res.json()
        answers.value = data.items
        totalCount.value = data.total
        currentPage.value = 1
      }
    } catch (err) {
      console.error('获取列表失败:', err)
    }
    // 同时刷新标签列表
    fetchAllTags()
  }

  // 加载下一页（无限滚动追加）
  const loadMoreAnswers = async () => {
    if (isLoadingMore.value || !hasMore.value) return
    isLoadingMore.value = true
    const nextPage = currentPage.value + 1
    try {
      const res = await apiFetch(buildListURL(nextPage))
      if (res.ok) {
        const data = await res.json()
        answers.value = [...answers.value, ...data.items]
        totalCount.value = data.total
        currentPage.value = nextPage
      }
    } catch (err) {
      console.error('加载更多失败:', err)
    } finally {
      isLoadingMore.value = false
    }
  }

  // tag 或 search 变化时重新拉取
  watch([activeCategory, searchQuery], () => {
    fetchAnswersList()
  })

  const selectAnswer = async (id, keepGroup = false) => {
    // 进入详情前先保存当前滚动位置
    if (typeof window !== 'undefined' && !keepGroup) {
      const el = document.getElementById('main-scroll-container')
      savedScrollY.value = el ? el.scrollTop : 0
    }
    try {
      const res = await apiFetch(`${API_BASE}/api/answers/${id}`)
      if (res.ok) {
        currentAnswer.value = await res.json()
        if (!keepGroup) {
          currentGroup.value = null
        }
        // 进入详情后滚动到顶部
        const el = document.getElementById('main-scroll-container')
        if (el) el.scrollTop = 0
      }
    } catch (err) {
      console.error('获取详情失败:', err)
    }
  }

  // 返回列表并恢复滚动位置
  const goBackToList = () => {
    currentAnswer.value = null
    currentGroup.value = null
    // 等 DOM 更新后恢复滚动位置
    setTimeout(() => {
      const el = document.getElementById('main-scroll-container')
      if (el) el.scrollTop = savedScrollY.value
    }, 30)
  }

  // 点击“所有回答” - 重置并滚回顶部
  const goHome = () => {
    currentAnswer.value = null
    currentGroup.value = null
    activeCategory.value = 'all'
    searchQuery.value = ''
    setTimeout(() => {
      const el = document.getElementById('main-scroll-container')
      if (el) el.scrollTop = 0
    }, 30)
  }

  const selectGroup = async (group) => {
    currentGroup.value = group
    if (group.answers && group.answers.length > 0) {
      // 先用前端现有的数据加载第一篇，保证响应速度
      selectAnswer(group.answers[0].answer_id, true)

      // 后台拉取完整的同一问题下的所有回答（解决分页导致的局部成组问题）
      try {
        const title = encodeURIComponent(group.answers[0].title || '')
        const qid = encodeURIComponent(group.answers[0].question_id || '')
        const res = await apiFetch(`${API_BASE}/api/group/answers?title=${title}&question_id=${qid}`)
        if (res.ok) {
          const fullAnswers = await res.json()
          if (fullAnswers && fullAnswers.length > 0) {
            currentGroup.value = {
              ...currentGroup.value,
              answers: fullAnswers,
              count: fullAnswers.length
            }
          }
        }
      } catch (err) {
        console.error('获取完整分组失败:', err)
      }
    }
  }

  const onArchiveSuccess = async (id) => {
    await fetchAnswersList()
    if (id) selectAnswer(id)
  }

  const confirmDelete = async () => {
    if (!itemToDelete.value) return

    if (itemToDeleteType.value === 'comment') {
      if (itemToDeleteCallback.value) await itemToDeleteCallback.value()
      itemToDelete.value = null
      return
    }

    const id = itemToDelete.value.answer_id
    itemToDelete.value = null

    try {
      const res = await apiFetch(`${API_BASE}/api/answers/${id}`, { method: 'DELETE' })
      if (res.ok) {
        let nextAnswerId = null
        if (currentGroup.value) {
          const updatedAnswers = currentGroup.value.answers.filter(a => a.answer_id !== id)
          if (updatedAnswers.length === 0) {
            currentGroup.value = null
          } else {
            currentGroup.value = { ...currentGroup.value, answers: updatedAnswers, count: updatedAnswers.length }
            nextAnswerId = updatedAnswers[0].answer_id
          }
        }
        
        if (currentAnswer.value?.answer_id === id) {
          if (nextAnswerId) {
            selectAnswer(nextAnswerId, true)
          } else {
            currentAnswer.value = null
          }
        }
        // 直接从本地列表移除，不重新拉取（避免重置分页）
        answers.value = answers.value.filter(a => a.answer_id !== id)
        totalCount.value = Math.max(0, totalCount.value - 1)
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
    s.has(groupKey) ? s.delete(groupKey) : s.add(groupKey)
    expandedGroups.value = s
  }

  // === Getters ===
  // tags 仍从已加载数据计算（够用了，后端 tag 筛选是精准的）
  const tags = computed(() => {
    const map = new Map()
    answers.value.forEach(a => {
      if (a.tag && !map.has(a.tag)) map.set(a.tag, a.tag_color || 'blue')
    })
    return Array.from(map.entries())
      .map(([name, color]) => ({ name, color }))
      .sort((a, b) => a.name.localeCompare(b.name))
  })

  // filteredAnswers 现在直接用已加载的 answers（过滤由后端完成）
  const filteredAnswers = computed(() => answers.value)

  // 分组逻辑保持不变
  const groupedAnswers = computed(() => {
    const list = filteredAnswers.value
    const groups = []
    const questionIdMap = new Map()
    const titleMap = new Map()

    for (const item of list) {
      const hasTitle = item.title && item.title.trim() !== ''
      const hasQid = item.question_id && item.question_id !== '0' && item.question_id !== 0

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

      let idx = hasQid ? questionIdMap.get(item.question_id) : undefined
      if (idx === undefined) idx = titleMap.get(item.title)

      if (idx !== undefined) {
        groups[idx].answers.push(item)
        if (hasQid && !questionIdMap.has(item.question_id)) questionIdMap.set(item.question_id, idx)
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
        if (hasQid) questionIdMap.set(item.question_id, newIdx)
      }
    }

    return groups.map(g => ({
      ...g,
      count: Math.max(g.answers.length, g.answers[0].group_count || 1),
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
    allTags,
    showGlobalAlert,
    globalAlertTitle,
    globalAlertMessage,
    toastVisible,
    toastMessage,
    toastType,
    currentPage,
    totalCount,
    isLoadingMore,
    hasMore,

    setViewMode,
    showAlert,
    showToast,
    apiFetch,
    fetchAnswersList,
    fetchAllTags,
    loadMoreAnswers,
    selectAnswer,
    goBackToList,
    goHome,
    selectGroup,
    onArchiveSuccess,
    confirmDelete,
    toggleGroup,

    tags,
    filteredAnswers,
    groupedAnswers
  }
})
