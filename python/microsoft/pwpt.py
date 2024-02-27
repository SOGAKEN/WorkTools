import os
import sys
import win32com.client
import pythoncom
from threading import Thread
from queue import Queue
from datetime import datetime


def open_office_application(extension):
    apps = {
        ".pptx": "PowerPoint.Application",
        ".docx": "Word.Application",
        ".xlsx": "Excel.Application",
    }
    app_name = apps.get(extension)
    if app_name:
        return win32com.client.Dispatch(app_name)
    else:
        raise ValueError("Unsupported file type")


def set_document_or_presentation_readonly(
    app, file_path, password, extension, result_queue
):
    try:
        if extension == ".pptx":
            doc = app.Presentations.Open(file_path, WithWindow=False)
            doc.WritePassword = password
            doc.Save()
            doc.Close()
        elif extension == ".docx":
            doc = app.Documents.Open(file_path)
            doc.WritePassword = password
            doc.SaveAs2(file_path, WritePassword=password)
            doc.Close()
        elif extension == ".xlsx":
            app.DisplayAlerts = False
            doc = app.Workbooks.Open(file_path)
            doc.SaveAs(file_path, Password="", WriteResPassword=password)
            doc.Close()
        result_queue.put("SUCCESS")
    except Exception as e:
        result_queue.put(f"ERROR: {e}")
    finally:
        app.Quit()


def worker(file_path, password, extension, result_queue):
    pythoncom.CoInitialize()
    try:
        app = open_office_application(extension)
        set_document_or_presentation_readonly(
            app, file_path, password, extension, result_queue
        )
    finally:
        pythoncom.CoUninitialize()


def set_readonly_with_timeout(file_path, password, timeout_seconds=30):
    result_queue = Queue()
    extension = os.path.splitext(file_path)[1].lower()
    thread = Thread(target=worker, args=(file_path, password, extension, result_queue))
    thread.start()
    thread.join(timeout=timeout_seconds)
    if thread.is_alive():
        print(
            f"WARNING: Processing of {file_path} timed out. The application may still be open."
        )
        return "TIMEOUT"
    else:
        return result_queue.get()


def process_directory_for_documents(directory, edit_password):
    results = []
    # `~$`で始まるファイルを除外する条件を追加
    total_files = sum(
        [
            len(
                [
                    file
                    for file in files
                    if (
                        file.endswith(".docx")
                        or file.endswith(".pptx")
                        or file.endswith(".xlsx")
                    )
                    and not file.startswith("~$")
                ]
            )
            for _, _, files in os.walk(directory)
        ]
    )
    print(f"合計で処理するファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in filter(
            lambda f: (
                f.endswith((".docx", ".pptx", ".xlsx")) and not f.startswith("~$")
            ),
            files,
        ):
            filepath = os.path.join(root, file)
            result = set_readonly_with_timeout(filepath, edit_password)
            results.append(
                {"NAME": os.path.basename(filepath), "RESULT": result, "PATH": filepath}
            )
            file_count += 1
            print_progress(file, result, file_count, total_files)

    if not results:
        print("指定されたディレクトリに対象のファイルが見つかりません。")
    return results


def print_progress(file_name, result, file_count, total_files):
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"[{current_time}][{result}] {file_count}/{total_files} | {file_name}")


def get_application_path():
    if getattr(sys, "frozen", False):
        return os.path.dirname(sys.executable)
    else:
        return os.path.dirname(os.path.abspath(__file__))


if __name__ == "__main__":
    edit_password = "your_edit_password"
    current_directory = get_application_path()
    os.chdir(current_directory)  # カレントディレクトリを変更
    print("Starting file processing...")
    results = process_directory_for_documents(current_directory, edit_password)
    if results:
        print(f"Results have been processed.")
    else:
        print("No files were processed.")
    input("Processing complete. Press Enter to exit...")
