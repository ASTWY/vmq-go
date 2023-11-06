<template>
    <DashboardLayout>
        <t-card>
            <t-tabs :default-value="1">
                <t-tab-panel :value="1" label="基本设置">
                    <t-card :bordered="false" title="通讯密钥" header-bordered>
                        <t-textarea v-model="settings.apiSecret" placeholder="请输入API密钥" />
                    </t-card>
                    <t-card :bordered="false" title="任意金额收款码" header-bordered>
                        <t-row>
                            <t-space align="center" class="space-border">
                                <t-card title="微信收款二维码">
                                    <t-space align="center" direction="vertical">
                                        <div>
                                            <img :src="wxQrcodeUrl" alt="微信二维码" style="width: 200px; height: 200px;">
                                        </div>
                                        <t-upload :request-method="requestMethod"></t-upload>
                                    </t-space>
                                </t-card>
                                <t-card title="支付宝收款二维码">
                                    <t-space align="center" direction="vertical">
                                        <div>
                                            <img :src="alipayQrcodeUrl" alt="支付宝二维码" style="width: 200px; height: 200px;" />
                                        </div>
                                        <t-upload :request-method="requestMethod"></t-upload>
                                    </t-space>
                                </t-card>
                            </t-space>
                        </t-row>
                    </t-card>
                    <t-button v-if="isChanged" theme="primary" size="large" @click="onUpdate" block
                        :disabled="!isChanged">保存</t-button>
                </t-tab-panel>
                <t-tab-panel :value="2" label="订单设置">
                    <t-card :bordered="false" title="订单有效时间" header-bordered>
                        <t-input-number v-model="settings.expire" theme="normal" autoWidth size="large">
                            <template #suffix><span>分钟</span></template>
                        </t-input-number>
                    </t-card>
                    <t-card :bordered="false" title="同金额最大订单数" header-bordered>
                        <t-input-number v-model="settings.orderMaxNum" theme="normal" autoWidth size="large">
                        </t-input-number>
                    </t-card>
                    <t-card :bordered="false" title="异步通讯地址" header-bordered>
                        <t-textarea v-model="settings.notifyUrl" placeholder="请输入异步通讯地址" />
                    </t-card>
                    <t-card :bordered="false" title="同步返回地址" header-bordered>
                        <t-textarea v-model="settings.returnUrl" placeholder="请输入异步通讯地址" />
                    </t-card>
                    <t-card :bordered="false" title="订单区分方式" header-bordered>
                        <t-radio-group variant="default-filled" :default-value="settings.orderType"
                            v-model:value="settings.orderType">
                            <t-radio-button value="1">金额递减</t-radio-button>
                            <t-radio-button value="2">金额递增</t-radio-button>
                        </t-radio-group>
                    </t-card>
                    <t-button v-if="isChanged" theme="primary" size="large" @click="onUpdate" :disabled="!isChanged"
                        block>保存</t-button>
                </t-tab-panel>
                <t-tab-panel :value="3" label="邮箱通知">
                    <t-card :bordered="false" title="STMP配置" header-bordered>
                        <template #actions>
                            <t-button theme="default" @click="sendTestEmail" block>发送测试邮件</t-button>
                        </template>
                        <t-space direction="vertical">
                            <!-- 收件人邮箱 -->
                            <t-input label="收件人:" v-model="settings.emailSMTPto" placeholder="请输入收件人邮箱" />
                            <!-- 发件人邮箱 -->
                            <t-input label="发件人:" v-model="settings.emailSMTPfrom" placeholder="请输入发件人邮箱" />
                            <!-- SMTP服务器 -->
                            <t-input label="SMTP服务器:" v-model="settings.emailSMTPhost" placeholder="请输入SMTP服务器" />
                            <!-- SMTP端口 -->
                            <t-input label="SMTP端口:" v-model="settings.emailSMTPport" placeholder="请输入SMTP端口" />
                            <!-- SMTP账号 -->
                            <t-input label="SMTP账号:" v-model="settings.emailSMTPuser" placeholder="请输入SMTP账号" />
                            <!-- SMTP密码 -->
                            <t-input label="SMTP密码:" v-model="settings.emailSMTPpwd" placeholder="请输入SMTP密码" />
                        </t-space>
                    </t-card>
                    <t-card :bordered="false" title="通知开关" header-bordered>
                        <!-- 通知开关  监控端掉线、订单异常、收款通知 -->
                        <t-checkbox-group v-model="checked" :options="['监控端掉线', '订单回调异常', '收款通知']"></t-checkbox-group>
                    </t-card>
                    <t-button v-if="isChanged" theme="primary" size="large" @click="onUpdate" :disabled="!isChanged"
                        block>保存</t-button>
                </t-tab-panel>
                <t-tab-panel :value="4" label="管理员账号">
                    <!-- 管理员账号密码修改 -->
                    <t-card :bordered="false" title="管理员账号密码修改" header-bordered>
                        <t-space direction="vertical" class="space-border">
                            <t-input size="large" v-model="settings.adminUser" clearable placeholder="请输入账户名">
                                <template #prefix-icon>
                                    <t-icon name="desktop" />
                                </template>
                            </t-input>
                            <t-input size="large" v-model="settings.adminPwd" type="password" placeholder="请输入密码">
                                <template #prefix-icon>
                                    <t-icon name="lock-on" />
                                </template>
                            </t-input>
                            <t-button v-if="isChanged" theme="primary" size="large" @click="onUpdate"
                                :disabled="!isChanged">保存</t-button>
                        </t-space>
                    </t-card>
                </t-tab-panel>
            </t-tabs>
        </t-card>
    </DashboardLayout>
