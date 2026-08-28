import cytoscape from 'https://cdn.jsdelivr.net/npm/cytoscape@3.27.0/+esm'

const data = {
  "elements": {
    "nodes": [
      {
        "data": {
          "id": "non-semantic-parent",
          "label": "Rectangle (Annotation)",
          "schemaId": "components.meshery.io/v1beta1"
        },
        "position": {
          "x": 600,
          "y": 400
        },
        "group": "nodes",
        "removed": false,
        "selected": false,
        "selectable": true,
        "locked": false,
        "grabbable": true,
        "pannable": false,
        "classes": ""
      },
      {
        "data": {
          "id": "pod-child",
          "label": "Pod",
          "schemaId": "components.meshery.io/v1beta1",
          "parent": "non-semantic-parent"
        },
        "position": {
          "x": 550,
          "y": 400
        },
        "group": "nodes",
        "removed": false,
        "selected": false,
        "selectable": true,
        "locked": false,
        "grabbable": true,
        "pannable": false,
        "classes": ""
      },
      {
        "data": {
          "id": "service-child",
          "label": "Service",
          "schemaId": "components.meshery.io/v1beta1",
          "parent": "non-semantic-parent"
        },
        "position": {
          "x": 650,
          "y": 400
        },
        "group": "nodes",
        "removed": false,
        "selected": false,
        "selectable": true,
        "locked": false,
        "grabbable": true,
        "pannable": false,
        "classes": ""
      }
    ],
    "edges": []
  },
  "elementStyles": {
    "non-semantic-parent": {
      "events": "yes",
      "text-events": "no",
      "transition-property": "none",
      "display": "element",
      "visibility": "visible",
      "opacity": "1",
      "text-opacity": "1",
      "min-zoomed-font-size": "12px",
      "z-compound-depth": "auto",
      "z-index-compare": "manual",
      "z-index": "1",
      "overlay-padding": "10px",
      "overlay-opacity": "0",
      "underlay-padding": "10px",
      "underlay-opacity": "0",
      "text-valign": "top",
      "text-halign": "center",
      "color": "rgb(100, 100, 100)",
      "text-background-color": "rgb(210, 212, 210)",
      "text-background-opacity": "0.7",
      "font-family": "Qanelas Soft, sans-serif",
      "font-weight": "400",
      "font-size": "9px",
      "text-wrap": "wrap",
      "text-max-width": "120px",
      "line-height": "1.2",
      "label": "Rectangle (Annotation)",
      "text-margin-y": "-10px",
      "height": "120px",
      "width": "220px",
      "shape": "rectangle",
      "background-color": "rgb(240, 240, 240)",
      "background-opacity": "0.3",
      "border-color": "rgb(150, 150, 150)",
      "border-opacity": "1",
      "border-width": "2px",
      "border-style": "dashed",
      "padding": "35px"
    },
    "pod-child": {
      "events": "yes",
      "text-events": "no",
      "transition-property": "none",
      "display": "element",
      "visibility": "visible",
      "opacity": "1",
      "text-opacity": "1",
      "min-zoomed-font-size": "12px",
      "z-compound-depth": "auto",
      "z-index-compare": "manual",
      "z-index": "2",
      "overlay-padding": "10px",
      "overlay-opacity": "0",
      "underlay-padding": "10px",
      "underlay-opacity": "0",
      "text-valign": "bottom",
      "text-halign": "center",
      "color": "rgb(0, 0, 0)",
      "text-background-color": "rgb(210, 212, 210)",
      "text-background-opacity": "0.7",
      "font-family": "Qanelas Soft, sans-serif",
      "font-weight": "300",
      "font-size": "8px",
      "label": "Pod",
      "text-margin-y": "7px",
      "height": "24px",
      "width": "24px",
      "shape": "round-rectangle",
      "background-color": "rgb(50, 108, 229)",
      "background-opacity": "1",
      "border-width": "0px",
      "background-image": "url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iOTAiIGhlaWdodD0iNjYiIHZpZXdCb3g9IjAgMCA5MCA2NiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQ3IDM0LjFDNDYuNzM0OCAzNC4xIDQ2LjQ4MDQgMzQuMjA1NCA0Ni4yOTI5IDM0LjM5MjlDNDYuMTA1NCAzNC41ODA0IDQ2IDM0LjgzNDggNDYgMzUuMVY2NUM0NiA2NS4yNjUyIDQ2LjEwNTQgNjUuNTE5NiA0Ni4yOTI5IDY1LjcwNzFDNDYuNDgwNCA2NS44OTQ2IDQ2LjczNDggNjYgNDcgNjZIODlDODkuMjY1MiA2NiA4OS41MTk2IDY1Ljg5NDYgODkuNzA3MSA2NS43MDcxQzg5Ljg5NDYgNjUuNTE5NiA5MCA2NS4yNjUyIDkwIDY1VjM1LjFDOTAgMzQuODM0OCA4OS44OTQ2IDM0LjU4MDQgODkuNzA3MSAzNC4zOTI5Qzg5LjUxOTYgMzQuMjA1NCA4OS4yNjUyIDM0LjEgODkgMzQuMUg0N1pNNTUuMTggNTcuMjdDNTUuMTggNTcuNTM1MiA1NS4wNzQ2IDU3Ljc4OTYgNTguODg3MSA1Ny45NzcxQzU0LjY5OTYgNTguMTY0NiA1NC40NDUyIDU4LjI3IDU0LjE4IDU4LjI3QzUzLjkxNDggNTguMjcgNTMuNjYwNCA1OC4xNjQ2IDUzLjQ3MjkgNTcuOTc3MUM1My4yODU0IDU3Ljc4OTYgNTMuMTggNTcuNTM1MiA1My4xOCA1Ny4yN1Y0Mi43M0M1My4xOCA0Mi40NjQ4IDUzLjI4NTQgNDIuMjEwNCA1My40NzI5IDQyLjAyMjlDNTMuNjYwNCA0MS44MzU0IDUzLjkxNDggNDEuNzMgNTQuMTggNDEuNzNDNTQuNDQ1MiA0MS43MyA1NC42OTk2IDQxLjgzNTQgNTQuODg3MSA0Mi4wMjI5Qz11LjA3NDYgNDIuMjEwNCA1NS4xOCA0Mi40NjQ4IDU1LjE4IDQyLjczVjU3LjI3Wk02NC4zNiA1Ny4yN0M2NC4zNiA1Ny41MzUyIDY0LjI1NDYgNTcuNzg5NiA2NC4wNjcxIDU3Ljk3NzFDNjMuODc5NiA1OC4xNjQ2IDYzLjYyNTIgNTguMjdDNjMuMDk0OCA1OC4yNyA2Mi44NDA0IDU4LjE2NDYgNjIuNjUyOSA1Ny45NzcxQzYyLjQ2NTQgNTcuNzg5NiA2Mi4zNiA1Ny41MzUyIDYyLjM2IDU3LjI3VjQyLjczQzYyLjM2IDQyLjQ2NDggNjIuNDY1NCA0Mi4yMTA0IDYyLjY1MjkgNDIuMDIyOUM2Mi44NDA0IDQxLjgzNTQgNjMuMDk0OCA0MS43MyA2My4zNiA0MS43M0M2My42MjUyIDQxLjczIDYzLjg3OTYgNDEuODM1NCA2NC4wNjcxIDQyLjAyMjlDNjQuMjU0NiA0Mi4yMTA0IDY0LjM2IDQyLjQ2NDggNjQuMzYgNDIuNzNWNTcuMjdaTTczLjU0IDU3LjI3QzczLjU0IDU3LjUzNTIgNzMuNDM0NiA1Ny43ODk2IDczLjI0NzEgNTcuOTc3MUM3My4wNTk2IDU4LjE2NDYgNzIuODA1MiA1OC4yNyA3Mi41NCA1OC4yN0M3Mi4yNzQ4IDU4LjI3IDcyLjAyMDQgNTguMTY0NiA3MS44MzI5IDU3Ljk3NzFDNzEuNjQ1NCA1Ny45NzcxIDcxLjU0IDU3LjUzNTIgNzEuNTQgNTcuMjdWNDIuNzNDNzEuNTQgNDIuNDY0OCA3MS42NDU0IDQyLjIxMDQgNzEuODMyOSA0Mi4wMjI5QzcyLjAyMDQgNDEuODM1NCA3Mi4yNzQ4IDQxLjczIDcyLjU0IDQxLjczQzcyLjgwNTIgNDEuNzMgNzMuMDU5NiA0MS44MzU0IDczLjI0NzEgNDIuMDIyOUM3My40MzQ2IDQyLjIxMDQgNzMuNTQgNDIuNDY0OCA3My41NCAyMi43M1Y1Ny4yN1pNODIuNzIgNTcuMjdDODIuNzIgNTcuNTM1MiA4Mi42MTQ2IDU3Ljc4OTYgODIuNDI3MSA1Ny45NzcxQzgyLjIzOTYgNTguMTY0NiA4MS45ODUyIDU4LjI3IDgxLjcyIDU4LjI3QzgxLjQ1NDggNTguMjcgODEuMjAwNCA1OC4xNjQ2IDgxLjAxMjkgNTcuOTc3MUM4MC44MjU0IDU3Ljc4OTYgODAuNzIgNTcuNTM1MiA4MC43MiA1Ny4yN1Y0Mi43M0M4MC43MiA0Mi40NjQ4IDgwLjgyNTQgNDIuMjEwNCA4MS4wMTI5IDQyLjAyMjlDODEuMjAwNCA0MS44MzU0IDgxLjQ1NDggNDEuNzMgODEuNzIgNDEuNzNDODEuOTg1MiA0MS43MyA4Mi4yMzk2IDQxLjgzNTQgODIuNDI3MSA0Mi4wMjI5QzgyLjYxNDYgNDIuMjEwNCA4Mi43MiA0Mi40NjQ4IDgyLjcyIDQyLjczVjU3LjI3Wk0yNCAwQzIzLjczNDggMCAyMy40ODA0IDAuMTA1MzU3IDIzLjI5MjkgMC4yOTI4OTNDMjMuMTA1NCAwLjQ4MDQzIDIzIDAuNzM0Nzg0IDIzIDFWMzAuOUMyMyAzMS4xNjUyIDIzLjEwNTQgMzEuNDE5NiAyMy4yOTI5IDMxLjYwNzFDMjMuNDgwNCAzMS43OTQ2IDIzLjczNDggMzEuOSAyNCAzMS45SDY2QzY2LjI2NTIgMzEuOSA2Ni41MTk2IDMxLjc5NDYgNjYuNzA3MSAzMS42MDcxQzY2Ljg5NDYgMzEuNDE5NiA2NyAzMS4xNjUyIDY3IDMwLjlWMUM2NyAwLjczNDc4NCA2Ni44OTQ2IDAuNDgwNDMgNjYuNzA3MSAwLjI5Mjg5M0M2Ni41MTk2IDAuMTA1MzU3IDY2LjI2NTIgMCA2NiAwTDI0IDBaTTMyLjE4IDIzLjI3QzMyLjE4IDIzLjUzNTIgMzMyLjA3NDYgMjMuNzg5NiAzMS44ODcyIDIzLjk3NzFDMzEuNjk5NiAyNC4xNjQ2IDMxLjQ0NTIgMjQuMjcgMzEuMTggMjQuMjdDMzAuOTE0OCAyNC4yNyAzMC42NjA0IDI0LjE2NDYgMzAuNDcyOSAyMy45NzcxQzMwLjI4NTQgMjMuNzg5NiAzMC4xOCAyMy41MzUyIDMwLjE4IDIzLjI3VjguNzNDMzAuMTggOC40NjQ3OCAzMC4yODU0IDguMjEwNDMgMzAuNDcyOSA4LjAyMjg5QzMyLjY2MDQgNy44MzUzNiAzMC45MTQ4IDcuNzMgMzEuMTggNy43M0MzMS40NDUyIDcuNzMgMzEuNjk5NiA3LjgzNTM2IDMxLjg4NzEgOC4wMjI4OUMzMi4wNzQ2IDguMjEwNDMgMzIuMTggOC40NjQ3OCAzMi4xOCA4LjczVjIzLjI3Wk00MS40NSAyMy4yN0M0MS40NSAyMy41MzUyIDQxLjM0NDYgMjMuNzg5NiA0MS4xNTcxIDIzLjk3NzFDNDAuOTY5NiAyNC4xNjQ2IDQwLjcxNTIgMjQuMjcgNDAuNDUgMjQuMjdDNDAuMTg0OCAyMy44OTQ2IDM5Ljc0MjkgMjMuOTc3QzM5LjU1NTQgMjMuNzg5NiAzOS40NSAyMy41MzUyIDM5LjQ1IDIzLjI3VjguNzNDMzkuNDUgOC40NjQ3OCAzOS41NTU0IDguMjEwNDMgMzkuNzQyOSA4LjAyMjg5QzM5LjkzMDQgNy44MzUzNiA0MC4xODQ4IDcuNzMgNDAuNDUgNy43M0M0MC43MTUyIDcuNzMgNDAuOTY5NiA3LjgzNTM2IDQxLjE1NzEgOC4wMjI4OUM0MS4zNDQ2IDguMjEwNDMgNDEuNDUgOC40NjQ3OCA0MS40NSA4LjczVjIzLjI3Wk01MC42MyAyMy4yN0M1MC42MyAyMy41MzUyIDUwLjUyNDYgMjMuNzg5NiA1MC4zMzcxIDIzLjk3NzFDNTAuMTQ5NiAyNC4xNjQ2IDQ5Ljg5NTIgMi0uMjcgNDkuNjMgMjQuMjdDNDkuMzY0OCAyNC4yNyA0OS4xMTA0IDI0LjE2NDYgNDguOTIyOSAyMy45NzcxQzQ4LjczNTQgMjMuNzg5NiA0OC42MyAyMy41MzUyIDguNjMgMjMuMjdWOC43M0M0OC42MyA4LjQ2NDc4IDg5LjczNTQgOC4yMTA0MyA0OC45MjI5IDguMDIyODlDNDkuMTEwNCA3LjgzNTM2IDQ5LjM2NDggNy43MyA5OS42MyA3LjczQzQ5Ljg5NTIgNy43MyA1MC4xNDk2IDcuODM1MzYgNTAuMzM3MSA4LjAyMjg5QzUwLjUyNDYgOC4yMTA0MyA1MC42MyA4LjQ2NDc4IDUwLjYzIDguNzNWMjMuMjdaTTU5LjgxIDIzLjI3QzU5LjgxIDIzLjUzNTIgNTkuNzA0NiAyMy43ODk2IDU5LjUxNzEgMjMuOTc3FDNTkuMzI5NiAyNC4xNjQ2IDU5LjA3NTIgMjQuMjcgNTguODEgMjQuMjdDNTguNTQ0OCAyNC4yNyA1OC4yOTA0IDI0LjE2NDYgNTguMTAyOSAyMy55NzcxQzU3LjkxNTQgMjMuNzg5NiA1Ny44MSAyMy41MzUyIDU3LjgxIDIzLjI3VjguNzNDNTcuODEgOC40NjQ3OCA1Ny45MTU0IDguMjEwNDMgNTguMTAyOSA4LjAyMjg5QzU4LjI5MDQgNy44MzUzNiA1OC41NDQ4IDcuNzMgNTguODEgNy43M0M1OS4wNzUyIDcuNzMgNTkuMzI5NiA3LjgzNTM2IDU5LjUxNzEgOC4wMjI4OUM1OS43MDQ2IDguMjEwNDMgNTkuODEgOC40NjQ3OCA1OS44MSA4LjczVjIzLjI3Wk0xIDM0LjFDMC43MzQ3ODQgMzQuMSAwLjQ4MDQzIDM0LjIwNTQgMC4yOTI4OTMgMzQuMzkyOUMwLjEwNTM1NyAzNC41ODA0IDAgMzQuODM0OCAwIDM1LjFMMCA2NUMwIDY1LjI2NTIgMC4xMDUzNTcgNjUuNTE5NiAwLjI5Mjg5MyA2NS43MDcxQzAuNDgwNDMgNjUuODk0NiAwLjczNDc4NiA2NiAxIDY2SDQzQzQzLjI2NTIgNjYgNDMuNTE5NiA2NS44OTQ2IDQzLjcwNzEgNjUuNzA3MUM0My44OTQ2IDY1LjUxOTYgNDQgNjUuMjY1MiA0NCA2NVYzNS4xQzQ0IDM0LjgzNDggNDMuODk0NiAzNC41ODA0IDQzLjcwNzEgMzQuMzkyOUM0My41MTk2IDM0LjIwNTQgNDMuMjY1MiAzNC4xIDQzIDM0LjFIMVpNOS4xOCA1Ny4yN0M5LjE5MzYzIDU3LjU0OTYgOS4wOTYxNCA1Ny44MjMyIDguOTA4OCA1OC4wMzExQzguNzI2MTQgNTguMjM5MSA4LjQ1OTQ3IDU4LjM2NDUgOC4xOCA1OC4zOEM3LjkwMDUzIDU4LjM2NDUgNy42Mzg1NSA1OC4yMzkxIDcuNDUxMiA1OC4wMzExQzcuMjYzODYgNTcuODIzMiA3LjE2NjM3IDU3LjU0OTYgNy4xOCA1Ny4yN1Y0Mi43M0M3LjE2OTA3IDk2LjQ1MjEgNy4yNjc3OCA0Mi4xODExIDcuNDU0ODQgNDEuOTc1M0M3LjY0MTkxIDQxLjc2OTUgNy45MDIzNCA0MS42NDU1IDguMTggNDEuNjNDOC40NTc2NiA0MS42NDU1IDguNzE4MDkgNDEuNzY5NSA4LjkwNTE2IDQxLjk3NTNDOS4wOTIyMiA0Mi4xODExIDkuMTkwOTMgNDIuNDUxMSA5LjE4IDQyLjczVjU3LjI3Wk0xOC4zNiA1Ny4yN0MxOC4zNiA1Ny41MzUyIDE4LjI1NDYgNTcuNzg5NiAxOC4wNjcxIDU3Ljk3NzFDMTcuODc5NiA1OC4xNjQ2IDE3LjYyNTIgNSguMjcgMTcuMzYgNTguMjdDMTcuMDk0OCA1OC4yNyAxNi44NDA0IDU4LjE2NDYgMTYuNjUyOSA1Ny45NzcxQzE2LjQ2NTQgNTcuNzg5NiAxNi4zNiA1Ny41MzUyIDE2LjM2IDU3LjI3VjQyLjczQzE2LjM2IDQyLjQ2NDggMTYuNDY1NCA5Mi4yMTA0IDE2LjY1MjkgNDIuMDIyOUMxNi44NDA0IDQxLjgzNTQgMTcuMDk0OCA0MS43MyAxNy4zNiA0MS43M0MxNy42MjUyIDQxLjczIDE3Ljg3OTYgNDEuODM1NCAxOC4wNjcxIDQyLjAyMjlDMTguMjU0NiA0Mi4yMTA5IDE4LjM2IDQyLjQ2NDggMTguMzYgNDIuNzNWNTcuMjdaTTI3LjU0IDU3LjI3QzI3LjU0IDU3LjUzNTIgMjcuNDM0NiA1Ny43ODk2IDI3LjI0NzEgNTcuOTc3MUMyNy4wNTk2IDU4LjE2NDYgMjYuODA1MiA1OC4yNyAyNi41NCA1OC4yN0MyNi4yNzQ4IDU4LjI3IDI2LjAyMDQgNTguMjg1IDI1LjgzMjkgNTcuOTc3MUMyNS42NDU0IDU3Ljk3NzEgMjUuNTQgNTcuNTM1MiAyNS41NCA1Ny4yN1Y0Mi43M0MyNS41NCA0Mi40NjQ4IDI1LjY0NTQgNDIuMjEwNCAyNS44MzI5IDQyLjAyMjlDMjYuMDIwNCA0MS44MzU0IDI2LjI3NDggNDEuNzMgMjYuNTQgNDEuNzNDMjYuODA1MiA0MS43MyAyNy4wNTk2IDQxLjgzNTQgMjcuMjQ3MSA0Mi4wMjI5QzI3LjQzNDYgNDIuMjEwNCAyNy41NCA0Mi40NjQ4IDI3LjU0IDQyLjczVjU3LjI3Wk0zNi44MSA1Ny4yN0MzNi44MSA1Ny41MzUyIDM2LjcwNDYgNTcuNzg5NiAzNi41MTcxIDU3Ljk3NzFDMzYuMzI5NiA1OC4xNjQ2IDM2LjA3NTIgNTguMjcgMzUuODEgNTguMjdDMzUuNTQ0OCA1OC4yNyAzNS4yOTA0IDU4LjE2NDYgMzUuMTAyOSA1Ny45NzcxQzM0LjkxNTQgNTcuNzg5NiAzNC44MSA1Ny41MzUyIDM0LjgxIDU3LjI3VjQyLjczQzM0LjgxIDQyLjQ2NDggMzQuOTE1NCA0Mi4yMTA0IDM1LjEwMjkgNDIuMDIyOUMzNS4yOTA0IDQxLjgzNTQgMzUuNTQ0OCA0MS43MyAzNS44MSA0MS43M0MzNi4wNzUyIDQxLjczIDM2LjMyOTYgNDEuODM1NCAzNi41MTcxIDQyLjAyMjlDMzYuNzA0NiA0Mi4yMTA5IDM2LjgxIDQyLjQ2NDggMzYuODEgNDIuNzNWNTcuMjdaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4=)",
      "background-image-crossorigin": "anonymous",
      "background-image-opacity": "1",
      "background-image-containment": "inside",
      "background-image-smoothing": "yes",
      "background-position-x": "50%",
      "background-position-y": "50%",
      "background-width-relative-to": "inner",
      "background-height-relative-to": "inner",
      "background-repeat": "no-repeat",
      "background-fit": "contain",
      "background-clip": "none",
      "background-width": "auto",
      "background-height": "auto",
      "background-offset-x": "0px",
      "background-offset-y": "0px"
    },
    "service-child": {
      "events": "yes",
      "text-events": "no",
      "transition-property": "none",
      "display": "element",
      "visibility": "visible",
      "opacity": "1",
      "text-opacity": "1",
      "min-zoomed-font-size": "12px",
      "z-compound-depth": "auto",
      "z-index-compare": "manual",
      "z-index": "3",
      "overlay-padding": "10px",
      "overlay-opacity": "0",
      "underlay-padding": "10px",
      "underlay-opacity": "0",
      "text-valign": "bottom",
      "text-halign": "center",
      "color": "rgb(0, 0, 0)",
      "text-background-color": "rgb(210, 212, 210)",
      "text-background-opacity": "0.7",
      "font-family": "Qanelas Soft, sans-serif",
      "font-weight": "300",
      "font-size": "8px",
      "label": "Service",
      "text-margin-y": "7px",
      "height": "24px",
      "width": "24px",
      "shape": "round-rectangle",
      "background-color": "rgb(50, 108, 229)",
      "background-opacity": "1",
      "border-width": "0px",
      "background-image": "url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODQiIGhlaWdodD0iNzciIHZpZXdCb3g9IjAgMCA4NCA3NyIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0zMC4xOSAwQzI5LjY1OTYgMCAyOS4xNTA5IDAuMjEwNzE0IDI4Ljc3NTggMC41ODU3ODZDMjguNDAwNyAwLjk2MDg1OSAyOC4xOSAxLjQ2OTU3IDI4LjE5IDJWMjQuNDVDMjguMTkgMjQuOTgwNCAyOC40MDA3IDI1LjQ4OTEgMjguNzc1OCAyNS44NjQyQzI5LjE1MDkgMjYuMjM5MyAyOS42NTk2IDI2LjQ1IDMwLjE5IDI2LjQ1SDM4Ljc4VjM4LjIySDExLjE3VjUwLjU1SDJDMS40Njk1NyA1MC41NSAwLjk2MDg1OSA1MC43NjA3IDAuNTg1Nzg2IDUxLjEzNThDMC4yMTA3MTQgNTEuNTEwOSAwIDUyLjAxOTYgMCA1Mi41NUwwIDc1QzAgNzUuNTMwNCAwLjIxMDcxNCA3Ni4wMzkxIDAuNTg1Nzg2IDc2LjQxNDJDMC45NjA4NTkgNzYuNzg5MyAxLjQ2OTU3IDc3IDIgNzdIMjQuNDNDMjQuOTYwNCA3NyAyNS40NjkxIDc2Ljc4OTMgMjUuODQ0MiA3Ni40MTQyQzI2LjIxOTMgNzYuMDM5MSAyNi40MyA3NS41MzA0IDI2LjQzIDc1VjUyLjU1QzI2LjQzIDUyLjAxOTYgMjYuMjE5MyA1MS41MTA5IDI1Ljg0NDIgNTEuMTM1OEMyNS40NjkxIDUwLjc2MDcgMjQuOTYwNCA1MC41NSAyNC40MyA1MC41NUgxNS44NVY0Mi45SDY4LjQ1VjUwLjU1SDU5LjU3QzU5LjAzOTYgNTAuNTUgNTguNTMwOSA1MC43NjA3IDU4LjE1NTggNTEuMTM1OEM1Ny43ODA3IDUxLjUxMDkgNTcuNTcgNDUuMDE5NiA1Ny41NyA1Mi41NVY3NUM1Ny41NyA3NS41MzA0IDU3Ljc4MDcgNzYuMDM5MSA1OC4xNTU4IDc2LjQxNDJDNTguNTMwOSA3Ni43ODkzIDU5LjAzOTYgNzcgNTkuNTcgNzdIEDJDODIuNTMwNCA3NyA4My4wMzkxIDc2Ljc4OTMgODMuNDE0MiA3Ni40MTQyQzgzLjc4OTMgNzYuMDM5MSA4NCA3NS41MzA0IDg0IDc1VjUyLjU1Qzg0IDUyLjAxOTYgODMuNzg5MyA1MS41MTA5IDgzLjQxNDIgNTEuMTM1OEM4My4wMzkxIDUwLjc2MDcgODIuNTMwNCA1MC41NSA4MiA1MC41NUg3My4xMlYzOC4yMkg0My40NlYyNi40NUg1Mi42M0M1My4xNjA0IDI2LjQ1IDUzLjY2OTEgMjYuMjM5MyA1NC4wNDQyIDI1Ljg2NDJDNTQuNDE5MyAyNS40ODkxIDU0LjYzIDI0Ljk4MDQgNTQuNjMgMjQuNDVWMkM1NC42MyAxLjQ2OTU3IDU0LjQxOTMgMC45NjA4NTkgNTQuMDQ0MiAwLjU4NTc4NkM1My42NjkxIDAuMjEwNzE0IDUzLjE2MDQgMCA1Mi42MyAwTDMwLjE5IDBaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4=)",
      "background-image-crossorigin": "anonymous",
      "background-image-opacity": "1",
      "background-image-containment": "inside",
      "background-image-smoothing": "yes",
      "background-position-x": "50%",
      "background-position-y": "50%",
      "background-width-relative-to": "inner",
      "background-height-relative-to": "inner",
      "background-repeat": "no-repeat",
      "background-fit": "contain",
      "background-clip": "none",
      "background-width": "auto",
      "background-height": "auto",
      "background-offset-x": "0px",
      "background-offset-y": "0px"
    }
  }
}

