import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import rny from '$lib/reg/MRX.123.456.789.rny/register.json';
import rnc from '$lib/reg/MRX.123.456.789.rnc/register.json';
import rnj from '$lib/reg/MRX.123.456.789.rnj/register.json';

export const fetchList = [
{type: "yaml", id: "IET", url: "demo04/IET.yaml" },
{type: "json", id: "cosmoslaundromat", url: "demo04/cosmos-laundromat.json" },
{type: "csv", id: "lostpast", url: "demo04/lostpast.csv" },
{type: "csv", id: "springwatch", url: "demo04/springwatch.csv" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.rny", reg: rny},
{mrxId: "MRX.123.456.789.rnc", reg: rnc},
{mrxId: "MRX.123.456.789.rnj", reg: rnj},
]

export const dataCache = [
{type: "yaml", id:"IET", post: {inputMRXID:"MRX.123.456.789.rny", data: "IET", outputMRXID:"MRX.123.456.789.rnf" ,mapping : "true"}},
{type: "json", id:"cosmoslaundromat", post: {inputMRXID:"MRX.123.456.789.rnj", data: "cosmoslaundromat", outputMRXID:"MRX.123.456.789.rnf" ,mapping : "true"}},
{type: "csv", id:"lostpast", post: {inputMRXID:"MRX.123.456.789.rnc", data: "lostpast", outputMRXID:"MRX.123.456.789.rnf" ,mapping : "true"}},
{type: "csv", id:"springwatch", post: {inputMRXID:"MRX.123.456.789.rnc", data: "springwatch", outputMRXID:"MRX.123.456.789.rnf" ,mapping : "true"}},
]

export const sources: DemoSource[] = [
{
		id: "IET",
		mrxId: "MRX.123.456.789.rny",
		name: "IET.yaml",
		clip: "",
		srcUrl: "demo04/IET.yaml",
		summary: "",
		srcReg: "yaml",
		xtra: [],
	  },
{
		id: "cosmoslaundromat",
		mrxId: "MRX.123.456.789.rnj",
		name: "cosmos-laundromat.json",
		clip: "",
		srcUrl: "demo04/cosmos-laundromat.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "lostpast",
		mrxId: "MRX.123.456.789.rnc",
		name: "lostpast.csv",
		clip: "",
		srcUrl: "demo04/lostpast.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
{
		id: "springwatch",
		mrxId: "MRX.123.456.789.rnc",
		name: "springwatch.csv",
		clip: "",
		srcUrl: "demo04/springwatch.csv",
		summary: "",
		srcReg: "csv",
		xtra: [],
	  },
]
