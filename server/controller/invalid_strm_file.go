package controller

import (
	"strconv"

	"github.com/MccRay-s/alist2strm/model/common/response"
	"github.com/MccRay-s/alist2strm/model/invalidstrm/request"
	"github.com/MccRay-s/alist2strm/service"
	"github.com/MccRay-s/alist2strm/utils"
	"github.com/gin-gonic/gin"
)

// 包级别的失效STRM文件控制器实例
var InvalidStrmFile = &InvalidStrmFileController{}

// InvalidStrmFileController 失效STRM文件控制器
type InvalidStrmFileController struct{}

// GetInvalidStrmFileList 获取失效STRM文件列表
// @Summary 获取失效STRM文件列表
// @Description 分页获取失效STRM文件列表，支持多种筛选条件
// @Tags 失效检测
// @Accept json
// @Produce json
// @Param page query int true "页码" minimum(1)
// @Param pageSize query int true "每页数量" minimum(1) maximum(100)
// @Param status query string false "状态筛选" Enums(pending,confirmed,ignored,processing)
// @Param reason query string false "失效原因筛选"
// @Param detectionType query string false "检测类型筛选" Enums(auto,manual)
// @Param keyword query string false "关键字搜索"
// @Param detectionTimeStart query string false "检测时间开始" format(date)
// @Param detectionTimeEnd query string false "检测时间结束" format(date)
// @Success 200 {object} Response{data=PagedResponse} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/invalid-strm-files [get]
func (c *InvalidStrmFileController) GetInvalidStrmFileList(ctx *gin.Context) {
	var req request.InvalidStrmFileListReq

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error("获取失效STRM文件列表参数绑定失败", "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("请求参数错误", ctx)
		return
	}

	// 参数验证
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 调用服务层
	invalidFiles, total, err := service.InvalidStrmFile.GetInvalidStrmFileList(&req)
	if err != nil {
		utils.Error("获取失效STRM文件列表失败", "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("获取失效STRM文件列表失败", ctx)
		return
	}

	// 返回分页数据
	utils.Info("获取失效STRM文件列表成功", "total", total, "request_id", ctx.GetString("request_id"))
	response.SuccessWithData(gin.H{
		"list":     invalidFiles,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	}, ctx)
}

// GetInvalidStrmFileDetail 获取失效STRM文件详情
// @Summary 获取失效STRM文件详情
// @Description 根据ID获取失效STRM文件的详细信息
// @Tags 失效检测
// @Accept json
// @Produce json
// @Param id path int true "失效文件ID"
// @Success 200 {object} Response{data=invalidstrm.InvalidStrmFile} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 404 {object} Response "文件不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/invalid-strm-files/{id} [get]
func (c *InvalidStrmFileController) GetInvalidStrmFileDetail(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error("失效STRM文件ID参数无效", "id", idStr, "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("无效的文件ID", ctx)
		return
	}

	// 调用服务层
	invalidFile, err := service.InvalidStrmFile.GetInvalidStrmFileByID(uint(id))
	if err != nil {
		utils.Error("获取失效STRM文件详情失败", "id", uint(id), "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("文件不存在", ctx)
		return
	}

	utils.Info("获取失效STRM文件详情成功", "id", uint(id), "request_id", ctx.GetString("request_id"))
	response.SuccessWithData(invalidFile, ctx)
}

// GetInvalidStrmFileStatistics 获取失效STRM文件统计信息
// @Summary 获取失效STRM文件统计信息
// @Description 获取失效STRM文件的统计信息，包括状态分布、失效原因等
// @Tags 失效检测
// @Accept json
// @Produce json
// @Param dateStart query string false "统计开始日期" format(date)
// @Param dateEnd query string false "统计结束日期" format(date)
// @Success 200 {object} Response{data=map[string]interface{}} "成功"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/invalid-strm-files/statistics [get]
func (c *InvalidStrmFileController) GetInvalidStrmFileStatistics(ctx *gin.Context) {
	var req request.InvalidStrmFileStatsReq

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error("获取失效STRM文件统计参数绑定失败", "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("请求参数错误", ctx)
		return
	}

	// 调用服务层
	stats, err := service.InvalidStrmFile.GetInvalidStrmFileStatistics(&req)
	if err != nil {
		utils.Error("获取失效STRM文件统计信息失败", "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("获取统计信息失败", ctx)
		return
	}

	utils.Info("获取失效STRM文件统计信息成功", "request_id", ctx.GetString("request_id"))
	response.SuccessWithData(stats, ctx)
}

// BatchProcessInvalidStrmFiles 批量处理失效STRM文件
// @Summary 批量处理失效STRM文件
// @Description 批量确认删除或忽略失效STRM文件
// @Tags 失效检测
// @Accept json
// @Produce json
// @Param request body request.InvalidStrmFileBatchProcessReq true "批量处理请求"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/invalid-strm-files/batch-process [post]
func (c *InvalidStrmFileController) BatchProcessInvalidStrmFiles(ctx *gin.Context) {
	var req request.InvalidStrmFileBatchProcessReq

	// 绑定请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error("批量处理失效STRM文件参数绑定失败", "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("请求参数错误", ctx)
		return
	}

	// 调用服务层
	err := service.InvalidStrmFile.BatchUpdateInvalidStrmFileStatus(&req)
	if err != nil {
		utils.Error("批量处理失效STRM文件失败", "count", len(req.IDs), "error", err.Error(), "request_id", ctx.GetString("request_id"))
		response.FailWithMessage("批量处理失败", ctx)
		return
	}

	utils.Info("批量处理失效STRM文件成功", "count", len(req.IDs), "action", string(req.Action), "request_id", ctx.GetString("request_id"))
	response.SuccessWithMessage("批量处理成功", ctx)
}
