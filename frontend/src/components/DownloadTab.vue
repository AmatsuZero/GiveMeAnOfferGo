<script lang="ts">
import { ElMessage } from 'element-plus';
import { DownloadTask, MergeFileType, PlaylistItem } from "../models";
import {Open, OpenSelectTsDir, TaskAdd } from "../../wailsjs/go/main/App";
import {main} from "../../wailsjs/go/models";

export default {
  props: {
    tabPane: "",
    config_save_dir: ""
  },

  name: "download-tab",

  methods: {
    clickClearTask: function () {

    },

    openTask: function (link: string) {
        Open(link);
    },

    clickSelectM3U8: function () {

    },

    clickClearMergeTS: function () {

    },

    dropM3U8File: function () {

    },

    clickNewTaskMuti: function () {

    },

    clickSelectTSDir: function () {

    },

    clickStartMergeTS: function () {},

    dropTSFiles: function () {

    },
  },
}
</script>
<script lang="ts" setup>
import { Link, CirclePlusFilled, RemoveFilled, Download } from "@element-plus/icons";
import {DownloadTask, MergeFileType, PlaylistItem} from "../models";
import {OpenSelectTsDir, TaskAdd, Play } from "../../wailsjs/go/main/App";
import {ElMessage, ElIcon} from "element-plus";
import { ref, reactive } from 'vue';
import {main} from "../../wailsjs/go/models";
import { EventsOn, EventsEmit } from "../../wailsjs/runtime";

let ts_urls = Array<string>();
const allVideos = Array<DownloadTask>();
let downloadSpeed = '0 MB/s';
let tsMergeType = MergeFileType.Speed;
let dlg_newTask_visible = ref(false);
let tsMergeMp4Path = '';
let tsMergeStatus = '';
let m3u8_urls = '';
let addTaskMessage = ref('');
let ts_dir = '';
let playlists = ref(Array<PlaylistItem>());
let playlistUri = ref('');
let toolTipVisible = ref(false);
let tsMergeProgress = 0;
let headers = ref("");

function clickOpenMergeTSDir () {
  OpenSelectTsDir("").then(files => {
    ts_urls = files
  })
}

function m3u8UrlChange() {
  playlists.value = [];
  playlistUri.value = '';
  addTaskMessage.value = "请输入M3U8视频源";
}

function clickNewTask() {
  dlg_newTask_visible.value = true;
  parserTask.taskName = '';
  parserTask.url = '';
  m3u8UrlChange();
}

const parserTask = reactive(new main.ParserTask());
parserTask.delOnComplete = true;

function clickNewTaskOK() {
  if (playlistUri.value.length > 0) {
    EventsEmit('variant-selected', playlistUri.value);
  } else {
    TaskAdd(parserTask).catch(err => {
      console.error(err)
    });

    addTaskMessage.value = "正在检查链接...";
  }
}

EventsOn("select-variant", (data) => {
  playlists.value = data["Info"];
  addTaskMessage.value = data["Message"];
});

EventsOn("task-notify-create", data => {
  const item = data as DownloadTask;
  allVideos.push(item);
  dlg_newTask_visible.value = false;
});

function clickPlayMergeMp4() {

}

function stopItem(task: DownloadTask) {
  EventsEmit('stop-live-stream-download', task.url);
}

function deleteTask(task: DownloadTask) {

}

function playTask(task: DownloadTask) {
  const file = task.isDone ? task.videoPath : task.url;
  Play(file).catch(e => console.error(e));
}

function onHeadersChange(value: string | number) {
  if (value === undefined || typeof value !== "string") {
    return;
  }

  const headers = value.split("\n");
  const obj = new Map<string, string>();
  for (const header of headers) {
      const arr = header.split(":")
      obj.set(arr[0].trim(), arr[1].trim())
  }
  parserTask.headers = Object.fromEntries(obj);
}

</script>

