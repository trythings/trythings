import React from 'react';
import Relay from 'react-relay';

import ActionButton from './ActionButton.js';
import AddTaskMole from './AddTaskMole.js';
import resetStyles from './resetStyles.js';

class AddTask extends React.Component {
	static propTypes = {
		refetch: React.PropTypes.func.isRequired,
		space: React.PropTypes.shape({
			// ...AddTaskMole.propTypes.space,
		}).isRequired,
		style: React.PropTypes.shape({
			flex: React.PropTypes.string,
			pointerEvents: React.PropTypes.string,
		}),
	};

	static styles = {
		container: {
			...resetStyles,

			justifyContent: 'flex-end',
			maxHeight: 360,
			flexDirection: 'column',
			overflow: 'visible',
		},
		addTaskMoleContainer: {
			...resetStyles,

			flex: '1 1 auto',
			flexDirection: 'column',
			paddingRight: 16,
			overflow: 'visible',
		},
		addTaskButtonContainer: {
			...resetStyles,

			paddingBottom: 24,
			paddingRight: 24,
			overflow: 'visible',
		},
		addTaskMole: {
			...resetStyles,
			flex: '1 1 auto',
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
		let style = AddTask.styles.container;
		if (this.props.style && this.props.style.flex) {
			style = {
				...style,
				flex: this.props.style.flex,
			};
		}

		if (this.props.style && this.props.style.pointerEvents) {
			style = {
				...style,
				pointerEvents: this.props.style.pointerEvents,
			};
		}

		return (
			<div style={style}>
				{this.state.isAddTaskFormVisible ?
					(
						<div style={AddTask.styles.addTaskMoleContainer}>
							<AddTaskMole
								autoFocus
								refetch={this.props.refetch}
								space={this.props.space}
								onCancelClick={this.onCancelClick}
								style={AddTask.styles.addTaskMole}
							/>
						</div>
					) :
					(
						<div style={AddTask.styles.addTaskButtonContainer}>
							<ActionButton onClick={this.onPlusClick} />
						</div>
					)
				}
			</div>
		);
	}
}

export default Relay.createContainer(AddTask, {
	fragments: {
		space: () => Relay.QL`
			fragment on Space {
				${AddTaskMole.getFragment('space')},
			},
		`,
	},
});
