// Description: Contains functions related to DOM elements and react components.

import React from 'react';

type HasClassElement = HTMLElement | null | undefined;

// recursively check if element or any of its parent has the class
export const hasClass = (element: HasClassElement, className: string): boolean => {
  try {
    if (typeof element?.className == 'string' && element?.className?.includes(className)) {
      return true;
    }
    if (element?.parentElement) {
      return hasClass(element.parentElement, className);
    }
  } catch (e) {
    console.error(`Error in hasClass while checking for ${className} in `, element, e);
  }
  return false;
};

// recursively go through component and its children and add the class to each of them
// This is required to prevent the clickaway listener from blocking the click event
// on the notification center IconButton. Create it as a HOC and use React.cloneElement to add the class.
type AddClassRecursivelyProps = {
  children: React.ReactNode;
  className: string;
};

export const AddClassRecursively = ({
  children,
  className,
}: AddClassRecursivelyProps): React.ReactNode => {
  return React.Children.map(children, (child) => {
    if (React.isValidElement(child)) {
      const element = child as React.ReactElement<any>;
      const existingClassName = (element.props as any).className || '';

      return React.cloneElement(
        element,
        {
          className: `${existingClassName} ${className}`.trim(),
        } as any,
        AddClassRecursively({
          children: (element.props as any).children,
          className,
        }),
      );
    }

    // if child is a svg or animated svg string
    return child;
  });
};
