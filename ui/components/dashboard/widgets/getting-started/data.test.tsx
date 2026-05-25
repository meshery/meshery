import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://cloud.meshery.io',
}));

vi.mock('@sistent/sistent', () => ({
  Link: ({ children, href }: { children?: React.ReactNode; href?: string }) => (
    <a href={href}>{children}</a>
  ),
}));

import { ActionName, stepsData } from './data';

describe('ActionName constants', () => {
  it('exposes the expected named actions', () => {
    expect(ActionName).toMatchObject({
      NEXT: 'Next',
      SETUP: 'Setup',
      SHARE: 'Share A Design',
      DONE: 'Done',
      EXPLORE: 'Explore',
      CATALOG: 'Choose Template',
      CREATE: 'From Scratch',
      IMPORT: 'Import',
      LEARN: 'Learn',
    });
  });
});

describe('stepsData', () => {
  it('contains the 5 documented steps in order', () => {
    expect(stepsData.map((s) => s.id)).toEqual([1, 2, 3, 4, 5]);
  });

  it('each step has a title, subTitle, and a non-empty journey', () => {
    stepsData.forEach((step) => {
      expect(typeof step.title).toBe('string');
      expect(typeof step.subTitle).toBe('string');
      expect(Array.isArray(step.journey)).toBe(true);
      expect(step.journey.length).toBeGreaterThan(0);
    });
  });

  it('uses MESHERY_CLOUD_PROD in the GitHub setup action', () => {
    // Step 2: GitOps your infra with Kanvas Snapshots
    const gitHubAction = stepsData[1].journey[0].primaryAction;
    expect(typeof gitHubAction).toBe('function');

    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (gitHubAction as () => void)();
    expect(openSpy).toHaveBeenCalledWith('https://cloud.meshery.io/connect/github/new', '_blank');
    openSpy.mockRestore();
  });

  it('routes the explore action to /extension/meshmap?mode=design', () => {
    // Step 1's "Explore" actionable journey is the second entry.
    const action = stepsData[0].journey[1].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith('/extension/meshmap?mode=design', '_self');
    openSpy.mockRestore();
  });

  it('the "Share A Design" action navigates to /catalog/content/my-designs', () => {
    const action = stepsData[2].journey[0].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith('/catalog/content/my-designs', '_self');
    openSpy.mockRestore();
  });

  it('the "From Scratch" action navigates to the meshmap design extension', () => {
    const action = stepsData[3].journey[1].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith('/extension/meshmap?mode=design', '_self');
    openSpy.mockRestore();
  });

  it('the "Choose Template" action navigates to the configuration catalog', () => {
    const action = stepsData[3].journey[2].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith('/configuration/catalog', '_self');
    openSpy.mockRestore();
  });

  it('the import action navigates to /configuration/designs', () => {
    const action = stepsData[3].journey[3].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith('/configuration/designs', '_self');
    openSpy.mockRestore();
  });

  it('the "Learn" action opens the learning path on cloud', () => {
    const action = stepsData[4].journey[0].primaryAction;
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);
    (action as () => void)();
    expect(openSpy).toHaveBeenCalledWith(
      'https://cloud.meshery.io/academy/learning-paths/11111111-1111-1111-1111-111111111111/mastering-meshery',
      '_blank',
    );
    openSpy.mockRestore();
  });
});
