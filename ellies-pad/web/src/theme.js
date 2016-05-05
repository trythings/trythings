import './Roboto.css';

// const deepPurple = {
// 	50: '#ede7f6',
// 	200: '#b39ddb',
// 	300: '#9575cd',
// 	400: '#7e57c2',
// 	500: '#673ab7',
// };

const indigo = {
	100: '#c5cae9',
	200: '#9fa8da',
	300: '#7986cb',
	400: '#5c6bc0',
};

const red = {
	A200: '#ff5252',
	A400: '#ff1744',
};

const grey = {
	100: '#f5f5f5',
	300: '#e0e0e0',
};

const white = '#ffffff';

const colors = {
	primary: {
		default: indigo[400],
		light: indigo[300],
		xlight: indigo[200],
	},
	accent: red.A200,

	canvas: grey[100],
	card: white,
	dividers: {
		dark: 'rgba(0, 0, 0, 0.12)',
		light: 'rgba(255, 255, 255, 0.12)',
	},
};

const elevation = {
	2: {
		boxShadow: [
			'0 1px 5px 0 rgba(0, 0, 0, 0.12)', // Ambient.
			'0 2px 2px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 3px 1px -2px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 2,
	},
	4: {
		boxShadow: [
			'0 1px 10px 0 rgba(0, 0, 0, 0.12)', // Ambient.
			'0 4px 5px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 2px 4px -1px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 4,
	},
	6: {
		boxShadow: [
			'0 1px 18px 0 rgba(0, 0, 0, 0.12)', // Ambient.
			'0 6px 10px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 3px 5px -1px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 6,
	},
	8: {
		boxShadow: [
			'0 3px 14px 2px rgba(0, 0, 0, 0.12)', // Ambient.
			'0 8px 10px 1px rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 5px 5px -3px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 8,
	},
	10: {
		boxShadow: [
			'0 4px 18px 3px rgba(0, 0, 0, 0.12)', // Ambient.
			'0 10px 14px 1px rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 6px 6px -3px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 10,
	},
	12: {
		boxShadow: [
			'0 5px 22px 4px rgba(0, 0, 0, 0.12)', // Ambient.
			'0 12px 17px 2px rgba(0, 0, 0, 0.14)', // Penumbra.
			'0 7px 8px -4px rgba(0, 0, 0, 0.20)', // Umbra.
		].join(','),
		zIndex: 12,
	},
};

const text = {
	fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
	light: {
		primary: {
			color: 'rgba(255, 255, 255, 1.00)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
		secondary: {
			color: 'rgba(255, 255, 255, 0.70)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
		disabled: {
			color: 'rgba(255, 255, 255, 0.50)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
	},
	dark: {
		primary: {
			color: 'rgba(0, 0, 0, 0.87)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
		secondary: {
			color: 'rgba(0, 0, 0, 0.54)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
		disabled: {
			color: 'rgba(0, 0, 0, 0.38)',
			fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
		},
	},
};

export default {
	colors,
	elevation,
	text,
};
