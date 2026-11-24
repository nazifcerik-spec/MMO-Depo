# Builds the web client and publishes it to GitHub Pages (DD-369).
#
#   powershell -ExecutionPolicy Bypass -File scripts/release-web.ps1 -Version 1.49.0
#   ... -NoPublish   build only (client build in C:\otter-build\client\build\web)
#
# The site is the `gh-pages` branch: https://nazifcerik-spec.github.io/MMO-Depo/
# Every release is an ordinary commit on top of the last (no force push; art
# that did not change is not stored again). The API answers that origin
# through CORS_ALLOWED_ORIGINS.
param(
    [Parameter(Mandatory = $true)][string]$Version,
    [switch]$NoPublish
)
$ErrorActionPreference = 'Stop'

if ($Version -notmatch '^\d+\.\d+\.\d+$') { throw 'Version must look like 1.49.0' }

$repo = Split-Path -Parent $PSScriptRoot
$work = 'C:\otter-build\client' # shared with release-apk.ps1
$api = 'https://new-world-mmo-api.onrender.com'
$site = 'https://nazifcerik-spec.github.io/MMO-Depo/'

robocopy "$repo\client" $work /MIR /XD build .dart_tool /NFL /NDL /NJH /NJS /NP | Out-Null
if ($LASTEXITCODE -ge 8) { throw "robocopy failed ($LASTEXITCODE)" }

Push-Location $work
try {
    # --pwa-strategy=none: no service worker, so a release is seen on reload.
    flutter build web --release --base-href /MMO-Depo/ --pwa-strategy=none `
        --build-name=$Version `
        --dart-define=API_BASE_URL=$api `
        --dart-define=APP_VERSION=$Version
    if ($LASTEXITCODE -ne 0) { throw 'flutter build failed' }
} finally {
    Pop-Location
}
$out = "$work\build\web"
Get-ChildItem "$out\canvaskit" -Recurse -Filter *.symbols | Remove-Item -Force # debugging only
Copy-Item "$out\index.html" "$out\404.html" # any path opens the game
Set-Content -Encoding ascii "$out\.nojekyll" ''
$mb = [math]::Round((Get-ChildItem $out -Recurse -File | Measure-Object Length -Sum).Sum / 1MB, 1)
Write-Host "built $out ($mb MB)"
if ($NoPublish) { return }

$pub = Join-Path $env:TEMP 'mena-web'
if (Test-Path $pub) { Remove-Item -Recurse -Force $pub }
$origin = git -C $repo remote get-url origin
git clone -q --depth 1 --branch gh-pages $origin $pub
if ($LASTEXITCODE -ne 0) { # first release: the branch does not exist yet
    New-Item -ItemType Directory -Force $pub | Out-Null
    git -C $pub init -q -b gh-pages
    git -C $pub remote add origin $origin
}
robocopy $out $pub /MIR /XD .git /NFL /NDL /NJH /NJS /NP | Out-Null
if ($LASTEXITCODE -ge 8) { throw "robocopy failed ($LASTEXITCODE)" }
git -C $pub add -A
git -C $pub -c core.safecrlf=false commit -q -m "MeNa World web $Version"
$pushed = $false
for ($i = 1; $i -le 4 -and -not $pushed; $i++) { # large uploads sometimes drop mid-way
    git -C $pub push -q origin gh-pages
    $pushed = $LASTEXITCODE -eq 0
    if (-not $pushed) { Start-Sleep 20 }
}
if (-not $pushed) { throw 'push failed' }
Write-Host "published $Version -> $site"
