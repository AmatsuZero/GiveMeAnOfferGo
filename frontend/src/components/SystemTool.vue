<script lang="ts">
import { BrowserOpenURL } from "../../wailsjs/runtime"
import { Link, Connection, Setting, Close } from "@element-plus/icons-vue";

export default {
  name: "SystemTool",
  data: {
    version: "",
  },
  methods: {
    clickClose: function () {

    },

    clickStartHookUrl: function () {

    },

    clickAClick: function (link: PointerEvent) {
      console.log(link)
      BrowserOpenURL("link.target.href")
    }
  }
}

</script>

<template>
  <el-header v-if="platform != 'darwin'" class="header" height='40px'>
    <el-row>
      <el-col :span="14">
        <div class="title">
          {{'M3U8-Downloader v' + version}}
        </div>
      </el-col>
    </el-row>

    <div class="systemTool">
      <el-dropdown trigger="click">
        <span class="el-dropdown-link">关于</span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item :icon="Link">
            <a class="mylink" href="https://tools.heisir.cn/HLSDownload/"
               v-on:click.stop.prevent="clickAClick">官网</a>
          </el-dropdown-item>
          <el-dropdown-item icon="el-icon-link">
            <a class="mylink" href="https://github.com/HeiSir2014/M3U8-Downloader"
               v-on:click.stop.prevent="clickAClick">开放源码</a>
          </el-dropdown-item>
          <el-dropdown-item icon="el-icon-connection">
            <a class="mylink" href="https://jq.qq.com/?_wv=1027&k=nhFrZBS0"
               v-on:click.stop.prevent="clickAClick">点击加群</a>
          </el-dropdown-item>
          <el-dropdown-item :icon="Connection">
            <a class="mylink" href="#" v-on:click.stop.prevent="clickStartHookUrl">资源嗅探</a>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      <el-button class="button" :icon="Setting" v-on:click="tabPane = 'setting'"></el-button>
      <el-button class="button" :icon="Close" v-on:click="clickClose"></el-button>
    </div>
  </el-header>
</template>

<style scoped>


</style>