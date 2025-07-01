import React from "react";
import ReactDOM from "react-dom";
import ExtensionComponent from "./components/ExtensionComponent/ExtensionComponent";
import "./styles.css";

if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker
      .register("./pub_static/service-worker.js")
      .then((registration) => {
        console.log("ServiceWorker registered: ", registration);
      })
      .catch((registrationError) => {
        console.log("ServiceWorker registration failed: ", registrationError);
      });
  });
}

ReactDOM.render(
  <React.StrictMode>
    <ExtensionComponent />
  </React.StrictMode>,
  document.getElementById("root"),
);
