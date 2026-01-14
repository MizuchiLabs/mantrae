import {
	durationFromMs,
	durationMs,
	timestampDate,
	timestampFromDate,
	type Duration,
	type Timestamp
} from '@bufbuild/protobuf/wkt';
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & {
	ref?: U | null;
};

// Helper function to truncate text with ellipsis
export function truncateText(text: string, maxLength: number = 30): string {
	return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
}

// Helper function to format array display with ellipsis
export function formatArrayDisplay(arr: string[] | undefined, maxItems: number = 2): string {
	if (!arr || arr.length === 0) return '';
	if (arr.length <= maxItems) return arr.join(', ');
	return `${arr.slice(0, maxItems).join(', ')} (+${arr.length - maxItems})`;
}

export function parseGoDuration(input: string) {
	const normalized = input
		.toLowerCase()
		.replace(/minutes?|mins?/g, 'm')
		.replace(/seconds?|secs?/g, 's')
		.replace(/hours?|hrs?/g, 'h')
		.replace(/\s+/g, '');

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
	return parts.join('');
}

// =============================================================================
// Timestamps
// =============================================================================

/** Convert Timestamp to Date (passthrough for convenience) */
export const toDate = (ts: Timestamp | null | undefined): Date =>
	ts ? timestampDate(ts) : new Date();

/** Create Timestamp from Date */
export const fromDate = (date: Date): Timestamp => timestampFromDate(date);

/** Create Timestamp for current time */
export const now = (): Timestamp => fromDate(new Date());

type DateStyle = 'full' | 'long' | 'short' | 'compact' | 'relative' | 'relative-short';

const formatters = {
	full: new Intl.DateTimeFormat(undefined, {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
		second: 'numeric'
	}),
	long: new Intl.DateTimeFormat(undefined, {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	}),
	short: new Intl.DateTimeFormat(undefined, {
		month: 'short',
		day: 'numeric',
		year: 'numeric',
		hour: 'numeric',
		minute: '2-digit'
	}),
	compact: new Intl.DateTimeFormat(undefined, {
		month: 'numeric',
		day: 'numeric',
		year: '2-digit'
	})
};

/** Format a Timestamp for display */
export function formatTs(
	ts: Timestamp | null | undefined,
	style: DateStyle = 'short',
	fallback = ''
): string {
	if (!ts) return fallback;

	if (style === 'relative' || style === 'relative-short') {
		return formatRelative(ts, style === 'relative-short');
	}

	return formatters[style].format(timestampDate(ts));
}

function formatRelative(ts: Timestamp, short = false): string {
	const diffMs = timestampDate(ts).getTime() - Date.now();
	const absDiff = Math.abs(diffMs);
	const isFuture = diffMs > 0;

	if (absDiff < 10_000) return 'just now';

	const sec = Math.floor(absDiff / 1000);
	const min = Math.floor(sec / 60);
	const hr = Math.floor(min / 60);
	const day = Math.floor(hr / 24);
	const month = Math.floor(day / 30);
	const year = Math.floor(day / 365);

	let str: string;
	if (sec < 60) str = short ? `${sec}s` : `${sec} second${sec !== 1 ? 's' : ''}`;
	else if (min < 60) str = short ? `${min}m` : `${min} minute${min !== 1 ? 's' : ''}`;
	else if (hr < 24) str = short ? `${hr}h` : `${hr} hour${hr !== 1 ? 's' : ''}`;
	else if (day < 30) str = short ? `${day}d` : `${day} day${day !== 1 ? 's' : ''}`;
	else if (month < 12) str = short ? `${month}mo` : `${month} month${month !== 1 ? 's' : ''}`;
	else str = short ? `${year}y` : `${year} year${year !== 1 ? 's' : ''}`;

	return isFuture ? `in ${str}` : `${str} ago`;
}

// =============================================================================
// Durations
// =============================================================================

/** Convert milliseconds to Duration */
export const msToD = (ms: number): Duration => durationFromMs(ms);

/** Convert seconds to Duration */
export const secToD = (sec: number): Duration => msToD(sec * 1000);

/** Convert Duration to milliseconds */
export const dToMs = (d: Duration | null | undefined): number => (d ? durationMs(d) : 0);

/** Convert Duration to seconds */
export const dToSec = (d: Duration | null | undefined): number => dToMs(d) / 1000;

type DurationStyle = 'short' | 'long' | 'compact' | 'clock';

const UNITS = [
	{ s: 'y', l: 'year', ms: 365 * 24 * 60 * 60 * 1000 },
	{ s: 'mo', l: 'month', ms: 30 * 24 * 60 * 60 * 1000 },
	{ s: 'w', l: 'week', ms: 7 * 24 * 60 * 60 * 1000 },
	{ s: 'd', l: 'day', ms: 24 * 60 * 60 * 1000 },
	{ s: 'h', l: 'hour', ms: 60 * 60 * 1000 },
	{ s: 'm', l: 'minute', ms: 60 * 1000 },
	{ s: 's', l: 'second', ms: 1000 }
] as const;

/** Format a Duration for display */
export function formatD(
	d: Duration | null | undefined,
	style: DurationStyle = 'short',
	maxUnits = 2
): string {
	const ms = dToMs(d);
	if (ms === 0) return style === 'long' ? '0 seconds' : '0s';

	const neg = ms < 0;
	let remaining = Math.abs(ms);

	if (style === 'clock') {
		const h = Math.floor(remaining / 3600000);
		const m = Math.floor((remaining % 3600000) / 60000);
		const s = Math.floor((remaining % 60000) / 1000);
		const result = h > 0 ? `${h}:${pad(m)}:${pad(s)}` : `${m}:${pad(s)}`;
		return neg ? `-${result}` : result;
	}

	const parts: string[] = [];
	for (const u of UNITS) {
		if (parts.length >= maxUnits) break;
		const count = Math.floor(remaining / u.ms);
		if (count > 0) {
			remaining %= u.ms;
			parts.push(style === 'long' ? `${count} ${u.l}${count !== 1 ? 's' : ''}` : `${count}${u.s}`);
		}
	}

	if (parts.length === 0) return style === 'long' ? '0 seconds' : '0s';

	const sep = style === 'compact' ? '' : style === 'long' ? ', ' : ' ';
	const result = parts.join(sep);
	return neg ? `-${result}` : result;
}

const pad = (n: number) => n.toString().padStart(2, '0');

// =============================================================================
// Enums
// =============================================================================

/** Convert protobuf enum to array (excludes 0/UNSPECIFIED) */
export function enumToArray<T extends object>(e: T): { name: string; value: number }[] {
	return Object.entries(e)
		.filter(([, v]) => typeof v === 'number' && v !== 0)
		.map(([name, value]) => ({ name, value: value as number }));
}

/** Convert SCREAMING_SNAKE_CASE enum to Title Case */
export function enumLabel(value: string | undefined): string {
	if (!value) return '';
	return value
		.toLowerCase()
		.split('_')
		.map((w) => w.charAt(0).toUpperCase() + w.slice(1))
		.join(' ');
}
