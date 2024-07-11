import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';

import ghj from '$lib/reg/MRX.123.456.789.ghj/register.json';

export const fetchList = [
{type: "json", id: "beach_camera", url: "demo06/beach_camera.json" },
{type: "json", id: "fire_camera", url: "demo06/fire_camera.json" },
{type: "json", id: "fish_camera", url: "demo06/fish_camera.json" },
{type: "json", id: "roar_camera", url: "demo06/roar_camera.json" },
]

export const regCache:MrxRegisterCache = [
{mrxId: "MRX.123.456.789.ghj", reg: ghj},
]

export const dataCache = [
{type: "json", id:"beach_camera", post: {inputMRXID:"MRX.123.456.789.ghj", data: "beach_camera", outputMRXID:"MRX.123.456.789.cxm" }},
{type: "json", id:"fire_camera", post: {inputMRXID:"MRX.123.456.789.ghj", data: "fire_camera", outputMRXID:"MRX.123.456.789.cxm" }},
{type: "json", id:"fish_camera", post: {inputMRXID:"MRX.123.456.789.ghj", data: "fish_camera", outputMRXID:"MRX.123.456.789.cxm" }},
{type: "json", id:"roar_camera", post: {inputMRXID:"MRX.123.456.789.ghj", data: "roar_camera", outputMRXID:"MRX.123.456.789.cxm" }},
]

export const sources: DemoSource[] = [
{
		id: "beach_camera",
		mrxId: "MRX.123.456.789.ghj",
		name: "beach_camera.json",
		clip: "",
		srcUrl: "demo06/beach_camera.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "fire_camera",
		mrxId: "MRX.123.456.789.ghj",
		name: "fire_camera.json",
		clip: "",
		srcUrl: "demo06/fire_camera.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "fish_camera",
		mrxId: "MRX.123.456.789.ghj",
		name: "fish_camera.json",
		clip: "",
		srcUrl: "demo06/fish_camera.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
{
		id: "roar_camera",
		mrxId: "MRX.123.456.789.ghj",
		name: "roar_camera.json",
		clip: "",
		srcUrl: "demo06/roar_camera.json",
		summary: "",
		srcReg: "json",
		xtra: [],
	  },
]
