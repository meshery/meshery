import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  Typography,
} from '@mui/material';
import ConnectionLoader from './ConnectionLoader';
import useConnectionLoader from '../hooks/useConnectionLoader';

const ConnectionRegistration = ({ open, onClose, onSuccess }) => {
  const [formData, setFormData] = useState({
    name: '',
    kind: '',
    type: '',
    subType: '',
    url: '',
    credentials: {},
  });
  
  const { isVerifying, verificationMessage, startVerification } = useConnectionLoader();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Start verification process
    startVerification('Verifying connection...');
    
    try {
      const response = await fetch('/api/integrations/connections', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...formData,
          status: 'discovered',
        }),
      });

      if (response.ok) {
        onSuccess?.();
        onClose();
      } else {
        const error = await response.text();
        console.error('Connection registration failed:', error);
      }
    } catch (error) {
      console.error('Connection registration error:', error);
    }
  };

  const handleInputChange = (field) => (event) => {
    setFormData(prev => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  if (isVerifying) {
    return (
      <Dialog open={open} maxWidth="sm" fullWidth>
        <DialogContent>
          <ConnectionLoader message={verificationMessage} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Register New Connection</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Connection Name"
              value={formData.name}
              onChange={handleInputChange('name')}
              required
              fullWidth
            />
            
            <FormControl fullWidth required>
              <InputLabel>Connection Kind</InputLabel>
              <Select
                value={formData.kind}
                onChange={handleInputChange('kind')}
                label="Connection Kind"
              >
                <MenuItem value="kubernetes">Kubernetes</MenuItem>
                <MenuItem value="prometheus">Prometheus</MenuItem>
                <MenuItem value="grafana">Grafana</MenuItem>
                <MenuItem value="meshery">Meshery</MenuItem>
              </Select>
            </FormControl>

            <TextField
              label="Connection URL"
              value={formData.url}
              onChange={handleInputChange('url')}
              fullWidth
              placeholder="https://example.com"
            />

            <FormControl fullWidth>
              <InputLabel>Connection Type</InputLabel>
              <Select
                value={formData.type}
                onChange={handleInputChange('type')}
                label="Connection Type"
              >
                <MenuItem value="platform">Platform</MenuItem>
                <MenuItem value="management">Management</MenuItem>
                <MenuItem value="monitoring">Monitoring</MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel>Connection Sub Type</InputLabel>
              <Select
                value={formData.subType}
                onChange={handleInputChange('subType')}
                label="Connection Sub Type"
              >
                <MenuItem value="orchestrator">Orchestrator</MenuItem>
                <MenuItem value="management">Management</MenuItem>
                <MenuItem value="monitoring">Monitoring</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button 
            type="submit" 
            variant="contained" 
            color="primary"
            disabled={!formData.name || !formData.kind}
          >
            Register Connection
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default ConnectionRegistration; 