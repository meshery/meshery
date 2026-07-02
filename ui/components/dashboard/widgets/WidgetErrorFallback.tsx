import { Button, ErrorIcon, Typography, useTheme } from '@sistent/sistent';
import { ErrorContainer } from '../style';

type WidgetErrorFallbackProps = {
  widgetTitle: string;
  message?: string;
  resetErrorBoundary?: () => void;
};

const WidgetErrorFallback = ({
  widgetTitle,
  message,
  resetErrorBoundary,
}: WidgetErrorFallbackProps) => {
  const theme = useTheme();

  return (
    <ErrorContainer>
      <ErrorIcon fill={theme.palette.background.error.default} />
      <Typography variant="h6" sx={{ color: theme.palette.text.error }}>
        Unable to load {widgetTitle}
      </Typography>
      <Typography variant="body1">
        {message ??
          'Unable to load. Confirm your organization selection and role assignments are appropriate and try again.'}
      </Typography>
      {resetErrorBoundary && (
        <Button
          variant="outlined"
          size="small"
          sx={{ marginTop: '1rem' }}
          onClick={resetErrorBoundary}
        >
          Try Again
        </Button>
      )}
    </ErrorContainer>
  );
};

export default WidgetErrorFallback;