<template>
  <el-tab-pane label="资源下载">
    <span slot="label">资源下载</span>
    <el-tabs type="border-card" v-model="tabPane">
      <el-tab-pane label="M3U8视频下载">
        <el-row :gutter="8">
          <el-col :span="3" :offset="1">
            <el-button class="mybutton" type="primary" :icon="CirclePlusFilled"
                       @click="clickNewTask">新建下载</el-button>
          </el-col>
          <el-col :span="4" :offset="1">
            <el-button class="mybutton" type="danger" :icon="RemoveFilled"
                       @click="clickClearTask">
              清空下载任务( {{allVideos.length}} )
            </el-button>
          </el-col>
          <el-col v-show="false" :span="4" :offset="11">
            <ElIcon class="speed"><Download/></ElIcon>
            {{downloadSpeed}}
          </el-col>
        </el-row>
        <ul class="list TaskList" style="overflow:auto">
          <li v-for="o in allVideos" :key="o.id" class="item">
            <div class="m3u8"></div>
            <div class="link"><input type="text" :value="o.url" readonly></div>

            <div class="name"><span>名称：</span><span class="value">
              {{o.taskName}}
            </span>
            </div>
            <div class="time"><span>时间：</span><span class="value">
              {{o.time}}
            </span>
            </div>
            <div class="status">
              <span>状态：</span><span class="value">
              {{o.status}}
            </span>
            </div>
            <div class="opt">
              <div class="top">
                <input class="opendir" type="button" value="打开文件夹" @click="openTask(o.videoPath)">
              </div>
              <div class="bottom">
                <input class="del" type="button" value="删除" @click="deleteTask(o)"/>
                <input class="play" type="button" value="播放" @click="playTask(o)"/>
                <input class="StartStop" type="button" value="停止" @click="stopItem(o)"/>
              </div>
            </div>
          </li>
          <el-alert v-if="allVideos.length === 0"
                    style="margin-top: 10px; height: 100px; line-height: 1;"
                    title="您还没有添加下载任务，在浏览器里嗅探到M3U8(HLS协议)视频流后，可以在这里缓存下载，快来试试吧。"
                    type="success"
                    effect="light" :closable="false" :center="true" show-icon/>
        </ul>
        <el-dialog title="新建下载任务"
                   width="60%"
                   v-model="dlg_newTask_visible"
                   :close-on-press-escape="false"
                   :close-on-click-modal="false"
        >
          <el-form label-width="80px" label-position="right" :model="parserTask">
            <el-form-item label="视频源">
              <div @drop="dropM3U8File" @dragover.prevent @dragenter.prevent>
                <el-input placeholder="输入在线网络视频源URL，或将M3U8文件拖拽至此" v-model="parserTask.url"
                          draggable="false" @input="m3u8UrlChange" :suffix-icon="Link">
                  <ElIcon><Link/></ElIcon>
                  <i slot="suffix" class="el-input__icon el-icon-folder-opened"
                     @click="clickSelectM3U8"></i>
                </el-input>
              </div>
            </el-form-item>
            <el-form-item v-if="playlists.length > 0" label="* 画质">
              <el-select v-model="playlistUri"
                         placeholder="选择视频源"
                         style="width: 100%;"
                         default-first-option="true">
                <el-option v-for="playlist in playlists" :key="playlist.uri"
                           :label="playlist.desc" :value="playlist.uri"/>
              </el-select>
            </el-form-item>
            <el-form-item label="任务名">
              <el-input type="text" placeholder="[可空] 默认当前时间戳" v-model="parserTask.taskName"/>
            </el-form-item>
            <el-form-item label="URL前缀"
                          v-if="parserTask.url.startsWith('file://')">
              <el-tooltip effect="light" placement="top-start">
                <div slot="content">
                  如果M3U8文件是从网上直接下载下来的，TS流是在网络上的且M3U8文件里没有URL(http)前缀，就需要填写<br /><br />
                  如果TS视频流在M3U8文件目录下，则不需要填写这块<br /></div>
                <el-input type="text"
                          placeholder="[可空] M3U8 URL前缀，例如：http://www.baidu.com/cdn/123456/"
                          :model-value="parserTask.prefix" />
              </el-tooltip>
            </el-form-item>
            <el-form-item label="附加头">
              <el-input type="textarea"
                        :rows="4"
                        v-model="headers"
                        placeholder="[可空] 请输入一行一个Header，例如:
                        Origin: http://www.host.com
                        Referer: http://www.host.com"
                        @change="onHeadersChange"
              />
            </el-form-item>
            <el-form-item label="私有KEY">
              <el-tooltip effect="light" placement="top-start" :visible="toolTipVisible">
                <template #content>
                  <div>
                    32位HEX格式KEY和32位IV值，IV值可空<br />示例：<br />ABABABABABBBBBBBABABABABABBBBBBBABABABABABBBBBBBABABABABABBBBBBB<br />或者：<br />ABABABABABBBBBBBABABABABABBBBBBB<br />
                  </div>
                </template>
                <el-input type="text"
                          placeholder="[可空] KEY和IV值(HEX格式)"
                          v-model="parserTask.keyIV"
                          @mouseenter="toolTipVisible = true"
                          @mouseleave="toolTipVisible = false"
                >
                </el-input>
              </el-tooltip>
            </el-form-item>
            <el-form-item label="合并完成">
              <el-checkbox v-model="parserTask.delOnComplete" label="删除TS片段"></el-checkbox>
            </el-form-item>
          </el-form>
          <template #footer>
            <div class="dialog-footer">
              <el-row :gutter="20">
                <el-col :span="11" :offset="4">
                  <el-alert :title="addTaskMessage" type="error" :closable="false"/>
                </el-col>
                <el-col :span="8" :offset="1">
                  <el-button type="primary" @click="clickNewTaskOK"
                             style="width: 100%;">确 定
                  </el-button>
                </el-col>
              </el-row>
            </div>
          </template>
        </el-dialog>
      </el-tab-pane>
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
            <el-input placeholder="[可空] 默认当前时间戳" v-model="parserTask.taskName" clearable
                      draggable="false"></el-input>
          </el-col>
        </el-row>
        <el-row :gutter="8" style="margin-bottom: 5px;">
          <el-col :span="7" :offset="3">
            <div>
              <el-radio style="line-height: 40px;" v-model="tsMergeType"
                        label="">快速合并
              </el-radio>
              <el-radio style="line-height: 40px;" v-model="tsMergeType"
                        label="transcoding">
                修复合并(慢|转码)</el-radio>
            </div>
          </el-col>
          <el-col :span="6" :offset="1">
            <span style="line-height: 40px;">
              共导入 {{ts_urls.length}} 个视频片段
            </span>
          </el-col>
          <el-col :span="3" :offset="1">
            <el-button class="mybutton" type="primary" @click="clickStartMergeTS"
                       :disabled="ts_urls.length===0">开始合并</el-button>
          </el-col>
          <el-col :span="2">
            <el-button class="mybutton" type="danger" @click="clickClearMergeTS"
                       :disabled="ts_urls.length===0">清空</el-button>
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
    </el-tabs>
  </el-tab-pane>
</template>

<style scoped>

</style>