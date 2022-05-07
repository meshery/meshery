import { useEffect, useState } from "react"
import axios from "axios"

export const useLogin = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    axios.get("http://127.0.0.1:7877/token/store")
      .then((obj) => {
        if (obj.status >= 200 && obj.status < 300) setIsLoggedIn(true)
        else setIsLoggedIn(false)
      })
      .catch((obj) => {
        setIsLoggedIn(false)
      })
  }, [])

  return isLoggedIn
}
