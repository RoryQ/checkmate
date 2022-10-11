const childProcess = require('child_process')
const os = require('os')
const process = require('process')


function chooseBinary(versionTag) {
    const platform = os.platform()
    const arch = os.arch()

    const binaryVersion = versionTag.substring(1);

    if (platform === 'linux' && arch === 'x64') {
        return `checkmate_${binaryVersion}_Linux_x86_64`
    }
    if (platform === 'linux' && arch === 'arm64') {
        return `checkmate_${binaryVersion}_Linux_arm64`
    }
    if (platform === 'windows' && arch === 'x64') {
        return `checkmate_${binaryVersion}_Windows_x86_64`
    }
    if (platform === 'windows' && arch === 'arm64') {
        return `checkmate_${binaryVersion}_Windows_arm64`
    }
    if (platform === 'darwin' && arch === 'x64') {
        return `checkmate_${binaryVersion}_Darwin_x86_64`
    }
    if (platform === 'darwin' && arch === 'arm64') {
        return `checkmate_${binaryVersion}_Darwin_arm64`
    }

    console.error(`Unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function downloadBinary(versionTag, binary) {
    const url = `https://github.com/RoryQ/checkmate/releases/download/${versionTag}/${binary}.tar.gz`
    let command = `curl -sL '${url}' | tar xvz`;
    let result = childProcess.spawnSync(command, {shell: true});
    return result.status
}

function main() {
    const versionTag = "v0.0.1"
    let status = downloadBinary(versionTag, chooseBinary("v0.0.1"))
    if (typeof status === 'number' && status > 0) {
        process.exit(status)
    }

    const spawnSyncReturns = childProcess.spawnSync(`${__dirname}/checkmate`, { stdio: 'inherit' })
    status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
