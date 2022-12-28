/* eslint-disable no-unused-vars */
import dataFetch from "../../../lib/data-fetch";
import { ctxUrl } from "../../../utils/multi-ctx";

export const verifyGrafanaConnection = (grafanaUrl) => {
  return new Promise((res, rej) => {
    if (grafanaUrl){
      pingGrafana(
        result => res("Grafana connected !"),
        error => rej("Grafana not connected ! " + error)
      )
      return
    }
    rej("Grafana not connected! "+ "Url not found")
  })
}


export const pingGrafanaWithNotification = (updateProgress, action, enqueueSnackbar) => {

  updateProgress({ showProgress : true })

  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      enqueueSnackbar("Grafana pinged!", { variant : "success",
        autoHideDuration : 2000,
        action });
    }
  }

  const errorCb = (error) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar("Grafana could not be pinged! : "+error, { variant : "error",
      autoHideDuration : 2000,
      action });
  }

  pingGrafana(successCb, errorCb)
}

export const pingGrafana = (successCb, errorCb) =>
  dataFetch(
    "/api/telemetry/metrics/grafana/ping",
    { credentials : "include" },
    successCb,
    errorCb
  );

export const verifyPrometheusConnection =  (prometheusUrl) => {
  console.log(prometheusUrl)
  return new Promise((res, rej) => {
    if (prometheusUrl !== ""){
      pingPrometheus(
        result => res("Prometheus connected !"),
        error => rej("Prometheus not connected ! " + error)
      )
      return
    } else
      rej("Prometheus not connected! "+ "Url not found")
  })
}


export const pingPrometheusWithNotification = (updateProgress, action, enqueueSnackbar) => {

  updateProgress({ showProgress : true })

  const successCb = (result) => {
    updateProgress({ showProgress : false });
    if (typeof result !== "undefined") {
      enqueueSnackbar("Prometheus pinged!", { variant : "success",
        autoHideDuration : 2000,
        action });
    }
  }

  const errorCb = (error) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar("Prometheus could not be pinged! : "+error, { variant : "error",
      autoHideDuration : 2000,
      action });
  }

  pingPrometheus(successCb, errorCb)
}

export const pingPrometheus = (successCb, errorCb) =>
  dataFetch(
    "/api/telemetry/metrics/ping",
    { credentials : "include" },
    successCb,
    errorCb
  );


export const fetchPromGrafanaScanData = (ctx) => {
  return new Promise((res, rej) => {
    dataFetch(
      ctxUrl('/api/system/meshsync/grafana', ctx),
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        let metricsUrls = { grafana : [], prometheus : [] }
        if (!result) res(metricsUrls);
        console.log()

        if (Array.isArray(result.prometheus)) {
          const urls = extractURLFromScanData(result.prometheus);
          metricsUrls.prometheus = urls
        }

        if (Array.isArray(result.grafana)) {
          const urls = extractURLFromScanData(result.grafana);
          metricsUrls.grafana = urls
        }
        res(metricsUrls)
      },
      (err) => rej("Unable to fetch grafana and prometheus scan data:"+err)
    )
  }
  )
}

/**
   * extractURLFromScanData scans the ingress urls from the
   * mesh scan data and returns an array of the response
   * @param {object[]} scannedData
   * @returns {string[]}
   */
const extractURLFromScanData = (scannedData) => {
  const result = [];
  scannedData.forEach(data => {
    // Add loadbalancer based url
    if (Array.isArray(data.status?.loadBalancer?.ingress)) {
      data.status.loadBalancer.ingress.forEach(lbdata => {
        let protocol = "http";

        // Iterate over ports exposed by the service
        if (Array.isArray(data.spec.ports)) {
          data.spec.ports.forEach(({ port }) => {
            if (port === 443) protocol = "https";

            // From kubernetes v1.19 docs
            // Hostname is set for load-balancer ingress points that are DNS based (typically AWS load-balancers)
            // IP is set for load-balancer ingress points that are IP based (typically GCE or OpenStack load-balancers)
            let address = lbdata.ip || lbdata.hostname;
            if (address) result.push(`${protocol}://${address}:${port}`);
          })
        }
      })
    }

    // Add clusterip based url
    // As per kubernetes v1.19 api, "None", "" as well as a valid ip is a valid clusterIP
    // Looking for valid ipv4 address
    if (data.spec.clusterIP?.match(/^((25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])(\.(?!$)|$)){4}$/g)?.[0]) {
      let protocol = "http";
      if (Array.isArray(data.spec.ports)) {
        data.spec.ports.forEach(({ port }) => {
          if (port === 443) protocol = "https";
          result.push(`${protocol}://${data.spec.clusterIP}:${port}`);
        })
      }
    }
  })

  return result
}



export const handleGrafanaConfigure = (grafanaURL, grafanaAPIKey, updateProgress, enqueueSnackbar, action, updateGrafanaConfig) => {

  if (
    grafanaURL === "" ||
      !(grafanaURL.toLowerCase().startsWith("http://") || grafanaURL.toLowerCase().startsWith("https://"))
  ) {
    return;
  }
  const data = { grafanaURL,
    grafanaAPIKey, };
  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join("&");

  updateProgress({ showProgress : true });
  dataFetch(
    "/api/telemetry/metrics/grafana/config",
    {
      method : "POST",
      credentials : "include",
      headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", },
      body : params,
    },
    (result) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Grafana was configured!", { variant : "success",
          autoHideDuration : 2000,
          action });
        updateGrafanaConfig({ grafana : { grafanaURL,
          grafanaAPIKey, }, });
      }
    },
    (err) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Grafana was not configured! :" + err, { variant : "error",
          autoHideDuration : 2000,
          action });
      }
    }

  )
}

export const handlePrometheusConfigure = (prometheusURL, updateProgress, enqueueSnackbar, action, updatePrometheusConfig) => {

  if (
    prometheusURL === "" ||
      !(prometheusURL.toLowerCase().startsWith("http://") || prometheusURL.toLowerCase().startsWith("https://"))
  ) {
    return;
  }
  const data = { prometheusURL, };
  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join("&");

  updateProgress({ showProgress : true });
  dataFetch(
    "/api/telemetry/metrics/config",
    {
      method : "POST",
      credentials : "include",
      headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", },
      body : params,
    },
    (result) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Prometheus was configured!", { variant : "success",
          autoHideDuration : 2000,
          action });
        updatePrometheusConfig({ prometheus : { prometheusURL,
          selectedPrometheusBoardsConfigs : [], }, });
      }
    },
    (err) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Prometheus was not configured! :" + err, { variant : "error",
          autoHideDuration : 2000,
          action });
      }
    }

  )
}

export const deleteMetricsComponentConfig = (componentName) => (successCb, errorCb) => dataFetch(
    `/api/telemetry/metrics${componentName === 'Grafana' ? "/grafana"
      : ""}/config`,
    {
      method : "DELETE",
      credentials : "include", },
    successCb,
    errorCb
)




