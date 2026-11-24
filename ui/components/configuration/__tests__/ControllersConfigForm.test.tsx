// What the rendered controllers editor is responsible for, asserted against the
// real Sistent controls rather than stubs: every control is tri-state, so a
// field left on Inherit stays absent from the document and the next precedence
// layer applies; the LoadBalancer-only service fields exist only while the
// effective service type is LoadBalancer; and the effective deployment mode
// governs which settings are live, which are inert, and what the form says
// about the ones it cannot apply.
//
// The document semantics under test are the client half of the layering
// described in docs/content/en/project/contributing/contributing-controllers-config.md:
// "merged" carries only explicitly-set fields, so a control that writes a value
// where the user asked to inherit silently pins that layer forever.

import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it } from 'vitest';
import ControllersConfigForm, {
  BUILT_IN_CONTROLLERS_CONFIG,
  type ControllersConfigDoc,
  type ControllersConfigFormProps,
} from '../ControllersConfigForm';
import type { DeploymentMode, DeploymentModeGovernance } from '../deploymentMode';

// --- harness ---------------------------------------------------------------

/**
 * The form is controlled, so the tests drive it through a holder that feeds
 * each onChange back in as the next `value`. `doc()` is then the document the
 * form would submit - which is the thing every assertion here is about.
 */
type Holder = { doc: () => ControllersConfigDoc; changes: () => number };

const renderForm = (
  props: Partial<ControllersConfigFormProps> & { initial?: ControllersConfigDoc } = {},
): Holder => {
  const { initial = {}, ...rest } = props;
  let latest: ControllersConfigDoc = initial;
  let changes = 0;

  const Harness = () => {
    const [value, setValue] = React.useState<ControllersConfigDoc>(initial);
    latest = value;
    return (
      <ControllersConfigForm
        value={value}
        onChange={(next) => {
          changes += 1;
          setValue(next);
        }}
        {...rest}
      />
    );
  };

  render(<Harness />);
  return { doc: () => latest, changes: () => changes };
};

const governance = (
  mode: DeploymentMode,
  scope: 'connection' | 'serverDefault' = 'connection',
): DeploymentModeGovernance => ({
  mode,
  sourceLabel: scope === 'connection' ? "this connection's override" : 'this server-wide default',
  unsaved: false,
  scope,
});

// --- locating controls -----------------------------------------------------
//
// The field labels are Typography, not <label for>, so they cannot be reached
// with getByLabelText. Each label sits in a Box next to its control inside one
// grid item, and every field of a section shares one grid container - which is
// what disambiguates the two "Replicas" fields without depending on DOM order.

const gridItemOf = (label: HTMLElement): HTMLElement =>
  label.parentElement!.parentElement! as HTMLElement;

const sectionGrid = (uniqueLabelInSection: string): HTMLElement =>
  gridItemOf(screen.getByText(uniqueLabelInSection)).parentElement! as HTMLElement;

const meshsyncGrid = () => sectionGrid('MeshSync version');
const brokerGrid = () => sectionGrid('Broker version');

const field = (label: string, scope?: HTMLElement): HTMLElement =>
  gridItemOf(scope ? within(scope).getByText(label) : screen.getByText(label));

const textbox = (label: string, scope?: HTMLElement) =>
  within(field(label, scope)).getByRole('textbox');
const spinbutton = (label: string, scope?: HTMLElement) =>
  within(field(label, scope)).getByRole('spinbutton');
const combobox = (label: string, scope?: HTMLElement) =>
  within(field(label, scope)).getByRole('combobox');

/** The control of a field whatever its kind - text, number or select. */
const control = (label: string, scope?: HTMLElement): HTMLElement => {
  const item = field(label, scope);
  return (item.querySelector('input:not([type="hidden"]), textarea, [role="combobox"]') ??
    within(item).getByRole('combobox')) as HTMLElement;
};

/**
 * Replaces the whole contents of a free-text control in one edit.
 *
 * The collection controls cannot be typed into character by character - see
 * scratchpad/blocker.md - so these tests exercise the parse contract the way a
 * paste does, in a single change, rather than asserting the broken path.
 */
const replaceText = async (input: HTMLElement, text: string) => {
  await userEvent.clear(input);
  if (text === '') return;
  await userEvent.click(input);
  await userEvent.paste(text);
};

