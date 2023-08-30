import { useState, useEffect } from "react";

/**
 * getWindowDimensions - Returns the width and height of the window
 * @returns {object} {width, height}
 */
export function getWindowDimensions() {
  const { innerWidth : width, innerHeight : height } = window;
  return {
    width,
    height
  };
}

/**
 * useWindowDimensions - Returns the width and height of the window
 * @returns {object} {width, height}
 * @example
 * const { width, height } = useWindowDimensions();
 * console.log(width, height);
 * // 1920 1080
 */
export function useWindowDimensions() {
  const [windowDimensions, setWindowDimensions] = useState(getWindowDimensions());

  useEffect(() => {
    let searchTimeout;

    function handleResize() {
      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }
      searchTimeout = setTimeout(() => {
        setWindowDimensions(getWindowDimensions());
      }, 500);
    }

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return windowDimensions;
}