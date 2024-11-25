<template>
  <div class="workshop-container">
    <!-- 操作区域 -->
    <el-card class="operation-card">
      <template #header>
        <div class="card-header">
          <span>车间管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>添加车间
          </el-button>
        </div>
      </template>

      <!-- 车间列表 -->
      <el-table
        v-loading="loading"
        :data="workshopList"
        border
        stripe
      >
        <el-table-column type="index" label="序号" width="80" />
        
        <el-table-column prop="name" label="车间名称" min-width="150" />
        
        <!-- <el-table-column prop="code" label="车间编号" width="120" /> -->
        
        <el-table-column prop="rtspUrl" label="RTSP地址" min-width="250" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column prop="updatedAt" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.updatedAt) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(row)"
              >
                编辑
              </el-button>
              <!-- <el-button 
                :type="row.status === 1 ? 'danger' : 'success'" 
                size="small" 
                @click="handleToggleStatus(row)"
              >
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button> -->
              <el-button 
                type="danger" 
                size="small" 
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加车间' : '编辑车间'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="车间名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入车间名称" />
        </el-form-item>
        
        <!-- <el-form-item label="车间编号" prop="code">
          <el-input v-model="formData.code" placeholder="请输入车间编号" />
        </el-form-item> -->
        
        <el-form-item label="RTSP地址" prop="rtspUrl">
          <el-input v-model="formData.rtspUrl" placeholder="请输入RTSP地址" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="formData.status"
            :active-value="1"
            :inactive-value="2"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import { 
  getWorkshopList,
  addWorkshop,
  updateWorkshop,
  deleteWorkshop,
  //toggleWorkshopStatus
} from '@/api/workshop'

// 状态定义
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const workshopList = ref([])
const formRef = ref(null)
const formData = ref({
  name: '',
  code: '',
  rtspUrl: '',
  status: 1
})

// 表单校验规则
const formRules = {
  name: [
    { required: true, message: '请输入车间名称', trigger: 'blur' },
    { min: 1, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  // code: [
  //   { required: true, message: '请输入车间编号', trigger: 'blur' },
  //   { pattern: /^[A-Za-z0-9-_]+$/, message: '只能包含字母、数字、下划线和横线', trigger: 'blur' }
  // ],
  rtspUrl: [
    { required: true, message: '请输入RTSP地址', trigger: 'blur' },
    { pattern: /^rtsp:\/\/.+/, message: 'RTSP地址格式不正确', trigger: 'blur' }
  ]
}

// 获取车间列表
const fetchWorkshopList = async () => {
  loading.value = true
  try {
    const { data } = await getWorkshopList()
    workshopList.value = data
  } catch (error) {
    ElMessage.error('获取车间列表失败')
  } finally {
    loading.value = false
  }
}

// 添加车间
const handleAdd = () => {
  dialogType.value = 'add'
  formData.value = {
    name: '',
    //code: '',
    rtspUrl: '',
    status: 1
  }
  dialogVisible.value = true
}

// 编辑车间
const handleEdit = (row) => {
  dialogType.value = 'edit'
  formData.value = { ...row }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    // 等待表单验证完成
    await formRef.value.validate()
    
    submitting.value = true
    if (dialogType.value === 'add') {
      await addWorkshop(formData.value)
      ElMessage.success('添加成功')
    } else {
      const { id, ...updateData } = formData.value
      await updateWorkshop(id, updateData)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchWorkshopList()
  } catch (error) {
    // 区分表单验证错误和API请求错误
    if (error.message) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error(dialogType.value === 'add' ? '添加失败' : '更新失败')
    }
  } finally {
    submitting.value = false
  }
}

// 切换状态
// const handleToggleStatus = async (row) => {
//   try {
//     await ElMessageBox.confirm(
//       `确认${row.status === 1 ? '禁用' : '启用'}该车间?`,
//       '提示',
//       { type: 'warning' }
//     )
//     await toggleWorkshopStatus(row.id)
//     ElMessage.success(`${row.status === 1 ? '禁用' : '启用'}成功`)
//     fetchWorkshopList()
//   } catch (error) {
//     if (error !== 'cancel') {
//       ElMessage.error(`${row.status === 1 ? '禁用' : '启用'}失败`)
//     }
//   }
// }

// 删除车间
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该车间?', '提示', {
      type: 'warning'
    })
    await deleteWorkshop(row.id)
    ElMessage.success('删除成功')
    fetchWorkshopList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 生命周期钩子
onMounted(() => {
  fetchWorkshopList()
})
</script>

<style scoped lang="scss">
.workshop-container {
  padding: 20px;
  
  .operation-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}

:deep(.el-dialog__body) {
  padding-top: 20px;
}
</style>
