# Shape Builder

> Create custom polygon shapes for Meshery components.

---

LLMS index: [llms.txt](/llms.txt)

---

Shape Builder is a visual tool for creating custom polygon shapes for Meshery components. It displays coordinates in the format Kanvas uses to render a polygon.

## Create a shape

1. Open [Shape Builder](https://shapes.meshery.io).
2. Draw at least three points, then select **Close Shape**.
   Check that each coordinate is between -1 and 1, including after scaling. If an x-coordinate exceeds 1, redraw the vertices within the leftmost 520 pixels of the canvas and check again. After your final redraw or scaling change, copy the checked coordinates displayed below the canvas.
3. In your Meshery component, set `styles.shape` to `polygon` and put the coordinates in `styles.shape-polygon-points`.

Polygon coordinates range from -1 to 1. For example, these styles describe a plus shape:

```yaml
styles:
  shape: polygon
  shape-polygon-points: "-0.33 -1 0.33 -1 0.33 -0.33 1 -0.33 1 0.33 0.33 0.33 0.33 1 -0.33 1 -0.33 0.33 -1 0.33 -1 -0.33 -0.33 -0.33"
```

See [Identifying Meshery Components](/guides/configuration-management/identifying-components/) for an overview of component visuals and built-in shapes. Visit the [Shape Builder repository](https://github.com/meshery-extensions/shape-builder) for more information.
