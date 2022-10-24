import {main} from "../wailsjs/go/models";
import ParserTask = main.ParserTask;

export class VideoItem {
    url: string = "";
    key: string = "";
    headers: string = "";
}

export enum DownloadTaskState {
    Done = "finish",
    Error = "error",
    Processing = "processing",
    Idle = "idle"
}

export class DownloadTask extends ParserTask {
    time = "";
    status = "";
    videoPath = "";
    isDone = false;
    state = DownloadTaskState.Idle;
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