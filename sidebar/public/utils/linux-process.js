'use strict'

const path = require('path')
const fs = require('fs')
const os = require('os')
const {spawn} = require('child_process')
const {spawnPromise} = require('./process-utils')

const FASTNODED_PATH = /fastnoded/
const ELECTRON = /\/electron\//

const memInstallPath = () => {
    // first, try to launch via $HOME/.local/share/fastnode/fastnoded, as this is a wrapper which handles restarts
    let homePath = path.join(os.homedir(), ".local", "share", "fastnode", "fastnoded")
    if (fs.existsSync(homePath)) {
        return homePath
    }

    // then, try to launch via /opt/fastnode/fastnoded
    let globalPath = "/opt/fastnode/fastnoded"
    if (fs.existsSync(globalPath)) {
        return globalPath
    }

    // return the path to fastnoded based on __dirname, e.g like $prefix/linux-unpacked/resources/app.asar/build/utils/fastnoded
    let dir = __dirname;
    while (dir.length > 0) {
        if (path.basename(dir) === "linux-unpacked") {
            return path.join(path.dirname(dir), "fastnoded")
        }

        dir = path.dirname(dir)
    }
    return "fastnoded"
}

const installPath = memInstallPath()

module.exports = {
    //resolves with a pid
    isFastnodeRunning() {
        return spawnPromise('/bin/ps', ['-axo', 'pid,command'], {encoding: 'utf8',}, 'ps_error')
            .then(stdout => {
                const procs = stdout.split('\n')
                const fastnodeprocs = procs.filter(s => FASTNODED_PATH.test(s) && !ELECTRON.test(s)) //filter out dev electron builds
                if (fastnodeprocs.length > 0) {
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
        let env = Object.create(process.env)
        env.FASTNODE_SKIP_ONBOARDING = '1'
        let fastnoded = spawn(installPath, ['--sidebar-restart'], {stdio: "ignore", env: env, detached: true})
        fastnoded.unref();
        return Promise.resolve()
    },

    stopFastnode() {
        return this.isFastnodeRunning()
            .then((res) => {
                if (res.running) {
                    //even if multiple Fastnode processes, for now, just assume 0-index
                    const pid = res.processes[0].trim().split(" ")[0]
                    return spawnPromise('kill', [pid], 'kill_error')
                } else {
                    return Promise.resolve()
                }
            })
    },
}