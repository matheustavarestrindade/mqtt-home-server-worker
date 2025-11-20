<script lang="ts">
	import type { SensorDataResponse } from '$lib/services/SensorsService';
	import MicroChart from './MicroChart.svelte';

	const { data }: { data: SensorDataResponse<'hidroponic-manager'> } = $props();

	const dataArrays = $derived.by(() => {
		const moistureData: number[] = [];
		const moistureSeverityData: number[] = [];

		const phData: number[] = [];
		const phSeverityData: number[] = [];

		const nData: number[] = [];
		const nSeverityData: number[] = [];

		const pData: number[] = [];
		const pSeverityData: number[] = [];

		const kData: number[] = [];
		const kSeverityData: number[] = [];

		const ecData: number[] = [];
		const ecSeverityData: number[] = [];

		const relayOnData: number[] = [];

		data.forEach((entry) => {
			moistureData.push(entry.moisture);
			moistureSeverityData.push(entry.moisture_severity);
			phData.push(entry.ph);
			phSeverityData.push(entry.ph_severity);
			nData.push(entry.nitrogen);
			nSeverityData.push(entry.nitrogen_severity);
			pData.push(entry.phosphorus);
			pSeverityData.push(entry.phosphorus_severity);
			kData.push(entry.potassium);
			kSeverityData.push(entry.potassium_severity);
			ecData.push(entry.conductivity);
			ecSeverityData.push(entry.conductivity_severity);
			relayOnData.push(entry.isOn ? 1 : 0);
		});

		const lastPHValue = phData[phData.length - 1] == undefined ? 'Unk.' : phData[phData.length - 1].toFixed(2);
		const lastNValue = nData[nData.length - 1] == undefined ? 'Unk.' : nData[nData.length - 1].toFixed(0);
		const lastPValue = pData[pData.length - 1] == undefined ? 'Unk.' : pData[pData.length - 1].toFixed(0);
		const lastKValue = kData[kData.length - 1] == undefined ? 'Unk.' : kData[kData.length - 1].toFixed(0);
		const lastECValue = ecData[ecData.length - 1] == undefined ? 'Unk.' : ecData[ecData.length - 1].toFixed(0);

		return [
			{
				name: 'pH ' + lastPHValue,
				data: phData,
				severityData: phSeverityData,
				color: '#10b981'
			},
			{
				name: 'Nitrogen ' + lastNValue,
				data: nData,
				severityData: nSeverityData,
				color: '#f59e0b'
			},
			{
				name: 'Phosphorus ' + lastPValue,
				data: pData,
				severityData: pSeverityData,
				color: '#8b5cf6'
			},
			{
				name: 'Potassium ' + lastKValue,
				data: kData,
				severityData: kSeverityData,
				color: '#ef4444'
			},
			{
				name: 'EC ' + lastECValue,
				data: ecData,
				severityData: ecSeverityData,
				color: '#14b8a6'
			}
		];
	});
</script>

<div class="mb-4 flex flex-wrap gap-4">
	{#each dataArrays as item}
		<MicroChart name={item.name} data={item.data} color={item.color} />
	{/each}
</div>
