import dataFetch from '../../../lib/data-fetch'
import { EVENT_TYPES } from '../../../lib/event-types';
import { commitMutation } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";
import AdapterStatusMutation from "../../graphql/mutations/__generated__/AdapterStatusMutation.graphql.js";
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


export const pingAdapterWithNotification = (notify, updateProgress, adapterLoc) => {
  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      notify({ message : "Adapter pinged!", event_type : EVENT_TYPES.SUCCESS })
    }
  }

  const errorCb = (err) => {
    updateProgress({ showProgress : false });
    notify({ message : "Adapter ping failed!", event_type : EVENT_TYPES.ERROR, details : err.toString() })
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

export const configureAdapterWithNotification = (notify, updateProgress, adapterLocation, updateAdaptersInfo) => {
  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      notify({ message : "Adapter was configured!", event_type : EVENT_TYPES.SUCCESS })
      updateAdaptersInfo({ meshAdapters : result });
    }
  }

  const errorCb = (err) => {
    updateProgress({ showProgress : false });
    notify({ message : "Adapter configuration failed!", event_type : EVENT_TYPES.ERROR, details : err.toString() })
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


export const handleDeleteAdapter =  (successCb, errorCb) => {
  const environment = createRelayEnvironment({});
  const mutation = AdapterStatusMutation;
  const vars = { input : { targetStatus : "DISABLED", targetPort : "10000", adapter : "meshery-istio" } };

  commitMutation(environment, {
    mutation,
    vars,
    onCompleted : successCb,
    onError : errorCb,
  });
};
