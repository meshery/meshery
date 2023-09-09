import axios from 'axios';

const dataFetch = async (url, options = {}, successFn, errorFn) => {
  // Set default error function if not provided
  if (!errorFn) {
    errorFn = (err) => {
      console.error(`Error fetching ${url} --DataFetch`, err);
    };
  }

  try {
    const response = await axios(url, options);

    // Handle redirects or unauthorized status
    if (response.status === 401 || (response.request && response.request.responseURL !== url)) {
      if (window.location.host.endsWith("3000")) {
        window.location = "/user/login"; // for local dev thru node server
      } else {
        window.location.reload(); // for use with Go server
      }
      return;
    }

    // Call the success function with the response data
    successFn(response.data);
  } catch (error) {
    // Handle the error
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      errorFn(error.response.data);
    } else if (error.request) {
      // The request was made but no response was received
      errorFn(error.request);
    } else {
      // Something happened in setting up the request that triggered an Error
      errorFn(error.message);
    }
  }
};

export default dataFetch;
