<script lang="ts">
	import { goto } from '$app/navigation';
	import { handleOIDCLogin, profileClient, userClient } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import logo from '$lib/images/logo.svg';
	import { profile } from '$lib/stores/profile';
	import { user } from '$lib/stores/user';
	import { ConnectError } from '@connectrpc/connect';
	import { toast } from 'svelte-sonner';

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
	const handleLogin = async () => {
		const isEmail = username.includes('@');

		try {
			const response = await userClient.loginUser({
				identifier: {
					case: isEmail ? 'email' : 'username',
					value: username
				},
				password: password
			});
			if (!response.user) throw new Error('Authentication failed');
			user.value = response.user;

			if (!profile.id) {
				const response = await profileClient.listProfiles({});
				profile.value = response.profiles[0];
			}
			await goto('/');
			toast.success('Logged in successfully!');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to login', { description: e.message });
		}
	};
	const onkeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleLogin();
		}
	};
</script>

<svelte:head>
	<title>Login - Mantrae</title>
	<meta
		name="description"
		content="Sign in to your Mantrae account to manage your reverse proxy configurations"
	/>
</svelte:head>

{#if !user.isLoggedIn()}
	{#await userClient.getOIDCStatus({}) then value}
		<form
			class="m-auto h-fit w-full max-w-sm overflow-hidden rounded-[calc(var(--radius)+.125rem)] border bg-muted shadow-md shadow-zinc-950/5 dark:[--color-muted:var(--color-zinc-900)]"
		>
			<div class="-m-px rounded-[calc(var(--radius)+.125rem)] border bg-card p-8 pb-6">
				<div class="text-center">
					<img src={logo} alt="logo" class="mx-auto h-8 w-fit" />
					<h1 class="mt-4 mb-1 text-xl font-semibold">Sign In to Mantrae</h1>
					<p class="text-sm">Welcome back! Sign in to continue</p>
				</div>

				{#if value.loginEnabled}
					<div class="mt-6 space-y-5">
						<div class="space-y-2">
							<Label for="username" class="block text-sm">Username</Label>
							<Input id="username" bind:value={username} {onkeydown} />
						</div>

						<div class="space-y-0.5">
							<div class="flex items-center justify-between">
								<Label for="pwd" class="text-title text-sm">Password</Label>
								<Button
									variant="link"
									size="sm"
									class="link intent-info variant-ghost text-xs text-muted-foreground"
									onclick={handleReset}
								>
									Forgot your Password?
								</Button>
							</div>
							<PasswordInput bind:value={password} {onkeydown} />
						</div>

						<Button class="w-full" type="submit" onclick={handleLogin}>Sign In</Button>
					</div>
				{/if}

				{#if value.oidcEnabled}
					{#if value.loginEnabled}
						<div class="my-6 grid grid-cols-[1fr_auto_1fr] items-center gap-3">
							<hr class="border-dashed" />
							<span class="text-xs text-muted-foreground">Or continue With</span>
							<hr class="border-dashed" />
						</div>
					{:else}
						<Separator class="my-5" />
					{/if}

					<div class="flex flex-col gap-4">
						<Button variant="outline" onclick={handleOIDCLogin}>
							Login with {value.provider || 'OIDC'}
						</Button>
					</div>
				{/if}
			</div>
		</form>
	{/await}
{/if}
