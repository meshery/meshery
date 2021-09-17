/* eslint-disable react/prop-types */
import React from "react";
import { Provider } from "react-redux";
import store from "../app/store";
import { Layout } from "@/components/Layout/Layout";
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
