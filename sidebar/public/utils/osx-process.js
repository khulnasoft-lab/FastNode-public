'use strict'

const { spawnSync } = require('child_process')
const { spawnPromise } = require('./process-utils')

const FASTNODED_PATH = /Fastnode\.app\/Contents\/MacOS\/(Fastnode\s|Fastnode$)/
const XCODE = /\/Xcode\//
const ELECTRON = /\/electron\//

const memInstallPath = () => {
  let installPaths = []
  return () => {
    if(installPaths.length === 0) {
      //compute paths
      installPaths = String(spawnSync('mdfind', [
        'kMDItemCFBundleIdentifier = "com.fastnode.Fastnode"',
      ]).stdout)
        .trim()
        .split('\n')
        .filter(path => !XCODE.test(path)) //filter out development paths
    }
    return installPaths[0]
  }
}

const installPath = memInstallPath()

module.exports = {
  //resolves with a pid
  isFastnodeRunning() {
    return spawnPromise('/bin/ps', [
      '-axo', 'pid,command',
    ], {
      encoding: 'utf8',
    }, 'ps_error')
      .then(stdout => {
        const procs = stdout.split('\n')
        const fastnodeprocs = procs.filter(s => FASTNODED_PATH.test(s) 
                                            && !ELECTRON.test(s) //filter out dev electron builds
                                            && !s.includes('Fastnode.app/Contents/Resources')) //filter out production electron app
        if(fastnodeprocs.length > 0) {
          return {
            processes: fastnodeprocs,
            running: true,
          }
        } else {
          return {
            running: false,
          }
        }
      })
  },

  startFastnode() {
    return spawnPromise('open', [
      '-a', installPath(), '--args', '"--sidebar-restart"',
    ])
  },

  stopFastnode() {
    return this.isFastnodeRunning()
      .then((res) => {
        if(res.running) {
          //even if multiple Fastnode processes, for now, just assume 0-index
          const pid = res.processes[0].trim().split(" ")[0]
          return spawnPromise('/bin/kill', [pid], 'kill_error')
        } else {
          return Promise.resolve()
        }
      })
  },
}