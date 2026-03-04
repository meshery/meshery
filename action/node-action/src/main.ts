import * as core from '@actions/core'
import puppeteer from 'puppeteer'
import { PuppeteerScreenRecorder } from 'puppeteer-screen-recorder';

async function run(): Promise<void> {
  try {
    const authToken = core.getInput("authToken")
    await puppeteerBrowserHandler(authToken)
    core.setOutput('mardownResult', 'node ran successfully')
  } catch (error) {
    if (error instanceof Error) core.setFailed(error.message)
  }
}

run()

async function puppeteerBrowserHandler(authToken: string): Promise<void> {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  const recorder = new PuppeteerScreenRecorder(page);
  page.goto("https://playground.meshery.io");
  // Set screen size
  await page.setViewport({ width: 1080, height: 1024 });


  await recorder.stop();
  await browser.close()
}