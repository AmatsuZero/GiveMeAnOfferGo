<script lang="ts" setup>
import {reactive} from 'vue'
import {OpenSelectM3U8, OpenConfigDir, TaskAdd, StartMergeTs} from '../../wailsjs/go/main/App'
import {main} from "../../wailsjs/go/models";
import ParserTask = main.ParserTask;
import {EventsEmit, EventsOn} from "../../wailsjs/runtime"

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  const config = new main.MergeFilesConfig()
  StartMergeTs(config).then(ret => {
    data.resultText = ret
  }).catch(err => {

  })
}

function openConfigDir() {
  OpenConfigDir().then(dir => {
    data.resultText = dir
  }).catch(err => {
    console.error(err);
  });
}

function ParseAndDownload() {
  const task = new ParserTask();
  task.url = data.name;
  TaskAdd(task).then(() => {

  }).catch( err => {
    console.error(err);
  });
}

EventsOn("select-variant", (msg) => {
  const resolution = Object.keys(msg["Info"]);

  EventsEmit("variant-selected", resolution[0]);
})

</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="greet">Greet</button>
      <button class="btn" @click="openConfigDir">Select</button>
      <button class="btn" @click="ParseAndDownload">Download</button>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 74px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
