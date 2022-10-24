// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {DownloadTask} from "../../../src/models";

export function Open(arg1:string):Promise<Error>;

export function OpenConfigDir():Promise<string>;

export function OpenSelectM3U8():Promise<string>;

export function OpenSelectTsDir(arg1:string):Promise<Array<string>>;

export function Play(arg1:string):Promise<Error>;

export function SniffLinks(arg1:string):Promise<Array<string>>;

export function StartMergeTs(arg1:main.MergeFilesConfig):Promise<Error>;

export function TaskAdd(arg1:main.ParserTask):Promise<Error>;

export function TaskAddMuti(arg1:Array<main.ParserTask>):Promise<Error>;

export function RemoveTaskNotifyItem(arg: DownloadTask): Promise<void>;

export function ClearTasks(): void;