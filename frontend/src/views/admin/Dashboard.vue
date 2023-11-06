<template>
    <DashboardLayout>
        <t-space direction="vertical">
            <t-row :gutter="16">
                <!-- 今日收入 -->
                <t-col :span="3">
                    <t-card title="今日收入">
                        <h1>￥{{ today.income }}</h1>
                    </t-card>
                </t-col>
                <!-- 今日订单 -->
                <t-col :span="3">
                    <t-card title="今日订单">
                        <h1>{{ today.orders }}</h1>
                    </t-card>
                </t-col>
                <!-- 今日成功订单 -->
                <t-col :span="3">
                    <t-card title="今日成功订单">
                        <h1>{{ today.successOrders }}</h1>
                    </t-card>
                </t-col>
                <!-- 今日失败订单 -->
                <t-col :span="3">
                    <t-card title="今日失败订单">
                        <h1>{{ today.failedOrders }}</h1>
                    </t-card>
                </t-col>
            </t-row>
            <t-row>
                <t-col :span="12">
                    <t-card>
                        <t-list split size="small">
                            <t-list-item>
                                <t-list-item-meta title="监控端状态" :description="monitorStatus" />
                            </t-list-item>
                            <t-list-item>
                                <t-list-item-meta title="最后心跳时间" :description="lastHeartbeatTime" />
                            </t-list-item>
                            <t-list-item>
                                <t-list-item-meta title="最后支付时间" :description="lastPaymentTime" />
                            </t-list-item>
                            <t-list-item>
                                <t-list-item-meta title="配置数据" :description="monitorConfig" />
                                <t-image :src="configQrcode" :style="{ width: '120px', height: '120px' }" />
                            </t-list-item>
                        </t-list>
                    </t-card>
                </t-col>
            </t-row>
        </t-space>
    </DashboardLayout>
</template>

<script>
import DashboardLayout from '@/components/Dashboard.vue'
import api from '@/api'
export default {
    name: 'DashboardView',
    // 组件逻辑和功能
    components: {
        // 引入组件
        DashboardLayout
    },
    data() {
        return {
            interval: null,
            // 组件数据
            today: {
                income: 0,
                orders: 0,
                successOrders: 0,
                failedOrders: 0
            },
            lastHeartbeat: 0,
            lastPayment: 0,
            apiSecret: '',
            host: window.location.host.replace('http://', '').replace('https://', ''),
            todayIncome: 0,
            todayOrders: 0,
            todaySuccessOrders: 0,
            todayFailedOrders: 0,
            sevenDaysIncome: 0
        }
    },
    computed: {
        // 计算属性
        monitorConfig() {
            return this.host + '/' + this.apiSecret
        },
        configQrcode() {
            return api.getQrcodeImg(this.monitorConfig)
        },
        monitorStatus() {
            // 如果最后心跳时间（时间戳）与当前时间相差 30 秒以上，说明监控程序已经断开连接
            if (this.lastHeartbeat + 30 < Math.floor(Date.now() / 1000)) {
                return '离线'
            }
            return '在线'
        },
        lastHeartbeatTime() {
            // 格式化最后心跳时间 2019-01-01 00:00:00
            return this.lastHeartbeat ? this.formatDate(this.lastHeartbeat) : ''
        },
        lastPaymentTime() {
            // 格式化最后支付时间 2019-01-01 00:00:00
            return this.lastPayment ? this.formatDate(this.lastPayment) : ''
        }
    },
    methods: {
        // 组件方法
        formatDate(timestamp) {
            // 首先根据时间戳长度判断是否为秒级时间戳
            if (timestamp.toString().length === 10) {
                timestamp = parseInt(timestamp)
            }
            if (timestamp.toString().length === 13) {
                timestamp = parseInt(timestamp / 1000)
            }
            // 时间戳转换为时间
            const date = new Date(timestamp * 1000)
            const year = date.getFullYear()
            const month = date.getMonth() + 1
            const day = date.getDate()
            const hour = date.getHours()
            const minute = date.getMinutes()
            const second = date.getSeconds()
            return `${year}-${month}-${day} ${hour}:${minute}:${second}`
        },
        getTodayData() {
            api.getTodayData().then((res) => {
                const data = res.data.data
                this.today.income = data.income
                this.today.orders = data.orders
                this.today.successOrders = data.successOrders
                this.today.failedOrders = data.failedOrders
            })
        }
    },
    mounted() {
        this.getTodayData()
        // 组件挂载后执行
        api.getSettings().then((res) => {
            const data = res.data.data
            this.apiSecret = data.apiSecret
            this.lastHeartbeat = data.lastHeart
            this.lastPayment = data.lastPay
        })
        // 启动定时器，每3秒更新一次心跳时间 如果路由切换，定时器会被销毁
        this.interval = setInterval(() => {
            api.getSettings().then((res) => {
                const data = res.data.data
                this.lastHeartbeat = data.lastHeart
                this.lastPayment = data.lastPay
            }).catch((err) => {
                console.log(err)
            })
        }, 3000)

    },
    beforeDestroy() {
        // 组件销毁前执行
        clearInterval(this.interval)
    }
};
</script>

<style scoped>
/* 组件样式 */
.t-space {
    width: calc(100%);
    margin: 20px auto;
}
</style>