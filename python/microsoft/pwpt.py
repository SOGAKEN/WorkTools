import os
import csv
import win32com.client
import pythoncom
from threading import Timer
from queue import Queue


def open_office_application(extension):
    """
    Officeアプリケーションを開くためのヘルパー関数。
    """
    apps = {
        ".pptx": "PowerPoint.Application",
        ".docx": "Word.Application",
        ".xlsx": "Excel.Application"
    }
    app_name = apps.get(extension)
    if app_name:
        if extension == ".docx":
            return win32com.client.Dispatch(app_name)  # Wordの場合はDispatchを使用
        else:
            return win32com.client.gencache.EnsureDispatch(app_name)
    else:
        raise ValueError("Unsupported file type")


def set_password_pptx(app, file_path, password):
    doc = app.Presentations.Open(file_path, WithWindow=False)
    doc.WritePassword = password
    doc.Save()
    doc.Close()


def set_password_docx(app, file_path, password):
    doc = app.Documents.Open(file_path)
    doc.WritePassword = password
    doc.SaveAs2(file_path, WritePassword=password)
    doc.Close()


def set_password_xlsx(app, file_path, password, result_queue):
    app.Visible = False
    doc = app.Workbooks.Open(file_path)
    if doc.WriteReserved:
        result_queue.put(
            ("PASS", "File is password protected or another error occurred"))
    else:
        doc.Password = password
        doc.SaveAs(file_path, Password='', WriteResPassword=password)
    doc.Close()


def set_password_office(file_path, password, result_queue):
    pythoncom.CoInitialize()
    try:
        extension = os.path.splitext(file_path)[1].lower()
        app = open_office_application(extension)
        if extension == ".pptx":
            set_password_pptx(app, file_path, password)
        elif extension == ".docx":
            set_password_docx(app, file_path, password)
        elif extension == ".xlsx":
            set_password_xlsx(app, file_path, password, result_queue)
        else:
            result_queue.put(("NG", "Unsupported file type"))
            return
        app.Quit()
        result_queue.put(("OK", "Password set successfully"))
    except Exception as e:
        result_queue.put(("ERROR", f"Error setting password: {e}"))
    finally:
        pythoncom.CoUninitialize()

# その他の関数は基本的に変更なしで、必要に応じてエラーハンドリングを追加または調整します。


if __name__ == "__main__":
    # メインの実行ロジックは変更なし
