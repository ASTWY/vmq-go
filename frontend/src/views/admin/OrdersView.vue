<template>
    <DashboardLayout>
        <!-- 添加一些过滤条件 -->
        <t-table row-key="id" :columns="columns" :data="data" table-layout="auto" :loading="isLoading"
            :expand-on-row-click="true" :pagination="pagination" bordered stripe lazy-load :filter-value="filterValue"
            @change="rehandleChange" @page-change="onPageChange" @filter-change="onFilterChange">
            <template #slot-type="{ row }">
                <!-- 微信或支付宝支付 由row.type决定 1: 微信 2: 支付宝 -->
                <t-tag shape="round" :theme="row.type === 1 ? 'success' : 'primary'" variant="light-outline">
                    <template #icon>
                        <t-icon :name="row.type === 1 ? 'logo-wechat-stroke' : 'logo-android'" />
                    </template>
                    {{ row.type === 1 ? "微信" : "支付宝" }}
                </t-tag>
            </template>
            <template #slot-state="{ row }">
                <t-tag shape="round" :theme="stateNameListMap[row.state].theme" variant="light-outline">
                    <template #icon>
                        <t-icon :name="stateNameListMap[row.state].icon" />
                    </template>
                    {{ stateNameListMap[row.state].label }}
                </t-tag>
            </template>
            <template #slot-param="{ row }">
                <!-- 回调参数 -->
                <div class="text-container">
                    <p class="text">{{ row.param }}</p>
                </div>
            </template>
            <template #slot-action="{ row }">
                <!-- 补单 -->
                <t-button v-if="row.state !== 3" theme="primary" size="small" shape="round"
                    @click="replenish(row)">补单</t-button>
                <!-- 删除 -->
                <t-button theme="danger" size="small" shape="round" @click="deleteOrder(row)">删除</t-button>
            </template>
            <template #expandedRow="{ row }">
                <div class="more-detail">
                    <p class="title"><b>订单参数:</b></p>
                    <p class="content">{{ row.param }}</p>
                    <br />
                    <p class="title"><b>回调结果:</b></p>
                    <p class="content">{{ row.notifyRes }}</p>
                </div>
            </template>
        </t-table>
    </DashboardLayout>
</template>

<script>
import isNumber from "lodash/isNumber";
import DashboardLayout from "@/components/Dashboard.vue";
import api from "@/api";

