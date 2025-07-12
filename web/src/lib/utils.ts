import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any }
	? Omit<T, "children">
	: T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & {
	ref?: U | null;
};

// Helper function to truncate text with ellipsis
export function truncateText(text: string, maxLength: number = 30): string {
	return text.length > maxLength ? text.substring(0, maxLength) + "..." : text;
}

// Helper function to format array display with ellipsis
export function formatArrayDisplay(
	arr: string[] | undefined,
	maxItems: number = 2,
): string {
	if (!arr || arr.length === 0) return "";
	if (arr.length <= maxItems) return arr.join(", ");
	return `${arr.slice(0, maxItems).join(", ")} (+${arr.length - maxItems})`;
}

export function parseGoDuration(input: string) {
	const normalized = input
		.toLowerCase()
		.replace(/minutes?|mins?/g, "m")
		.replace(/seconds?|secs?/g, "s")
		.replace(/hours?|hrs?/g, "h")
		.replace(/\s+/g, "");

	const regex = /(?<value>[-+]?\d+)(?<unit>ns|us|µs|ms|s|m|h)/g;
	const parts: string[] = [];
	let match;

	while ((match = regex.exec(normalized)) !== null) {
		const groups = match.groups as { value: string; unit: string } | undefined;
		if (!groups) return null;

		const { value, unit } = groups;
		parts.push(`${parseInt(value)}${unit}`);
	}

	if (parts.length === 0) return null;
	return parts.join("");
}
