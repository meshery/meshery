import React from "react";
import BaseInput from "./CustomBaseInput";

const CustomFileWidget = props => {
  const inputType = "file";

  /**
   * @param {React.ChangeEvent<HTMLInputElement>} event
   * @return {Promise<string>} - file in data url string
   */
  const processFile = event => {
    const f = event.target.files[0];
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = event => resolve(event.target.result);
      reader.readAsDataURL(f);
    }).then(props.onChange);
  };

  return (
    <BaseInput
      {...props}
      options={{ ...props.options, inputType}}
      onChange={processFile}
    />
  );
};

export default CustomFileWidget;
