# 样式架构说明

## 文件结构
```
src/assets/styles/
├── index.css       # 样式入口文件
├── variables.css   # CSS变量定义
├── global.css      # 全局样式
├── components.css  # 组件样式
└── README.md       # 说明文档
```

## 使用方式

### 1. CSS变量使用
在组件中使用预定义的CSS变量：
```css
.my-component {
  padding: var(--spacing-md);
  border-radius: var(--border-radius-sm);
  transition: var(--transition-normal);
}
```

### 2. 组件样式类
使用预定义的组件样式类：
```vue
<n-button class="fixed-button theme-toggle-button">
  切换主题
</n-button>
```

## 样式变量说明

### 间距系统
- `--spacing-xs`: 8px
- `--spacing-sm`: 12px  
- `--spacing-md`: 16px
- `--spacing-lg`: 20px
- `--spacing-xl`: 24px
- `--spacing-2xl`: 40px

### 圆角系统
- `--border-radius-sm`: 6px
- `--border-radius-md`: 12px
- `--border-radius-lg`: 16px

### 按钮尺寸
- `--button-size-sm`: 30px
- `--button-size-md`: 36px
- `--button-size-lg`: 42px
