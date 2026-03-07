// @ts-nocheck
import { createTheme } from '@sistent/sistent';
import createEmotionCache from '../lib/createEmotionCache';

// Create a custom theme
const theme = createTheme({
  // Your theme customization here
});

function createPageContext(emotionCache = createEmotionCache()) {
  return {
    theme,
    emotionCache,
  };
}

let pageContext;

export default function getPageContext(emotionCache) {
  // Make sure to create a new context for every server-side request
  if (typeof window === 'undefined') {
    return createPageContext(emotionCache);
  }

  // Reuse context on the client-side
  if (!pageContext) {
    pageContext = createPageContext(emotionCache);
  }

  return pageContext;
}
