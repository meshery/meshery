import { shuffle } from 'lodash';
import { Designs } from '../samples/seededDesigns';
import { test } from './fixtures/pages';

/**
 * @typedef {import("@playwright/test").Page} Page
 * @typedef {import("./fixtures/pages").DesignerPage} DesignerPage
 * @typedef {import("./fixtures/pages").Canvas} Canvas
 */

const randomDesigns = shuffle(Object.values(Designs));

randomDesigns.forEach(({ id, name, expectations }) => {
  // This test will load the design from the route
  test(`Load Design ${name} from route` /** @param {{ page: Page, designerPage: DesignerPage }} */, async ({
    designsPage,
  }) => {
    await designsPage.loadDesign(id);

    if (!expectations) {
      console('✔ Design Loaded Succesfully , no more expectations defined');
      return;
    }
  });
});
