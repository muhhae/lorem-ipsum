import { $ } from 'bun'

async function bun_build() {
    const start = performance.now()
    await $`bun build:linux`.quiet()
    return performance.now() - start
}
async function node_build() {
    const start = performance.now()
    await $`npm run build:linux`.quiet()
    return performance.now() - start
}

let node_total_time = 0
let bun_total_time = 0

async function main() {
    const bun_time = await bun_build()
    console.log(`bun: ${bun_time} ms`)

    const node_time = await node_build()
    console.log(`node: ${node_time} ms\n`)

    bun_total_time += bun_time
    node_total_time += node_time
}

for (let i = 0; i < 10; i++)
    await main()

console.log('------------------')
console.log(`bun average: ${bun_total_time / 10} ms`)
console.log(`node average: ${node_total_time / 10} ms`)
console.log('------------------')

