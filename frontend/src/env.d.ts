/// <reference types="vite/client" />

declare module '*.vue' {
    import Vue from "vue";
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
    const component: Vue.DefineComponent<{}, {}, any>
    export default component
}
