<script lang="ts">
export default {
  name: "task-list"
}
</script>
<script lang="ts" setup>
import { Link, CirclePlusFilled, RemoveFilled, Download, ArrowDown, VideoPlay, Stopwatch, FolderOpened, Refresh } from "@element-plus/icons";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {DownloadTask, DownloadTaskState, PlaylistItem} from "../models";
import {computed, reactive, ref} from "vue";
import {ClearTasks, Open, Play, RemoveTaskNotifyItem, TaskAdd} from "../../wailsjs/go/main/App";
import {main} from "../../wailsjs/go/models";

let headers = ref('');
let playlists = ref(Array<PlaylistItem>());
let playlistUri = ref('');
let toolTipVisible = ref(false);
let addTaskMessage = ref('');
let dlg_newTask_visible = ref(false);
const allVideos = ref(Array<DownloadTask>());
let downloadSpeed = '0 MB/s';

const search = ref('')
const filterTableData = computed(() =>
    allVideos.value.filter(
        (data) =>
            !search.value ||
            data.taskName.toLowerCase().includes(search.value.toLowerCase()) ||
            data.url.toLowerCase().includes(search.value.toLowerCase())
    )
);

const tableRowClassName = ({
                             row,
                             rowIndex,
                           }: {
  row: DownloadTask
  rowIndex: number
}) => {
  switch (row.state) {
    case DownloadTaskState.Processing:
      return 'warning-row';
    case DownloadTaskState.Idle:
      return 'primary-row';
    case DownloadTaskState.Error:
      return 'error-row'
    case DownloadTaskState.Done:
      return 'success-row'
  }
  return ''
}

EventsOn("task-notify-update", data => updateTaskItem(data));
EventsOn("task-notify-end", data => updateTaskItem(data));

function deleteTask(task: DownloadTask) {
  stopItem(task);
  const idx = allVideos.value.findIndex(video => video.url === task.url);
  if (idx === -1) {
    return;
  }
  allVideos.value.splice(idx, 1);
  RemoveTaskNotifyItem(task);
}

function playTask(task: DownloadTask) {
  const file = task.isDone ? task.videoPath : task.url;
  Play(file).catch(e => console.error(e));
}

function stopItem(task: DownloadTask) {
  EventsEmit('task-stop', task.url);
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
    TaskAdd(parserTask).catch(() => {
      addTaskMessage.value = "资源解析失败";
    });

    addTaskMessage.value = "正在检查链接...";
  }
}

EventsOn("select-variant", (data) => {
  playlists.value = data["Info"];
  addTaskMessage.value = data["Message"];
});

EventsOn("task-notify-create", (data) => {
  const item = data as DownloadTask;
  allVideos.value.push(item);
  dlg_newTask_visible.value = false;
});

EventsOn('task-add-reply', (message) => {
  const code = message['code'];
  if (code === -1) {
    addTaskMessage.value = message['message'];
  }
});

EventsOn("task-notify-update", data => updateTaskItem(data));

EventsOn("task-notify-end", data => updateTaskItem(data));

function updateTaskItem(data: any) {
  const item = data as DownloadTask;
  const idx = allVideos.value.findIndex(video => video.url === item.url);
  if (idx === -1) {
    return
  }
  allVideos.value[idx].isDone = item.isDone;
  allVideos.value[idx].status = item.status;
  allVideos.value[idx].videoPath = item.videoPath;
  allVideos.value[idx].state = item.state;
}

function clickClearTask() {
  allVideos.value.forEach(video => deleteTask(video));
  ClearTasks();
}

function openTask(link: string) {
  Open(link);
}

function dropM3U8File() {

}

function clickSelectM3U8() {

}

function handleDownloadTask(ev: any) {
  if (ev['play'] !== undefined) {
    playTask(ev['play'] as DownloadTask)
  } else if (ev['open'] !== undefined) {
    const item = ev['open'] as DownloadTask;
    openTask(item.videoPath);
  } else if (ev['remove'] !== undefined) {
    deleteTask(ev['remove'] as DownloadTask);
  } else if (ev['stop'] !== undefined) {
    stopItem(ev['stop'] as DownloadTask);
  } else if (ev['retry'] !== undefined) {
    const item = ev['retry'] as DownloadTask;
    deleteTask(item);
    TaskAdd(item);
  }
}

</script>

<template>
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
    <el-table
        :data="filterTableData"
        :row-class-name="tableRowClassName"
        style="width: 100%">
      <el-table-column label="创建时间" prop="time"/>
      <el-table-column label="任务名" prop="taskName"/>
      <el-table-column label="状态" prop="status"/>
      <el-table-column align="right">
        <template #header>
          <el-input v-model="search" size="small" placeholder="搜索任务" />
        </template>
        <template #default="scope">
          <el-dropdown size="small" type="primary" @command="handleDownloadTask">
          <span class="el-dropdown-link">
            更多
            <el-icon class="el-icon--right">
              <ArrowDown />
            </el-icon>
          </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :icon="FolderOpened" :command="{open: scope.row}">打开</el-dropdown-item>
                <el-dropdown-item :icon="VideoPlay" :command="{play: scope.row}">播放</el-dropdown-item>
                <el-dropdown-item :icon="RemoveFilled" :command="{remove: scope.row}">删除</el-dropdown-item>
                <el-dropdown-item :icon="Stopwatch" :command="{stop: scope.row}">停止</el-dropdown-item>
                <el-dropdown-item :icon="Refresh" :command="{retry:scope.row}">重试</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
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
</template>

<style scoped>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-9);
}

.el-table .error-row {
  --el-table-tr-bg-color: var(--el-color-error-light-9);
}

.el-table .primary-row {
  --el-table-tr-bg-color: var(--el-color-primary-light-9);
}

.el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
</style>