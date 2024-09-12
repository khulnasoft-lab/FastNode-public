import React from 'react'

import { track } from '../../utils/analytics'

import '../../assets/setup/start.css'

import FastnodeLogo from '../FastnodeLogo'

class Start extends React.Component {
  constructor(props) {
    super(props)
  }

  componentDidMount() {
    track({event: "onboarding_start_step_mounted"})
  }

  render() {
    return <div className="setup__start">
      <h2 className="setup__title showup__animation showup__animation--delay-2"> Welcome to Fastnode! </h2>
      <FastnodeLogo/>
      <p className="setup__text showup__animation showup__animation--delay-2">
        Fastnode helps you code smarter by plugging into your editor to show you code snippets and documentation related to what you’re working on.
      </p>
      <p className="setup__text showup__animation showup__animation--delay-2">
        To work properly, Fastnode just needs to set up a few editor plugins.
      </p>
      <button className="setup__button showup__animation showup__animation--delay-2" onClick={ this.props.advance }>
        Continue
      </button>
    </div>
  }
}

export default Start
