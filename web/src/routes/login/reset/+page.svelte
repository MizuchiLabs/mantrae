<script lang="ts">
	import * as InputOTP from '$lib/components/ui/input-otp/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { REGEXP_ONLY_DIGITS } from 'bits-ui';
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';

	let token = $state('');

	// Clean and trim the value
	const onPaste = (value: string) => value.replace(/[^0-9]/g, '').trim();
	const onComplete = async () => {
		const username = page.url.searchParams.get('username');
		if (!username || !token) return;
		try {
			await api.verifyOTP(username, token);
			toast.success('Token verified successfully!', {
				description: 'Please reset your password!'
			});
		} catch (error) {
			let err = error as Error;
			toast.error('Failed to verify Token: ', {
				description: err.message
			});
		}
	};
	// Handle Ctrl+V paste
	const handleKeyDown = async (e: KeyboardEvent) => {
		if (e.ctrlKey && e.key === 'v') {
			e.preventDefault();
			try {
				const text = await navigator.clipboard.readText();
				const cleaned = onPaste(text);
				token = cleaned.slice(0, 6); // Limit to max length
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
			bind:value={token}
			{onComplete}
			{onPaste}
			onkeydown={handleKeyDown}
		>
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
	</Card.Content>
</Card.Root>
