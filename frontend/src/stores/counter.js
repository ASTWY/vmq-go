import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', {
  state: () => {
    const storedToken = localStorage.getItem('counterToken')
    return {
      token: storedToken || '',
    }
  },
  getters: {
    isLogin: (state) => !!state.token
  },
  actions: {
    // 登录
    login(token) {
      this.token = token
      localStorage.setItem('counterToken', token)
    },
    // 登出
    logout() {
      this.token = ''
      localStorage.removeItem('counterToken')
    },
  },
  beforeUnmount() {
    // 在组件销毁之前存储状态
    localStorage.setItem('counterToken', this.token)
  }
})
