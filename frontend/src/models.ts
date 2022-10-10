export class VideoItem {
    url: string = "";
    key: string = "";
    headers: string = "";
}

export class DownloadTask {
    id = "";
    url = "";
    taskName = "";
    time = "";
    status = "";
    videoPath = "";
}

export enum MergeFileType {
    Speed = "speed",
    TransCoding = "transcoding"
}

export class PlaylistItem {
    uri = "";
    desc = "";
}

export class InterceptItem {
    url = "";
    status = 0;
}