import { FileHistory } from '@/models/file-history.js'
import type { WhereOptions } from 'sequelize'
import { Op } from 'sequelize'
import { logger } from '@/utils/logger.js'


export class FileHistoryService {
  /**
   * 创建文件历史
   */
  async create(data:App.FileHistory.Create): Promise<FileHistory> {
    try {
      const fileHistory = await FileHistory.create(data as any)
      logger.info.info('创建文件历史成功:', { id: fileHistory.id, fileName: fileHistory.fileName })
      return fileHistory
    }
    catch (error) {
      logger.error.error('创建文件历史失败:', error)
      throw error
    }
  }

  /**
   * 分页查询文件历史
   */
  async findByPage(query: App.FileHistory.Query): Promise<App.Common.PaginationResult<FileHistory>> {
    try {
      const { page = 1, pageSize = 10, keyword, fileType, fileSuffix, startTime, endTime } = query
      const where: WhereOptions<Models.FileHistoryAttributes> = {}

      // 关键字搜索
      if (keyword) {
        Object.assign(where, {
          [Op.or]: [
            { fileName: { [Op.like]: `%${keyword}%` } },
            { sourcePath: { [Op.like]: `%${keyword}%` } },
            { targetFilePath: { [Op.like]: `%${keyword}%` } },
          ],
        } as WhereOptions<Models.FileHistoryAttributes>)
      }

      // 文件类型过滤
      if (fileType)
        where.fileType = fileType

      // 文件后缀过滤
      if (fileSuffix)
        where.fileSuffix = fileSuffix

      // 时间范围过滤
      if (startTime || endTime) {
        where.createdAt = {}
        if (startTime)
          Object.assign(where.createdAt, { [Op.gte]: startTime })
        if (endTime)
          Object.assign(where.createdAt, { [Op.lte]: endTime })
      }

      const { count, rows } = await FileHistory.findAndCountAll({
        where,
        offset: (page - 1) * pageSize,
        limit: pageSize,
        order: [['createdAt', 'DESC']],
      })

      logger.debug.debug('分页查询文件历史:', { page, pageSize, total: count })
      return {
        list: rows,
        total: count,
        page,
        pageSize,
      }
    }
    catch (error) {
      logger.error.error('分页查询文件历史失败:', error)
      throw error
    }
  }

  /**
   * 根据ID查询文件历史
   */
  async findById(id: number): Promise<FileHistory | null> {
    try {
      const fileHistory = await FileHistory.findByPk(id)
      if (!fileHistory)
        logger.warn.warn('查询文件历史失败: 历史记录不存在', { id })
      return fileHistory
    }
    catch (error) {
      logger.error.error('查询文件历史失败:', error)
      throw error
    }
  }

  /**
   * 检查文件是否已存在
   */
  async checkFileExists(sourcePath: string, fileName: string): Promise<boolean> {
    try {
      const count = await FileHistory.count({
        where: {
          sourcePath,
          fileName,
        },
      })
      return count > 0
    }
    catch (error) {
      logger.error.error('检查文件是否存在失败:', error)
      throw error
    }
  }

  /**
   * 批量删除文件历史记录
   * @param ids ID列表
   * @returns 删除的记录数
   */
  async bulkDelete(ids: number[]): Promise<number> {
    return await FileHistory.bulkDelete(ids)
  }

  /**
   * 清空所有文件历史记录
   * @returns 删除的记录数
   */
  async clearAll(): Promise<number> {
    return await FileHistory.clearAll()
  }
}

export const fileHistoryService = new FileHistoryService() 