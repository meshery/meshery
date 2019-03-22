import IconButton from '@material-ui/core/IconButton';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import NoSsr from '@material-ui/core/NoSsr';
import { Badge, Drawer, Tooltip, Divider, Typography, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Icon } from '@material-ui/core';
import MesheryEventViewer from './MesheryEventViewer';
import NotificationsIcon from '@material-ui/icons/Notifications';
import { withStyles } from '@material-ui/core/styles';
import { eventTypes } from '../lib/event-types';
import classNames from 'classnames';
import amber from '@material-ui/core/colors/amber';


const styles = theme => ({
    sidelist: {
        width: 350,
    },
    notficiationDrawer: {
        backgroundColor: '#FFFFFF',
    },
    notificationTitle: {
        textAlign: 'center',
    },
    icon: {
        fontSize: 20,
    },
    iconVariant: {
        opacity: 0.9,
        marginRight: theme.spacing(1),
        marginTop: theme.spacing(1) * 3/4,
    },
    error: {
        backgroundColor: theme.palette.error.dark,
    },
    info: {
        backgroundColor: theme.palette.primary.dark,
    },
    warning: {
        backgroundColor: amber[700],
    },
    message: {
        display: 'flex',
        // alignItems: 'center',
    },
});

class MesheryNotification extends React.Component {

  state = {
    events: [{
        event_type: 0,
        summary: 'Logs from Istio vet',
        details: "E0322 17:14:30.743301 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:30.744773 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:31.043917 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:31.090949 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:31.301201 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:31.746774 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:31.749316 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:32.048922 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:32.093981 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:32.304137 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:32.751371 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:32.753722 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:33.053525 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:33.097251 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:33.306781 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:33.754353 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:33.757715 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:34.057151 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:34.100946 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:34.311764 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:34.756976 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:34.761242 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:35.063689 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:35.104594 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:35.315015 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:35.759662 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:35.764446 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:36.068087 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:36.110670 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:36.319696 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:36.762761 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:36.767865 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:37.071776 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:37.114354 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:37.324978 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:37.765897 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:37.772022 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:38.075389 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:38.117829 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:38.332164 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:38.769871 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:38.774717 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:39.080048 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:39.121658 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:39.341770 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:39.773311 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:39.777642 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:40.083805 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:40.127004 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:40.346784 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:40.777380 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:40.781142 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:41.087095 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:41.130036 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:41.351192 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:41.782450 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:41.785032 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:42.090001 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:42.134209 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:42.355410 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:42.785458 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:42.788062 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:43.093638 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:43.138497 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:43.358071 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:43.790472 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:43.791139 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:44.096118 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:44.142283 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:44.360539 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:44.794385 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:44.797170 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:45.098564 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:45.148091 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:45.364081 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:45.800851 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:45.800852 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:46.102626 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:46.152523 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:46.368102 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:46.805038 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:46.806826 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:47.108613 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:47.155841 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:47.371455 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:47.812234 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:47.812917 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:48.113219 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:48.159237 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:48.374430 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:48.819833 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:48.819869 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:49.116401 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:49.162755 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:49.377330 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:49.823225 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:49.824215 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:50.120060 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:50.165993 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:50.381746 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\n"
    },{
        event_type: 1,
        summary: 'Logs from Istio vet',
        details: "E0322 17:14:30.743301 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:30.744773 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:31.043917 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:31.090949 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:31.301201 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:31.746774 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:31.749316 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:32.048922 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:32.093981 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:32.304137 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:32.751371 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:32.753722 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:33.053525 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:33.097251 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:33.306781 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:33.754353 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:33.757715 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:34.057151 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:34.100946 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:34.311764 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:34.756976 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:34.761242 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:35.063689 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:35.104594 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:35.315015 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:35.759662 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:35.764446 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:36.068087 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:36.110670 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:36.319696 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:36.762761 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:36.767865 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:37.071776 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:37.114354 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:37.324978 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:37.765897 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:37.772022 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:38.075389 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:38.117829 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:38.332164 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:38.769871 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:38.774717 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:39.080048 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:39.121658 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:39.341770 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:39.773311 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:39.777642 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:40.083805 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:40.127004 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory."
    },{
        event_type: 2,
        summary: 'Logs from Istio vet',
        details: "E0322 17:14:30.743301 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized\nE0322 17:14:30.744773 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized\nE0322 17:14:31.043917 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized\nE0322 17:14:31.090949 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized\nE0322 17:14:31.301201 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized\nE0322 17:14:31.746774 1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86:"
    }],
    open: false,
    dialogShow: false,
    k8sConfig: {
        inClusterConfig: false,
        k8sfile: '', 
        contextName: '', 
        meshLocationURL: '', 
        reconfigureCluster: true,
    },
    createStream: false,
  }

  handleToggle = () => {
    this.setState(state => ({ open: !state.open }));
  };

