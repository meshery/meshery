import React from 'react';

import { Button } from './Button';

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
  title : 'Components/Button',
  component : Button,
};

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template = (args) => <Button {...args} />;

export const Primary = Template.bind({});
// More on args: https://storybook.js.org/docs/react/writing-stories/args
Primary.args = {
  variant : 'contained',
  color : 'primary',
};

export const Secondary = Template.bind({});
Secondary.args = {
  variant : 'outlined',
  color : 'secondary'
};

export const Large = Template.bind({});
Large.args = {
  size : 'large',
  color : 'error'
};

export const Small = Template.bind({});
Small.args = {
  size : 'small',
  color : 'success'
};