import dynamic from 'next/dynamic'

const WasmComponent = dynamic({
  loader : async () => {
    const go = new globalThis.Go();
    WebAssembly.instantiateStreaming(fetch("static/main.wasm"), go.importObject).then((result) => {
      go.run(result.instance)
    });
  },
})

export default function Index() {
  return (
    <div>
      <WasmComponent/>
    </div>
  )
}