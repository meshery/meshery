import dataFetch from '../../lib/data-fetch';
import { EVENT_TYPES } from '../../lib/event-types';

/**
 * fetch the adapters that are available
 *
 */
export const fetchAvailableAdapters = () => {
  return new Promise((res, rej) =>
    dataFetch(
      '/api/system/adapters',

      {
        method: 'GET',
        credentials: 'include',
      },

      (result) => {
        if (typeof result !== 'undefined') {
          const options = result.map((res) => ({
            value: res.adapter_location,
            label: res.adapter_location,
            name: res.name,
            version: res.version,
          }));
          res(options);
        }
      },

      (err) => rej(err),
    ),
  );
};

export const pingAdapterWithNotification = (notify, updateProgress, adapterLoc) => {
  const successCb = (result) => {
    updateProgress({ showProgress: false });
    if (typeof result !== 'undefined') {
      notify({ message: 'Adapter pinged!', event_type: EVENT_TYPES.SUCCESS });
    }
  };

  const errorCb = (err) => {
    updateProgress({ showProgress: false });
    notify({
      message: 'Adapter ping failed!',
      event_type: EVENT_TYPES.ERROR,
      details: err.toString(),
    });
  };

  pingAdapter(adapterLoc, successCb, errorCb);
};

export const pingAdapter = (adapterLoc, successCb, errorCb) => {
  dataFetch(
    `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
    {
      credentials: 'include',
    },
    successCb,
    errorCb,
  );
};

export const configureAdapterWithNotification = (
  notify,
  updateProgress,
  adapterLocation,
  updateAdaptersInfo,
) => {
  const successCb = (result) => {
    updateProgress({ showProgress: false });
    if (typeof result !== 'undefined') {
      // Server-side event handles notification via Notification Center
      updateAdaptersInfo({ meshAdapters: result });
    }
  };

  const errorCb = () => {
    updateProgress({ showProgress: false });
    // Server-side event handles error notification via Notification Center
  };

  configureAdapter(successCb, errorCb, adapterLocation);
};

export const configureAdapter = (successCb, errorCb, adapterLocation) => {
  const data = { meshLocationURL: adapterLocation };

  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join('&');

  return dataFetch(
    '/api/system/adapter/manage',
    {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      body: params,
    },
    successCb,
    errorCb,
  );
};

export const handleDeleteAdapter = (successCb, errorCb) => (adapterLoc) => {
  return dataFetch(
    `/api/system/adapter/manage?adapter=${encodeURIComponent(adapterLoc)}`,
    {
      method: 'DELETE',
      credentials: 'include',
    },
    successCb,
    errorCb,
  );
};
