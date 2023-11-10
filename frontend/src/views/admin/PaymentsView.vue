<template>
  <DashboardLayout>
    <!-- 添加一些过滤条件 -->
    <t-table row-key="id" :columns="columns" :data="tableData" table-layout="auto" :loading="isLoading"
      :pagination="pagination" bordered stripe lazy-load :filter-value="filterValue" @change="rehandleChange"
      @page-change="onPageChange" @filter-change="onFilterChange">
      <template #slot-type="{ row }">
        <!-- 微信或支付宝支付 由row.type决定 1: 微信 2: 支付宝 -->
        <t-tag shape="round" :theme="row.type === '1' ? 'success' : 'primary'" variant="light-outline">
          <template #icon>
            <t-icon :name="row.type === '1' ? 'logo-wechat-stroke' : 'logo-android'" />
            <!-- <t-icon v-if="row.type === '1'" :name="'logo-wechat-stroke'" /> -->
          </template>
          {{ row.type === "1" ? "微信" : "支付宝" }}
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
    </t-table>
  </DashboardLayout>
</template>

<script>
import DashboardLayout from "@/components/Dashboard.vue";
import api from "@/api";
import isNumber from "lodash/isNumber";

export default {
  name: "DashboardView",
  // 组件逻辑和功能
  components: {
    // 引入组件
    DashboardLayout,
  },
  data() {
    return {
      tableData: [],
      isLoading: false,
      pagination: {
        current: 1,
        pageSize: 10,
        defaultPageSize: 10,
        total: 100,
        defaultCurrent: 1,
      },
      stateNameListMap: {
        0: {
          theme: "warning",
          label: "未匹配到订单",
          icon: "error-circle",
        },
        1: {
          theme: "success",
          label: "已匹配到订单",
          icon: "check-circle",
        },
      },
      columns: [
        {
          colKey: "id",
          title: "ID",
        },
        {
          colKey: "price",
          title: "金额",
        },
        {
          colKey: "type",
          title: "类型",
          // 单选过滤配置
          filter: {
            // 过滤行中的列标题别名
            // label: '申请状态 A',
            type: "single",
            list: [
              { label: "支付宝", value: "1" },
              { label: "微信", value: "2" },
            ],
          },
          cell: "slot-type",
        },
        {
          colKey: "status",
          title: "订单状态",
          cell: "slot-state",
          // 多选过滤配置
          filter: {
            type: "multiple",
            resetValue: [],
            list: [
              { label: "全选", checkAll: true },
              { label: "未匹配到订单", value: 0 },
              { label: "已匹配到订单", value: 1 },
            ],
            // 是否显示重置取消按钮，一般情况不需要显示
            // showConfirmAndReset: true,
          },
        },
        {
          colKey: "createtime",
          title: "创建时间",
          cell: (h, { col, row }) => {
            // 1698388484 -> 2021-03-31 16:21:24
            let time = new Date(
              row.createtime * 1000 + 8 * 3600 * 1000
            ).toISOString();
            return h("span", time.slice(0, 10) + " " + time.slice(11, 19));
          },
        },
        {
          colKey: "orderId",
          title: "订单号",
        },
        {
          colKey: "metadata",
          title: "元数据",
        },
      ],
      filterValue: {},
      timer: setTimeout(() => { }, 0),
    };
  },
  methods: {
    getPaylogs(page, pageSize) {
      this.isLoading = true;
      api.getPaylogs(page, pageSize).then((res) => {
        let data = res.data.data;
        this.tableData = data.list;
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
            this.getPaylogs(this.pagination.current, this.pagination.pageSize);
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
        this.getPaylogs(pageInfo.current, pageInfo.pageSize);
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
    this.getPaylogs(
      this.pagination.defaultCurrent,
      this.pagination.defaultPageSize
    );
  },
};
</script>
