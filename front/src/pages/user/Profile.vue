<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 左侧个人信息 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>个人信息</span>
          </template>
          <el-avatar :src="userStore.userInfo.avatar || ''" size="large">
            {{ userStore.userInfo.username?.slice(0, 1) }}
          </el-avatar>
          <div class="info-item">
            <span class="label">用户名：</span>
            <span>{{ userStore.userInfo.username }}</span>
          </div>
          <div class="info-item">
            <span class="label">角色：</span>
            <span>{{ userStore.userInfo.role === 'admin' ? '管理员' : '普通用户' }}</span>
          </div>
          <div class="info-item">
            <span class="label">注册时间：</span>
            <span>{{ userStore.userInfo.createdAt }}</span>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧修改密码 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>修改密码</span>
          </template>

          <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-width="100px">
            <el-form-item label="原密码" prop="old_password">
              <el-input v-model="pwdForm.old_password" type="password" placeholder="请输入原密码"></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input v-model="pwdForm.new_password" type="password" placeholder="请输入新密码（6-20位）"></el-input>
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirm_password">
              <el-input v-model="pwdForm.confirm_password" type="password" placeholder="请再次输入新密码"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleChangePwd">提交修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { changePassword } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()
const pwdFormRef = ref(null)

// 密码表单
const pwdForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 表单校验规则
const pwdRules = ref({
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '新密码长度在6-20位之间', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdForm.value.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

// 修改密码
const handleChangePwd = () => {
  pwdFormRef.value.validate((valid) => {
    if (!valid) return
    changePassword(pwdForm.value).then(() => {
      ElMessage.success('密码修改成功，请重新登录')
      userStore.clearUser()
      router.push('/login')
    })
  })
}
</script>

<style scoped>
.profile-container {
  padding: 20px 0;
}

.info-item {
  margin: 10px 0;
  font-size: 14px;
}

.label {
  color: #999;
  margin-right: 10px;
}
</style>