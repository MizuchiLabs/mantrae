<script>
	import * as Avatar from '$lib/components/ui/avatar/index';
	import { Button } from '$lib/components/ui/button/index';
	import { logout, profiles } from '$lib/api';
	import Profile from './profile.svelte';
	import { goto } from '$app/navigation';
	import InfoModal from '../modals/info.svelte';
	import Warning from '../modals/warning.svelte';
	import { ArrowLeft, LogOut } from 'lucide-svelte';

	const handleLogout = () => {
		logout();
		goto('/login');
	};
</script>

<nav class="flex h-16 items-center justify-between border-b bg-primary-foreground">
	<div class="ml-4 flex flex-row items-center gap-2">
		<Profile />
		{#if $profiles?.length === 0 || !$profiles}
			<span class="ml-2 flex items-center gap-2 text-sm text-muted-foreground">
				<ArrowLeft size="1rem" />
				No profiles configured, create one here
			</span>
		{:else}
			<InfoModal />
		{/if}
		<Warning />
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
		<Button variant="ghost" on:click={handleLogout} size="icon">
			<LogOut size="1rem" />
		</Button>
	</div>
</nav>
