<script lang="ts" setup>
import SystemTool from './components/SystemTool.vue'
import 'element-plus/dist/index.css'
import "./assets/css/start.css"
import "./assets/css/style.css"
import DownloadTab from "./components/DownloadTab.vue";
import SinfferTab from "./components/SnifferTab.vue";
import SettingsTab from "./components/SettingsTab.vue";
import FeedbackTab from "./components/FeedbackTab.vue";
import AboutTab from "./components/AboutTab.vue";


</script>
<template>
  <el-container>
    <SystemTool/>
    <el-main class="main">
      <el-tabs>
        <DownloadTab/>
        <SinfferTab/>
        <SettingsTab/>
        <FeedbackTab/>
        <AboutTab/>
      </el-tabs>
    </el-main>
  </el-container>
</template>

<style>

</style>
