package invalidstrm

import (
	"time"
)

// InvalidStrmFileStatus 失效文件状态枚举
type InvalidStrmFileStatus string

const (
	StatusPending    InvalidStrmFileStatus = "pending"    // 待处理
	StatusConfirmed  InvalidStrmFileStatus = "confirmed"  // 已确认删除
	StatusIgnored    InvalidStrmFileStatus = "ignored"    // 用户选择忽略
	StatusProcessing InvalidStrmFileStatus = "processing" // 处理中
)

// InvalidStrmFileReason 失效原因枚举
type InvalidStrmFileReason string

const (
	ReasonFileNotFound     InvalidStrmFileReason = "file_not_found"      // 源文件不存在
	ReasonStrmFileNotFound InvalidStrmFileReason = "strm_file_not_found" // STRM文件不存在
	ReasonURLInvalid       InvalidStrmFileReason = "url_invalid"         // URL无效
	ReasonAccessDenied     InvalidStrmFileReason = "access_denied"       // 访问被拒绝
	ReasonServerError      InvalidStrmFileReason = "server_error"        // 服务器错误
	ReasonNetworkError     InvalidStrmFileReason = "network_error"       // 网络错误
)

// InvalidStrmFile 失效STRM文件模型
type InvalidStrmFile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// 关联的文件历史记录
	FileHistoryID uint `json:"fileHistoryId" gorm:"not null;index;uniqueIndex:idx_file_history_detection"`

	// 检测信息
	DetectionTime time.Time             `json:"detectionTime" gorm:"not null;index"`                             // 检测时间
	DetectionType string                `json:"detectionType" gorm:"not null;type:varchar(50)"`                  // 检测类型：auto(自动), manual(手动)
	Reason        InvalidStrmFileReason `json:"reason" gorm:"not null;type:varchar(50)"`                         // 失效原因
	ErrorMessage  string                `json:"errorMessage" gorm:"type:text"`                                   // 错误详情
	Status        InvalidStrmFileStatus `json:"status" gorm:"not null;type:varchar(20);index;default:'pending'"` // 处理状态

	// 文件信息快照（冗余存储，便于查询和显示）
	FileName       string `json:"fileName" gorm:"not null;type:varchar(200);index"`
	SourcePath     string `json:"sourcePath" gorm:"not null;type:varchar(500);index"`
	TargetFilePath string `json:"targetFilePath" gorm:"not null;type:varchar(500);index"`
	FileSize       int64  `json:"fileSize" gorm:"not null"`

	// STRM 相关信息
	StrmURL        string `json:"strmUrl" gorm:"type:text"`    // STRM文件中的URL
	HttpStatusCode *int   `json:"httpStatusCode" gorm:"index"` // HTTP响应状态码（如果是网络检测）

	// 处理信息
	ProcessedAt   *time.Time `json:"processedAt" gorm:"index"`             // 处理时间
	ProcessedBy   string     `json:"processedBy" gorm:"type:varchar(100)"` // 处理人（用户名或系统）
	ProcessResult string     `json:"processResult" gorm:"type:text"`       // 处理结果说明
}

// TableName 表名
func (InvalidStrmFile) TableName() string {
	return "invalid_strm_files"
}

// GetReasonDescription 获取失效原因的中文描述
func (r InvalidStrmFileReason) GetDescription() string {
	switch r {
	case ReasonFileNotFound:
		return "源文件不存在"
	case ReasonStrmFileNotFound:
		return "STRM文件不存在"
	case ReasonURLInvalid:
		return "URL格式无效"
	case ReasonAccessDenied:
		return "访问被拒绝"
	case ReasonServerError:
		return "服务器错误"
	case ReasonNetworkError:
		return "网络连接错误"
	default:
		return "未知错误"
	}
}

// GetStatusDescription 获取状态的中文描述
func (s InvalidStrmFileStatus) GetDescription() string {
	switch s {
	case StatusPending:
		return "待处理"
	case StatusConfirmed:
		return "已确认删除"
	case StatusIgnored:
		return "已忽略"
	case StatusProcessing:
		return "处理中"
	default:
		return "未知状态"
	}
}
