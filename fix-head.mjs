import { readFileSync, writeFileSync, readdirSync } from 'fs';

const CSS_PATH = '/Users/hsf/.trae-cn/design_libraries/dl_builtin_apple/colors_and_type.css';
const PAGES_DIR = '/Users/hsf/Desktop/AI/XYFamily/xyfamily-admin-design/pages';

const cssContent = readFileSync(CSS_PATH, 'utf-8');
const TAILWIND_CDN = 'https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4.3.1/dist/index.global.js';

function buildHead(title) {
  return `<!DOCTYPE html>
<html lang="zh-CN" class="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title}</title>
    <style id="theme-vars">
${cssContent}
    </style>
    <script src="${TAILWIND_CDN}"></script>
    <style type="text/tailwindcss">
  @theme inline {
    --color-background: var(--background);
    --color-foreground: var(--foreground);
    --color-card: var(--card);
    --color-card-foreground: var(--card-foreground);
    --color-popover: var(--popover);
    --color-popover-foreground: var(--popover-foreground);
    --color-primary: var(--primary);
    --color-primary-foreground: var(--primary-foreground);
    --color-secondary: var(--secondary);
    --color-secondary-foreground: var(--secondary-foreground);
    --color-muted: var(--muted);
    --color-muted-foreground: var(--muted-foreground);
    --color-accent: var(--accent);
    --color-accent-foreground: var(--accent-foreground);
    --color-destructive: var(--destructive);
    --color-destructive-foreground: var(--destructive-foreground);
    --color-success: var(--success);
    --color-success-foreground: var(--success-foreground);
    --color-border: var(--border);
    --color-input: var(--input);
    --color-ring: var(--ring);
    --color-chart-1: var(--chart-1);
    --color-chart-2: var(--chart-2);
    --color-chart-3: var(--chart-3);
    --color-chart-4: var(--chart-4);
    --color-chart-5: var(--chart-5);
    --color-sidebar: var(--sidebar);
    --color-sidebar-foreground: var(--sidebar-foreground);
    --color-sidebar-primary: var(--sidebar-primary);
    --color-sidebar-primary-foreground: var(--sidebar-primary-foreground);
    --color-sidebar-accent: var(--sidebar-accent);
    --color-sidebar-accent-foreground: var(--sidebar-accent-foreground);
    --radius-sm: var(--brand-radius-sm);
    --radius-md: var(--brand-radius-md);
    --radius-lg: var(--brand-radius-lg);
  }
  @layer base {
    body { background: var(--background); color: var(--foreground); }
    td, th { @apply break-words; word-break: break-all; word-break: auto-phrase; }
    th { @apply whitespace-nowrap; }
  }
    </style>
    <style>
      .no-scrollbar::-webkit-scrollbar { display: none; }
      .no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
    </style>
`;
}

const files = readdirSync(PAGES_DIR).filter(f => f.endsWith('.html'));
for (const file of files) {
  const filePath = `${PAGES_DIR}/${file}`;
  let content = readFileSync(filePath, 'utf-8');
  
  const titleMatch = content.match(/<title>(.*?)<\/title>/);
  const title = titleMatch ? titleMatch[1] : file.replace('.html', '');
  
  const bodyStartMatch = content.match(/<\/(?:head|HEAD)>\s*([\s\S]*?)<\/(?:html|HTML)>/);
  if (!bodyStartMatch) {
    console.log(`SKIP ${file}: cannot parse body`);
    continue;
  }
  const bodyContent = bodyStartMatch[1].trim();
  
  const newHtml = buildHead(title) + bodyContent + '\n</html>\n';
  writeFileSync(filePath, newHtml);
  console.log(`FIXED ${file}`);
}
console.log(`Done. Fixed ${files.length} files.`);
