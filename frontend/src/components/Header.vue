<script lang="ts" setup>
import {BrowserOpenURL, Quit, Environment} from "../../wailsjs/runtime"
import { Connection, Link, Setting, Close } from "@element-plus/icons";

const info = await Environment();
let platform = info.platform;
let version = "1.0.0";
let tabPane = "";

function clickClose() {
  Quit();
}

function clickStartHookUrl () {

}

function clickLink (link: string) {
  BrowserOpenURL(link);
}

</script>
<template>
  <el-header v-if="platform !== 'darwin'" class="header" height='40px'>
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
               v-on:click.stop.prevent="clickLink">官网</a>
          </el-dropdown-item>
          <el-dropdown-item :icon="Link">
            <a class="mylink" href="https://github.com/AmatsuZero/GiveMeAnOfferGo/tree/m3u8"
               v-on:click.stop.prevent="clickLink">开放源码</a>
          </el-dropdown-item>
          <el-dropdown-item>
            <el-icon>
              <Connection/>
            </el-icon>
            <a class="mylink" href="https://github.com/AmatsuZero/GiveMeAnOfferGo/issues"
               v-on:click.stop.prevent="clickLink">点击加群</a>
          </el-dropdown-item>
          <el-dropdown-item :icon="Connection">
            <a class="mylink" href="#" v-on:click.stop.prevent="clickStartHookUrl">资源嗅探</a>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      <el-button :icon="Setting" class="button" v-on:click="tabPane = 'setting'"></el-button>
      <el-button class="button" :icon="Close" v-on:click="clickClose"></el-button>
    </div>
  </el-header>
</template>

<style scoped>


</style>