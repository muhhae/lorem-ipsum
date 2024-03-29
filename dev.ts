import { $ } from 'bun'
import os from 'os'

let current_os = os.platform()

switch (current_os) {
    case 'win32':
        await $`air -build.bin '.\\app.exe' -build.cmd 'bun run build:windows'`
        break
    case 'linux':
        await $`air -build.bin './app' -build.cmd 'bun run build:linux'`
        break
}
