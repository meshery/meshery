# Shape Builder

> Create custom polygon shapes for Meshery components.

---

LLMS index: [llms.txt](/pr-preview/pr-22895/llms.txt)

---

Shape Builder is a visual tool for creating custom polygon shapes for Meshery components. It exports coordinates in the format Kanvas uses to render a polygon.

## Create a shape

1. Open [Shape Builder](https://shapes.meshery.io).
2. Draw a polygon and export its matrix notation.
3. Use the exported coordinates in a `polygon` shape in your Meshery component.

Polygon coordinates range from -1 to 1. For example, this coordinate sequence describes a plus shape:

```text
-0.33 -1 0.33 -1 0.33 -0.33 1 -0.33 1 0.33 0.33 0.33 0.33 1 -0.33 1 -0.33 0.33 -1 0.33 -1 -0.33 -0.33 -0.33
```

See [Identifying Meshery Components](/pr-preview/pr-22895/guides/configuration-management/identifying-components/) for an overview of component visuals and built-in shapes. Visit the [Shape Builder repository](https://github.com/meshery-extensions/shape-builder) for more information.
