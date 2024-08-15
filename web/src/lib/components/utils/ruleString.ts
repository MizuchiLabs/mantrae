// Regex patterns for traefik rules

function validateSingleRule(rule: string): boolean {
	const rulePatterns = {
		Header: /^!?Header\(`[^`]+`,\s*`[^`]+`\)$/,
		HeaderRegexp: /^!?HeaderRegexp\(`[^`]+`,\s*`[^`]+`\)$/,
		Host: /^!?Host\(`[^`]+`\)$/,
		HostRegexp: /^!?HostRegexp\(`[^`]+`\)$/,
		Method: /^!?Method\(`[^`]+`\)$/,
		Path: /^!?Path\(`[^`]+`\)$/,
		PathPrefix: /^!?PathPrefix\(`[^`]+`\)$/,
		PathRegexp: /^!?PathRegexp\(`[^`]+`\)$/,
		Query: /^!?Query\(`[^`]+`,\s*`[^`]+`\)$/,
		QueryRegexp: /^!?QueryRegexp\(`[^`]+`,\s*`[^`]+`\)$/,
		ClientIP: /^!?ClientIP\(`[^`]+`\)$/
	};

	for (const pattern of Object.values(rulePatterns)) {
		if (pattern.test(rule.trim())) {
			return true;
		}
	}
	return false;
}

export function ValidateRule(rule: string | undefined): boolean {
	if (rule === '' || rule === undefined) return true;
	const ruleParts = rule.split(/&&|\|\|/);
	return ruleParts.every(validateSingleRule);
}

export function RuleDescription(rules: string) {
	if (rules === '' || rules === undefined) return 'No specific routing rules applied';
	let description = '';
	const hostRule = '';
	const pathRules = [];
	const otherConditions = [];

	const rulePatterns = {
		Host: /Host\(`(.*?)`\)/,
		HostRegexp: /HostRegexp\(`(.*?)`\)/,
		Path: /Path\(`(.*?)`\)/,
		PathPrefix: /PathPrefix\(`(.*?)`\)/,
		PathRegexp: /PathRegexp\(`(.*?)`\)/,
		Method: /Method\(`(.*?)`\)/,
		Query: /Query\(`(.*?)`, `(.*?)`\)/,
		QueryRegexp: /QueryRegexp\(`(.*?)`, `(.*?)`\)/,
		ClientIP: /ClientIP\(`(.*?)`\)/,
		Header: /Header\(`(.*?)`, `(.*?)`\)/,
		HeaderRegexp: /HeaderRegexp\(`(.*?)`, `(.*?)`\)/
	};

	const conditions = rules.split(/(&&|\|\|)/);

	// Handle specific combinations first
	conditions.forEach((condition, index) => {
		let formattedCondition = condition.trim();

		// Check for negation
		const isNegated = formattedCondition.startsWith('!');
		if (isNegated) {
			formattedCondition = formattedCondition.substring(1).trim();
		}

		// Match against known patterns
		for (const [type, pattern] of Object.entries(rulePatterns)) {
			const match = formattedCondition.match(pattern);
			if (match) {
				let result;
				switch (type) {
					case 'Host':
					case 'HostRegexp':
						result = match[1] + (type === 'HostRegexp' ? '*' : '');
						break;
					case 'Path':
					case 'PathPrefix':
					case 'PathRegexp':
						result = ` Path: ${match[1]}` + (type === 'PathRegexp' ? '*' : '');
						break;
					case 'Method':
						result = '';
						description = `[${match[1]}] ${description}`; // Add at the start of the description
						break;
					case 'Query':
					case 'QueryRegexp':
						result = `?${match[1]}=${match[2]}` + (type === 'QueryRegexp' ? '*' : '');
						break;
					case 'ClientIP':
						result = `from ${match[1]}`;
						break;
					case 'Header':
					case 'HeaderRegexp':
						result = `Header: ${match[1]}=${match[2]}` + (type === 'HeaderRegexp' ? '*' : '');
						break;
					default:
						result = formattedCondition;
				}

				if (isNegated) {
					result = `not ${result}`;
				}

				description += result;
				break;
			}
		}
	});

	// Handle the case when no specific rule is matched
	if (description === '') {
		description = 'No specific routing rules applied';
	}

	return description;
}
