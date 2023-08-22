import axios from "axios";

const instance = axios.create({
  withCredentials : true // for pushing client-cookies in all requests to server
});

instance.interceptors.response.use(response => {
  if (response.request.responseURL.includes("/auth/login")) {
    window.location = "/auth/login";
    window.onbeforeunload = null;
  }
  return response;
});


export default instance;