/** Picks an option from a Sistent/MUI select by its visible text. */
const choose = async (control: HTMLElement, option: string | RegExp) => {
  await userEvent.click(control);
  const listbox = await screen.findByRole('listbox');
  await userEvent.click(within(listbox).getByRole('option', { name: option }));
};

let user: ReturnType<typeof userEvent.setup>;
beforeEach(() => {
  user = userEvent.setup();
});

// --- tri-state: Inherit is absence, not a value ----------------------------

describe('tri-state inherit / override', () => {
  it('mirrors the server built-in defaults', () => {
    // The doc names this mirror as a thing that must agree with
    // connections.BuiltInControllersConfig(); the form reads every "Inherit
    // (...)" hint out of it, so a drift here mislabels every unset control.
    expect(BUILT_IN_CONTROLLERS_CONFIG).toEqual({
      operator: { deploymentMode: 'embedded' },
      meshsync: {
        replicas: 1,
        redactSecrets: false,
        brokerContentDedup: false,
        debugLogging: false,
      },
      broker: { replicas: 1, service: { type: 'ClusterIP' } },
    });
  });

  it('starts every control on Inherit with an empty document', () => {
    renderForm();
    expect(combobox('Deployment mode')).toHaveTextContent('Inherit (embedded)');
    expect(textbox('Operator version')).toHaveValue('');
    expect(spinbutton('Replicas', meshsyncGrid())).toHaveValue(null);
    expect(combobox('Secret redaction')).toHaveTextContent('Inherit (Disabled)');
    expect(combobox('Service type')).toHaveTextContent('Inherit (ClusterIP)');
    expect(combobox('Watched resources (discovery scope)')).toHaveTextContent('Inherit');
  });

  it('names the inherited value it would fall back to', () => {
    renderForm({
      inheritedLayers: [{ meshsync: { version: 'v0.8.0', outputNamespaces: ['kube-system'] } }],
    });
    expect(textbox('MeshSync version')).toHaveAttribute('placeholder', 'Inherit (v0.8.0)');
    expect(textbox('Output namespaces')).toHaveAttribute('placeholder', 'Inherit (kube-system)');
    // No layer sets it, so there is nothing to name.
    expect(textbox('Output resources')).toHaveAttribute('placeholder', 'Inherit (all)');
  });

  it('round-trips a text field: overriding stores it, emptying releases it', async () => {
    const holder = renderForm();
    await user.type(textbox('MeshSync version'), 'v0.8.42');
    expect(holder.doc()).toEqual({ meshsync: { version: 'v0.8.42' } });

    await user.clear(textbox('MeshSync version'));
    // The parent object is pruned too: an all-inherit section must not survive
    // as an empty husk in the stored document.
    expect(holder.doc()).toEqual({});
  });

  it('round-trips a number field', async () => {
    const holder = renderForm();
    await user.type(spinbutton('Replicas', meshsyncGrid()), '3');
    expect(holder.doc()).toEqual({ meshsync: { replicas: 3 } });

    await user.clear(spinbutton('Replicas', meshsyncGrid()));
    expect(holder.doc()).toEqual({});
  });

  it('round-trips a tri-state boolean, distinguishing Inherit from Disabled', async () => {
    const holder = renderForm();
    await choose(combobox('Secret redaction'), 'Enabled');
    expect(holder.doc()).toEqual({ meshsync: { redactSecrets: true } });

    await choose(combobox('Secret redaction'), 'Disabled');
    // Disabled is an override of false - not the same as inheriting a false.
    expect(holder.doc()).toEqual({ meshsync: { redactSecrets: false } });

    await choose(combobox('Secret redaction'), /^Inherit/);
    expect(holder.doc()).toEqual({});
  });

  it('round-trips an enum select', async () => {
    const holder = renderForm();
    await choose(combobox('Service type'), 'NodePort');
    expect(holder.doc()).toEqual({ broker: { service: { type: 'NodePort' } } });

    await choose(combobox('Service type'), /^Inherit/);
    expect(holder.doc()).toEqual({});
  });

  it('round-trips a list field, trimming and dropping blanks', async () => {
    const holder = renderForm();
    await replaceText(textbox('Output namespaces'), 'kube-system, ,default ');
    expect(holder.doc()).toEqual({ meshsync: { outputNamespaces: ['kube-system', 'default'] } });

    await user.clear(textbox('Output namespaces'));
    expect(holder.doc()).toEqual({});
  });

  it('replaces a collection whole rather than merging it element-wise', async () => {
    // Collections merge atomically across layers, and the control has to behave
    // the same way: a layer that sets one replaces the lower layer's value.
    const holder = renderForm({
      initial: { meshsync: { outputResources: ['pod', 'deployment'] } },
      inheritedLayers: [{ meshsync: { outputResources: ['service'] } }],
    });
    expect(textbox('Output resources')).toHaveValue('pod, deployment');
    await replaceText(textbox('Output resources'), 'configmap');
    expect(holder.doc()).toEqual({ meshsync: { outputResources: ['configmap'] } });
  });

  it('round-trips broker service annotations', async () => {
    const holder = renderForm();
    await replaceText(
      textbox('Service annotations'),
      'service.beta.kubernetes.io/aws-load-balancer-internal=true\nteam=platform',
    );
    expect(holder.doc()).toEqual({
      broker: {
        service: {
          annotations: {
            'service.beta.kubernetes.io/aws-load-balancer-internal': 'true',
            team: 'platform',
          },
        },
      },
    });

    await user.clear(textbox('Service annotations'));
    expect(holder.doc()).toEqual({});
  });

  it('round-trips the watch list through both of its mutually exclusive modes', async () => {
    const holder = renderForm();
    const watchMode = () => combobox('Watched resources (discovery scope)');

    await choose(watchMode(), 'Whitelist (watch only these)');
    expect(holder.doc()).toEqual({ meshsync: { watchList: { whitelist: [] } } });

    await user.click(screen.getByRole('button', { name: 'Add resource' }));
    expect(holder.doc()).toEqual({
      meshsync: {
        watchList: { whitelist: [{ resource: '', events: ['ADDED', 'MODIFIED', 'DELETED'] }] },
      },
    });

    // Switching mode replaces the collection wholesale - a document carrying
    // both a whitelist and a blacklist is rejected by the server.
    await choose(watchMode(), 'Blacklist (default scope minus these)');
    expect(holder.doc()).toEqual({ meshsync: { watchList: { blacklist: [] } } });

    await replaceText(
      field('Watched resources (discovery scope)').querySelector('textarea')!,
      'secrets.v1.\nevents.v1.',
    );
    expect(holder.doc()).toEqual({
      meshsync: { watchList: { blacklist: ['secrets.v1.', 'events.v1.'] } },
    });

    await choose(watchMode(), 'Inherit');
    expect(holder.doc()).toEqual({});
  });

  it('releases one field without disturbing its siblings', async () => {
    const holder = renderForm({
      initial: { meshsync: { version: 'v0.8.42', replicas: 3 } },
    });
    await user.clear(textbox('MeshSync version'));
    expect(holder.doc()).toEqual({ meshsync: { replicas: 3 } });
  });
});

