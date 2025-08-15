# Giscus 评论系统配置指南

Giscus 是基于 GitHub Discussions 的评论系统，功能强大且易于配置。

## 配置步骤

### 1. 准备 GitHub 仓库
- 确保你有一个公开的 GitHub 仓库
- 在仓库设置中启用 Discussions 功能：
  - 进入仓库 → Settings → General
  - 在 Features 部分勾选 "Discussions"

### 2. 安装 Giscus GitHub App
- 访问 [Giscus GitHub App](https://github.com/apps/giscus)
- 点击 "Install" 安装到你的 GitHub 账户
- 选择要安装的仓库（建议选择 "All repositories" 或选择特定仓库）

### 3. 获取配置参数
- 访问 [giscus.app](https://giscus.app/zh-CN)
- 按照页面指引填写你的仓库信息
- 选择合适的配置选项
- 复制生成的配置参数

### 4. 修改配置文件
编辑 `frontend/src/config/giscus.ts` 文件：

```typescript
export const giscusConfig = {
  // 你的 GitHub 仓库
  repo: 'your-username/your-repo', // 例如：'octocat/Hello-World'
  
  // 从 giscus.app 获取的仓库 ID
  repoId: 'R_kgDOH...',
  
  // Discussion 分类
  category: 'General',
  
  // 从 giscus.app 获取的分类 ID
  categoryId: 'DIC_kwDOH...',
  
  // 其他配置可根据需要调整
  mapping: 'pathname',
  strict: false,
  reactionsEnabled: true,
  emitMetadata: false,
  inputPosition: 'top',
  theme: 'light',
  lang: 'zh-CN',
  enabled: true,
}
```

### 5. 主题选项
可选的主题包括：
- `light` - 浅色主题（默认）
- `dark` - 深色主题
- `preferred_color_scheme` - 跟随系统主题
- `transparent_dark` - 透明深色主题
- `dark_dimmed` - 暗淡深色主题
- `dark_high_contrast` - 高对比度深色主题
- `light_high_contrast` - 高对比度浅色主题
- `dark_protanopia` - 深色红绿色盲友好主题
- `light_protanopia` - 浅色红绿色盲友好主题
- `dark_tritanopia` - 深色蓝黄色盲友好主题
- `light_tritanopia` - 浅色蓝黄色盲友好主题

### 6. 页面映射选项
- `pathname` - 使用页面路径作为 discussion 标题（推荐）
- `url` - 使用页面完整 URL 作为 discussion 标题
- `title` - 使用页面标题作为 discussion 标题
- `og:title` - 使用页面 og:title 作为 discussion 标题
- `specific` - 使用特定术语
- `number` - 使用特定 discussion 编号

## Giscus 的优势

1. **基于 Discussions** - 比 Issues 更适合评论讨论
2. **功能丰富** - 支持 reactions、回复、编辑等
3. **无需数据库** - 所有数据存储在 GitHub
4. **SEO 友好** - 评论内容可被搜索引擎索引
5. **多主题支持** - 包括无障碍友好主题
6. **实时更新** - 支持实时评论更新

## 注意事项

1. **仓库必须是公开的** - Giscus 需要访问公开仓库的 Discussions
2. **启用 Discussions 功能** - 确保仓库的 Discussions 功能已启用
3. **GitHub App 权限** - 确保 Giscus App 有权限访问你的仓库
4. **首次加载** - 第一次访问时可能需要几秒钟加载评论系统

## 禁用评论系统

如果不需要评论功能，可以在配置文件中设置：
```typescript
enabled: false
```

## 故障排除

如果评论系统无法正常显示：
1. 检查仓库名称和 ID 是否正确
2. 确认仓库是公开的且启用了 Discussions
3. 确认已安装 Giscus GitHub App
4. 检查分类和分类 ID 是否正确
5. 检查浏览器控制台是否有错误信息
6. 确认网络可以访问 giscus.app

## 更多信息

- [Giscus 官网](https://giscus.app/zh-CN)
- [Giscus GitHub 仓库](https://github.com/giscus/giscus)
- [GitHub Discussions 文档](https://docs.github.com/en/discussions)