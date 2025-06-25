<script lang="ts">
	import * as InputOTP from '$lib/components/ui/input-otp/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { REGEXP_ONLY_DIGITS } from 'bits-ui';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { userClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';

	let otp = $state('');

	// Clean and trim the value
	const pasteTransformer = (value: string) => value.replace(/[^0-9]/g, '').trim();
	const onComplete = async () => {
		const username = page.url.searchParams.get('username');
		if (!username || !otp) return;
		const isEmail = username.includes('@');

		try {
			await userClient.verifyOTP({
				identifier: {
					case: isEmail ? 'email' : 'username',
					value: username
				},
				otp
			});
			const verified = await userClient.verifyJWT({});
			if (verified.user) {
				user.value = verified.user;
				await goto('/');
			}
			toast.success('Token verified successfully!', {
				description: 'Please reset your password!'
			});
		} catch (err) {
			let e = ConnectError.from(err);
			toast.error('Failed to verify Token', { description: e.message });
		}
	};
	// Handle Ctrl+V paste
	const handleKeyDown = async (e: KeyboardEvent) => {
		if (e.ctrlKey && e.key === 'v') {
			e.preventDefault();
			try {
				const text = await navigator.clipboard.readText();
				const cleaned = pasteTransformer(text);
				otp = cleaned.slice(0, 6); // Limit to max length
			} catch (err) {
				console.error('Failed to read clipboard:', err);
			}
		}
	};
</script>

<Card.Root class="w-[400px]">
	<Card.Header>
		<Card.Title>Reset Token</Card.Title>
		<Card.Description>Please enter the one-time password sent to your email</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col items-center gap-2">
		<InputOTP.Root
			maxlength={6}
			pattern={REGEXP_ONLY_DIGITS}
			bind:value={otp}
			{onComplete}
			{pasteTransformer}
			onkeydown={handleKeyDown}
		>
			{#snippet children({ cells })}
				<InputOTP.Group>
					{#each cells.slice(0, 3) as cell (cell)}
						<InputOTP.Slot {cell} />
					{/each}
				</InputOTP.Group>
				<InputOTP.Separator />
				<InputOTP.Group>
					{#each cells.slice(3, 6) as cell (cell)}
						<InputOTP.Slot {cell} />
					{/each}
				</InputOTP.Group>
			{/snippet}
		</InputOTP.Root>
	</Card.Content>
</Card.Root>
