import type { Component } from 'svelte';
import type { IconProps } from '@lucide/svelte';

type IconComponent = Component<IconProps, Record<string, never>, ''>;

export type BulkAction<T> = {
	type: 'button' | 'select';
	label: string;
	icon?: IconComponent;
	variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
	class?: string;
	disabled?: boolean;
	// For button type
	onClick?: (selectedItems: T[]) => Promise<void> | void;
	// For select type
	options?: {
		label: string;
		value: string;
		onClick: (selectedItems: T[], value: string) => Promise<void> | void;
	}[];
};
