<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { toast } from 'svelte-sonner';
	import { api } from '$lib/api';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';

	let username = '';
	let password = '';
	let remember = false;
	const handleReset = async () => {
		if (username.length > 0) {
			// await api.sendResetEmail(username);
		} else {
			toast.error('Please enter a username!');
		}
	};
	const handleSubmit = async () => {
		await api.login(username, password, remember);
	};
	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleSubmit();
		}
	};
</script>

<Card.Root class="w-[400px]">
	<Card.Header>
		<Card.Title>Login</Card.Title>
		<Card.Description>Login to your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="grid w-full items-center gap-4" on:keydown={handleKeydown} aria-hidden>
			<div class="flex flex-col gap-2">
				<Label for="username">Username</Label>
				<Input id="username" bind:value={username} />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="password">Password</Label>
				<PasswordInput bind:password />
				<div class="mt-1 flex flex-row items-center justify-between">
					<div class="items-top flex items-center justify-end gap-2">
						<Checkbox id="remember" bind:checked={remember} />
						<div class="grid gap-2 leading-none">
							<Label for="terms1" class="text-sm">Remember me</Label>
						</div>
					</div>
					<button class="text-xs text-muted-foreground" on:click={handleReset}>
						Forgot password?
					</button>
				</div>
			</div>

			<div class="mt-4 flex flex-col">
				<Button type="submit" on:click={handleSubmit}>Login</Button>
			</div>
		</div>
	</Card.Content>
</Card.Root>
