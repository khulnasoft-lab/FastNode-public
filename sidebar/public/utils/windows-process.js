'use strict'

const { execSync, spawn } = require('child_process')
const { spawnPromise } = require('./process-utils')
const { join } = require('path')

const KEY_BAT = `"${join(__dirname, 'read-key.bat')}"`
const FALLBACK_INSTALL_PATH = process.env.ProgramW6432
  ? join(process.env.ProgramW6432, 'Fastnode')
  : 'C:\\Program Files\\Fastnode'

const memInstallPath = () => {
  try {
    const registryPath = String(execSync(KEY_BAT)).trim()
    return () => {
      if (registryPath !== 'not found') return registryPath
      return FALLBACK_INSTALL_PATH
    }
  } catch (err) {
    console.error('error finding registry', err)
    return () => { return FALLBACK_INSTALL_PATH }
  }
}

const installPath = memInstallPath()
const FASTNODE_EXE_PATH = join(installPath(), 'fastnoded.exe')

module.exports = {
  isFastnodeRunning() {
    return spawnPromise('tasklist', 'tasklist_error')
      .then(stdout => {
        const procs = stdout.split('\n')
        const fastnodeprocs = procs.filter(proc => proc.indexOf("fastnoded.exe") !== -1)
        if (fastnodeprocs.length > 0) {
          return {
            processes: fastnodeprocs,
            running: true
          }
        }
        return { running: false }
      })
  },

  startFastnode() {
    var env = Object.create(process.env)
    env.FASTNODE_SKIP_ONBOARDING = '1'
    spawn(`${FASTNODE_EXE_PATH}`, ['--sidebar-restart'], { detached: true, env: env })
    return Promise.resolve()
  },

  stopFastnode() {
    return this.isFastnodeRunning()
      .then((res) => {
        if (res.running) {
          return spawnPromise('taskkill', ['/im', 'fastnoded.exe', '/f'], 'taskkill_err')
        } else {
          return Promise.resolve()
        }
      })
  },
}
