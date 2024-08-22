import { test as base } from '@playwright/test';
import { MeshMapTutorial } from './meshmap_tutorial';

/*
  This is safe place if you wanted to skip some process
  and not including it on videos, make sure to just extend the behavior here
  and called it inside initSetup and dont forget
  to cleanup the behavior inside cleanUp function
*/
export class Recorder extends MeshMapTutorial {
  constructor(page) {
    super(page);
  }

  async initSetup() {
    await super.skipTutorial();
  }

  async cleanUp() {
    await super.turnOnTutorialSettings();
  }
}

export const test = base.extend({
  recorderPage: [
    async ({ browser }, use) => {
      const initBrowser = await browser.newPage();

      const recorderPage = new Recorder(initBrowser);
      await recorderPage.goto();
      await recorderPage.initSetup();

      await use();

      await recorderPage.cleanUp();

      await initBrowser.close();
    },
    { scope: 'worker' },
  ],
});

export { expect } from '@playwright/test';
