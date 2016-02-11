import './Roboto.css';

const deepPurple = {
	50: '#ede7f6',
	500: '#673ab7',
};

const red = {
	A200: '#ff5252',
	A400: '#ff1744',
};

const grey = {
	300: '#e0e0e0',
};

const white = '#ffffff';

const colors = {
	primary: deepPurple[500],
	accent: red.A400,
	accentLight: red.A200,

	canvas: grey[300],
	card: white,
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
		dividers: {
			color: 'rgba(255, 255, 255, 0.12)',
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
		dividers: {
			color: 'rgba(0, 0, 0, 0.12)',
		},
	},
};

export default {
	colors,
	elevation: {
		2: {
			boxShadow: [
				'0 1px 5px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 2px 2px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 1px -2px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			zIndex: 2,
		},
		6: {
			boxShadow: [
				'0 1px 18px 0 rgba(0, 0, 0, 0.12)', // Ambient.
				'0 6px 10px 0 rgba(0, 0, 0, 0.14)', // Penumbra.
				'0 3px 5px -1px rgba(0, 0, 0, 0.20)', // Umbra.
			].join(','),
			zIndex: 6,
		},
	},
	text,
};
