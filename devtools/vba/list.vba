Sub GenerateList()
    Dim ws As Worksheet
    Dim lastRow As Long
    Dim i As Long
    Dim startNum As Long
    Dim endNum As Long
    Dim j As Long
    Dim keyword As String
    Dim outputRange As Range
    
    Set ws = ActiveSheet
    lastRow = ws.Cells(ws.Rows.Count, "B").End(xlUp).Row
    
    For i = 1 To lastRow
        keyword = LCase(ws.Cells(i, "A").Value)
        startNum = ws.Cells(i, "B").Value
        endNum = ws.Cells(i, "C").Value
        
        If keyword = "between" Then
            For j = startNum To endNum
                ws.Cells(ws.Rows.Count, "D").End(xlUp).Offset(1, 0).Value = j
            Next j
        ElseIf keyword = "only" Then
            ws.Cells(ws.Rows.Count, "D").End(xlUp).Offset(1, 0).Value = startNum
        End If
    Next i
    
    Set outputRange = ws.Range("D1:D" & ws.Cells(ws.Rows.Count, "D").End(xlUp).Row)
    
    ' 先に適用されている条件付き書式をクリアします
    outputRange.FormatConditions.Delete
    
    ' 重複値を青色でハイライトする条件付き書式を追加します
    outputRange.FormatConditions.Add Type:=xlExpression, Formula1:="=COUNTIF($D:$D, D1)>1"
    outputRange.FormatConditions(outputRange.FormatConditions.Count).SetFirstPriority
    With outputRange.FormatConditions(1)
        .Font.Color = -16776961  '青色
        .TintAndShade = 0
    End With
End Sub
