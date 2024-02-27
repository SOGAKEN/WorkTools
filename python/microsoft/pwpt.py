import win32com.client as win32
import os
import csv
import sys
import threading
import queue
import pythoncom  # COMライブラリの初期化とクリーンアップに必要
from datetime import datetime
from pptx import Presentation
from pptx.exc import PackageNotFoundError
from openpyxl import load_workbook
from openpyxl.utils.exceptions import InvalidFileException


def get_application_path():
    """アプリケーションのパスを取得します。"""
    if getattr(sys, "frozen", False):
        path = os.path.dirname(sys.executable)
    else:
        path = os.path.dirname(os.path.abspath(__file__))
    return path


def get_csv_path(base_name="process_results.csv"):
    """CSVファイルのパスを生成します。"""
    base_path = os.path.join(get_application_path(), base_name)
    if os.path.exists(base_path):
        timestamp = datetime.now().strftime("_%Y%m%d%H%M%S")
        csv_path = base_path.replace(".csv", f"{timestamp}.csv")
    else:
        csv_path = base_path
    return csv_path


def can_open_file(filepath):
    """指定されたファイルが開けるかどうかを判定します。"""
    try:
        if filepath.endswith(".pptx"):
            Presentation(filepath)
        elif filepath.endswith(".xlsx"):
            load_workbook(filepath)
        return True
    except (PackageNotFoundError, InvalidFileException):
        # ファイルが開けない、またはパスワードで保護されている場合
        return False


def set_write_password_with_timeout(filepath, edit_password, timeout=30):
    """ファイルに書き込みパスワードを設定します。タイムアウト機能付き。"""

    def target(result_queue):
        pythoncom.CoInitialize()  # スレッドでCOMライブラリを初期化
        try:
            if filepath.endswith(".docx"):
                word = win32.gencache.EnsureDispatch("Word.Application")
                doc = word.Documents.Open(filepath)
                doc.Password = edit_password
                result_queue.put("OK")
                doc.Save()
                doc.Close(False)
                word.Quit()
            elif filepath.endswith(".pptx"):
                powerpoint = win32.gencache.EnsureDispatch(
                    "PowerPoint.Application")
                presentation = powerpoint.Presentations.Open(
                    filepath, WithWindow=False)
                presentation.WritePassword = edit_password
                result_queue.put("OK")
                presentation.Save()
                presentation.Close()
                powerpoint.Quit()
            elif filepath.endswith(".xlsx"):
                excel = win32.gencache.EnsureDispatch("Excel.Application")
                workbook = excel.Workbooks.Open(filepath)
                workbook.Password = edit_password
                result_queue.put("OK")
                workbook.Save()
                workbook.Close(False)
                excel.Quit()
        except Exception as e:
            result_queue.put("NG")
        finally:
            pythoncom.CoUninitialize()

    result_queue = queue.Queue()
    thread = threading.Thread(target=target, args=(result_queue,))
    thread.start()
    thread.join(timeout)
    if thread.is_alive():
        thread.join()  # Ensure thread has finished
        return "TIMEOUT"
    return result_queue.get()


def process_directory_for_documents(directory, edit_password):
    """ディレクトリ内のファイルを処理します。"""
    results = []
    total_files = sum(
        [
            len(files)
            for _, _, files in os.walk(directory)
            if any(file.endswith((".docx", ".pptx", ".xlsx")) for file in files)
        ]
    )
    print(f"合計で処理するファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in filter(
            lambda f: f.endswith(".docx") or f.endswith(
                ".pptx") or f.endswith(".xlsx"),
            files,
        ):
            filepath = os.path.join(root, file)
            if not can_open_file(filepath):
                print(f"ファイルが開けないためスキップします: {filepath}")
                results.append(
                    {
                        "NAME": os.path.basename(filepath),
                        "RESULT": "SKIP",
                        "PATH": filepath,
                    }
                )
                continue
            result = set_write_password_with_timeout(filepath, edit_password)
            results.append(
                {"NAME": os.path.basename(
                    filepath), "RESULT": result, "PATH": filepath}
            )
            file_count += 1
            print_progress(file, result, file_count, total_files)

    if not results:
        print("指定されたディレクトリに対象のファイルが見つかりません。")
    return results


def print_progress(file_name, result, file_count, total_files):
    """処理の進捗を表示します。"""
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"[{current_time}][{result}] {file_count}/{total_files} | {file_name}")


def write_results_to_csv(results, csv_path):
    """結果をCSVに書き込みます。"""
    with open(csv_path, "w", newline="", encoding="utf-8-sig") as csvfile:
        fieldnames = ["NAME", "RESULT", "PATH"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for result in results:
            writer.writerow(result)
    print(f"CSVに結果を書き込みました: {csv_path}")


if __name__ == "__main__":
    edit_password = "your_edit_password"  # ここに適切なパスワードを設定してください
    current_directory = get_application_path()
    results = process_directory_for_documents(current_directory, edit_password)
    if results:
        output_csv_path = get_csv_path()
        write_results_to_csv(results, output_csv_path)
        print(f"処理結果は{output_csv_path}に保存されました。")
    else:
        print("処理するファイルが見つかりませんでした。")

    input("処理が完了しました。エンターキーを押して終了してください...")
