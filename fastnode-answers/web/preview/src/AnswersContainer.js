import React from "react";

import AnswersPage from "@khulnasoft-lab/fastnode-answers-renderer";

class AnswersContainer extends React.Component {
  render() {
    if (this.props.input && this.props.input.content) {
      return <AnswersPage source={this.props.input} />;
    }
    return null;
  }
}

export default AnswersContainer;
