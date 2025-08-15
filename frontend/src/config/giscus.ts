// Giscus 评论系统配置
export const giscusConfig = {
  // 你的 GitHub 仓库，格式：用户名/仓库名
  repo: "bestk/tgState",

  // 仓库 ID (从 giscus.app 获取)
  repoId: "R_kgDONBiMHA",

  // Discussion 分类 (从 giscus.app 获取)
  category: "General",

  // Discussion 分类 ID (从 giscus.app 获取)
  categoryId: "DIC_kwDONBiMHM4CuNL-",

  // 页面 ↔️ discussion 映射关系
  // 'pathname' - 使用页面路径
  // 'url' - 使用页面完整 URL
  // 'title' - 使用页面标题
  // 'og:title' - 使用页面 og:title
  mapping: "pathname" as const,

  // 是否启用严格标题匹配
  strict: false,

  // 是否启用 reactions
  reactionsEnabled: true,

  // 是否发出 reactions
  emitMetadata: false,

  // 输入框位置
  // 'top' - 在评论上方
  // 'bottom' - 在评论下方
  inputPosition: "top" as const,

  // 主题
  // 'light' - 浅色主题
  // 'dark' - 深色主题
  // 'preferred_color_scheme' - 跟随系统主题
  // 'transparent_dark' - 透明深色主题
  // 'dark_dimmed' - 暗淡深色主题
  theme: "light" as const,

  // 语言
  lang: "en" as const,

  // 是否启用评论系统
  enabled: true,
};
