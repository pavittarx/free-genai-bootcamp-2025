import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import IndexPage from './index.vue'

describe('Index Page', () => {
  it('renders the page title', () => {
    const wrapper = mount(IndexPage)
    expect(wrapper.text()).toContain('Language Portal Dashboard')
  })

  it('has navigation links', () => {
    const wrapper = mount(IndexPage)
    const dashboardLink = wrapper.find('a[href="/dashboard"]')
    const studySessionsLink = wrapper.find('a[href="/study-sessions"]')
    
    expect(dashboardLink.exists()).toBe(true)
    expect(studySessionsLink.exists()).toBe(true)
  })
})
