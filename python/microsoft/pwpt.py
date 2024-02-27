import os
import sys
import csv
import win32com.client
import pythoncom
from threading import Thread
from queue import Queue
from datetime import datetime
from pywinauto.application import Application
from pywinauto.findwindows import find_window
import time


class DialogWatcher(Thread):
    def __init__(self, titles, button_title="キャンセル"):
        super().__init__(daemon=True)
        self.titles = titles
        self.button_title = button_title
        self.running = True

    def run(self):
        while self.running:
            for title in self.titles:
                try:
                    hwnd = find_window(title=title)
                    if hwnd:
                        app = Application().connect(handle=hwnd)
                        dialog = app.window(handle=hwnd)
                        dialog[self.button_title].click()
                        print(f"Clicked '{self.button_title}' on '{title}' dialog")
                except Exception as e:
                    print(f"Error handling dialog: {e}")
            time.sleep(5)  # 5秒ごとに確認

    def stop(self):
        self.running = False


def open_office_application(extension):
    apps = {
        ".pptx": "PowerPoint.Application",
        ".docx": "Word.Application",
        ".xlsx": "Excel.Application",
    }
    app_name = apps.get(extension)
    if app_name:
        if extension == ".xlsx":  # Excelの場合のみgencache.EnsureDispatchを使用
            return win32com.client.gencache.EnsureDispatch(app_name)
        else:
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


def set_readonly(file_path, password, extension):
    result_queue = Queue()
    thread = Thread(target=worker, args=(file_path, password, extension, result_queue))
    thread.start()
    thread.join()  # タイムアウト指定を削除
    return result_queue.get()


def process_directory_for_documents(directory, edit_password):
    results = []
    total_files = sum(
        [
            len(
                [
                    file
                    for file in files
                    if (
                        file.endswith((".docx", ".pptx", ".xlsx"))
                        and not file.startswith("~$")
                    )
                ]
            )
            for _, _, files in os.walk(directory)
        ]
    )
    print(f"合計で処理するファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith((".docx", ".pptx", ".xlsx")) and not file.startswith("~$"):
                filepath = os.path.join(root, file)
                result = set_readonly(
                    filepath, edit_password, os.path.splitext(file)[1].lower()
                )
                results.append(
                    {
                        "NAME": os.path.basename(filepath),
                        "RESULT": result,
                        "PATH": filepath,
                    }
                )
                file_count += 1
                print_progress(file, result, file_count, total_files)

    if not results:
        print("指定されたディレクトリに対象のファイルが見つかりません。")
    return results


def print_progress(file_name, result, file_count, total_files):
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"[{current_time}][{result}] {file_count}/{total_files} | {file_name}")


def write_results_to_csv(results, output_csv_path):
    with open(output_csv_path, "w", newline="", encoding="utf-8-sig") as csvfile:
        fieldnames = ["NAME", "RESULT", "PATH"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for result in results:
            writer.writerow(result)


def get_application_path():
    if getattr(sys, "frozen", False):
        return os.path.dirname(sys.executable)
    else:
        return os.path.dirname(os.path.abspath(__file__))


if __name__ == "__main__":
    edit_password = "your_edit_password"
    current_directory = get_application_path()
    os.chdir(current_directory)
    print("Starting file processing...")

    dialog_watcher = DialogWatcher(titles=["確認ウィンドウ1", "確認ウィンドウ2"])
    dialog_watcher.start()

    results = process_directory_for_documents(current_directory, edit_password)

    dialog_watcher.stop()
    dialog_watcher.join()

    if results:
        print(f"Results have been processed.")
        output_csv_path = os.path.join(current_directory, "results.csv")
        write_results_to_csv(results, output_csv_path)
        print(f"Results have been saved to {output_csv_path}.")
    else:
        print("No files were processed.")
    input("Processing complete. Press Enter to exit...")
