<script lang="ts">
	import { activeProfile } from '$lib/api';
</script>

<footer class="flex flex-row items-center justify-end sm:px-6">
	{#if $activeProfile?.client?.dynamic?.version}
		<span class="text-xs text-muted-foreground">
			Traefik v{$activeProfile?.client?.dynamic?.version}
		</span>
	{/if}
</footer>
