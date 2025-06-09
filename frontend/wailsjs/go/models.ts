export namespace main {
	
	export class BuildInfo {
	    version: string;
	    commit: string;
	    buildTime: string;
	
	    static createFrom(source: any = {}) {
	        return new BuildInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.commit = source["commit"];
	        this.buildTime = source["buildTime"];
	    }
	}
	export class ProxyStatus {
	    enabled: boolean;
	    server: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.server = source["server"];
	        this.error = source["error"];
	    }
	}

}

