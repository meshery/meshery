import fetch from 'isomorphic-unfetch'

const dataFetch = (url, options = {}, successFn, errorFn) => {
  fetch(url, options)
    .then(res => {
      if (res.ok) {
        return res.json();
      } else {
        res.text().then(errorFn);
      }
    }).then(successFn)
    .catch(errorFn);
}

export default dataFetch;