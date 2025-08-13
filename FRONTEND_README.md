# TGState 前端重构

使用 Vite + Vue 3 + Shadcn UI 重构的现代化前端界面。

## 功能特性

- 🎨 现代化的UI设计，基于Shadcn UI组件库
- 📱 完全响应式设计，支持移动端
- 🚀 使用Vite构建，开发体验更佳
- ⚡ Vue 3 Composition API，性能更优
- 📤 支持拖拽上传和多文件上传
- 📊 实时上传进度显示
- 🔄 大文件分片上传支持
- 📋 一键复制链接功能
- 🖼️ 图片文件支持HTML和Markdown格式复制

## 开发环境

### 前端开发
```bash
cd frontend
npm install
npm run dev
```

### 后端开发
```bash
go run main.go
```

### 同时启动前后端（Windows）
```bash
dev-frontend.bat
```

## 生产构建

### 构建前端并编译Go程序（Windows）
```bash
build-frontend.bat
```

### 构建前端并编译Go程序（Linux/Mac）
```bash
chmod +x build-frontend.sh
./build-frontend.sh
```

### 手动构建
```bash
# 构建前端
cd frontend
npm install
npm run build
cd ..

# 构建Go程序
go build -o tgstate
```

## 项目结构

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/           # Shadcn UI组件
│   │   └── FileUpload.vue # 主上传组件
│   ├── lib/
│   │   └── utils.ts      # 工具函数
│   ├── services/
│   │   └── api.ts        # API服务
│   ├── App.vue           # 主应用组件
│   ├── main.ts           # 应用入口
│   └── style.css         # 全局样式
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全的JavaScript
- **Vite** - 下一代前端构建工具
- **Tailwind CSS** - 实用优先的CSS框架
- **Shadcn UI** - 高质量的Vue组件库
- **Axios** - HTTP客户端
- **Lucide Vue** - 图标库

## API兼容性

新前端完全兼容现有的Go后端API：
- `/api` - 文件上传
- `/api/chunk` - 分片上传
- `/api/merge` - 分片合并
- `/d/` - 文件下载
- `/s/` - 短链重定向

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 注意事项

1. 构建后的文件会输出到 `assets/dist/` 目录
2. Go程序会优先使用构建后的前端，如果不存在则回退到原始模板
3. 开发时前端运行在3000端口，通过代理访问后端8088端口
4. 生产环境下前端文件会被嵌入到Go二进制文件中