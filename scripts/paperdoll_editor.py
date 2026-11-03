"""Paper doll position editor (DD-347), development only.

    python scripts/paperdoll_editor.py [port]     # http://127.0.0.1:5175

Pick equipment for each slot, move it on the figure and save: the rules go to
client/assets/paperdoll.json, no code changes. A rule can be saved for the one
item, its type (<slot>.<subclass>, e.g. all axes) or its slot. The file ships with the next build.
"""
import io
import json
import mimetypes
import os
import re
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

import yaml

ROOT = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..'))
ASSETS = os.path.join(ROOT, 'client', 'assets')
METADATA = os.path.join(ASSETS, 'paperdoll.json')
SLOTS = ['mainhand', 'offhand', 'helm', 'chest', 'gloves', 'legs', 'boots', 'necklace', 'ring', 'trinket']
LAYERS = ('front', 'back', 'none')
ITEMS = []  # read once: the content files are large
FIELDS = {'avatarAsset': str, 'anchorPoint': str, 'pivot': list, 'offsetX': (int, float), 'offsetY': (int, float),
          'scale': (int, float), 'rotation': (int, float), 'zIndex': int, 'renderingLayer': str,
          'flipX': bool, 'hidesHair': bool, 'mirrorAt': str}


def items():
    if ITEMS:
        return ITEMS
    with io.open(os.path.join(ROOT, 'content', 'i18n', 'tr.yaml'), encoding='utf-8') as f:
        names = yaml.safe_load(f)
    out = []
    for name in ('items.yaml', 'gear.yaml', 'sets.yaml'):
        with io.open(os.path.join(ROOT, 'content', name), encoding='utf-8') as f:
            for it in (yaml.safe_load(f) or {}).get('items') or []:
                if it.get('equip_slot'):
                    out.append({'id': it['code'], 'slot': it['equip_slot'], 'subclass': it.get('item_subclass', ''),
                                'level': it.get('min_level', 1), 'name': names.get('item.' + it['code'], it['code'])})
    out.sort(key=lambda i: (i['slot'], i['subclass'], i['level']))
    ITEMS.extend(out)
    return ITEMS


def figures():
    return sorted(f[:-5] for f in os.listdir(os.path.join(ASSETS, 'bodies')) if f.endswith('.webp'))


def validate(meta):
    """Refuses anything the client could not read; returns the error or None."""
    if not isinstance(meta, dict):
        return 'not an object'
    anchors = meta.get('anchors')
    if not isinstance(anchors, dict) or not anchors:
        return 'anchors missing'
    def point(p):
        return isinstance(p, list) and len(p) == 2 and all(isinstance(v, (int, float)) for v in p)
    for name, p in anchors.items():
        if not point(p):
            return 'anchor %s is not [x, y]' % name
    for scope in ('slots', 'subclasses', 'items'):
        group = meta.get(scope)
        if not isinstance(group, dict):
            return scope + ' missing'
        for key, rule in group.items():
            if not re.fullmatch(r'[a-z0-9_.]+', key) or not isinstance(rule, dict):
                return 'bad rule %s' % key
            for field, value in rule.items():
                kind = FIELDS.get(field)
                if kind is None or not isinstance(value, kind) or (kind is int and isinstance(value, bool)):
                    return '%s.%s is invalid' % (key, field)
            if 'pivot' in rule and not point(rule['pivot']):
                return '%s.pivot is not [x, y]' % key
            if rule.get('renderingLayer', 'front') not in LAYERS:
                return '%s.renderingLayer is invalid' % key
            for field in ('anchorPoint', 'mirrorAt'):
                if field in rule and rule[field] not in anchors:
                    return '%s.%s is unknown' % (key, field)
    return None


