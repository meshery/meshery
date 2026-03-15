import { createTheme } from '@sistent/sistent';
import type { Theme } from '@sistent/sistent';
import createCache from '@emotion/cache';
import type { EmotionCache } from '@emotion/cache';

// Create a custom theme
const theme = createTheme({
  // Your theme customization here
});

// Create emotion cache
const createEmotionCache = (): EmotionCache => {
  return createCache({ key: 'css', prepend: true });
};

type PageContextType = {
  theme: Theme;
  emotionCache: EmotionCache;
};

function createPageContext(): PageContextType {
  return {
    theme,
    emotionCache: createEmotionCache(),
  };
}

let pageContext: PageContextType | undefined;

export default function getPageContext(): PageContextType {
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
