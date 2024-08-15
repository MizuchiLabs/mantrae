<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import type { Selected } from 'bits-ui';
	import { activeProfile, updateProfile } from '$lib/api';
	import type { HttpMiddleware, TCPMiddleware } from '$lib/types/middlewares';

	export let httpMiddleware: HttpMiddleware | undefined;
	export let tcpMiddleware: TCPMiddleware | undefined;

	const update = () => {
		if (httpMiddleware === undefined && tcpMiddleware === undefined) return;
		if (httpMiddleware !== undefined) {
			httpMiddleware.name = httpMiddleware.name.trim();
			if (httpMiddleware.name === '') return;
			activeProfile.update((p) => ({
				...p,
				instance: {
					...p.instance,
					dynamic: {
						...p.instance.dynamic,
						httpmiddlewares: [...(p.instance.dynamic?.httpmiddlewares || []), { ...httpMiddleware }]
					}
				}
			}));
		}
		if (tcpMiddleware !== undefined) {
			tcpMiddleware.name = tcpMiddleware.name.trim();
			if (tcpMiddleware.name === '') return;
			activeProfile.update((p) => ({
				...p,
				instance: {
					...p.instance,
					dynamic: {
						...p.instance.dynamic,
						tcpmiddlewares: [...(p.instance.dynamic?.tcpmiddlewares || []), { ...tcpMiddleware }]
					}
				}
			}));
		}
		updateProfile($activeProfile.name, $activeProfile);
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-orange-400">
			<iconify-icon icon="fa6-solid:pencil" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[520px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title>Middleware</Card.Title>
				<Card.Description>
					Make changes to your Middleware here. Click save when you're done.
				</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-2">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="name" class="text-right">Name</Label>
					{#if httpMiddleware !== undefined}
						<Input
							id="name"
							name="name"
							type="text"
							bind:value={httpMiddleware.name}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Name of the middleware"
							required
						/>
					{:else if tcpMiddleware !== undefined}
						<Input
							id="name"
							name="name"
							type="text"
							bind:value={tcpMiddleware.name}
							class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
							placeholder="Name of the middleware"
							required
						/>
					{/if}
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
