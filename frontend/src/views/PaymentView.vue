<template>
    <t-dialog :footer="false" theme="info" :header="false" :preventScrollThrough="false" showOverlay
        v-model:visible="visible" :on-close="close">
        <p>{{ msg }}</p>
    </t-dialog>
    <div class="xh-title">收银台</div>
    <div class="qrbox clearfix">
        <div class="left">
            <div class="qrcon">
                <h5>
                    <img v-if="type === 2" src="../assets/image/alipay-logo.png" style="height:30px;" />
                    <img v-if="type === 1" src="../assets/image/wechat-logo.png" style="height:30px;" />
                </h5>
                <div class="title">剩余支付时间：{{ title }}</div>
                <div class="title"><strong style="color: red;">请务必按照当前页面显示金额支付</strong></div>
                <div class="price">￥{{ price }}</div>
                <div align="center" style="position:relative;">
                    <img :src="payQrcode" style="width: 230px;height: 230px;margin-bottom: 10px;" alt="">
                </div>
                <div class="bottom">
                    <p v-if="type === 2">请使用支付宝APP扫一扫扫描二维码支付<br>微信端请点右上角使用自带浏览器打开</p>
                    <p v-if="type === 1">请使用微信扫一扫扫描二维码支付</p>
                </div>
            </div>
        </div>
        <div class="sys">
            <img v-if="type === 1" src="../assets/image/wechat-sys.png" alt="">
            <img v-if="type === 2" src="../assets/image/alipay-sys.png" alt="">
        </div>
    </div>
</template>


<script>
import api from '@/api'
export default {
    data() {
        return {
            title: '倒计时',
            payId: '',
            orderID: '',
            visible: false,
            msg: '',
            price: '0',
            type: 1,
            payUrl: '',
            expectDate: 0,
            timer: null,
            payResultTimer: null,
        }
    },
    computed: {
        payQrcode() {
            return api.getQrcodeImg(this.payUrl)
        },
    },
    methods: {
        getOrder() {
            api.getOrder(this.orderID).then(res => {
                if (res.data.code == 1) {
                    this.price = res.data.data.reallyPrice
                    this.payId = res.data.data.payId
                    this.type = res.data.data.payType
                    this.payUrl = res.data.data.payUrl
                    this.expectDate = res.data.data.expectDate
                    this.payQrcode = res.data.data.payQrcode
                }
            })
        },
        close() {
            console.log('close')
            this.visible = false;
            // 跳转到来源页面 从哪里来回哪里去
            window.history.back();
            // 清除定时器
            clearInterval(this.timer);
            clearInterval(this.payResultTimer);
        },
        // 倒计时
        countDown() {
            // 1.获取当前时间 秒级时间戳
            const now = parseInt(new Date().getTime() / 1000);
            // 2.获取结束时间 秒级时间戳
            const endTime = new Date(this.expectDate).getTime();
            // 3.计算时间差
            const time = endTime - now;
            // 4.计算时分秒
            let hour = Math.floor(time / 60 / 60 % 24);
            let minute = Math.floor(time / 60 % 60);
            let second = Math.floor(time % 60);
            // 5.格式化时分秒
            hour = hour < 10 ? '0' + hour : hour;
            minute = minute < 10 ? '0' + minute : minute;
            second = second < 10 ? '0' + second : second;
            // 6.拼接时间 如果小时为0 不显示小时
            this.title = hour === '00' ? `${minute}:${second}` : `${hour}:${minute}:${second}`;
            // 7.判断时间是否到了
            if (time <= 0) {
                // 停止倒计时
                clearInterval(this.timer);
                this.title = "订单已过期"
                // 提示用户
                this.msg = '支付超时，请重新下单';
                this.payUrl = 'error';
                this.visible = true;
            }
        },
        // 获取支付结果
        getPayResult() {
            api.getOrderStatus(this.orderID).then(res => {
                if (res.data.code == 1) {
                    let status = res.data.data.state // -1过期 0未支付 1已支付 
                    let returnUrl = res.data.data.returnUrl // 支付成功后跳转的页面
                    if (status == -1) {
                        this.msg = '支付超时，请重新下单';
                        this.visible = true;
                    } else if (status == 1) {
                        this.msg = '支付成功';
                        this.visible = true;
                        setTimeout(() => {
                            window.location.href = returnUrl
                        }, 1000)
                    }
                }
            })
        }
    },
    mounted() {
        // 1.获取订单号 
        // 获取订单号
        this.orderID = this.$route.params.id;
        this.getOrder();
        // 2.开启定时器
        this.timer = setInterval(this.countDown, 1000);
        // 3.获取支付结果 定时器
        this.payResultTimer = setInterval(this.getPayResult, 1000);
    }
}

</script>

<style>
* {
    margin: 0;
    padding: 0;
}

body {
    background: #f2f2f4;
}

.clearfix:after {
    content: ".";
    display: block;
    height: 0;
    clear: both;
    visibility: hidden;
}

.clearfix {
    display: inline-block;
}

* html .clearfix {
    height: 1%;
}

.clearfix {
    display: block;
}

.xh-title {
    height: 75px;
    line-height: 75px;
    text-align: center;
    font-size: 30px;
    font-weight: 300;
    border-bottom: 2px solid #eee;
    background: #fff;
}

.qrbox {
    max-width: 900px;
    margin: 0 auto;
    padding: 85px 20px 20px 50px;
}

.qrbox .left {
    width: 40%;
    float: left;
    display: block;
    margin: 0px auto;
}

.qrbox .left .qrcon {
    border-radius: 10px;
    background: #fff;
    overflow: visible;
    text-align: center;
    padding-top: 25px;
    color: #555;
    box-shadow: 0 3px 3px 0 rgba(0, 0, 0, .05);
    vertical-align: top;
    -webkit-transition: all .2s linear;
    transition: all .2s linear;
}

.qrbox .left .qrcon .logo {
    width: 100%;
}

.qrbox .left .qrcon .title {
    font-size: 16px;
    margin: 10px auto;
    width: 90%;
}

.qrbox .left .qrcon .price {
    font-size: 22px;
    margin: 0px auto;
    width: 100%;
}

.qrbox .left .qrcon .bottom {
    border-radius: 0 0 10px 10px;
    width: 100%;
    background: #32343d;
    color: #f2f2f2;
    padding: 15px 0px;
    text-align: center;
    font-size: 14px;
}

.qrbox .sys {
    width: 60%;
    float: right;
    text-align: center;
    padding-top: 20px;
    font-size: 12px;
    color: #ccc
}

.qrbox img {
    max-width: 100%;
}

@media (max-width : 767px) {
    .qrbox {
        padding: 20px;
    }

    .qrbox .left {
        width: 90%;
        float: none;
    }

    .qrbox .sys {
        display: none;
    }
}
</style>