<script lang="ts" setup>
import { BrowserOpenURL } from "../../wailsjs/runtime"
import { Link, Connection, Setting, Close } from "@element-plus/icons-vue";
  function clickAClick(link: PointerEvent) {
    console.log(link)
    BrowserOpenURL("link.target.href")
  }
</script>

<template>
    <el-dropdown trigger="click" slot="dropdown">
      <span class="el-dropdown-link">关于</span>
        <el-dropdown-menu>
          <template #dropdown>
            <el-dropdown-item :icon="Link">
              <a class="mylink" href="https://github.com/AmatsuZero/GiveMeAnOfferGo/wiki"
                 v-on:click.stop.prevent="clickAClick">官网</a>
            </el-dropdown-item>
            <el-dropdown-item :icon="Link">
              <a class="mylink" href="https://github.com/AmatsuZero/GiveMeAnOfferGo/tree/m3u8"
                 v-on:click.stop.prevent="clickAClick">开放源码</a>
            </el-dropdown-item>
            <el-dropdown-item :icon="Link">
              <a class="mylink" href="https://github.com/AmatsuZero/GiveMeAnOfferGo/issues"
                 v-on:click.stop.prevent="clickAClick">点击加群</a>
            </el-dropdown-item>
            <el-dropdown-item :icon="Connection">
              <a class="mylink" href="#" v-on:click.stop.prevent="clickStartHookUrl">资源嗅探</a>
            </el-dropdown-item>
          </template>

        </el-dropdown-menu>
    </el-dropdown>
  <el-button class="button" :icon="Setting" v-on:click="tabPane = 'setting'"></el-button>
  <el-button class="button" :icon="Close" v-on:click="clickClose"></el-button>
</template>

<style scoped>


</style>