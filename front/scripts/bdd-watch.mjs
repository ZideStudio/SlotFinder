import { watch } from 'node:fs';
import { resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawn, spawnSync } from 'node:child_process';

const LOG_PREFIX = '[playwright-bdd watch]';
const DEBOUNCE_MS = 200;
const GRACEFUL_SHUTDOWN_MS = 10_000;
const isWindows = process.platform === 'win32';
const npxCommand = isWindows ? 'npx.cmd' : 'npx';
const root = resolve(fileURLToPath(new URL('..', import.meta.url)));

const RELEVANT_FILE = /(^|\/)(tests\/.*\.(feature|steps\.[jt]sx?)|playwright\.config\.(js|ts|mjs|cjs))$/;

const children = new Set();
let activeBddGen = null;
let debounceTimer = null;

const spawnChild = (args) => {
  // Detached on POSIX so the whole process group (npx -> node -> browsers) can be killed at once.
  const child = spawn(npxCommand, args, { cwd: root, stdio: 'inherit', detached: !isWindows });
  children.add(child);
  child.on('exit', () => children.delete(child));
  return child;
};

const killTree = (child, signal = 'SIGTERM') => {
  if (!child.pid || child.exitCode !== null || child.signalCode !== null) {
    return;
  }

  try {
    if (isWindows) {
      spawnSync('taskkill', ['/pid', String(child.pid), '/T', '/F'], { stdio: 'ignore' });
    } else {
      process.kill(-child.pid, signal);
    }
  } catch {
    // Process already gone.
  }
};

const runBddGen = () => {
  if (activeBddGen) {
    killTree(activeBddGen);
  }

  console.log(`\n${LOG_PREFIX} Regenerating BDD tests...`);
  activeBddGen = spawnChild(['bddgen']);
};

const scheduleBddGen = () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(runBddGen, DEBOUNCE_MS);
};

// Recursive watching is supported natively on Linux since Node 20.
const watcher = watch(root, { recursive: true }, (_eventType, filename) => {
  const relativePath = filename?.replace(/\\/g, '/');
  if (relativePath && !relativePath.startsWith('node_modules/') && RELEVANT_FILE.test(relativePath)) {
    scheduleBddGen();
  }
});
watcher.on('error', (error) => console.error(`${LOG_PREFIX} Watcher error:`, error.message));

const waitForExit = (child, timeoutMs) =>
  new Promise((resolveWait) => {
    if (child.exitCode !== null || child.signalCode !== null) {
      return resolveWait();
    }
    const timer = setTimeout(resolveWait, timeoutMs);
    child.once('exit', () => {
      clearTimeout(timer);
      resolveWait();
    });
  });

let isShuttingDown = false;

const shutdown = async (exitCode = 0) => {
  if (isShuttingDown) {
    return;
  }
  isShuttingDown = true;

  watcher.close();
  clearTimeout(debounceTimer);

  // Playwright runs its webServer (port 3004) in its own process group: it must be able to
  // clean it up itself, so ask it to stop gracefully before falling back to a forced kill.
  const remaining = [...children];
  remaining.forEach((child) => killTree(child, 'SIGINT'));
  await Promise.all(remaining.map((child) => waitForExit(child, GRACEFUL_SHUTDOWN_MS)));

  process.exit(exitCode);
};

// SIGHUP is sent when the terminal is killed; the 'exit' hook is a last-resort safety net.
['SIGINT', 'SIGTERM', 'SIGHUP'].forEach((signal) => process.on(signal, () => shutdown(0)));
process.on('exit', () => children.forEach((child) => killTree(child, 'SIGKILL')));

console.log(`${LOG_PREFIX} Watching feature and step files for changes...`);
runBddGen();

const uiProcess = spawnChild(['playwright', 'test', '--ui', '--ui-host=0.0.0.0', '--ui-port=9323']);
uiProcess.on('exit', (code) => shutdown(code ?? 0));
