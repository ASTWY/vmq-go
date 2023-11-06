module.exports = {
    base: '/vmq-go/',
    title: 'V免签Go',
    description: '使用文档',
    port: '8080',
    head: [
        ['link', { rel: 'icon', href: 'logo.png' }]
    ],
    markdown: {
        lineNumbers: true
    },
    theme: 'antdocs',
    themeConfig: {
        nav: [
            { text: '首页', link: '/' },
            { text: '指南', link: '/guide/' },
            { text: '配置', link: '/config/' },
            { text: '更新日志', link: '/changelog/' },
            { text: 'Github', link: 'https://github.com/astwy/vmq-go' },
        ],
        logo: '/assets/logo.png',
        sidebar: 'auto',
        // nav: require("./nav.js"),
        nextLinks: true,
        prevLinks: true,
        docsRepo: 'astwy/vmq-go',
        docsBranch: 'docs/docs',
        editLinks: true,
        editLinkText: '在 GitHub 上编辑此页 ！',
        lastUpdated: '上次更新',
        smoothScroll: true,
    }
}