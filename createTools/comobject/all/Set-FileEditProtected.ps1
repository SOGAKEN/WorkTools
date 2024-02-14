param (
    [String]$directoryPath,
    [String]$outputCsv = "output.csv",
    [String]$password
)

function Set-ReadOnly {
    param (
        [String]$filePath,
        [String]$fileType
    )

    $result = "OK"
    try {
        switch ($fileType) {
            "Word" {
                $word = New-Object -ComObject Word.Application
                $word.Visible = $false
                $document = $word.Documents.Open($filePath)
                if ($document.ProtectionType -eq -1) {
                    $document.Protect(2, $true, $password)
                    $document.Save()
                } else {
                    $result = "pass"
                }
                $document.Close()
                $word.Quit()
                [System.Runtime.InteropServices.Marshal]::ReleaseComObject($word) | Out-Null
            }
            "Excel" {
                $excel = New-Object -ComObject Excel.Application
                $excel.Visible = $false
                $workbook = $excel.Workbooks.Open($filePath)
                if ($workbook.HasPassword) {
                    $result = "pass"
                } else {
                    $workbook.Protect($password, $True, $True)
                    $workbook.Save()
                }
                $workbook.Close()
                $excel.Quit()
                [System.Runtime.InteropServices.Marshal]::ReleaseComObject($excel) | Out-Null
            }
            "PowerPoint" {
                $ppt = New-Object -ComObject PowerPoint.Application
                $presentation = $ppt.Presentations.Open($filePath)
                # PowerPointの保護状態の確認方法は限られています。
                # ここでは、保存時にエラーが発生しないかどうかで判断します（非推奨）。
                try {
                    $presentation.SaveAs($filePath, -2, $password)
                } catch {
                    $result = "pass"
                }
                $presentation.Close()
                $ppt.Quit()
                [System.Runtime.InteropServices.Marshal]::ReleaseComObject($ppt) | Out-Null
            }
        }
    } catch {
        $result = "NG"
    }

    return $result
}

# CSVヘッダー
"Name,Path,Result" | Out-File $outputCsv

Get-ChildItem -Path $directoryPath -Recurse | Where-Object { $_.Extension -match "\.(docx|xlsx|pptx)$" } | ForEach-Object {
    $fileType = switch ($_.Extension) {
        ".docx" { "Word" }
        ".xlsx" { "Excel" }
        ".pptx" { "PowerPoint" }
    }

    $result = Set-ReadOnly -filePath $_.FullName -fileType $fileType

    # CSVに結果を追加
    $csvLine = "$($_.Name),$($_.FullName),$result"
    $csvLine | Out-File $outputCsv -Append
}

[System.GC]::Collect()
[System.GC]::WaitForPendingFinalizers()
