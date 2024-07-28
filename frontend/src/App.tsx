import { useEffect, useState } from "react"
import reactLogo from "./assets/react.svg"
import viteLogo from "/vite.svg"
import "./App.css"

import { createPromiseClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import { FileService } from "./gen/file/v1/file_connect"

function App() {
  const [count, setCount] = useState(0)
  const [filenames, setFilenames] = useState<string[]>([])

  const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
  })

  // Create a client to connect to the server
  const client = createPromiseClient(FileService, transport)

  useEffect(() => {
    client.listFiles({}).then((response) => {
      const filenames = response.filenames.map((filename) => filename)
      setFilenames(filenames)
    })
  }, [])

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
      <ul>
        {filenames?.map((filename) => <li key={filename}>{filename}</li>)}
      </ul>
    </>
  )
}

export default App
