/// <reference types="svelte" />
/// <reference types="vite/client" />

interface Country {
  name: string
  flag: string
}

type EveryFlagCountries = {
  countries: Country[]
}

interface RegisterFormData {
  name: string
  email: string
  password: string
  password_confirmation: string
}

interface HomePageProps {
  user: User
}

interface User {
  id: number
  name: string
  email: string
}

interface UserProfileProps {
  user: User
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

// Form validation errors type
interface FormErrors {
  [key: string]: string[]
}
