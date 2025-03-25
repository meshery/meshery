import { useEffect } from "react"
import { useState } from "react"
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
      
    }, 2000)
  }, [])
  return isServerAvailable
}
