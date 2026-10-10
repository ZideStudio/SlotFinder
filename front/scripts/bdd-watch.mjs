import { existsSync, readdirSync, watch } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawn } from 'node:child_process';

const root = resolve(fileURLToPath(new URL('..', import.meta.url)));
const npxCommand = process.platform === 'win32' ? 'npx.cmd' : 'npx';
const testsDir = join(root, 'tests');
const configNames = ['playwright.config.js', 'playwright.config.ts', 'playwright.config.mjs', 'playwright.config.cjs'];

const isWindows = process.platform === 'win32';
let activeBddGen = null;
let debounceTimer = null;

const terminateProcessTree = (child, signal = 'SIGINT') => {
  if (!child || child.killed || child.exitCode !== null || child.signalCode !== null) {
    return;
  }

  try {
    if (isWindows) {
      spawn('taskkill', ['/pid', String(child.pid), '/T', '/F'], {
        stdio: 'ignore',
      });
      return;
    }

    process.kill(-child.pid, signal);
  } catch {
    if (!child.killed) {
      child.kill(signal);
    }
  }
};

const isRelevantChange = (filePath) => {
  if (!filePath) {
    return false;
  }

  const normalized = filePath.replace(/\\/g, '/');
  return /(^|\/)tests\/.*\.(feature|steps\.[jt]sx?)$/.test(normalized)
    || /(^|\/)playwright\.config\.(js|ts|mjs|cjs)$/.test(normalized);
};

const runBddGen = () => {
  if (activeBddGen && !activeBddGen.killed) {
    activeBddGen.kill('SIGINT');
  }

  console.log('\n[playwright-bdd watch] Regenerating BDD tests...');
  activeBddGen = spawn(npxCommand, ['bddgen'], {
    cwd: root,
    stdio: 'inherit',
    env: process.env,
    detached: !isWindows,
  });
};

const scheduleBddGen = () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(runBddGen, 200);
};

const watchDirectory = (directory, onChange) => {
  try {
    const watcher = watch(directory, { persistent: true }, (_eventType, filename) => {
      if (!filename) {
        return;
      }

      const relativePath = relative(root, resolve(directory, filename)).replace(/\\/g, '/');
      if (isRelevantChange(relativePath)) {
        onChange();
      }
    });

    watcher.on('error', (error) => {
      console.error(`[playwright-bdd watch] Watcher error for ${directory}:`, error.message);
    });
  } catch (error) {
    console.error(`[playwright-bdd watch] Unable to watch ${directory}:`, error.message);
  }
};

const watchDirectoryRecursively = (directory) => {
  if (!existsSync(directory)) {
    return;
  }

  const stack = [directory];
  const seen = new Set();

  while (stack.length > 0) {
    const currentDir = stack.pop();
    if (!currentDir || seen.has(currentDir)) {
      continue;
    }
    seen.add(currentDir);

    watchDirectory(currentDir, scheduleBddGen);

    try {
      const entries = readdirSync(currentDir, { withFileTypes: true });
      for (const entry of entries) {
        if (!entry.isDirectory()) {
          continue;
        }

        stack.push(join(currentDir, entry.name));
      }
    } catch {
      // Ignore transient directory read errors while the tree is changing.
    }
  }
};

const watchConfigFiles = () => {
  for (const configName of configNames) {
    const configPath = join(root, configName);
    if (!existsSync(configPath)) {
      continue;
    }

    watchDirectory(dirname(configPath), scheduleBddGen);
  }
};

watchDirectoryRecursively(testsDir);
watchConfigFiles();

console.log('[playwright-bdd watch] Watching feature and step files for changes...');
runBddGen();

const uiProcess = spawn(npxCommand, ['playwright', 'test', '--ui', '--ui-host=0.0.0.0', '--ui-port=9323'], {
  cwd: root,
  stdio: 'inherit',
  env: process.env,
  detached: !isWindows,
});

uiProcess.on('exit', (code) => {
  process.exit(code ?? 0);
});

const shutdown = () => {
  terminateProcessTree(activeBddGen);
  terminateProcessTree(uiProcess);

  setTimeout(() => {
    process.exit(0);
  }, 100);
};

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);
