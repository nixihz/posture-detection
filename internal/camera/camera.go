package camera

import (
	"fmt"
	"log"
	"time"

	"posture-detector/internal/config"

	"gocv.io/x/gocv"
)

// Camera 封装了摄像头操作
type Camera struct {
	device *gocv.VideoCapture
}

// NewCamera 创建新的摄像头实例
func NewCamera() (*Camera, error) {
	log.Println("正在初始化摄像头...")
	device, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, fmt.Errorf("无法打开摄像头: %v", err)
	}

	cfg := config.GetConfig()

	// 设置摄像头参数
	device.Set(gocv.VideoCaptureFrameWidth, float64(cfg.Camera.Width))
	device.Set(gocv.VideoCaptureFrameHeight, float64(cfg.Camera.Height))
	device.Set(gocv.VideoCaptureFPS, float64(cfg.Camera.FPS))
	device.Set(gocv.VideoCaptureAutoFocus, boolToFloat64(cfg.Camera.Autofocus))
	device.Set(gocv.VideoCaptureAutoExposure, boolToFloat64(cfg.Camera.Autoexposure))
	device.Set(gocv.VideoCaptureBrightness, float64(cfg.Camera.Brightness))
	device.Set(gocv.VideoCaptureContrast, float64(cfg.Camera.Contrast))

	// 验证摄像头是否真正打开
	if !device.IsOpened() {
		return nil, fmt.Errorf("摄像头未能正确打开")
	}

	// 等待摄像头初始化
	time.Sleep(2 * time.Second)

	// 读取一帧以验证摄像头是否正常工作
	img := gocv.NewMat()
	defer img.Close()
	if ok := device.Read(&img); !ok {
		device.Close()
		return nil, fmt.Errorf("无法读取摄像头画面")
	}

	if img.Empty() {
		device.Close()
		return nil, fmt.Errorf("从摄像头读取到空帧")
	}

	log.Printf("摄像头初始化成功，图像尺寸: %dx%d", img.Cols(), img.Rows())

	return &Camera{
		device: device,
	}, nil
}

// ReadFrame 读取一帧图像
func (c *Camera) ReadFrame() (gocv.Mat, error) {
	if !c.device.IsOpened() {
		return gocv.NewMat(), fmt.Errorf("摄像头未打开")
	}

	img := gocv.NewMat()

	// 尝试多次读取帧，以处理偶尔的读取失败
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if ok := c.device.Read(&img); !ok {
			if i == maxRetries-1 {
				img.Close()
				return gocv.NewMat(), fmt.Errorf("无法读取摄像头帧，请检查摄像头连接")
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 检查帧是否为空或无效
		if img.Empty() {
			if i == maxRetries-1 {
				img.Close()
				return gocv.NewMat(), fmt.Errorf("摄像头返回空帧")
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 成功读取到有效帧
		break
	}

	return img, nil
}

// Close 释放资源
func (c *Camera) Close() {
	if c.device != nil {
		c.device.Close()
	}
}

// boolToFloat64 将布尔值转换为float64
func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
