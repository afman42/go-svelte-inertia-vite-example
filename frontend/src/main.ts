import { createInertiaApp, type ResolvedComponent } from '@inertiajs/svelte'
import { mount } from 'svelte'
import './app.css'

createInertiaApp({
  resolve: (name: string) => {
    const pages = import.meta.glob('./pages/**/*.svelte', { eager: true })
    return pages[`./pages/${name}.svelte`] as
      | ResolvedComponent
      | Promise<ResolvedComponent>
  },
  setup({ el, App, props }) {
    if (el == null) return

    // Global error handling
    window.addEventListener('error', (event) => {
      console.error('Global error:', event.error)
    })

    window.addEventListener('unhandledrejection', (event) => {
      console.error('Unhandled promise rejection:', event.reason)
    })

    mount(App, { target: el, props })
  }
})
