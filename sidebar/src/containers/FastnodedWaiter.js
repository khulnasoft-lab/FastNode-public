import React from 'react'
import { connect } from 'react-redux'
import ErrorOverlay from '../components/ErrorOverlay'
import * as system from '../actions/system'
import * as account from '../actions/account'
import * as settings from '../actions/settings'

/**
 * FastnodedWaiter wraps children and waits to render them
 * until after fastnoded signals that it is ready via polling
 * through getFastnodedReady.
 * After getting that signal, it stops polling
 */
class FastnodedWaiter extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      ready: this.props.fastnodedReady && typeof this.props.os !== 'undefined'
    }
  }

  UNSAFE_componentWillMount() {
    const {
      getFastnodedReady,
      getSystemInfo
    } = this.props
    getSystemInfo().then(() => {
      this.poller = setInterval(() => {
        getFastnodedReady()
      }, 300)
    })
  }

  componentDidUpdate() {
    const {
      fastnodedReady,
      getDefaultTheme,
      checkIfOnline,
      getUser,
    } = this.props
    if(fastnodedReady !== this.state.ready) {
      clearInterval(this.poller)
      // fetch env variables from fastnoded since it's now ready
      Promise.all([
        getDefaultTheme(),
        getUser(),
        checkIfOnline(),
      ]).then((res) => {
        setTimeout(() => {
          this.setState({ ready: fastnodedReady })
        }, 200)
      })
    }
  }

  render() {
    if(!this.state.ready) {
      //render some overlay
      return <ErrorOverlay
        title="Initializing..."
        spinner={true}
      />
    }
    return this.props.children
  }
}

const mapStateToProps = (state, ownProps) => ({
  ...ownProps,
  fastnodedReady: state.system.fastnodedReady,
  os: state.system.os,
})

const mapDispatchToProps = dispatch => ({
  getFastnodedReady: () => dispatch(system.getFastnodedReady()),
  getUser: () => dispatch(account.getUser()),
  getDefaultTheme: () => dispatch(settings.getDefaultTheme()),
  checkIfOnline: () => dispatch(system.checkIfOnline()),
  getSystemInfo: () => dispatch(system.getSystemInfo()),
})

export default connect(mapStateToProps, mapDispatchToProps)(FastnodedWaiter)
