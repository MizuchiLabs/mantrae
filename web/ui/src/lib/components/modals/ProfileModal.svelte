<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { RotateCcw } from '@lucide/svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { CopyInput } from '../ui/input-group';
	import { profile } from '$lib/api/profiles.svelte';
	import { setting } from '$lib/api/settings.svelte';

	interface Props {
		data?: Profile;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let profileData = $state({} as Profile);
	$effect(() => {
		if (data) profileData = { ...data };
	});
	$effect(() => {
		if (!open) profileData = {} as Profile;
	});

	const serverURL = setting.get('server_url');
	const createMutation = profile.create();
	const updateMutation = profile.update();
	function onsubmit() {
		if (profileData.id) {
			updateMutation.mutate({ ...profileData });
		} else {
			createMutation.mutate({ ...profileData });
		}
	}

	let newProfile = $derived(profile.get(profileData.id));
	function regenerate() {
		updateMutation.mutate({ ...profileData, regenerateToken: true });
		if (newProfile.isSuccess && newProfile.data?.token) {
			profileData.token = newProfile.data.token;
		}
	}

	let connString = $derived(
		`${serverURL.data?.value}/api/${profileData.name}?token=${profileData.token}`
	);
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-100 overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{profileData?.id ? 'Edit' : 'Create'} Profile</Dialog.Title>
			<Dialog.Description>Configure your profile settings</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-4">
			<div class="space-y-2">
				<Label for="name" class="text-sm font-medium">Name</Label>
				<Input
					id="name"
					bind:value={profileData.name}
					placeholder="traefik-site"
					class="transition-colors"
				/>
				<p class="text-xs text-muted-foreground">A descriptive name for this profile</p>
			</div>

			<div class="space-y-2">
				<Label for="description" class="text-sm font-medium">Description</Label>
				<Input
					id="description"
					bind:value={profileData.description}
					placeholder="Site description"
				/>
				<p class="text-xs text-muted-foreground">Optional description for this profile</p>
			</div>

			{#if profileData.id}
				<div class="space-y-2">
					<Label for="token" class="text-sm font-medium">Connection Token</Label>
					<div class="flex gap-2">
						<CopyInput value={profileData.token} readonly />
						<Button variant="outline" size="icon" onclick={regenerate} title="Regenerate token">
							<RotateCcw class="h-4 w-4" />
						</Button>
					</div>
					<p class="text-xs text-muted-foreground">
						Used in the connection URL to connect to this profile with your traefik instance
						<span class="underline">
							{connString}
						</span>
					</p>
				</div>
			{/if}

			<Separator />

			<Button type="submit" class="w-full">{profileData.id ? 'Update' : 'Create'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
