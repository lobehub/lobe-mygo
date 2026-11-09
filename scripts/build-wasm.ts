import { execFile } from 'node:child_process';
import { createHash } from 'node:crypto';
import { watch } from 'node:fs';
import { copyFile, mkdir, readdir, readFile, rename, rm, writeFile } from 'node:fs/promises';
import { join, relative, sep } from 'node:path';
import { promisify } from 'node:util';

import { createReloadServer } from './reload-server';

const run = promisify(execFile);
const root = join(import.meta.dirname, '..');
const outDir = join(root, 'public', 'mygo');
const prod = process.argv.includes('--prod');
const watching = process.argv.includes('--watch');
const reloadPort = Number(process.env.MYGO_RELOAD_PORT ?? 5199);
const reloadUrl = watching ? `http://localhost:${reloadPort}/events` : null;

const go = (args: string[], env: Record<string, string> = {}) =>
  run('go', args, { cwd: root, env: { ...process.env, ...env }, maxBuffer: 1 << 24 });

async function build(): Promise<string> {
  const started = performance.now();
  await go(['run', './cmd/gen-api', '-out', 'catalog']);
  await mkdir(outDir, { recursive: true });

  const tmp = join(outDir, `.catalog-${process.pid}.wasm`);
  const flags = prod ? ['-tags', 'mygo_noinspector', '-trimpath', '-ldflags=-s -w'] : [];
  await go(['build', ...flags, '-o', tmp, './catalog'], { CGO_ENABLED: '0', GOARCH: 'wasm', GOOS: 'js' });

  const bytes = await readFile(tmp);
  const id = createHash('sha256').update(bytes).digest('hex').slice(0, 12);
  const name = prod ? `catalog.${id}.wasm` : 'catalog.wasm';
  for (const file of await readdir(outDir)) {
    if (/^catalog(\.[0-9a-f]{12})?\.wasm$/.test(file) && file !== name) await rm(join(outDir, file));
  }
  await rename(tmp, join(outDir, name));

  const { stdout: goroot } = await go(['env', 'GOROOT']);
  await copyFile(join(goroot.trim(), 'lib', 'wasm', 'wasm_exec.js'), join(outDir, 'wasm_exec.js'));

  const page = await readFile(join(root, 'web', 'mygo', 'index.html'), 'utf8');
  await writeFile(
    join(outDir, 'index.html'),
    page
      .replace('__WASM__', JSON.stringify(name))
      .replace('__DEV__', JSON.stringify(!prod))
      .replace('__RELOAD__', JSON.stringify(reloadUrl)),
  );
  if (prod) {
    await writeFile(
      join(root, 'public', '_headers'),
      [
        '/mygo/catalog.*.wasm',
        '  Cache-Control: public, max-age=31536000, immutable',
        '/*',
        `  Link: </mygo/${name}>; rel=prefetch; as=fetch; crossorigin`,
        '/mygo',
        '  ! Link',
        '',
      ].join('\n'),
    );
  }
  const mb = (bytes.length / 1024 / 1024).toFixed(1);
  console.log(`[mygo] ${name} ${mb} MB in ${Math.round(performance.now() - started)} ms`);
  return id;
}

const failure = (error: unknown) => {
  const e = error as { message?: string; stderr?: string };
  return (e.stderr || e.message || String(error)).trim();
};

const ignored = (path: string) =>
  ['node_modules', '.ryugi', 'public', '.react-router', 'dist', '.git'].some((dir) => path.split(sep).includes(dir));

async function main() {
  if (!watching) {
    await build();
    return;
  }

  const server = createReloadServer(reloadPort);
  console.log(`[mygo] reload events on ${reloadUrl}`);

  let timer: ReturnType<typeof setTimeout> | undefined;
  let building = false;
  let pending = false;
  const rebuild = async () => {
    if (building) {
      pending = true;
      return;
    }
    building = true;
    server.broadcast({ type: 'building' });
    try {
      server.broadcast({ id: await build(), type: 'built' });
    } catch (error) {
      const message = failure(error);
      console.error(`[mygo] build failed:\n${message}`);
      server.broadcast({ message, type: 'error' });
    } finally {
      building = false;
      if (pending) {
        pending = false;
        void rebuild();
      }
    }
  };

  await rebuild();
  watch(root, { recursive: true }, (_event, file) => {
    if (!file || !file.endsWith('.go') || ignored(relative(root, join(root, file)))) return;
    clearTimeout(timer);
    timer = setTimeout(rebuild, 200);
  });
  console.log('[mygo] watching .go files');
}

main().catch((error) => {
  console.error(`[mygo] ${failure(error)}`);
  process.exit(1);
});
