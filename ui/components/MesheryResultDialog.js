import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { styled, Grid, Table, Typography } from '@layer5/sistent';
import { NoSsr, TableRow, TableCell, TableBody } from '@mui/material';
import MesheryChartDialog from './MesheryChartDialog';

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
                {createTableRow('Internal IP Address', node.internal_ip)}
                {createTableRow('Hostname', node.hostname)}
                {createTableRow('Allocatable CPU', node.allocatable_cpu)}
                {createTableRow('Allocatable Memory', node.allocatable_memory)}
                {createTableRow('Capacity CPU', node.capacity_cpu)}
                {createTableRow('Capacity Memory', node.capacity_memory)}
                {createTableRow('Architecture', node.architecture)}
                {createTableRow('Operating system', node.operating_system)}
                {createTableRow('OS Image', node.os_image)}
                {createTableRow('Container runtime version', node.container_runtime_version)}
                {createTableRow('Kubelet version', node.kubelet_version)}
                {createTableRow('Kubeproxy version', node.kubeproxy_version)}
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
    if (rowData.runner_results) {
      let percentiles = '';
      if (
        rowData.runner_results.DurationHistogram &&
        rowData.runner_results.DurationHistogram.Percentiles
      ) {
        percentiles = rowData.runner_results.DurationHistogram.Percentiles.map((p) =>
          createTableRow(`p${p.Percentile} Response Time`, p.Value),
        );
      }

      contents = [
        createTableRow('URL', rowData.runner_results.URL),
        createTableRow('Requested QPS', rowData.runner_results.RequestedQPS),
        createTableRow('Actual QPS', rowData.runner_results.ActualQPS),
        createTableRow('Threads', rowData.runner_results.NumThreads),
        createTableRow('Connections', rowData.runner_results.SocketCount),
        createTableRow('Requested Duration', rowData.runner_results.RequestedDuration),
        createTableRow(
          'Actual Duration',
          (rowData.runner_results.ActualDuration / 1000000000).toFixed(1),
        ),
        createTableRow('Average Response Time', rowData.runner_results.DurationHistogram.Avg),
        ...percentiles,
        createTableRow('Maximum Response Time', rowData.runner_results.DurationHistogram.Max),
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
            <Grid container spacing={1}>
              {rowData && rowData.runner_results && rowData.runner_results.kubernetes && (
                <Grid item xs={12} sm={6}>
                  {renderKubernetesInfo(rowData.runner_results.kubernetes)}
                </Grid>
              )}
              {rowData && (
                <Grid item xs={12} sm={6}>
                  {renderLoadProfile(rowData)}
                </Grid>
              )}
              {rowData && rowData.runner_results && rowData.runner_results['detected-meshes'] && (
                <Grid item xs={12} sm={6}>
                  {renderMeshesInfo(rowData.runner_results['detected-meshes'])}
                </Grid>
              )}
            </Grid>
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
