import { wrapper } from "@/core/store/store";
import { Provider } from "react-redux";
import { AppThemeProvider } from "@/core/providers/AppThemeProvider";

export default function App({ Component, pageProps, ...rest }) {
  console.log('rest: ', rest);
  const { store, props } = wrapper.useWrappedStore(rest);

  return (
    <Provider store={store}>
      <AppThemeProvider>
        <Component {...props.pageProps} />
      </AppThemeProvider>
    </Provider>
  )
};