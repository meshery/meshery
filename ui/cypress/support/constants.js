import { id } from "./helpers";

export const DESIGNER = "designer"
export const VISUALIZER = "visualizer"
export const extension = {
  path: "/api/provider/extension*",
  alias: "extensionFileLoad"
};

export const designEndpoint = {
  path: "/api/pattern*",
  alias: "designEp",
  wait: "@designEp",
  absolutePath: "/api/pattern"
}

export const MESHMAP_PATH = "/extension/meshmap";

export const CANVAS_CONTAINER_ID = "cy-canvas-container"

export const TIME = {
  SMALL: 500,
  MEDIUM: 1000,
  LARGE: 1500,
  XLARGE: 2000
}

export const canvasContainer = {
  query: id(CANVAS_CONTAINER_ID),
  alias: "canvas"
}

/**
 * Selection and general Event Binding Layer
 */
export const canvasLayer0 = {
  query: '[data-id="layer0-selectbox"]',
  alias: "layer0"
}

/**
 * drag and drop Layer
 */
export const canvasLayer1 = {
  query: '[data-id="layer1-drag"]',
  alias: "layer1"
}

/**
 * Node and Element Layer
 */
export const canvaslayer2 = {
  query: '[data-id="layer2-node"]',
  alias: "layer2"
}
