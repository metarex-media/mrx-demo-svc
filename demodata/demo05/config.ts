import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import vqc from '$lib/reg/MRX.123.456.789.vqc/register.json';

export const fetchList = [
{type: "xml", id: "PulsarReport1Fail", url: "demo05/PulsarReport1Fail.xml" },
{type: "xml", id: "PulsarReport2Fail", url: "demo05/PulsarReport2Fail.xml" },
{type: "xml", id: "PulsarReport3Fail", url: "demo05/PulsarReport3Fail.xml" },
{type: "xml", id: "PulsarReport4Pass", url: "demo05/PulsarReport4Pass.xml" },
{type: "xml", id: "PulsarReport5Pass", url: "demo05/PulsarReport5Pass.xml" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.vqc", reg: vqc},
]

export const dataCache = [
{type: "xml", id:"PulsarReport1Fail", post: {inputMRXID:"MRX.123.456.789.vqc", data: "PulsarReport1Fail", outputMRXID:"tograph" }},
{type: "xml", id:"PulsarReport2Fail", post: {inputMRXID:"MRX.123.456.789.vqc", data: "PulsarReport2Fail", outputMRXID:"tograph" }},
{type: "xml", id:"PulsarReport3Fail", post: {inputMRXID:"MRX.123.456.789.vqc", data: "PulsarReport3Fail", outputMRXID:"tograph" }},
{type: "xml", id:"PulsarReport4Pass", post: {inputMRXID:"MRX.123.456.789.vqc", data: "PulsarReport4Pass", outputMRXID:"tograph" }},
{type: "xml", id:"PulsarReport5Pass", post: {inputMRXID:"MRX.123.456.789.vqc", data: "PulsarReport5Pass", outputMRXID:"tograph" }},
]

export const sources: DemoSource[] = [
{
		id: "PulsarReport1Fail",
		mrxId: "MRX.123.456.789.vqc",
		name: "PulsarReport1Fail.xml",
		clip: "",
		srcUrl: "demo05/PulsarReport1Fail.xml",
		summary: "",
		srcReg: "xml",
		xtra: [],
	  },
{
		id: "PulsarReport2Fail",
		mrxId: "MRX.123.456.789.vqc",
		name: "PulsarReport2Fail.xml",
		clip: "",
		srcUrl: "demo05/PulsarReport2Fail.xml",
		summary: "",
		srcReg: "xml",
		xtra: [],
	  },
{
		id: "PulsarReport3Fail",
		mrxId: "MRX.123.456.789.vqc",
		name: "PulsarReport3Fail.xml",
		clip: "",
		srcUrl: "demo05/PulsarReport3Fail.xml",
		summary: "",
		srcReg: "xml",
		xtra: [],
	  },
{
		id: "PulsarReport4Pass",
		mrxId: "MRX.123.456.789.vqc",
		name: "PulsarReport4Pass.xml",
		clip: "",
		srcUrl: "demo05/PulsarReport4Pass.xml",
		summary: "",
		srcReg: "xml",
		xtra: [],
	  },
{
		id: "PulsarReport5Pass",
		mrxId: "MRX.123.456.789.vqc",
		name: "PulsarReport5Pass.xml",
		clip: "",
		srcUrl: "demo05/PulsarReport5Pass.xml",
		summary: "",
		srcReg: "xml",
		xtra: [],
	  },
]
