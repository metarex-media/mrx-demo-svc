import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import mxf from '$lib/reg/MRX.123.456.789.mxf/register.json';

export const fetchList = [
{type: "mxf", id: "badISXD", url: "demo13/badISXD.mxf" },
{type: "mxf", id: "goodISXD", url: "demo13/goodISXD.mxf" },
{type: "mxf", id: "veryBadISXD", url: "demo13/veryBadISXD.mxf" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.mxf", reg: mxf},
]

export const dataCache = [
{type: "mxf", id:"badISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "badISXD", outputMRXID:"MXFToGraph" }},
{type: "mxf", id:"badISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "badISXD", outputMRXID:"MXFToReport" }},
{type: "mxf", id:"goodISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "goodISXD", outputMRXID:"MXFToGraph" }},
{type: "mxf", id:"goodISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "goodISXD", outputMRXID:"MXFToReport" }},
{type: "mxf", id:"veryBadISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "veryBadISXD", outputMRXID:"MXFToGraph" }},
{type: "mxf", id:"veryBadISXD", post: {inputMRXID:"MRX.123.456.789.mxf", data: "veryBadISXD", outputMRXID:"MXFToReport" }},
]

export const sources: DemoSource[] = [
{
		id: "badISXD",
		mrxId: "MRX.123.456.789.mxf",
		name: "badISXD.mxf",
		clip: "",
		srcUrl: "demo13/badISXD.mxf",
		summary: "",
		srcReg: "mxf",
		xtra: [],
	  },
{
		id: "goodISXD",
		mrxId: "MRX.123.456.789.mxf",
		name: "goodISXD.mxf",
		clip: "",
		srcUrl: "demo13/goodISXD.mxf",
		summary: "",
		srcReg: "mxf",
		xtra: [],
	  },
{
		id: "veryBadISXD",
		mrxId: "MRX.123.456.789.mxf",
		name: "veryBadISXD.mxf",
		clip: "",
		srcUrl: "demo13/veryBadISXD.mxf",
		summary: "",
		srcReg: "mxf",
		xtra: [],
	  },
]
