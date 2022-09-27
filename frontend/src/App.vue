<script lang="ts" setup>
import Header from './components/Header.vue'
import 'element-plus/dist/index.css'
import "./assets/css/start.css"
import "./assets/css/style.css"
import DownloadTab from "./components/DownloadTab.vue";
import SnifferTab from "./components/SnifferTab.vue";
import SettingsTab from "./components/SettingsTab.vue";
import FeedbackTab from "./components/FeedbackTab.vue";
import AboutTab from "./components/AboutTab.vue";

</script>
<template>
  <el-container direction="vertical">
    <Suspense>
      <Header/>
    </Suspense>
    <el-main class="main">
      <el-tabs tab-position="left" type="border-card" style="height: 100%;">
        <DownloadTab/>
        <SnifferTab/>
        <SettingsTab/>
        <FeedbackTab/>
        <AboutTab/>
      </el-tabs>
    </el-main>
  </el-container>
</template>

<style>

</style>
