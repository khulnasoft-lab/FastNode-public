import React from 'react'
import { ENTERPRISE } from '../utils/enterprise'

const TopBar = () => {
  return <div className="header">
    { ENTERPRISE
      ? <div className="fastnode-enterprise-logo"/>
      : <div className="fastnode-logo"/>
    }
  </div>
}

export default TopBar
