const cp = require('child_process')
const os = require('os')
const process = require('process')

function logDebug(...data) {
    if (!!process.env.CHECKMATE_DEBUG) {
        console.debug(...data)
    }
}

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
    logDebug(url)
    const result = cp.execSync(`curl --silent --location --remote-header-name  ${url} | tar xvz`)
    logDebug(result.toString())
    return result.status
}

function determineVersion() {
    if (!!process.env.INPUT_VERSION) {
        return process.env.INPUT_VERSION
    }
    const result = cp.execSync(`curl --silent --location "https://api.github.com/repos/RoryQ/checkmate/releases/latest" | jq  -r ".. .tag_name? // empty"`)
    return result.toString().trim();
}

function main() {
    logDebug(process.env.INPUT_PATHS)
    logDebug("started")
    const versionTag = determineVersion()
    let status = downloadBinary(versionTag, chooseBinary(versionTag))
    if (typeof status === 'number' && status > 0) {
        process.exit(status)
    }

    console.log(`     _____________/\\/\\________________________________/\\/\\______________________________________/\\/\\_________________
    ___/\\/\\/\\/\\__/\\/\\__________/\\/\\/\\______/\\/\\/\\/\\__/\\/\\__/\\/\\__/\\/\\/\\__/\\/\\____/\\/\\/\\______/\\/\\/\\/\\/\\____/\\/\\/\\___ 
   _/\\/\\________/\\/\\/\\/\\____/\\/\\/\\/\\/\\__/\\/\\________/\\/\\/\\/\\____/\\/\\/\\/\\/\\/\\/\\______/\\/\\______/\\/\\______/\\/\\/\\/\\/\\_  
  _/\\/\\________/\\/\\__/\\/\\__/\\/\\________/\\/\\________/\\/\\/\\/\\____/\\/\\__/\\__/\\/\\__/\\/\\/\\/\\______/\\/\\______/\\/\\_______   
 ___/\\/\\/\\/\\__/\\/\\__/\\/\\____/\\/\\/\\/\\____/\\/\\/\\/\\__/\\/\\__/\\/\\__/\\/\\______/\\/\\__/\\/\\/\\/\\/\\____/\\/\\/\\______/\\/\\/\\/\\_    
________________________________________________________________________________________________________________     `)

    const spawnSyncReturns = cp.spawnSync(`./checkmate`, { stdio: 'inherit' })
    logDebug(spawnSyncReturns)
    status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