</template>


<script>
import MD5 from 'crypto-js/md5'
import DashboardLayout from '@/components/Dashboard.vue'
import api from '@/api'



export default {
    name: 'SettingsView',
    // 组件逻辑和功能
    components: {
        // 引入组件
        DashboardLayout
    },
    data() {
        return {
            // 组件数据
            settings: {}, // 系统设置
            settingsCopy: {}, // 系统设置副本 用于比较是否有修改
            checked: [], // 通知开关
        }
    },
    watch: {
        // 监听数据变化 当数据变化时执行
        checked: {
            handler: function (val, oldVal) {
                if (val.includes('监控端掉线')) {
                    this.settings.monitorNotice = '1'
                } else {
                    this.settings.monitorNotice = '0'
                }
                if (val.includes('订单回调异常')) {
                    this.settings.errorNotice = '1'
                } else {
                    this.settings.errorNotice = '0'
                }
                if (val.includes('收款通知')) {
                    this.settings.payNotice = '1'
                } else {
                    this.settings.payNotice = '0'
                }
            },
            deep: true // 深度监听
        }
    },
    computed: {
        // 计算属性
        isChanged() {
            return JSON.stringify(this.settings) !== JSON.stringify(this.settingsCopy)
        },
        wxQrcodeUrl() {
            let url = api.getQrcodeImg(this.settings.wechatPay)
            return url
        },
        alipayQrcodeUrl() {
            let url = api.getQrcodeImg(this.settings.aliPay)
            return url
        },
    },
    methods: {
        // 组件方法
        onUpdate() {
            // 更新系统设置
            // 取出settingsCopy中的值，与settings比较，如果不同则更新
            let data = {}
            for (let key in this.settingsCopy) {
                if (this.settingsCopy[key] !== this.settings[key]) {
                    if (key === 'adminPwd') {
                        this.settings[key] = MD5(this.settings[key]).toString()
                        continue
                    }
                    data[key] = this.settings[key]
                }
            }
            for (let key in data) {
                api.updateSetting({
                    key: key,
                    value: data[key]
                }).then(data => {
                    this.getSettings()
                })
            }
        },
        getSettings() {
            // 获取系统设置
            api.getSettings().then(res => {
                let data = res.data.data
                this.settings = data
                if (this.settings['monitorNotice'] === '1') {
                    this.checked = [...this.checked, '监控端掉线']
                }
                if (this.settings['errorNotice'] === '1') {
                    this.checked = [...this.checked, '订单回调异常']
                }
                if (this.settings['payNotice'] === '1') {
                    this.checked = [...this.checked, '收款通知']
                }
                this.settingsCopy = JSON.parse(JSON.stringify(data))
            })
        },
        requestMethod(file) {
            // 自定义上传方法
            return new Promise((resolve) => {
                let fromData = new FormData()
                fromData.append('file', file.raw)
                api.postQrcode(fromData).then(res => {
                    const content = res.data.data.content
                    // 判断content开头是否为http
                    if (content.startsWith('https://qr.alipay.com/')) {
                        this.settings.aliPay = content
                    } else if (content.startsWith('wxp://')) {
                        this.settings.wechatPay = content
                    }
                    resolve({ status: 'success' })
                })
            })
        },
        sendTestEmail() {
            // 发送测试邮件
            api.sendTestEmail().then(res => {
                this.$message.success('发送成功')
            })
        }

    },
    mounted() {
        // 组件挂载后执行
        this.getSettings()
    },
};
</script>
