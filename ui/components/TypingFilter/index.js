import {
  Autocomplete,
  Chip,
  ContentFilterIcon,
  CrossCircleIcon,
  InputAdornment,
  useTheme,
  Divider,
} from '@layer5/sistent';
import { InputField, Root, DropDown } from './style';
import React, { useState, useEffect } from 'react';
import { iconSmall } from 'css/icons.styles';

function transformData(data) {
  const result = {};

  data.forEach(({ type, value }) => {
    const key = type.toLowerCase();

    if (result[key]) {
      // If the key already exists and is an array, push the new value
      if (Array.isArray(result[key])) {
        result[key].push(value);
      } else {
        // If it's not an array (e.g., status: string), convert to array
        result[key] = [result[key], value];
      }
    } else {
      // For STATUS, store as string; for others, start with array
      result[key] = key === 'status' ? value : [value];
    }
  });

  return result;
}

/**
 * Filter Schema Object
 *
 * The `filterSchema` object defines available filter options for the TypingFilter component.
 * It provides information about different filter categories, their descriptions, and possible values.
 *
 * @typedef {object} FilterSchema
 * @property {object} [CATEGORY_NAME] - An object representing a filter category.
 * @property {string} [CATEGORY_NAME.value] - The filter name used for filtering within the category.
 * @property {string} [CATEGORY_NAME.description] - Description of the filter category.
 * @property {string} [CATEGORY_NAME.type] - The data type of the filter (optional).
 * @property {string[]} [CATEGORY_NAME.values] - Possible values for the filter (optional).
 *
 * @example
 * // Example filter schema with multiple filter categories
 * const filterSchema = {
 *   SEVERITY: {
 *     value: "severity",
 *     description: "Filter by severity",
 *     values: ["Low", "Medium", "High"],
 *   },
 *   STATUS: {
 *     value: "status",
 *     description: "Filter by status",
 *     type: "string",
 *     values: ["Open", "Closed", "In Progress"],
 *   },
 *   CUSTOM_FILTER: {
 *     value: "custom",
 *     description: "Custom filter description",
 *     type: "number",
 *   },
 *   // Add more filter categories as needed
 * };
 */

/**
 * TypingFilter Component
 *
 * A component for real-time filtering and selection with typing. It provides a user-friendly
 * interface for filtering data based on user input.
 *
 * @component
 * @param {object} props - Component props.
 * @param {FilterSchema} filterSchema - The schema defining available filter options.
 * @param {function} handleFilter - A callback function to handle filter changes.
 * @param {string} placeholder - Placeholder text for the input field.
 * @param {object[]} defaultFilters - An array of default filters to initialize the component.
 * @returns {JSX.Element} - A React JSX element representing the TypingFilter component.
 */
