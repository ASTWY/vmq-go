<!-- LeftNavigation.vue -->
<template>
    <div class="left-navigation">
        <!-- 导航栏内容 -->
        <t-menu theme="light" :default-value="this.$route.name" style="margin-right: 40px" height="550px">
            <template #logo>
                <h1><t-avatar shape="round">V</t-avatar> V免签Go 管理后台</h1>
            </template>
            <!-- 循环渲染路由  t-menu-item  -->
            <t-menu-item v-for="item in routers" :key="item.name" :value="item.name"
                @click="$router.push({ name: item.name })">
                <template #icon>
                    <t-icon :name="item.meta.icon" />
                </template>
                {{ item.meta.title }}
            </t-menu-item>
        </t-menu>
    </div>
</template>

<script>
// 从router中引入routers来获取路由信息
import { routers } from '@/router/index'; // 导入路由配置


export default {
    name: 'LeftNavigation',
    // 组件逻辑和功能
    components: {
        // 引入组件
    },
    data() {
        return {
            // 侧边栏路由 删去meta中没有title的路由
            routers: routers.filter((item) => item.meta && item.meta.title),
            // 当前路由
        }
    },
};
</script>
