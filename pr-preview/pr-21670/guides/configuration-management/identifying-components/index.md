# Identifying Meshery Components

> A guide to help you identify and understand the various component icons, shapes, and visual styles used across the Meshery UI.

Source: /pr-preview/pr-21670/guides/configuration-management/identifying-components/

Ever wondered what the different icons and shapes in Meshery represent? Whether you're looking at a dashboard, a settings page, or a design, you'll encounter a rich library of visual elements. This guide is here to help you understand what they mean.

The [components](/pr-preview/pr-21670/concepts/logical/components/) in Meshery fall into two fundamental categories, distinguished by whether they can be orchestrated (managed) by Meshery during deployment:

- **Semantic Components (Orchestratable):** These represent actual infrastructure resources that Meshery can understand and manage during deployment. Examples include Kubernetes resources (like Pods and Services), databases, and other infrastructure components. Meshery will actively manage their lifecycle during deployment.

- **Non-semantic Components (Annotation):** These are visual elements used for documentation and organization, such as text boxes, arrows, shapes, and comments. Meshery ignores these during deployment as they don't represent actual infrastructure.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Visual Customization</div>


All components, whether semantic or non-semantic, support rich visual customization. For example, you can change the color of a Kubernetes Pod icon, modify its shape, or customize its background - it's all configurable!
</div>


## Semantic Components

These components represent real infrastructure that Meshery can manage. They can be either built-in (like Kubernetes components) or custom components that you [create](/pr-preview/pr-21670/guides/configuration-management/creating-models/).

### Kubernetes Components

While Kubernetes components are commonly used, they follow the same principles as all other semantic components. They have a default distinct visual style to help you instantly recognize them:

- **Uniform Color Scheme:** Kubernetes component icons typically use a **distinctive blue background** as a standard identifier.
- **Standardized Icon Structure:** The fundamental structure is consistent: an outer container shape with the blue background, encompassing a unique inner white symbol.
- **Meaningful Inner Symbols:** The white symbol inside each icon is the crucial unique identifier for that specific Kubernetes Kind, often inspired by the core function of the resource.

<a href="./images/k8s_component.gif" target="_blank">
  <img src="./images/k8s_component.gif" style="width:50%; height:auto;" alt="Kubernetes components in Meshery">
</a>

### Integrated Technologies

Meshery supports various technologies (like AWS, Prometheus, Istio, KEDA, etc.) with their official icons. These components have the same orchestratable capabilities as Kubernetes components.

<a href="./images/AWS-models.png" target="_blank">
  <img src="./images/AWS-models.png" style="width:50%; height:auto;" alt="AWS models in Meshery">
</a>

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Exploring All Integrations</div>


This guide covers the visual style of components. For a complete catalog of all technologies that Meshery integrates, visit the integrations directory. <strong><a href='/pr-preview/pr-21670/extensions/models/'>Explore All Integrations</a></strong>
</div>


## Non-semantic Components

These components help you document and organize your designs without affecting the actual infrastructure. They include:

- Text boxes and comments for documentation
- Shapes and containers for visual grouping
- Lines and arrows for showing relationships
- Labels and tags for organization

While these components are ignored during deployment, they support the same visual customization options as semantic components.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Edge Components</div>


To learn more about edge components and their visual styles, visit the <strong><a href='/pr-preview/pr-21670/guides/configuration-management/edges-guide/'>Edge Components Guide</a></strong>
</div>


## Foundational Elements

In addition to the rich library of predefined icons, Meshery also provides a core set of foundational geometric shapes. These are not tied to any specific technology but are intended for representing abstract concepts or for simple annotations within your designs where a specific component icon doesn't apply.

<a href="./images/shapes.png" target="_blank">
  <img src="./images/shapes.png" style="width:50%; height:auto;" alt="Generic shapes palette in Meshery">
</a>

<details>
<summary><strong>Foundational Geometric Shapes</strong></summary>
<br>
These are the quintessential building blocks for many diagrams – your circles, squares, triangles, and basic polygons. They offer simple, clean, and universally understood forms for a wide range of uses.

