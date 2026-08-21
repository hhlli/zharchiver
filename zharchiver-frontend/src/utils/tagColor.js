/**
 * 统一的标签颜色工具模块
 * 所有组件通过 import 此模块来处理标签颜色，不再各自维护独立的颜色映射。
 *
 * 颜色存储格式约定（存入数据库的值）：
 *   - 纯色：Hex 字符串，如 "#3b82f6"
 *   - 渐变色：CSS linear-gradient 字符串，如 "linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)"
 *   - 历史遗留：颜色名关键词，如 "blue"（仅用于兼容旧数据，新数据不再写入此格式）
 */

/** 历史遗留颜色关键词到 Hex 的映射，仅用于向下兼容旧数据库数据 */
export const LEGACY_COLOR_MAP = {
  blue:   '#3b82f6',
  red:    '#ef4444',
  green:  '#10b981',
  yellow: '#f59e0b',
  purple: '#8b5cf6',
  pink:   '#ec4899',
  indigo: '#6366f1',
  teal:   '#14b8a6',
  orange: '#f97316',
  cyan:   '#06b6d4',
  slate:  '#64748b',
  gray:   '#6b7280',
}

const DEFAULT_HEX = '#3b82f6'

/**
 * 将任意格式的 colorVal 解析为可以直接赋值给 CSS `background` 属性的字符串。
 * @param {string} colorVal - 数据库中存储的颜色值
 * @returns {string} CSS background 值
 */
export function resolveTagBackground(colorVal) {
  if (!colorVal) return DEFAULT_HEX
  // 原生 CSS 格式，直接透传
  if (
    colorVal.startsWith('#') ||
    colorVal.startsWith('linear-gradient') ||
    colorVal.startsWith('rgb') ||
    colorVal.startsWith('hsl')
  ) return colorVal
  // 兼容旧关键词
  return LEGACY_COLOR_MAP[colorVal] || DEFAULT_HEX
}

/**
 * 生成用于标签胶囊（badge）的完整 style 对象（背景/文字/边框）。
 * @param {string} colorVal - 数据库中存储的颜色值
 * @returns {object} Vue :style 对象
 */
export function resolveTagStyle(colorVal) {
  const resolved = resolveTagBackground(colorVal)

  // 渐变色：用渐变作为背景，白色文字
  if (resolved.startsWith('linear-gradient')) {
    return {
      background: resolved,
      color: 'white',
      borderColor: 'transparent',
    }
  }

  // 纯色：透明度 10% 背景 + 原色文字 + 原色 20% 边框
  return {
    backgroundColor: resolved + '1A',
    color: resolved,
    borderColor: resolved + '33',
  }
}

