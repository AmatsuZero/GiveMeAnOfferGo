export namespace main {
	export class MergeFilesConfig {
	    files: string[];
	    mergeType: string;
	    taskName: string;
		output: string;
	
	    static createFrom(source: any = {}) {
	        return new MergeFilesConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = source["files"] ?? [];
	        this.mergeType = source["mergeType"];
	        this.taskName = source["taskName"];
			this.output = source["output"];
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

