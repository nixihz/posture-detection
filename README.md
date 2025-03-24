# 姿态检测系统

这是一个基于 OpenCV 的实时姿态检测系统，可以检测用户的坐姿状态并提供实时反馈。

## 功能特点

- 实时人脸检测和姿态分析
- 多维度姿态评估：
  - 头部位置检测
  - 坐姿距离检测
  - 坐姿高度检测
  - 侧视图姿势检测
- 实时反馈和提醒
- 可配置的检测参数
- 低资源占用

## 系统要求

- Go 1.21 或更高版本
- OpenCV 4.x
- 摄像头设备

## 安装

1. 克隆仓库：
```bash
git clone [repository-url]
cd posture-detector
```

2. 安装依赖：
```bash
go mod tidy
```

3. 运行程序：
```bash
go run cmd/main.go
```

## 配置说明

配置文件 `config.yaml` 包含以下设置：

### 摄像头设置
- `width`: 摄像头分辨率宽度
- `height`: 摄像头分辨率高度
- `fps`: 帧率
- `autofocus`: 自动对焦
- `autoexposure`: 自动曝光
- `brightness`: 亮度
- `contrast`: 对比度

### 显示设置
- `show_window`: 是否显示窗口
- `window_width`: 窗口宽度
- `window_height`: 窗口高度
- `window_title`: 窗口标题

### 检测器设置
- `face_cascade`: 正面人脸检测模型
- `profile_cascade`: 侧脸检测模型
- `min_face_size`: 最小人脸尺寸
- `max_face_size`: 最大人脸尺寸
- `scale_factor`: 检测缩放因子
- `min_neighbors`: 最小邻居数

### 提醒设置
- `enabled`: 是否启用提醒
- `interval`: 提醒间隔（秒）

## 模型文件说明

所有模型文件位于 `models/` 目录下，这些是 OpenCV 的 Haar Cascade 分类器模型：

### 1. haarcascade_frontalface_default.xml
- 用途：正面人脸检测
- 特点：
  - 最常用的人脸检测模型
  - 对正面人脸有较好的检测效果
  - 检测速度快，资源占用低
- 适用场景：
  - 用户正对摄像头
  - 光线条件良好
  - 需要快速检测

### 2. haarcascade_frontalface_alt.xml
- 用途：正面人脸检测（替代模型）
- 特点：
  - 比默认模型更严格
  - 误检率更低
  - 检测速度稍慢
- 适用场景：
  - 需要更准确的人脸检测
  - 环境光线复杂
  - 对误检要求严格

### 3. haarcascade_profileface.xml
- 用途：侧脸检测
- 特点：
  - 专门用于检测侧脸
  - 对侧面角度的人脸有较好效果
  - 检测速度适中
- 适用场景：
  - 用户侧对摄像头
  - 需要检测侧面姿态
  - 多角度人脸检测

## 使用说明

1. 启动程序后，确保摄像头正常工作
2. 调整坐姿，使面部在摄像头视野范围内
3. 系统会实时分析您的姿态并显示结果
4. 当检测到不良姿态时，会通过提醒系统通知您

## 注意事项

- 确保摄像头权限已开启
- 保持适当的光线条件
- 建议保持正面朝向摄像头
- 如果检测不准确，可以调整坐姿或光线

## 开发说明

- 使用 Go 语言开发
- 基于 OpenCV 进行图像处理
- 使用 Haar Cascade 分类器进行人脸检测
- 支持实时姿态分析和反馈

## 许可证

[添加许可证信息] 