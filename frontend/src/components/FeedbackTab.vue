<template>
  <el-tab-pane label="交流反馈">
    <span slot="label"><i class="el-icon-chat-dot-square"></i> 交流反馈</span>
    <el-row :gutter="8">
      <el-col :span="22">
        <el-alert title="警示：【天网恢恢疏而不漏】遵纪守法是每位公民应尽的义务，严禁在此发布违规违法内容，否则后果自负！违规内容举报QQ：1586462"
                  type="warning" :center="true" :closable="false" show-icon>
        </el-alert>
      </el-col>
      <el-col :span="2">
        <el-button class="mybutton" type="success" @click="clickRefreshComment">刷新
        </el-button>
<!--        </el-alert>-->
      </el-col>
    </el-row>
    <section id="comments" class="vcomment"></section>
  </el-tab-pane>
</template>

<script>
export default {
  name: "FeedbackTab"
}
</script>

<style scoped>

</style>