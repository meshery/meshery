# Structured Data Formatter Documentation

2. **Dynamic Formatter**: Since data can vary significantly in structure, it is not practical to create a specific formatter for each kind. Dynamic formatters analyze the schema's structure and apply custom-defined rules for formatting:
   - Text strings are rendered using the _BodySectionRenderer_ (more on this later).
   - Arrays are rendered using the _ArrayRenderer_.
   - Key-value pairs are rendered using the _KeyValueRenderer_.
   - Nested objects are recursively rendered.

## BodySectionRenderer

The BodySectionRenderer is responsible for formatting and rendering raw text strings into React components. During this process, it parses the string to replace external links with `<Link>` components and checks if the link matches predefined sites to render the link accordingly.

## ArrayRenderer

The ArrayRenderer is responsible for rendering an array of items in a recursive manner, presenting them as a bulletized list using the _StructuredDataFormatter_.

## KeyValueRenderer

Object properties with string values are considered key-value pairs and are rendered as such.

## Usage

You can start using the StructuredDataFormatter for your structured data by using the _FormatStructuredData_ component ,
the component takes two props the _data_ to be formatted and optional _propertyFormatter_ object , which the formatter uses to format the specific properties of the structured data using the propertyFormatter .

## Example :

```javascript
export const FormattedMetadata = ({ event, classes }) => {
  const PropertyFormatters = {
    Doc: (value, _metadata) => LinkFormatters.DOC.formatter(value),
    error: (value) => <ErrorMetadataFormatter metadata={value} event={event} classes={classes} />,
    dryRunResponse: (value) => <DryRunResponse response={value} />,
  };
  if (!event || !event.metadata || isEmptyAtAllDepths(event.metadata)) {
    return <EmptyState event={event} />;
  }

  return <FormatStructuredData propertyFormatters={PropertyFormatters} data={event.metadata} />;
};
```
