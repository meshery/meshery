import { describe, expect, it } from 'vitest';
import { TABS_SCROLL_BUTTONS_CLASS } from './constants';

describe('dashboard constants', () => {
  it('exports the MUI Tabs scroll buttons class literal', () => {
    expect(TABS_SCROLL_BUTTONS_CLASS).toBe('MuiTabs-scrollButtons');
  });
});
