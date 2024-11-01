<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="address" class="text-right">Address</Label>
	<Input
		id="address"
		name="address"
		type="text"
		bind:value={middleware.content.address}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="https://example.com/auth"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="trust-forward-header" class="text-right">Trust Forward Header</Label>
	<Switch
		id="trust-forward-header"
		bind:checked={middleware.content.trustForwardHeader}
		class="col-span-3"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="auth-response-headers-regex" class="text-right">Auth Response Headers Regex</Label>
	<Input
		id="auth-response-headers-regex"
		name="auth-response-headers-regex"
		type="text"
		bind:value={middleware.content.authResponseHeadersRegex}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="^X-"
		{disabled}
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="auth-request-headers" class="text-right">Auth Request Headers</Label>
	<Input
		id="auth-request-headers"
		name="auth-request-headers"
		type="text"
		bind:value={middleware.content.authRequestHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="X-CustomHeader"
		{disabled}
	/>
</div>
<ArrayInput
	bind:items={middleware.content.authResponseHeaders}
	label="Auth Response Headers"
	placeholder="X-Auth-User"
	{disabled}
/>
<ArrayInput
	bind:items={middleware.content.authRequestHeaders}
	label="Auth Request Headers"
	placeholder="Accept"
	{disabled}
/>
<ArrayInput
	bind:items={middleware.content.addAuthCookiesToResponse}
	label="Add Auth Cookies To Response"
	placeholder="Session-Cookie"
	{disabled}
/>
