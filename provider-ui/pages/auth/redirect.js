import { useEffect } from "react";
import { PROVIDER_URL } from "../../lib/data-fetch";

const Redirect = () => {
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const returnTo = params.get("return_to");
    let redirectURL;
    if (returnTo) {
      redirectURL = new URL(`${returnTo}?${params.toString()}`);
    } else {
      redirectURL = new URL(`${PROVIDER_URL}/?${params.toString()}`);
    }

    window.location.href = redirectURL.toString();
  }, []);
  return <></>;
};

export default Redirect;
