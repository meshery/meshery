import dataFetch from "../../../lib/data-fetch"

/**
  * fetch the adapters that are available
  *
*/
export const fetchAvailableAdapters = () => {
  return new Promise((res, rej) =>

    dataFetch(
      "/api/system/adapters",

      {
        method : "GET",
        credentials : "include", },

      (result) => {
        if (typeof result !== "undefined") {
          const options = result.map((res) => ({
            value : res.adapter_location,
            label : res.adapter_location,
            name : res.name,
            version : res.version
          }));
          res(options)
        }
      },

      (err) => rej(err)

    )
  )
};


export const pingAdapterWithNotification = (updateProgress, enqueueSnackbar, action, adapterLoc) => {

  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      enqueueSnackbar("Adapter pinged!", { variant : "success",
        autoHideDuration : 2000,
        action });
    }
  }

  const errorCb = (err) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar("Adapter ping failed! : "+err, { variant : "error",
      autoHideDuration : 2000,
      action });
  }

  pingAdapter(adapterLoc, successCb, errorCb)
}


export const pingAdapter = (adapterLoc, successCb, errorCb) => {

  dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials : "include", },
      successCb,
      errorCb
  );
}

export const configureAdapterWithNotification = (enqueueSnackbar, updateProgress, action, adapterLocation, updateAdaptersInfo) => {

  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      enqueueSnackbar("Adapter was configured!", {
        variant : "success",
        "data-cy" : "adapterSuccessSnackbar",
        autoHideDuration : 2000,
        action
      });
      updateAdaptersInfo({ meshAdapters : result });
    }
  }

  const errorCb = (err) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar("Adapter configuration failed! : "+err, { variant : "error",
      autoHideDuration : 2000,
      action });
  }

  configureAdapter(successCb, errorCb, adapterLocation)
}

export const configureAdapter = (successCb, errorCb, adapterLocation) => {

  const data = { meshLocationURL : adapterLocation };

  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join("&");

  return dataFetch(
    "/api/system/adapter/manage",
    {

      method : "POST",
      credentials : "include",
      headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", },
      body : params,
    },
    successCb,
    errorCb
  );
}


export const handleDeleteAdapter =  (successCb, errorCb) => (adapterLoc) => {

  return  dataFetch(
      `/api/system/adapter/manage?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        method : "DELETE",
        credentials : "include",
      },
      successCb,
      errorCb
  );
};
