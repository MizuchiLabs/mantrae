<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { toast } from 'svelte-sonner';
	import { api, loading } from '$lib/api';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { goto } from '$app/navigation';

	let username = $state('');
	let password = $state('');
	let remember = $state(false);
	const handleReset = () => {
		if (username.length > 0) {
			api.sendResetEmail(username);
			goto(`/login/reset?username=${username}`);
		} else {
			toast.error('Please enter a username!');
		}
	};
	const handleSubmit = async () => {
		await api.login(username, password, remember);
	};
</script>

<Card.Root class="w-[400px]">
	<Card.Header>
		<Card.Title>Login</Card.Title>
		<Card.Description>Login to your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<form onsubmit={handleSubmit} class="space-y-4">
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
					<button class="text-xs text-muted-foreground" type="button" onclick={handleReset}>
						Forgot password?
					</button>
				</div>
			</div>

			<Separator />

			<Button type="submit" class="w-full" disabled={$loading}>Login</Button>
		</form>
	</Card.Content>
</Card.Root>
