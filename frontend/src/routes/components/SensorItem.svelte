<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardTitle, CardHeader, CardDescription, CardContent } from '$lib/components/ui/card';
	import { Activity, Battery, Clock, MapPin, TriangleAlert, Wifi, WifiOff } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

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

	const { sensor, children }: { sensor: Sensor; children?: Snippet<[]> } = $props();
	const SignalIcon = sensor.signal > 70 ? Wifi : sensor.signal > 30 ? Wifi : WifiOff;

	function getStatusColor(status: string) {
		switch (status) {
			case 'online':
				return 'bg-green-500';
			case 'offline':
				return 'bg-red-500';
			case 'warning':
				return 'bg-yellow-500';
			default:
				return 'bg-gray-500';
		}
	}
</script>

<Card class="min-width-[2O0px] flex-1 relative overflow-hidden">
	<!-- Status indicator bar -->
	<div class={`absolute top-0 left-0 h-full w-1 ${getStatusColor(sensor.status)}`}></div>
	<CardHeader class="pb-3">
		<div class="flex items-start justify-between">
			<div>
				<CardTitle class="text-lg">{sensor.name}</CardTitle>
				<CardDescription class="mt-1">{sensor.id}</CardDescription>
			</div>
			<Badge variant={sensor.status === 'online' ? 'default' : sensor.status === 'warning' ? 'secondary' : 'destructive'}>
				{sensor.status}
			</Badge>
		</div>
	</CardHeader>
	<CardContent class="flex flex-1 flex-col space-y-4">
		<!-- Sensor Info -->
		<div class="space-y-2">
			<div class="flex items-center text-sm text-muted-foreground">
				<MapPin class="mr-2 h-4 w-4" />
				{sensor.location}
			</div>
			<div class="flex items-center text-sm text-muted-foreground">
				<Activity class="mr-2 h-4 w-4" />
				{sensor.type}
			</div>
		</div>
		<!-- Status Indicators -->
		<div class="grid grid-cols-2 gap-2 border-t pt-2">
			<div class="flex items-center space-x-2">
				<Battery class={`h-4 w-4 ${sensor.battery < 25 ? 'text-red-500' : sensor.battery < 50 ? 'text-yellow-500' : 'text-green-500'}`} />
				<span class="text-sm font-medium">{sensor.battery}%</span>
			</div>
			<div class="flex items-center space-x-2">
				<SignalIcon class={`h-4 w-4 ${sensor.signal < 30 ? 'text-red-500' : sensor.signal < 70 ? 'text-yellow-500' : 'text-green-500'}`} />
				<span class="text-sm font-medium">{sensor.signal}%</span>
			</div>
		</div>
		<!-- Last Message -->
		<div class="flex items-center text-xs text-muted-foreground">
			<Clock class="mr-2 h-3 w-3" />
			Last message: {sensor.lastMessage}
		</div>
		<!-- Notifications -->
		{#if sensor.notifications.length > 0}
			<div class="space-y-1 rounded-lg bg-yellow-50 p-3 dark:bg-yellow-950">
				<div class="flex items-center">
					<TriangleAlert class="mr-2 h-4 w-4 text-yellow-600 dark:text-yellow-400" />
					<span class="text-xs font-medium text-yellow-800 dark:text-yellow-200">Notifications</span>
				</div>
				{#each sensor.notifications as notification}
					<p class="ml-6 text-xs text-yellow-700 dark:text-yellow-300">{notification}</p>
				{/each}
			</div>
		{/if}
		{@render children?.()}
		<!-- Actions -->
		<div class="mt-auto flex space-x-2 pt-2">
			<Button variant="outline" size="sm" class="flex-1">View Details</Button>
			<Button variant="ghost" size="sm">Configure</Button>
		</div>
	</CardContent>
</Card>
