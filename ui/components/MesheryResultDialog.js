import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, Table, TableRow, TableCell, TableBody, Typography } from '@material-ui/core';
import MesheryChartDialog from './MesheryChartDialog';

const defaultToolbarSelectStyles = {
  iconButton: {
    marginRight: '24px',
    top: '50%',
    display: 'inline-block',
    position: 'relative',
  },
  icon: { color: '#000' },
  inverseIcon: { transform: 'rotate(90deg)' },
  row: { borderBottom: 'none' },
};

function MesheryResultDialog(props) {
  const { classes, rowData, close } = props;
  const [dialogOpen, setDialogOpen] = useState(true);

  const handleDialogClose = () => {
    setDialogOpen(false);
    close();
  };

  const createTableRow = (key, val) => {
    return (
      <TableRow>
        <TableCell align="right" component="th" scope="row" className={classes.row}>
          {key}:
        </TableCell>
        <TableCell className={classes.row}>{val}</TableCell>
      </TableRow>
    );
  };

  const renderKubernetesInfo = (kubernetes) => {
    return (
      <NoSsr>
        <Typography className={classes.title} variant="h6" id="tableTitle">
          Environment
        </Typography>
        <Table className={classes.table} size="small" aria-label="Environment">
          <TableBody>
            {createTableRow('Kubernetes API Server', kubernetes.server_version)}
            {kubernetes.nodes.map((node, ind) => (
              <NoSsr key={node.uniqueID}>
                <TableRow>
                  <TableCell colSpan={2} className={classes.row} align="center">
                    <strong>
                      Node
                      {ind + 1}
                    </strong>
                  </TableCell>
                </TableRow>
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
        <Typography className={classes.title} variant="h6" id="tableTitle">
          Service Mesh
          {meshes.length > 1 ? 'es' : ''}
        </Typography>
        <Table className={classes.table} size="small" aria-label="Service Mesh">
          <TableBody>
            {meshes.map((mesh, ind) => (
              <NoSsr key={mesh}>
                {meshes.length > 1 ? (
                  <TableRow>
                    <TableCell colSpan={2} className={classes.row} align="center">
                      <strong>
                        Service Mesh
                        {ind + 1}
                      </strong>
                    </TableCell>
                  </TableRow>
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

      // let reqDuration = rowData.runner_results.RequestedDuration.substring(0, rowData.runner_results.RequestedDuration.length-1);
      // switch(rowData.runner_results.RequestedDuration.slice(-1)){
      //   case 's':
      //     reqDuration += ' seconds';
      //     break;
      //   case 'm':
      //     reqDuration += ' minutes';
      //     break;
      //   case 'h':
      //     reqDuration += ' hours';
      //     break;
      // }

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
        <Typography className={classes.title} variant="h6" id="tableTitle">
          Load Profile
        </Typography>
        <Table className={classes.table} size="small" aria-label="Load Profile">
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
  classes: PropTypes.object.isRequired,
  rowData: PropTypes.object.isRequired,
  close: PropTypes.func.isRequired,
};

export default withStyles(defaultToolbarSelectStyles, { name: 'MesheryResultDialog' })(
  MesheryResultDialog,
);
