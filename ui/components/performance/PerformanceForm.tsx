import React from 'react';
import {
  AccordionDetails,
  AccordionSummary,
  Autocomplete,
  Button,
  CustomTooltip,
  FormControlLabel,
  FormLabel,
  Grid2,
  MenuItem,
  RadioGroup,
  TextField,
  Typography,
} from '@sistent/sistent';
import {
  ExpandMore as ExpandMoreIcon,
  HelpOutlineOutlined as HelpOutlineOutlinedIcon,
} from '@/assets/icons';
import { durationOptions } from '../../lib/prePopulatedOptions';
import { CustomTextTooltip } from '../meshery-mesh-interface/PatternService/CustomTextTooltip';
import { ExpansionPanelComponent, FormContainer, HelpIcon, RadioButton } from './style';
import {
  infoCRTCertificates,
  infoFlags,
  infoloadGenerators,
  loadGenerators,
} from './performance-helpers';

interface PerformanceFormProps {
  profileName: string;
  meshName: string;
  selectedMesh: string;
  meshModels: string[];
  url: string;
  urlError: boolean;
  c: string | number;
  qps: string | number;
  t: string;
  tValue: string;
  tError: string;
  headers: string;
  cookies: string;
  contentType: string;
  reqBody: string;
  additionalOptions: string;
  jsonError: boolean;
  caCertificate: { name?: string; file?: string };
  metadata: any;
  loadGenerator: string;
  handleChange: (name: string) => (event: any) => void;
  handleDurationChange: (event: any, newValue: any) => void;
  handleInputDurationChange: (event: any, newValue: any) => void;
  handleCertificateUpload: (event: any) => void;
}

/**
 * PerformanceForm renders the full performance-test configuration form: the
 * basic fields (name/mesh/url/concurrency/qps/duration), the Advanced Options
 * accordion (headers/cookies/content type/body/additional options/CA cert),
 * and the load-generator radio group.
 *
 * Extracted from `performance/index.tsx` in Phase 5.a so the entry point
 * stays under the 600-line size budget. All state and handlers live in the
 * parent and are passed through; no behaviour changes.
 */
