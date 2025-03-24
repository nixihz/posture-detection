package detector

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"path/filepath"

	"posture-detector/internal/config"
	"posture-detector/internal/notify"

	"gocv.io/x/gocv"
)

// PoseResult 包含姿态检测的详细结果
type PoseResult struct {
	IsCorrect       bool   // 姿势是否正确
	HasPerson       bool   // 是否检测到人
	FacePosition    string // 面部位置描述
	ErrorMessage    string // 错误信息
	SitDistance     string // 坐姿距离描述
	SitHeight       string // 坐姿高度描述
	SidePosture     string // 侧视图姿势描述
	SideViewPosture string
}

// PoseDetector 姿态检测器
type PoseDetector struct {
	faceCascade    *gocv.CascadeClassifier
	profileCascade *gocv.CascadeClassifier
	window         *gocv.Window
	config         *config.Config
}

type Config struct {
	// ... existing code ...
	SideView struct {
		// ... existing code ...
		Hunchback struct {
			Enabled        bool    `yaml:"enabled"`
			AngleThreshold float64 `yaml:"angle_threshold"` // 驼背角度阈值
			MinConfidence  float64 `yaml:"min_confidence"`  // 最小置信度
		} `yaml:"hunchback"`
	} `yaml:"side_view"`
}

// NewPoseDetector 创建新的姿态检测器
func NewPoseDetector() (*PoseDetector, error) {
	log.Println("正在初始化姿态检测器...")
	cfg := config.GetConfig()

	// 加载人脸检测模型
	faceCascade := gocv.NewCascadeClassifier()

	faceModelPath := filepath.Join("models", "haarcascade_frontalface_default.xml")
	if ok := faceCascade.Load(faceModelPath); !ok {
		return nil, fmt.Errorf("加载人脸检测模型失败: %s", faceModelPath)
	}

	// 加载侧脸检测模型
	profileCascade := gocv.NewCascadeClassifier()

	profileModelPath := filepath.Join("models", "haarcascade_profileface.xml")
	if ok := profileCascade.Load(profileModelPath); !ok {
		return nil, fmt.Errorf("加载侧脸检测模型失败: %s", profileModelPath)
	}

	// 创建窗口
	window := gocv.NewWindow("姿态检测")
	window.ResizeWindow(640, 480)
	window.MoveWindow(0, 0)

	return &PoseDetector{
		faceCascade:    &faceCascade,
		profileCascade: &profileCascade,
		window:         window,
		config:         cfg,
	}, nil
}

