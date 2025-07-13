package request

import (
	"github.com/MccRay-s/alist2strm/model/invalidstrm"
)

// InvalidStrmFileListReq 失效STRM文件列表请求
type InvalidStrmFileListReq struct {
	Page     int `json:"page" form:"page" binding:"required,min=1"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,min=1"` // 每页数量

	// 筛选条件
	Status        *invalidstrm.InvalidStrmFileStatus `json:"status" form:"status"`               // 状态筛选
	Reason        *invalidstrm.InvalidStrmFileReason `json:"reason" form:"reason"`               // 失效原因筛选
	DetectionType *string                            `json:"detectionType" form:"detectionType"` // 检测类型筛选
	Keyword       string                             `json:"keyword" form:"keyword"`             // 关键字搜索（文件名、路径）

	// 时间范围筛选
	DetectionTimeStart *string `json:"detectionTimeStart" form:"detectionTimeStart"` // 检测时间开始
	DetectionTimeEnd   *string `json:"detectionTimeEnd" form:"detectionTimeEnd"`     // 检测时间结束
}

// InvalidStrmFileBatchProcessReq 批量处理失效STRM文件请求
type InvalidStrmFileBatchProcessReq struct {
	IDs    []uint                            `json:"ids" binding:"required,min=1"` // 要处理的文件ID列表
	Action invalidstrm.InvalidStrmFileStatus `json:"action" binding:"required"`    // 处理动作：confirmed(确认删除), ignored(忽略)
	Reason string                            `json:"reason"`                       // 处理原因说明
}

// InvalidStrmFileDetectionReq 启动失效检测请求
type InvalidStrmFileDetectionReq struct {
	DetectionType string   `json:"detectionType" binding:"required"` // 检测类型：auto, manual
	FileIDs       []uint   `json:"fileIds"`                          // 指定检测的文件历史ID列表（可选，为空则检测所有STRM文件）
	CheckTypes    []string `json:"checkTypes"`                       // 检测类型：file_exists(文件存在性), url_valid(URL有效性), http_access(HTTP访问)
}

// InvalidStrmFileStatsReq 失效文件统计请求
type InvalidStrmFileStatsReq struct {
	// 可以添加时间范围等筛选条件
	DateStart *string `json:"dateStart" form:"dateStart"` // 统计开始日期
	DateEnd   *string `json:"dateEnd" form:"dateEnd"`     // 统计结束日期
}
