<template>
  <el-tab-pane label="M3U8批量下载">
    <el-form label-width="130px" style="margin-right: 65px;">
      <el-form-item label="批量视频源" label-position="right">
        <el-input type="textarea" :rows="6" placeholder="请输入一行一个M3U8视频源，格式：视频源----任务名（可空）,例如：
                https://host/index1.m3u8----第一个视频
                https://host/index2.m3u8----第二个视频
                https://host/index3.m3u8" v-model="m3u8_urls"></el-input>
      </el-form-item>
      <el-form-item label="附加头" label-position="right">
        <el-input type="textarea" :rows="4" placeholder="[可空] 请输入一行一个Header，例如：
                Origin: http://www.host.com
                Referer: http://www.host.com" v-model="parserTask.headers"/>
      </el-form-item>
      <el-form-item label="合并完成">
        <el-checkbox v-model="parserTask.delOnComplete" label="删除TS片段"/>
      </el-form-item>
      <el-form-item>
        <el-button class="mybutton" type="primary" @click="clickNewTaskMuti">批量下载
        </el-button>
      </el-form-item>
    </el-form>
  </el-tab-pane>
</template>

<script lang="ts">
export default {
  name: "MultiTasks"
}
</script>
<script lang="ts" setup>
import {reactive} from "vue";
import {main} from "../../wailsjs/go/models";

let m3u8_urls = '';
const parserTask = reactive(new main.ParserTask());
parserTask.delOnComplete = true;

function clickNewTaskMuti() {

}


</script>

<style scoped>

</style>