<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { User } from '$lib/types/base';
	import type { Selected } from 'bits-ui';
	import { Bot, UserIcon } from 'lucide-svelte';

	export let user: User;
	export let pw = '';
	export let pwconfirm = '';
	export let disabled = false;

	const userTypes: Selected<string>[] = [
		{ label: 'User', value: 'user' },
		{ label: 'Machine', value: 'machine' }
	];
	let userType: Selected<string> | undefined = userTypes.find((u) => u.value === user.type);
	const setUserType = async (type: Selected<string> | undefined) => {
		if (type === undefined) return;
		user.type = type.value.toLowerCase();
	};
</script>

<Card.Root class="mt-4">
	<Card.Header>
		<Card.Title class="flex items-center justify-between gap-2">
			User
			{#if user.type === 'user'}
				<UserIcon />
			{:else}
				<Bot />
			{/if}
		</Card.Title>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Type -->
		{#if !disabled}
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
		{/if}

		<!-- Username -->
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="name" class="text-right">Name</Label>
			<Input
				name="name"
				type="text"
				bind:value={user.username}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Name"
				required
			/>
		</div>
		{#if user.type === 'user'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="Email" class="text-right">Email</Label>
				<Input
					name="Email"
					type="email"
					bind:value={user.email}
					class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
					placeholder="Email"
				/>
			</div>
		{/if}

		{#if disabled}
			<div class="mt-4 flex flex-row items-center justify-end text-sm">
				Leave blank to keep the same
			</div>
		{/if}

		<!-- Password -->
		{#if user.type === 'user'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="Password" class="text-right">Password</Label>
				<Input
					name="Password"
					type="password"
					bind:value={pw}
					class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
					placeholder="New Password"
					required
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="Password" class="text-right">Confirm Password</Label>
				<Input
					name="Password"
					type="password"
					bind:value={pwconfirm}
					class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
					placeholder="Confirm Password"
					required
				/>
			</div>
		{/if}
	</Card.Content>
</Card.Root>
