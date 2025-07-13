package service

import (
	"github.com/MccRay-s/alist2strm/model/invalidstrm"
	invalidStrmRequest "github.com/MccRay-s/alist2strm/model/invalidstrm/request"
	"github.com/MccRay-s/alist2strm/repository"
	"github.com/MccRay-s/alist2strm/utils"
)

// 包级别的失效STRM文件服务实例
var InvalidStrmFile = &InvalidStrmFileService{}

// InvalidStrmFileService 失效STRM文件服务
type InvalidStrmFileService struct{}

// GetInvalidStrmFileList 获取失效STRM文件分页列表
func (s *InvalidStrmFileService) GetInvalidStrmFileList(req *invalidStrmRequest.InvalidStrmFileListReq) ([]*invalidstrm.InvalidStrmFile, int64, error) {
	utils.Debug("获取失效STRM文件列表", "page", req.Page, "pageSize", req.PageSize, "keyword", req.Keyword)

	// 调用 repository 层获取数据
	invalidFiles, total, err := repository.InvalidStrmFile.GetList(req)
	if err != nil {
		utils.Error("获取失效STRM文件列表失败", "error", err.Error())
		return nil, 0, err
	}

	utils.Debug("获取失效STRM文件列表成功", "total", total, "count", len(invalidFiles))

	return invalidFiles, total, nil
}

// GetInvalidStrmFileByID 根据ID获取失效STRM文件详情
func (s *InvalidStrmFileService) GetInvalidStrmFileByID(id uint) (*invalidstrm.InvalidStrmFile, error) {
	utils.Debug("获取失效STRM文件详情", "id", id)

	invalidFile, err := repository.InvalidStrmFile.GetByID(id)
	if err != nil {
		utils.Error("获取失效STRM文件详情失败", "id", id, "error", err.Error())
		return nil, err
	}

	utils.Debug("获取失效STRM文件详情成功", "id", id, "fileName", invalidFile.FileName)

	return invalidFile, nil
}

// GetInvalidStrmFileStatistics 获取失效STRM文件统计信息
func (s *InvalidStrmFileService) GetInvalidStrmFileStatistics(req *invalidStrmRequest.InvalidStrmFileStatsReq) (map[string]interface{}, error) {
	utils.Debug("获取失效STRM文件统计信息")

	stats, err := repository.InvalidStrmFile.GetStatistics(req)
	if err != nil {
		utils.Error("获取失效STRM文件统计信息失败", "error", err.Error())
		return nil, err
	}

	utils.Debug("获取失效STRM文件统计信息成功")

	return stats, nil
}

// BatchUpdateInvalidStrmFileStatus 批量更新失效STRM文件状态
func (s *InvalidStrmFileService) BatchUpdateInvalidStrmFileStatus(req *invalidStrmRequest.InvalidStrmFileBatchProcessReq) error {
	utils.Info("批量更新失效STRM文件状态", "count", len(req.IDs), "action", string(req.Action))

	// 这里可以添加业务逻辑验证
	// 例如：检查状态转换是否合法、权限验证等

	err := repository.InvalidStrmFile.BatchUpdateStatus(req.IDs, req.Action, "admin", req.Reason)
	if err != nil {
		utils.Error("批量更新失效STRM文件状态失败", "count", len(req.IDs), "error", err.Error())
		return err
	}

	utils.Info("批量更新失效STRM文件状态成功", "count", len(req.IDs), "action", string(req.Action))

	return nil
}

// DeleteInvalidStrmFiles 删除失效STRM文件记录
func (s *InvalidStrmFileService) DeleteInvalidStrmFiles(ids []uint) error {
	utils.Info("删除失效STRM文件记录", "count", len(ids))

	err := repository.InvalidStrmFile.DeleteByIDs(ids)
	if err != nil {
		utils.Error("删除失效STRM文件记录失败", "count", len(ids), "error", err.Error())
		return err
	}

	utils.Info("删除失效STRM文件记录成功", "count", len(ids))

	return nil
}
