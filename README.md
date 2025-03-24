# 姿态检测系统

基于 OpenCV 和 Go 语言开发的实时姿态检测系统，用于监测和提醒用户保持正确的坐姿。

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

- Go 1.16 或更高版本
- OpenCV 4.5.0 或更高版本
- 摄像头设备

## 安装

1. 克隆仓库：
```bash
git clone https://github.com/nixihz/posture-detection.git
cd posture-detection
```

2. 安装依赖：
```bash
go mod download
```

3. 运行程序：
```bash
go run cmd/main.go
```

## 配置

系统使用 YAML 配置文件进行设置，配置文件位于 `config/config.yaml`。主要配置项包括：

- 检测器参数（人脸检测、姿态检测等）
- 摄像头设置（分辨率、帧率等）
- 显示设置（窗口大小、标题等）
- 通知设置（是否启用、提醒间隔等）

## 模型文件

系统使用以下模型文件进行姿态检测：

- `models/haarcascade_frontalface_default.xml`: 主要的人脸检测模型
- `models/haarcascade_frontalface_alt.xml`: 备用人脸检测模型
- `models/haarcascade_profileface.xml`: 侧视图人脸检测模型

## 编程范式

本项目采用 Vibe Coding 编程范式，这是一种注重代码可读性、可维护性和开发体验的编程方法。主要特点包括：

- 清晰的代码结构和模块化设计
- 统一的命名规范和代码风格
- 完善的错误处理和日志记录
- 灵活的配置管理
- 详细的文档和注释

## 提示词系统

系统使用结构化的提示词来指导各个组件的运行和问题诊断。详细提示词文档请参考 [docs/prompt.txt](docs/prompt.txt)。

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

MIT License 