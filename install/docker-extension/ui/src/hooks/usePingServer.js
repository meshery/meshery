import { useEffect, useState } from "react"
import ping from "web-pingjs"

export const usePingServer = () => {
  const [isServerAvailable, setIsServerAvailable] = useState(false)

  useEffect(() => {
    let isMounted = true

    const checkServerAvailability = () => {
      ping("http://127.0.0.1:9081/api/system/version").then((res) => {
        if (isMounted) {
          setIsServerAvailable(true)
        }
      }).catch(() => {
        if (isMounted) {
          setIsServerAvailable(false)
        }
      })
    }

    checkServerAvailability()
    const intervalId = setInterval(checkServerAvailability, 2000)

    return () => {
      isMounted = false
      clearInterval(intervalId)
    }
  }, [])

  return isServerAvailable
}
