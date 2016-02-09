const deepPurple = {
	50: '#ede7f6',
	500: '#673ab7',
};

const red = {
	A400: '#ff1744',
};

const grey = {
	300: '#e0e0e0',
};

const colors = {
	primary1: deepPurple['500'],
	// primary2,
	// primary3,
	accent1: red.A400,
	canvas: grey['300'],
};

const text = {
	light: {
		color: '#ffffff',
		opacity: {
			primary: '1.00',
			secondary: '0.70',
			disabled: '0.50',
			dividers: '0.12',
		},
	},
	dark: {
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
	text,
};
