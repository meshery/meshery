import dataFetch from "@/lib/data-fetch";

export const fetchMesheryServerDetails = () =>
  new Promise((res, rej) =>
    dataFetch(
      "/api/system/version",
      {
        credentials: "same-origin",
        method: "GET",
      },
      (result) => res(result),
      (err) => rej(err)
    )
  );
