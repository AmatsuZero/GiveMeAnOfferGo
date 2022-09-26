<template>
  <el-tab-pane>
    <span slot="label"><i class="el-icon-view"></i> 资源嗅探</span>
    <el-container>
      <el-container>
        <el-header class="navigatorTop">
          <el-input type="text" v-model="navigatorInput" placeholder="请输入 url" clearable @change="clickNaviagte"></el-input>
        </el-header>
        <el-main style="padding: 0;">
          <webview id="browser" style="height: 100%;" :src="navigatorUrl"
                   :useragent="currentUserAgent"></webview>
        </el-main>
      </el-container>
      <el-aside class="videoPanel" v-if="browserVideoUrls && browserVideoUrls.length" width="300px">
        <ul>
          <li class="videoItem" v-for="item in browserVideoUrls">
            <el-row :gutter="10">
              <el-col :span="6"><el-label>{{item.type}}</el-label></el-col>
              <el-col :span="18"><el-input type="text" readonly :value="item.url"></el-input></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="6"><el-label>HEADER</el-label></el-col>
              <el-col :span="18"><el-input type="text" readonly :value="item.headers"></el-input></el-col>
            </el-row>
            <el-row v-if="item.key" :gutter="10">
              <el-col :span="6"><el-label>KEY</el-label></el-col>
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

<script lang="ts">
export default {
  name: "SnifferTab",
  data: {
    browserVideoUrls: Array<string>(),
    currentUserAgent: "",
    navigatorInput: "",
    navigatorUrl: ""
  }
}
</script>

<style scoped>

</style>