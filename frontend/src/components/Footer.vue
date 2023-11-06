<template>
    <div class="footer">
        <t-link theme="default" size="large" href="https://github.com/astwy/vmq-go">
            <template #suffix-icon>
                <t-icon name="jump" />
            </template>
            <t-icon name="logo-github" /> vmq-go
        </t-link>
    </div>
</template>

<script>
export default {
    name: 'Footer',
};
</script>

<style scoped>
.footer {
    text-align: center;
    padding: 20px;
    color: white;
}
</style>