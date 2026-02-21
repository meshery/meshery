import React from 'react';
import BaseInput from './CustomBaseInput';
import { useRef } from 'react';
import { Colors } from '@/themes/app';

const CustomFileWidget = (props) => {
  const inputType = 'file';
  const fileInputUpdated = useRef(false);

  /**
   * @param {React.ChangeEvent<HTMLInputElement>} event
   * @return {Promise<string>} - file in data url string
   */
  const processFile = (event) => {
    fileInputUpdated.current = true;
    const f = event.target.files[0];
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (event) => resolve(event.target.result);
      reader.onerror = (error) => reject(error);
      reader.readAsDataURL(f);
    })
      .then(props.onChange)
      .catch((error) => {
        console.error('Error processing file:', error);
      });
  };

  if (!fileInputUpdated.current && props.options && props.options.default) {
    return (
      <div style={{ position: 'relative' }}>
        <BaseInput {...props} options={{ ...props.options, inputType }} onChange={processFile} />
        {/* Display default value as a mask */}
        <label
          htmlFor={props.id}
          style={{
            paddingInline: '10px',
            paddingBlock: '0',
            minWidth: '50%',
            position: 'absolute',
            bottom: '10px',
            left: '6rem',
            backgroundColor: Colors.keppelGreen,
          }}
        >
          {props.options.default}
        </label>
      </div>
    );
  }

  return <BaseInput {...props} options={{ ...props.options, inputType }} onChange={processFile} />;
};

export default CustomFileWidget;
