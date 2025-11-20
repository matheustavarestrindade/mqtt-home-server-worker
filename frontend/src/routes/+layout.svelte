<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<main class="container m-auto items-center">
	{@render children?.()}
</main>
