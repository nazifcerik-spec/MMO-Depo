"""Shared helpers for the content generators: merge names into the
content/i18n/<lang>.yaml files, keeping each "# --- kind ---" section sorted
(a new kind gets its own section at the end)."""
import io
import os
import re

I18N = os.path.join(os.path.dirname(__file__), '..', 'content', 'i18n')
FILES = {'en': 'en', 'tr': 'tr', 'zh': 'zh-CN', 'es': 'es'}
KEY = re.compile(r'^([\w.\-]+):\s*"')


def _quote(v):
    return '"' + v.replace('\\', '\\\\').replace('"', '\\"') + '"'


def merge(names):
    """names: {lang: {key: value}} with lang in en/tr/zh/es."""
    for lang, mapping in names.items():
        path = os.path.join(I18N, FILES[lang] + '.yaml')
        with io.open(path, encoding='utf-8', newline='') as f:
            raw = f.read()
        crlf = '\r\n' in raw
        lines = raw.replace('\r\n', '\n').rstrip('\n').split('\n')
        for key, value in sorted(mapping.items()):
            line = '%s: %s' % (key, _quote(value))
            idx = {KEY.match(l).group(1): i for i, l in enumerate(lines) if KEY.match(l)}
            if key in idx:
                lines[idx[key]] = line
                continue
            kind = key.split('.', 1)[0]
            same = sorted(k for k in idx if k.split('.', 1)[0] == kind)
            if same:
                after = [k for k in same if k > key]
                lines.insert(idx[after[0]] if after else idx[same[-1]] + 1, line)
            else:
                lines += ['', '# --- %s ---' % kind, line]
        out = '\n'.join(lines) + '\n'
        with io.open(path, 'w', encoding='utf-8', newline='') as f:
            f.write(out.replace('\n', '\r\n') if crlf else out)
