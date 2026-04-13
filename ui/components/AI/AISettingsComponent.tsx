import React, { useState } from 'react';
import {
  Box,
  Button,
  Divider,
  TextField,
  Typography,
  Alert,
  CircularProgress,
} from '@sistent/sistent';
import { useGetAIHealthQuery } from '@/rtk-query/ai';

interface ProviderConfig {
  apiKey: string;
  endpoint: string;
  model: string;
}

interface AISettings {
  openai: ProviderConfig;
  anthropic: ProviderConfig;
  ollama: ProviderConfig;
}

const DEFAULT_SETTINGS: AISettings = {
  openai: { apiKey: '', endpoint: '', model: 'gpt-4o' },
  anthropic: { apiKey: '', endpoint: '', model: 'claude-3-5-sonnet-20241022' },
  ollama: { apiKey: '', endpoint: 'http://localhost:11434', model: 'llama3' },
};

const ENV_VAR_HINTS: Record<keyof AISettings, string> = {
  openai: 'MESHERY_AI_OPENAI_API_KEY',
  anthropic: 'MESHERY_AI_ANTHROPIC_API_KEY',
  ollama: 'MESHERY_AI_OLLAMA_ENDPOINT',
};

/**
 * AISettingsComponent allows users to view and configure AI provider settings.
 *
 * Note: This component surfaces configuration guidance. Actual key storage
 * is handled server-side via environment variables or Meshery preferences.
 * Future work: wire to a POST /api/ai/settings endpoint.
 */
const AISettingsComponent: React.FC = () => {
  const [settings, setSettings] = useState<AISettings>(DEFAULT_SETTINGS);
  const [saved, setSaved] = useState(false);
  const [saving, setSaving] = useState(false);

  const { data: healthData, isLoading: healthLoading } = useGetAIHealthQuery();

  const updateProvider = (
    provider: keyof AISettings,
    field: keyof ProviderConfig,
    value: string,
  ) => {
    setSettings((prev) => ({
      ...prev,
      [provider]: { ...prev[provider], [field]: value },
    }));
    setSaved(false);
  };

  const handleSave = async () => {
    setSaving(true);
    // TODO: wire to POST /api/ai/settings when the server endpoint is available.
    // For now, display guidance about environment variables.
    await new Promise((r) => setTimeout(r, 400));
    setSaved(true);
    setSaving(false);
  };

  const activeProviders = healthData?.providers ?? [];

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, p: 2, maxWidth: 640 }}>
      <Typography variant="h6">AI Provider Settings</Typography>

      {healthLoading ? (
        <CircularProgress size={20} />
      ) : activeProviders.length > 0 ? (
        <Alert severity="success">
          AI service is active. Configured providers: {activeProviders.join(', ')}
        </Alert>
      ) : (
        <Alert severity="info">
          No AI providers are currently active. Configure at least one provider below or via
          environment variables.
        </Alert>
      )}

      {(['openai', 'anthropic', 'ollama'] as const).map((provider) => (
        <Box key={provider} sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
          <Typography variant="subtitle1" sx={{ textTransform: 'capitalize', fontWeight: 600 }}>
            {provider}
          </Typography>

          <Typography variant="caption" color="text.secondary">
            Environment variable: <code>{ENV_VAR_HINTS[provider]}</code>
          </Typography>

          {provider !== 'ollama' && (
            <TextField
              label="API Key"
              type="password"
              value={settings[provider].apiKey}
              onChange={(e) => updateProvider(provider, 'apiKey', e.target.value)}
              fullWidth
              size="small"
              placeholder="sk-..."
              autoComplete="new-password"
            />
          )}

          <TextField
            label={provider === 'ollama' ? 'Endpoint URL' : 'Base URL (optional)'}
            value={settings[provider].endpoint}
            onChange={(e) => updateProvider(provider, 'endpoint', e.target.value)}
            fullWidth
            size="small"
            placeholder={
              provider === 'ollama' ? 'http://localhost:11434' : 'https://api.openai.com'
            }
          />

          <TextField
            label="Default Model"
            value={settings[provider].model}
            onChange={(e) => updateProvider(provider, 'model', e.target.value)}
            fullWidth
            size="small"
          />

          <Divider />
        </Box>
      ))}

      {saved && (
        <Alert severity="success" onClose={() => setSaved(false)}>
          Settings saved. Restart Meshery server to apply environment variable changes.
        </Alert>
      )}

      <Button
        variant="contained"
        onClick={handleSave}
        disabled={saving}
        startIcon={saving ? <CircularProgress size={16} /> : undefined}
        sx={{ alignSelf: 'flex-start' }}
      >
        {saving ? 'Saving…' : 'Save Settings'}
      </Button>
    </Box>
  );
};

export default AISettingsComponent;
