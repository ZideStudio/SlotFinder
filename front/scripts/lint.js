import concurrently from 'concurrently';

const { result } = concurrently(['npm:lint:*(!fix)'], {
  prefixColors: 'auto',
  maxProcesses: process.env.CI ? 1 : undefined,
});

// oxlint-disable-next-line jest/require-hook
result.then(
  () => process.exit(0),
  () => process.exit(1),
);
