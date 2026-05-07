import { createTheme } from '@sistent/sistent';
import createCache, { EmotionCache } from '@emotion/cache';

type PageContext = {
  theme: ReturnType<typeof createTheme>;
  emotionCache: EmotionCache;
};

// Create a custom theme
const theme = createTheme({
  // Your theme customization here
});

// Create emotion cache
const createEmotionCache = () => {
  return createCache({ key: 'css', prepend: true });
};

function createPageContext(): PageContext {
  return {
    theme,
    emotionCache: createEmotionCache(),
  };
}

let pageContext: PageContext | undefined;

export default function getPageContext(): PageContext {
  // Make sure to create a new context for every server-side request
  if (typeof window === 'undefined') {
    return createPageContext();
  }

  // Reuse context on the client-side
  if (!pageContext) {
    pageContext = createPageContext();
  }

  return pageContext;
}
