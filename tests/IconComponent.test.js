import { shallowMount } from '@vue/test-utils';
import IconComponent from '../src/components/IconComponent.vue';
import theme from '../src/styles/theme.css';

jest.mock('../src/styles/theme.css', () => ({
  colors: {
    primary: '#123456'
  }
}));

describe('IconComponent.vue', () => {
  it('should apply the theme primary color to the icon', () => {
    const wrapper = shallowMount(IconComponent, {
      computed: {
        iconColor: () => '#000000'
      }
    });
    expect(wrapper.vm.getThemeColor).toBe('#123456');
    expect(wrapper.attributes('style')).toContain('color: #123456');
  });

  it('should handle missing theme color gracefully', () => {
    jest.mock('../src/styles/theme.css', () => ({
      colors: {}
    }));
    const wrapper = shallowMount(IconComponent, {
      computed: {
        iconColor: () => '#000000'
      }
    });
    expect(wrapper.vm.getThemeColor).toBeUndefined();
  });

  it('should fallback to iconColor if theme color is not defined', () => {
    jest.mock('../src/styles/theme.css', () => ({
      colors: {}
    }));
    const wrapper = shallowMount(IconComponent, {
      computed: {
        iconColor: () => '#000000'
      }
    });
    expect(wrapper.attributes('style')).toContain('color: #000000');
  });

  it('should correctly import theme.css', () => {
    const wrapper = shallowMount(IconComponent);
    expect(wrapper.vm.$options.computed.getThemeColor()).toBe('#123456');
  });
});