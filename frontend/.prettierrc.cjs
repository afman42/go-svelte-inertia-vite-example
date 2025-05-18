module.exports = {
  useTabs: false,
  printWidth: 80,
  tabWidth: 2,
  semi: false,
  trailingComma: 'none',
  singleQuote: true,
  plugins: [require('prettier-plugin-svelte')],
  overrides: [
    {
      files: '**/*.ts',
      options: { parser: 'typescript' }
    },
    {
      files: '**/*.js',
      options: { parser: 'javascript' }
    }
  ]
}
