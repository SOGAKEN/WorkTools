param (
    [String]$excelFilePath
)

$excel = New-Object -ComObject Excel.Application
$excel.Visible = $False
$workbook = $excel.Workbooks.Open($excelFilePath)

# ブック保護に使用するパスワードを設定
$password = "YourPassword"

$workbook.Protect($password, $True, $True)
$workbook.Save()
$workbook.Close()

$excel.Quit()
[System.Runtime.InteropServices.Marshal]::ReleaseComObject($excel)
