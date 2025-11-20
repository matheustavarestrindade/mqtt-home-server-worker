<script lang="ts">
	import type { SensorDataResponse } from '$lib/services/SensorsService';
	import MicroChart from './MicroChart.svelte';

	const { data }: { data: SensorDataResponse<'water-level-meter'> } = $props();

	const dataArrays = $derived.by(() => {
		const waterLevelDatas: number[] = [];

		data.forEach((entry) => waterLevelDatas.push(entry.average_water_level_cm));
		const lastWaterLevelValue = waterLevelDatas[waterLevelDatas.length - 1] == undefined ? 'Unk.' : waterLevelDatas[waterLevelDatas.length - 1].toFixed(2);

		return [
			{
				name: 'Water Level ' + lastWaterLevelValue + ' cm',
				data: waterLevelDatas,
				color: '#3b82f6'
			}
		];
	});
</script>

<div class="mb-4 flex flex-wrap gap-4">
	{#each dataArrays as item}
		<MicroChart name={item.name} data={item.data} color={item.color} />
	{/each}
</div>
