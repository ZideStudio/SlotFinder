const config = {
  plugins: [
    "stylelint-scss"
  ],
  extends: ['@pplancq/stylelint-config/prettier'],
  rules: {
    'selector-not-notation': null,
    'color-function-alias-notation': null,
    // BEM pattern enforcement using stylelint-scss plugin's selector-class-pattern rule
    // Supports patterns like: .ds-tag, .ds-tag__element, .ds-tag--modifier, .util-flex
    // resolveNestedSelectors: true allows validation of SCSS nested selectors with &
    "scss/selector-class-pattern": [
      "^[a-z]([a-z0-9-]+)?(__[a-z0-9-]+)?(--[a-z0-9-]+)?$",
      {
        "resolveNestedSelectors": true,
        "message": "Expected class selector to follow BEM pattern: .block__element--modifier (e.g., .ds-tag, .ds-tag__icon, .ds-tag--filled, .util-flex)",
      }
    ],
  },
};

export default config;