import { createTheme } from '@mui/material/styles';
import createCache from '@emotion/cache';

// Create a custom theme
const theme = createTheme({
  // Your theme customization here
});

// Create emotion cache
const createEmotionCache = () => {
  return createCache({ key: 'css', prepend: true });
};

function createPageContext() {
  return {
    theme,
    emotionCache: createEmotionCache(),
  };
}

let pageContext;

export default function getPageContext() {
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
