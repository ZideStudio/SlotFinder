export default {
  '*.{js,jsx,mjs,ts,tsx}': 'oxlint --react-plugin --import-plugin --jsx-a11y-plugin --fix',
  '*.{ts,tsx}': 'tsc-files --noEmit',
  '*.{scss,css}': 'stylelint --fix',
  '*.!(js|jsx|mjs|ts|tsx|scss|css)': 'prettier --write --ignore-unknown',
  '(package|package-lock).json': 'npm run package:check',
};
