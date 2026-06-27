/**
 * Generates server/meshmodel/exoscale-icons from meshery.io SVG assets.
 * Each icon SVG already includes the Exoscale red (#da291c) background + white glyph.
 * svgComplete is set so Kanvas renders the full branded icon (fixes #15466).
 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '..');
const MODEL_VERSION = '1.0.0';
const MODEL_DIR = path.join(
  ROOT,
  'server',
  'meshmodel',
  'exoscale-icons',
  MODEL_VERSION,
  'v1.0.0',
);
const COMPONENTS_DIR = path.join(MODEL_DIR, 'components');
const BASE_URL =
  'https://raw.githubusercontent.com/meshery/meshery.io/master/assets/images/custom-integration/exoscale-icons/components';

const PRIMARY_COLOR = '#DA291C';
const SECONDARY_COLOR = '#000000';

const EMPTY_ID = '00000000-0000-0000-0000-000000000000';

function escapeSvgForJson(svg) {
  return svg.replace(/\r\n/g, '\n');
}

function toWhiteVariant(svg) {
  return svg
    .replace(/style="fill:#da291c"/gi, 'style="fill:transparent"')
    .replace(/style="fill:#fff"/gi, 'style="fill:#ffffff"');
}

function toPascalCase(name) {
  return name
    .split('-')
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join('');
}

function buildModelMetadata(svgColor, svgWhite, svgComplete = svgColor) {
  return {
    capabilities: null,
    isAnnotation: true,
    primaryColor: PRIMARY_COLOR,
    secondaryColor: SECONDARY_COLOR,
    shape: 'circle',
    styleOverrides: '',
    svgColor,
    svgComplete,
    svgWhite,
  };
}

function buildComponentJson({ fileName, svgRaw }) {
  const baseName = fileName.replace(/\.svg$/i, '');
  const kind = toPascalCase(baseName);
  const displayName = baseName
    .split('-')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');
  const svgColor = escapeSvgForJson(svgRaw);
  const svgWhite = escapeSvgForJson(toWhiteVariant(svgRaw));
  const svgComplete = svgColor;
  const modelMetadata = buildModelMetadata(svgColor, svgWhite, svgComplete);

  return {
    id: EMPTY_ID,
    schemaVersion: 'components.meshery.io/v1beta1',
    version: 'v1.0.0',
    displayName,
    description: `Exoscale ${displayName} icon`,
    format: 'JSON',
    model: {
      id: EMPTY_ID,
      schemaVersion: 'models.meshery.io/v1beta2',
      version: 'v1.0.0',
      name: 'exoscale-icons',
      displayName: 'Exoscale Icons',
      description: 'A curated collection of high-quality icons designed by Exoscale.',
      status: 'enabled',
      registrant: {
        id: EMPTY_ID,
        name: 'meshery',
        type: 'registry',
        subType: '',
        kind: 'meshery',
        status: 'discovered',
        created_at: '0001-01-01T00:00:00Z',
        updated_at: '0001-01-01T00:00:00Z',
        deleted_at: null,
        schemaVersion: '',
      },
      connection_id: EMPTY_ID,
      category: {
        id: EMPTY_ID,
        name: 'Platform',
      },
      subCategory: 'Cloud Provider',
      metadata: modelMetadata,
      model: { version: MODEL_VERSION },
      components_count: 0,
      relationships_count: 0,
      created_at: '0001-01-01T00:00:00Z',
      updated_at: '0001-01-01T00:00:00Z',
      components: null,
      relationships: null,
    },
    modelReference: {
      displayName: '',
      id: EMPTY_ID,
      model: { version: '' },
      name: '',
      registrant: { kind: '' },
      version: '',
    },
    styles: {
      'background-opacity': 0,
      data: { label: '' },
      primaryColor: PRIMARY_COLOR,
      secondaryColor: SECONDARY_COLOR,
      shape: 'circle',
      svgColor,
      svgComplete,
      svgWhite,
    },
    capabilities: [],
    status: 'enabled',
    metadata: {
      configurationUISchema: '',
      genealogy: '',
      instanceDetails: null,
      isAnnotation: true,
      isNamespaced: false,
      published: false,
    },
    configuration: null,
    component: {
      version: 'core.meshery.io/v1alpha1',
      kind,
      schema: '',
    },
    createdAt: '0001-01-01T00:00:00Z',
    updatedAt: '0001-01-01T00:00:00Z',
    deletedAt: null,
  };
}

async function fetchComponentList() {
  const res = await fetch(
    'https://api.github.com/repos/meshery/meshery.io/contents/assets/images/custom-integration/exoscale-icons/components?ref=master',
  );
  if (!res.ok) {
    throw new Error(`Failed to list components: ${res.status}`);
  }
  const items = await res.json();
  return items.filter((item) => item.name.endsWith('.svg')).map((item) => item.name);
}

async function fetchSvg(fileName) {
  const res = await fetch(`${BASE_URL}/${fileName}`);
  if (!res.ok) {
    throw new Error(`Failed to fetch ${fileName}: ${res.status}`);
  }
  return res.text();
}

async function main() {
  fs.mkdirSync(COMPONENTS_DIR, { recursive: true });

  const logoSvg = await fetch(
    'https://raw.githubusercontent.com/meshery/meshery.io/master/assets/images/custom-integration/exoscale-icons/icons/exoscale-logo.svg',
  ).then((r) => r.text());

  const logoColor = escapeSvgForJson(logoSvg);
  const logoWhite = escapeSvgForJson(toWhiteVariant(logoSvg));
  const modelMetadata = buildModelMetadata(logoColor, logoWhite, logoColor);

  const modelJson = {
    id: EMPTY_ID,
    schemaVersion: 'models.meshery.io/v1beta2',
    version: 'v1.0.0',
    name: 'exoscale-icons',
    displayName: 'Exoscale Icons',
    description: 'A curated collection of high-quality icons designed by Exoscale.',
    status: 'enabled',
    registrant: {
      id: EMPTY_ID,
      name: 'meshery',
      type: 'registry',
      subType: '',
      kind: 'meshery',
      status: 'discovered',
      created_at: '0001-01-01T00:00:00Z',
      updated_at: '0001-01-01T00:00:00Z',
      deleted_at: null,
      schemaVersion: '',
    },
    connection_id: EMPTY_ID,
    category: {
      id: EMPTY_ID,
      name: 'Platform',
    },
    subCategory: 'Cloud Provider',
    metadata: modelMetadata,
    model: { version: MODEL_VERSION },
    components_count: 0,
    relationships_count: 0,
    created_at: '0001-01-01T00:00:00Z',
    updated_at: '0001-01-01T00:00:00Z',
    components: null,
    relationships: null,
  };

  fs.writeFileSync(path.join(MODEL_DIR, 'model.json'), JSON.stringify(modelJson, null, 2));

  const files = await fetchComponentList();
  for (const fileName of files) {
    const svgRaw = await fetchSvg(fileName);
    const component = buildComponentJson({ fileName, svgRaw });
    const outName = `${fileName.replace(/\.svg$/i, '')}.json`;
    fs.writeFileSync(path.join(COMPONENTS_DIR, outName), JSON.stringify(component, null, 2));
    console.log(`Wrote ${outName}`);
  }

  console.log(`Generated ${files.length} components in ${MODEL_DIR}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
