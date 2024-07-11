import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import rnf from '$lib/reg/MRX.123.456.789.rnf/register.json';

export const fetchList = [
{type: "csv", id: "IET", url: "demo12/IET.csv" },
{type: "csv", id: "laundromat", url: "demo12/laundromat.csv" },
{type: "csv", id: "lostpast", url: "demo12/lostpast.csv" },
{type: "csv", id: "springwatch", url: "demo12/springwatch.csv" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.rnf", reg: rnf},
]

export const dataCache = [
{type: "csv", id:"IET", post: {inputMRXID:"MRX.123.456.789.rnf", data: "IET", outputMRXID:"generateFFmpeg" ,title : "IET"}},
{type: "csv", id:"laundromat", post: {inputMRXID:"MRX.123.456.789.rnf", data: "laundromat", outputMRXID:"generateFFmpeg" ,title : "LM"}},
{type: "csv", id:"lostpast", post: {inputMRXID:"MRX.123.456.789.rnf", data: "lostpast", outputMRXID:"generateFFmpeg" ,title : " LP"}},
{type: "csv", id:"springwatch", post: {inputMRXID:"MRX.123.456.789.rnf", data: "springwatch", outputMRXID:"generateFFmpeg" ,title : "SW"}},
]

export const sources: DemoSource[] = [
{
		id: "IET",
		mrxId: "MRX.123.456.789.rnf",
		name: "IET.csv",
		clip: "",
		srcUrl: "demo12/IET.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
{
		id: "laundromat",
		mrxId: "MRX.123.456.789.rnf",
		name: "laundromat.csv",
		clip: "",
		srcUrl: "demo12/laundromat.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
{
		id: "lostpast",
		mrxId: "MRX.123.456.789.rnf",
		name: "lostpast.csv",
		clip: "",
		srcUrl: "demo12/lostpast.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
{
		id: "springwatch",
		mrxId: "MRX.123.456.789.rnf",
		name: "springwatch.csv",
		clip: "",
		srcUrl: "demo12/springwatch.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
]
