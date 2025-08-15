// Giscus 评论系统配置
export const giscusConfig = {
  // 你的 GitHub 仓库，格式：用户名/仓库名
  repo: import.meta.env.VITE_GISCUS_REPO || 'your-username/repo',

  // 仓库 ID (从 giscus.app 获取)
  repoId: import.meta.env.VITE_GISCUS_REPO_ID,

  // Discussion 分类 (从 giscus.app 获取)
  category: import.meta.env.VITE_GISCUS_CATEGORY || "General",

  // Discussion 分类 ID (从 giscus.app 获取)
  categoryId: import.meta.env.VITE_GISCUS_CATEGORY_ID || "DIC_kwDONBiMHM4CuNL-",

  // 页面 ↔️ discussion 映射关系
  mapping: (import.meta.env.VITE_GISCUS_MAPPING || "pathname") as
    | "pathname"
    | "url"
    | "title"
    | "og:title",

  // 是否启用严格标题匹配
  strict: import.meta.env.VITE_GISCUS_STRICT === "true",

  // 是否启用 reactions
  reactionsEnabled: import.meta.env.VITE_GISCUS_REACTIONS_ENABLED !== "false",

  // 是否发出 reactions
  emitMetadata: import.meta.env.VITE_GISCUS_EMIT_METADATA === "true",

  // 输入框位置
  inputPosition: (import.meta.env.VITE_GISCUS_INPUT_POSITION || "top") as
    | "top"
    | "bottom",

  // 主题
  theme: (import.meta.env.VITE_GISCUS_THEME || "light") as
    | "light"
    | "dark"
    | "preferred_color_scheme"
    | "transparent_dark"
    | "dark_dimmed",

  // 语言
  lang: (import.meta.env.VITE_GISCUS_LANG || "en") as string,

  // 是否启用评论系统
  enabled: import.meta.env.VITE_GISCUS_ENABLED !== "false",
};
