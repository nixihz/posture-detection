# 坐姿检测系统

这是一个基于 OpenCV 和 Go 语言开发的坐姿检测系统，可以实时检测用户的坐姿并通过 Apple Watch 发送提醒。

## 功能特点

- 实时摄像头捕捉
- 姿态检测
- 智能提醒系统
- 支持 Apple Watch 通知

## 系统要求

- Go 1.16 或更高版本
- OpenCV 4.x
- macOS 或 Linux 系统
- 摄像头设备

## 安装步骤

1. 安装 OpenCV：
   ```bash
   # macOS
   brew install opencv
   
   # Linux
   sudo apt-get install libopencv-dev
   ```

2. 安装 Go 依赖：
   ```bash
   go mod download
   ```

3. 下载模型文件：
   ```bash
   # 需要下载 haarcascade_frontalface_default.xml 文件到 models 目录
   ```

## 运行程序

```bash
go run cmd/main.go
```

## 使用说明

1. 运行程序后，系统会自动打开摄像头
2. 程序会实时检测你的坐姿
3. 如果检测到不良坐姿，会通过 Apple Watch 发送提醒
4. 按 Ctrl+C 可以退出程序

## 注意事项

- 确保摄像头可用且未被其他程序占用
- 保持适当的光照条件
- 确保面部在摄像头视野范围内

## 开发计划

- [ ] 添加更精确的姿态检测算法
- [ ] 优化提醒机制
- [ ] 添加数据统计功能
- [ ] 改进用户界面 