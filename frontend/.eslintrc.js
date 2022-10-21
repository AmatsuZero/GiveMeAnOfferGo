module.exports = {
    root: true,
    env: {
        browser: true,
        node: true,
        es6: true,
    },
    parser: 'vue-eslint-parser', // .vue 文件解析器，帮助ESLint解析vue文件
    plugins: ['@typescript-eslint'],
    extends: [
        'eslint:recommended',  // ESLint官方预定义规则集
        'plugin:vue/recommended',  // Vue规则集
        'plugin:@typescript-eslint/recommended',  // TypeScript规则集
        'plugin:prettier/recommended',  // Prettier规则集
    ],
    // 自定义规则，配置后会覆盖extends中已有的规则，官方规则配置手册: https://eslint.bootcss.com/docs/rules/
    rules: {
    },
    parserOptions: {
        parser: '@typescript-eslint/parser', // .ts 文件解析器，帮助ESLint解析typescript文件
    }
};