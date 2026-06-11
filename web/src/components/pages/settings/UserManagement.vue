<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { toast } from 'vue-sonner'
import { request } from '@/api/api'
import { CONSTANT } from '@/constant'

interface User {
  id: number
  username: string
  role: string
  created_on: string
}

const users = ref<User[]>([])
const loading = ref(false)
const dialogOpen = ref(false)
const newUser = ref({ username: '', password: '', role: 'user' })

// 从 JWT 解析当前用户名
const getCurrentUsername = () => {
  try {
    const token = localStorage.getItem(CONSTANT.STORE_TOKEN_NAME) || ''
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.username || ''
  } catch { return '' }
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const rsp = await request.get('/users/list')
    if (rsp.data.code === 200) {
      users.value = rsp.data.data || []
    }
  } finally {
    loading.value = false
  }
}

const handleAdd = async () => {
  if (!newUser.value.username || !newUser.value.password) {
    toast.error('请填写用户名和密码')
    return
  }
  try {
    const rsp = await request.post('/users/add', newUser.value)
    if (rsp.data.code === 200) {
      toast.success('添加用户成功')
      dialogOpen.value = false
      newUser.value = { username: '', password: '', role: 'user' }
      fetchUsers()
    } else {
      toast.error(rsp.data.msg)
    }
  } catch {
    toast.error('添加失败')
  }
}

const handleDelete = async (user: User) => {
  if (!confirm(`确定删除用户 "${user.username}"？`)) return
  try {
    const rsp = await request.post('/users/delete', { id: user.id })
    if (rsp.data.code === 200) {
      toast.success('删除成功')
      fetchUsers()
    } else {
      toast.error(rsp.data.msg)
    }
  } catch {
    toast.error('删除失败')
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">用户管理</h2>
      <Dialog v-model:open="dialogOpen">
        <DialogTrigger as-child>
          <Button>新增用户</Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>新增用户</DialogTitle>
          </DialogHeader>
          <div class="space-y-4 mt-4">
            <div>
              <label class="text-sm font-medium">用户名</label>
              <Input v-model="newUser.username" placeholder="请输入用户名" class="mt-1" />
            </div>
            <div>
              <label class="text-sm font-medium">密码</label>
              <Input v-model="newUser.password" type="password" placeholder="请输入密码" class="mt-1" />
            </div>
            <div>
              <label class="text-sm font-medium">角色</label>
              <select v-model="newUser.role" class="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm mt-1">
                <option value="admin">管理员 (admin)</option>
                <option value="user">普通用户 (user)</option>
              </select>
            </div>
            <Button @click="handleAdd" class="w-full">确认添加</Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>用户列表</CardTitle>
        <CardDescription>当前系统共有 {{ users.length }} 个用户</CardDescription>
      </CardHeader>
      <CardContent>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b text-left text-muted-foreground">
              <th class="py-3 px-4">ID</th>
              <th class="py-3 px-4">用户名</th>
              <th class="py-3 px-4">角色</th>
              <th class="py-3 px-4">创建时间</th>
              <th class="py-3 px-4 text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading" class="border-b">
              <td colspan="5" class="py-8 text-center text-muted-foreground">加载中...</td>
            </tr>
            <tr v-else-if="users.length === 0" class="border-b">
              <td colspan="5" class="py-8 text-center text-muted-foreground">暂无用户</td>
            </tr>
            <tr v-for="user in users" :key="user.id" class="border-b hover:bg-muted/50">
              <td class="py-3 px-4">{{ user.id }}</td>
              <td class="py-3 px-4">{{ user.username }}</td>
              <td class="py-3 px-4">
                <Badge :variant="user.role === 'admin' ? 'default' : 'outline'">
                  {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                </Badge>
              </td>
              <td class="py-3 px-4">{{ user.created_on }}</td>
              <td class="py-3 px-4 text-right">
                <Button
                  size="sm"
                  variant="outline"
                  class="text-red-500 border-red-300 hover:bg-red-50"
                  @click="handleDelete(user)"
                  :disabled="user.username === getCurrentUsername()"
                >删除</Button>
              </td>
            </tr>
          </tbody>
        </table>
      </CardContent>
    </Card>
  </div>
</template>
