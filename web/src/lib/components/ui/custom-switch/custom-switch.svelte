<script lang="ts">
	import { Check, X } from '@lucide/svelte';

	interface props {
		checked?: boolean;
		disabled?: boolean;
		size?: 'sm' | 'md' | 'lg';
		variant?: 'default' | 'text';
		textLabels?: { checked: string; unchecked: string };
		onCheckedChange?: (checked: boolean) => void;
	}

	let {
		checked = $bindable(),
		disabled = false,
		size = 'md',
		onCheckedChange,
		variant = 'default',
		textLabels
	}: props = $props();

	const sizeClasses = {
		sm: variant === 'text' ? 'h-5' : 'h-5 w-9',
		md: variant === 'text' ? 'h-6' : 'h-6 w-11',
		lg: variant === 'text' ? 'h-8' : 'h-8 w-14'
	};

	const thumbSizeClasses = {
		sm: variant === 'text' ? 'h-4 px-2 min-w-8' : 'h-4 w-4',
		md: variant === 'text' ? 'h-5 px-2 min-w-10' : 'h-5 w-5',
		lg: variant === 'text' ? 'h-7 px-3 min-w-12' : 'h-7 w-7'
	};

	const iconSizes = {
		sm: 10,
		md: 12,
		lg: 16
	};

	// Calculate approximate width needed for text variant
	const getTextWidth = () => {
		if (!textLabels) return { checked: 0, unchecked: 0 };
		const maxLength = Math.max(textLabels.checked.length, textLabels.unchecked.length);
		return {
			sm: Math.max(maxLength * 6 + 64),
			md: Math.max(maxLength * 7 + 80),
			lg: Math.max(maxLength * 8 + 96)
		}[size];
	};

	// Dynamic translation based on variant and container width
	const getTranslation = () => {
		if (variant === 'text') {
			return checked ? 'translate-x-[calc(100%-0.125rem)]' : 'translate-x-0.5';
		}
		switch (size) {
			case 'sm':
				return checked ? 'translate-x-4' : 'translate-x-0.5';
			case 'md':
				return checked ? 'translate-x-5.5' : 'translate-x-0.5';
			case 'lg':
				return checked ? 'translate-x-6' : 'translate-x-0.5';
			default:
				return checked ? 'translate-x-5.5' : 'translate-x-0.5';
		}
	};

	const handletoggle = () => {
		if (disabled) return;
		checked = !checked;
		onCheckedChange?.(checked);
	};
</script>

<button
	type="button"
	role="switch"
	aria-checked={checked}
	{disabled}
	onclick={handletoggle}
	class="focus:ring-primary relative inline-flex items-center rounded-full transition-colors duration-200 ease-in-out focus:outline-none
		{sizeClasses[size]}
		{checked ? 'bg-primary' : 'bg-muted'}
		{disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}
	"
	style={variant === 'text' && textLabels ? `width: ${getTextWidth()}px;` : ''}
	aria-label="Toggle switch"
>
	{#if variant === 'text' && textLabels}
		<!-- Background text labels -->
		<div class="absolute inset-0 flex items-center justify-between px-2">
			<span class="text-muted-foreground text-xs font-medium transition-colors duration-200">
				{textLabels.unchecked}
			</span>
			<span class="text-muted-foreground text-xs font-medium transition-colors duration-200">
				{textLabels.checked}
			</span>
		</div>
	{/if}

	<span
		class="pointer-events-none flex items-center justify-center rounded-full shadow-lg ring-0 transition-all duration-200 ease-in-out
			{thumbSizeClasses[size]}
			{checked ? 'bg-primary-foreground' : 'bg-background border-border/20 border'}
			{getTranslation()}
		"
	>
		{#if variant === 'default'}
			<div class="transition-opacity duration-150 {checked ? 'opacity-100' : 'opacity-0'}">
				{#if checked}
					<Check size={iconSizes[size]} class="stroke-[2.5] text-green-500" />
				{/if}
			</div>
			<div class="absolute transition-opacity duration-150 {checked ? 'opacity-0' : 'opacity-100'}">
				{#if !checked}
					<X size={iconSizes[size]} class="stroke-[2.5] text-red-500" />
				{/if}
			</div>
		{:else if variant === 'text' && textLabels}
			<span class="text-primary text-xs font-bold whitespace-nowrap">
				{checked ? textLabels.checked : textLabels.unchecked}
			</span>
		{/if}
	</span>
</button>
