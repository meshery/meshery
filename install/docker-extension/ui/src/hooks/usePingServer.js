import { useEffect } from "react"
import { useState } from "react"
import axios from "axios"
import ping from "web-pingjs"

export const usePingServer = (path, { host, port }) => {
  const [isServerAvailable, setIsServerAvailable] = useState(false)
  useEffect(() => {
    setInterval(() => {
      ping("http://127.0.0.1:9081/api/system/version").then((res) => {
        setIsServerAvailable(true)
      }).catch(() => {
        setIsServerAvailable(false)
      })
      // axios.get("http://localhost:9081/api/system/versionsa"
      // )
      //   .then((obj) => {
      //     console.log(obj)
      //     if (obj.status >= 200 && obj.status < 300) setIsServerAvailable(true)
      //     else setIsServerAvailable(false)
      //   })
      //   .catch((obj) => {
      //     setIsServerAvailable(false)
      //     console.log(obj)
      //   })
    }, 2000)
  }, [])
  return isServerAvailable
}
