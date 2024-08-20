<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { profile, updateMiddleware } from '$lib/api';
	import type { Middleware } from '$lib/types/middlewares';

	export let middleware: Middleware;
	let name = middleware.name.split('@')[0];

	const update = () => {
		let oldName = middleware.name;
		middleware.name = name + '@' + middleware.provider;
		updateMiddleware($profile, middleware, oldName);
	};

	const onKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			update();
		}
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
					<Input
						id="name"
						name="name"
						type="text"
						bind:value={name}
						on:keydown={onKeydown}
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						placeholder="Name of the middleware"
						required
					/>
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
