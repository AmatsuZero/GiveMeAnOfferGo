export namespace main {
	
	export class DownloadTaskUIItem {
	    url: string;
	    taskName: string;
	    prefix: string;
	    delOnComplete: boolean;
	    keyIV: string;
	    headers: string;
	    // Go type: time.Time
	    time: any;
	    status: string;
	    isDone: boolean;
	    videoPath: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadTaskUIItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.taskName = source["taskName"];
	        this.prefix = source["prefix"];
	        this.delOnComplete = source["delOnComplete"];
	        this.keyIV = source["keyIV"];
	        this.headers = source["headers"];
	        this.time = this.convertValues(source["time"], null);
	        this.status = source["status"];
	        this.isDone = source["isDone"];
	        this.videoPath = source["videoPath"];
	        this.state = source["state"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MergeFilesConfig {
	    files: string[];
	    mergeType: string;
	    taskName: string;
	
	    static createFrom(source: any = {}) {
	        return new MergeFilesConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = source["files"];
	        this.mergeType = source["mergeType"];
	        this.taskName = source["taskName"];
	    }
	}
	export class ParserTask {
	    url: string;
	    taskName: string;
	    prefix: string;
	    delOnComplete: boolean;
	    keyIV: string;
	    headers: string;
	
	    static createFrom(source: any = {}) {
	        return new ParserTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.taskName = source["taskName"];
	        this.prefix = source["prefix"];
	        this.delOnComplete = source["delOnComplete"];
	        this.keyIV = source["keyIV"];
	        this.headers = source["headers"];
	    }
	}

}

