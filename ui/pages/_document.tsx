import React from 'react';
import Document, {
  Html,
  Head,
  Main,
  NextScript,
  DocumentContext,
  DocumentInitialProps,
} from 'next/document';
import createEmotionServer from '@emotion/server/create-instance';
import createCache from '@emotion/cache';
import { PureHtmlLoadingScreen } from '@/components/shared/LoadingState/LoadingComponentServer';

// Create emotion cache for SSR
const createEmotionCache = () => {
  return createCache({ key: 'css', prepend: true });
};

interface MyDocumentProps extends DocumentInitialProps {
  emotionStyleTags: React.ReactElement[];
}

class MesheryDocument extends Document<MyDocumentProps> {
  static async getInitialProps(ctx: DocumentContext): Promise<MyDocumentProps> {
    const originalRenderPage = ctx.renderPage;

    const cache = createEmotionCache();
    const { extractCriticalToChunks } = createEmotionServer(cache);

    ctx.renderPage = () =>
      originalRenderPage({
        enhanceApp: (App: React.ComponentType<any>) =>
          function EnhanceApp(props) {
            return <App emotionCache={cache} {...props} />;
          },
      });

    const initialProps = await Document.getInitialProps(ctx);

    // Extract critical CSS from the rendered page
    const emotionStyles = extractCriticalToChunks(initialProps.html);
    const emotionStyleTags = emotionStyles.styles.map((style) => (
      <style
        data-emotion={`${style.key} ${style.ids.join(' ')}`}
        key={style.key}
        dangerouslySetInnerHTML={{ __html: style.css }}
      />
    ));

    return {
      ...initialProps,
      emotionStyleTags,
    };
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
          <link rel="icon" href="/static/img/icons/favicon.png" />

          {/* Inject Emotion styles for SSR - prevents FOUC */}
          {this.props.emotionStyleTags}

          {/* Preload Qanelas Soft font for loading screen */}
          <link
            rel="preload"
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
              @font-face {
                font-family: 'Qanelas Soft';
                src: url('/static/fonts/qanelas-soft/QanelasSoftRegular.otf') format('opentype');
                font-weight: 400;
                font-display: swap;
              }
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

          {/* Fallback displayed when JavaScript is disabled */}
          <noscript
            dangerouslySetInnerHTML={{
              __html: `
                <style type="text/css">
                  #PRE_REACT_LOADER {
                    display: none !important;
                  }
                  .noscript-overlay {
                    position: fixed;
                    inset: 0;
                    z-index: 999999;
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    justify-content: safe center;
                    overflow-y: auto;
                    background-color: #1e2529;
                    color: #ffffff;
                    font-family: 'Qanelas Soft', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
                    text-align: center;
                    padding: 24px;
                    box-sizing: border-box;
                  }
                  .noscript-card {
                    max-width: 540px;
                    width: 100%;
                    background-color: #263238;
                    border: 1px solid rgba(255, 255, 255, 0.12);
                    border-radius: 12px;
                    padding: 36px 28px;
                    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
                  }
                  .noscript-logo {
                    width: 64px;
                    height: 64px;
                    margin-bottom: 20px;
                  }
                  .noscript-title {
                    font-size: 1.5rem;
                    font-weight: 600;
                    color: #00B39F;
                    margin: 0 0 16px 0;
                  }
                  .noscript-desc {
                    font-size: 0.95rem;
                    line-height: 1.6;
                    color: #c4cdd5;
                    margin: 0 0 16px 0;
                  }
                  .noscript-instruction {
                    font-size: 0.95rem;
                    font-weight: 500;
                    color: #ffffff;
                    margin: 0 0 24px 0;
                  }
                  .noscript-links {
                    display: flex;
                    gap: 12px;
                    justify-content: center;
                    flex-wrap: wrap;
                  }
                  .noscript-btn {
                    display: inline-block;
                    padding: 8px 16px;
                    border-radius: 6px;
                    font-size: 0.875rem;
                    font-weight: 500;
                    text-decoration: none;
                    color: #00D3A9;
                    border: 1px solid rgba(0, 211, 169, 0.4);
                    background-color: rgba(0, 211, 169, 0.08);
                  }
                </style>
                <div class="noscript-overlay">
                  <div class="noscript-card">
                    <img
                      class="noscript-logo"
                      src="/static/img/meshery-logo/meshery-logo.png"
                      alt="Meshery Logo"
                    />
                    <h1 class="noscript-title">JavaScript is required</h1>
                    <p class="noscript-desc">
                      Meshery is an interactive cloud native management plane and requires JavaScript to load and function properly.
                    </p>
                    <p class="noscript-instruction">
                      Please enable JavaScript in your browser settings and reload this page to continue.
                    </p>
                    <div class="noscript-links">
                      <a
                        class="noscript-btn"
                        href="https://docs.meshery.io"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Documentation
                      </a>
                      <a
                        class="noscript-btn"
                        href="https://meshery.io"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Website
                      </a>
                      <a
                        class="noscript-btn"
                        href="https://github.com/meshery/meshery"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        GitHub
                      </a>
                    </div>
                  </div>
                </div>
              `,
            }}
          />
          {/* Pre-React script - must be sync to run before React hydration */}
          {/* eslint-disable-next-line @next/next/no-sync-scripts */}
          <script src="/loadingMessages.js"></script>

          <script async src="/static/wasm/wasm_exec.js"></script>

          <PureHtmlLoadingScreen id={'PRE_REACT_LOADER'} message="" />
          <Main />
          <NextScript />

          <script
            dangerouslySetInnerHTML={{
              __html: `
                (function () {
                  try {
                    if (window.Loader?.initializePreReactLoader) {
                      window.Loader.initializePreReactLoader();
                    }
                  } catch (e) {
                    console.log("Failed to initialize pre-react loader", e)
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
