<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { downloadBackup, uploadBackup } from '$lib/api';

	let fileInput: HTMLInputElement;
	const handleFileUpload = (event: Event) => {
		const file = (event.target as HTMLInputElement).files?.[0];
		if (file) {
			uploadBackup(file);
		}
		fileInput.value = '';
	};
</script>

<div class="container mx-auto flex flex-col items-center justify-center gap-4 py-4">
	<Card.Root class="w-[600px]">
		<Card.Header>
			<Card.Title class="text-lg font-bold">Settings</Card.Title>
		</Card.Header>
		<Card.Content class="mt-4">
			<div class="grid grid-cols-4 items-center justify-between gap-2">
				<Label for="backup" class="col-span-1">Backup & Restore</Label>
				<div class="col-span-3 flex w-full gap-2">
					<input
						type="file"
						accept=".tar.gz"
						class="hidden"
						on:change={handleFileUpload}
						bind:this={fileInput}
						required
					/>
					<Button variant="ghost" class="w-full bg-orange-400" on:click={() => fileInput.click()}>
						<iconify-icon icon="fa6-solid:upload" width="16" height="16" />
					</Button>
					<Button variant="default" class="w-full" on:click={() => downloadBackup()}>
						<iconify-icon icon="fa6-solid:download" width="16" height="16" />
					</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>
