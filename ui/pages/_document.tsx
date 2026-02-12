import React from 'react';
import Document, {
  Html,
  Head,
  Main,
  NextScript,
  DocumentContext,
  DocumentInitialProps,
} from 'next/document';
import { PureHtmlLoadingScreen } from '@/components/LoadingComponents/LoadingComponentServer';

class MesheryDocument extends Document {
  static async getInitialProps(ctx: DocumentContext): Promise<DocumentInitialProps> {
    const initialProps = await Document.getInitialProps(ctx);
    return { ...initialProps };
  }

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

          {/* Preload Qanelas Soft font for loading screen */}
          <link
            rel="preload"
            href="/static/fonts/qanelas-soft/QanelasSoftRegular.otf"
            as="font"
            type="font/otf"
            crossOrigin="anonymous"
          />
          <link
            href="/static/fonts/qanelas-soft/QanelasSoftRegular.otf"
            as="font"
            type="font/otf"
            crossOrigin="anonymous"
          />

          {/* Google Tag Manager */}
          {/* eslint-disable-next-line @next/next/next-script-for-ga */}
          <script
            dangerouslySetInnerHTML={{
              __html: `(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start': new Date().getTime(),event:'gtm'});var f=d.getElementsByTagName(s)[0], j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src= 'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
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
            {`
              .hide-scrollbar::-webkit-scrollbar {
                width: 0 !important;
              }
              .reduce-scrollbar-width::-webkit-scrollbar {
                width: 0.3em !important;
              }
            `}
          </style>
        </Head>
        <body style={{ margin: 0 }}>
          {/* Google Tag Manager (noscript) */}
          <noscript
            dangerouslySetInnerHTML={{
              __html: `<iframe src="https://www.googletagmanager.com/ns.html?id=GTM-TFLZDSQ" height="0" width="0" style="display:none;visibility:hidden"></iframe>`,
            }}
          />
          {/* End Google Tag Manager (noscript) */}
          {/* Pre-React script */}
          {/* eslint-disable-next-line @next/next/no-sync-scripts */}
          <script src="/loadingMessages.js"></script>

          <PureHtmlLoadingScreen id={'PRE_REACT_LOADER'} message="" />
          <Main />
          <NextScript />

          <script
            dangerouslySetInnerHTML={{
              __html: `
                (function () {
                  const loaderId = "PRE_REACT_LOADER-text-message"

                  try {
                    var el = document.getElementById(loaderId)
                    if (!el) return;

                    el.textContent = window.Loader.PersistedRandomLoadingMessage()
                  } catch (e) {
                    console.log("Failed to set loading message",e)
                  }
                })();
              `,
            }}
          />
        </body>
      </Html>
    );
  }
}

export default MesheryDocument;
