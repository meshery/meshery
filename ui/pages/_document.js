import React from "react";
import PropTypes from "prop-types";
import Document, { Head, Main, NextScript, Html } from "next/document";
import flush from "styled-jsx/server";

/**
 * setupGA setups up the google analtyics in meshery UI
 *
 * This function may not be invoked here (server side) as "window" object does
 * not exist on the server side. Hence this function is stringified and
 * forced to execute on the client side only.
 *
 */
function setupGA() {
  window.dataLayer = window.dataLayer || [];

  function gtag() {
    dataLayer.push(arguments);
  }

  fetch("/api/user/prefs", { credentials: "include" })
    .then((res) => res.json())
    .then((res) => {
      if (res?.anonymousUsageStats) {
        gtag("js", new Date());
        gtag("config", "G-8Q51RLT8TZ", {
          page_path: window.location.pathname,
        });
      }
    })
    .catch((err) => console.error(err));
}

class MesheryDocument extends Document {
  render() {
    return (
      <Html lang="en" dir="ltr">
        <Head>
          <meta charSet="utf-8" />
          <link rel="icon" href="/static/favicon.png" />
          <link rel="shortcut icon" href="/ui/public/static/img/meshery-logo/meshery-logo.svg" />
          <title>Meshery</title>
          <script async src={`https://www.googletagmanager.com/gtag/js?id=G-8Q51RLT8TZ`} />

          <script
            dangerouslySetInnerHTML={{
              __html: `${"" + setupGA}; setupGA();`,
            }}
          />

          {/**
           * For hiding the scrollbar without losing the scroll functionality
           * add the class "hide-scrollbar" to hide scrollbar for that element
           * Only applicable for Chrome, safari and newest version of Opera
           */}
          <style type="text/css">
            {
              "\
            .hide-scrollbar::-webkit-scrollbar {\
              width: 0 !important;\
            }\
          "
            }
          </style>
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

MesheryDocument.getInitialProps = (ctx) => {
  // resolution order
  //
  // on the server:
  // 1. app.getInitialProps
  // 2. page.getInitialProps
  // 3. document.getInitialProps
  // 4. app.render
  // 5. page.render
  // 6. document.render
  //
  // on the server with error:
  // 1. document.getInitialProps
  // 2. app.render
  // 3. page.render
  // 4. document.render
  //
  // on the client
  // 1. app.getInitialProps
  // 2. page.getInitialProps
  // 3. app.render
  // 4. page.render

  // render app and page and get the context of the page with collected side effects.
  let pageContext;
  const page = ctx.renderPage((Component) => {
    const WrappedComponent = (props) => {
      pageContext = props.pageContext;
      return <Component {...props} />;
    };

    WrappedComponent.propTypes = {
      pageContext: PropTypes.object.isRequired,
    };

    return WrappedComponent;
  });

  let css;
  // it might be undefined, e.g. after an error.
  if (pageContext) {
    css = pageContext.sheetsRegistry.toString();
  }

  return {
    ...page,
    pageContext,
    // styles fragment is rendered after the app and page rendering finish.
    styles: (
      <React.Fragment>
        <style
          id="jss-server-side"
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{ __html: css }}
        />
        {flush() || null}
      </React.Fragment>
    ),
  };
};

export default MesheryDocument;
