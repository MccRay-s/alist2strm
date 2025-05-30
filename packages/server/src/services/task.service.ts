import { Task } from '@/models/task.js'
import type { WhereOptions } from 'sequelize'
import { Op } from 'sequelize'
import { logger } from '@/utils/logger.js'
import { EventEmitter } from 'node:events'
import { taskQueue } from './task-queue.service.js'
import { taskScheduler } from './task-scheduler.service.js'

// 任务执行器事件发射器
const taskEmitter = new EventEmitter()

export class TaskService {
  /**
   * 创建任务
   */
  async create(data: App.Task.Create): Promise<Task> {
    try {
      const task = await Task.create({
        ...data,
        running: false,
        lastRunAt: null,
      })
      logger.info.info('创建任务成功:', { id: task.id, name: task.name })
      
      // 如果任务启用且有 cron 表达式，则调度任务
      if (task.enabled && task.cron) {
        await taskScheduler.scheduleTask(task)
      }
      return task
    }
    catch (error) {
      logger.error.error('创建任务失败:', error)
      throw error
    }
  }

  /**
   * 更新任务
   */
  async update(id: number, data: App.Task.Update): Promise<Task | null> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('更新任务失败: 任务不存在', { id })
        throw new Error(`任务【¥${id}】不存在`)
      }
      await task.update(data)
      logger.info.info('更新任务成功:', { id, name: task.dataValues.name })
      // 如果任务启用且有 cron 表达式，则调度任务
      if (task.dataValues.enabled && task.dataValues.cron) {
        await taskScheduler.scheduleTask(task)
      }
      else {
        taskScheduler.unscheduleTask(id)
      }
      return task
    }
    catch (error) {
      logger.error.error('更新任务失败:', error)
      throw error
    }
  }

  /**
   * 删除任务
   */
  async delete(id: number): Promise<boolean> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('删除任务失败: 任务不存在', { id })
        throw new Error(`任务【¥${id}】不存在`)
      }

      // 取消任务调度
      taskScheduler.unscheduleTask(id)
      
      await task.destroy()
      logger.info.info('删除任务成功:', { id, name: task.name })
      return true
    }
    catch (error) {
      logger.error.error('删除任务失败:', error)
      throw error
    }
  }

  /**
   * 分页查询任务
   */
  async findByPage(query: App.Task.Query): Promise<App.Common.PaginationResult<Task>> {
    try {
      const { page = 1, pageSize = 10, keyword, enabled, running, sortBy, sortOrder} = query
      const where: WhereOptions<Models.TaskAttributes> = {}

      // 关键字搜索
      if (keyword) {
        Object.assign(where, {
          [Op.or]: [
            { name: { [Op.like]: `%${keyword}%` } },
            { sourcePath: { [Op.like]: `%${keyword}%` } },
            { targetPath: { [Op.like]: `%${keyword}%` } },
          ],
        } as WhereOptions<Models.TaskAttributes>)
      }
      if (typeof enabled === 'boolean')
        where.enabled = enabled
      if (typeof running === 'boolean')
        where.running = running

      const { count, rows } = await Task.findAndCountAll({
        where,
        offset: (page - 1) * pageSize,
        limit: pageSize,
        order: [[sortBy || 'createdAt', sortOrder ||'DESC']],
      })

      logger.debug.debug('分页查询任务:', { page, pageSize, total: count })
      return {
        list: rows,
        total: count,
        page,
        pageSize,
      }
    }
    catch (error) {
      logger.error.error('分页查询任务失败:', error)
      throw error
    }
  }

  /**
   * 查询所有任务
   */
  async findAll(name?: string): Promise<Task[]> {
    try {
      const where: WhereOptions<Models.TaskAttributes> = {}
      if (name) {
        where.name = {
          [Op.like]: `%${name}%`,
        }
      }
      const tasks = await Task.findAll({
        where,
        order: [['createdAt', 'DESC']],
      })
      logger.debug.debug('查询所有任务:', { total: tasks.length })
      return tasks
    }
    catch (error) {
      logger.error.error('查询所有任务失败:', error)
      throw error
    }
  }

  /**
   * 根据ID查询任务
   */
  async findById(id: number): Promise<Task | null> {
    try {
      const task = await Task.findByPk(id)
      if (!task){
        logger.warn.warn('查询任务失败: 任务不存在', { id })
        return null
      }
      return task
    }
    catch (error) {
      logger.error.error('查询任务失败:', error)
      throw error
    }
  }

  /**
   * 更新任务运行状态
   */
  async updateRunningStatus(id: number, running: boolean): Promise<Task | null> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('更新任务运行状态失败: 任务不存在', { id })
        return null
      }
      await task.update({
        running,
        lastRunAt: running ? new Date() : task.lastRunAt,
      })
      logger.info.info('更新任务运行状态成功:', { id, name: task.name, running })
      return task
    }
    catch (error) {
      logger.error.error('更新任务运行状态失败:', error)
      throw error
    }
  }

  /**
   * 执行任务
   */
  async execute(id: number): Promise<boolean> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('执行任务失败: 任务不存在', { id })
        throw new Error(`任务【¥${id}】不存在`)
      }

      if (task.running) {
        logger.warn.warn('执行任务失败: 任务正在运行', { id })
        throw new Error(`任务【¥${id}】正在运行`)
      }
      return await taskQueue.addTask(id)
    }
    catch (error) {
      logger.error.error('执行任务失败:', error)
      throw error
    }
  }

  /**
   * 停止任务
   */
  async stop(id: number): Promise<boolean> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('停止任务失败: 任务不存在', { id })
        throw new Error(`任务【¥${id}】不存在`)
      }

      if (!task.running) {
        return true
      }
      return await taskQueue.stopTask(id)
    }
    catch (error) {
      logger.error.error('停止任务失败:', error)
      throw error
    }
  }

  /**
   * 重置任务状态
   */
  async resetStatus(id: number): Promise<boolean> {
    try {
      const task = await Task.findByPk(id)
      if (!task) {
        logger.warn.warn('重置任务状态失败: 任务不存在', { id })
        throw new Error(`任务【¥${id}】不存在`)
      }
      // 如何检测是否在运行？ 不使用状态的情况
      const jobInfo = await taskQueue.getTaskStatus(id)
      if (jobInfo && jobInfo.status && ['pending','running'].includes(jobInfo.status)) {
        logger.warn.warn('重置任务状态失败: 任务正在运行中！', { id })
        throw new Error(`任务【¥${id}】正在运行`)
      }
      
      await task.update({
        running: false,
        lastRunAt: null,
      })
      logger.info.info('重置任务状态成功:', { id, name: task.name })
      return true
    }
    catch (error) {
      logger.error.error('重置任务状态失败:', error) 
      throw error
    }
  }

  /**
   * 获取任务状态
   */
  getTaskStatus(id: number) {
    return taskQueue.getTaskStatus(id)
  }

  /**
   * 订阅任务进度
   */
  onProgress(id: number, callback: (progress: number) => void): void {
    taskEmitter.on(`task:${id}:progress`, callback)
  }

  /**
   * 取消订阅任务进度
   */
  offProgress(id: number, callback: (progress: number) => void): void {
    taskEmitter.off(`task:${id}:progress`, callback)
  }
}

export const taskService = new TaskService() 