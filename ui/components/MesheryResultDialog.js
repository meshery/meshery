import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, Grid, Table, TableRow, TableCell, TableBody, Typography,
} from '@material-ui/core';
import MesheryChartDialog from './MesheryChartDialog';


const defaultToolbarSelectStyles = {
  iconButton : {
    marginRight : '24px',
    top : '50%',
    display : 'inline-block',
    position : 'relative',
  },
  icon : { color : '#000', },
  inverseIcon : { transform : 'rotate(90deg)', },
  row : { borderBottom : 'none', },
};

class MesheryResultDialog extends React.Component {
    state = { dialogOpen : true, }

    handleDialogClose = () => {
      this.setState({ dialogOpen : false });
      this.props.close();
    }

    createTableRow(key, val) {
      const { classes } = this.props;
      return (
        <TableRow>
          <TableCell align="right" component="th" scope="row" className={classes.row}>
            {key}
            :
          </TableCell>
          <TableCell className={classes.row}>{val}</TableCell>
        </TableRow>
      );
    }

    renderKubernetesInfo(kubernetes) {
      const { classes } = this.props;
      return (
        <NoSsr>
          <Typography className={classes.title} variant="h6" id="tableTitle">
            Environment
          </Typography>
          <Table className={classes.table} size="small" aria-label="Environment">
            <TableBody>
              {this.createTableRow('Kubernetes API Server', kubernetes.server_version)}
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
                  {this.createTableRow('Internal IP Address', node.internal_ip)}
                  {this.createTableRow('Hostname', node.hostname)}
                  {this.createTableRow('Allocatable CPU', node.allocatable_cpu)}
                  {this.createTableRow('Allocatable Memory', node.allocatable_memory)}
                  {this.createTableRow('Capacity CPU', node.capacity_cpu)}
                  {this.createTableRow('Capacity Memory', node.capacity_memory)}
                  {this.createTableRow('Architecture', node.architecture)}
                  {this.createTableRow('Operating system', node.operating_system)}
                  {this.createTableRow('OS Image', node.os_image)}
                  {this.createTableRow('Container runtime version', node.container_runtime_version)}
                  {this.createTableRow('Kubelet version', node.kubelet_version)}
                  {this.createTableRow('Kubeproxy version', node.kubeproxy_version)}
                </NoSsr>
              ))}
            </TableBody>
          </Table>
        </NoSsr>
      );
    }

    renderMeshesInfo(detectedMeshes) {
      const { classes } = this.props;
      const meshes = Object.keys(detectedMeshes);
      return (
        <NoSsr>
          <Typography className={classes.title} variant="h6" id="tableTitle">
            Service Mesh
            {meshes.length > 1
              ? 'es'
              : ''}
          </Typography>
          <Table className={classes.table} size="small" aria-label="Service Mesh">
            <TableBody>
              {meshes.map(([mesh, version], ind) => (
                <NoSsr key={mesh.uniqueID}>
                  {meshes.length > 1
                    ? (
                      <TableRow>
                        <TableCell colSpan={2} className={classes.row} align="center">
                          <strong>
                            Service Mesh
                            {ind + 1}
                          </strong>
                        </TableCell>
                      </TableRow>
                    )
                    : ''}
                  {this.createTableRow('Name', mesh)}
                  {this.createTableRow('Version', version)}
                </NoSsr>
              ))
              }
            </TableBody>
          </Table>
        </NoSsr>
      );
    }

    renderLoadProfile(rowData) {
      const { classes } = this.props;
      let contents = '';
      if (rowData.runner_results) {
        let percentiles = '';
        if (rowData.runner_results.DurationHistogram && rowData.runner_results.DurationHistogram.Percentiles) {
          percentiles = rowData.runner_results.DurationHistogram.Percentiles.map((p) => this.createTableRow(`p${p.Percentile} Response Time`, p.Value));
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
          this.createTableRow('URL', rowData.runner_results.URL),
          this.createTableRow('Requested QPS', rowData.runner_results.RequestedQPS),
          this.createTableRow('Actual QPS', rowData.runner_results.ActualQPS),
          this.createTableRow('Threads', rowData.runner_results.NumThreads),
          this.createTableRow('Connections', rowData.runner_results.SocketCount),
          this.createTableRow('Requested Duration', rowData.runner_results.RequestedDuration),
          this.createTableRow('Actual Duration', (rowData.runner_results.ActualDuration / 1000000000).toFixed(1)),
          this.createTableRow('Average Response Time', rowData.runner_results.DurationHistogram.Avg),
          ...percentiles,
          this.createTableRow('Maximum Response Time', rowData.runner_results.DurationHistogram.Max),
        ];
      }
      return (
        <NoSsr>
          <Typography className={classes.title} variant="h6" id="tableTitle">
            Load Profile
          </Typography>
          <Table className={classes.table} size="small" aria-label="Load Profile">
            <TableBody>
              {/* {this.createTableRow('Name', rowData.name)} */}
              {contents}
            </TableBody>
          </Table>
        </NoSsr>
      );
    }

    render() {
      const { rowData } = this.props;

      return (
        <NoSsr>
          <MesheryChartDialog
            title={`Details${rowData
              ? ` - ${rowData.name}`
              : ''}`}
            handleClose={this.handleDialogClose}
            open={this.state.dialogOpen}
            content={(
              <div>
                <Grid container spacing={1}>
                  {rowData && rowData.runner_results && rowData.runner_results.kubernetes
                && (
                  <Grid item xs={12} sm={6}>
                    {this.renderKubernetesInfo(rowData.runner_results.kubernetes)}
                  </Grid>
                )}
                  {rowData
                && (
                  <Grid item xs={12} sm={6}>
                    {this.renderLoadProfile(rowData)}
                  </Grid>
                )}
                  {rowData && rowData.runner_results && rowData.runner_results['detected-meshes']
                && (
                  <Grid item xs={12} sm={6}>
                    {this.renderMeshesInfo(rowData.runner_results['detected-meshes'])}
                  </Grid>
                )}
                </Grid>
              </div>
            )}
          />
        </NoSsr>
      );
    }
}


MesheryResultDialog.propTypes = {
  // classes: PropTypes.object.isRequired,
  rowData : PropTypes.object.isRequired,
  close : PropTypes.func.isRequired, };

export default withStyles(defaultToolbarSelectStyles, { name : 'MesheryResultDialog' })(MesheryResultDialog);
