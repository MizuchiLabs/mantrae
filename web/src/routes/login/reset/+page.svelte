<script lang="ts">
	import * as InputOTP from '$lib/components/ui/input-otp/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { REGEXP_ONLY_DIGITS } from 'bits-ui';
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

	let token = $state('');
	let password = $state('');

	const onSubmit = async () => {
		try {
			await api.resetPassword(page.params.username, token, password);
			goto('/login');
			toast.success('Password reset successfully');
		} catch (error) {
			let err = error as Error;
			toast.error('Failed to reset password: ', {
				description: err.message
			});
		}
	};
</script>

<Card.Root class="w-[400px]">
	<Card.Header>
		<Card.Title>Reset Token</Card.Title>
		<Card.Description>Please enter your reset token and set your new password</Card.Description>
	</Card.Header>
	<Card.Content>
		<form onsubmit={onSubmit} class="flex flex-col gap-4">
			<div class="space-y-1 self-center">
				<Label for="token">Token</Label>
				<InputOTP.Root maxlength={6} pattern={REGEXP_ONLY_DIGITS} bind:value={token}>
					{#snippet children({ cells })}
						<InputOTP.Group>
							{#each cells.slice(0, 3) as cell}
								<InputOTP.Slot {cell} />
							{/each}
						</InputOTP.Group>
						<InputOTP.Separator />
						<InputOTP.Group>
							{#each cells.slice(3, 6) as cell}
								<InputOTP.Slot {cell} />
							{/each}
						</InputOTP.Group>
					{/snippet}
				</InputOTP.Root>
			</div>

			<div class="space-y-1">
				<Label for="password">New Password</Label>
				<PasswordInput bind:password />
			</div>

			<Button type="submit">Reset</Button>
		</form>
	</Card.Content>
</Card.Root>
