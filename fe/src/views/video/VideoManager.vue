<template>
  <div class="video-manager">
    <!-- 查询表单 -->
    <el-form :inline="true" :model="queryForm" class="query-form">
      <el-form-item label="车间">
        <el-select v-model="queryForm.workshopId" placeholder="请选择车间" clearable style="width: 240px;">
          <el-option
            v-for="item in workshops"
            :key="item.id"
            :value="item.id"
            :label="item.name"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="时间范围" required>
        <el-row :gutter="10">
          <el-col :span="11">
            <el-date-picker
              v-model="queryForm.startTime"
              type="datetime"
              placeholder="开始时间"
              style="width: 100%"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-col>
          <el-col :span="2" class="text-center">
            <span class="separator">至</span>
          </el-col>
          <el-col :span="11">
            <el-date-picker
              v-model="queryForm.endTime"
              type="datetime"
              placeholder="结束时间"
              style="width: 100%"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 操作按钮 -->
    <div class="operation-container">
      <el-button type="danger" @click="handleBatchDelete" :disabled="!selectedVideos.length">
        批量删除
      </el-button>
    </div>

    <!-- 视频列表 -->
    <el-table
      v-loading="loading"
      :data="videoList"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="captureId" label="采集ID" width="100" />
      <el-table-column prop="fileName" label="文件名" />
      <el-table-column prop="workshop.name" label="车间" />
      <el-table-column label="文件大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.fileSize) }}
        </template>
      </el-table-column>
      <el-table-column label="开始时间" width="240">
    <template #default="{ row }">
      {{ formatDateTime(row.startTime) }}
    </template>
  </el-table-column>
  <el-table-column label="结束时间" width="240">
    <template #default="{ row }">
      {{ formatDateTime(row.endTime) }}
    </template>
  </el-table-column>
      <el-table-column prop="duration" label="时长" width="100">
        <template #default="{ row }">
          {{ formatDuration(row.duration) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="text" @click="handlePreview(row)">预览</el-button>
          <el-button type="text" @click="handleDownload(row)">下载</el-button>
          <el-button type="text" @click="handleDelete(row)" class="delete-btn">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      class="pagination"
      :current-page="queryForm.page"
      :page-size="queryForm.pageSize"
      :total="total"
      layout="total, prev, pager, next, jumper"
      @current-change="handlePageChange"
    />

    <!-- 视频预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="视频预览"
      width="80%"
      :destroy-on-close="true"
    >
      <video-player
        v-if="previewVisible"
        :src="previewUrl"
        :options="{
          autoplay: false,
          controls: true,
        }"
        @error="handleVideoError"
      />
    </el-dialog>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex'
import VideoPlayer from './VideoPlayer.vue'
import { 
  getWorkshopList, 
  downloadVideo, 
  startRecording, 
  stopRecording 
} from '@/api/video'
import { formatDateTime } from '@/utils/format'

export default {
  name: 'VideoManager',
  
  components: {
    VideoPlayer
  },
  
  data() {
    return {
      queryForm: {
        page: 1,
        pageSize: 10,
        workshopId: null,
        startTime: '',
        endTime: ''
      },
      timeRange: [],
      defaultTime: ['00:00:00', '23:59:59'],
      workshops: [],
      selectedVideos: [],
      loading: false,
      previewVisible: false,
      previewUrl: '',
      isRecording: false,
      selectedWorkshop: null,
      shortcuts: [
        {
          text: '最近一小时',
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000)
            return [start, end]
          }
        },
        {
          text: '今天',
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setHours(0, 0, 0, 0)
            return [start, end]
          }
        },
        {
          text: '昨天',
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            start.setHours(0, 0, 0, 0)
            end.setTime(end.getTime() - 3600 * 1000 * 24)
            end.setHours(23, 59, 59, 999)
            return [start, end]
          }
        },
        {
          text: '最近一周',
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [start, end]
          }
        },
        {
          text: '最近一个月',
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setMonth(start.getMonth() - 1)
            return [start, end]
          }
        }
      ],
      formatDateTime,
    }
  },
  
  computed: {
    ...mapState('video', ['videoList', 'total'])
  },
  
  created() {
    this.fetchWorkshops()
    this.fetchVideos()
  },
  
  methods: {
    ...mapActions('video', ['getVideoList', 'deleteVideo', 'batchDeleteVideos']),
    
    // 获取车间列表
    async fetchWorkshops() {
      try {
        const { data } = await getWorkshopList()
        this.workshops = data
      } catch (error) {
        this.$message.error('获取车间列表失败')
      }
    },
    
    // 获取视频列表
    async fetchVideos() {
      this.loading = true
      try {
        const params = {
          ...this.queryForm,
          startTime: this.queryForm.startTime,
          endTime: this.queryForm.endTime
        }
        await this.getVideoList(params)
      } catch (error) {
        this.$message.error('获取视频列表失败')
      } finally {
        this.loading = false
      }
    },
    
    // 处理查询
    handleSearch() {
      this.queryForm.page = 1
      this.fetchVideos()
    },
    
    // 重置查询
    resetQuery() {
      this.queryForm = {
        page: 1,
        pageSize: 10,
        workshopId: null,
        startTime: '',
        endTime: ''
      }
      this.timeRange = []
      this.fetchVideos()
    },
    
    // 处理分页
    handlePageChange(page) {
      this.queryForm.page = page
      this.fetchVideos()
    },
    
    // 处理选择
    handleSelectionChange(selection) {
      this.selectedVideos = selection
    },
    
    // 处理预览
    async handlePreview(row) {
      console.log('API URL:', process.env.VUE_APP_API_URL)
      const baseUrl = process.env.VUE_APP_API_URL
      const cleanPath = row.filePath.replace(/^storage[/\\]/, '').replace(/\\/g, '/')
      const encodedPath = encodeURIComponent(cleanPath)
      
      this.previewUrl = `${baseUrl}/storage/${encodedPath}`
      this.previewVisible = true
    },
    
    // 处理下载
    async handleDownload(row) {
      try {
        const blob = await downloadVideo(row.id)
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = row.fileName
        link.click()
        window.URL.revokeObjectURL(url)
      } catch (error) {
        this.$message.error('下载失败')
      }
    },
    
    // 处理删除
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该视频?', '提示', {
          type: 'warning'
        })
        await this.deleteVideo(row.id)
        this.$message.success('删除成功')
        this.fetchVideos()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    
    // 处理批量删除
    async handleBatchDelete() {
      try {
        await this.$confirm(`确认删除选中的 ${this.selectedVideos.length} 个视频?`, '提示', {
          type: 'warning'
        })
        const ids = this.selectedVideos.map(item => item.id)
        await this.batchDeleteVideos(ids)
        this.$message.success('删除成功')
        this.fetchVideos()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    
    // 开始录制
    async handleStartRecord() {
      try {
        await startRecording({
          workshopId: this.selectedWorkshop
        })
        this.isRecording = true
        this.$message.success('开始录制')
      } catch (error) {
        this.$message.error('开始录制失败')
      }
    },
    
    // 停止录制
    async handleStopRecord() {
      try {
        await stopRecording(this.selectedWorkshop)
        this.isRecording = false
        this.$message.success('停止录制')
        this.fetchVideos()
      } catch (error) {
        this.$message.error('停止录制失败')
      }
    },
    
    // 格式化时长
    formatDuration(seconds) {
      const h = Math.floor(seconds / 3600)
      const m = Math.floor((seconds % 3600) / 60)
      const s = Math.floor(seconds % 60)
      return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
    },
    
    // 格式化文件大小
    formatFileSize(bytes) {
      if (!bytes) return '0 MB'
      const mb = bytes / (1024 * 1024)
      return `${mb.toFixed(2)} MB`
    },
    
    // 处理视频错误
    handleVideoError(error) {
      this.$message.error('视频加载失败：' + error.message)
    },
    
    // 禁用未来日期
    disabledDate(time) {
      return time.getTime() > Date.now()
    }
  }
}
</script>

<style scoped>
.video-manager {
  padding: 20px;
}

.query-form {
  margin-bottom: 20px;
}

.operation-container {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.delete-btn {
  color: #F56C6C;
}

:deep(.el-dialog__body) {
  padding: 10px;
  background: #000;
}

/* 确保时间选择器面板可见 */
.el-picker-panel {
  z-index: 2000 !important;
}

.text-center {
  text-align: center;
  line-height: 32px;
}

.separator {
  color: #909399;
}
</style> 