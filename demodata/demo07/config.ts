import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import bat from '$lib/reg/MRX.123.456.789.bat/register.json';

export const fetchList = [
{type: "json", id: "fire", url: "demo07/fire.json" },
{type: "json", id: "good", url: "demo07/good.json" },
{type: "json", id: "puncture", url: "demo07/puncture.json" },
{type: "json", id: "red", url: "demo07/red.json" },
{type: "json", id: "unknown", url: "demo07/unknown.json" },
{type: "json", id: "yellow", url: "demo07/yellow.json" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.bat", reg: bat},
]

export const dataCache = [
{type: "json", id:"fire", post: {inputMRXID:"MRX.123.456.789.bat", data: "fire", outputMRXID:"fill" }},
{type: "json", id:"fire", post: {inputMRXID:"MRX.123.456.789.bat", data: "fire", outputMRXID:"fault" }},
{type: "json", id:"fire", post: {inputMRXID:"MRX.123.456.789.bat", data: "fire", outputMRXID:"stagger" }},
{type: "json", id:"good", post: {inputMRXID:"MRX.123.456.789.bat", data: "good", outputMRXID:"fill" }},
{type: "json", id:"good", post: {inputMRXID:"MRX.123.456.789.bat", data: "good", outputMRXID:"fault" }},
{type: "json", id:"good", post: {inputMRXID:"MRX.123.456.789.bat", data: "good", outputMRXID:"stagger" }},
{type: "json", id:"puncture", post: {inputMRXID:"MRX.123.456.789.bat", data: "puncture", outputMRXID:"fill" }},
{type: "json", id:"puncture", post: {inputMRXID:"MRX.123.456.789.bat", data: "puncture", outputMRXID:"fault" }},
{type: "json", id:"puncture", post: {inputMRXID:"MRX.123.456.789.bat", data: "puncture", outputMRXID:"stagger" }},
{type: "json", id:"red", post: {inputMRXID:"MRX.123.456.789.bat", data: "red", outputMRXID:"fill" }},
{type: "json", id:"red", post: {inputMRXID:"MRX.123.456.789.bat", data: "red", outputMRXID:"fault" }},
{type: "json", id:"red", post: {inputMRXID:"MRX.123.456.789.bat", data: "red", outputMRXID:"stagger" }},
{type: "json", id:"unknown", post: {inputMRXID:"MRX.123.456.789.bat", data: "unknown", outputMRXID:"fill" }},
{type: "json", id:"unknown", post: {inputMRXID:"MRX.123.456.789.bat", data: "unknown", outputMRXID:"fault" }},
{type: "json", id:"unknown", post: {inputMRXID:"MRX.123.456.789.bat", data: "unknown", outputMRXID:"stagger" }},
{type: "json", id:"yellow", post: {inputMRXID:"MRX.123.456.789.bat", data: "yellow", outputMRXID:"fill" }},
{type: "json", id:"yellow", post: {inputMRXID:"MRX.123.456.789.bat", data: "yellow", outputMRXID:"fault" }},
{type: "json", id:"yellow", post: {inputMRXID:"MRX.123.456.789.bat", data: "yellow", outputMRXID:"stagger" }},
]

export const sources: DemoSource[] = [
{
		id: "fire",
		mrxId: "MRX.123.456.789.bat",
		name: "fire.json",
		clip: "",
		srcUrl: "demo07/fire.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "good",
		mrxId: "MRX.123.456.789.bat",
		name: "good.json",
		clip: "",
		srcUrl: "demo07/good.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "puncture",
		mrxId: "MRX.123.456.789.bat",
		name: "puncture.json",
		clip: "",
		srcUrl: "demo07/puncture.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "red",
		mrxId: "MRX.123.456.789.bat",
		name: "red.json",
		clip: "",
		srcUrl: "demo07/red.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "unknown",
		mrxId: "MRX.123.456.789.bat",
		name: "unknown.json",
		clip: "",
		srcUrl: "demo07/unknown.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "yellow",
		mrxId: "MRX.123.456.789.bat",
		name: "yellow.json",
		clip: "",
		srcUrl: "demo07/yellow.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
]