// --- precedence, as the form states it -------------------------------------

describe('source indicators', () => {
  const layered = () =>
    renderForm({
      showSourceIndicators: true,
      inheritLabel: 'Server default',
      initial: { meshsync: { version: 'v0.8.42' } },
      inheritedLayers: [{ meshsync: { replicas: 5 } }, BUILT_IN_CONTROLLERS_CONFIG],
    });

  it('names the layer each field resolves from', () => {
    layered();
    expect(field('MeshSync version')).toHaveTextContent('Override');
    expect(field('Replicas', meshsyncGrid())).toHaveTextContent('Server default');
    expect(field('Broker version')).toHaveTextContent('Built-in default');
  });

  it('omits the chips entirely on the layer that has nothing above it', () => {
    renderForm({ initial: { meshsync: { version: 'v0.8.42' } } });
    expect(field('MeshSync version')).not.toHaveTextContent('Override');
  });
});

// --- LoadBalancer-only fields ----------------------------------------------

describe('LoadBalancer service fields', () => {
  const lbFieldsPresent = () => ({
    class: screen.queryByText('Load balancer class') !== null,
    ranges: screen.queryByText('Load balancer source ranges') !== null,
  });

  it('are absent while the effective service type is not LoadBalancer', () => {
    renderForm({ inheritedLayers: [BUILT_IN_CONTROLLERS_CONFIG] });
    expect(lbFieldsPresent()).toEqual({ class: false, ranges: false });
  });

  it('appear when this layer selects LoadBalancer', async () => {
    renderForm();
    await choose(combobox('Service type'), 'LoadBalancer');
    expect(lbFieldsPresent()).toEqual({ class: true, ranges: true });
  });

  it('appear when a lower layer supplies LoadBalancer and this one inherits', () => {
    renderForm({ inheritedLayers: [{ broker: { service: { type: 'LoadBalancer' } } }] });
    expect(combobox('Service type')).toHaveTextContent('Inherit (LoadBalancer)');
    expect(lbFieldsPresent()).toEqual({ class: true, ranges: true });
  });

  it('clear their stored values when the type moves away from LoadBalancer', async () => {
    const holder = renderForm({
      initial: {
        broker: {
          service: {
            type: 'LoadBalancer',
            loadBalancerClass: 'service.k8s.aws/nlb',
            loadBalancerSourceRanges: ['10.0.0.0/8'],
            externalEndpointOverride: 'broker.example.com:4222',
          },
        },
      },
    });
    await choose(combobox('Service type'), /^ClusterIP/);
    // Left behind they would be invisible in the form and rejected by the
    // server, with nothing on screen to clear.
    expect(holder.doc()).toEqual({
      broker: {
        service: { type: 'ClusterIP', externalEndpointOverride: 'broker.example.com:4222' },
      },
    });
    expect(lbFieldsPresent()).toEqual({ class: false, ranges: false });
  });

  it('keep their values when the type returns to Inherit over a LoadBalancer layer', async () => {
    const holder = renderForm({
      initial: {
        broker: { service: { type: 'LoadBalancer', loadBalancerClass: 'service.k8s.aws/nlb' } },
      },
      inheritedLayers: [{ broker: { service: { type: 'LoadBalancer' } } }],
    });
    await choose(combobox('Service type'), /^Inherit/);
    expect(holder.doc()).toEqual({
      broker: { service: { loadBalancerClass: 'service.k8s.aws/nlb' } },
    });
    expect(lbFieldsPresent()).toEqual({ class: true, ranges: true });
  });

  it('clear when the type returns to Inherit over a non-LoadBalancer layer', async () => {
    const holder = renderForm({
      initial: {
        broker: { service: { type: 'LoadBalancer', loadBalancerSourceRanges: ['10.0.0.0/8'] } },
      },
      inheritedLayers: [BUILT_IN_CONTROLLERS_CONFIG],
    });
    await choose(combobox('Service type'), /^Inherit/);
    expect(holder.doc()).toEqual({});
  });
});

