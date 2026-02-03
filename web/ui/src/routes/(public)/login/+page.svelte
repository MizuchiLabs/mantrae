<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import logo from '$lib/assets/logo.svg';
	import { user } from '$lib/api/users.svelte';
	import { BackendURL } from '$lib/config';

	let username = $state('');
	let password = $state('');

	const login = user.login();
	const oidcStatus = user.oidc();

	const handleLogin = async () => {
		const isEmail = username.includes('@');
		login.mutate({
			password: password,
			identifier: {
				case: isEmail ? 'email' : 'username',
				value: username
			}
		});
	};
	const onkeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleLogin();
		}
	};
	function handleOIDCLogin() {
		window.location.href = `${BackendURL}/oidc/login`;
	}
</script>

<svelte:head>
	<title>Login</title>
	<meta
		name="description"
		content="Sign in to your Mantrae account to manage your reverse proxy configurations"
	/>
</svelte:head>

<form
	class="m-auto h-fit w-full max-w-sm overflow-hidden rounded-[calc(var(--radius)+.125rem)] border bg-muted shadow-md shadow-zinc-950/5 dark:[--color-muted:var(--color-zinc-900)]"
>
	<div class="-m-px rounded-[calc(var(--radius)+.125rem)] border bg-card p-8 pb-6">
		<div class="text-center">
			<img src={logo} alt="logo" class="mx-auto h-8 w-fit" />
			<h1 class="mt-4 mb-1 text-xl font-semibold">Sign In to Mantrae</h1>
			<p class="text-sm">Welcome back! Sign in to continue</p>
		</div>

		{#if oidcStatus.isSuccess}
			{#if oidcStatus.data.loginEnabled}
				<div class="mt-6 space-y-5">
					<div class="space-y-2">
						<Label for="username" class="block text-sm">Username</Label>
						<Input id="username" bind:value={username} {onkeydown} />
					</div>

					<div class="space-y-0.5">
						<Label for="pwd" class="text-title text-sm">Password</Label>
						<PasswordInput bind:value={password} {onkeydown} />
					</div>

					<Button class="w-full" type="submit" onclick={handleLogin}>Sign In</Button>
				</div>
			{/if}

			{#if oidcStatus.data.loginEnabled && oidcStatus.data.oidcEnabled}
				<div class="my-6 grid grid-cols-[1fr_auto_1fr] items-center gap-3">
					<hr class="border-dashed" />
					<span class="text-xs text-muted-foreground">Or continue With</span>
					<hr class="border-dashed" />
				</div>

				<Separator class="my-5" />
			{/if}

			{#if oidcStatus.data.oidcEnabled}
				<div class="flex flex-col gap-4">
					<Button variant="outline" onclick={handleOIDCLogin}>
						Login with {oidcStatus.data.provider || 'OIDC'}
					</Button>
				</div>
			{/if}
		{/if}
	</div>
</form>
