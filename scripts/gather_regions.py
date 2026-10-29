"""Place every gathering resource on a map (DD-317): python scripts/gather_regions.py

A resource of gathering level L (1..100) feeds items used at character
level (L-1)*10+1 (the recipe rule of DD-316), so it grows in the open
region (zone or wilderness) with the highest level not above that. The
region goes into the resource's `gather` spec in content/gathering.yaml;
rerun after adding resources or regions.
"""
import io
import os
import re

import yaml

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'content')


def main():
    regions = yaml.safe_load(io.open(os.path.join(ROOT, 'regions.yaml'), encoding='utf-8'))['regions']
    open_regions = sorted((r.get('min_level', 1), r['code']) for r in regions
                          if r.get('kind') in ('zone', 'wilderness'))

    def region_for(level):
        item_level = (level - 1) * 10 + 1
        best = open_regions[0][1]
        for lv, code in open_regions:
            if lv <= item_level:
                best = code
        return best

    path = os.path.join(ROOT, 'gathering.yaml')
    raw = io.open(path, encoding='utf-8', newline='').read()
    crlf = '\r\n' in raw
    s = raw.replace('\r\n', '\n')
    count = 0

    def place(m):
        nonlocal count
        count += 1
        kind, level, tier = m.group(1), int(m.group(2)), m.group(3)
        return 'gather: { kind: %s, level: %d, tier: %s, region: %s }' % (kind, level, tier, region_for(level))

    s = re.sub(r'gather: \{ kind: (\w+), level: (\d+), tier: (\d+)(?:, region: \w+)? \}', place, s)
    io.open(path, 'w', encoding='utf-8', newline='').write(s.replace('\n', '\r\n') if crlf else s)
    print('placed %d resources' % count)


if __name__ == '__main__':
    main()