class Handler(BaseHTTPRequestHandler):
    def send(self, code, body, kind='application/json; charset=utf-8'):
        self.send_response(code)
        self.send_header('Content-Type', kind)
        self.send_header('Content-Length', str(len(body)))
        self.send_header('Cache-Control', 'no-store')
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        path = self.path.split('?')[0]
        if path == '/':
            self.send(200, PAGE.encode('utf-8'), 'text/html; charset=utf-8')
        elif path == '/api/state':
            with io.open(METADATA, encoding='utf-8') as f:
                meta = json.load(f)
            self.send(200, json.dumps({'metadata': meta, 'items': items(), 'figures': figures(), 'slots': SLOTS}).encode('utf-8'))
        elif path.startswith('/assets/'):
            target = os.path.normpath(os.path.join(ASSETS, path[len('/assets/'):]))
            if not target.startswith(ASSETS + os.sep) or not os.path.isfile(target):
                self.send(404, b'{}')
                return
            with open(target, 'rb') as f:
                self.send(200, f.read(), mimetypes.guess_type(target)[0] or 'application/octet-stream')
        else:
            self.send(404, b'{}')

    def do_POST(self):
        if self.path != '/api/metadata':
            self.send(404, b'{}')
            return
        try:
            meta = json.loads(self.rfile.read(int(self.headers.get('Content-Length', 0))))
        except ValueError:
            self.send(400, json.dumps({'error': 'invalid json'}).encode('utf-8'))
            return
        err = validate(meta)
        if err:
            self.send(400, json.dumps({'error': err}).encode('utf-8'))
            return
        tmp = METADATA + '.tmp'
        with io.open(tmp, 'w', encoding='utf-8', newline='\n') as f:
            json.dump(meta, f, ensure_ascii=False, indent=2)
            f.write('\n')
        os.replace(tmp, METADATA)
        self.send(200, b'{"ok":true}')

    def log_message(self, *args):
        pass


