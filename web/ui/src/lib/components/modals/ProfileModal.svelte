<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { RotateCcw, Trash2, TriangleAlert } from '@lucide/svelte';
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
	const deleteMutation = profile.delete();
	function onsubmit() {
		if (profileData.id) {
			updateMutation.mutate({ ...profileData });
		} else {
			createMutation.mutate({ ...profileData });
		}
		open = false;
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

			<div class="flex w-full gap-2">
				{#if profileData.id}
					<Popover.Root>
						<Popover.Trigger>
							<Button type="button" variant="destructive">
								<Trash2 size={16} />
							</Button>
						</Popover.Trigger>
						<Popover.Content class="w-80" align="end">
							<div class="space-y-4">
								<div class="flex items-start gap-3">
									<div
										class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-destructive/10"
									>
										<TriangleAlert class="h-4 w-4 text-destructive" />
									</div>
									<div class="flex-1 space-y-2">
										<h4 class="text-sm leading-none font-semibold">Delete Profile</h4>
										<p class="text-sm leading-relaxed text-muted-foreground">
											Are you sure you want to delete this profile?
										</p>
									</div>
								</div>
								<div class="flex items-center justify-end gap-2">
									<Button type="button" variant="outline">Cancel</Button>
									<Button
										type="button"
										variant="destructive"
										onclick={() => {
											deleteMutation.mutate({ id: profileData.id });
											open = false;
										}}
									>
										Delete
									</Button>
								</div>
							</div>
						</Popover.Content>
					</Popover.Root>
				{/if}
				<Button type="submit" class="flex-1">
					{profileData.id ? 'Update' : 'Create'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
