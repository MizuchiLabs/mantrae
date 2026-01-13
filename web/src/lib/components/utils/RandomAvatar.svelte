<script lang="ts">
	const SIZE = 36;

	interface Props {
		size?: number;
		name?: string;
		square?: boolean;
		colors?: string[];
	}
	let {
		size = 32,
		name = 'Clara Barton',
		square = false,
		colors = ['#92A1C6', '#146A7C', '#F0AB3D', '#C271B4', '#C20D90']
	}: Props = $props();

	const getNumber = (name: string): number => {
		return Array.from(name).reduce((sum, char) => sum + char.charCodeAt(0), 0);
	};

	const getDigit = (number: number, position: number): number => {
		return Math.floor((number / Math.pow(10, position)) % 10);
	};

	const getUnit = (number: number, range: number, index?: number): number => {
		const value = number % range;
		return index && getDigit(number, index) % 2 === 0 ? -value : value;
	};

	const getRandomColor = (number: number, colors: string[]): string => {
		return colors[number % colors.length];
	};

	const getContrast = (hexcolor: string): 'black' | 'white' => {
		const hex = hexcolor.startsWith('#') ? hexcolor.slice(1) : hexcolor;

		const r = parseInt(hex.slice(0, 2), 16);
		const g = parseInt(hex.slice(2, 4), 16);
		const b = parseInt(hex.slice(4, 6), 16);

		const yiq = (r * 299 + g * 587 + b * 114) / 1000;
		return yiq >= 128 ? 'black' : 'white';
	};

	const maskId = `mask__beam_${Math.random().toString(36).slice(2)}`;

	const generateData = (name: string, colors: string[]) => {
		const numFromName = getNumber(name);
		const wrapperColor = getRandomColor(numFromName, colors);

		const preTranslateX = getUnit(numFromName, 10, 1);
		const wrapperTranslateX = preTranslateX < 5 ? preTranslateX + SIZE / 9 : preTranslateX;

		const preTranslateY = getUnit(numFromName, 10, 2);
		const wrapperTranslateY = preTranslateY < 5 ? preTranslateY + SIZE / 9 : preTranslateY;

		return {
			wrapperColor,
			faceColor: getContrast(wrapperColor),
			backgroundColor: getRandomColor(numFromName + 13, colors),
			wrapperTranslateX,
			wrapperTranslateY,
			wrapperRotate: getUnit(numFromName, 360),
			wrapperScale: 1 + getUnit(numFromName, SIZE / 12) / 10,
			isMouthOpen: getDigit(numFromName, 2) % 2 === 0,
			isCircle: getDigit(numFromName, 1) % 2 === 0,
			eyeSpread: getUnit(numFromName, 5),
			mouthSpread: getUnit(numFromName, 3),
			faceRotate: getUnit(numFromName, 10, 3),
			faceTranslateX:
				wrapperTranslateX > SIZE / 6 ? wrapperTranslateX / 2 : getUnit(numFromName, 8, 1),
			faceTranslateY:
				wrapperTranslateY > SIZE / 6 ? wrapperTranslateY / 2 : getUnit(numFromName, 7, 2)
		};
	};

	let data = $derived(generateData(name, colors));
</script>

<svg
	viewBox="0 0 {SIZE} {SIZE}"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
	width={size}
	height={size}
	data-testid="avatar-beam"
>
	<mask id={maskId} maskUnits="userSpaceOnUse" x={0} y={0} width={SIZE} height={SIZE}>
		<rect width={SIZE} height={SIZE} rx={square ? undefined : SIZE * 2} fill="white" />
	</mask>
	<g mask="url(#{maskId})">
		<rect width={SIZE} height={SIZE} fill={data.backgroundColor} />
		<rect
			x="0"
			y="0"
			width={SIZE}
			height={SIZE}
			transform="translate({data.wrapperTranslateX} {data.wrapperTranslateY}) rotate({data.wrapperRotate} {SIZE /
				2} {SIZE / 2}) scale({data.wrapperScale})"
			fill={data.wrapperColor}
			rx={data.isCircle ? SIZE : SIZE / 6}
		/>
		<g
			transform="translate({data.faceTranslateX} {data.faceTranslateY}) rotate({data.faceRotate} {SIZE /
				2} {SIZE / 2})"
		>
			{#if data.isMouthOpen}
				<path
					d="M15 {19 + data.mouthSpread}c2 1 4 1 6 0"
					stroke={data.faceColor}
					fill="none"
					stroke-linecap="round"
				/>
			{:else}
				<path d="M13,{19 + data.mouthSpread} a1,0.75 0 0,0 10,0" fill={data.faceColor} />
			{/if}

			<rect
				x={14 - data.eyeSpread}
				y={14}
				width={1.5}
				height={2}
				rx={1}
				stroke="none"
				fill={data.faceColor}
			/>
			<rect
				x={20 + data.eyeSpread}
				y={14}
				width={1.5}
				height={2}
				rx={1}
				stroke="none"
				fill={data.faceColor}
			/>
		</g>
	</g>
</svg>