<style>

    .extension-guides-container {
      display: flex;
      flex-wrap: wrap;
      gap: 20px;
      justify-content: space-between;
      margin-top: 2rem;

    }
    .extension-guide-card {
       display: flex;
     flex-direction: column;
       max-width: 20rem;
       max-height: 30rem;
       gap: 1rem;
  }
    .extension-guide-svg-container {
       height: auto;
       display: flex;
       flex-direction: column;
       align-items: center;
       gap: 0.455rem;
     flex-basis: 30%;
  }
    .extension-guide-svg-container img {
     width: 50%;
     height: auto;
  }
    .extension-guide-details {
       display: flex;
       flex-direction: column;
       flex-basis: 77%;
       gap: 10px;
     text-align: center;
  }
    @media (max-width: 767px) {
      .extension-guide-container {
        flex-direction: column;
      }
      .extension-guide-svg-container {
        gap: 0.3rem;
      }
      .extension-guide-svg-container img {
        width: 40%;
      }
      .extenion-guide-card{
        max-width: 30rem;
        flex-direction: column;
      }
     }
  </style>








<div class="extension-guides-container">
  

  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Circle.svg" alt="Circle">
        
        <div style="text-align:center;">Circle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>The default shape for nodes, displayed as an oval. This shape is useful for general-purpose nodes where no special distinction is required.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Triangle.svg" alt="Triangle">
        
        <div style="text-align:center;">Triangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A simple three-sided shape, often used to represent hierarchical data or directional flows. It is useful for indicating movement, direction, or hierarchical relationships within a network.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Heptagon.svg" alt="Heptagon">
        
        <div style="text-align:center;">Heptagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A polygon with seven straight sides and angles. This shape is less common but can be used to represent specific or unique data types that require distinction from other polygonal shapes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Pentagon.svg" alt="Pentagon">
        
        <div style="text-align:center;">Pentagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A five-sided polygon often used for specific classifications of data. Its unique shape helps to quickly distinguish it from other polygonal shapes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundRectangle.svg" alt="Round Rectangle">
        
        <div style="text-align:center;">Round Rectangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>Similar to a rectangle but with rounded corners, providing a softer, more approachable look. This shape is often used to visually distinguish nodes from standard rectangular nodes, indicating a special status or different type of data.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundTriangle.svg" alt="Round Triangle">
        
        <div style="text-align:center;">Round Triangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A triangular shape with rounded corners, providing a softer look while still indicating hierarchy or direction.</div>
        
      </div>
    </div>
  
</div>


</details>

<details>
<summary><strong>Common Representational & Symbolic Shapes</strong></summary>
<br>
This group includes shapes that, by common convention, often evoke a more specific symbolic meaning, like using a "Barrel" for data storage or an "Actor" for a user role. Leveraging these established visual metaphors can make your custom diagrams more intuitive.

<style>

    .extension-guides-container {
      display: flex;
      flex-wrap: wrap;
      gap: 20px;
      justify-content: space-between;
      margin-top: 2rem;

    }
    .extension-guide-card {
       display: flex;
     flex-direction: column;
       max-width: 20rem;
       max-height: 30rem;
       gap: 1rem;
  }
    .extension-guide-svg-container {
       height: auto;
       display: flex;
       flex-direction: column;
       align-items: center;
       gap: 0.455rem;
     flex-basis: 30%;
  }
    .extension-guide-svg-container img {
     width: 50%;
     height: auto;
  }
    .extension-guide-details {
       display: flex;
       flex-direction: column;
       flex-basis: 77%;
       gap: 10px;
     text-align: center;
  }
    @media (max-width: 767px) {
      .extension-guide-container {
        flex-direction: column;
      }
      .extension-guide-svg-container {
        gap: 0.3rem;
      }
      .extension-guide-svg-container img {
        width: 40%;
      }
      .extenion-guide-card{
        max-width: 30rem;
        flex-direction: column;
      }
     }
  </style>








