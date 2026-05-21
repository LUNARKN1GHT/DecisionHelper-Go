export namespace main {
	
	export class Criterion {
	    id: string;
	    name: string;
	    weight: number;
	
	    static createFrom(source: any = {}) {
	        return new Criterion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.weight = source["weight"];
	    }
	}
	export class Score {
	    option: string;
	    criterion_id: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new Score(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.option = source["option"];
	        this.criterion_id = source["criterion_id"];
	        this.value = source["value"];
	    }
	}
	export class Decision {
	    id: string;
	    title: string;
	    created_at: string;
	    options: string[];
	    criteria: Criterion[];
	    scores: Score[];
	
	    static createFrom(source: any = {}) {
	        return new Decision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.created_at = source["created_at"];
	        this.options = source["options"];
	        this.criteria = this.convertValues(source["criteria"], Criterion);
	        this.scores = this.convertValues(source["scores"], Score);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class OptionResult {
	    option: string;
	    score: number;
	    rank: number;
	
	    static createFrom(source: any = {}) {
	        return new OptionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.option = source["option"];
	        this.score = source["score"];
	        this.rank = source["rank"];
	    }
	}

}

