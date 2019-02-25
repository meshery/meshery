import fetch from 'isomorphic-unfetch'

const dataFetch = (url, options = {}, successFn, errorFn) => {
  fetch(url, options)
    .then(res => {
      if (res.status === 401 || res.redirected){
        window.location = "/login"; // for local dev thru node server
        // window.location.reload(); // for use with Go server
      }
      if (res.ok) {
        return res.json();
      } else {
        res.text().then(errorFn);
      }
    }).then(successFn)
    .catch(errorFn);
}

export default dataFetch;