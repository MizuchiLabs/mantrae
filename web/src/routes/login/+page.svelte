<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { login, loggedIn, sendResetEmail } from '$lib/api';
	import { Eye, EyeOff } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let username = '';
	let password = '';
	let showPassword = false;
	let remember = false;
	const handleReset = async () => {
		if (username.length > 0) {
			await sendResetEmail(username);
		} else {
			toast.error('Please enter a username!');
		}
	};
	const handleSubmit = async () => {
		await login(username, password, remember);
	};
	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleSubmit();
		}
	};
</script>

{#if !$loggedIn}
	<Card.Root class="w-[400px]">
		<Card.Header>
			<Card.Title>Login</Card.Title>
			<Card.Description>Login to your account</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="grid w-full items-center gap-4" on:keydown={handleKeydown} aria-hidden>
				<div class="flex flex-col space-y-1.5">
					<Label for="username">Username</Label>
					<Input id="username" bind:value={username} />
				</div>
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
					<div class="flex flex-row items-center justify-between gap-1">
						<div class="items-top flex items-center justify-end gap-2">
							<Checkbox id="remember" bind:checked={remember} />
							<div class="grid gap-1.5 leading-none">
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
{/if}
