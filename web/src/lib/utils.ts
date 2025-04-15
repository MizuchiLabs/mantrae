import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function safeClone<T>(obj: T): T {
	try {
		return JSON.parse(JSON.stringify(obj));
	} catch (e) {
		console.warn('Failed to clone object:', e);
		return obj;
	}
}

export async function tryLoad<T>(fn: () => Promise<T>, fallback: T): Promise<T> {
	try {
		return await fn();
	} catch (err: unknown) {
		const error = err instanceof Error ? err.message : String(err);
		console.warn('Failed to load:', error);
		return fallback;
	}
}

export function cleanupFormData(data: unknown): unknown {
	// Base cases for non-objects
	if (data === null || data === undefined || data === '') {
		return null;
	}

	if (typeof data !== 'object') {
		// Keep primitive values like numbers and booleans
		return data;
	}

	// Handle arrays
	if (Array.isArray(data)) {
		// Filter out empty values and clean each remaining item
		const filtered = data
			.filter((item) => item !== null && item !== undefined && item !== '')
			.map((item) => cleanupFormData(item))
			.filter((item) => item !== null);

		return filtered.length > 0 ? filtered : null;
	}

	// Handle objects
	const result: Record<string, unknown> = {};
	let hasValidProperty = false;

	for (const [key, value] of Object.entries(data)) {
		// Skip default values for specific properties
		if (
			(key === 'depth' && value === 0) ||
			(key === 'requestHost' && value === false) ||
			(key === 'excludedIPs' && Array.isArray(value) && value.length === 0)
		) {
			continue;
		}

		const cleanValue = cleanupFormData(value);

		// Only include meaningful values
		if (cleanValue !== null) {
			result[key] = cleanValue;
			hasValidProperty = true;
		}
	}

	// Return null for empty objects to remove them entirely
	return hasValidProperty ? result : null;
}
