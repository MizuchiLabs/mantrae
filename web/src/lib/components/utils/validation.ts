// Some extra zod validations
import { z } from 'zod';

const ipv4Regex =
	/^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$/;
const ipv4CidrRegex = /^(3[0-2]|[12]?[0-9])$/;

const ipv6Regex =
	/^(([a-f0-9]{1,4}:){7}|::([a-f0-9]{1,4}:){0,6}|([a-f0-9]{1,4}:){1}:([a-f0-9]{1,4}:){0,5}|([a-f0-9]{1,4}:){2}:([a-f0-9]{1,4}:){0,4}|([a-f0-9]{1,4}:){3}:([a-f0-9]{1,4}:){0,3}|([a-f0-9]{1,4}:){4}:([a-f0-9]{1,4}:){0,2}|([a-f0-9]{1,4}:){5}:([a-f0-9]{1,4}:){0,1})([a-f0-9]{1,4}|(((25[0-5])|(2[0-4][0-9])|(1[0-9]{2})|([0-9]{1,2}))\.){3}((25[0-5])|(2[0-4][0-9])|(1[0-9]{2})|([0-9]{1,2})))$/;
const ipv6CidrRegex = /^(12[0-8]|1[01][0-9]|[1-9]?[0-9])$/;

const timeUnitRegex = /^(0|[1-9]\d*)(ns|us|µs|ms|s|m|h)?$/;

// Custom Zod schema to validate either IPv4 or IPv6 with or without CIDR
export const CustomIPSchema = z.string().refine(
	(value) => {
		const [ipAddress, mask] = value.split('/');
		if (!mask) {
			return ipv4Regex.test(ipAddress) || ipv6Regex.test(ipAddress);
		} else {
			return (
				(ipv4Regex.test(ipAddress) && ipv4CidrRegex.test(mask)) ||
				(ipv6Regex.test(ipAddress) && ipv6CidrRegex.test(mask))
			);
		}
	},
	{
		message: 'Invalid IP address or CIDR notation'
	}
);

export const CustomIPSchemaOptional = z.string().refine(
	(value) => {
		const [ipAddress, mask] = value.split('/');
		if (!ipAddress) return true;
		if (!mask) {
			return ipv4Regex.test(ipAddress) || ipv6Regex.test(ipAddress);
		} else {
			return (
				(ipv4Regex.test(ipAddress) && ipv4CidrRegex.test(mask)) ||
				(ipv6Regex.test(ipAddress) && ipv6CidrRegex.test(mask))
			);
		}
	},
	{
		message: 'Invalid IP address or CIDR notation'
	}
);

export const CustomTimeUnitSchema = z.string().refine((value) => timeUnitRegex.test(value), {
	message: 'Invalid time unit'
});

export const CustomTimeUnitSchemaOptional = z
	.string()
	.optional()
	.refine(
		(value) => {
			if (!value) return true;
			return timeUnitRegex.test(value);
		},
		{
			message: 'Invalid time unit'
		}
	);
