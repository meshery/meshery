import React from 'react';
import useStyles from '../../assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { FormatStructuredData } from '../DataFormatter';

const KeyValue = ({ property, value }) => (
  <>
    <div
      style={{
        display: 'flex',
        margin: '0.3rem 0',
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
        {property}
      </p>
      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{value}</p>
    </div>
  </>
);

const Title = ({ title }) => (
  <>
    <p
      style={{
        marginTop: '0px',
        fontSize: '19px',
        fontWeight: 'bold',
      }}
    >
      {title}
    </p>
  </>
);

const ModelContents = ({ model }) => (
  <>
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
        <KeyValue property="Version" value={model.version} />
        <KeyValue property="Registrant" value={model.hostname} />
        <KeyValue
          property="Components"
          value={model.components === null ? '0' : model.components.length}
        />
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
        <KeyValue property="Category" value={model.category?.name} />
        <KeyValue property="Duplicates" value={model.duplicates} />
        <KeyValue
          property="Relationships"
          value={model.relationships === null ? '0' : model.relationships.length}
        />
      </div>
    </div>
  </>
);

const ChildTypeTitle = ({ title }) => (
  <p
    style={{
      fontSize: '18px',
      margin: '0.7rem 0',
    }}
  >
    {title}
  </p>
);

const ChildTitle = ({ title }) => (
  <>
    <p
      style={{
        fontSize: '18px',
        fontWeight: 'bold',
        margin: '0.4rem 0',
      }}
    >
      {title}
    </p>
  </>
);

const ComponentContents = ({ components }) => (
  <div>
    {components.map((component, index) => (
      <div
        key={index}
        style={{
          margin: '0.9rem 0',
        }}
      >
        <ChildTitle title={component.displayName} />
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
            <KeyValue property="API Version" value={component.apiVersion} />
          </div>
          <div style={{ display: 'flex', width: '50%' }}>
            <KeyValue property="Sub Category" value={component.kind} />
          </div>
        </div>
      </div>
    ))}
  </div>
);

const RelationshipContents = ({ relationships }) => (
  <div>
    {relationships.map((rela, index) => (
      <div
        key={index}
        style={{
          margin: '0.9rem 0',
        }}
      >
        <ChildTitle title={rela.kind} />
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
            <KeyValue property="API Version" value={rela.apiVersion} />
          </div>
          <div style={{ display: 'flex', width: '50%' }}>
            <KeyValue property="Sub Type" value={rela.subType} />
          </div>
        </div>
      </div>
    ))}
  </div>
);

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
          <Title title={regi.hostname} />
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
              <KeyValue property="Models" value={regi.summary?.models} />
              <KeyValue property="Components" value={regi.summary?.components} />
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <KeyValue property="Relationships" value={regi.summary?.relationships} />
              <KeyValue property="Policies" value={regi.summary?.policies} />
            </div>
          </div>
          {show.model.displayName && <hr style={{ margin: '1rem 0' }} />}
        </div>
      )}
      {(view === MODELS || view === REGISTRANTS) && (
        <>
          {show.model.displayName && (
            <div>
              <Title title={show.model.displayName} />
              <ModelContents model={show.model} />
            </div>
          )}
          {show.components.length !== 0 && (
            <div>
              <hr style={{ margin: '1rem 0' }} />
              <ChildTypeTitle title="Components" />
              <ComponentContents components={show.components} />
            </div>
          )}
          {show.relationships.length !== 0 && (
            <div>
              <hr style={{ marginTop: '1rem 0' }} />
              <ChildTypeTitle title="Relationships" />
              <RelationshipContents relationships={show.relationships} />
            </div>
          )}
        </>
      )}
      {view === COMPONENTS && comp.displayName && (
        <div>
          <ChildTitle title={comp.displayName} />
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
              <KeyValue property="API Version" value={comp.apiVersion} />
              <KeyValue property="Model Name" value={comp.model?.displayName} />
              <KeyValue property="Kind" value={comp.kind} />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <KeyValue property="Registrant" value={comp.displayhostname} />
              <KeyValue property="Duplicates" value={comp.duplicates} />
            </div>
          </div>
        </div>
      )}
      {view === RELATIONSHIPS && rela.kind && (
        <div>
          <ChildTitle title={rela.kind} />
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
              <KeyValue property="API Version" value={rela.apiVersion} />
              <KeyValue property="Model Name" value={rela.model?.displayName} />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
              <KeyValue property="Sub Type" value={rela.subType} />
              <KeyValue property="Registrant" value={rela.displayhostname} />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MeshModelDetails;