const TypingFilter = ({ filterSchema, placeholder, handleFilter, defaultFilters }) => {
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [selectedFilters, setSelectedFilters] = useState(defaultFilters);

  useEffect(() => {
    setSelectedFilters(defaultFilters);
  }, [defaultFilters.length]);

  const getOptions = () => {
    if (inputValue.includes(':')) {
      const [filterTypeInput, searchValue = ''] = inputValue.split(':');
      const filterTypeInputLower = filterTypeInput.toLowerCase().trim();

      // Find the matching schema key
      const schemaKey = Object.keys(filterSchema).find(
        (key) => filterSchema[key].value.toLowerCase() === filterTypeInputLower,
      );

      if (schemaKey) {
        const schema = filterSchema[schemaKey];

        // If the filter has predefined values, filter them by the search term
        if (schema.values && schema.values.length > 0) {
          return schema.values
            .filter((value) => value.toLowerCase().includes(searchValue.toLowerCase().trim()))
            .map((value) => ({
              type: schemaKey,
              value,
              label: `${schema.value}: ${value}`,
            }));
        } else {
          // For custom input filters like AUTHOR
          const customValue = searchValue.trim();
          if (customValue) {
            return [
              {
                type: schemaKey,
                value: customValue,
                label: `${schema.value}: ${customValue}`,
              },
            ];
          }
        }
      }
      return [];
    }

    // When no colon, show all filter types
    return Object.entries(filterSchema).map(([key, data]) => ({
      type: key,
      value: data.value,
      label: `${data.value}: ${data.description}`,
    }));
  };

  const handleSelect = (option) => {
    if (!option) return;

    if (!inputValue.includes(':')) {
      setInputValue(`${option.value}: `);
    } else {
      const newFilter = {
        type: option.type,
        value: option.value,
        label: `${filterSchema[option.type].value}: ${option.value}`,
      };

      const existingFilterIndex = selectedFilters.findIndex(
        (f) => f.type === option.type && f.value === option.value,
      );

      if (existingFilterIndex !== -1) {
        return;
      }

      if (filterSchema[option.type].multiple === false) {
        setSelectedFilters((prev) => [...prev.filter((f) => f.type !== option.type), newFilter]);
      } else {
        setSelectedFilters((prev) => [...prev, newFilter]);
      }

      handleFilter(transformData([...selectedFilters, newFilter]));
      setInputValue('');
    }
  };

  const handleKeyDown = (event) => {
    if (event.key === 'Enter' && inputValue.includes(':')) {
      const [filterType, value] = inputValue.split(':');
      const schemaKey = Object.keys(filterSchema).find(
        (key) => filterSchema[key].value.toLowerCase() === filterType.toLowerCase(),
      );
      if (schemaKey && !filterSchema[schemaKey].values?.length && value.trim()) {
        handleSelect({
          type: schemaKey,
          value: value.trim(),
          label: `${filterSchema[schemaKey].value}: ${value.trim()}`,
        });
      }
    }
  };

  const clearAllFilters = () => {
    setSelectedFilters([]);
    setInputValue('');
    handleFilter({});
  };

  const handleDeleteChip = (details) => {
    setSelectedFilters((prev) => prev.filter((filter) => filter !== details?.option));
    handleFilter(transformData(selectedFilters.filter((filter) => filter !== details?.option)));
  };

  const theme = useTheme();
  return (
    <Root className="mui-fixed">
      <Autocomplete
        style={{ paddingLeft: '0px' }}
        size="small"
        multiple
        open={open}
        onOpen={() => setOpen(true)}
        onClose={(event, reason) => {
          if (reason === 'escape' || reason === 'blur' || reason === 'toggleInput') {
            setOpen(false);
          }
        }}
        PaperComponent={DropDown}
        disableCloseOnSelect
        freeSolo={inputValue.includes(':')}
        inputValue={inputValue}
        onInputChange={(_, newValue) => {
          setInputValue(newValue);
          if (newValue) setOpen(true);
        }}
        options={getOptions()}
        value={selectedFilters}
        noOptionsText="Sorry we dont currently support this filter"
        onChange={(_, __, reason, details) => {
          if (reason === 'removeOption' && details?.option) {
            handleDeleteChip(details);
          } else if (reason === 'selectOption' && details?.option) {
            handleSelect(details.option);
          }
        }}
        getOptionLabel={(option) => option.label}
        isOptionEqualToValue={(option, value) =>
          option.type === value.type && option.value === value.value
        }
        renderOption={(props, option, { index }) => {
          const options = getOptions();
          return (
            <React.Fragment key={option.label}>
              <li
                {...props}
                style={{
                  padding: '0.25rem 2rem',
                  margin: '0.25rem 0.5rem',
                  borderRadius: '0.5rem',
                }}
              >
                {option.label}
              </li>
              {index < options.length - 1 && <Divider />}
            </React.Fragment>
          );
        }}
        renderTags={(value, getTagProps) =>
          value.map((option, index) => (
            <Chip
              {...getTagProps({ index })}
              key={option.type}
              label={`${filterSchema[option.type].value}: ${option.value}`}
              style={{ margin: '0.15rem', maxWidth: '80%', height: 'auto' }}
              size="small"
              sx={{
                height: 'auto',
                '& .MuiChip-label': {
                  display: 'block',
                  whiteSpace: 'normal',
                  lineHeight: '1.5rem',
                  fontSize: '0.75rem',
                  marginBlock: '0.1rem',
                },
              }}
            />
          ))
        }
        clearIcon={
          <CrossCircleIcon
            fill={theme.palette.icon.default}
            {...iconSmall}
            onClick={() => clearAllFilters()}
          />
        }
        renderInput={(params) => {
          const customInputProps = {
            ...params.InputProps,
            startAdornment: (
              <>
                <InputAdornment
                  position="start"
                  style={{ marginInline: '0.25rem', height: 'auto' }}
                >
                  <ContentFilterIcon fill={theme.palette.icon.default} />
                </InputAdornment>
                {params.InputProps.startAdornment}
              </>
            ),
          };

          return (
            <InputField
              {...params}
              placeholder={placeholder}
              onKeyDown={handleKeyDown}
              InputProps={customInputProps}
            />
          );
        }}
      />
    </Root>
  );
};

export default TypingFilter;
