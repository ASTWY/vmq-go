import axios from 'axios';
import { useCounterStore } from '@/stores/counter'
import router from '@/router'


const apiURL = '/api'


axios.defaults.baseURL = apiURL
// axios 请求拦截器 
// 每次请求前，如果存在token则在请求头中携带token
axios.interceptors.request.use(
    config => {
        const store = useCounterStore()
        if (store.token) {
            config.headers.Authorization = "Bearer " + store.token
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)
// 当返回401时，清除token并跳转到登录页面 否则弹出错误信息
axios.interceptors.response.use(
    response => {
        return response
    },
    error => {
        const store = useCounterStore()
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    store.logout()
                    router.push({ name: 'login' })
            }
        }
        return Promise.reject(error.response.data)
    }
)

const api = {
    // 获取captcha
    getCaptcha: () => {
        return axios.get('/captcha')
    },
    // 获取qrcode json
    getQrcode: (content) => {
        return axios.get('/qrcode?content=' + content)
    },
    // 获取qrcode img
    getQrcodeImg: (content) => {
        return apiURL + '/qrcode?content=' + content + '&format=image'
    },
    // 解析qrcode
    postQrcode: (params) => {
        return axios.post('/qrcode', params)
    },
    // 获取订单
    getOrder: (orderID) => {
        return axios.get('/order/' + orderID)
    },
    // 获取订单状态
    getOrderStatus: (orderID) => {
        return axios.get('/order/' + orderID + '/state')
    },
    // login
    login: (params) => {
        return axios.post('/admin/login', params)
    },
    // logout
    logout: () => {
        return axios.get('/admin/logout')
    },
    // 获取今日数据 
    getTodayData: () => {
        return axios.get('/admin/order/today')
    },
    // 获取所有设置
    getSettings: () => {
        return axios.get('/admin/settings')
    },
    // 更新设置
    updateSetting: (params) => {
        return axios.put('/admin/setting', params)
    },
    // 发送测试邮件
    sendTestEmail: () => {
        return axios.post('/admin/setting/email')
    },
    // 获取订单列表
    getOrders: (page, pageSize) => {
        return axios.get('/admin/orders' + '?page=' + page + '&pageSize=' + pageSize)
    },
    // 补单
    replenish: (orderID) => {
        return axios.post('/admin/order/' + orderID + '/replenish')
    },
    deleteOrder: (orderID) => {
        return axios.delete('/admin/order/' + orderID)
    },
    // 获取所有收款码
    getQrcodes: () => {
        return axios.get('/admin/qrcodes')
    },
    // 添加收款码
    addQrcode: (params) => {
        return axios.post('/admin/qrcode', params)
    },
    // 删除收款码
    deleteQrcode: (qrcodeID) => {
        return axios.delete('/admin/qrcode/' + qrcodeID)
    },
    // 获取支付记录
    getPaylogs: (page, pageSize) => {
        return axios.get('/admin/paylogs' + '?page=' + page + '&pageSize=' + pageSize)
    },
}

export default api;