import win32com.client as win32
import os
import csv
import sys
import threading
import queue
import pythoncom  # COMライブラリの初期化とクリーンアップに必要
from datetime import datetime

def get_application_path():
    """アプリケーションのパスを取得します。"""
    if getattr(sys, 'frozen', False):
        path = os.path.dirname(sys.executable)
    else:
        path = os.path.dirname(os.path.abspath(__file__))
    return path

def get_csv_path(base_name='process_results.csv'):
    """CSVファイルのパスを生成します。"""
    base_path = os.path.join(get_application_path(), base_name)
    if os.path.exists(base_path):
        timestamp = datetime.now().strftime("_%Y%m%d%H%M%S")
        csv_path = base_path.replace('.csv', f'{timestamp}.csv')
    else:
        csv_path = base_path
    return csv_path

def set_document_or_presentation_readonly_with_timeout(filepath, edit_password, timeout=30):
    """ドキュメントまたはプレゼンテーションを読み取り専用に設定します。タイムアウト機能付き。"""
    def target(result_queue):
        pythoncom.CoInitialize()  # スレッドでCOMライブラリを初期化
        try:
            if filepath.endswith('.docx'):
                word = win32.gencache.EnsureDispatch('Word.Application')
                doc = word.Documents.Open(filepath)
                if doc.ProtectionType == win32.constants.wdNoProtection:
                    doc.Protect(Type=win32.constants.wdAllowOnlyReading, NoReset=True, Password=edit_password)
                    result_queue.put('OK')
                else:
                    result_queue.put('PASS')  # 既に保護されている場合
                doc.Save()
                doc.Close(False)
                word.Quit()
            elif filepath.endswith('.pptx'):
                powerpoint = win32.gencache.EnsureDispatch('PowerPoint.Application')
                try:
                    presentation = powerpoint.Presentations.Open(filepath, WithWindow=False)
                    # 書き込みパスワードを設定
                    presentation.PasswordEncryptionProvider = "Office Standard"
                    presentation.PasswordEncryptionAlgorithm = "RC4"
                    presentation.PasswordEncryptionKeyLength = 40
                    presentation.WritePassword = edit_password
                    # ファイルを保存して閉じる
                    presentation.Save()
                    presentation.Close()
                    result_queue.put('OK')
                except Exception as e:
                    # PowerPointファイルがパスワードで保護されている場合、ここで例外が発生
                    result_queue.put('PASS')
                finally:
                    powerpoint.Quit()
        except Exception as e:
            result_queue.put('NG')
        finally:
            pythoncom.CoUninitialize()

    result_queue = queue.Queue()
    thread = threading.Thread(target=target, args=(result_queue,))
    thread.start()
    thread.join(timeout)
    if thread.is_alive():
        thread.join()  # Ensure thread has finished
        return 'TIMEOUT'
    return result_queue.get()

def process_directory_for_documents(directory, edit_password):
    """ディレクトリ内のdocxおよびpptxファイルを処理します。"""
    results = []
    total_files = sum([len([file for file in files if file.endswith('.docx') or file.endswith('.pptx')]) for _, _, files in os.walk(directory)])
    print(f"合計で処理するファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in filter(lambda f: f.endswith('.docx') or f.endswith('.pptx'), files):
            filepath = os.path.join(root, file)
            result = set_document_or_presentation_readonly_with_timeout(filepath, edit_password)
            results.append({'NAME': os.path.basename(filepath), 'RESULT': result, 'PATH': filepath})
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
    with open(csv_path, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = ['NAME', 'RESULT', 'PATH']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for result in results:
            writer.writerow(result)
    print(f"CSVに結果を書き込みました: {csv_path}")

if __name__ == '__main__':
    edit_password = 'your_edit_password'  # ここに適切なパスワードを設定してください
    current_directory = get_application_path()
    results = process_directory_for_documents(current_directory, edit_password)
    if results:
        output_csv_path = get_csv_path()
        write_results_to_csv(results, output_csv_path)
        print(f'処理結果は{output_csv_path}に保存されました。')
    else:
        print("処理するファイルが見つかりませんでした。")

    input("処理が完了しました。エンターキーを押して終了してください...")
