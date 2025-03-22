package alert

import (
	"log"
	"time"
)

// AlertManager 管理提醒的发送
type AlertManager struct {
	lastAlertTime time.Time
}

// NewAlertManager 创建新的提醒管理器
func NewAlertManager() *AlertManager {
	return &AlertManager{
		lastAlertTime: time.Now(),
	}
}

// SendAlert 发送提醒
func (am *AlertManager) SendAlert(message string) {
	// 限制提醒频率，至少间隔5秒
	if time.Since(am.lastAlertTime) < 5*time.Second {
		return
	}

	// TODO: 实现Apple Watch通知
	log.Printf("发送提醒: %s", message)
	am.lastAlertTime = time.Now()
}
