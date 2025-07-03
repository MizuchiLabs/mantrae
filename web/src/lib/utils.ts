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
