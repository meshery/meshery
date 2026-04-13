import React, { useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
  Alert,
} from '@sistent/sistent';
import {
  useGetAIProvidersQuery,
  useAiDesignGenerateMutation,
  useGetAIHealthQuery,
} from '@/rtk-query/ai';

interface AIChatComponentProps {
  /** Called when a design has been successfully generated. */
  onDesignGenerated?: (design: object, warnings: string[]) => void;
}

/**
 * AIChatComponent provides a natural-language interface for generating
 * Meshery Designs. Users enter a description of their infrastructure and
 * select an AI provider; the component calls POST /api/ai/design/generate
 * and surfaces the result (or forwards it to the parent via onDesignGenerated).
 */
const AIChatComponent: React.FC<AIChatComponentProps> = ({ onDesignGenerated }) => {
  const [prompt, setPrompt] = useState('');
  const [selectedProvider, setSelectedProvider] = useState('');
  const [selectedModel, setSelectedModel] = useState('');
  const [result, setResult] = useState<{ design: object; warnings: string[] } | null>(null);

  const { data: healthData } = useGetAIHealthQuery();
  const { data: providersData, isLoading: loadingProviders } = useGetAIProvidersQuery();
  const [generateDesign, { isLoading: generating, error: generateError }] =
    useAiDesignGenerateMutation();

  const isAIAvailable = healthData?.status === 'ok';
  const providers = providersData?.providers ?? [];

  const handleGenerate = async () => {
    if (!prompt.trim() || !selectedProvider) return;

    try {
      const response = await generateDesign({
        prompt: prompt.trim(),
        provider: selectedProvider,
        model: selectedModel || undefined,
      }).unwrap();

      const warnings = response.warnings ?? [];
      setResult({ design: response.design, warnings });
      onDesignGenerated?.(response.design, warnings);
    } catch {
      // error is exposed via generateError from RTK-Query
    }
  };

  if (!isAIAvailable) {
    return (
      <Box sx={{ p: 2 }}>
        <Alert severity="info">
          AI service is not configured. Configure an AI provider (OpenAI, Anthropic, or Ollama) to
          enable natural-language design generation.
        </Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, p: 2 }}>
      <Typography variant="h6">Generate Design from Natural Language</Typography>

      <TextField
        label="Describe your infrastructure"
        placeholder="e.g. Deploy an nginx web server with a Redis cache in Kubernetes"
        multiline
        minRows={3}
        maxRows={8}
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        fullWidth
        disabled={generating}
      />

      <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
        <FormControl sx={{ minWidth: 160 }} disabled={generating || loadingProviders}>
          <InputLabel id="ai-provider-label">AI Provider</InputLabel>
          <Select
            labelId="ai-provider-label"
            value={selectedProvider}
            label="AI Provider"
            onChange={(e) => setSelectedProvider(e.target.value as string)}
          >
            {providers.map((p) => (
              <MenuItem key={p.name} value={p.name}>
                {p.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <TextField
          label="Model (optional)"
          placeholder="e.g. gpt-4o"
          value={selectedModel}
          onChange={(e) => setSelectedModel(e.target.value)}
          sx={{ minWidth: 180 }}
          disabled={generating}
        />
      </Box>

      {generateError && (
        <Alert severity="error">
          {'data' in generateError
            ? String((generateError as { data: unknown }).data)
            : 'Failed to generate design. Check your provider configuration.'}
        </Alert>
      )}

      {result && result.warnings.length > 0 && (
        <Alert severity="warning">
          Design generated with warnings:
          <ul style={{ margin: 0, paddingLeft: 16 }}>
            {result.warnings.map((w, i) => (
              <li key={i}>{w}</li>
            ))}
          </ul>
        </Alert>
      )}

      {result && result.warnings.length === 0 && (
        <Alert severity="success">Design generated successfully.</Alert>
      )}

      <Box sx={{ display: 'flex', gap: 1 }}>
        <Button
          variant="contained"
          onClick={handleGenerate}
          disabled={!prompt.trim() || !selectedProvider || generating}
          startIcon={generating ? <CircularProgress size={16} /> : undefined}
        >
          {generating ? 'Generating…' : 'Generate Design'}
        </Button>

        {result && (
          <Button
            variant="outlined"
            onClick={() => {
              setResult(null);
              setPrompt('');
            }}
          >
            Clear
          </Button>
        )}
      </Box>
    </Box>
  );
};

export default AIChatComponent;
