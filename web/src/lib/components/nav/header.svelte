<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import { api, profiles } from '$lib/api';
	import Profile from './profile.svelte';
	import InfoModal from '../modals/info.svelte';
	import { ArrowLeft, LogOut } from 'lucide-svelte';
</script>

<nav class="flex h-16 items-center justify-between border-b bg-primary-foreground pl-4">
	<div class="flex flex-row items-center gap-2">
		<Profile />
		{#if !$profiles}
			<span class="ml-2 flex items-center gap-2 text-sm text-muted-foreground">
				<ArrowLeft size="1rem" />
				No profiles configured, create one here
			</span>
		{:else}
			<InfoModal />
		{/if}
	</div>

	<div class="mr-2 flex flex-row items-center gap-2">
		<Button variant="ghost" onclick={api.logout} size="icon">
			<LogOut size="1rem" />
		</Button>
	</div>
</nav>
