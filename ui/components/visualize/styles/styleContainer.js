export default class GraphStyle {
  static getLayout() {
    //Right now using cose. Adding more options in future
    return {
      name: 'cose',
      animate: false,
      fit: false,
      nodeDimensionsIncludeLabels: true
    };
  }

  static renderNode() {
    // Icon path is assumed to be of 32x32 in this example. You may auto calculate this if you wish.
    const iconPath = 'M31.4,16.5c0.8,0.8,0.8,2,0,2.8l-6,6c-0.8,0.8-2,0.8-2.8,0l-2.9-2.9c0.6,3,0.2,6.2-1.3,9 c-0.2,0.3-0.5,0.5-0.9,0.5c-0.3,0-0.5-0.1-0.7-0.3l-7.6-7.5l-1.4,1.4C7.9,25.7,8,25.8,8,26c0,1.1-0.9,2-2,2s-2-0.9-2-2 c0-1.1,0.9-2,2-2c0.2,0,0.3,0.1,0.5,0.1l1.4-1.4l-7.5-7.5c-0.5-0.5-0.4-1.3,0.2-1.6c2.7-1.5,5.9-1.9,8.9-1.3L6.6,9.4 c-0.8-0.8-0.8-2,0-2.8l6-6C13,0.2,13.5,0,14,0c0.5,0,1,0.2,1.4,0.6L19.9,5l4.4-4.4c0.8-0.8,2-0.8,2.8,0c0,0,0,0,0,0l4.3,4.3 c0.8,0.8,0.8,2,0,2.8L27,12.1L31.4,16.5z M14,18c-2.5-2.5-6-3.5-9.4-2.8l12.2,12.2C17.5,24.1,16.5,20.5,14,18z M9.4,8l3.7,3.7 l4.6-4.6L14,3.4L9.4,8z M14.6,14.6c0.6,0.4,1.1,0.8,1.6,1.3c0.5,0.5,0.9,1,1.3,1.6L28.6,6.3l-2.9-2.9L14.6,14.6z M28.6,18l-3.7-3.7 l-4.6,4.6l3.7,3.7L28.6,18z';
    const iconColor = '#ffffff';
    const size = 32; // may need to calculate this yourself
    const iconResize = 22; // adjust this for more "padding" (bigger number = more smaller icon)
  
    const width = size;
    const height = size;
    const scale = (size - iconResize) / size;
    const iconTranslate = iconResize / 2 / scale;
    const backgroundColor = `#33362F`;
  
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}">
        <rect x="0" y="0" width="${width}" height="${height}" fill="${backgroundColor}"></rect>
        <path d="${iconPath}" fill="${iconColor}" transform="scale(${scale}) translate(${iconTranslate}, ${iconTranslate}) "></path>
      </svg>`;
    return {
      svg: 'data:image/svg+xml;base64,' + btoa(svg),
      width,
      height,
    };
  }
  //Styling nodes and edges
  //Work needed to be done. Feel free to play around with it
  static getStylesheetContainer() {
    return [
      {
        selector: "node",
        style: {
          "background-color": "rgb(255,255,255)",
          "background-width": "80%",
          "background-height": "80%",
          "background-position-x": "1px",
          "background-position-y": "1px",
          "border-width": "1px",
          ghost: "yes",
          "ghost-offset-x": "1px",
          "ghost-offset-y": "1px",
          "ghost-opacity": "0.4",
          height: "25px",
          width: "25px",
          "z-index": "10"
        }
      },
      {
        selector: "node[app = 'istio']",
        style: {
          'background-image': (ele) => this.renderNode(ele).svg,
          'background-fit': 'contain',
          'background-opacity': 0, // do not show the bg colour
          'border-width': 0, // no border that would increase node size
          'background-clip': 'none', // let image go beyond node shape (also better performance)
          "z-index": "15"
        }
      },
      {
        selector: "node[app = 'linkerd']",
        style: {
          'background-image': (ele) => this.renderNode(ele).svg,
          'background-opacity': 0,
          'background-clip': 'none',
          "z-index": "15"
        }
      },
      {
        selector: "node[?isGroup]",
        style: {
          "background-color": "rgb(255,255,255)"
        }
      },
      {
        selector: "node:selected",
        style: {
          "border-width": "3px"
        }
      },
      {
        selector: "node.mousehighlight",
        style: {
          "font-size": "11px"
        }
      },
      {
        selector: "node.mousehighlight[^isGroup]",
        style: {
        }
      },
      {
        selector: "node.mousedim",
        style: {
          opacity: "0.6"
        }
      },
      {
        selector: "edge",
        style: {
          "curve-style": "bezier",
          "font-family": "Verdana,Arial,Helvetica,sans-serif,pficon",
          "font-size": "6px",
          label: "edge",
          "line-style": "solid",
          "target-arrow-shape": "vee",
          "text-events": "yes",
          "text-outline-color": "rgb(255,255,255)",
          "text-outline-width": "1px",
          "text-wrap": "wrap",
          width: "2px"
        }
      },
      {
        selector: "edge:selected",
        style: {
          width: "4px",
          label: "edge"
        }
      },
      {
        selector: "edge[protocol = \"tcp\"]",
        style: {
          "target-arrow-shape": "triangle-cross",
          "line-style": "solid"
        }
      },
      {
        selector: "edge.mousehighlight",
        style: {
          "font-size": "10px"
        }
      },
      {
        selector: "edge.mousehover",
        style: {
          label: "edge"
        }
      },
      {
        selector: "edge.mousedim",
        style: {
          opacity: "0.3"
        }
      },
      {
        selector: ".find[^isGroup]",
        style: {
          "overlay-color": "rgb(240,171,0)",
          "overlay-padding": "8px",
          "overlay-opacity": "0.5"
        }
      },
      {
        selector: 'node.highlight',
        style: {
          'border-color': '#123',
          'border-width': '2px'
        }
      },
      {
        selector: 'node.semitransp',
        style:{ 'opacity': '0.5' }
      },
      {
        selector: 'edge.highlight',
        style: { 'mid-target-arrow-color': '#FFF' }
      },
      {
        selector: 'edge.semitransp',
        style:{ 'opacity': '0.2' }
      }
    ]
  }
}