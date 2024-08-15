<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: HttpMiddleware;
	middleware.circuitBreaker = middleware.circuitBreaker ?? {
		expression: '',
		checkPeriod: '',
		fallbackDuration: '',
		recoveryDuration: ''
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="expression" class="text-right">Expression</Label>
	<Input
		id="expression"
		name="expression"
		type="text"
		bind:value={middleware.circuitBreaker.expression}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Expression"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="check-period" class="text-right">Check Period</Label>
	<Input
		id="check-period"
		name="check-period"
		type="text"
		bind:value={middleware.circuitBreaker.checkPeriod}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Check Period"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="fallback-duration" class="text-right">Fallback Duration</Label>
	<Input
		id="fallback-duration"
		name="fallback-duration"
		type="text"
		bind:value={middleware.circuitBreaker.fallbackDuration}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Fallback Duration"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="recovery-duration" class="text-right">Recovery Duration</Label>
	<Input
		id="recovery-duration"
		name="recovery-duration"
		type="text"
		bind:value={middleware.circuitBreaker.recoveryDuration}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Recovery Duration"
	/>
</div>
