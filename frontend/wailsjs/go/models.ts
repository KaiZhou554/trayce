export namespace main {
	
	export class DeleteResult {
	    deleted: number;
	    backupPath: string;
	
	    static createFrom(source: any = {}) {
	        return new DeleteResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deleted = source["deleted"];
	        this.backupPath = source["backupPath"];
	    }
	}

}

export namespace trayicons {
	
	export class TrayIconEntry {
	    id: string;
	    iconGuid: string;
	    publisher: string;
	    executablePath: string;
	    iconBase64: string;
	    status: string;
	    isSpecialPath: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TrayIconEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.iconGuid = source["iconGuid"];
	        this.publisher = source["publisher"];
	        this.executablePath = source["executablePath"];
	        this.iconBase64 = source["iconBase64"];
	        this.status = source["status"];
	        this.isSpecialPath = source["isSpecialPath"];
	    }
	}

}

