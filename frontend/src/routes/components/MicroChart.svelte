<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import type ApexCharts from 'apexcharts';
	import { browser } from '$app/environment';

	let chartElement: HTMLDivElement;
	let chart: ApexCharts | undefined = $state();

	type Props = {
		name: string;
		data: number[];
		severityData?: number[];
		color: string;
		severityColor?: string;
	};

	const { name, data, color, severityData, severityColor }: Props = $props();

	const options = {
		series: [
			{
				name,
				data
			},
			{
				name: 'Severity',
				data: severityData || [],
				color: severityColor || '#FF4560'
			}
		],
		chart: {
			type: 'area',
			sparkline: {
				enabled: true
			},
			height: 80,
			width: '100%'
		},
		stroke: { curve: 'smooth', width: 2 },
		fill: { opacity: 0.3, colors: [color] },
		colors: [color]
	};

	const mountApexChart = async () => {
		if (!browser) return;
		const ApexChartsModule = (await import('apexcharts')).default;
		chart = new ApexChartsModule(chartElement, options);
		chart.render();
	};

	onMount(() => {
		mountApexChart();
	});

	onDestroy(() => {
		if (chart) {
			chart.destroy();
		}
	});
</script>

<div class="flex min-w-[120px] flex-1 flex-col">
	<div bind:this={chartElement}></div>
	<span class="mt-1 text-sm font-medium text-gray-500">{name}</span>
</div>
