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
