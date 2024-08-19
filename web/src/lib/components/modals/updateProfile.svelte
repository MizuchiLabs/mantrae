<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { deleteProfile, profiles, updateProfile } from '$lib/api';

	export let name: string;
	let oldName = $profiles[name].name;

	const update = () => {
		updateProfile(oldName, $profiles[name]);
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};
</script>

<Dialog.Root>
	<Dialog.Trigger class="flex flex-row justify-end">
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-orange-400">
			<iconify-icon icon="fa6-solid:pencil" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Update profile</Dialog.Title>
		</Dialog.Header>
		<div class="grid gap-4 py-4" on:keydown={onKeydown} aria-hidden>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input
					name="name"
					type="text"
					class="col-span-3"
					placeholder="Your profile name"
					bind:value={$profiles[name].name}
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="url" class="text-right">URL</Label>
				<Input
					name="url"
					type="text"
					class="col-span-3"
					bind:value={$profiles[name].url}
					placeholder="URL of your client"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input
					name="username"
					type="text"
					class="col-span-3"
					bind:value={$profiles[name].username}
					placeholder="Username of your client"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="password" class="text-right">Password</Label>
				<Input
					name="password"
					type="password"
					class="col-span-3"
					bind:value={$profiles[name].password}
					placeholder="Password of your client"
					required
				/>
			</div>
		</div>
		<Dialog.Close class="flex w-full flex-row gap-2">
			<Button type="submit" class="w-full bg-red-400" on:click={() => deleteProfile(name)}>
				Delete
			</Button>
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
