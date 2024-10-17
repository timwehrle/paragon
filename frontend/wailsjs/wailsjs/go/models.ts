export namespace main {
	
	export class VolumeInfo {
	    drive: string;
	    type: string;
	    volumeLabel: string;
	    fileSystem: string;
	    serialNumber: string;
	
	    static createFrom(source: any = {}) {
	        return new VolumeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.drive = source["drive"];
	        this.type = source["type"];
	        this.volumeLabel = source["volumeLabel"];
	        this.fileSystem = source["fileSystem"];
	        this.serialNumber = source["serialNumber"];
	    }
	}

}

