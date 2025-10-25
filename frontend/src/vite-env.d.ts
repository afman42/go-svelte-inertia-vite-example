/// <reference types="svelte" />
/// <reference types="vite/client" />

interface Country {
  name: string
  flag: string
}

type EveryFlagCountries = {
  countries: Country[]
}

interface HomePageProps {
  user: string
}

interface NewCountryForm {
  name: string
  code: string
}

interface InertiaPageProps {
  url: string
  props: Record<string, any>
  component: string
}
