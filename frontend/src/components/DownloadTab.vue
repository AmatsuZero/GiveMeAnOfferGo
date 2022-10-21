<script lang="ts">
export default {
  props: {
    tabPane: "",
    config_save_dir: ""
  },

  name: "download-tab",
}
</script>
<script lang="ts" setup>
import List from './Tasklist.vue';
import Multi from './MultiTasks.vue';
import Merge from  './MergeTask.vue';
</script>
<template>
  <el-tab-pane label="资源下载">
    <span slot="label">资源下载</span>
    <el-tabs type="border-card" v-model="tabPane">
      <List/>
      <Multi/>
      <Merge/>
    </el-tabs>
  </el-tab-pane>
</template>

<style scoped>

</style>