Param(
    [string]$log_dir
)
$log_dir = Convert-Path $log_dir

{{range $i, $v := .}}
echo TestId::{{$v.Id}}
$log_path = Join-Path $log_dir "{{$v.Id}}"
Invoke-Command  -ScriptBlock { 
    {{$v.Text}} 
} | Out-File $log_path -Encoding UTF8{{end}}
