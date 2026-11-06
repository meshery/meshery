import { describe, expect, it } from 'vitest';
import type { Command } from './commands';
import { filterCommands } from './filter';

describe('filterCommands', () => {
  const sampleCommands: Command[] = [
    {
      id: 'action:create-design',
      title: 'Create new Design',
      section: 'Actions',
      keywords: ['design', 'create', 'pattern', 'new', 'configurator'],
      href: '/configuration/designs/configurator',
      disabled: false,
    },
    {
      id: 'action:import-design',
      title: 'Import Design',
      section: 'Actions',
      keywords: ['design', 'import', 'upload', 'pattern'],
      href: '/configuration/designs?action=import',
      disabled: false,
    },
    {
      id: 'nav:/',
      title: 'Dashboard',
      section: 'Navigate',
      keywords: [],
      href: '/',
      disabled: false,
    },
    {
      id: 'nav:/configuration/designs',
      title: 'Designs',
      section: 'Navigate',
      subtitle: 'Configuration',
      keywords: ['Configuration'],
      href: '/configuration/designs',
      disabled: false,
    },
    {
      id: 'nav:/configuration/catalog',
      title: 'Catalog',
      section: 'Navigate',
      subtitle: 'Configuration',
      keywords: ['Configuration'],
      href: '/configuration/catalog',
      disabled: false,
    },
    {
      id: 'nav:/management/connections',
      title: 'Connections',
      section: 'Navigate',
      subtitle: 'Lifecycle',
      keywords: ['Lifecycle'],
      href: '/management/connections',
      disabled: false,
    },
    {
      id: 'settings',
      title: 'Settings',
      section: 'Settings',
      keywords: ['settings', 'configure', 'system'],
      href: '/settings',
      disabled: false,
    },
  ];

  it('returns all commands unchanged when query is empty or whitespace', () => {
    expect(filterCommands(sampleCommands, '')).toEqual(sampleCommands);
    expect(filterCommands(sampleCommands, '   ')).toEqual(sampleCommands);
  });

  it('requires EVERY token to match title, keywords, or subtitle', () => {
    // "config design" -> both tokens match "Designs" (subtitle: "Configuration", title: "Designs")
    const results = filterCommands(sampleCommands, 'config design');
    expect(results.some((c) => c.id === 'nav:/configuration/designs')).toBe(true);

    // "config nonexisting" -> one token fails, so excluded
    const failedResults = filterCommands(sampleCommands, 'config nonexisting');
    expect(failedResults).toHaveLength(0);
  });

  it('ranks title-prefix > title-word-prefix > substring', () => {
    // Query "des"
    // "Designs" has title-prefix (starts with "des") -> rank 3
    // "Import Design" and "Create new Design" have title-word-prefix ("Design" starts with "des") -> rank 2
    // "Catalog" with keyword "design" (if any) would be substring -> rank 1
    const results = filterCommands(sampleCommands, 'des');
    expect(results.length).toBeGreaterThan(1);
    expect(results[0].title).toBe('Designs');
    expect(results.slice(1).map((c) => c.title)).toContain('Create new Design');
    expect(results.slice(1).map((c) => c.title)).toContain('Import Design');
  });

  it('preserves stable order when ranks are tied', () => {
    // "Configuration" matches subtitle for Designs and Catalog with same rank
    const results = filterCommands(sampleCommands, 'lifecycle');
    expect(results).toHaveLength(1);
    expect(results[0].title).toBe('Connections');

    const configResults = filterCommands(sampleCommands, 'Configuration');
    // Designs was index 3, Catalog was index 4 in sampleCommands
    expect(configResults.map((c) => c.title)).toEqual(['Designs', 'Catalog']);
  });

  it('matches tokens across title, subtitle, and keywords', () => {
    // Matches subtitle
    expect(filterCommands(sampleCommands, 'Lifecycle')).toHaveLength(1);
    // Matches keyword
    expect(filterCommands(sampleCommands, 'upload')).toHaveLength(1);
    expect(filterCommands(sampleCommands, 'upload')[0].id).toBe('action:import-design');
    // Matches title
    expect(filterCommands(sampleCommands, 'Dashboard')).toHaveLength(1);
  });
});
