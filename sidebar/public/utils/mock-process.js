'use strict'

/**
 * Dummy implementation for unsupported OSs.
 */

module.exports = {
  isFastnodeRunning() {
    return new Promise((resolve, reject) => {
      resolve({isRunning:false})
    })
  },

  startFastnode() {
    return new Promise((resolve, reject) => {
      resolve('MOCK')
    })
  },

  stopFastnode() {
    return new Promise((resolve, reject) => {
      resolve('MOCK')
    })
  },
}