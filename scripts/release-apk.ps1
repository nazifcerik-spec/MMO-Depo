# Builds the Android APK and publishes it as the in-app update (DD-304).
#
#   powershell -ExecutionPolicy Bypass -File scripts/release-apk.ps1 -Version 1.25.0
#   ... -NoPublish   build only (cable installs: adb install -r <apk>)
#
# The `releases` branch holds nothing but latest.json + the APK in a single
# commit that every release replaces (force push; never main). The game reads
# it through raw.githubusercontent.com. The APK is signed on this machine, so
# releases must be built here: the phone only accepts updates signed with the
# key of the installed app.
param(
    [Parameter(Mandatory = $true)][string]$Version,
    [switch]$NoPublish
)
$ErrorActionPreference = 'Stop'

if ($Version -notmatch '^\d+\.\d+\.\d+$') { throw 'Version must look like 1.25.0' }
$v = $Version.Split('.') | ForEach-Object { [int]$_ }
$buildNumber = $v[0] * 10000 + $v[1] * 100 + $v[2] # Android versionCode, always rising

$repo = Split-Path -Parent $PSScriptRoot
$work = 'C:\otter-build\client' # Gradle fails under the project's non-ASCII path
$api = 'https://new-world-mmo-api.onrender.com'
$feedBase = 'https://raw.githubusercontent.com/nazifcerik-spec/MMO-Depo/releases/'

robocopy "$repo\client" $work /MIR /XD build .dart_tool /NFL /NDL /NJH /NJS /NP | Out-Null
if ($LASTEXITCODE -ge 8) { throw "robocopy failed ($LASTEXITCODE)" }

Push-Location $work
try {
    flutter build apk --release --target-platform android-arm,android-arm64 `
        --build-name=$Version --build-number=$buildNumber `
        --dart-define=API_BASE_URL=$api `
        --dart-define=APP_VERSION=$Version `
        --dart-define=UPDATE_FEED_URL=${feedBase}latest.json
    if ($LASTEXITCODE -ne 0) { throw 'flutter build failed' }
} finally {
    Pop-Location
}
$apk = "$work\build\app\outputs\flutter-apk\app-release.apk"
Write-Host "built $apk ($([math]::Round((Get-Item $apk).Length / 1MB, 1)) MB)"
if ($NoPublish) { return }

$pub = Join-Path $env:TEMP 'mena-release'
if (Test-Path $pub) { Remove-Item -Recurse -Force $pub }
New-Item -ItemType Directory $pub | Out-Null
$name = "mena-world-$Version.apk"
Copy-Item $apk "$pub\$name"
$size = (Get-Item "$pub\$name").Length
[ordered]@{ version = $Version; apk = $name; size = $size } | ConvertTo-Json |
    Set-Content -Encoding ascii "$pub\latest.json"

$origin = git -C $repo remote get-url origin
git -C $pub init -q -b releases
git -C $pub add -A
git -C $pub commit -q -m "MeNa World $Version"
# Large uploads sometimes drop mid-way: try a few times.
$pushed = $false
for ($i = 1; $i -le 4 -and -not $pushed; $i++) {
    git -C $pub push -q --force $origin releases
    $pushed = $LASTEXITCODE -eq 0
    if (-not $pushed) { Start-Sleep 20 }
}
if (-not $pushed) { throw 'push failed' }
Write-Host "published $Version -> ${feedBase}latest.json"
