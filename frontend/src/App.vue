<script lang="ts" setup>
import SystemTool from './components/SystemTool.vue'
import 'element-plus/dist/index.css'
import "./assets/css/start.css"
import "./assets/css/style.css"

</script>
<template>
  <body class="heisir">
  <div id="app">
  <el-container>
    <el-header v-if="platform !== 'darwin'" class="header" height='40px'>
      <el-row>
        <el-col :span="14">
          <div class="title">{{'M3U8-Downloader v'}}</div>
        </el-col>
        <SystemTool class="systemTool"/>
      </el-row>
    </el-header>
  </el-container>
  </div>
  </body>
</template>

<style>

</style>
