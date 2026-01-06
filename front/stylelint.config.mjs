const config = {
  plugins: [
    "stylelint-selector-bem-pattern"
  ],
  extends: ['@pplancq/stylelint-config/prettier'],
  rules: {
    "selector-class-pattern": null,
    'selector-not-notation': null,
    'color-function-alias-notation': null,
    "plugin/selector-bem-pattern": {
      "componentName": "[a-z0-9-]+", 
      "componentSelectors": {
        "initial": "^\\.{componentName}(?:__[a-z0-9-]+)?(?:--[a-z0-9-]+)?$"
      },
      "utilitySelectors": "^\\.util-[a-z]+$"
    },
  },
};

export default config;