PAGE = r'''<!doctype html><html lang="tr"><head><meta charset="utf-8"><title>Paper Doll Editörü</title><style>
body{margin:0;background:#101114;color:#e8eaee;font:14px system-ui;display:flex;gap:20px;padding:16px;flex-wrap:wrap}
canvas{background:#1a1c22;border:1px solid #3a3f4c;border-radius:8px;cursor:move;touch-action:none}
.panel{width:430px}h3{margin:14px 0 6px;font-size:13px;color:#e4a756;letter-spacing:.5px}
.row{display:flex;align-items:center;gap:8px;margin:5px 0}.row label{width:82px;color:#a8acb8;flex:none}
.row input[type=range]{flex:1}.row input[type=number]{width:70px}
select,input,button{background:#242731;color:#e8eaee;border:1px solid #3a3f4c;border-radius:5px;padding:4px 6px;font:inherit}
select{flex:1;min-width:0}button{cursor:pointer;padding:6px 12px}button.main{background:#e4a756;color:#111;border:0;font-weight:600}
.slot{display:flex;gap:6px;align-items:center;margin:3px 0}.slot input[type=radio]{flex:none}.slot span{width:70px;flex:none;color:#a8acb8}
#status{margin-left:8px;color:#5fb973}#status.err{color:#e85b62}.muted{color:#6e727e;font-size:12px}
</style></head><body>
<div><canvas id="c" width="400" height="640"></canvas>
<div class="row"><label>Figür</label><select id="figure"></select></div>
<div class="row"><label>Yüz katmanı</label><input type="checkbox" id="showHead" checked><span class="muted">yüz 1, saç 1</span>
<label style="width:auto;margin-left:12px">Çapa</label><input type="checkbox" id="showAnchors" checked></div></div>
<div class="panel"><h3>EKİPMAN</h3><div id="slots"></div>
<h3>DÜZENLENEN: <span id="target">-</span></h3>
<div class="row"><label>Kapsam</label><select id="scope"><option value="items">Yalnız bu eşya</option><option value="subclasses">Bu tür (yuva + alt sınıf)</option><option value="slots">Bu yuva</option></select></div>
<div class="row"><label>Çapa noktası</label><select id="anchorPoint"></select></div>
<div class="row"><label>Eşi</label><select id="mirrorAt"><option value="">Yok</option></select></div><div class="row"><label>Katman</label><select id="renderingLayer"><option value="front">Önde</option><option value="back">Arkada</option><option value="none">Çizilmez</option></select></div>
<div id="numbers"></div>
<div class="row"><label>Yatay çevir</label><input type="checkbox" id="flipX"><label style="margin-left:16px">Saçı gizle</label><input type="checkbox" id="hidesHair"></div>
<div class="row"><label>Görsel</label><input id="avatarAsset" style="flex:1"></div>
<div class="row"><button class="main" id="save">Kaydet</button><button id="reset">Bu kapsamı sıfırla</button><span id="status"></span></div>
<p class="muted">Tuvalde sürükleyerek konumu değiştir. Kural önceliği: eşya, sonra tür, sonra yuva. Kayıt <code>client/assets/paperdoll.json</code> dosyasına gider; oyuna sonraki derlemeyle girer.</p>
</div><script>
const $=id=>document.getElementById(id),W=400,H=640,ctx=$('c').getContext('2d');
let meta,allItems,slots,active='mainhand';const equipped={},images={};
const NUMS=[['offsetX','X kayma',-.5,.5,.005,0],['offsetY','Y kayma',-.5,.5,.005,0],['scale','Ölçek',.05,1.5,.01,.5],['rotation','Açı',-180,180,1,0],['pivotX','Tutma X',0,1,.01,.5],['pivotY','Tutma Y',0,1,.01,.5],['zIndex','Z sırası',-50,100,1,0]];
function img(src){if(!images[src]){const i=new Image();i.onload=draw;i.onerror=()=>{i.failed=true;draw()};i.src='/'+src;images[src]=i}return images[src]}
const ready=i=>i.complete&&i.naturalWidth>0;
function item(id){return allItems.find(i=>i.id===id)}
function rule(it){return {...(meta.slots[it.slot]||{}),...(meta.subclasses[it.slot+'.'+it.subclass]||{}),...(meta.items[it.id]||{})}}
function figureParts(){const m=$('figure').value.match(/^(.+)_(male|female)_(.+)$/);return m?{race:m[1],body:m[2]}:{race:'',body:''}}
function anchor(name){const f=figureParts(),o=meta.anchorOverrides||{};return (o[f.race+'_'+f.body]||{})[name]||(o[f.body]||{})[name]||meta.anchors[name]}
function tinted(src,color,fade){const i=img(src);if(!ready(i))return null;const c=document.createElement('canvas');c.width=c.height=256;const g=c.getContext('2d');g.drawImage(i,0,0,256,256);
if(color){g.globalCompositeOperation='multiply';g.fillStyle=color;g.fillRect(0,0,256,256);g.globalCompositeOperation='destination-in';g.drawImage(i,0,0,256,256)}
if(fade){g.globalCompositeOperation='destination-in';const l=g.createLinearGradient(0,0,0,256);l.addColorStop(0,'#fff');l.addColorStop(fade[0],'#fff');l.addColorStop(fade[1],'rgba(0,0,0,0)');g.fillStyle=l;g.fillRect(0,0,256,256)}return c}
function drawHead(hideHair){if(!$('showHead').checked)return;const b=figureParts().body;if(!b)return;const s=W*.34,x=.505*W-s/2,y=.122*H-s/2,p='assets/heads/'+b;
[[p+'_face1_base.webp',null,[.76,.84]],[p+'_face1_brows.webp','#784b2d'],[p+'_face1_eyes.webp','#4682dc'],...(hideHair?[]:[[p+'_hair1.webp','#784b2d',[.9,1]]])].forEach(([src,col,f])=>{const c=tinted(src,col,f);if(c)ctx.drawImage(c,x,y,s,s)})}
function layers(){const out=[];for(const s of slots){const id=equipped[s];if(!id)continue;const it=item(id),r=rule(it);if(!Object.keys(r).length||r.renderingLayer==='none'||!r.avatarAsset)continue;const a=anchor(r.anchorPoint);if(!a)continue;out.push({it,r,a,src:r.avatarAsset.replace('{id}',id)})}return out.sort((p,q)=>(p.r.zIndex||0)-(q.r.zIndex||0))}
function drawLayer(l){const i=img(l.src);if(!ready(i))return;const size=(l.r.scale??.5)*W,pv=l.r.pivot||[.5,.5];const one=(a,side)=>{const px=side<0?1-pv[0]:pv[0];ctx.save();ctx.translate((a[0]+side*(l.r.offsetX||0))*W,(a[1]+(l.r.offsetY||0))*H);ctx.rotate(side*(l.r.rotation||0)*Math.PI/180);ctx.translate(-px*size,-pv[1]*size);if(!!l.r.flipX!==(side<0)){ctx.translate(size,0);ctx.scale(-1,1)}ctx.drawImage(i,0,0,size,size);ctx.restore()};one(l.a,1);if(l.r.mirrorAt&&anchor(l.r.mirrorAt))one(anchor(l.r.mirrorAt),-1)}
function draw(){ctx.clearRect(0,0,W,H);const all=layers(),hide=all.some(l=>l.r.hidesHair);all.filter(l=>l.r.renderingLayer==='back').forEach(drawLayer);
const f=img('assets/bodies/'+$('figure').value+'.webp');if(ready(f))ctx.drawImage(f,0,0,W,H);const front=all.filter(l=>l.r.renderingLayer!=='back');
front.filter(l=>(l.r.zIndex||0)<10).forEach(drawLayer);drawHead(hide);front.filter(l=>(l.r.zIndex||0)>=10).forEach(drawLayer);
if($('showAnchors').checked){ctx.font='10px system-ui';for(const[n,p]of Object.entries(meta.anchors)){const a=anchor(n);ctx.fillStyle=n===currentRule().anchorPoint?'#e4a756':'#5da9e9';ctx.beginPath();ctx.arc(a[0]*W,a[1]*H,3,0,7);ctx.fill();ctx.fillText(n,a[0]*W+5,a[1]*H-4)}}}
function currentItem(){return equipped[active]?item(equipped[active]):null}
function currentRule(){const it=currentItem();return it?rule(it):{...(meta.slots[active]||{})}}
function scopeKey(){const it=currentItem(),s=$('scope').value;return s==='items'?it&&it.id:s==='subclasses'?it&&it.slot+'.'+it.subclass:active}
function setField(k,v){const key=scopeKey();if(!key)return;const g=meta[$('scope').value];g[key]=g[key]||{};g[key][k]=v;refresh()}
function setNum(k,v){const r=currentRule();if(k==='pivotX')setField('pivot',[v,(r.pivot||[.5,.5])[1]]);else if(k==='pivotY')setField('pivot',[(r.pivot||[.5,.5])[0],v]);else setField(k,v)}
function refresh(){const r=currentRule(),it=currentItem();$('target').textContent=it?it.name+' ('+it.id+')':active;$('anchorPoint').value=r.anchorPoint||'';$('mirrorAt').value=r.mirrorAt||'';$('renderingLayer').value=r.renderingLayer||'front';$('flipX').checked=!!r.flipX;$('hidesHair').checked=!!r.hidesHair;$('avatarAsset').value=r.avatarAsset||'';
for(const[k,,,,,d]of NUMS){const v=k==='pivotX'?(r.pivot||[.5,.5])[0]:k==='pivotY'?(r.pivot||[.5,.5])[1]:(r[k]??d);$('r_'+k).value=v;$('n_'+k).value=v}draw()}
async function start(){const s=await(await fetch('/api/state')).json();meta=s.metadata;allItems=s.items;slots=s.slots;meta.subclasses=meta.subclasses||{};meta.items=meta.items||{};
$('figure').innerHTML=s.figures.map(f=>`<option>${f}</option>`).join('');$('figure').value=s.figures.includes('human_male_plate3')?'human_male_plate3':s.figures[0];
$('anchorPoint').innerHTML=Object.keys(meta.anchors).map(a=>`<option>${a}</option>`).join('');$('mirrorAt').innerHTML+=Object.keys(meta.anchors).map(a=>`<option>${a}</option>`).join('');
$('slots').innerHTML=slots.map(sl=>`<div class="slot"><input type="radio" name="active" value="${sl}" ${sl===active?'checked':''}><span>${sl}</span><select data-slot="${sl}"><option value="">-</option>${allItems.filter(i=>i.slot===sl).map(i=>`<option value="${i.id}">${i.subclass} · Sv ${i.level} · ${i.name}</option>`).join('')}</select></div>`).join('');
document.querySelectorAll('[name=active]').forEach(r=>r.onchange=()=>{active=r.value;refresh()});document.querySelectorAll('[data-slot]').forEach(c=>c.onchange=()=>{equipped[c.dataset.slot]=c.value;active=c.dataset.slot;document.querySelector(`[name=active][value=${active}]`).checked=true;refresh()});
$('numbers').innerHTML=NUMS.map(([k,l,a,b,st])=>`<div class="row"><label>${l}</label><input type="range" id="r_${k}" min="${a}" max="${b}" step="${st}"><input type="number" id="n_${k}" min="${a}" max="${b}" step="${st}"></div>`).join('');
for(const[k]of NUMS){const parse=v=>k==='zIndex'?parseInt(v,10)||0:parseFloat(v)||0;$('r_'+k).oninput=e=>setNum(k,parse(e.target.value));$('n_'+k).onchange=e=>setNum(k,parse(e.target.value))}
$('anchorPoint').onchange=e=>setField('anchorPoint',e.target.value);$('mirrorAt').onchange=e=>{const k=scopeKey();if(!k)return;const g=meta[$('scope').value];g[k]=g[k]||{};if(e.target.value)g[k].mirrorAt=e.target.value;else delete g[k].mirrorAt;refresh()};$('renderingLayer').onchange=e=>setField('renderingLayer',e.target.value);$('flipX').onchange=e=>setField('flipX',e.target.checked);$('hidesHair').onchange=e=>setField('hidesHair',e.target.checked);$('avatarAsset').onchange=e=>setField('avatarAsset',e.target.value);
$('figure').onchange=draw;$('showHead').onchange=draw;$('showAnchors').onchange=draw;$('scope').onchange=refresh;
$('reset').onclick=()=>{const k=scopeKey();if(k)delete meta[$('scope').value][k];refresh()};
$('save').onclick=async()=>{const st=$('status');const res=await fetch('/api/metadata',{method:'POST',body:JSON.stringify(meta)});const b=await res.json();st.className=res.ok?'':'err';st.textContent=res.ok?'Kaydedildi':b.error};
let drag=null;const c=$('c');c.onpointerdown=e=>{const r=currentRule();drag={x:e.offsetX,y:e.offsetY,ox:r.offsetX||0,oy:r.offsetY||0};c.setPointerCapture(e.pointerId)};
c.onpointermove=e=>{if(!drag||!scopeKey())return;const g=meta[$('scope').value],k=scopeKey();g[k]=g[k]||{};g[k].offsetX=+(drag.ox+(e.offsetX-drag.x)/W).toFixed(3);g[k].offsetY=+(drag.oy+(e.offsetY-drag.y)/H).toFixed(3);refresh()};c.onpointerup=()=>drag=null;
const first=allItems.find(i=>i.id==='axe_sunforged');if(first){equipped.mainhand=first.id;document.querySelector('[data-slot=mainhand]').value=first.id}refresh()}
start();
</script></body></html>'''


if __name__ == '__main__':
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 5175
    print('http://127.0.0.1:%d' % port)
    ThreadingHTTPServer(('127.0.0.1', port), Handler).serve_forever()
