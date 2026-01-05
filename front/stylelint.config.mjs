const config = {
  extends: ['@pplancq/stylelint-config/prettier'],
  rules: {
    'selector-not-notation': null,
    'color-function-alias-notation': null,
    'selector-class-pattern': '^[a-z]([a-z0-9]*(-{1,2}|_{1,2})?)*[a-z0-9]$'
  },
};

export default config;
