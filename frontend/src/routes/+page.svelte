<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Activity, TriangleAlert, CircleCheck, Bell } from 'lucide-svelte';
	import SensorCard from './components/SensorCard.svelte';
	import SensorItem from './components/SensorItem.svelte';
	import { getUserSensorData, getUserSensors } from '$lib/remote/user.remote';
	import HydroponicMangerCharts from './components/HydroponicMangerCharts.svelte';
	import WaterLevelCharts from './components/WaterLevelCharts.svelte';

	const userSensors = await getUserSensors();

	interface Sensor {
		id: string;
		name: string;
		type: string;
		location: string;
		status: 'online' | 'offline' | 'warning';
		lastMessage: string;
		battery: number;
		signal: number;
		notifications: string[];
	}

	$inspect(userSensors);

	const HOURS_TO_FETCH = 2;

	const sensors: Sensor[] = userSensors.map((sensor) => {
		const lastSeen = new Date(sensor.last_seen).getTime();
		const now = Date.now();
		const timeDiff = now - lastSeen;

		return {
			id: sensor.fuse_id,
			name: sensor.name,
			type: sensor.type,
			location: sensor.location,
			status: timeDiff < 10 * 60 * 1000 ? 'online' : timeDiff < 30 * 60 * 1000 ? 'warning' : 'offline',
			lastMessage: sensor.last_seen.toLocaleString(),
			battery: sensor.battery_percent,
			signal: sensor.wifi_strength,
			notifications: []
		};
	});
</script>

<div class="flex-1 space-y-6 p-8 pt-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Sensor Dashboard</h1>
			<p class="mt-2 text-muted-foreground">Monitor your connected sensors and their real-time status</p>
		</div>
		<Button>
			<Activity class="mr-2 h-4 w-4" />
			Refresh All
		</Button>
	</div>

	<!-- Stats Overview -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
		<SensorCard name="Total Sensors" value={sensors.length} Icon={Activity} description="Active devices" />
		<SensorCard
			name="Online Sensors"
			value={sensors.filter((s) => s.status === 'online').length}
			Icon={CircleCheck}
			iconColor="green"
			description={`${Math.round((sensors.filter((s) => s.status === 'online').length / sensors.length) * 100)}% uptime`}
		/>
		<SensorCard
			name="Active Alerts"
			value={sensors.filter((s) => s.notifications.length > 0).length}
			Icon={TriangleAlert}
			iconColor="yellow"
			description="Requires attention"
		/>
		<SensorCard name="Data Points Today" value="1,247" Icon={Bell} description="+12% from yesterday" />
	</div>

	<!-- Sensors Grid -->
	<div class="flex flex-wrap gap-6">
		{#each sensors as sensor}
			<SensorItem {sensor}>
				{#if sensor.type === 'hydroponic-manager'}
					{#await getUserSensorData({ fuseId: sensor.id, from: new Date(Date.now() - HOURS_TO_FETCH * 60 * 60 * 1000), to: new Date() }) then sensorData}
						{#if 'error' in sensorData}
							<p class="text-sm text-red-500">Error loading data: {sensorData.error}</p>
						{:else}
							<HydroponicMangerCharts data={sensorData} />
						{/if}
					{/await}
				{:else if sensor.type === 'water-level-meter'}
					{#await getUserSensorData({ fuseId: sensor.id, from: new Date(Date.now() - HOURS_TO_FETCH * 60 * 60 * 1000), to: new Date() }) then sensorData}
						{#if 'error' in sensorData}
							<p class="text-sm text-red-500">Error loading data: {sensorData.error}</p>
						{:else}
							<WaterLevelCharts data={sensorData} />
						{/if}
					{/await}
				{:else}
					<p class="text-sm text-muted-foreground">No data available for this sensor type.</p>
				{/if}
			</SensorItem>
		{/each}
	</div>
</div>
