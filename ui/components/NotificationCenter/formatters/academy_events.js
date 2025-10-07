import React from 'react';
import {
  Box,
  Divider,
  Link,
  List,
  ListItem,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from '@sistent/sistent';

export const AcademyEventsFormatter = ({ event }) => {
  const result = event.metadata.result;
  const quiz = result.quiz;
  if (!result || !quiz || !quiz.parent) {
    return (
      <Typography variant="body1" color="error">
        Could not load quiz results. The event data is incomplete.
      </Typography>
    );
  }

  const quizDetails = [
    { label: 'Chapter', value: quiz.parent.title },
    { label: 'Description', value: quiz.description },
    { label: 'Time limit', value: quiz.time_limit },
    { label: 'Attempted at', value: new Date(result.attempted_at).toLocaleString() },
  ];

  return (
    <Box>
      <Typography variant="h6">Quiz details</Typography>
      <List disablePadding>
        <ListItem sx={{ py: 0.5 }}>
          <Link variant="body1" target="_blank" rel="noopener noreferrer" href={quiz.permalink}>
            Go to quiz
          </Link>
        </ListItem>
        {quizDetails.map(({ label, value }) => (
          <ListItem key={label} sx={{ py: 0.5 }}>
            <Typography variant="body1">
              <Box component="span" fontWeight="bold">
                {label}:
              </Box>{' '}
              {value || (
                <Box component="span" fontStyle="italic">
                  none
                </Box>
              )}
            </Typography>
          </ListItem>
        ))}
      </List>

      <Divider sx={{ my: 2 }} />

      <Typography variant="h6">Result</Typography>

      <List disablePadding>
        <ListItem sx={{ py: 0.5 }} divider>
          <Box>
            <Typography variant="body1">
              {`Score: ${Number(result.percentage_scored).toFixed(2)}% (Pass mark: ${result.pass_percentage}%)`}
            </Typography>
            <Typography variant="body1" color={result.passed ? 'green' : 'red'}>
              {result.passed ? '✅ Passed' : '❌ Not passed'}
            </Typography>
          </Box>
        </ListItem>
        <ListItem sx={{ display: 'block', py: 0.5 }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Question</TableCell>
                <TableCell>Result</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {Object.entries(result.correct_submissions).map(([question, isCorrect]) => (
                <TableRow key={question}>
                  <TableCell>{question}</TableCell>
                  <TableCell>{isCorrect ? '✅' : '❌'}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </ListItem>
      </List>
    </Box>
  );
};
