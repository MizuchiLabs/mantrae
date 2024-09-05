<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button';
	import { newUser, type User } from '$lib/types/base';
	import type { Selected } from 'bits-ui';
	import { createUser } from '$lib/api';

	let user: User = newUser();
	const userTypes: Selected<string>[] = [
		{ label: 'Regular', value: 'user' },
		{ label: 'Machine', value: 'machine' }
	];

	const create = async () => {
		if (user.username === '' || user.password === '' || user.type === '') return;
		await createUser(user);
		user = newUser();
		userType = userTypes[0];
	};

	let userType: Selected<string> | undefined = userTypes[0];
	const setUserType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		user.type = type.value.toLowerCase();
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<div class="flex w-full flex-row items-center justify-between">
			<Button class="flex items-center gap-2 bg-red-400 text-black">
				<span>Add User</span>
				<iconify-icon icon="fa6-solid:plus" />
			</Button>
		</div>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[500px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title>User</Card.Title>
				<Card.Description>Add a new regular user or machine user</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-2">
				<div class="grid grid-cols-4 items-center gap-4 space-y-2">
					<Label for="current" class="text-right">Type</Label>
					<Select.Root onSelectedChange={setUserType} selected={userType}>
						<Select.Trigger class="col-span-3">
							<Select.Value placeholder="Select a type" />
						</Select.Trigger>
						<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
							{#each userTypes as type}
								<Select.Item value={type.value} label={type.label}>
									{type.label}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="Username" class="text-right">Username</Label>
					<Input
						name="Username"
						type="text"
						placeholder="Username"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={user.username}
						required
					/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="Password" class="text-right">Password</Label>
					<Input
						name="Password"
						type="password"
						placeholder="Password"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={user.password}
						required
					/>
				</div>
				{#if user.type === 'user'}
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="Email" class="text-right">Email</Label>
						<Input
							name="Email"
							type="email"
							placeholder="Email"
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							bind:value={user.email}
						/>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
