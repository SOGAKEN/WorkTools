Sub GenerateList()
    Dim ws As Worksheet
    Dim lastRow As Long
    Dim i As Long
    Dim startNum As Long
    Dim endNum As Long
    Dim j As Long
    Dim keyword As String
    
    Set ws = ActiveSheet
    lastRow = ws.Cells(ws.Rows.Count, "B").End(xlUp).Row
    
    For i = 1 To lastRow
        keyword = LCase(ws.Cells(i, "A").Value)
        startNum = ws.Cells(i, "B").Value
        endNum = ws.Cells(i, "C").Value
        
        If keyword = "in range" Then
            For j = startNum To endNum
                ws.Cells(ws.Rows.Count, "J").End(xlUp).Offset(1, 0).Value = j
            Next j
        ElseIf keyword = "equal to" Then
            ws.Cells(ws.Rows.Count, "J").End(xlUp).Offset(1, 0).Value = startNum
        End If
    Next i
End Sub
