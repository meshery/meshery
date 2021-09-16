/* eslint-disable react/prop-types */
import React from "react";
import { Provider } from "react-redux";
import store from "../app/store";
import { Layout } from "@/components/Layout/Layout";
import { NoSsr } from "@material-ui/core";
import { StylesProvider } from "@material-ui/styles";

function MyApp({ Component, pageProps }) {
  return (
    <NoSsr>
      <StylesProvider injectFirst>
        <Provider store={store}>
          <Layout>
            <Component {...pageProps} />
          </Layout>
        </Provider>
      </StylesProvider>
    </NoSsr>
  );
}

export default MyApp;