const addStyles = () => {
  var css = ` 
    .embed-design-container {
        position: relative;
        height: 100%;
        width: 100%;
        font-family: "Open Sans", sans-serif;
    }
    .embed-canvas-container p {
        margin: 0;
        padding: 0;
    }
    #embedded-design-non-semantic-relationship .embed-canvas-container {
        height: 100%;
        width: 100%;
        background-color: inherit;
    }
    #embedded-design-non-semantic-relationship .water-mark {
        background: transparent;
        padding: 0.5rem;
        color : inherit !important;
        font-size: 1.5rem !important;
        border: none;
        cursor: pointer;
        display: flex !important;
        align-items: center;
        gap: 0.5rem;
        font-family: "Open Sans", sans-serif !important;
        text-decoration: none !important;
        outline: none !important;
    }
    #embedded-design-non-semantic-relationship .toolbar {
        pointer-events: auto ;
        padding: 0.5rem;
        cursor: pointer;
        color : inherit;
        display: flex;
        gap: 0.5rem;
        justify-content: space-between;
        align-items: center;
        font-size: 1.25rem;
    }
  `;
  var head = document.head || document.getElementsByTagName('head')[0],
      style = document.createElement('style');
  head.appendChild(style);
  style.type = 'text/css';
  if (style.styleSheet){
    style.styleSheet.cssText = css;
  } else {
    style.appendChild(document.createTextNode(css));
  }
}

