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
    
    outputRange.FormatConditions.Add Type:=xlCellValue, Operator:=xlDuplicate
    outputRange.FormatConditions(outputRange.FormatConditions.Count).SetFirstPriority
    outputRange.FormatConditions(1).DupeUnique = xlDuplicate
    With outputRange.FormatConditions(1).Interior
        .PatternColorIndex = xlAutomatic
        .Color = 65535  'Yellow color
        .TintAndShade = 0
    End With
End Sub
