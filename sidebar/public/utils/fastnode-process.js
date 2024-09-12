'use strict'

const { platform } = require('os')

function getSupport() {
  switch(platform()) {
    case 'darwin': return require('./osx-process.js')
    case 'win32': return require('./windows-process.js')
    case 'linux': return require('./linux-process.js')
    default: return require('./mock-process.js')
  }
}

const system = getSupport()

const fastnodeProcess = {
  launchFastnode() {
    return this.stopFastnode()
      .then(() => this.startFastnode())
  },

  isFastnodeRunning() {
    return system.isFastnodeRunning()
  },

  startFastnode() {
    return system.startFastnode()
  },

  stopFastnode() {
    return system.stopFastnode()
  },
}

module.exports = {
  fastnodeProcess,
}