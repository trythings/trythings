const noInherit = {
	// all: 'initial',

	azimuth: 'initial',
	borderCollapse: 'initial',
	borderSpacing: 'initial',
	captionSide: 'initial',
	color: 'initial',
	// cursor: 'initial',
	direction: 'initial',
	elevation: 'initial',
	emptyCells: 'initial',
	fontFamily: 'initial',
	fontSize: 'initial',
	fontStyle: 'initial',
	fontVariant: 'initial',
	fontWeight: 'initial',
	// font: 'initial',
	letterSpacing: 'initial',
	lineHeight: 'initial',
	listStyleImage: 'initial',
	listStylePosition: 'initial',
	listStyleType: 'initial',
	// listStyle: 'initial',
	orphans: 'initial',
	pitchRange: 'initial',
	pitch: 'initial',
	quotes: 'initial',
	richness: 'initial',
	speakHeader: 'initial',
	speakNumeral: 'initial',
	speakPunctuation: 'initial',
	// speak: 'initial',
	speechRate: 'initial',
	stress: 'initial',
	textAlign: 'initial',
	textIndent: 'initial',
	textTransform: 'initial',
	visibility: 'initial',
	voiceFamily: 'initial',
	volume: 'initial',
	whiteSpace: 'initial',
	widows: 'initial',
	wordSpacing: 'initial',
};

const defaults = {
	backgroundColor: 'transparent',

	// border* can actually be broken down further.
	borderBottom: 'none',
	borderLeft: 'none',
	borderRight: 'none',
	borderTop: 'none',

	boxSizing: 'border-box',
	display: 'flex',

	marginBottom: 0,
	marginLeft: 0,
	marginRight: 0,
	marginTop: 0,

	outline: 0,
	overflow: 'hidden',

	paddingBottom: 0,
	paddingLeft: 0,
	paddingRight: 0,
	paddingTop: 0,

	position: 'relative',
	resize: 'none',
};

export default {
	...noInherit,
	...defaults,
};
