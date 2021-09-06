/* eslint-disable react/prop-types */
import React from "react";
import "../styles/globals.css";
import { Provider } from "react-redux";
import store from "../app/store";
import { Layout } from "@/components/Layout";
import { NoSsr } from "@material-ui/core";

function MyApp({ Component, pageProps }) {
  return (
    <NoSsr>
      <Provider store={store}>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </Provider>
    </NoSsr>
  );
}

export default MyApp;
