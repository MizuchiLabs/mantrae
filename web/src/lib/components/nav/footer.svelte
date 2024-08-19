<script lang="ts">
	import { version } from '$lib/api';
</script>

<footer class="flex flex-row items-center justify-end sm:px-6">
	{#if $version}
		<span class="text-xs text-muted-foreground">
			Traefik v{$version}
		</span>
	{/if}
</footer>
