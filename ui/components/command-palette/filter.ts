import type { Command } from './commands';

/**
 * Filters and ranks commands based on user query string.
 *
 * Rules:
 * - trim + lowercase
 * - split on whitespace
 * - EVERY token must match title, keywords, or subtitle (substring match)
 * - Ranking: title-prefix of the full query > title-word-prefix > substring
 * - Stable order otherwise
 * - Empty/whitespace query returns all commands unchanged
 */
export function filterCommands(commands: Command[] = [], query: string = ''): Command[] {
  const q = (query || '').trim().toLowerCase();
  if (!q) {
    return commands;
  }

  const tokens = q.split(/\s+/).filter(Boolean);
  if (tokens.length === 0) {
    return commands;
  }

  const matches: { cmd: Command; rank: number; index: number }[] = [];

  for (let i = 0; i < commands.length; i++) {
    const cmd = commands[i];
    const title = (cmd.title || '').toLowerCase();
    const subtitle = (cmd.subtitle || '').toLowerCase();
    const keywords = (cmd.keywords || []).map((k) => k.toLowerCase());

    // EVERY token must match title, keywords or subtitle (substring)
    const matchesAllTokens = tokens.every((token) => {
      return (
        title.includes(token) || subtitle.includes(token) || keywords.some((k) => k.includes(token))
      );
    });

    if (!matchesAllTokens) {
      continue;
    }

    // Rank: title-prefix of the full query > title-word-prefix > substring
    let rank = 1;
    const titleWords = title.split(/\s+/).filter(Boolean);

    if (title.startsWith(q)) {
      rank = 3;
    } else if (titleWords.some((word) => word.startsWith(q))) {
      rank = 2;
    }

    matches.push({ cmd, rank, index: i });
  }

  // Stable sort: higher rank first, fallback to original index
  matches.sort((a, b) => b.rank - a.rank || a.index - b.index);

  return matches.map((m) => m.cmd);
}