const PerformanceForm: React.FC<PerformanceFormProps> = ({
  profileName,
  meshName,
  selectedMesh,
  meshModels,
  url,
  urlError,
  c,
  qps,
  t,
  tValue,
  tError,
  headers,
  cookies,
  contentType,
  reqBody,
  additionalOptions,
  jsonError,
  caCertificate,
  metadata,
  loadGenerator,
  handleChange,
  handleDurationChange,
  handleInputDurationChange,
  handleCertificateUpload,
}) => {
  return (
    <Grid2 container spacing={1} size="grow">
      <Grid2 size={{ xs: 12, md: 6 }}>
        <TextField
          id="profileName"
          name="profileName"
          label="Profile Name"
          fullWidth
          value={profileName}
          margin="normal"
          variant="outlined"
          onChange={handleChange('profileName')}
          inputProps={{
            maxLength: 300,
          }}
          InputProps={{
            endAdornment: (
              <CustomTooltip title="Create a profile providing a name, if a profile name is not provided, a random one will be generated for you.">
                <HelpOutlineOutlinedIcon style={{ color: '#929292' }} />
              </CustomTooltip>
            ),
          }}
        />
      </Grid2>

      <Grid2 size={{ xs: 12, md: 6 }}>
        <TextField
          select
          id="meshName"
          name="meshName"
          label="Technology"
          fullWidth
          value={meshName === '' && selectedMesh !== '' ? selectedMesh : meshName}
          margin="normal"
          variant="outlined"
          onChange={handleChange('meshName')}
        >
          <MenuItem key="mh_-_none" value="None">
            None
          </MenuItem>
          {meshModels &&
            meshModels.map((mesh) => (
              <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>
                {mesh}
              </MenuItem>
            ))}
        </TextField>
      </Grid2>
      <Grid2 size={{ xs: 12 }}>
        <TextField
          required
          id="url"
          name="url"
          label="URL to test"
          type="url"
          fullWidth
          value={url}
          error={urlError}
          helperText={urlError ? 'Please enter a valid URL along with protocol' : ''}
          margin="normal"
          variant="outlined"
          onChange={handleChange('url')}
          InputProps={{
            endAdornment: (
              <CustomTooltip title="The Endpoint where the load will be generated and the performance test will run against.">
                <HelpOutlineOutlinedIcon style={{ color: '#929292' }} />
              </CustomTooltip>
            ),
          }}
        />
      </Grid2>
      <Grid2 size={{ xs: 12, md: 4 }}>
        <TextField
          required
          id="c"
          name="c"
          label="Concurrent requests"
          type="number"
          fullWidth
          value={c}
          inputProps={{ min: '0', step: '1' }}
          margin="normal"
          variant="outlined"
          onChange={handleChange('c')}
          InputLabelProps={{ shrink: true }}
          InputProps={{
            endAdornment: (
              <CustomTooltip title="Load Testing tool will create this many concurrent request against the endpoint.">
                <HelpOutlineOutlinedIcon style={{ color: '#929292' }} />
              </CustomTooltip>
            ),
          }}
        />
      </Grid2>
      <Grid2 size={{ xs: 12, md: 4 }}>
        <TextField
          required
          id="qps"
          name="qps"
          label="Queries per second"
          type="number"
          fullWidth
          value={qps}
          inputProps={{ min: '0', step: '1' }}
          margin="normal"
          variant="outlined"
          onChange={handleChange('qps')}
          InputLabelProps={{ shrink: true }}
          InputProps={{
            endAdornment: (
              <CustomTooltip title="The Number of queries/second. If not provided then the MAX number of queries/second will be requested">
                <HelpOutlineOutlinedIcon style={{ color: '#929292' }} />
              </CustomTooltip>
            ),
          }}
        />
      </Grid2>
      <Grid2 size={{ xs: 12, md: 4 }}>
        <CustomTooltip
          title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}
        >
          <Autocomplete
            required
            id="t"
            name="t"
            freeSolo
            label="Duration*"
            fullWidth
            variant="outlined"
            classes={{ root: tError }}
            value={tValue}
            inputValue={t}
            onChange={handleDurationChange}
            onInputChange={handleInputDurationChange}
            options={durationOptions}
            style={{ marginTop: '16px', marginBottom: '8px' }}
            renderInput={(params) => <TextField {...params} label="Duration*" variant="outlined" />}
            InputProps={{
              endAdornment: (
                <CustomTooltip title="Default duration is 30 seconds">
                  <HelpOutlineOutlinedIcon style={{ color: '#929292' }} />
                </CustomTooltip>
              ),
            }}
          />
        </CustomTooltip>
      </Grid2>
      <Grid2 size={{ xs: 12, md: 12 }}>
        <ExpansionPanelComponent>
          <AccordionSummary expanded={true} expandIcon={<ExpandMoreIcon />}>
            <Typography align="center" color="textSecondary" variant="h6">
              Advanced Options
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Grid2 container spacing={1} size="grow">
              <Grid2 size={{ xs: 12 }}>
                <TextField
                  id="headers"
                  name="headers"
                  label='Request Headers e.g. {"host":"bookinfo.meshery.io"}'
                  fullWidth
                  value={headers}
                  multiline
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('headers')}
                ></TextField>
              </Grid2>
              <Grid2 size={{ xs: 12 }}>
                <TextField
                  id="cookies"
                  name="cookies"
                  label='Request Cookies e.g. {"yummy_cookie":"choco_chip"}'
                  fullWidth
                  value={cookies}
                  multiline
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('cookies')}
                ></TextField>
              </Grid2>
              <Grid2 size={{ xs: 12 }}>
                <TextField
                  id="contentType"
                  name="contentType"
                  label="Content Type e.g. application/json"
                  fullWidth
                  value={contentType}
                  multiline
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('contentType')}
                ></TextField>
              </Grid2>
              <Grid2 size={{ xs: 12 }}>
                <TextField
                  id="cookies"
                  name="cookies"
                  label='Request Body e.g. {"method":"post","url":"http://bookinfo.meshery.io/test"}'
                  fullWidth
                  value={reqBody}
                  multiline
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('reqBody')}
                ></TextField>
              </Grid2>
              <Grid2 container size="grow">
                <Grid2 size={{ xs: 6 }}>
                  <TextField
                    id="additional_options"
                    name="additional_options"
                    label="Additional Options e.g. { `requestPerSecond`: 20 }"
                    fullWidth
                    error={jsonError}
                    helperText={jsonError ? 'Please enter a valid JSON string' : ''}
                    value={
                      additionalOptions.length > 150
                        ? `${additionalOptions.slice(0, 150)} .....`
                        : additionalOptions
                    }
                    multiline
                    margin="normal"
                    variant="outlined"
                    size="small"
                    onChange={handleChange('additional_options')}
                  />
                </Grid2>
                <Grid2 size={{ xs: 6 }}>
                  <label
                    htmlFor="upload-additional-options"
                    style={{ paddingLeft: '0.7rem', paddingTop: '8px' }}
                    fullWidth
                  >
                    <Button
                      variant="outlined"
                      onChange={handleChange('additional_options')}
                      aria-label="Upload Button"
                      component="span"
                      style={{ margin: '0.5rem', marginTop: '1.15rem' }}
                    >
                      <input
                        id="upload-additional-options"
                        type="file"
                        accept={'.json'}
                        name="upload-button"
                        hidden
                        data-cy="additional-options-upload-button"
                      />
                      Browse
                    </Button>
                    <CustomTooltip title={infoFlags} interactive>
                      <HelpIcon />
                    </CustomTooltip>
                  </label>
                </Grid2>
              </Grid2>
              <Grid2 container size="grow">
                <Grid2 size={{ xs: 6 }}>
                  <TextField
                    size="small"
                    variant="outlined"
                    margin="mormal"
                    fullWidth
                    label={caCertificate?.name || 'Upload SSL Certificate e.g. .crt file'}
                    style={{ width: '100%', margin: '0.5rem 0' }}
                    value={metadata?.ca_certificate.name}
                  />
                </Grid2>
                <Grid2 size={{ xs: 6 }}>
                  <label
                    htmlFor="upload-cacertificate"
                    style={{ paddingLeft: '0.7rem', paddingTop: '8px' }}
                  >
                    <Button
                      variant="outlined"
                      aria-label="Upload Button"
                      onChange={handleChange('caCertificate')}
                      component="span"
                      style={{ margin: '0.5rem' }}
                    >
                      <input
                        id="upload-cacertificate"
                        type="file"
                        accept={'.crt'}
                        name="upload-button"
                        hidden
                        data-cy="cacertificate-upload-button"
                        onChange={handleCertificateUpload}
                      />
                      Browse
                    </Button>
                    <CustomTooltip title={infoCRTCertificates} interactive>
                      <HelpIcon />
                    </CustomTooltip>
                  </label>
                </Grid2>
              </Grid2>
            </Grid2>
          </AccordionDetails>
        </ExpansionPanelComponent>
      </Grid2>
      <Grid2 size={{ xs: 12, md: 4 }}>
        <FormContainer component="loadGenerator">
          <FormLabel
            component="loadGenerator"
            style={{
              display: 'flex',
              alignItems: 'center',
              flexWrap: 'wrap',
            }}
          >
            Load generator
            <CustomTextTooltip title={infoloadGenerators} interactive>
              <HelpIcon />
            </CustomTextTooltip>
          </FormLabel>
          <RadioGroup
            aria-label="loadGenerator"
            name="loadGenerator"
            value={loadGenerator}
            onChange={handleChange('loadGenerator')}
            row
          >
            {loadGenerators.map((lg, index) => (
              <FormControlLabel
                key={index}
                value={lg}
                disabled={lg === 'wrk2'}
                control={<RadioButton color="primary" />}
                label={lg}
              />
            ))}
          </RadioGroup>
        </FormContainer>
      </Grid2>
    </Grid2>
  );
};

export default PerformanceForm;
