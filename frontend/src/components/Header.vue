<template>
    <h1>{{ title }}</h1>
</template>

<script>
export default {
    name: 'Header',
    // 组件逻辑和功能
    components: {
        // 引入组件
    },
    data() {
        return {
            // 组件数据 从路由中获取title
            title: this.$route.meta.title,
        }
    },
};
</script>