import React from "react";

export const CustomFieldTemplate=(props) => {
  const {  classNames, help,  description,  children } = props;
  return (
    <div className={classNames}>
      {description}
      {children}
      {help}
    </div>
  )
}