<div class="extension-guides-container">
  

  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Barrel.svg" alt="Barrel">
        
        <div style="text-align:center;">Barrel</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A shape resembling a 3D barrel, with both the top and bottom edges appearing rounded. This shape can be used to represent data storage, such as databases or repositories.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Tag.svg" alt="Tag">
        
        <div style="text-align:center;">Tag</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A shape resembling a label or tag, useful for categorized or labeled data nodes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundTag.svg" alt="Round Tag">
        
        <div style="text-align:center;">Round Tag</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A tag shape with rounded edges, offering a softer visual cue for categorized nodes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/TextBox.svg" alt="TextBox">
        
        <div style="text-align:center;">TextBox</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A rectangular shape representing a text entry or display area, useful in diagrams to show nodes that contain text information or comments.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/GenericNode.svg" alt="Generic Node">
        
        <div style="text-align:center;">Generic Node</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A simple, often circular or rectangular shape used as a placeholder for unspecified or generic data points.</div>
        
      </div>
    </div>
  
</div>


</details>

<details>
<summary><strong>Specialized & Decorative Geometric Shapes</strong></summary>
<br>
When your diagrams require a more distinct visual style or an element for emphasis, this collection offers a variety of options, from complex polygons to stylized forms like "Crescent," "Star," or "XWing." These shapes are less commonly used in Meshery and have no universally accepted meaning.

<style>

    .extension-guides-container {
      display: flex;
      flex-wrap: wrap;
      gap: 20px;
      justify-content: space-between;
      margin-top: 2rem;

    }
    .extension-guide-card {
       display: flex;
     flex-direction: column;
       max-width: 20rem;
       max-height: 30rem;
       gap: 1rem;
  }
    .extension-guide-svg-container {
       height: auto;
       display: flex;
       flex-direction: column;
       align-items: center;
       gap: 0.455rem;
     flex-basis: 30%;
  }
    .extension-guide-svg-container img {
     width: 50%;
     height: auto;
  }
    .extension-guide-details {
       display: flex;
       flex-direction: column;
       flex-basis: 77%;
       gap: 10px;
     text-align: center;
  }
    @media (max-width: 767px) {
      .extension-guide-container {
        flex-direction: column;
      }
      .extension-guide-svg-container {
        gap: 0.3rem;
      }
      .extension-guide-svg-container img {
        width: 40%;
      }
      .extenion-guide-card{
        max-width: 30rem;
        flex-direction: column;
      }
     }
  </style>








