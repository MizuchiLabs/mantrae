<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	onMount(() => {
		api.listProfiles();
	});
</script>

<div class="mx-auto max-w-7xl"></div>
