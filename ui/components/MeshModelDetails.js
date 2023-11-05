import React from 'react';
import useStyles from '../assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../constants/navigator';

const MeshModelDetails = ({ view, show, rela, regi, comp }) => {
  const StyleClass = useStyles();

  return (
    <div
      className={
        (view === MODELS && !show.model.displayName) ||
        (view === COMPONENTS && !comp.displayName) ||
        (view === RELATIONSHIPS && !rela.kind) ||
        (view === REGISTRANTS && !regi.hostname)
          ? StyleClass.emptyDetailsContainer
          : StyleClass.detailsContainer
      }
    >
      {((view === MODELS && !show.model.displayName) ||
        (view === COMPONENTS && !comp.displayName) ||
        (view === RELATIONSHIPS && !rela.kind) ||
        (view === REGISTRANTS && !regi.hostname)) && (
        <p style={{ color: '#969696' }}>No {view} selected</p>
      )}
      {view === REGISTRANTS && regi.hostname && (
        <div>
          <p
            style={{
              fontSize: '20px',
              marginTop: '0',
              fontWeight: 'bold',
            }}
          >
            {regi.hostname}
          </p>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
            }}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                width: '50%',
                paddingRight: '1rem',
              }}
            >
              <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                Models
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{regi.summary?.models}</p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Components
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                {regi.summary?.components}
              </p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                Relationships
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                {regi.summary?.relationships}
              </p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Policies
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                {regi.summary?.policies}
              </p>
            </div>
          </div>
          {show.model.displayName && <hr />}
        </div>
      )}
      {(view === MODELS || view === REGISTRANTS) && (
        <>
          {show.model.displayName && (
            <div>
              <p
                style={{
                  marginTop: '0px',
                  fontSize: '18px',
                  fontWeight: 'bold',
                }}
              >
                {show.model.displayName}
              </p>
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'row',
                }}
              >
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    width: '50%',
                    paddingRight: '1rem',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontSize: '14px',
                        fontWeight: '600',
                      }}
                    >
                      Version
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.version}
                    </p>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontSize: '14px',
                        fontWeight: '600',
                      }}
                    >
                      Registrant
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.hostname}
                    </p>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontSize: '14px',
                        fontWeight: '600',
                      }}
                    >
                      Components
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.components === null ? '0' : show.model.components.length}
                    </p>
                  </div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontSize: '14px',
                        fontWeight: '600',
                      }}
                    >
                      Category
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.category?.name}
                    </p>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontSize: '14px',
                        fontWeight: '600',
                      }}
                    >
                      Duplicates
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.duplicates}
                    </p>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                    }}
                  >
                    <p
                      style={{
                        padding: '0',
                        margin: '0 0.5rem 0 0',
                        fontWeight: '600',
                        fontSize: '14px',
                      }}
                    >
                      Relationships
                    </p>
                    <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                      {show.model.relationships === null ? '0' : show.model.relationships.length}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )}
          {show.components.length !== 0 && (
            <div>
              <hr style={{ margin: '1rem 0' }} />
              <p
                style={{
                  fontSize: '18px',
                  margin: '0.7rem 0',
                }}
              >
                Components
              </p>
              {show.components.map((component, index) => (
                <div
                  key={index}
                  style={{
                    margin: '0.9rem 0',
                  }}
                >
                  <p
                    style={{
                      fontSize: '18px',
                      fontWeight: 'bold',
                      margin: '0.4rem 0',
                    }}
                  >
                    {component.displayName}
                  </p>
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'row',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        width: '50%',
                        paddingRight: '1rem',
                      }}
                    >
                      <p
                        style={{
                          padding: '0',
                          margin: '0 0.5rem 0 0',
                          fontSize: '14px',
                          fontWeight: '600',
                        }}
                      >
                        API Version
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {component.apiVersion}
                      </p>
                    </div>
                    <div style={{ display: 'flex', width: '50%' }}>
                      <p
                        style={{
                          padding: '0',
                          margin: '0 0.5rem 0 0',
                          fontSize: '14px',
                          fontWeight: '600',
                        }}
                      >
                        Sub Category
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {component.kind}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
          {show.relationships.length !== 0 && (
            <div>
              <hr style={{ marginTop: '1rem' }} />
              <p
                style={{
                  fontSize: '18px',
                  margin: '0.7rem 0',
                }}
              >
                Relationships
              </p>
              {show.relationships.map((rela, index) => (
                <div key={index}>
                  <p
                    style={{
                      fontSize: '20px',
                      margin: '0',
                      fontWeight: 'bold',
                    }}
                  >
                    {rela.kind}
                  </p>
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'row',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        width: '50%',
                        paddingRight: '1rem',
                      }}
                    >
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        API Version
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {rela.apiVersion}
                      </p>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Sub Type
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{rela.subType}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}
      {view === COMPONENTS && comp.displayName && (
        <div>
          <p
            style={{
              fontSize: '20px',
              fontWeight: 'bold',
              marginTop: '0',
            }}
          >
            {comp.displayName}
          </p>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
            }}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                width: '50%',
                paddingRight: '1rem',
              }}
            >
              <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                API Version
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.apiVersion}</p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Model Name
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                {comp.model?.displayName}
              </p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Kind
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.kind}</p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Registrant
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.displayhostname}</p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Duplicates
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.duplicates}</p>
            </div>
          </div>
        </div>
      )}
      {view === RELATIONSHIPS && rela.kind && (
        <div>
          <p
            style={{
              fontSize: '20px',
              marginTop: '0',
              fontWeight: 'bold',
            }}
          >
            {rela.kind}
          </p>
          <p style={{ fontWeight: '600', margin: '0' }}>Description</p>
          <p style={{ margin: '0', fontSize: '14px' }}>{rela.metadata?.description}</p>
          <div
            style={{
              display: 'flex',
              marginTop: '12px',
              flexDirection: 'row',
            }}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                width: '50%',
                paddingRight: '1rem',
              }}
            >
              <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                API Version
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{rela.apiVersion}</p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontWeight: '600',
                  fontSize: '16px',
                }}
              >
                Model Name
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                {rela.model?.displayName}
              </p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                Sub Type
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{rela.subType}</p>
              <p
                style={{
                  padding: '0',
                  margin: '0',
                  marginTop: '12px',
                  fontSize: '16px',
                  fontWeight: '600',
                }}
              >
                Registrant
              </p>
              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{rela.displayhostname}</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MeshModelDetails;
