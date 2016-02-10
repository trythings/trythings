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
	primary: deepPurple['500'],
	accent: red.A400,
	accentLight: red.A200,

	canvas: grey['300'],
	card: white,
};

const text = {
	light: {
		primary: 'rgba(255, 255, 255, 1.00)',
		secondary: 'rgba(255, 255, 255, 0.70)',
		disabled: 'rgba(255, 255, 255, 0.50)',
		dividers: 'rgba(255, 255, 255, 0.12)',

		color: '#ffffff',
		opacity: {
			primary: '1.00',
			secondary: '0.70',
			disabled: '0.50',
			dividers: '0.12',
		},
	},
	dark: {
		primary: 'rgba(0, 0, 0, 0.87)',
		secondary: 'rgba(0, 0, 0, 0.54)',
		disabled: 'rgba(0, 0, 0, 0.38)',
		dividers: 'rgba(0, 0, 0, 0.12)',

		color: '#000000',
		opacity: {
			primary: '0.87',
			secondary: '0.54',
			disabled: '0.38',
			dividers: '0.12',
		},
	},
};

export default {
	colors,
	fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
	text,
};
