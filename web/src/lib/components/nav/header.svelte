<script>
	import * as Avatar from '$lib/components/ui/avatar/index';
	import { Button } from '$lib/components/ui/button/index';
	import { logout, profiles } from '$lib/api';
	import Profile from './profile.svelte';
	import { goto } from '$app/navigation';
	import InfoTraefik from '../modals/infoTraefik.svelte';

	const handleLogout = () => {
		logout();
		goto('/login');
	};
</script>

<nav class="flex h-16 items-center justify-between border-b bg-primary-foreground">
	<div class="ml-6 flex flex-row items-center gap-4">
		<Profile />
		{#if $profiles?.length === 0 || !$profiles}
			<span class="flex items-center gap-1 text-sm text-muted-foreground">
				<iconify-icon icon="fa6-solid:arrow-left" />
				No profiles configured, create one here
			</span>
		{:else}
			<InfoTraefik />
		{/if}
	</div>

	<div class="mr-2 flex flex-row items-center gap-2">
		<!-- TODO: Add notifications -->
		<!-- <Button variant="outline" class="h-8 w-8 rounded-full"> -->
		<!-- 	<iconify-icon icon="fa6-solid:bell" /> -->
		<!-- </Button> -->
		<!-- <Avatar.Root> -->
		<!-- 	<Avatar.Image src="" alt="@user" /> -->
		<!-- 	<Avatar.Fallback>AD</Avatar.Fallback> -->
		<!-- </Avatar.Root> -->
		<Button variant="ghost" on:click={handleLogout} class="h-8 w-8 rounded-full">
			<iconify-icon icon="fa6-solid:right-from-bracket" />
		</Button>
	</div>
</nav>
