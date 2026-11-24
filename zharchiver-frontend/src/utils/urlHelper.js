export function getArchiveTargetId() {
  if (typeof window === 'undefined') return null

  // 1. 检查原生 window.location.search (?id=xxx 或 ?answer_id=xxx 或 ?startapp=xxx)
  const searchParams = new URLSearchParams(window.location.search)
  let id = searchParams.get('id') || searchParams.get('answer_id') || searchParams.get('startapp')
  if (id) return id.trim()

  // 2. 检查 Telegram WebApp 原生 initDataUnsafe.start_param
  if (window.Telegram?.WebApp?.initDataUnsafe?.start_param) {
    return String(window.Telegram.WebApp.initDataUnsafe.start_param).trim()
  }

  // 3. 检查 window.location.hash 中的参数 (如 /#/?id=xxx 或 #tgWebAppStartParam=xxx)
  if (window.location.hash) {
    const hash = window.location.hash
    const qIndex = hash.indexOf('?')
    if (qIndex !== -1) {
      const hashParams = new URLSearchParams(hash.substring(qIndex + 1))
      id = hashParams.get('id') || hashParams.get('answer_id') || hashParams.get('startapp')
      if (id) return id.trim()
    }
    const cleanHash = hash.startsWith('#') ? hash.slice(1) : hash
    const hashParams = new URLSearchParams(cleanHash)
    id = hashParams.get('tgWebAppStartParam') || hashParams.get('id') || hashParams.get('answer_id') || hashParams.get('startapp')
    if (id) return id.trim()
  }

  return null
}

export function getAuthTokenFromUrl() {
  if (typeof window === 'undefined') return null
  const searchParams = new URLSearchParams(window.location.search)
  let token = searchParams.get('token')
  if (token) return token.trim()

  if (window.location.hash) {
    const hash = window.location.hash
    const qIndex = hash.indexOf('?')
    if (qIndex !== -1) {
      const hashParams = new URLSearchParams(hash.substring(qIndex + 1))
      token = hashParams.get('token')
      if (token) return token.trim()
    }
    const cleanHash = hash.startsWith('#') ? hash.slice(1) : hash
    const hashParams = new URLSearchParams(cleanHash)
    token = hashParams.get('token')
    if (token) return token.trim()
  }
  return null
}