  handleClose = event => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ open: false });
  };

  static getDerivedStateFromProps(props, state){
    if (JSON.stringify(props.k8sConfig) !== JSON.stringify(state.k8sConfig)) {
        return {
            createStream: true,
            k8sConfig: props.k8sConfig,
        };
    }
    return null;
  }

  componentDidUpdate(){
      const {createStream, k8sConfig} = this.state;
    if (k8sConfig.k8sfile === '' && k8sConfig.meshLocationURL === '') {
        this.closeEventStream();
    }
    if (createStream && k8sConfig.k8sfile !== '' && k8sConfig.meshLocationURL !== '') {
        this.startEventStream();
    }
  }

  async startEventStream() {
    this.closeEventStream();
    this.eventStream = new EventSource("/api/events");
    this.eventStream.onmessage = this.handleEvents;
    this.eventStream.onerror = this.handleError;
    this.setState({createStream: false});
  }

  handleEvents = e => {
    const {events} = this.state;
    const data = JSON.parse(e.data);
    events.push(data);
    this.setState({events});
  }

  handleError = e => {
    // attempting to reestablish connection
    this.startEventStream();
  }

  closeEventStream() {
      if(this.eventStream && this.eventStream.close){
        this.eventStream.close();
      }
  }

  deleteEvent = ind => () => {
      const {events} = this.state;
      if (events[ind]){
        events.splice(ind, 1);
      }
      this.setState({events, dialogShow: false});
  }

  clickEvent = (event,ind) => () => {
    const {events} = this.state;
    let fInd = -1;
    events.forEach((ev, i) => {
        if (ev.event_type === event.event_type && ev.summary === event.summary && ev.details === event.details) {
            fInd = i;
        }
    });
    if (fInd === ind) {
        this.setState({ open:true, dialogShow: true, ev: event, ind });
    }
  }

  handleDialogClose = event => {
    this.setState({ dialogShow: false });
  };

  viewEventDetails = () => {
      const {classes} = this.props;
    const {ev, ind, dialogShow} = this.state;
    if (ev && typeof ind !== 'undefined') {
        console.log(`decided icon class: ${JSON.stringify(eventTypes[ev.event_type]?eventTypes[ev.event_type].icon:eventTypes[0].icon)}`);
        const Icon = eventTypes[ev.event_type]?eventTypes[ev.event_type].icon:eventTypes[0].icon;
    return (
    <Dialog
            fullWidth
            maxWidth='md'
            open={dialogShow}
            onClose={this.handleDialogClose}
            aria-labelledby="event-dialog-title"
            >
            <DialogTitle id="event-dialog-title">
                <span id="client-snackbar" className={classes.message}>
                    <Icon className={classNames(classes.icon, classes.iconVariant)} fontSize='large' />
                    {ev.summary}
                </span>
            </DialogTitle>
            <Divider light variant="fullWidth" />
            <DialogContent>
                <DialogContentText>
                {ev.details.split('\n').map(det => {
                    return (
                        <div>{det}</div>
                    );
                })} 
                </DialogContentText>
            </DialogContent>
            <Divider light variant="fullWidth" />
            <DialogActions>
                <Button onClick={this.deleteEvent(ind)} color="secondary" variant="outlined">
                Dismiss
                </Button>
                <Button onClick={this.handleDialogClose} color="primary" variant="outlined">
                Close
                </Button>
            </DialogActions>
        </Dialog>
        );
    }
    return null;
}

  render() {
    const {classes} = this.props;
    const { open, events, ev, ind } = this.state;
    const self = this;

    let toolTipMsg = `There are ${events.length} events`;
    switch (events.length) {
        case 0:
            toolTipMsg = `There are no events`;
        case 1:
            toolTipMsg = `There is 1 event`;
    }

    return (
      <div>
        <NoSsr>
        <Tooltip title={toolTipMsg}>
        <IconButton 
            buttonRef={node => {
                this.anchorEl = node;
              }}
            color="inherit" onClick={this.handleToggle}>
            <Badge badgeContent={events.length} color="secondary">
                <NotificationsIcon />
            </Badge>
        </IconButton>
        </Tooltip>

        <Drawer anchor="right" open={open} onClose={this.handleClose} 
        classes={{
            paper: classes.notficiationDrawer,
        }}
        >
          <div
            tabIndex={0}
            role="button"
            onClick={this.handleClose}
            onKeyDown={this.handleClose}
          >
            <div className={classes.sidelist}>
            <div className={classes.notificationTitle}>
                <Typography variant="subtitle2">Notifications</Typography>
            </div>
            <Divider light />
            {events && events.map((event, ind) => (
                <MesheryEventViewer eventVariant={event.event_type} eventSummary={event.summary} 
                    deleteEvent={self.deleteEvent(ind)} 
                    onClick={self.clickEvent(event, ind)} />
            ))}
            </div>
          </div>
        </Drawer>
        {this.viewEventDetails(ev, ind)}
        </NoSsr>
    </div>
    )
  }
}


const mapStateToProps = state => {
    const k8sConfig = state.get("k8sConfig").toObject();
    return {k8sConfig};
  }

export default withStyles(styles)(connect(
    mapStateToProps
)(MesheryNotification));