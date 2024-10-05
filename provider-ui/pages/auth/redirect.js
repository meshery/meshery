
import { useEffect } from "react";


const Redirect = () => {
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const returnTo = params.get("return_to");
    let redirectURL;
    if (returnTo) {
      redirectURL = new URL(`${returnTo}?${params.toString()}`)
    } else {
      redirectURL = new URL(`https://meshery.layer5.io/?${params.toString()}`)
    }

    window.location.href = redirectURL.toString();

  }, [])
  return (
    <></>
  )
}

export default Redirect;