// --- deployment-mode gating ------------------------------------------------

describe('deployment mode gating (per-connection editor)', () => {
  const disabledState = (control: HTMLElement) =>
    control.hasAttribute('disabled') || control.getAttribute('aria-disabled') === 'true';

  it('leaves every setting live in Operator mode', () => {
    renderForm({ deploymentMode: governance('operator') });
    expect(screen.queryByText('Not applied')).toBeNull();
    expect(disabledState(textbox('MeshSync version'))).toBe(false);
    expect(disabledState(textbox('Broker version'))).toBe(false);
    expect(screen.queryByTestId('controllers-config-inert-broker')).toBeNull();
  });

  it('treats every setting as live when no mode governs the editor', () => {
    renderForm();
    expect(screen.queryByTestId('controllers-config-mode-banner')).toBeNull();
    expect(screen.queryByText('Not applied')).toBeNull();
    expect(disabledState(textbox('Broker version'))).toBe(false);
  });

  it('renders the in-cluster settings inert in Embedded mode', () => {
    renderForm({ deploymentMode: governance('embedded') });

    // Every path the doc's mode table calls inert in embedded mode, by the
    // label the user reads it under.
    const inertFields: [string, HTMLElement | undefined][] = [
      ['Operator version', undefined],
      ['MeshSync version', undefined],
      ['Replicas', meshsyncGrid()],
      ['Secret redaction', undefined],
      ['Broker content dedup', undefined],
      ['Debug logging', undefined],
      ['Watched resources (discovery scope)', undefined],
      ['Broker version', undefined],
      ['Replicas', brokerGrid()],
      ['Service type', undefined],
      ['External endpoint override', undefined],
      ['Service annotations', undefined],
    ];
    for (const [label, scope] of inertFields) {
      expect(disabledState(control(label, scope)), `${label} should be inert`).toBe(true);
      expect(field(label, scope)).toHaveTextContent('Not applied');
    }
  });

  it('keeps the settings embedded mode can still apply live', () => {
    renderForm({ deploymentMode: governance('embedded') });

    expect(disabledState(combobox('Deployment mode'))).toBe(false);
    expect(field('Deployment mode')).not.toHaveTextContent('Not applied');
    for (const label of ['Output namespaces', 'Output resources']) {
      expect(disabledState(textbox(label)), `${label} should stay live`).toBe(false);
      expect(field(label)).not.toHaveTextContent('Not applied');
    }
  });

  it('says why in the form body, and marks a fully inert section on its heading', () => {
    renderForm({ deploymentMode: governance('embedded') });

    expect(screen.getByTestId('controllers-config-inert-broker')).toHaveTextContent(
      'Embedded mode installs no Broker on this cluster',
    );
    expect(screen.getByTestId('controllers-config-inert-meshsync')).toHaveTextContent(
      'Only the output filters below take effect',
    );
    // Broker is inert end to end; MeshSync still applies its output filters, so
    // claiming the whole section is dead would be wrong.
    expect(screen.getByTestId('controllers-config-section-broker')).toHaveTextContent(
      'Not applied in Embedded mode',
    );
    expect(screen.getByTestId('controllers-config-section-meshsync')).not.toHaveTextContent(
      'Not applied in Embedded mode',
    );
    expect(screen.getByTestId('controllers-config-mode-banner')).toHaveTextContent(
      'Effective deployment mode: Embedded (in Meshery Server)',
    );
  });

  it('keeps values the mode cannot apply, names them dormant, and clears only those', async () => {
    const holder = renderForm({
      deploymentMode: governance('embedded'),
      initial: {
        broker: { version: '2.10.1', replicas: 3, service: { type: 'NodePort' } },
        meshsync: { version: 'v0.8.42', outputNamespaces: ['kube-system'] },
      },
    });

    // Dormant, not dropped: the user chose them and they go live again in
    // Operator mode, so they stay on screen.
    expect(textbox('Broker version')).toHaveValue('2.10.1');
    const brokerNotice = screen.getByTestId('controllers-config-inert-broker');
    expect(brokerNotice).toHaveTextContent('Clear 3 dormant values');
    expect(screen.getByTestId('controllers-config-inert-meshsync')).toHaveTextContent(
      'Clear 1 dormant value',
    );

    await user.click(within(brokerNotice).getByRole('button', { name: /Clear 3 dormant values/ }));
    expect(holder.doc()).toEqual({
      meshsync: { version: 'v0.8.42', outputNamespaces: ['kube-system'] },
    });
  });

  it('suppresses the secret-redaction warning where the setting does not decide it', () => {
    renderForm({ deploymentMode: governance('embedded') });
    // In embedded mode the Meshery Server process environment decides redaction,
    // so a warning pointing at this control would name the wrong knob.
    expect(screen.queryByText(/Secret redaction is disabled/)).toBeNull();
  });

  it('warns about disabled secret redaction where the setting does decide it', () => {
    renderForm({ deploymentMode: governance('operator') });
    expect(screen.getByText(/Secret redaction is disabled/)).toBeInTheDocument();
  });

  it('honours the disabled prop independently of the mode', () => {
    renderForm({ deploymentMode: governance('operator'), disabled: true });
    expect(disabledState(textbox('Broker version'))).toBe(true);
    expect(screen.queryByText('Not applied')).toBeNull();
  });
});

describe('deployment mode gating (server-wide defaults editor)', () => {
  const disabledState = (control: HTMLElement) =>
    control.hasAttribute('disabled') || control.getAttribute('aria-disabled') === 'true';

  it('annotates rather than disables: a connection may override the mode', () => {
    renderForm({ deploymentMode: governance('embedded', 'serverDefault') });

    expect(disabledState(textbox('Broker version'))).toBe(false);
    expect(disabledState(textbox('MeshSync version'))).toBe(false);
    expect(screen.queryByText('Not applied')).toBeNull();
    expect(screen.queryByTestId('controllers-config-inert-broker')).toBeNull();
    expect(screen.getByTestId('controllers-config-mode-banner')).toHaveTextContent(
      'Default deployment mode: Embedded (in Meshery Server)',
    );
    expect(
      screen.getByText(/Meshery Broker settings reach only connections running/),
    ).toBeVisible();
  });

  it('leaves nothing dormant, because nothing stored on this layer is dead', () => {
    renderForm({
      deploymentMode: governance('embedded', 'serverDefault'),
      initial: { broker: { version: '2.10.1' } },
    });
    expect(screen.queryByRole('button', { name: /dormant/ })).toBeNull();
  });
});
