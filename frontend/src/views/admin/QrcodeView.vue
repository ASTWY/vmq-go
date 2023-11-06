<template>
    <DashboardLayout>
        <t-tabs :default-value="1">
            <t-tab-panel :value="1" label="微信">
                <t-table :data="tableDataWechat" :columns="columns">
                    <template #slot-qrcode="{ row }">
                        <t-image fit="contain" :src="row.img" :style="{ width: '100px', height: '100px' }" />
                    </template>
                    <template #slot-action="{ row }">
                        <t-button type="primary" @click="handleClick(row)">删除</t-button>
                    </template>
                    <template #footerSummary>
                        <div class="t-table__row-filter-inner">
                            <t-upload ref="uploadRef" v-model="files" :request-method="requestMethod"
                                :on-fail="handleRequestFail"><t-button>添加收款二维码</t-button></t-upload>
                        </div>
                    </template>
                </t-table>
            </t-tab-panel>
            <t-tab-panel :value="2" label="支付宝">
                <t-table :data="tableDataAlipay" :columns="columns">
                    <template #slot-qrcode="{ row }">
                        <t-image fit="contain" :src="row.img" :style="{ width: '100px', height: '100px' }" />
                    </template>
                    <template #slot-action="{ row }">
                        <t-button type="primary" @click="handleClick(row)">删除</t-button>
                    </template>
                    <template #footerSummary>
                        <div class="t-table__row-filter-inner">
                            <t-upload ref="uploadRef" v-model="files" :request-method="requestMethod"
                                :on-fail="handleRequestFail"><t-button>添加收款二维码</t-button></t-upload>
                        </div>
                    </template>
                </t-table>
            </t-tab-panel>
        </t-tabs>
        <t-dialog :visible="visible" :close-btn="true" confirm-btn="确认" cancel-btn="取消" :on-close="close"
            :on-confirm="onConfirm">
            <template #header>添加收款二维码</template>
            <template #body>
                <t-input label="支付链接：" v-model="form.payUrl" readonly />
                <t-input label="价格：" suffix="元" v-model="form.price" />
            </template>
        </t-dialog>
    </DashboardLayout>
</template>

<script>
import api from '@/api'
import DashboardLayout from '@/components/Dashboard.vue'

export default {
    name: 'QrcodeView',
    components: {
        DashboardLayout,
    },
    data() {
        return {
            visible: false,
            form: {
                payUrl: '',
                price: '',
                type: '-1',
            },
            files: [],
            tabledata: [],
            columns: [{
                colKey: 'id',
                title: 'ID',
                width: 80,
            },
            {
                colKey: 'price',
                title: '金额',
            },
            {
                colKey: 'payUrl',
                title: '支付链接',
            },
            {
                title: '支付二维码',
                cell: 'slot-qrcode'
            },
            {
                title: '操作',
                cell: 'slot-action',
            }
            ]
        }
    },
    computed: {
        tableDataWechat() {
            return this.tabledata.filter(item => item.type === 1)
        },
        tableDataAlipay() {
            return this.tabledata.filter(item => item.type === 2)
        },
    },
    methods: {
        close() {
            this.visible = false
        },
        onConfirm() {
            api.addQrcode(
                this.form
            ).then(() => {
                this.visible = false
                this.form.payUrl = ''
                this.form.price = ''
                this.form.type = '-1'
                this.files = []
                this.getQrcodeList()
            })
        },
        getQrcodeList() {
            api.getQrcodes().then(res => {
                let list = res.data.data
                list.forEach(item => {
                    item.img = api.getQrcodeImg(item.payUrl)
                })
                this.tabledata = list
            })
        },
        handleClick(row) {
            api.deleteQrcode(row.id).then(() => {
                this.getQrcodeList()
            })
        },
        requestMethod(file) {
            console.log(file)
            return new Promise((resolve) => {
                let formData = new FormData()
                formData.append('file', file.raw)
                api.postQrcode(formData).then(res => {
                    console.log(res.data.data)
                    let content = res.data.data.content
                    // 如果content开头是wxp://或者https://qr.alipay.com/，则是正确的二维码
                    if (content.startsWith('wxp://')) {
                        this.form.payUrl = content
                        this.form.type = '1'
                        this.visible = true
                        resolve({ status: 'success' })
                    } else if (content.startsWith('https://qr.alipay.com/')) {
                        this.form.payUrl = content
                        this.form.type = '2'
                        this.visible = true
                        rsolve({ status: 'success' })
                    } else {
                        resolve({ status: 'fail', error: '不是正确的二维码' })
                    }
                })
            });
        },
        handleRequestFail(err) {
            console.log(err);
        },
    },
    mounted() {
        this.getQrcodeList()
    }
}
</script>
