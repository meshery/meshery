import React from 'react';
import PropTypes from 'prop-types';
import Document, { Head, Main, NextScript, Html } from 'next/document';
import { createStyleRegistry } from 'styled-jsx';

const registry = createStyleRegistry();
const flush = registry.flush();
class MesheryDocument extends Document {
  render() {
    return (
      <Html lang="en" dir="ltr">
        <Head>
          <meta charSet="utf-8" />
          {/**
           * content="no-referrer" included to avoid 403 errors on Google avatars
           */}
          <meta name="referrer" content="no-referrer" />
          <link rel="icon" href="/static/favicon.png" />

          {/* Google Tag Manager */}
          {/* eslint-disable-next-line @next/next/next-script-for-ga */}
          <script
            dangerouslySetInnerHTML={{
              __html: `(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start': new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0], j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src= 'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
            // Fetch user preferences
            fetch("/api/user/prefs", { credentials: 'include' })
              .then((res) => res.json())
              .then((res) => {
                if (res && res?.anonymousUsageStats === false) {
                  // User opted out of tracking, disable GTM
                  w[l] = []; // Clear the dataLayer array
                  w['ga-disable-'+i] = true; // Disable Google Analytics tracking
                }
              }).catch((err) => {
                console.error("error while fetching user prefs for googletagmanag",err);
              });
          })(window,document,'script','dataLayer','GTM-TFLZDSQ');`,
            }}
          />
          {/* End Google Tag Manager */}

          {/**
           * For hiding the scrollbar without losing the scroll functionality
           * add the class "hide-scrollbar" to hide scrollbar for that element
           * Only applicable for Chrome, safari and newest version of Opera
           */}
          <style type="text/css">
            {
              '\
            .hide-scrollbar::-webkit-scrollbar {\
              width: 0 !important;\
            }\
          .reduce-scrollbar-width::-webkit-scrollbar {\
              width: 0.3em !important;\
            }\
          '
            }
          </style>
        </Head>
        <body>
          {/* Google Tag Manager (noscript) */}
          <noscript
            dangerouslySetInnerHTML={{
              __html: `<iframe src="https://www.googletagmanager.com/ns.html?id=GTM-TFLZDSQ" height="0" width="0" style="display:none;visibility:hidden"></iframe>`,
            }}
          />
          {/* End Google Tag Manager (noscript) */}

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
        {flush || null}
      </React.Fragment>
    ),
  };
};

export default MesheryDocument;
