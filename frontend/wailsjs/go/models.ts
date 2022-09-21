export namespace main {
	
	export class MergeFilesConfig {
	    files: string[];
	    age: string;
	    taskName: string;
	
	    static createFrom(source: any = {}) {
	        return new MergeFilesConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = source["files"];
	        this.age = source["age"];
	        this.taskName = source["taskName"];
	    }
	}

}

