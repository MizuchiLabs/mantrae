<script lang="ts">
	import { activeProfile } from '$lib/api';
</script>

<footer class="flex flex-row items-center justify-end sm:px-6">
	{#if $activeProfile?.instance?.dynamic?.version}
		<span class="text-xs text-muted-foreground">
			Traefik v{$activeProfile?.instance?.dynamic?.version}
		</span>
	{/if}
</footer>
