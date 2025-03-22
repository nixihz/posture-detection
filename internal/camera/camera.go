package camera

import (
	"fmt"
	"log"
	"time"

	"gocv.io/x/gocv"
)

// Camera 封装了摄像头操作
type Camera struct {
	capture *gocv.VideoCapture
	frame   gocv.Mat
}

// NewCamera 创建一个新的摄像头实例
func NewCamera() (*Camera, error) {
	log.Println("正在尝试打开摄像头...")
	capture, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, fmt.Errorf("无法打开摄像头: %v", err)
	}

	// 设置摄像头参数
	log.Println("正在设置摄像头参数...")
	capture.Set(gocv.VideoCaptureFrameWidth, 640)
	capture.Set(gocv.VideoCaptureFrameHeight, 480)
	capture.Set(gocv.VideoCaptureFPS, 30)

	// 验证摄像头是否真正打开
	if !capture.IsOpened() {
		return nil, fmt.Errorf("摄像头未能正确打开")
	}

	// 等待摄像头初始化
	time.Sleep(2 * time.Second)

	// 尝试读取一帧来验证摄像头是否正常工作
	frame := gocv.NewMat()
	defer frame.Close()

	if ok := capture.Read(&frame); !ok {
		return nil, fmt.Errorf("无法从摄像头读取图像")
	}

	if frame.Empty() {
		return nil, fmt.Errorf("从摄像头读取到空帧")
	}

	log.Printf("摄像头初始化成功，图像尺寸: %dx%d", frame.Cols(), frame.Rows())

	return &Camera{
		capture: capture,
		frame:   gocv.NewMat(),
	}, nil
}

// ReadFrame 读取一帧图像
func (c *Camera) ReadFrame() (gocv.Mat, error) {
	if !c.capture.IsOpened() {
		return c.frame, fmt.Errorf("摄像头未打开")
	}

	// 尝试多次读取
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if ok := c.capture.Read(&c.frame); ok {
			if !c.frame.Empty() {
				return c.frame, nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return c.frame, fmt.Errorf("无法读取摄像头帧，请检查摄像头连接")
}

// Close 释放资源
func (c *Camera) Close() {
	if c.capture != nil {
		c.capture.Close()
	}
	if !c.frame.Empty() {
		c.frame.Close()
	}
}
