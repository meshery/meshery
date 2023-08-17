import React from "react";
import BaseInput from "./CustomBaseInput";

const CustomFileWidget = (props) => {
  const inputType = "file";

  /**
   * @param {React.ChangeEvent<HTMLInputElement>} event
   * @return {Promise<string>} - file in data url string
   */
  const processFile = (event) => {
    const f = event.target.files[0];
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (event) => resolve(event.target.result);
      reader.onerror = (error) => reject(error);
      reader.readAsDataURL(f);
    })
      .then(props.onChange)
      .catch((error) => {
        console.error("Error processing file:", error);
      });
  };

  return <BaseInput {...props} options={{ ...props.options, inputType }} onChange={processFile} />;
};

export default CustomFileWidget;