export default {
    name: "DashboardView",
    components: {
        DashboardLayout,
    },
    data() {
        return {
            data: [],
            isLoading: false,
            pagination: {
                current: 1,
                pageSize: 10,
                defaultPageSize: 10,
                total: 100,
                defaultCurrent: 1,
            },
            stateNameListMap: {
                "-1": {
                    theme: "danger",
                    label: "已过期",
                    icon: "close-circle",
                },
                0: {
                    theme: "warning",
                    label: "未支付",
                    icon: "time",
                },
                1: {
                    theme: "success",
                    label: "已支付",
                    icon: "minus-circle",
                },
                2: {
                    theme: "danger",
                    label: "回调失败",
                    icon: "info-circle",
                },
                3: {
                    theme: "success",
                    label: "已完成",
                    icon: "check-circle",
                },
            },
            columns: [
                {
                    colKey: "id",
                    title: "ID",
                },
                {
                    colKey: "orderId",
                    title: "订单号",
                },
                {
                    colKey: "payId",
                    title: "商户订单号",
                },
                {
                    colKey: "type",
                    title: "订单类型",
                    // 单选过滤配置
                    filter: {
                        // 过滤行中的列标题别名
                        // label: '申请状态 A',
                        type: "single",
                        list: [
                            { label: "微信", value: 1 },
                            { label: "支付宝", value: 2 },
                        ],
                    },
                    cell: "slot-type",
                },
                {
                    colKey: "price",
                    title: "订单金额",
                },
                {
                    colKey: "reallyPrice",
                    title: "实际支付金额",
                },
                {
                    colKey: "state",
                    title: "订单状态",
                    cell: "slot-state",
                    // 多选过滤配置
                    filter: {
                        type: "multiple",
                        resetValue: [],
                        list: [
                            { label: "全选", checkAll: true },
                            { label: "已过期", value: -1 },
                            { label: "未支付", value: 0 },
                            { label: "已支付", value: 1 },
                            { label: "回调失败", value: 2 },
                            { label: "已完成", value: 3 },
                        ],
                        // 是否显示重置取消按钮，一般情况不需要显示
                        // showConfirmAndReset: true,
                    },
                },
                {
                    colKey: "createDate",
                    title: "创建时间",
                    cell: (h, { col, row }) => {
                        // 1698388484 -> 2021-03-31 16:21:24
                        let time = new Date(
                            row.createDate * 1000 + 8 * 3600 * 1000
                        ).toISOString();
                        return h("span", time.slice(0, 10) + " " + time.slice(11, 19));
                    },
                },
                {
                    colKey: "closeDate",
                    title: "关闭时间",
                    cell: (h, { col, row }) => {
                        // 1698388484 -> 2021-03-31 16:21:24
                        let time = new Date(
                            row.createDate * 1000 + 8 * 3600 * 1000
                        ).toISOString();
                        return h("span", time.slice(0, 10) + " " + time.slice(11, 19));
                    },
                },
                {
                    colKey: "action",
                    title: "操作",
                    cell: "slot-action",
                },
            ],
            filterValue: {},
            timer: setTimeout(() => { }, 0),
        };
    },
    methods: {
        getOrdes(page, pageSize) {
            this.isLoading = true;
            api.getOrders(page, pageSize).then((res) => {
                let data = res.data.data;
                this.data = data.list;
                this.pagination.total = data.total;
                this.isLoading = false;
            });
        },
        replenish(row) {
            api.replenish(row.orderId).then(() => {
                // 一秒后重新获取数据
                setTimeout(() => {
                    // 如果有过滤条件，需要重新过滤数据
                    if (Object.keys(this.filterValue).length > 0) {
                        this.filterData(this.filterValue);
                    } else {
                        this.getOrdes(this.pagination.current, this.pagination.pageSize);
                    }
                }, 1000);
            });
        },
        deleteOrder(row) {
            api.deleteOrder(row.orderId).then(() => {
                // 一秒后重新获取数据
                setTimeout(() => {
                    // 如果有过滤条件，需要重新过滤数据
                    if (Object.keys(this.filterValue).length > 0) {
                        this.filterData(this.filterValue);
                    } else {
                        this.getOrdes(this.pagination.current, this.pagination.pageSize);
                    }
                }, 1000);
            });
        },
        onPageChange(pageInfo) {
            // {current: 2, previous: 1, pageSize: 10}
            this.pagination.current = pageInfo.current;
            this.pagination.pageSize = pageInfo.pageSize;
            // 如果有过滤条件，需要重新过滤数据
            if (Object.keys(this.filterValue).length > 0) {
                this.filterData(this.filterValue);
            } else {
                this.getOrdes(pageInfo.current, pageInfo.pageSize);
            }
        },
        rehandleChange(changeParams, triggerAndData) { },
        filterData(filters) {
            this.timer = setTimeout(() => {
                clearTimeout(this.timer); // 清除定时器
                // 根据过滤条件过滤数据
                this.isLoading = true;
                api
                    .getOrders(this.pagination.current, this.pagination.pageSize)
                    .then((res) => {
                        let data = res.data.data;
                        let newData = data.list.filter((item) => {
                            let result = true;
                            if (isNumber(filters.type)) {
                                result = item.type === filters.type;
                            }
                            if (result && filters.state && filters.state.length) {
                                result = filters.state.includes(item.state);
                            }
                            return result;
                        });
                        this.data = newData;
                        this.isLoading = false;
                    });
            }, 100);
        },
        onFilterChange(filters, ctx) {
            this.filterValue = {
                ...filters,
                state: filters.state || [],
            };
            this.filterData(filters);
        },
    },
    mounted() {
        this.getOrdes(
            this.pagination.defaultCurrent,
            this.pagination.defaultPageSize
        );
    },
};
</script>

<style scoped>
.text-container {
    /* 设置容器宽度，适应你的需求 */
    overflow: hidden;
    /* 超出容器宽度的文本内容将被隐藏 */
    white-space: nowrap;
    /* 禁止文本换行 */
    text-overflow: ellipsis;
    /* 当文本超出容器宽度时，用省略号来表示被截断的部分 */
}

.text-container:hover .text {
    white-space: normal;
    /* 鼠标悬浮该容器时，文本内容换行显示 */
    text-overflow: unset;
    /* 取消省略号显示 */
}

:deep([class*="t-table-expandable-icon-cell"]) .t-icon {
    background-color: transparent;
}

.link {
    cursor: pointer;
    margin-right: 15px;
}

.more-detail {
    line-height: 22px;

    >p {
        display: inline-block;
        margin: 4px 0;
    }

    >p.title {
        width: 120px;
    }
}
</style>
