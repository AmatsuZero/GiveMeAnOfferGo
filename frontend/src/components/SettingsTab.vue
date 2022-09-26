<script lang="ts">
export default {
  name: "SettingsTab"
}
</script>

<template>
  <el-tab-pane label="软件设置">
    <span slot="label"><i class="el-icon-s-tools"></i> 软件设置</span>
    <el-row :gutter="8">
      <el-col :span="2" :offset="1">
        <span style="line-height: 40px;float: right;">视频存储</span>
      </el-col>
      <el-col :span="18">
        <el-input placeholder="请单击右侧按钮选择存储文件夹" v-model="config_save_dir" :readonly="true">
        </el-input>
      </el-col>
      <el-col :span="2">
        <el-button class="mybutton" type="primary" icon="el-icon-folder-opened"
                   @click="clickOpenConfigDir">浏览</el-button>
      </el-col>
    </el-row>
    <br />
    <el-row :gutter="8">
      <el-col :span="2" :offset="1">
        <span style="line-height: 40px;float: right;">日志</span>
      </el-col>
      <el-col :span="4">
        <el-button class="mybutton" type="primary" icon="el-icon-view"
                   @click="clickOpenLogDir">查看日志目录</el-button>
      </el-col>
      <el-col :span="4" :offset="1">
        <span style="line-height: 40px;float: right;">HTTP代理</span>
      </el-col>

      <el-col :span="9">
        <el-input placeholder="请输入HTTP代理(如：http://127.0.0.1:7890 )" v-model="config_proxy"
                  @change="proxyChange"></el-input>
      </el-col>
    </el-row>
    <br>
    <el-divider><i class="el-icon-coffee"></i></el-divider>
    <el-row :gutter="8">
      <el-col :span="4" :offset="7">
        <el-card class="paycode alipay" shadow="always">
          <img class="image" src="icon/alipay_code.png">
          <div style="padding: 14px;word-break: keep-all;">
            <span>支付宝赞助</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4" :offset="2">
        <el-card class="paycode wechat" shadow="always">
          <img class="image" src="icon/wechat_code.png">
          <div style="padding: 14px;word-break: keep-all;">
            <span>微信赞助</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="8">
      <el-col :span="12" :offset="6">
        <el-alert title="开发不易，如果您觉得此款软件对您有帮助，可以请小哥哥喝杯咖啡吗?" type="error" :center="true"
                  :closable="false">
        </el-alert>
      </el-col>
    </el-row>
  </el-tab-pane>
</template>

<style scoped>

</style>