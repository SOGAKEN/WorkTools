param (
    [String]$password
)

$directoryPath = Split-Path -Parent $MyInvocation.MyCommand.Path

$word = New-Object -ComObject Word.Application
$word.Visible = $false

Get-ChildItem -Path $directoryPath -Filter *.docx | ForEach-Object {
    $document = $word.Documents.Open($_.FullName)
    $document.Password = $password
    $document.SaveAs([ref] $_.FullName, [ref] 12, [ref] $false, [ref] $password) # 12 は docx 形式を意味します
    $document.Close()
}

$word.Quit()
[System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
