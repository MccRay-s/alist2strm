package main

import (
	"github.com/MccRay-s/alist2strm/controller"
	"github.com/MccRay-s/alist2strm/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.RequestID())    // 请求ID中间件
	r.Use(middleware.AccessLogger()) // 访问日志中间件
	r.Use(gin.Recovery())            // 错误恢复中间件

	// API 路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "alist2strm server is running",
			})
		})

		// 公开路由（不需要认证）
		public := api.Group("/user")
		{
			public.POST("/login", controller.User.Login)       // 用户登录
			public.POST("/register", controller.User.Register) // 用户注册
		}

		// Emby 图片公开路由（不需要认证）
		api.GET("/emby/items/:item_id/images/:image_type", controller.Emby.GetImage) // 获取Emby图片

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.JWTAuth()) // 应用JWT认证中间件
		{
			// 用户相关路由
			user := auth.Group("/user")
			{
				user.GET("/me", controller.User.Me)            // 获取当前用户信息
				user.GET("/:id", controller.User.GetUserInfo)  // 获取指定用户信息
				user.PUT("/:id", controller.User.UpdateUser)   // 更新用户信息
				user.GET("/list", controller.User.GetUserList) // 获取用户列表
			}

			// 配置相关路由
			config := auth.Group("/config")
			{
				config.POST("/", controller.Config.Create)                   // 创建配置
				config.GET("/:id", controller.Config.GetConfigInfo)          // 获取指定配置信息
				config.GET("/code/:code", controller.Config.GetConfigByCode) // 根据代码获取配置
				config.PUT("/:id", controller.Config.UpdateConfig)           // 更新配置信息
				config.DELETE("/:id", controller.Config.DeleteConfig)        // 删除配置
				config.GET("/list", controller.Config.GetConfigList)         // 获取配置列表
			}

			// 任务相关路由
			task := auth.Group("/task")
			{
				task.POST("/", controller.Task.Create)                     // 创建任务
				task.GET("/:id", controller.Task.GetTaskInfo)              // 获取指定任务信息
				task.PUT("/:id", controller.Task.UpdateTask)               // 更新任务信息
				task.DELETE("/:id", controller.Task.DeleteTask)            // 删除任务
				task.GET("/list", controller.Task.GetTaskList)             // 获取任务列表（分页）
				task.GET("/all", controller.Task.GetAllTasks)              // 获取所有任务（不分页）
				task.GET("/stats", controller.Task.GetTaskStats)           // 获取任务统计数据
				task.PUT("/:id/toggle", controller.Task.ToggleTaskEnabled) // 切换任务启用状态
				task.PUT("/:id/reset", controller.Task.ResetTaskStatus)    // 重置任务运行状态
				task.POST("/:id/execute", controller.Task.ExecuteTask)     // 执行任务（支持同步/异步）
			}

			// 任务日志相关路由
			taskLog := auth.Group("/task-log")
			{
				taskLog.GET("/:id", controller.TaskLogControllerInstance.GetTaskLogInfo)                      // 获取指定任务日志信息
				taskLog.GET("/", controller.TaskLogControllerInstance.GetTaskLogList)                         // 获取任务日志列表（分页）
				taskLog.GET("/stats/processing", controller.TaskLogControllerInstance.GetFileProcessingStats) // 获取文件处理统计数据
			}

			// 文件历史相关路由
			fileHistory := auth.Group("/file-history")
			{
				fileHistoryController := &controller.FileHistoryController{}
				fileHistory.GET("/", fileHistoryController.GetFileList)           // 获取主文件分页列表
				fileHistory.GET("/:id", fileHistoryController.GetFileHistoryInfo) // 获取文件历史详情
			}

			// 失效STRM文件检测相关路由
			invalidStrmFile := auth.Group("/invalid-strm-files")
			{
				invalidStrmFile.GET("/", controller.InvalidStrmFile.GetInvalidStrmFileList)                     // 获取失效STRM文件列表
				invalidStrmFile.GET("/statistics", controller.InvalidStrmFile.GetInvalidStrmFileStatistics)     // 获取统计信息
				invalidStrmFile.GET("/:id", controller.InvalidStrmFile.GetInvalidStrmFileDetail)                // 获取失效STRM文件详情
				invalidStrmFile.POST("/batch-process", controller.InvalidStrmFile.BatchProcessInvalidStrmFiles) // 批量处理失效文件
			}

			// AList 相关路由
			alist := auth.Group("/alist")
			{
				alist.POST("/test", controller.AList.TestConnection) // 测试AList连接
			}

			// Emby 相关需认证路由
			emby := auth.Group("/emby")
			{
				emby.GET("/test", controller.Emby.TestConnection)                    // 测试Emby服务可用性
				emby.GET("/libraries", controller.Emby.GetLibraries)                 // 获取Emby媒体库列表
				emby.GET("/latest", controller.Emby.GetLatestMedia)                  // 获取Emby最新入库列表
				emby.POST("/libraries/:id/refresh", controller.Emby.RefreshLibrary)  // 刷新指定媒体库
				emby.POST("/libraries/refresh", controller.Emby.RefreshAllLibraries) // 刷新所有媒体库
			}
		}

	}

	return r
}