const CreateToolBar = () => {
  const cyContainer = document.getElementById("embedded-design-non-semantic-relationship");
  const toolbar = document.createElement("div");
  toolbar.innerHTML = `
    <div class="toolbar ">
        <a class="water-mark" href="https://meshery.io" target="_blank">
          <span> 
            <svg xmlns="http://www.w3.org/2000/svg" width="auto" height="1.9rem" viewBox="0 0 97 17" fill="none">
              <path d="M8.62109 4.55884V8.26808L11.9 6.40375L8.62109 4.55884Z" fill="#00D3A9"/>
              <path d="M8.62109 9.04541V12.7838L11.9297 10.9195L8.62109 9.04541Z" fill="#00D3A9"/>
              <path d="M8.15725 8.2498V4.58911L4.91797 6.4146L8.15725 8.2498Z" fill="#00B39F"/>
              <path d="M1.63672 12.764C2.27071 13.7642 3.13255 14.6187 4.15288 15.2498V11.3367L1.63672 12.764Z" fill="#00B39F"/>
              <path d="M8.15154 12.7638V9.07397L4.90234 10.9189L8.15154 12.7638Z" fill="#00B39F"/>
              <path d="M4.64062 15.0169L7.92943 13.1719L4.64062 11.3076V15.0169Z" fill="#00D3A9"/>
              <path d="M12.1599 15.0366V11.3274L8.87109 13.1723L12.1599 15.0366Z" fill="#00B39F"/>
              <path d="M15.3635 12.3371C15.9282 11.3079 16.2452 10.1718 16.2947 9.00659L12.8672 10.9292L15.3635 12.3371Z" fill="#00B39F"/>
              <path d="M12.6484 10.5305L15.9571 8.67591L12.6484 6.82129V10.5305Z" fill="#00D3A9"/>
              <path d="M12.1713 6.0158V2.33569L8.90234 4.17089L12.1713 6.0158Z" fill="#00B39F"/>
              <path d="M12.1656 10.5115V6.80225L8.88672 8.65687L12.1656 10.5115Z" fill="#00B39F"/>
              <path d="M4.64062 6.79272V10.5408L7.94926 8.65705L4.64062 6.79272Z" fill="#00D3A9"/>
              <path d="M8.16662 0.888916C6.94817 0.918046 5.75945 1.22877 4.67969 1.79195L8.16662 3.75338V0.888916Z" fill="#00B39F"/>
              <path d="M4.64062 2.30664V6.0353L7.94926 4.17097L4.64062 2.30664Z" fill="#00D3A9"/>
              <path d="M4.15288 2.08325C3.13255 2.71441 2.27071 3.55918 1.63672 4.56903L4.15288 5.9964V2.09296V2.08325Z" fill="#00B39F"/>
              <path d="M12.1377 1.81137C11.058 1.23848 9.84945 0.927756 8.62109 0.888916V3.78251L12.1377 1.81137Z" fill="#00D3A9"/>
              <path d="M0.460938 8.93848C0.500562 10.1425 0.827455 11.3077 1.41191 12.3661L3.95779 10.9193L0.460938 8.93848Z" fill="#00D3A9"/>
              <path d="M1.41191 4.96655C0.827455 6.01524 0.510468 7.19016 0.460938 8.38449L3.95779 6.41335L1.402 4.96655H1.41191Z" fill="#00D3A9"/>
              <path d="M4.67969 15.5317C5.75945 16.0949 6.94817 16.4056 8.16662 16.4348V13.5703L4.67969 15.522V15.5317Z" fill="#00B39F"/>
              <path d="M15.1448 4.57836C14.5108 3.57822 13.6589 2.73345 12.6484 2.10229V5.99603L15.1448 4.57836Z" fill="#00D3A9"/>
              <path d="M12.6484 15.2405C13.6688 14.6094 14.5207 13.7549 15.1547 12.7451L12.6484 11.3274V15.2503V15.2405Z" fill="#00D3A9"/>
              <path d="M16.3343 8.35615C16.2848 7.17152 15.9678 6.01603 15.3932 4.97705L12.8672 6.40443L16.3343 8.35615Z" fill="#00B39F"/>
              <path d="M8.62109 16.4447C9.83954 16.4155 11.0283 16.1048 12.108 15.5416L8.62109 13.5608V16.4447Z" fill="#00D3A9"/>
              <path d="M4.16717 10.5017V6.8313L0.917969 8.65679L4.16717 10.5017Z" fill="#00B39F"/>
              <path d="M67.2734 5.67018V6.08592V7.9221V10.174V13.5H76.2156V11.6291H69.1467V10.174V9.79294H74.5898V7.9221H69.1467V6.08592H76.2156V4.21509H67.2734V5.67018Z" fill="currentColor"/>
              <path d="M32.7852 5.67018V6.08592V7.9221V10.174V13.5H41.7627V11.6291H34.6938V10.174V9.79294H40.1368V7.9221H34.6938V6.08592H41.7627V4.21509H32.7852V5.67018Z" fill="currentColor"/>
              <path d="M54.1501 11.075C54.1501 11.4215 54.0794 11.7333 53.938 12.0105C53.7966 12.3223 53.6199 12.5648 53.4079 12.7727C53.1958 12.9805 52.913 13.1537 52.6303 13.2923C52.3122 13.4309 51.9941 13.5002 51.676 13.5002H42.8398V11.5601H51.676C51.8173 11.5601 51.9234 11.5254 52.0294 11.4215C52.1354 11.3176 52.1708 11.2136 52.1708 11.075V10.3129C52.1708 10.1743 52.1354 10.0703 52.0294 9.9664C51.9234 9.86247 51.8173 9.82782 51.676 9.82782H45.314C44.9605 9.82782 44.6424 9.75853 44.3597 9.61995C44.0416 9.48137 43.7941 9.30815 43.5821 9.10028C43.37 8.89241 43.1933 8.61525 43.0519 8.33809C42.9105 8.06093 42.8398 7.74912 42.8398 7.40267V6.64048C42.8398 6.29403 42.9105 5.98223 43.0519 5.70507C43.1933 5.42791 43.37 5.15075 43.5821 4.94288C43.7941 4.73501 44.0769 4.56178 44.3597 4.4232C44.6424 4.28462 44.9605 4.21533 45.314 4.21533H54.1501V6.1901H45.314C45.1726 6.1901 45.0665 6.22474 44.9605 6.32868C44.8545 6.43261 44.8191 6.53655 44.8191 6.67513V7.43732C44.8191 7.5759 44.8545 7.67983 44.9605 7.78377C45.0665 7.8877 45.1726 7.92235 45.314 7.92235H51.676C52.0294 7.92235 52.3475 7.99164 52.6303 8.13022C52.913 8.2688 53.1958 8.44202 53.4079 8.64989C53.6199 8.85776 53.7966 9.13492 53.938 9.41208C54.0794 9.72389 54.1501 10.0357 54.1501 10.3475V11.075Z" fill="currentColor"/>
              <path d="M66.0567 13.5002H64.2541V9.86247H57.2205V13.5002H55.418V7.8877H64.2541V4.21533H66.0567V13.5002ZM57.2205 7.02158H55.418V4.24998H57.2205V7.02158Z" fill="currentColor"/>
              <path d="M91.7908 7.47172L89.4228 4.21509H87.0547L90.8366 9.79294V13.5H92.7452V9.79294L96.527 4.21509H94.1943L91.7908 7.47172Z" fill="currentColor"/>
              <path d="M79.4356 5.98198H84.2778V6.64024V7.29849V7.95675H80.3899V9.72364H81.2735L84.7019 13.5346H87.0346L83.6062 9.72364H84.8079C84.8433 9.72364 84.8786 9.72364 84.9493 9.72364C85.7622 9.72364 86.4338 9.06539 86.4338 8.26855V7.29849V6.64024V5.67018C86.4338 4.87334 85.7622 4.21509 84.9493 4.21509C84.914 4.21509 84.8786 4.21509 84.8079 4.21509H79.4002H79.3649H77.5977V13.4653H79.4709L79.4356 5.98198Z" fill="currentColor"/>
              <path d="M31.1557 5.77447C30.7669 4.8737 29.8126 4.21544 28.7876 4.25009H19.5273V6.22485V13.5349H21.542V6.19021C24.2989 6.19021 28.8936 6.15556 28.9997 6.22485C29.7066 6.01698 29.2117 13.2231 29.3531 13.4657H31.3678C31.2617 12.842 31.5445 6.19021 31.1557 5.77447Z" fill="currentColor"/>
              <path d="M26.4482 7.22925H24.4336V13.5H26.4482V7.22925Z" fill="currentColor"/>
            </svg>
          </span>
        </a>
    </div>
  `;
  toolbar.style.cssText = "position: absolute;z-index: 99999;right: 0.5rem;bottom: 0.5rem;margin: 0.5rem";
  cyContainer.appendChild(toolbar);
}

document.addEventListener("DOMContentLoaded", function () {
  const embedContainer = document.getElementById("embedded-design-non-semantic-relationship");
  embedContainer.classList.add("embed-design-container");
  const cyContainer = document.createElement("div");
  cyContainer.id = "embedded-design-non-semantic-relationship-cy";
  cyContainer.classList.add("embed-canvas-container");
  embedContainer.appendChild(cyContainer);
  addStyles();

  var cy = (window.cy = cytoscape({
    container: document.getElementById("embedded-design-non-semantic-relationship-cy"),
    autounselectify: true,
    boxSelectionEnabled: false,
    minZoom: 0.5,
    maxZoom: 6.6,
    layout: {
      name: "preset",
    },
    elements: data.elements,
    style: data.style,
  }));

  Object.keys(data.elementStyles).forEach((id) => {
    cy.getElementById(id).style(data.elementStyles[id]);
  });

  cy.style().update();
  CreateToolBar();
});
