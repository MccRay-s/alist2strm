import { TaskLog } from '@/models/task-log.js'
import type { WhereOptions } from 'sequelize'
import { Op } from 'sequelize'
import { logger } from '@/utils/logger.js'


export class TaskLogService {
  /**
   * 创建任务日志
   */
  async create(data: App.Task.LogCreate): Promise<TaskLog> {
    try {
      const taskLog = await TaskLog.create(data)
      logger.info.info('创建任务日志成功:', { id: taskLog.id, taskId: taskLog.taskId })
      return taskLog
    }
    catch (error) {
      logger.error.error('创建任务日志失败:', error)
      throw error
    }
  }

  /**
   * 更新任务日志
   */
  async update(id: number, data: App.Task.LogUpdate): Promise<TaskLog | null> {
    try {
      const taskLog = await TaskLog.findByPk(id)
      if (!taskLog) {
        logger.warn.warn('更新任务日志失败: 日志不存在', { id })
        throw new Error(`日志【¥${id}】不存在`)
      }
      await taskLog.update(data)
      logger.info.info('更新任务日志成功:', { id, taskId: taskLog.taskId })
      return taskLog
    }
    catch (error) {
      logger.error.error('更新任务日志失败:', error)
      throw error
    }
  }

  /**
   * 分页查询任务日志
   */
  async findByPage(query: App.Task.LogQuery): Promise<App.Common.PaginationResult<TaskLog>> {
    try {
      const { page = 1, pageSize = 10, taskId, status, startTime, endTime } = query
      const where: WhereOptions<Models.TaskLogAttributes> = {}

      if (taskId)
        where.taskId = taskId

      if (status)
        where.status = status

      if (startTime || endTime) {
        where.startTime = {}
        if (startTime)
          Object.assign(where.startTime, { [Op.gte]: startTime })
        if (endTime)
          Object.assign(where.startTime, { [Op.lte]: endTime })
      }

      const { count, rows } = await TaskLog.findAndCountAll({
        where,
        offset: (page - 1) * pageSize,
        limit: pageSize,
        order: [['startTime', 'DESC']],
      })

      logger.debug.debug('分页查询任务日志:', { page, pageSize, total: count })
      return {
        list: rows,
        total: count,
        page,
        pageSize,
      }
    }
    catch (error) {
      logger.error.error('分页查询任务日志失败:', error)
      throw error
    }
  }

  /**
   * 根据任务ID查询最新日志
   */
  async findLatestByTaskId(taskId: number): Promise<TaskLog | null> {
    try {
      const taskLog = await TaskLog.findOne({
        where: { taskId },
        order: [['startTime', 'DESC']],
      })
      return taskLog
    }
    catch (error) {
      logger.error.error('查询任务最新日志失败:', error)
      throw error
    }
  }

  /**
   * 根据ID查询任务日志
   */
  async findById(id: number): Promise<TaskLog | null> {
    try {
      const taskLog = await TaskLog.findByPk(id)
      if (!taskLog)
        logger.warn.warn('查询任务日志失败: 日志不存在', { id })
      return taskLog
    }
    catch (error) {
      logger.error.error('查询任务日志失败:', error)
      throw error
    }
  }

  /**
   * 根据任务ID查询日志
   */
  async findByTaskId(taskId: number): Promise<TaskLog[]> {
    try {
      const logs = await TaskLog.findAll({
        where: { taskId },
        order: [['startTime', 'DESC']],
      })
      return logs
    }
    catch (error) {
      logger.error.error('查询任务日志失败:', error)
      throw error
    }
  }
}

// 导出单例实例
export const taskLogService = new TaskLogService()