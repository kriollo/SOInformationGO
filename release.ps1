param(
    [string]$Version
)

# --- 1. Control de versión ---
if (-not $Version) {
    $lastTag = git tag --sort=-v:refname | Select-Object -First 1
    if ($lastTag -match '^v(\d+)\.(\d+)\.(\d+)$') {
        $major = [int]$matches[1]
        $minor = [int]$matches[2]
        $patch = [int]$matches[3] + 1
        $Version = "v$major.$minor.$patch"
    } else {
        $Version = "v1.0.0"
    }
    Write-Host "Usando versión: $Version"
}

# --- 2. Compilar binarios ---
Write-Host "Compilando binario para Windows..."
go build -o build/systeminfoGo.exe main.go
Write-Host "Compilando binario para Linux..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build/systeminfoGo main.go

# --- 3. Generar CHANGELOG.md ---
$changelogFile = "CHANGELOG.md"
if (Test-Path $changelogFile) {
    $oldChangelog = Get-Content $changelogFile
} else {
    $oldChangelog = @()
}
$lastTag = git tag --sort=-v:refname | Select-Object -First 1
if ($lastTag) {
    $log = git log $lastTag..HEAD --pretty=format:"* %s"
} else {
    $log = git log --pretty=format:"* %s"
}
$today = Get-Date -Format "yyyy-MM-dd"
$header = "## $Version - $today"
$changelogEntry = @($header) + $log + "" + $oldChangelog
Set-Content $changelogFile $changelogEntry

# --- 4. Git add, commit, push ---
git add .
git commit -m "Release $Version"
git push origin main

# --- 5. Crear y subir tag ---
git tag $Version
git push origin $Version

# --- 6. Crear release en GitHub y subir binarios + changelog ---
$releaseTitle = "Release $Version"
$releaseNotes = Get-Content $changelogFile | Select-Object -First ($log.Count + 2) | Out-String

Write-Host "Creando release $Version en GitHub..."
gh release create $Version `
  build/systeminfoGo.exe build/systeminfoGo `
  --title "$releaseTitle" `
  --notes "$releaseNotes"

Write-Host "`n¡Release $Version creada, binarios y changelog subidos!"