// DetectPose 检测姿态
func (d *PoseDetector) DetectPose(img gocv.Mat) (*PoseResult, error) {
	// 创建显示图像
	displayImg := img.Clone()
	defer displayImg.Close()

	// 调整图像大小
	gocv.Resize(displayImg, &displayImg, image.Point{X: 640, Y: 480}, 0, 0, gocv.InterpolationLinear)

	// 转换为灰度图
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(displayImg, &gray, gocv.ColorBGRToGray)

	// 直方图均衡化
	gocv.EqualizeHist(gray, &gray)

	// 检测人脸
	faces := d.faceCascade.DetectMultiScaleWithParams(
		gray,
		d.config.Detector.ScaleFactor,
		d.config.Detector.MinNeighbors,
		0,
		image.Point{X: d.config.Detector.MinFaceSize, Y: d.config.Detector.MinFaceSize},
		image.Point{X: d.config.Detector.MaxFaceSize, Y: d.config.Detector.MaxFaceSize},
	)

	// 如果没有检测到正面人脸，尝试检测侧脸
	if len(faces) == 0 {
		faces = d.profileCascade.DetectMultiScaleWithParams(
			gray,
			d.config.Detector.ScaleFactor,
			d.config.Detector.MinNeighbors,
			0,
			image.Point{X: d.config.Detector.MinFaceSize, Y: d.config.Detector.MinFaceSize},
			image.Point{X: d.config.Detector.MaxFaceSize, Y: d.config.Detector.MaxFaceSize},
		)
	}

	result := &PoseResult{
		HasPerson: len(faces) > 0,
	}

	if result.HasPerson {
		// 获取最大的人脸区域
		maxFace := faces[0]
		for _, face := range faces {
			if face.Dx()*face.Dy() > maxFace.Dx()*maxFace.Dy() {
				maxFace = face
			}
		}

		// 计算人脸位置和姿态
		faceCenterX := maxFace.Min.X + maxFace.Dx()/2
		faceCenterY := maxFace.Min.Y + maxFace.Dy()/2
		imgWidth := displayImg.Cols()
		imgHeight := displayImg.Rows()

		// 判断头部位置
		if faceCenterY < imgHeight/3 {
			result.FacePosition = "头部位置过高"
		} else if faceCenterY > imgHeight*2/3 {
			result.FacePosition = "头部位置过低"
		} else {
			result.FacePosition = "头部位置正常"
		}

		// 判断坐姿距离
		faceSize := float64(maxFace.Dx() * maxFace.Dy())
		imgSize := float64(imgWidth * imgHeight)
		faceRatio := faceSize / imgSize

		if faceRatio < 0.05 {
			result.SitDistance = "距离过远"
		} else if faceRatio > 0.15 {
			result.SitDistance = "距离过近"
		} else {
			result.SitDistance = "距离适中"
		}

		// 判断坐姿高度
		if faceCenterY < imgHeight/4 {
			result.SitHeight = "坐姿过高"
		} else if faceCenterY > imgHeight*3/4 {
			result.SitHeight = "坐姿过低"
		} else {
			result.SitHeight = "坐姿高度正常"
		}

		// 判断侧视图姿势
		if faceCenterX < imgWidth/3 {
			result.SidePosture = "身体前倾"
		} else if faceCenterX > imgWidth*2/3 {
			result.SidePosture = "身体后仰"
		} else {
			result.SidePosture = "坐姿端正"
		}

		// 综合判断姿势是否正确
		result.IsCorrect = result.FacePosition == "头部位置正常" &&
			result.SitDistance == "距离适中" &&
			result.SitHeight == "坐姿高度正常" &&
			result.SidePosture == "坐姿端正"

		// 在图像上绘制检测结果
		gocv.Rectangle(&displayImg, maxFace, color.RGBA{R: 0, G: 255, B: 0, A: 255}, 2)

		// 在左上角显示状态信息
		statusColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		if !result.IsCorrect {
			statusColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		}

		// 使用更小的字体和更紧凑的布局
		gocv.PutText(&displayImg, result.FacePosition,
			image.Point{X: 5, Y: 20}, gocv.FontHersheyPlain, 1.0, statusColor, 1)
		gocv.PutText(&displayImg, result.SitDistance,
			image.Point{X: 5, Y: 40}, gocv.FontHersheyPlain, 1.0, statusColor, 1)
		gocv.PutText(&displayImg, result.SitHeight,
			image.Point{X: 5, Y: 60}, gocv.FontHersheyPlain, 1.0, statusColor, 1)
		gocv.PutText(&displayImg, result.SidePosture,
			image.Point{X: 5, Y: 80}, gocv.FontHersheyPlain, 1.0, statusColor, 1)
	}

	// 检测侧视图姿势
	if d.config.Detector.EnableSideView {
		// 检测侧脸
		profileFaces := d.detectProfileFace(img)
		if len(profileFaces) > 0 {
			// 获取最大的侧脸
			maxFace := profileFaces[0]
			for _, face := range profileFaces {
				if face.Dx()*face.Dy() > maxFace.Dx()*maxFace.Dy() {
					maxFace = face
				}
			}

			// 检测驼背
			if d.config.Detector.EnableHunchbackDetection {
				isHunchback := d.detectHunchback(img, maxFace)
				if isHunchback {
					result.SideViewPosture = "存在驼背"
					// 发送提醒
					notify.SendNotification("检测到驼背，请保持正确的坐姿")
				} else {
					result.SideViewPosture = "坐姿端正"
				}
			}

			// 绘制侧脸检测框
			gocv.Rectangle(&displayImg, maxFace, color.RGBA{0, 255, 0, 255}, 2)
		}
	}

	// 显示图像
	if d.window != nil {
		d.window.IMShow(displayImg)
		d.window.WaitKey(1)
	}

	return result, nil
}

// Close 释放资源
func (d *PoseDetector) Close() {
	if d.window != nil {
		d.window.Close()
	}
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 检测驼背
func (d *PoseDetector) detectHunchback(img gocv.Mat, faceRect image.Rectangle) bool {
	// 获取人脸区域
	faceROI := img.Region(faceRect)
	defer faceROI.Close()

	// 转换为灰度图
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(faceROI, &gray, gocv.ColorBGRToGray)

	// 使用边缘检测
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(gray, &edges, 50, 150)

	// 查找轮廓
	contours := gocv.FindContours(edges, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	// 分析轮廓
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)

		// 计算轮廓的面积
		area := gocv.ContourArea(contour)

		// 过滤掉太小的区域
		if area < 100 {
			continue
		}

		// 计算轮廓的主方向
		angle := d.calculateContourAngle(contour)

		// 如果角度超过阈值，判定为驼背
		if angle > d.config.Detector.HunchbackAngleThreshold {
			return true
		}
	}

	return false
}

// 计算轮廓的主方向角度
func (d *PoseDetector) calculateContourAngle(contour gocv.PointVector) float64 {
	if contour.Size() < 2 {
		return 0
	}

	// 计算轮廓的主方向
	var sumX, sumY float64
	for i := 0; i < contour.Size(); i++ {
		point := contour.At(i)
		sumX += float64(point.X)
		sumY += float64(point.Y)
	}
	centerX := sumX / float64(contour.Size())
	centerY := sumY / float64(contour.Size())

	// 计算与水平线的夹角
	angle := math.Atan2(centerY, centerX) * 180 / math.Pi
	if angle < 0 {
		angle += 360
	}

	return angle
}

// 检测侧脸
func (d *PoseDetector) detectProfileFace(img gocv.Mat) []image.Rectangle {
	cfg := config.GetConfig()

	// 检测侧脸
	profileFaces := d.profileCascade.DetectMultiScaleWithParams(
		img,
		cfg.Detector.ScaleFactor,
		cfg.Detector.MinNeighbors,
		0,
		image.Point{X: cfg.Detector.MinFaceSize, Y: cfg.Detector.MinFaceSize},
		image.Point{X: cfg.Detector.MaxFaceSize, Y: cfg.Detector.MaxFaceSize},
	)

	return profileFaces
}
