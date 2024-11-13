<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { User } from '$lib/types/base';
	import { Crown, UserIcon } from 'lucide-svelte';
	import { z } from 'zod';

	export let user: User;
	export let disabled = false;
	let pw = '';
	let pwconfirm = '';

	const schema = z.object({
		username: z.string().trim().min(1, 'Username is required').max(255),
		email: z.union([z.literal(''), z.literal(null), z.string().trim().email().max(255)]),
		passCheck: z
			.object({
				password: z.string().trim().nullish(),
				confirm: z.string().trim().nullish()
			})
			.refine((data) => data.password === data.confirm, {
				message: "Passwords don't match"
			})
	});

	let errors: Record<any, string[] | undefined> = {};
	export const validate = () => {
		try {
			schema.parse({
				username: user.username,
				email: user.email,
				passCheck: {
					password: pw,
					confirm: pwconfirm
				}
			});
			if (pw !== '') user.password = pw;
			errors = {};
			return true;
		} catch (error) {
			if (error instanceof z.ZodError) {
				errors = error.flatten().fieldErrors;
			}
			return false;
		}
	};
</script>

<Card.Root class="mt-4">
	<Card.Header>
		<Card.Title class="flex items-center justify-between gap-2">
			User
			{#if user.isAdmin}
				<Crown />
			{:else}
				<UserIcon />
			{/if}
		</Card.Title>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Admin -->
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="admin" class="text-right">Make Admin</Label>
			<Switch id="admin" bind:checked={user.isAdmin} class="col-span-3" />
		</div>

		<!-- Username -->
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="name" class="text-right">Name</Label>
			<Input
				name="name"
				type="text"
				bind:value={user.username}
				on:input={validate}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Name"
				required
			/>
			{#if errors?.username}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.username}</div>
			{/if}
		</div>

		<!-- Email -->
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="Email" class="text-right">Email</Label>
			<Input
				name="Email"
				type="email"
				bind:value={user.email}
				on:input={validate}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Email"
			/>
			{#if errors?.email}
				<div class="col-span-4 text-right text-sm text-red-500">{errors.email}</div>
			{/if}
		</div>

		{#if disabled}
			<div class="mt-4 flex flex-row items-center justify-end text-sm">
				Leave blank to keep the same
			</div>
		{/if}

		<!-- Password -->
		<div class="grid grid-cols-4 items-center gap-4">
			<Label for="Password" class="text-right">Password</Label>
			<Input
				name="Password"
				type="password"
				bind:value={pw}
				on:input={validate}
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
				on:input={validate}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Confirm Password"
				required
			/>
			{#if errors?.passCheck}
				<div class="col-span-4 text-right text-sm text-red-500">
					{errors.passCheck}
				</div>
			{/if}
		</div>
	</Card.Content>
</Card.Root>
