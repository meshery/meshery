import React from 'react';
import { NoSsr } from '@mui/material';
import Popup from '../Popup';
import { withNotify } from '../../utils/hooks/useNotification';
import DashboardMeshModelGraph from './charts/DashboardMeshModelGraph.js';
import ConnectionStatsChart from './charts/ConnectionCharts.js';
import MesheryConfigurationChart from './charts/MesheryConfigurationCharts.js';
import { Provider } from 'react-redux';
import { store } from '@/store/index';
import { useTheme, styled, blue, Grid, charcoal } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';

const Root = styled('div')(() => ({
  marginTop: '1.5rem',
}));

const useStyles = (theme) => ({
  datatable: {
    boxShadow: 'none',
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  link: {
    cursor: 'pointer',
    textDecoration: 'none',
  },
  metricsButton: { width: '240px' },
  alreadyConfigured: { textAlign: 'center' },
  margin: { margin: theme.spacing(1) },
  colorSwitchBase: {
    color: blue[300],
    '&.Mui-checked': {
      color: blue[500],
      '& + .MuiSwitch-track': { backgroundColor: blue[500] },
    },
  },
  fileLabel: { width: '100%' },
  fileLabelText: {},
  inClusterLabel: { paddingRight: theme.spacing(2) },
  alignCenter: { textAlign: 'center' },
  icon: { width: theme.spacing(2.5) },
  istioIcon: { width: theme.spacing(1.5) },
  settingsIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
    padding: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  redirectButton: {
    marginLeft: '-.5em',
    color: charcoal[10],
  },
  dashboardSection: {
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    marginBottom: theme.spacing(2),
  },
});

const Section = styled('div')(() => {
  const theme = useTheme();
  return {
    padding: theme.spacing(2),
    backgroundColor: theme.palette.background.default,
    borderRadius: '8px',
    height: '100%',
  };
});

const Overview = () => {
  const theme = useTheme();
  const classes = useStyles(theme);

  return (
    <NoSsr>
      <Popup />
      <Provider store={store}>
        <Root sx={{ paddingLeft: '0px' }}>
          <UsesSistent>
            <Grid container spacing={2}>
              <Grid item xs={12} md={12}>
                <DashboardMeshModelGraph classes={classes} />
              </Grid>

              <Grid item xs={12} md={6}>
                <Section>
                  <ConnectionStatsChart classes={classes} />
                </Section>
              </Grid>

              <Grid item xs={12} md={6} sx={{ paddingLeft: '0px' }}>
                <Section>
                  <MesheryConfigurationChart classes={classes} />
                </Section>
              </Grid>
            </Grid>
          </UsesSistent>
        </Root>
      </Provider>
    </NoSsr>
  );
};

export default withNotify(Overview);

// import React from 'react';
// import { NoSsr } from '@mui/material';
// import Popup from '../Popup';
// import { withNotify } from '../../utils/hooks/useNotification';
// import DashboardMeshModelGraph from './charts/DashboardMeshModelGraph.js';
// import ConnectionStatsChart from './charts/ConnectionCharts.js';
// import MesheryConfigurationChart from './charts/MesheryConfigurationCharts.js';
// import { Provider } from 'react-redux';
// import { store } from '@/store/index';
// import { useTheme, styled, blue, Grid, Paper, charcoal } from '@layer5/sistent';
// import { UsesSistent } from '../SistentWrapper';

// const Root = styled('div')(() => ({
//   marginTop: '1.5rem',
// }));

// const useStyles = (theme) => ({
//   datatable: {
//     boxShadow: 'none',
//   },
//   chip: {
//     marginRight: theme.spacing(1),
//     marginBottom: theme.spacing(1),
//   },
//   buttons: {
//     display: 'flex',
//     justifyContent: 'flex-end',
//   },
//   button: {
//     marginTop: theme.spacing(3),
//     marginLeft: theme.spacing(1),
//   },
//   link: {
//     cursor: 'pointer',
//     textDecoration: 'none',
//   },
//   metricsButton: { width: '240px' },
//   alreadyConfigured: { textAlign: 'center' },
//   margin: { margin: theme.spacing(1) },
//   colorSwitchBase: {
//     color: blue[300],
//     '&.Mui-checked': {
//       color: blue[500],
//       '& + .MuiSwitch-track': { backgroundColor: blue[500] },
//     },
//   },
//   fileLabel: { width: '100%' },
//   fileLabelText: {},
//   inClusterLabel: { paddingRight: theme.spacing(2) },
//   alignCenter: { textAlign: 'center' },
//   icon: { width: theme.spacing(2.5) },
//   istioIcon: { width: theme.spacing(1.5) },
//   settingsIcon: {
//     width: theme.spacing(2.5),
//     paddingRight: theme.spacing(0.5),
//   },
//   addIcon: {
//     width: theme.spacing(2.5),
//     paddingRight: theme.spacing(0.5),
//   },
//   cardHeader: { fontSize: theme.spacing(2) },
//   card: {
//     height: '100%',
//     marginTop: theme.spacing(2),
//     padding: theme.spacing(2),
//   },
//   cardContent: { height: '100%' },
//   redirectButton: {
//     marginLeft: '-.5em',
//     color: charcoal[10],
//   },
//   dashboardSection: {
//     padding: theme.spacing(2),
//     borderRadius: 4,
//     height: '100%',
//     marginBottom: theme.spacing(2),
//   },
// });

// const Overview = () => {
//   const theme = useTheme();
//   const classes = useStyles(theme);

//   return (
//     <NoSsr>
//       <Popup />
//       <Provider store={store}>
//         <Root sx={{ paddingLeft: '0px' }}>
//           <UsesSistent>
//             <Grid container spacing={2}>
//               <Grid item xs={12} md={12}>
//                 <DashboardMeshModelGraph classes={classes} />
//               </Grid>

//               <Grid item xs={12}>
//                 <Grid container spacing={2}>
//                   <Grid
//                     item
//                     xs={12}
//                     md={6}
//                     style={{
//                       padding: theme.spacing(2),
//                       backgroundColor: charcoal[10],
//                       borderRadius: '8px',
//                       height: '100%',
//                     }}
//                   >
//                     <ConnectionStatsChart classes={classes} />
//                   </Grid>
//                   <Grid
//                     item
//                     xs={12}
//                     md={6}
//                     style={{
//                       padding: theme.spacing(2),
//                       backgroundColor: charcoal[10],
//                       borderRadius: '8px',
//                       height: '100%',
//                     }}
//                   >
//                     <MesheryConfigurationChart classes={classes} />
//                   </Grid>
//                 </Grid>
//               </Grid>
//             </Grid>
//           </UsesSistent>
//         </Root>
//       </Provider>
//     </NoSsr>
//   );
// };

// export default withNotify(Overview);
