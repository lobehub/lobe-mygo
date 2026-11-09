// Cloudflare's drop-trailing-slash html_handling redirects /mygo/index.html to /mygo, which
// costs every demo frame a round trip; the Vite dev server only serves the explicit file.
const FRAME_PATH = import.meta.env.DEV ? '/mygo/index.html' : '/mygo';

export const mygoFrameSrc = (demo: string, scrollable = false) =>
  `${FRAME_PATH}?demo=${encodeURIComponent(demo)}${scrollable ? '&scroll=1' : ''}`;
