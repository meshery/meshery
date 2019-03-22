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
        details:'E0322 15:36:13.527222       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:13.531716       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:13.533791       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:13.536143       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:14.145852       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:14.530597       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:14.535708       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:14.539411       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:14.539435       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:15.150523       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:15.534957       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:15.538212       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:15.544614       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:15.544617       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:16.153258       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:16.538276       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:16.542704       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:16.548064       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:16.549022       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:17.159057       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:17.546703       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:17.546787       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:17.551159       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:17.552997       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:18.163434       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:18.551544       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:18.553520       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:18.556371       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:18.557319       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:19.167420       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:19.555362       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:19.558165       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:19.560440       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:19.560677       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:20.173810       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:20.558248       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:20.561875       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:20.565781       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:20.565813       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:21.180833       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:21.561676       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:21.564826       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:21.572087       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:21.574039       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:22.185122       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:22.566219       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:22.569715       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:22.578885       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:22.580484       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:23.190101       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:23.570063       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:23.572272       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:23.581990       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:23.584327       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:24.196438       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:24.574578       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:24.575706       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:24.584055       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:24.587750       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:25.200813       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:25.579326       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:25.579382       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:25.587105       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:25.590630       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:26.205093       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:26.582808       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:26.584030       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:26.590255       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:26.594133       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:27.209299       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:27.587160       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:27.587207       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:27.595738       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:27.598892       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:28.213123       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:28.591186       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:28.591573       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:28.601403       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:28.601868       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:29.215590       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:29.595580       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:29.596439       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:29.603760       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:29.613005       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:30.221467       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:30.598657       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:30.600329       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:30.606405       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:30.615848       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:31.228572       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:31.606617       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:31.606695       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:31.608962       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:31.618401       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:32.232185       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:32.613212       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:32.613290       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:32.614967       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:32.620820       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:33.235261       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵'
    },{
        event_type: 1,
        summary: 'Logs from Istio vet',
        details:'E0322 15:36:13.527222       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:13.531716       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:13.533791       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:13.536143       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:14.145852       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:14.530597       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:14.535708       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:14.539411       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:14.539435       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:15.150523       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:15.534957       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:15.538212       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:15.544614       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:15.544617       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:16.153258       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:16.538276       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:16.542704       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:16.548064       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:16.549022       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:17.159057       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:17.546703       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:17.546787       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:17.551159       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:17.552997       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:18.163434       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:18.551544       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:18.553520       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:18.556371       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:18.557319       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:19.167420       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:19.555362       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:19.558165       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:19.560440       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:19.560677       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:20.173810       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:20.558248       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:20.561875       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:20.565781       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:20.565813       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:21.180833       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:21.561676       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:21.564826       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:21.572087       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:21.574039       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:22.185122       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:22.566219       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:22.569715       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:22.578885       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:22.580484       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:23.190101       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:23.570063       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:23.572272       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:23.581990       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:23.584327       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:24.196438       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:24.574578       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:24.575706       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:24.584055       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:24.587750       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:25.200813       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:25.579326       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:25.579382       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:25.587105       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:25.590630       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:26.205093       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:26.582808       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:26.584030       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:26.590255       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:26.594133       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:27.209299       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:27.587160       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:27.587207       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:27.595738       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:27.598892       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:28.213123       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:28.591186       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:28.591573       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:28.601403       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:28.601868       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace: Unauthorized↵E0322 15:36:29.215590       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Pod: Unauthorized↵E0322 15:36:29.595580       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:29.596439       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: Unauthorized↵E0322 15:36:29.603760       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Service: Unauthorized↵E0322 15:36:29.613005       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Namespace:'
    },{
        event_type: 2,
        summary: 'Logs from Istio vet',
        details:'E0322 15:36:13.527222       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.Endpoints: Unauthorized↵E0322 15:36:13.531716       1 reflector.go:205] github.com/aspenmesh/istio-vet/vendor/k8s.io/client-go/informers/factory.go:86: Failed to list *v1.ConfigMap: '
    }],
    badgeContent: 1,
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
    this.setState({createStream: false});
  }

  handleEvents = e => {
    const {events} = this.state;
    const data = JSON.parse(e.data);
    events.push(data);
    this.setState({events, badgeContent: events.length});
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
      this.setState({events, dialogShow: false, badgeContent: events.length});
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
                {ev.details}
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
    const { open, badgeContent, events, ev, ind } = this.state;
    const self = this;

    return (
      <div>
        <NoSsr>
        <Tooltip title={`There are ${badgeContent} events`}>
        <IconButton 
            buttonRef={node => {
                this.anchorEl = node;
              }}
            color="inherit" onClick={this.handleToggle}>
            <Badge badgeContent={badgeContent} color="secondary">
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