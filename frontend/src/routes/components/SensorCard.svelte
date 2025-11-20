<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import type { Icon as IconLucide } from 'lucide-svelte';

	type SensorCardProps = {
		name: string;
		value: string | number;
		Icon: typeof IconLucide;
		iconColor?: 'green' | 'yellow' | 'red';
		description?: string;
	};

	const { name, value, Icon, iconColor, description }: SensorCardProps = $props();

	const colorClassName = iconColor
		? iconColor === 'green'
			? 'text-green-500'
			: iconColor === 'yellow'
				? 'text-yellow-500'
				: 'text-red-500'
		: 'text-muted-foreground';
</script>

<Card>
	<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
		<CardTitle class="text-sm font-medium">{name}</CardTitle>
		<Icon class={'h-4 w-4 ' + colorClassName} />
	</CardHeader>
	<CardContent>
		<div class="text-2xl font-bold">{value}</div>
		<p class="text-xs text-muted-foreground">{description}</p>
	</CardContent>
</Card>
