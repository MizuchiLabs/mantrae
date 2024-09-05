<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { User } from '$lib/types/base';
	import { updateUser } from '$lib/api';

	export let u: User;
	let open = false;

	const update = async () => {
		if (u.username === '' || u.password === '') return;
		await updateUser(u);
		open = false;
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class="w-full">
		<Button class="w-full bg-orange-400 text-black">Edit</Button>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[500px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					User
					{#if u.type === 'user'}
						<iconify-icon icon="fa6-solid:user" class="text-green-400" />
					{:else}
						<iconify-icon icon="fa6-solid:robot" class="text-blue-400" />
					{/if}
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-2" on:keydown={onKeydown} aria-hidden>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="Username" class="text-right">Username</Label>
						<Input
							name="Username"
							type="text"
							bind:value={u.username}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Username"
							required
						/>
					</div>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="Password" class="text-right">Password</Label>
						<Input
							name="Password"
							type="password"
							bind:value={u.password}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Password"
							required
						/>
					</div>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="Email" class="text-right">Email</Label>
						<Input
							name="Email"
							type="email"
							bind:value={u.email}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Email"
						/>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Update</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
