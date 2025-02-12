import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// Mock Nuxt runtime configuration
vi.mock('nuxt/app', () => ({
  useNuxtApp: vi.fn(),
  defineNuxtPlugin: vi.fn(),
  NuxtLink: {
    template: '<a v-bind="$attrs"><slot /></a>',
    props: ['to']
  }
}))

// Configure Vue Test Utils to use the mocked NuxtLink
config.global.components = {
  NuxtLink: {
    template: '<a v-bind="$attrs"><slot /></a>',
    props: ['to']
  }
}
