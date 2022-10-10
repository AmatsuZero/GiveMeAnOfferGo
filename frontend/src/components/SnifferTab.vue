<script lang="ts">
import { VideoItem } from "../models"

export default {
  name: "SnifferTab",
  data() {
    return {
      browserVideoUrls: Array<VideoItem>(),
      currentUserAgent: "",
      navigatorUrl: ""
    }
  }
}
</script>

<script lang="ts" setup>
import { ref } from 'vue';
import {InterceptItem} from "../models";
import {EventsOn} from "../../wailsjs/runtime";
import { SniffLinks } from "../../wailsjs/go/main/App"

const interceptURLs = ref(Array<InterceptItem>());

const navigatorInput = ref("")

function clickNavigate() {
  SniffLinks(navigatorInput.value);
}

EventsOn('intercept-url', function (item) {
  const i = new InterceptItem();
  i.url = item["url"];
  i.status = item["status"];
  interceptURLs.value.push(i);
});

</script>

<template>
  <el-tab-pane label="资源嗅探">
    <span slot="label"><i class="el-icon-view"></i> 资源嗅探</span>
    <el-container>
      <el-container>
        <el-header class="navigatorTop">
          <el-input type="text" v-model="navigatorInput" placeholder="请输入 url" clearable @change="clickNavigate"></el-input>
        </el-header>
        <el-main style="padding: 0;">
          <el-table :data="interceptURLs" style="width: 100%">
            <el-table-column prop="status" label="Status"/>
            <el-table-column prop="url" label="Url"/>
          </el-table>
        </el-main>
      </el-container>
      <el-aside class="videoPanel" v-if="browserVideoUrls && browserVideoUrls.length" width="300px">
        <ul>
          <li class="videoItem" v-for="item in browserVideoUrls">
            <el-row :gutter="10">
              <el-col :span="6"><label>{{item.type}}</label></el-col>
              <el-col :span="18"><el-input type="text" readonly :value="item.url"></el-input></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="6"><label>HEADER</label></el-col>
              <el-col :span="18"><el-input type="text" readonly :value="item.headers"></el-input></el-col>
            </el-row>
            <el-row v-if="item.key" :gutter="10">
              <el-col :span="6"><label>KEY</label></el-col>
              <el-col :span="18"><el-input type="text" readonly :value="item.key"></el-input></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :offset="8" :span="8"><el-button>浏览</el-button></el-col>
              <el-col :span="8"><el-button>下载</el-button></el-col>
            </el-row>
          </li>
        </ul>
      </el-aside>
    </el-container>
  </el-tab-pane>
</template>

<style scoped>

</style>