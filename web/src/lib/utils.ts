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
// export function cleanupFormData(data: Record<string, unknown>|null): Record<string, unknown>|null {
//   if (data === null || typeof data !== 'object') {
//     return data;
//   }

//   // Handle arrays
//   if (Array.isArray(data)) {
//     // Filter out empty strings from arrays
//     const filtered = data.filter(item => item !== '');

//     // If the array is now empty, return null to remove it
//     if (filtered.length === 0) {
//       return null;
//     }

//     // Clean each item in the array
//     return filtered.map(item => cleanupFormData(item));
//   }

//   // Handle objects
//   const result = {};
//   let hasValue = false;

//   for (const [key, value] of Object.entries(data)) {
//     // Clean the value recursively
//     const cleanValue = cleanupFormData(value);

//     // Only include non-empty values
//     if (cleanValue !== null &&
//         cleanValue !== undefined &&
//         cleanValue !== '' &&
//         !(typeof cleanValue === 'object' && Object.keys(cleanValue).length === 0) &&
//         !(Array.isArray(cleanValue) && cleanValue.length === 0)) {

//       result[key] = cleanValue;
//       hasValue = true;
//     }
//   }

//   // If the object has no valid properties, return null to remove it
//   return hasValue ? result : null;
// }
