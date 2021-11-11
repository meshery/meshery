/* eslint-disable react/prop-types */
import React from "react";
import { Provider } from "react-redux";
import store from "../app/store";
import { Layout } from "@/components/Layout/Layout";
import { NoSsr } from "@mui/material";
import { SnackbarProvider } from 'notistack';
import Slide from '@mui/material/Slide';

function MyApp({ Component, pageProps }) {
  return (
    <NoSsr>
    <SnackbarProvider
      hideIconVariant={false}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
    }}
    TransitionComponent={Slide}
    >
      <Provider store={store}>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </Provider>
      </SnackbarProvider>
    </NoSsr>
  );
}

export default MyApp;
