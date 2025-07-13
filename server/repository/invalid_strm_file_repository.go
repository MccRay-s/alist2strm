package repository

import (
	"time"

	"github.com/MccRay-s/alist2strm/database"
	"github.com/MccRay-s/alist2strm/model/invalidstrm"
	invalidStrmRequest "github.com/MccRay-s/alist2strm/model/invalidstrm/request"
)

type InvalidStrmFileRepository struct{}

// 包级别的全局实例
var InvalidStrmFile = &InvalidStrmFileRepository{}

// GetList 获取失效STRM文件列表
func (r *InvalidStrmFileRepository) GetList(req *invalidStrmRequest.InvalidStrmFileListReq) ([]*invalidstrm.InvalidStrmFile, int64, error) {
	db := database.DB
	var invalidFiles []*invalidstrm.InvalidStrmFile
	var total int64

	// 构建查询
	query := db.Model(&invalidstrm.InvalidStrmFile{})

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 失效原因筛选
	if req.Reason != nil {
		query = query.Where("reason = ?", *req.Reason)
	}

	// 检测类型筛选
	if req.DetectionType != nil {
		query = query.Where("detection_type = ?", *req.DetectionType)
	}

	// 关键字搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("file_name LIKE ? OR source_path LIKE ? OR target_file_path LIKE ?", keyword, keyword, keyword)
	}

	// 时间范围筛选
	if req.DetectionTimeStart != nil {
		if startTime, err := time.Parse("2006-01-02", *req.DetectionTimeStart); err == nil {
			query = query.Where("detection_time >= ?", startTime)
		}
	}
	if req.DetectionTimeEnd != nil {
		if endTime, err := time.Parse("2006-01-02", *req.DetectionTimeEnd); err == nil {
			// 结束时间加一天，表示当天23:59:59
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("detection_time < ?", endTime)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("detection_time DESC").Offset(offset).Limit(req.PageSize).Find(&invalidFiles).Error; err != nil {
		return nil, 0, err
	}

	return invalidFiles, total, nil
}

// GetByID 根据ID获取失效STRM文件
func (r *InvalidStrmFileRepository) GetByID(id uint) (*invalidstrm.InvalidStrmFile, error) {
	db := database.DB
	var invalidFile invalidstrm.InvalidStrmFile

	if err := db.First(&invalidFile, id).Error; err != nil {
		return nil, err
	}

	return &invalidFile, nil
}

// Create 创建失效STRM文件记录
func (r *InvalidStrmFileRepository) Create(invalidFile *invalidstrm.InvalidStrmFile) error {
	return database.DB.Create(invalidFile).Error
}

// BatchCreate 批量创建失效STRM文件记录
func (r *InvalidStrmFileRepository) BatchCreate(invalidFiles []*invalidstrm.InvalidStrmFile) error {
	if len(invalidFiles) == 0 {
		return nil
	}
	return database.DB.CreateInBatches(invalidFiles, 100).Error
}

// UpdateByID 根据ID更新失效STRM文件记录
func (r *InvalidStrmFileRepository) UpdateByID(id uint, updateFields map[string]interface{}) error {
	if id == 0 {
		return nil
	}

	return database.DB.Model(&invalidstrm.InvalidStrmFile{}).Where("id = ?", id).Updates(updateFields).Error
}

// BatchUpdateStatus 批量更新状态
func (r *InvalidStrmFileRepository) BatchUpdateStatus(ids []uint, status invalidstrm.InvalidStrmFileStatus, processedBy, processResult string) error {
	if len(ids) == 0 {
		return nil
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"status":         status,
		"processed_at":   &now,
		"processed_by":   processedBy,
		"process_result": processResult,
		"updated_at":     now,
	}

	return database.DB.Model(&invalidstrm.InvalidStrmFile{}).Where("id IN ?", ids).Updates(updateData).Error
}

// GetByFileHistoryID 根据文件历史ID获取失效记录
func (r *InvalidStrmFileRepository) GetByFileHistoryID(fileHistoryID uint) (*invalidstrm.InvalidStrmFile, error) {
	db := database.DB
	var invalidFile invalidstrm.InvalidStrmFile

	if err := db.Where("file_history_id = ?", fileHistoryID).First(&invalidFile).Error; err != nil {
		return nil, err
	}

	return &invalidFile, nil
}

// DeleteByIDs 根据ID列表删除记录
func (r *InvalidStrmFileRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return database.DB.Where("id IN ?", ids).Delete(&invalidstrm.InvalidStrmFile{}).Error
}

// GetStatistics 获取失效文件统计信息
func (r *InvalidStrmFileRepository) GetStatistics(req *invalidStrmRequest.InvalidStrmFileStatsReq) (map[string]interface{}, error) {
	db := database.DB

	// 构建基础查询
	query := db.Model(&invalidstrm.InvalidStrmFile{})

	// 时间范围筛选
	if req.DateStart != nil {
		if startTime, err := time.Parse("2006-01-02", *req.DateStart); err == nil {
			query = query.Where("detection_time >= ?", startTime)
		}
	}
	if req.DateEnd != nil {
		if endTime, err := time.Parse("2006-01-02", *req.DateEnd); err == nil {
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("detection_time < ?", endTime)
		}
	}

	// 统计各种状态的数量
	var stats struct {
		Total      int64 `json:"total"`
		Pending    int64 `json:"pending"`
		Confirmed  int64 `json:"confirmed"`
		Ignored    int64 `json:"ignored"`
		Processing int64 `json:"processing"`
	}

	// 总数
	query.Count(&stats.Total)

	// 各状态数量
	db.Model(&invalidstrm.InvalidStrmFile{}).Where("status = ?", invalidstrm.StatusPending).Count(&stats.Pending)
	db.Model(&invalidstrm.InvalidStrmFile{}).Where("status = ?", invalidstrm.StatusConfirmed).Count(&stats.Confirmed)
	db.Model(&invalidstrm.InvalidStrmFile{}).Where("status = ?", invalidstrm.StatusIgnored).Count(&stats.Ignored)
	db.Model(&invalidstrm.InvalidStrmFile{}).Where("status = ?", invalidstrm.StatusProcessing).Count(&stats.Processing)

	// 失效原因统计
	var reasonStats []struct {
		Reason string `json:"reason"`
		Count  int64  `json:"count"`
	}
	db.Model(&invalidstrm.InvalidStrmFile{}).
		Select("reason, COUNT(*) as count").
		Group("reason").
		Scan(&reasonStats)

	result := map[string]interface{}{
		"statusStats": stats,
		"reasonStats": reasonStats,
	}

	return result, nil
}
