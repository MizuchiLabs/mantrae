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
	import { user } from '$lib/stores/user';
	import type { OAuthStatus } from '$lib/types';
	import { onMount } from 'svelte';

	let username = $state('');
	let password = $state('');
	let remember = $state(false);
	let oauthStatus: OAuthStatus = $state({ enabled: false, provider: '' });
	const handleReset = async () => {
		if (username.length > 0) {
			const resetPromise = api
				.sendResetEmail(username)
				.then(() => {
					goto(`/login/reset?username=${username}`);
					return 'Reset email sent successfully!';
				})
				.catch((error) => {
					throw new Error(error.message || 'Failed to send reset email');
				});

			toast.promise(resetPromise, {
				loading: 'Sending reset email...',
				success: (message) => message,
				error: (error) => (error as Error).message
			});
		} else {
			toast.error('Please enter a username!');
		}
	};
	const handleSubmit = async () => {
		const loginPromise = api
			.login(username, password, remember)
			.then(() => {
				return 'Logged in successfully!';
			})
			.catch((error) => {
				throw new Error(error.message || 'Login failed');
			});

		toast.promise(loginPromise, {
			loading: 'Logging in...',
			success: (message) => message,
			error: (error) => (error as Error).message
		});
	};
	const handleOIDCLogin = () => {
		window.location.href = '/api/oidc/login';
	};
	onMount(async () => {
		oauthStatus = await api.oauthStatus();
	});
</script>

{#if !user.isLoggedIn()}
	<Card.Root class="w-[400px]">
		<Card.Header>
			<Card.Title>Login</Card.Title>
			<Card.Description>Login to your account</Card.Description>
		</Card.Header>
		<Card.Content>
			<form onsubmit={handleSubmit} class="space-y-4">
				{#if !oauthStatus.loginDisabled}
					<div class="flex flex-col gap-2">
						<Label for="username">Username</Label>
						<Input id="username" bind:value={username} />
					</div>

					<div class="flex flex-col gap-2">
						<Label for="password">Password</Label>
						<PasswordInput bind:value={password} />
						<div class="mt-1 flex flex-row items-center justify-between">
							<div class="items-top flex items-center justify-end gap-2">
								<Checkbox id="remember" bind:checked={remember} />
								<div class="grid gap-2 leading-none">
									<Label for="terms1" class="text-sm">Remember me</Label>
								</div>
							</div>
							<button class="text-muted-foreground text-xs" type="button" onclick={handleReset}>
								Forgot password?
							</button>
						</div>
					</div>

					<Separator />

					<Button type="submit" class="w-full" disabled={$loading}>Login</Button>
				{/if}

				{#if oauthStatus.enabled}
					<Button variant="outline" class="w-full" onclick={handleOIDCLogin}>
						Login with {oauthStatus.provider || 'OIDC'}
					</Button>
				{/if}
			</form>
		</Card.Content>
	</Card.Root>
{/if}
