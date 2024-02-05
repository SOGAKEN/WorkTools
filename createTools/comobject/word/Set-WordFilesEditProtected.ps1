param (
    [String]$password,
    [ValidateSet('ReadOnly', 'Comments', 'Forms', 'TrackedChanges', 'None')]
    [String]$editingRestriction = 'None'
)

function Set-DocumentProtection {
    param (
        [Object]$document,
        [String]$password,
        [String]$editingRestriction
    )

    switch ($editingRestriction) {
        'ReadOnly' { $restrictionType = 2 }
        'Comments' { $restrictionType = 4 }
        'Forms' { $restrictionType = 3 }
        'TrackedChanges' { $restrictionType = 5 }
        'None' { $restrictionType = 0 }
    }

    if ($restrictionType -ne 0) {
        $document.Protect($restrictionType, $true, $password)
    }
}

$directoryPath = Split-Path -Parent $MyInvocation.MyCommand.Path

$word = New-Object -ComObject Word.Application
$word.Visible = $false

Get-ChildItem -Path $directoryPath -Filter *.docx | ForEach-Object {
    $document = $word.Documents.Open($_.FullName)
    
    Set-DocumentProtection -document $document -password $password -editingRestriction $editingRestriction
    
    $document.Save()
    $document.Close()
}

$word.Quit()
[System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
[System.GC]::Collect()
[System.GC]::WaitForPendingFinalizers()
