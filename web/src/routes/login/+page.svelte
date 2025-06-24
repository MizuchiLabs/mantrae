<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import { profileClient, userClient } from '$lib/api';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';
	import { ConnectError } from '@connectrpc/connect';
	import { profile } from '$lib/stores/profile';
	import { handleOIDCLogin } from '$lib/api';

	let username = $state('');
	let password = $state('');

	const handleReset = async () => {
		if (username.length > 0) {
			toast.error('Please enter a username!');
			return;
		}
		const isEmail = username.includes('@');

		try {
			await userClient.sendOTP({
				identifier: {
					case: isEmail ? 'email' : 'username',
					value: username
				}
			});

			await goto(`/login/reset?username=${username}`);
			toast.success('Code sent successfully!');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to send reset code', { description: e.message });
		}
	};
	const handleSubmit = async () => {
		const isEmail = username.includes('@');

		try {
			await userClient.loginUser({
				identifier: {
					case: isEmail ? 'email' : 'username',
					value: username
				},
				password: password
			});

			const verified = await userClient.verifyJWT({});
			if (verified.user) {
				user.value = verified.user;
				if (!profile.id) {
					const response = await profileClient.listProfiles({});
					profile.value = response.profiles[0];
				}
				await goto('/');
			}
			toast.success('Logged in successfully!');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to login', { description: e.message });
		}
	};

	// const handleOIDCLogin = () => {
	// 	window.location.href = '/oidc/login';
	// };
</script>

{#if !user.isLoggedIn()}
	<Card.Root class="max-w-md">
		<Card.Header class="flex flex-col items-center text-center">
			<Card.Title class="text-2xl font-bold">Welcome back</Card.Title>
			<Card.Description>Login to your account</Card.Description>
		</Card.Header>
		<Card.Content>
			<form onsubmit={handleSubmit} class="p-4">
				{#await userClient.getOIDCStatus({}) then value}
					{#if value.loginEnabled}
						<div class="flex flex-col gap-4">
							<div class="grid gap-3">
								<Label for="username">Username</Label>
								<Input id="username" bind:value={username} />
							</div>

							<div class="grid gap-3">
								<div class="flex items-center">
									<Label for="password">Password</Label>
									<button
										class="text-muted-foreground ml-auto text-xs hover:underline"
										type="button"
										onclick={handleReset}
									>
										Forgot your password?
									</button>
								</div>
								<PasswordInput bind:value={password} />
							</div>

							<Button type="submit" class="w-full">Login</Button>

							{#if value.oidcEnabled}
								<div
									class="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t"
								>
									<span class="bg-background text-muted-foreground relative z-10 px-2">
										Or continue with
									</span>
								</div>
							{/if}
						</div>
					{/if}

					{#if value.oidcEnabled}
						<Button variant="outline" class="mt-3 w-full" onclick={handleOIDCLogin}>
							Login with {value.provider || 'OIDC'}
						</Button>
					{/if}
				{/await}
			</form>
		</Card.Content>
	</Card.Root>
{/if}
