import { http } from './http'

export class TaskAPI {
  private baseUrl = '/task'

  /**
   * 获取任务统计信息
   * @param timeRange 时间范围：day-日, month-月, year-年
   */
  async getStats(timeRange: string = 'day') {
    return http.get<Api.Task.Stats>(`${this.baseUrl}/stats`, {
      params: { timeRange },
    })
  }

  /**
   * 创建任务
   */
  async create(data: Api.Task.Create) {
    return http.post(this.baseUrl, data)
  }

  /**
   * 更新任务
   */
  async update(id: number, data: Api.Task.Update) {
    return http.put(`${this.baseUrl}/${id}`, data)
  }

  /**
   * 删除任务
   */
  async delete(id: number) {
    return http.delete(`${this.baseUrl}/${id}`)
  }

  /**
   * 获取所有任务
   */
  async findAll(query: { name?: string }) {
    return http.get<Api.Task.Record[]>(`${this.baseUrl}/all`, { params: query })
  }

  /**
   * 执行任务
   */
  async execute(id: number) {
    return http.post(`${this.baseUrl}/${id}/execute?async=true`)
  }

  /**
   * 获取任务日志
   */
  findLogs(query: Api.Task.LogQuery) {
    return http.get<Api.Common.PaginationResponse<Api.Task.Log>>(`${this.baseUrl}/${query.taskId}/logs`, { params: query })
  }

  /**
   * 重置任务状态
   */
  async resetStatus(id: number) {
    return http.put(`${this.baseUrl}/${id}/reset`)
  }

  /**
   * 启用禁用任务
   */
  async toggleEnabled(id: number) {
    return http.put(`${this.baseUrl}/${id}/toggle`)
  }
}

export const taskAPI = new TaskAPI()
