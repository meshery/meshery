import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import GetKubernetesNodeIcon from './utils';

vi.mock('@/constants/common', () => ({
  FALLBACK_KUBERNETES_IMAGE_PATH: '/fallback.svg',
  KUBERNETES: 'kubernetes',
}));

const normalizeStaticImagePath = vi.fn((src?: string) => src);

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args: any[]) => normalizeStaticImagePath(...args),
}));

const componentIconMock = vi.fn();

vi.mock('@sistent/sistent', () => ({
  componentIcon: (...args: any[]) => componentIconMock(...args),
}));

vi.mock('css/icons.styles', () => ({
  iconXLarge: { width: 32, height: 32 },
}));

describe('GetKubernetesNodeIcon', () => {
  beforeEach(() => {
    normalizeStaticImagePath.mockReset();
    componentIconMock.mockReset();
    componentIconMock.mockReturnValue('/static/img/pod.svg');
    normalizeStaticImagePath.mockImplementation((src) => src);
  });

  it('renders an image with the normalized icon src for the given kind', () => {
    render(<GetKubernetesNodeIcon kind="Pod" />);
    expect(componentIconMock).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'pod', model: 'kubernetes' }),
    );
    const img = screen.getByRole('img', { name: 'Pod' });
    expect(img).toHaveAttribute('src', '/static/img/pod.svg');
  });

  it('passes through a custom model to componentIcon', () => {
    render(<GetKubernetesNodeIcon kind="MyCRD" model="my-model" />);
    expect(componentIconMock).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'mycrd', model: 'my-model' }),
    );
  });

  it('falls back to the kubernetes fallback path when the image errors', () => {
    componentIconMock.mockReturnValue('/static/img/missing.svg');
    render(<GetKubernetesNodeIcon kind="Pod" />);
    const img = screen.getByRole('img', { name: 'Pod' }) as HTMLImageElement;
    fireEvent.error(img);
    expect(img.src).toContain('/fallback.svg');
  });

  it('uses the default size style when none is provided', () => {
    componentIconMock.mockReturnValue('/static/img/pod.svg');
    render(<GetKubernetesNodeIcon kind="Pod" />);
    const img = screen.getByRole('img', { name: 'Pod' });
    expect(img).toHaveStyle({ width: '32px', height: '32px' });
  });

  it('applies the size override style when provided', () => {
    componentIconMock.mockReturnValue('/static/img/pod.svg');
    render(<GetKubernetesNodeIcon kind="Pod" size={{ width: 10, height: 10 }} />);
    const img = screen.getByRole('img', { name: 'Pod' });
    expect(img).toHaveStyle({ width: '10px', height: '10px' });
  });
});
