import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
  styled,
  Grid2,
  Table,
  Typography,
  NoSsr,
  TableRow,
  TableCell,
  TableBody,
} from '@sistent/sistent';
import MesheryChartDialog from './general/MesheryChartDialog';

const StyledTableRow = styled(TableRow)(() => ({
  borderBottom: 'none',
}));

const StyledTableCell = styled(TableCell)(() => ({
  borderBottom: 'none',
}));

const StyledTitle = styled(Typography)(({ theme }) => ({
  marginBottom: theme.spacing(2),
}));

function MesheryResultDialog(props) {
  const { rowData, close } = props;
  const [dialogOpen, setDialogOpen] = useState(true);
  // The server sends camelCase; the snake_case fallback keeps older payloads working.
  const runnerResults = rowData?.runnerResults ?? rowData?.runner_results;

  const handleDialogClose = () => {
    setDialogOpen(false);
    close();
  };

  const createTableRow = (key, val) => {
    return (
      <StyledTableRow>
        <StyledTableCell align="right" component="th" scope="row">
          {key}:
        </StyledTableCell>
        <StyledTableCell>{val}</StyledTableCell>
      </StyledTableRow>
    );
  };

  const renderKubernetesInfo = (kubernetes) => {
    return (
      <NoSsr>
        <StyledTitle variant="h6" id="tableTitle">
          Environment
        </StyledTitle>
        <Table size="small" aria-label="Environment">
          <TableBody>
            {createTableRow('Kubernetes API Server', kubernetes.server_version)}
            {kubernetes.nodes.map((node, ind) => (
              <NoSsr key={node.uniqueID}>
                <StyledTableRow>
                  <StyledTableCell colSpan={2} align="center">
                    <strong>
                      Node
                      {ind + 1}
                    </strong>
                  </StyledTableCell>
                </StyledTableRow>
                {createTableRow('Internal IP Address', node.internalIp)}
                {createTableRow('Hostname', node.hostname)}
                {createTableRow('Allocatable CPU', node.allocatableCpu)}
                {createTableRow('Allocatable Memory', node.allocatableMemory)}
                {createTableRow('Capacity CPU', node.capacityCpu)}
                {createTableRow('Capacity Memory', node.capacityMemory)}
                {createTableRow('Architecture', node.architecture)}
                {createTableRow('Operating system', node.operatingSystem)}
                {createTableRow('OS Image', node.osImage)}
                {createTableRow('Container runtime version', node.containerRuntimeVersion)}
                {createTableRow('Kubelet version', node.kubeletVersion)}
                {createTableRow('Kubeproxy version', node.kubeproxyVersion)}
              </NoSsr>
            ))}
          </TableBody>
        </Table>
      </NoSsr>
    );
  };

  const renderMeshesInfo = (detectedMeshes) => {
    const meshes = Object.keys(detectedMeshes);
    return (
      <NoSsr>
        <StyledTitle variant="h6" id="tableTitle">
          Service Mesh
          {meshes.length > 1 ? 'es' : ''}
        </StyledTitle>
        <Table size="small" aria-label="Service Mesh">
          <TableBody>
            {meshes.map((mesh, ind) => (
              <NoSsr key={mesh}>
                {meshes.length > 1 ? (
                  <StyledTableRow>
                    <StyledTableCell colSpan={2} align="center">
                      <strong>
                        Service Mesh
                        {ind + 1}
                      </strong>
                    </StyledTableCell>
                  </StyledTableRow>
                ) : (
                  ''
                )}
                {createTableRow('Name', mesh)}
                {createTableRow('Version', detectedMeshes[mesh])}
              </NoSsr>
            ))}
          </TableBody>
        </Table>
      </NoSsr>
    );
  };

  const renderLoadProfile = (rowData) => {
    let contents = '';
    if (runnerResults) {
      let percentiles = '';
      if (runnerResults.DurationHistogram && runnerResults.DurationHistogram.Percentiles) {
        percentiles = runnerResults.DurationHistogram.Percentiles.map((p) =>
          createTableRow(`p${p.Percentile} Response Time`, p.Value),
        );
      }

      contents = [
        createTableRow('URL', runnerResults.URL),
        createTableRow('Requested QPS', runnerResults.RequestedQPS),
        createTableRow('Actual QPS', runnerResults.ActualQPS),
        createTableRow('Threads', runnerResults.NumThreads),
        createTableRow('Connections', runnerResults.SocketCount),
        createTableRow('Requested Duration', runnerResults.RequestedDuration),
        createTableRow('Actual Duration', (runnerResults.ActualDuration / 1000000000).toFixed(1)),
        ...(runnerResults.DurationHistogram
          ? [createTableRow('Average Response Time', runnerResults.DurationHistogram.Avg)]
          : []),
        ...percentiles,
        ...(runnerResults.DurationHistogram
          ? [createTableRow('Maximum Response Time', runnerResults.DurationHistogram.Max)]
          : []),
      ];
    }
    return (
      <NoSsr>
        <StyledTitle variant="h6" id="tableTitle">
          Load Profile
        </StyledTitle>
        <Table size="small" aria-label="Load Profile">
          <TableBody>{contents}</TableBody>
        </Table>
      </NoSsr>
    );
  };

  return (
    <NoSsr>
      <MesheryChartDialog
        title={`Details${rowData ? ` - ${rowData.name}` : ''}`}
        handleClose={handleDialogClose}
        open={dialogOpen}
        content={
          <div>
            <Grid2 container spacing={1} size="grow">
              {rowData && runnerResults && runnerResults.kubernetes && (
                <Grid2 size={{ xs: 12, sm: 6 }}>
                  {renderKubernetesInfo(runnerResults.kubernetes)}
                </Grid2>
              )}
              {rowData && <Grid2 size={{ xs: 12, sm: 6 }}>{renderLoadProfile(rowData)}</Grid2>}
              {rowData && runnerResults && runnerResults['detected-meshes'] && (
                <Grid2 size={{ xs: 12, sm: 6 }}>
                  {renderMeshesInfo(runnerResults['detected-meshes'])}
                </Grid2>
              )}
            </Grid2>
          </div>
        }
      />
    </NoSsr>
  );
}

MesheryResultDialog.propTypes = {
  rowData: PropTypes.object.isRequired,
  close: PropTypes.func.isRequired,
};

export default MesheryResultDialog;
