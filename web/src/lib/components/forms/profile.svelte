<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import HoverInfo from '../utils/hoverInfo.svelte';
	import type { Profile } from '$lib/types/base';

	export let profile: Profile;
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
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="password" class="text-right">Password</Label>
		<Input
			name="password"
			type="password"
			class="col-span-3"
			bind:value={profile.password}
			placeholder="Basic auth password"
			required
		/>
	</div>
</div>
