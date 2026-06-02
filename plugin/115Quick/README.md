# 115Quick Chrome Extension

## 图标生成

由于 Chrome 扩展需要 PNG 格式的图标，请使用以下方法生成图标：

1. 访问 https://convertio.co/svg-png/ 或其他 SVG 转 PNG 工具
2. 上传 `public/icons/icon.svg` 文件
3. 分别生成以下尺寸：
   - 16x16 像素 -> 保存为 `icon16.png`
   - 48x48 像素 -> 保存为 `icon48.png`
   - 128x128 像素 -> 保存为 `icon128.png`
4. 将生成的 PNG 文件放入 `public/icons/` 目录

## 安装依赖

```bash
npm install
```

## 开发模式

```bash
npm run dev
```

## 构建扩展

```bash
npm run build
```

构建完成后，`dist` 目录即为可发布的 Chrome 扩展。

## 安装扩展

1. 打开 Chrome 浏览器，访问 `chrome://extensions/`
2. 开启"开发者模式"
3. 点击"加载已解压的扩展程序"
4. 选择构建后的 `dist` 目录
