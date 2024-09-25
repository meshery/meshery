---
layout: page
title: Contributing to Meshery UI - Sistent
permalink: project/contributing/contributing-ui-sistent
abstract: How to contribute to the Meshery's web-based UI using sistent design system.
language: en
display-title: false
type: project
category: contributing
list: include
---

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisite Reading</strong></p>
  <ol><li><a ahref="contributing-ui">Contributing to Meshery UI</a></li></ol>
</div>

## <a name="contributing-ui-sistent">Contributing to Meshery UI - Sistent</a>

`Sistent` is the Layer5 Design System - an open source design system that offers building blocks to create consistent, accessible, and user-friendly interfaces. It's aimed at developers who want to design applications aligned with the same brand and ensure a uniform user experience across different products.

Sistent leverages Material UI libraries and provides a custom theme on top of it for a consistent look and feel. It includes components, icons, and design tokens that developers can readily integrate into their applications. By using Sistent, developers can save time and effort while maintaining a high-quality user experience throughout Meshery.

**Example**
- This `Modal` is a Sistent Standard Modal used in validation of design.

<a href="{{ site.baseurl }}/assets/img/sistent/sistent-modal.png">
<img style= "width: 600px;" src="{{ site.baseurl }}/assets/img/sistent/sistent-modal.png" />
</a>

### How to use Sistent in Meshery UI

**Table of Contents**

- [Usage](#usage)
- [How to use Sistent tokens/theme colors](#how-to-use-tokenscolors-from-sistent-theme)
- [Examples](#examples)

The Sistent design system includes a variety of base components such as `Button`, `Textfield`, `Checkbox`, and more, which can be found [here](https://github.com/layer5io/sistent/tree/master/src/base)
 Additionally, it provides custom components like `Modal`, `TransferList`, and others, which are available [here](https://github.com/layer5io/sistent/tree/master/src/custom)

###  Usage

- **Import any base or custom component from `@layer5/sistent`**. Here's how you can do it:

```javascript
import {Button} from `@layer5/sistent`

function MyComponent() {
    return (
    <div>
      <Button
        variant="contained"
        onClick={onClick}
      >
      {/* Text to display */}
      </Button>
    </div>
  );
}

export default MyComponent;
```

- **Wrap the Sistent component with `UsesSistent` wrapper**.<br/>
This `UsesSistent` wrapper ensures that the correct theme is applied to the Sistent component based on the current theme of the Meshery UI.
Here how you can do it:

```javascript
import { UsesSistent } from '<path>/SistentWrapper';

function MyComponent() {
    return (
    <div>
    <UsesSistent>
      <Button
        variant="contained"
        onClick={onClick}
      >
      {/* Text to display */}
      </Button>
    </UsesSistent>
    </div>
  );
}

export default MyComponent;
```
### How to use `Tokens/Colors` from Sistent theme

Let's start with a few of the common terms that we will come across frequently, as understanding what they mean will inform us of applicable use cases and proper procedures that should not be overlooked.

- **Theme:**
 A theme provides a cohesive and consistent look and feel for a product, achieved through harmonious color palettes, legible fonts, and layout patterns. Sistent specifies both light and dark themes.

- **Value:**
A value is a unique visual attribute assigned to a token via themes, such as hex codes or RGBA values, used to highlight specific colors. Avoid referencing exact values directly to ensure consistency; instead, use tokens to manage and implement reusable values.

- **Tokens:**
Tokens serve as a shared language between design and development, detailing how to build user interfaces. Tokens represent context (background, text, component), role (success, warning, brand, inverse), and modifiers (secondary, tertiary, hover) derived from the [color palette](https://github.com/layer5io/sistent/blob/master/src/theme/palette.ts).

- **Role:**
Roles specify the context for applying colors. Different roles can share the same value but will have different use cases due to the token structure. These values can vary depending on the current theme.

#### How to use these tokens/colors from sistent theme.

- Import `useTheme` hook from `@layer5/sistent` to access the current theme
- Use the `UsesSistent` wrapper to ensure the Sistent theme is applied to your components.
- Utilize the `useTheme` hook to access theme properties and apply them to your components, such as setting the background style of the `Button`.

```javascript
import { UsesSistent } from '<path>/SistentWrapper';
import {Button, useTheme} from `@layer5/sistent`

function MyComponent() {
 const theme = useTheme()
    return (
    <div>
    <UsesSistent>
    {% raw %}
      <Button
        variant="contained"
        onClick={onClick}
        style = {{
            background: theme.palette.background.default,
            text: theme.palette.text.secondary
        }}
      >
      {/* Text to display */}
      </Button>
       {% endraw %}
    </UsesSistent>
    </div>
  );
}

export default MyComponent;
```

{% include alert.html type="info" title="NOTE" content="The <a href='https://github.com/layer5io/sistent/blob/master/src/theme/palette.ts'>tokens</a> and their roles are specific to their use. For example, use the background palette for background styles and the text palette for text styles." %}
