package detector

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"gocv.io/x/gocv"
)

// PoseResult 包含姿态检测的详细结果
type PoseResult struct {
	IsCorrect    bool   // 姿势是否正确
	HasPerson    bool   // 是否检测到人
	FacePosition string // 面部位置描述
	ErrorMessage string // 错误信息
}

// PoseDetector 姿态检测器
type PoseDetector struct {
	classifier gocv.CascadeClassifier
	window     *gocv.Window
}

// NewPoseDetector 创建新的姿态检测器
func NewPoseDetector() (*PoseDetector, error) {
	classifier := gocv.NewCascadeClassifier()
	classifierPath := filepath.Join("models", "haarcascade_frontalface_alt.xml")
	if !classifier.Load(classifierPath) {
		return nil, fmt.Errorf("无法加载分类器: %s", classifierPath)
	}

	window := gocv.NewWindow("姿态检测")
	window.ResizeWindow(800, 600)

	return &PoseDetector{
		classifier: classifier,
		window:     window,
	}, nil
}

// DetectPose 检测姿态
func (pd *PoseDetector) DetectPose(img gocv.Mat) (*PoseResult, error) {
	result := &PoseResult{
		IsCorrect: true,
		HasPerson: false,
	}

	// 转换为灰度图像
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// 调整图像大小以提高检测效果
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(gray, &resized, image.Point{}, 1.2, 1.2, gocv.InterpolationLinear)

	// 调整图像对比度和亮度以提高检测率
	gocv.EqualizeHist(resized, &resized)

	// 在图像上绘制网格线以帮助调试
	imgCopy := img.Clone()
	defer imgCopy.Close()

	// 绘制水平和垂直中心线
	centerY := img.Rows() / 2
	centerX := img.Cols() / 2
	gocv.Line(&imgCopy,
		image.Point{X: 0, Y: centerY},
		image.Point{X: img.Cols(), Y: centerY},
		color.RGBA{R: 255, G: 255, B: 0, A: 255}, 1)
	gocv.Line(&imgCopy,
		image.Point{X: centerX, Y: 0},
		image.Point{X: centerX, Y: img.Rows()},
		color.RGBA{R: 255, G: 255, B: 0, A: 255}, 1)

	// 使用更宽松的检测参数
	rects := pd.classifier.DetectMultiScaleWithParams(
		resized,
		1.05,                      // scaleFactor
		2,                         // minNeighbors
		0,                         // flags
		image.Point{X: 30, Y: 30}, // minSize
		image.Point{},             // maxSize
	)

	// 调整检测到的矩形区域的大小以匹配原始图像
	for i := range rects {
		rects[i].Min.X = int(float64(rects[i].Min.X) / 1.2)
		rects[i].Min.Y = int(float64(rects[i].Min.Y) / 1.2)
		rects[i].Max.X = int(float64(rects[i].Max.X) / 1.2)
		rects[i].Max.Y = int(float64(rects[i].Max.Y) / 1.2)
	}

	if len(rects) == 0 {
		result.HasPerson = false
		result.FacePosition = "未检测到人脸"
		// 显示提示信息
		gocv.PutText(&imgCopy, "未检测到人脸", image.Point{X: 10, Y: 30},
			gocv.FontHersheyPlain, 1.2, color.RGBA{R: 255, G: 0, B: 0, A: 255}, 2)

		// 显示图像尺寸信息
		sizeInfo := fmt.Sprintf("图像尺寸: %dx%d", img.Cols(), img.Rows())
		gocv.PutText(&imgCopy, sizeInfo, image.Point{X: 10, Y: 60},
			gocv.FontHersheyPlain, 1.2, color.RGBA{R: 255, G: 255, B: 0, A: 255}, 2)

		// 显示检测参数
		detInfo := fmt.Sprintf("检测参数: scale=1.05, neighbors=2, minSize=30x30")
		gocv.PutText(&imgCopy, detInfo, image.Point{X: 10, Y: 90},
			gocv.FontHersheyPlain, 1.2, color.RGBA{R: 255, G: 255, B: 0, A: 255}, 2)
	} else {
		result.HasPerson = true

		// 获取图像中心
		centerY := img.Rows() / 2
		faceY := rects[0].Min.Y

		// 计算面部位置
		if faceY < centerY-50 {
			result.FacePosition = "头部位置过高"
			result.IsCorrect = false
		} else if faceY > centerY+50 {
			result.FacePosition = "头部位置过低"
			result.IsCorrect = false
		} else {
			result.FacePosition = "头部位置正常"
		}

		// 检查头部是否倾斜
		faceWidth := rects[0].Dx()
		faceHeight := rects[0].Dy()
		if float64(faceWidth)/float64(faceHeight) < 0.7 || float64(faceWidth)/float64(faceHeight) > 1.3 {
			result.FacePosition = "头部倾斜"
			result.IsCorrect = false
		}

		// 绘制人脸检测框
		gocv.Rectangle(&imgCopy, rects[0], color.RGBA{R: 0, G: 255, B: 0, A: 255}, 2)

		// 显示姿态信息
		textColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		if !result.IsCorrect {
			textColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		}
		gocv.PutText(&imgCopy, result.FacePosition, image.Point{X: 10, Y: 30},
			gocv.FontHersheyPlain, 1.2, textColor, 2)

		// 显示人脸框信息
		faceInfo := fmt.Sprintf("人脸位置: (%d,%d) 尺寸: %dx%d",
			rects[0].Min.X, rects[0].Min.Y, faceWidth, faceHeight)
		gocv.PutText(&imgCopy, faceInfo, image.Point{X: 10, Y: 60},
			gocv.FontHersheyPlain, 1.2, color.RGBA{R: 255, G: 255, B: 0, A: 255}, 2)

		// 显示检测参数
		detInfo := fmt.Sprintf("检测参数: scale=1.05, neighbors=2, minSize=30x30")
		gocv.PutText(&imgCopy, detInfo, image.Point{X: 10, Y: 90},
			gocv.FontHersheyPlain, 1.2, color.RGBA{R: 255, G: 255, B: 0, A: 255}, 2)
	}

	// 显示处理后的图像
	pd.window.IMShow(imgCopy)
	pd.window.WaitKey(1)

	return result, nil
}

// Close 关闭检测器
func (pd *PoseDetector) Close() {
	if pd.window != nil {
		pd.window.Close()
	}
	pd.classifier.Close()
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
