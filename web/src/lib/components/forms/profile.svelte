<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import HoverInfo from '../utils/hoverInfo.svelte';
	import type { Profile } from '$lib/types/base';
	import { Eye, EyeOff } from 'lucide-svelte';

	interface Props {
		profile: Profile;
	}

	let { profile = $bindable() }: Props = $props();
	let showPassword = $state(false);
</script>

<div class="grid gap-4 py-4">
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="name" class="text-right">Name</Label>
		<Input
			name="name"
			type="text"
			bind:value={profile.name}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="Your profile name"
			required
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="url" class="text-right">URL</Label>
		<Input
			name="url"
			type="text"
			class="col-span-3"
			bind:value={profile.url}
			placeholder="URL of your traefik instance"
			required
		/>
	</div>
	<div class="flex flex-row items-center justify-end gap-2">
		<Label for="tls" class="flex items-center gap-0.5">
			Verify Certificate
			<HoverInfo
				text="If your Traefik instance uses a self-signed certificate, you can enable/disable certificate verification here."
			/>
		</Label>
		<Checkbox name="tls" bind:checked={profile.tls} required />
	</div>

	<span class="mt-2 flex flex-row items-center gap-1 border-b border-gray-200 pb-2">
		<span class="font-bold">Basic Authentication</span>
		<HoverInfo
			text="If your Traefik instance requires basic authentication, you can enter your username and password here."
		/>
	</span>

	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="username" class="text-right">Username</Label>
		<Input
			name="username"
			type="text"
			class="col-span-3"
			bind:value={profile.username}
			placeholder="Basic auth username"
			required
		/>
	</div>
	<div class="relative grid grid-cols-4 items-center gap-4">
		<Label for="password" class="text-right">Password</Label>
		<div class="col-span-3 flex flex-row items-center justify-end gap-1">
			{#if showPassword}
				<Input
					name="password"
					type="text"
					class="col-span-3 pr-10"
					bind:value={profile.password}
					placeholder="Basic auth password"
					required
				/>
			{:else}
				<Input
					name="password"
					type="password"
					class="col-span-3 pr-10"
					bind:value={profile.password}
					placeholder="Basic auth password"
					required
				/>
			{/if}
			<Button
				variant="ghost"
				size="icon"
				class="absolute hover:bg-transparent hover:text-red-400"
				on:click={() => (showPassword = !showPassword)}
			>
				{#if showPassword}
					<Eye size="1rem" />
				{:else}
					<EyeOff size="1rem" />
				{/if}
			</Button>
		</div>
	</div>
</div>
