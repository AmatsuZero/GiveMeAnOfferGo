<template>
  <el-tab-pane label="合并视频片段">
    <el-row :gutter="8" style="margin-bottom: 5px;">
      <el-col :span="2" :offset="1">
        <span style="line-height: 40px;float: right;">视频片段</span>
      </el-col>
      <el-col :span="20">
        <div @drop="dropTSFiles" @dragover.prevent @dragenter.prevent>
          <el-input placeholder="选择一个包含TS流的目录或将所有TS文件拖拽至此" v-model="ts_dir"
                    clearable draggable="false">
            <i slot="suffix" class="el-input__icon el-icon-folder-opened"
               @click="clickSelectTSDir"></i>
          </el-input>
        </div>
      </el-col>
    </el-row>
    <el-row :gutter="8" style="margin-bottom: 5px;">
      <el-col :span="2" :offset="1">
        <span style="line-height: 40px;float: right;">合并名称</span>
      </el-col>
      <el-col :span="20">
        <el-input placeholder="[可空] 默认当前时间戳" v-model="mergeConfig.taskName" clearable
                  draggable="false"></el-input>
      </el-col>
    </el-row>
    <el-row :gutter="8" style="margin-bottom: 5px;">
      <el-col :span="7" :offset="3">
        <div>
          <el-radio style="line-height: 40px;" v-model="mergeConfig.mergeType"
                    label="">快速合并
          </el-radio>
          <el-radio style="line-height: 40px;" v-model="mergeConfig.mergeType"
                    label="transcoding">
            修复合并(慢|转码)</el-radio>
        </div>
      </el-col>
      <el-col :span="6" :offset="1">
            <span style="line-height: 40px;">
              共导入 {{mergeConfig.files.length}} 个视频片段
            </span>
      </el-col>
      <el-col :span="3" :offset="1">
        <el-button class="mybutton" type="primary" @click="clickStartMergeTS"
                   :disabled="mergeConfig.files.length===0">开始合并</el-button>
      </el-col>
      <el-col :span="2">
        <el-button class="mybutton" type="danger" @click="clickClearMergeTS"
                   :disabled="mergeConfig.files.length===0">清空</el-button>
      </el-col>
    </el-row>

    <el-row :gutter="8" style="margin-top: 10px;margin-bottom: 5px;"
            v-if="tsMergeStatus">
      <el-col :span="2" :offset="1">
            <span v-if="tsMergeStatus !=='success'"
                  style="line-height: 100px;float: right;">合并进度</span>
        <span v-if="tsMergeStatus ==='success'"
              style="line-height: 40px;float: right;">查看文件</span>
      </el-col>
      <el-col :span="6" v-if="tsMergeStatus !=='success'">
        <el-progress type="circle" :percentage="tsMergeProgress" status="success"
                     :width="100"></el-progress>
      </el-col>
      <el-col :span="4" v-else>
        <el-button class="mybutton" type="success" @click="clickOpenMergeTSDir">
          打开文件夹
        </el-button>
      </el-col>
      <el-col :span="4" v-if="tsMergeMp4Path">
        <el-button class="mybutton" type="success" @click="clickPlayMergeMp4">播放视频
        </el-button>
      </el-col>
    </el-row>

    <el-row :gutter="8" v-if="tsMergeStatus !=='success' && tsMergeStatus">
      <el-col :span="2" :offset="1">
        <span style="line-height: 40px;float: right;">日志</span>
      </el-col>
      <el-col :span="20">
        <el-input type="textarea" :rows="4" readonly placeholder=""
                  v-model="tsMergeStatus">
        </el-input>
      </el-col>
    </el-row>
  </el-tab-pane>
</template>

<script lang="ts">
export default {
  name: "MergeTask"
}
</script>
<script lang="ts" setup>

import {MergeFileType} from "../models";
import {reactive, ref} from "vue";
import {OpenSelectTsDir} from "../../wailsjs/go/main/App";
import {main} from "../../wailsjs/go/models";

let tsMergeMp4Path = '';
let tsMergeStatus = '';
let ts_dir = '';
let tsMergeProgress = 0;
let headers = ref("");
const mergeConfig = reactive(new main.MergeFilesConfig());

function clickOpenMergeTSDir () {
  OpenSelectTsDir("").then(files => {
    mergeConfig.files = files
  })
}

function clickPlayMergeMp4() {

}

function clickClearMergeTS() {

}

function clickSelectTSDir () {

}

function clickStartMergeTS() {}

function dropTSFiles() {

}
</script>

<style scoped>

</style>