<template>
    <div>
        <t-layout>
            <t-aside>
                <leftNavigation />
            </t-aside>
            <t-layout>
                <t-header>
                    <Header />
                </t-header>
                <t-content class="t-content">
                    <slot></slot>
                </t-content>
                <t-footer>
                    <Footer />
                </t-footer>
            </t-layout>
        </t-layout>
    </div>
</template>

<script>
import leftNavigation from '@/components/LeftNavigation.vue'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'

export default {
    name: 'DashboardLayout',
    components: {
        leftNavigation,
        Header,
        Footer
    },
};
</script>

<style scoped>
/* t-content 边缘留白 适应手机端 */
.t-content {
    padding: 24px;
}
</style>