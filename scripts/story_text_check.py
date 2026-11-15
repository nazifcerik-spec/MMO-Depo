"""Lists scene text keys missing or unused per language (DD-367):
python scripts/story_text_check.py [lang ...]  (default: en tr)."""
import glob
import io
import os
import sys

import yaml

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'content')


def keys():
    with io.open(os.path.join(ROOT, 'story.yaml'), encoding='utf-8') as f:
        book = yaml.safe_load(f)
    out = []
    for sc in book['scenes']:
        c = sc['code']
        out += [c, c + '.body'] + ([c + '.after'] if sc.get('after') else [])
        out += [c + '.x.' + e['key'] for e in sc.get('extras', [])]
        for o in sc.get('options', []):
            out += [c + '.' + o['code'], c + '.' + o['code'] + '.result']
    return out


def text(lang):
    m = {}
    for f in sorted(glob.glob(os.path.join(ROOT, 'story', lang, '*.yaml'))):
        with io.open(f, encoding='utf-8') as fh:
            m.update(yaml.safe_load(fh) or {})
    return m


want = keys()
for lang in sys.argv[1:] or ['en', 'tr']:
    have = text(lang)
    missing = [k for k in want if k not in have]
    unused = sorted(set(have) - set(want))
    print('%s: %d keys, %d missing, %d unused' % (lang, len(have), len(missing), len(unused)))
    for k in missing:
        print('  missing', k)
    for k in unused:
        print('  unused', k)
