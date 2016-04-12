import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskCard from './AddTaskCard.js';
import resetStyles from './resetStyles.js';

class AddTask extends React.Component {
	static propTypes = {
		space: React.PropTypes.shape({
			// ...AddTaskCard.propTypes.space,
		}).isRequired,
	};

	static styles = {
		addTaskButton: {
			...resetStyles,

			overflow: 'visible',
			paddingRight: 24,
			paddingBottom: 24,
			position: 'absolute',
			right: 0,
			bottom: 0,
		},
	};

	state = {
		isAddTaskFormVisible: false,
	};

	onCancelClick = () => {
		this.setState({ isAddTaskFormVisible: false });
	};

	onPlusClick = () => {
		this.setState({ isAddTaskFormVisible: true });
	};

	render() {
		if (this.state.isAddTaskFormVisible) {
			return (
				<AddTaskCard
					autoFocus
					space={this.props.space}
					onCancelClick={this.onCancelClick}
				/>
			);
		}

		return (
			<div style={AddTask.styles.addTaskButton}>
				<ActionButton onClick={this.onPlusClick} />
			</div>
		);
	}
}

export default Relay.createContainer(AddTask, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				${AddTaskCard.getFragment('space')},
			},
		`,
	},
});
