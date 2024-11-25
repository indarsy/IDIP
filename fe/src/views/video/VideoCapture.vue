<template>
  <div class="video-capture">
    <el-card class="config-card">
      <template #header>
        <div class="card-header">
          <span class="title">
            <el-icon><VideoCamera /></el-icon>
            视频采集配置
          </span>
        </div>
      </template>
      
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef" 
        label-width="100px"
        class="config-form"
      >
        <!-- 第一行：选择车间 -->
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="选择车间" prop="workshopId">
              <el-select 
                v-model="form.workshopId" 
                placeholder="请选择车间"
                @change="handleWorkshopChange"
                style="width: 100%"
              >
                <el-option
                  v-for="item in workshopList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 第二行：视频预览 -->
        <el-form-item label="实时预览">
    <div class="preview-wrapper">
      <div 
        class="preview-container" 
        v-if="currentWorkshop?.rtspUrl"
      >
        <VideoPlayer 
          ref="playerRef"
          :form="{ rtsp_url: currentWorkshop.rtspUrl }"
        />
      </div>
      <div v-else class="preview-placeholder">
        <el-empty description="请选择车间以预览视频" />
      </div>
    </div>
  </el-form-item>

        <!-- 第三行：采集时间、间隔和按钮 -->
        <el-row :gutter="20" class="bottom-row">
          <!-- 采集时间设置 -->
          <el-col :span="9">
            <el-form-item label="采集时间" required>
              <el-row :gutter="10">
                <el-col :span="11">
                  <el-form-item prop="startTime" class="no-margin">
                    <el-date-picker
                      v-model="form.startTime"
                      type="datetime"
                      placeholder="开始时间"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="2" class="text-center">
                  <span class="separator">至</span>
                </el-col>
                <el-col :span="11">
                  <el-form-item prop="endTime" class="no-margin">
                    <el-date-picker
                      v-model="form.endTime"
                      type="datetime"
                      placeholder="结束时间"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>

          <!-- 采集间隔设置 -->
          <el-col :span="11">
            <el-form-item label="采集间隔" prop="interval">
              <el-row :gutter="10">
                <el-col :span="form.interval === 'custom' ? 16 : 24">
                  <el-select 
                    v-model="form.interval" 
                    placeholder="请选择间隔"
                    style="width: 100%"
                  >
                    <el-option label="30分钟" :value="30" />
                    <el-option label="1小时" :value="60" />
                    <el-option label="2小时" :value="120" />
                    <el-option label="自定义" value="custom" />
                  </el-select>
                </el-col>
                <el-col :span="8" v-if="form.interval === 'custom'">
                  <el-input-number
                    v-model="customInterval"
                    :min="1"
                    :max="1440"
                    placeholder="分钟"
                    style="width: 100%"
                    controls-position="right"
                  />
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>

          <!-- 提交按钮 -->
          <el-col :span="4">
            <el-form-item label=" " label-width="0">
              <el-button 
                type="primary" 
                @click="handleSubmit" 
                :loading="submitting"
                class="submit-btn"
              >
                <el-icon><VideoPlay /></el-icon>
                开始采集
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 任务列表卡片 -->
    <el-card class="task-card">
      <template #header>
        <div class="card-header">
          <span class="title">
            <el-icon><List /></el-icon>
            采集任务列表
          </span>
          <el-button 
            type="primary" 
            link 
            @click="fetchCaptureList"
            :loading="loading"
          >
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table 
        :data="captureList" 
        v-loading="loading"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="workshop.name" label="车间" width="580" />
        <el-table-column prop="startTime" label="开始时间" width="280">
          <template #default="{ row }">
            {{ formatDateTime(row.startTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="endTime" label="结束时间" width="280">
          <template #default="{ row }">
            {{ formatDateTime(row.endTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="interval" label="采集间隔" width="120" align="center">
          <template #default="{ row }">
            {{ row.interval }}分钟
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="140" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="dark">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              link 
              @click="handleCancel(row)"
              v-if="row.status === 'waiting' || row.status === 'running'"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import VideoPlayer from '@/components/VideoPlayer.vue'
import { formatDateTime } from '@/utils/format'
import { 
  getWorkshopList, 
  // downloadVideo, 
  // startRecording, 
  // stopRecording 
} from '@/api/video'

export default {
  name: 'VideoCapture',
  components: {
    VideoPlayer
  },
  setup() {
    const store = useStore()
    const formRef = ref(null)
    const loading = ref(false)
    const submitting = ref(false)
    const customInterval = ref(30)
    const workshopList = ref([])
    const captureList = ref([])
    const playerRef = ref(null)

    const form = ref({
      workshopId: '',
      startTime: '',
      endTime: '',
      interval: 30
    })

    const rules = {
      workshopId: [
        { required: true, message: '请选择车间', trigger: 'change' }
      ],
      startTime: [
        { required: true, message: '请选择开始时间', trigger: 'change' },
        { 
          validator: (rule, value, callback) => {
            if (value && value < new Date()) {
              callback(new Error('开始时间不能早于当前时间'))
            } else {
              callback()
            }
          },
          trigger: 'change'
        }
      ],
      endTime: [
        { required: true, message: '请选择结束时间', trigger: 'change' },
        {
          validator: (rule, value, callback) => {
            if (value && form.value.startTime && value <= form.value.startTime) {
              callback(new Error('结束时间必须晚于开始时间'))
            } else {
              callback()
            }
          },
          trigger: 'change'
        }
      ],
      interval: [
        { required: true, message: '请选择采集间隔', trigger: 'change' }
      ]
    }

    // 获取车间列表
    const fetchWorkshops = async () => {
      try {
        const { data } = await getWorkshopList()
        workshopList.value = data || []
      } catch (error) {
        console.error('获取车间列表失败:', error)
        ElMessage.error('获取车间列表失败')
        workshopList.value = []
      }
    }

    // 当前选中的车间
    const currentWorkshop = computed(() => {
      if (!workshopList.value?.length || !form.value.workshopId) return null
      return workshopList.value.find(w => w.id === form.value.workshopId)
    })

    // 获取采集任务列表
    const fetchCaptureList = async () => {
      loading.value = true
      try {
        const workshopId = form.value.workshopId || undefined
        const response = await store.dispatch('capture/fetchList', workshopId)
        
        // 确保返回的数据是数组
        captureList.value = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        console.error('获取采集任务列表失败:', error)
        ElMessage.error('获取采集任务列表失败')
        captureList.value = []
      } finally {
        loading.value = false
      }
    }

    // 处理车间选择变化
    const handleWorkshopChange = async () => {
      // 先清理旧的视频流
      if (playerRef.value) {
        await playerRef.value.cleanupConnection()
      }
      
      // 等待一小段时间确保资源被释放
      await new Promise(resolve => setTimeout(resolve, 100))
      
      if (form.value.workshopId) {
        fetchCaptureList()
      }
    }

    // 处理表单提交
    const handleSubmit = async () => {
      if (!formRef.value) return
      
      try {
        await formRef.value.validate()
        
        submitting.value = true
        const captureData = {
          ...form.value,
          interval: form.value.interval === 'custom' ? customInterval.value : form.value.interval
        }
        
        await store.dispatch('capture/create', captureData)
        ElMessage.success('创建采集任务成功')
        fetchCaptureList()
        
        // 重置表单
        formRef.value.resetFields()
      } catch (error) {
        console.error('创建采集任务失败:', error)
        ElMessage.error(error.message || '创建采集任务失败')
      } finally {
        submitting.value = false
      }
    }

    // 取消任务
    const handleCancel = async (row) => {
      try {
        await ElMessageBox.confirm('确认取消该采集任务?', '提示', {
          type: 'warning'
        })
        
        await store.dispatch('capture/cancel', row.id)
        ElMessage.success('取消成功')
        fetchCaptureList()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('取消任务失败:', error)
          ElMessage.error('取消任务失败')
        }
      }
    }

    // 状态显示
    const getStatusType = (status) => {
      const types = {
        waiting: 'info',
        running: 'primary',
        completed: 'success',
        failed: 'danger',
        cancelled: 'warning'
      }
      return types[status] || 'info'
    }

    const getStatusText = (status) => {
      const texts = {
        waiting: '等待中',
        running: '进行中',
        completed: '已完成',
        failed: '失败',
        cancelled: '已取消'
      }
      return texts[status] || status
    }

    onMounted(async () => {
      try {
        loading.value = true
        // 并行获取车间列表和采集任务列表
        await Promise.all([
          fetchWorkshops(),
          fetchCaptureList() // 不需要等待选择车间就获取所有任务
        ])
      } catch (error) {
        console.error('初始化数据失败:', error)
        ElMessage.error('加载数据失败')
      } finally {
        loading.value = false
      }
    })

    return {
      formRef,
      form,
      rules,
      loading,
      submitting,
      customInterval,
      workshopList,
      currentWorkshop,
      captureList,
      handleWorkshopChange,
      handleSubmit,
      handleCancel,
      fetchCaptureList,
      formatDateTime,
      getStatusType,
      getStatusText,
      playerRef
    }
  }
}
</script>

<style scoped>
.video-capture {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 84px);
}

.config-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
}

.config-form {
  padding: 20px 0;
}

.preview-wrapper {
  width: 100%;
  margin: 0;
  padding: 0;
  background-color: transparent;
}

.preview-container {
  width: 100%;
  aspect-ratio: 16/9;
  min-height: 400px;
  max-height: 600px;
  background-color: #1a1a1a;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.preview-placeholder {
  width: 100%;
  aspect-ratio: 16/9;
  min-height: 400px;
  max-height: 600px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px dashed #dcdfe6;
}

.text-center {
  text-align: center;
  line-height: 32px;
}

.separator {
  color: #909399;
}

.no-margin {
  margin-bottom: 0;
}

.submit-btn {
  width: 120px;
}

:deep(.el-form-item__content) {
  margin: 0 !important;
  padding: 0 !important;
}

:deep(.el-input-number .el-input__wrapper) {
  padding-right: 30px;
}

:deep(.el-table) {
  border-radius: 4px;
  overflow: hidden;
}

/* 调整表单项间距 */
:deep(.el-form-item) {
  margin-bottom: 22px;
  overflow: visible;
}

/* 优化预览区域样式 */
:deep(.el-form-item__label) {
  font-weight: 500;
}

/* 调整日期选择器样式 */
:deep(.el-date-editor.el-input) {
  width: 100%;
}

/* 确保视频播放器填充整个容器 */
:deep(.video-js) {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
}

/* 移除所有可能的内边距 */
:deep(.el-form-item__content) {
  margin: 0 !important;
  padding: 0 !important;
}

:deep(.rtsp-player),
:deep(.video-container),
:deep(video) {
  margin: 0;
  padding: 0;
}

.bottom-row {
  margin-top: 20px;
  align-items: flex-start;
}

.submit-btn {
  width: 100%;  /* 让按钮填充整个容器 */
  margin-top: 4px;  /* 微调按钮位置以对齐其他元素 */
}

/* 确保预览容器的高度合适 */
.preview-container,
.preview-placeholder {
  width: 100%;
  aspect-ratio: 16/9;
  min-height: 360px;
  max-height: 500px;
}

/* 调整表单项间距 */
:deep(.el-form-item) {
  margin-bottom: 18px;
}

/* 最后一行的表单项去掉底部间距 */
.bottom-row :deep(.el-form-item) {
  margin-bottom: 0;
}

/* 调整选择器的最小宽度 */
:deep(.el-select) {
  min-width: 160px;
}

/* 确保选择框内容完整显示 */
:deep(.el-select .el-input__wrapper) {
  padding: 0 30px 0 12px;
}

/* 调整输入框内容的显示 */
:deep(.el-select .el-input__inner) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
