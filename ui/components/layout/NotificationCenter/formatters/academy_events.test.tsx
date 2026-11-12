import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Divider: () => <hr data-testid="divider" />,
  Link: ({ children, href, ...props }: any) => (
    <a data-testid="quiz-link" href={href} {...props}>
      {children}
    </a>
  ),
  List: ({ children, ...props }: any) => <ul {...props}>{children}</ul>,
  ListItem: ({ children, ...props }: any) => <li {...props}>{children}</li>,
  Table: ({ children, ...props }: any) => <table {...props}>{children}</table>,
  TableBody: ({ children, ...props }: any) => <tbody {...props}>{children}</tbody>,
  TableCell: ({ children, ...props }: any) => <td {...props}>{children}</td>,
  TableHead: ({ children, ...props }: any) => <thead {...props}>{children}</thead>,
  TableRow: ({ children, ...props }: any) => <tr {...props}>{children}</tr>,
  Typography: ({ children, color, variant, ...props }: any) => (
    <span data-color={color} data-variant={variant} {...props}>
      {children}
    </span>
  ),
}));

import { AcademyEventsFormatter } from './academy_events';

const buildEvent = (overrides: any = {}) => ({
  metadata: {
    result: {
      quiz: {
        parent: { title: 'Service Mesh 101' },
        description: 'Quiz on mesh basics',
        time_limit: '15m',
        permalink: 'https://academy.example/quiz/1',
      },
      attempted_at: '2024-01-02T10:00:00Z',
      percentage_scored: 85.5,
      pass_percentage: 70,
      passed: true,
      correct_submissions: {
        'What is a service mesh?': true,
        'Is the sky blue?': false,
      },
      ...overrides,
    },
  },
});

describe('AcademyEventsFormatter', () => {
  it('renders an error fallback when quiz data is incomplete', () => {
    render(
      <AcademyEventsFormatter
        event={{ metadata: { result: { quiz: { parent: undefined as any } } } } as any}
      />,
    );

    expect(screen.getByText(/could not load quiz results/i)).toBeInTheDocument();
  });

  it('renders the quiz details, score, and per-question table when complete', () => {
    render(<AcademyEventsFormatter event={buildEvent() as any} />);

    // Heading
    expect(screen.getByText('Quiz details')).toBeInTheDocument();
    // Link
    expect(screen.getByTestId('quiz-link')).toHaveAttribute(
      'href',
      'https://academy.example/quiz/1',
    );
    // Quiz detail entries
    expect(screen.getByText('Chapter:')).toBeInTheDocument();
    expect(screen.getByText('Service Mesh 101')).toBeInTheDocument();
    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('Quiz on mesh basics')).toBeInTheDocument();
    expect(screen.getByText('Time limit:')).toBeInTheDocument();
    expect(screen.getByText('15m')).toBeInTheDocument();
    expect(screen.getByText('Attempted at:')).toBeInTheDocument();
    // Result section
    expect(screen.getByText(/Score: 85\.50%/)).toBeInTheDocument();
    expect(screen.getByText(/Pass mark: 70%/)).toBeInTheDocument();
    // Question rows
    expect(screen.getByText('What is a service mesh?')).toBeInTheDocument();
    expect(screen.getByText('Is the sky blue?')).toBeInTheDocument();
  });

  it('shows an italicised "none" placeholder for missing detail values', () => {
    render(
      <AcademyEventsFormatter
        event={
          {
            metadata: {
              result: {
                quiz: {
                  parent: { title: 'Chapter X' },
                  description: undefined,
                  time_limit: undefined,
                  permalink: '',
                },
                attempted_at: 0,
                percentage_scored: 0,
                pass_percentage: 50,
                passed: false,
                correct_submissions: {},
              },
            },
          } as any
        }
      />,
    );

    expect(screen.getAllByText('none').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText(/Not passed/)).toBeInTheDocument();
  });
});
