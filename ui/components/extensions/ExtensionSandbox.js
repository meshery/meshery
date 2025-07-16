import React, { useCallback, useEffect, useState, useMemo } from 'react';
import { useLazyGetExtensionsByTypeQuery, useLazyGetFullPageExtensionsQuery } from '../rtk-query/extension';

// Custom hooks to handle the API calls properly
export const useGetCapabilities = () => {
  const [getExtensionsByType] = useLazyGetExtensionsByTypeQuery();
  
  const getCapabilities = useCallback(async (type) => {
    try {
      const result = await getExtensionsByType(type);
      return result;
    } catch (error) {
      console.error('Error getting capabilities:', error);
      return null;
    }
  }, [getExtensionsByType]);
  
  return getCapabilities;
};

export const useGetFullPageExtensions = () => {
  const [getFullPageExtensions] = useLazyGetFullPageExtensionsQuery();
  
  const getFullPageExtensionsFunc = useCallback(async () => {
    try {
      const result = await getFullPageExtensions();
      return result;
    } catch (error) {
      console.error('Error getting full page extensions:', error);
      return null;
    }
  }, [getFullPageExtensions]);
  
  return getFullPageExtensionsFunc;
};

const ExtensionSandbox = ({ 
  type, 
  Extension, 
  extensionType, 
  onExtensionLoad,
  onExtensionError,
  ...props 
}) => {
  const [extensions, setExtensions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Use the custom hooks
  const getCapabilities = useGetCapabilities();
  const getFullPageExtensionsFunc = useGetFullPageExtensions();

  // Load extensions based on type
  useEffect(() => {
    const loadExtensions = async () => {
      setLoading(true);
      setError(null);
      
      try {
        let result;
        if (type === 'full-page') {
          result = await getFullPageExtensionsFunc();
        } else {
          result = await getCapabilities(type);
        }
        
        if (result && result.data) {
          setExtensions(result.data);
          onExtensionLoad?.(result.data);
        }
      } catch (err) {
        setError(err.message);
        onExtensionError?.(err);
      } finally {
        setLoading(false);
      }
    };

    if (type) {
      loadExtensions();
    }
  }, [type, getCapabilities, getFullPageExtensionsFunc, onExtensionLoad, onExtensionError]);

  // Render extensions
  const renderExtensions = useMemo(() => {
    if (loading) {
      return <div>Loading extensions...</div>;
    }

    if (error) {
      return <div>Error loading extensions: {error}</div>;
    }

    if (!extensions || extensions.length === 0) {
      return <div>No extensions found</div>;
    }

    return extensions.map((extension, index) => (
      <div key={extension.id || index} className="extension-item">
        {Extension ? (
          <Extension 
            extension={extension} 
            type={extensionType}
            {...props}
          />
        ) : (
          <div>
            <h3>{extension.name}</h3>
            <p>{extension.description}</p>
          </div>
        )}
      </div>
    ));
  }, [extensions, loading, error, Extension, extensionType, props]);

  return (
    <div className="extension-sandbox">
      {renderExtensions}
    </div>
  );
};

// Helper functions that can be used outside of React components
export const extensionHelpers = {
  // Non-hook utility functions
  validateExtension: (extension) => {
    return extension && extension.name && extension.type;
  },
  
  formatExtensionData: (extensionData) => {
    if (!extensionData) return null;
    
    return {
      id: extensionData.id,
      name: extensionData.name,
      description: extensionData.description,
      type: extensionData.type,
      version: extensionData.version,
      enabled: extensionData.enabled || false,
    };
  },
  
  sortExtensions: (extensions, sortBy = 'name') => {
    if (!Array.isArray(extensions)) return [];
    
    return extensions.sort((a, b) => {
      if (sortBy === 'name') {
        return a.name.localeCompare(b.name);
      }
      if (sortBy === 'type') {
        return a.type.localeCompare(b.type);
      }
      return 0;
    });
  },
};

// Higher-order component for extension management
export const withExtensionCapabilities = (WrappedComponent) => {
  return function WithExtensionCapabilities(props) {
    const getCapabilities = useGetCapabilities();
    const getFullPageExtensions = useGetFullPageExtensions();
    
    const extensionProps = {
      getCapabilities,
      getFullPageExtensions,
      ...props,
    };
    
    return <WrappedComponent {...extensionProps} />;
  };
};

export default ExtensionSandbox;