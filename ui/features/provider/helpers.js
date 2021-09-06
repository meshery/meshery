import dataFetch from "@/lib/data-fetch";

export const fetchProviderDetails = () => new Promise((res, rej) => {
    dataFetch(
      "/api/provider/capabilities",
      {
        credentials: "same-origin",
        method: "GET",
      },
      (result) => {
        res(result)
      },
      (err) => rej(err)
    );
}) 