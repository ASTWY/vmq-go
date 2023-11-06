<template>
    <div class="login">
        <t-card title="登录" header-bordered>
            <t-form ref="form" :data="formData" :colon="true" :label-width="0" @submit="onSubmit">
                <t-form-item name="username">
                    <t-input size="large" v-model="formData.username" clearable placeholder="请输入账户名">
                        <template #prefix-icon>
                            <desktop-icon />
                        </template>
                    </t-input>
                </t-form-item>

                <t-form-item name="password">
                    <t-input size="large" v-model="formData.password" type="password" placeholder="请输入密码">
                        <template #prefix-icon>
                            <lock-on-icon />
                        </template>
                    </t-input>
                </t-form-item>

                <t-form-item name="captcha">
                    <t-row>
                        <t-col span="8">
                            <t-input size="large" v-model="formData.captcha" placeholder="请输入验证码"></t-input>
                        </t-col>
                        <t-col span="4">
                            <t-image fit="contain" :src="captchaImg" @click="getCaptcha" :style="{ height: '40px' }" />
                        </t-col>
                    </t-row>
                </t-form-item>

                <t-form-item>
                    <t-button theme="primary" type="submit" block>登录</t-button>
                </t-form-item>
            </t-form>
        </t-card>
    </div>
</template>

<script>
import { DesktopIcon, LockOnIcon } from 'tdesign-icons-vue-next';
import { MD5 } from 'crypto-js';
import { useCounterStore } from '@/stores/counter'
import api from '@/api'

const store = useCounterStore()

export default {
    name: 'LoginView',
    components: {
        DesktopIcon,
        LockOnIcon
    },
    data() {
        return {
            formData: {
                username: '',
                password: '',
                captcha: '',
                captcha_id: ''
            },
            captchaImg: ''
        }
    },
    methods: {
        onSubmit() {
            // 使用axios发送请求
            let postData = Object.assign({}, this.formData)
            postData.password = MD5(this.formData.password).toString()
            api.login(postData).then((res) => {
                const data = res.data.data
                // 将token存储到store中
                store.login(data.token)
                this.$message('success', { content: '登录成功', duration: 3000 })
                this.$router.push({ name: 'dashboard' })
            }).catch((err) => {
                this.$message('error', { content: err.message, duration: 3000 })
                this.getCaptcha()
            })
        },
        // 获取验证码
        getCaptcha() {
            // 使用axios发送请求
            api.getCaptcha().then((res) => {
                const data = res.data.data
                this.captchaImg = data.captcha
                this.formData.captcha_id = data.id
            })
        }
    },
    mounted() {
        if (store.isLogin) {
            this.$router.push({ name: 'dashboard' })
        }
        this.getCaptcha()
    }
}
</script>


<style scoped>
/* 使login屏幕居中 */
.login {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
}
</style>