export enum MergeFileType {
    Speed = "speed",
    TransCoding = "transcoding"
}

export default {
    props: {
        version:'',
        m3u8_url: '',
        m3u8_urls: '',
        ts_dir:'',
        ts_urls:[],
        m3u8_url_prefix:'',
        dlg_header_visible: false,
        dlg_newtask_visible: false,
        config_save_dir:'',
        config_ffmpeg:'',
        config_proxy:'',
        headers:'',
        myKeyIV:'',
        myLocalKeyIV:'',
        taskName:'',
        taskIsDelTs:true,
        allVideos:[],
        tabPane:'',
        tsMergeType: MergeFileType.Speed,
        tsMergeProgress:0,
        tsMergeStatus:'',
        tsMergeMp4Path:'',
        tsMergeMp4Dir:'',
        tsTaskName:'',
        downloadSpeed:'0 MB/s',
        playlists:[],
        playlistUri:'',
        addTaskMessage:'',
        navigatorInput:'',
        //navigatorUrl:'about:blank',
        navigatorUrl:'https://haokan.baidu.com/?sfrom=baidu-top',
        currentUserAgent:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
        browserVideoUrls:[],
        platform:''
    },

    created() {

    }
}