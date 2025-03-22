package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"posture-detector/internal/alert"
	"posture-detector/internal/camera"
	"posture-detector/internal/detector"
)

func checkCameraPermission() error {
	// 在macOS上，我们需要检查是否有摄像头权限
	// 这里我们只是简单地尝试打开摄像头
	testCam, err := camera.NewCamera()
	if err != nil {
		return fmt.Errorf("摄像头访问被拒绝。请确保：\n"+
			"1. 在系统偏好设置 > 安全性与隐私 > 隐私 > 相机中允许程序访问摄像头\n"+
			"2. 没有其他程序正在使用摄像头\n"+
			"3. 摄像头设备正常工作\n"+
			"错误详情: %v", err)
	}
	testCam.Close()
	return nil
}

func main() {
	// 检查摄像头权限
	if err := checkCameraPermission(); err != nil {
		log.Fatal(err)
	}

	// 初始化摄像头
	log.Println("正在初始化摄像头...")
	cam, err := camera.NewCamera()
	if err != nil {
		log.Fatalf("初始化摄像头失败: %v", err)
	}
	defer cam.Close()

	// 初始化姿态检测器
	log.Println("正在初始化姿态检测器...")
	poseDetector, err := detector.NewPoseDetector()
	if err != nil {
		log.Fatalf("初始化姿态检测器失败: %v", err)
	}
	defer poseDetector.Close()

	// 初始化提醒管理器
	alertManager := alert.NewAlertManager()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 开始检测
	log.Println("开始检测姿态...")
	log.Println("按 Ctrl+C 退出程序")

	frameCount := 0
	lastFrameTime := time.Now()
	consecutiveErrors := 0
	maxConsecutiveErrors := 5

	for {
		select {
		case <-sigChan:
			log.Println("程序正在退出...")
			return
		default:
			// 读取摄像头帧
			frame, err := cam.ReadFrame()
			if err != nil {
				log.Printf("读取帧错误: %v", err)
				consecutiveErrors++

				if consecutiveErrors >= maxConsecutiveErrors {
					log.Println("连续错误次数过多，正在重新初始化摄像头...")
					cam.Close()
					time.Sleep(2 * time.Second)

					cam, err = camera.NewCamera()
					if err != nil {
						log.Printf("重新初始化摄像头失败: %v", err)
						continue
					}
					consecutiveErrors = 0
				}

				time.Sleep(500 * time.Millisecond)
				continue
			}
			consecutiveErrors = 0

			// 检测姿态
			result, err := poseDetector.DetectPose(frame)
			if err != nil {
				log.Printf("检测姿态错误: %v", err)
				continue
			}

			frameCount++
			if frameCount%30 == 0 {
				log.Printf("已处理 %d 帧", frameCount)
				log.Printf("检测结果: 检测到人=%v, 姿势=%v, 面部位置=%v",
					result.HasPerson,
					map[bool]string{true: "正常", false: "不正常"}[result.IsCorrect],
					result.FacePosition)

				// 计算帧率
				now := time.Now()
				elapsed := now.Sub(lastFrameTime)
				fps := float64(30) / elapsed.Seconds()
				log.Printf("当前帧率: %.2f FPS", fps)
				lastFrameTime = now
			}

			// 发送提醒
			if !result.IsCorrect {
				alertManager.SendAlert(result.FacePosition)
			}

			// 控制帧率
			time.Sleep(16 * time.Millisecond) // 约60FPS
		}
	}
}
