<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Eye, EyeOff } from 'lucide-svelte';
	import { page } from '$app/stores';

	let token = $page.url.searchParams.get('token');
	let password = $state('');
	let showPassword = $state(false);

	const handleSubmit = async () => {
		if (!token) return;
		await resetPassword(token, password);
	};
	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleSubmit();
		}
	};
</script>

<Card.Root class="w-[400px]">
	<Card.Header>
		<Card.Title>Reset Password</Card.Title>
		<Card.Description>Set your new password</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="grid w-full items-center gap-4" onkeydown={handleKeydown} aria-hidden>
			<div class="flex flex-col space-y-1.5">
				<Label for="password">Password</Label>
				<div class="flex flex-row items-center justify-end gap-1">
					{#if showPassword}
						<Input id="password" type="text" bind:value={password} />
					{:else}
						<Input id="password" type="password" bind:value={password} />
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
			<div class="mt-4 flex flex-col">
				<Button type="submit" on:click={handleSubmit}>Save</Button>
			</div>
		</div>
	</Card.Content>
</Card.Root>
