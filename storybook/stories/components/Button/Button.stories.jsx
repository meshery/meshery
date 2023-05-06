import { Button } from './Button';

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
  title : 'Components/Button',
  component : Button,
  tags: ['autodocs'],
};

export const Primary = {
  args: {
    variant : 'contained',
    color : 'primary',
  },
};

export const Secondary = {
  args: {
    variant : 'outlined',
    color : 'secondary'
  },
};

export const Large = {
  args: {
    size : 'large',
    color : 'error'
  },
};

export const Small = {
  args: {
    size : 'small',
    color : 'success'
  },
};