<div class="extension-guides-container">
  

  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/BottomRoundRectangle.svg" alt="Bottom Round Rectangle">
        
        <div style="text-align:center;">Bottom Round Rectangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A unique shape where only the bottom corners are rounded. This can be useful to indicate hierarchical relationships or to create a visual distinction from other rounded shapes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/ConcaveHexagon.svg" alt="Concave Hexagon">
        
        <div style="text-align:center;">Concave Hexagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A six-sided shape with one or more inward curves, creating a star-like appearance, useful for special types of nodes that need to stand out distinctly.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Crescent.svg" alt="Crescent">
        
        <div style="text-align:center;">Crescent</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A shape resembling a crescent moon, which can be used to represent phases, growth, or transition nodes within a graph.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/CutRectangle.svg" alt="Cut Rectangle">
        
        <div style="text-align:center;">Cut Rectangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A rectangle with corners that are cut off or beveled. This shape gives a more modern and distinct appearance, often used to highlight important nodes or different types of data.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Diamond.svg" alt="Diamond">
        
        <div style="text-align:center;">Diamond</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A four-sided shape oriented at 45 degrees, resembling a diamond. It is typically used for decision points in flowcharts or to represent special conditions or states.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/LeftTriangle.svg" alt="Left Triangle">
        
        <div style="text-align:center;">Left Triangle</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A triangular shape pointing to the left, often used to indicate directional flow or hierarchical steps in processes or data sequences.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Octagon.svg" alt="Octagon">
        
        <div style="text-align:center;">Octagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>An eight-sided polygon, providing another option for representing complex data types. The additional sides make it stand out more compared to hexagons or pentagons.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Parallelogram.svg" alt="Parallelogram">
        
        <div style="text-align:center;">Parallelogram</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A four-sided shape with opposite sides parallel, used to represent data processing steps or functions with directional flow.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Rhomboid.svg" alt="Rhomboid">
        
        <div style="text-align:center;">Rhomboid</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A quadrilateral with opposite sides parallel but not perpendicular, creating a skewed rectangle look. Useful for depicting relationships that aren&#39;t straightforward or direct.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RightRhomboid.svg" alt="Right Rhomboid">
        
        <div style="text-align:center;">Right Rhomboid</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A rhomboid where one angle is a right angle, offering a distinctive look for special types of data connections or relationships.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundDiamond.svg" alt="Round Diamond">
        
        <div style="text-align:center;">Round Diamond</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A diamond shape with rounded corners, used to indicate special conditions or decision points while maintaining a softer visual appeal.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundHeptagon.svg" alt="Round Heptagon">
        
        <div style="text-align:center;">Round Heptagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A seven-sided polygon with rounded corners, less common but useful for unique data representations that require distinction.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundHexagon.svg" alt="Round Hexagon">
        
        <div style="text-align:center;">Round Hexagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A six-sided polygon with rounded corners, providing a less angular look for complex data types.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/RoundOctagon.svg" alt="Round Octagon">
        
        <div style="text-align:center;">Round Octagon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>An eight-sided polygon with rounded corners, offering another distinct and less angular representation for complex data.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/SlantedParallelogram.svg" alt="Slanted Parallelogram">
        
        <div style="text-align:center;">Slanted Parallelogram</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A four-sided shape with opposite sides parallel but slanted, creating a dynamic appearance often used to indicate processes or flows that have a directional aspect.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Star.svg" alt="Star">
        
        <div style="text-align:center;">Star</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A shape with multiple points radiating outwards, typically used for highlighting special or notable nodes. It draws immediate attention, making it ideal for critical data points or focal nodes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Trapezoid.svg" alt="Trapezoid">
        
        <div style="text-align:center;">Trapezoid</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A four-sided shape with one pair of parallel sides, used to indicate hierarchical relationships or steps in a process.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Vee.svg" alt="Vee">
        
        <div style="text-align:center;">Vee</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A shape resembling the letter &#34;V&#34;, often used to show directionality or a splitting point in flowcharts.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/XWing.svg" alt="XWing">
        
        <div style="text-align:center;">XWing</div>
      </div>
      <div class="extension-guide-details">
        
          <div>This shape resembles the Star Wars X-Wing fighter, potentially useful for visually representing nodes that have multiple connections or interactions, symbolizing complexity and centrality.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Polygon.svg" alt="Polygon">
        
        <div style="text-align:center;">Polygon</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A custom point polygon allows the creation of irregular, user-defined polygonal shapes for graph nodes.</div>
        
      </div>
    </div>
  
    <div class="extension-guide-card">
      <div class="extension-guide-svg-container">
        
          <img src="/pr-preview/pr-21670/images/shapes/Plus.svg" alt="Plus(Polygon)">
        
        <div style="text-align:center;">Plus(Polygon)</div>
      </div>
      <div class="extension-guide-details">
        
          <div>A polygonal shape with any number of points, allowing for highly customizable node representations to fit specific data needs.</div>
        
      </div>
    </div>
  
</div>


</details>

## Component Visuals in Different Contexts

To see how this works in practice, let's take the Kubernetes `Deployment` component as an example. Its appearance adapts to different views in the Meshery UI:

1.The full component shape as it appears in a design:

<a href="./images/deployment-shape.png">
    <img src="./images/deployment-shape.png" style="width:30%; height:auto;" alt="Deployment Component Shape">
</a>

2.The simplified icon as it appears in a component selection panel:

<a href="./images/deployment-icon.png">
    <img src="./images/deployment-icon.png" style="width:50%; height:auto;" alt="Deployment icon in a component selection panel">
</a>

3.The icon as seen in a cluster resource overview:

<a href="./images/deployment-dashboard.png">
    <img src="./images/deployment-dashboard.png" style="width:50%; height:auto;" alt="Deployment component in a cluster resource overview">
</a>
