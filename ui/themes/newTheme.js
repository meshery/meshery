import { createTheme } from '@mui/material/styles';
// import { iconMedium } from '../css/icons.styles';
import { blueGrey } from '@mui/material/colors';

export let darkTheme = createTheme({
  typography: {
    fontFamily: ['Qanelas Soft', 'Roboto', 'Helvectica', 'Arial', 'sans-serif'].join(','),
    h5: {
      fontWeight: 'bolder',
      fontSize: 26,
      color: '#FFF',
      letterSpacing: 0.5,
    },
    p: {
      color: '#FFF',
    },
    h6: {
      color: '#FFF',
    },
  },
  palette: {
    mode: 'dark',
    primary: blueGrey,
    secondary: {
      main: '#EE5351',
      primeColor: '#303030',
      dark: '#121212',
      titleText: '#FBFBFB',
      text: '#FFF',
      text2: '#7494a1',
      text3: '#FFF',
      textMain: '#F6F8F8',
      titleBackground: '#000',
      mainBackground: '#202020',
      mainBackground2: '#303030',
      elevatedComponents: '#202020',
      elevatedComponents2: '#303030',
      elevatedComponents3: '#303030',
      lightText: 'rgba(255, 255, 255, 0.54)',
      icon: 'rgba(255, 255, 255, 0.54)',
      icon2: '#E6E6E6',
      iconMain: 'green',
      disabledIcon: 'rgba(255, 255, 255, 0.26)',
      chevron: 'rgb(255, 255, 255, 0.2)',
      link: 'rgba(255, 255, 255, 0.7)',
      link2: '#05FFCD',
      headerColor: '#202020',
      sideBar: '#1A1A1A',
      drawer: '#252E31',
      drawerHover: '#202020',
      img: 'invert(0.8)',
      appBar: '#363636',
      number: '#eee',
      completeInvert: 'invert(1)',
      canvas: '#1A1A1A',
      brightness: 'brightness(200%)',
      switcherButtons: '#1e1e1e',
      honeyComb: '#303030',
      filterChipBackground: '#222222',
      searchBackground: '#294957',
      searchBorder: '#396679',
      tabs: '#202020',
      modalTabs: '#363636',
      tabHover: '#212121',
      confirmationModal: '#111111',
      focused: '#00b39f',
      primaryModalText: '#FFF',
      default: '#9FAFB6',
      success: '#00D3A9',
      primary: '#86B2C6',
      warning: '#EBC017',
      error: '#F91313',
      lightError: '#B32700',
      penColorPrimary: '#E6E6E6',
      penColorSecondary: '#E6E6E6',
      toolbarBg2: '#464646',
      menuBg: '#363636',
      menuActionBg: '#222',
      menuSelectedBg: 'rgba(0, 179, 159, 0.25)',
      menuActionText: '#FBFBFB',
      menuItemBorder: '#979797',
      pinball: '#222222',
      innertableBg1: 'rgb(255, 255, 255, 0.1)',
      innertableBg2: 'rgb(255, 255, 255, 0.05)',
    },
  },
  shape: { borderRadius: 8 },
  breakpoints: {
    values: {
      xs: 0,
      sm: 600,
      md: 960,
      lg: 1280,
      xl: 1920,
    },
  },
});

darkTheme = {
  ...darkTheme,
  components: {
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundColor: '#303030',
        },
      },
      variants: [
        {
          props: { variant: 'elevation' },
          style: {
            boxShadow: '0px 0px 0px 0px rgba(0,0,0,0.0)',
          },
        },
      ],
    },
  },
};

let theme = createTheme({
  typography: {
    fontFamily: ['Qanelas Soft', 'Roboto', 'Helvectica', 'Arial', 'sans-serif'].join(','),
    h5: {
      fontWeight: 'bolder',
      fontSize: 26,
      color: '#FFF',
      letterSpacing: 0.5,
    },
    p: {
      color: '#FFF',
    },
    h6: {
      color: '#FFF',
    },
  },
  palette: {
    mode: 'light',
    primary: blueGrey,
    secondary: {
      main: '#EE5351',
      primeColor: '#ebeff1',
      dark: '#455a64',
      titleText: '#7494A1',
      text: '#000',
      text2: 'rgba(57, 102, 121, .9)',
      text3: '#333333',
      textMain: '#3C494F',
      titleBackground: 'rgba(57, 102, 121, .1)',
      mainBackground: '#396679',
      mainBackground2: '#FFF',
      elevatedComponents: '#FFF',
      elevatedComponents2: '#eaeff1',
      elevatedComponents3: '#FFF',
      lightText: 'rgba(0, 0, 0, 0.54)',
      icon: 'rgba(0, 0, 0, 0.54)',
      icon2: 'gray',
      iconMain: 'yellow',
      disabledIcon: 'rgba(0, 0, 0, 0.26)',
      chevron: '#FFF',
      link: '#000',
      link2: '#00b39F',
      headerColor: '#eeeeee',
      sideBar: '#FFF',
      drawer: '#FFF',
      drawerHover: '#f2f5f7',
      img: 'none',
      appBar: '#FFF',
      number: '#607d8b',
      completeInvert: 'none',
      canvas: '#fff',
      brightness: 'none',
      switcherButtons: '#335c6d',
      honeyComb: '#F0F0F0',
      filterChipBackground: '#CCCCCC',
      searchBackground: '#fff',
      searchBorder: '#CCCCCC',
      tabs: '#eeeeee87',
      modalTabs: '#dddddd',
      tabHover: '#e3e3e3',
      confirmationModal: 'rgb(234, 235, 236)',
      focused: '#607d8b',
      primaryModalText: '#FFF',
      default: '#51636B',
      success: '#00B39F',
      primary: '#477E96',
      warning: '#F0A303',
      error: '#8F1F00',
      lightError: '#8F1F00',
      penColorPrimary: '#3C494F',
      penColorSecondary: '#677E88',
      toolbarBg1: '#FFFFFF',
      menuBg: '#EBEFF1',
      menuItemBg: '#EBEFF1',
      menuActionBg: '#51636B',
      menuSelectedBg: '#CCC',
      menuActionText: '#FBFBFB',
      menuItemBorder: '#979797',
      pinball: '#d3d3d3',
      innertableBg1: '#ffffff',
      innertableBg2: '#ECECED',
    },
  },
  shape: { borderRadius: 8 },
  breakpoints: {
    values: {
      xs: 0,
      sm: 600,
      md: 960,
      lg: 1280,
      xl: 1920,
    },
  },
});

theme = {
  ...theme,
  components: {
    // MuiPaper: {
    //     styleOverrides: {
    //         root: {
    //             backgroundColor: "#fff",
    //         },
    //     },
    //     variants: [
    //         {
    //             props: { variant: 'elevation' },
    //             style: {
    //                 boxShadow: '0px 0px 0px 0px rgba(0,0,0,0.0)',
    //             },
    //         },
    //     ],
    // },
  },
};

export default theme;
