<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { toast } from 'svelte-sonner';
	import { userClient } from '$lib/api';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';
	import type { OAuthStatus } from '$lib/types';
	import { token } from '$lib/stores/common';
	import { ConnectError } from '@connectrpc/connect';

	let username = $state('');
	let password = $state('');
	let remember = $state(false);
	// let oauthStatus: OAuthStatus = $state({ enabled: false, provider: '' });
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
			const res = await userClient.loginUser({
				identifier: {
					case: isEmail ? 'email' : 'username',
					value: username
				},
				password: password,
				remember: remember
			});
			token.value = res.token ?? null;

			if (!token.value) {
				throw new Error('No token received');
			}

			const verified = await userClient.verifyJWT({ token: token.value });
			if (verified.userId) {
				await goto('/');
			}
			toast.success('Logged in successfully!');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to login', { description: e.message });
		}
	};
	// const handleOIDCLogin = () => {
	// 	window.location.href = '/api/oidc/login';
	// };
	// onMount(async () => {
	// 	oauthStatus = await api.oauthStatus();
	// });
</script>

{#if !user.isLoggedIn()}
	<Card.Root class="w-[400px]">
		<Card.Header>
			<Card.Title>Login</Card.Title>
			<Card.Description>Login to your account</Card.Description>
		</Card.Header>
		<Card.Content>
			<form onsubmit={handleSubmit} class="space-y-4">
				<!-- {#if !oauthStatus.loginDisabled} -->
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

				<Button type="submit" class="w-full">Login</Button>
				<!-- {/if} -->

				<!-- {#if oauthStatus.enabled} -->
				<!-- 	<Button variant="outline" class="w-full" onclick={handleOIDCLogin}> -->
				<!-- 		Login with {oauthStatus.provider || 'OIDC'} -->
				<!-- 	</Button> -->
				<!-- {/if} -->
			</form>
		</Card.Content>
	</Card.Root>
{/if}
