import win32com.client as win32
import os
import csv
import sys
import threading
import queue
from datetime import datetime

def get_application_path():
    """アプリケーションのパスを取得します。"""
    if getattr(sys, 'frozen', False):
        path = os.path.dirname(sys.executable)
    else:
        path = os.path.dirname(os.path.abspath(__file__))
    print(f"アプリケーションパス: {path}")
    return path

def get_csv_path(base_name='word_process_results.csv'):
    """CSVファイルのパスを生成します。"""
    base_path = os.path.join(get_application_path(), base_name)
    if os.path.exists(base_path):
        timestamp = datetime.now().strftime("_%Y%m%d%H%M%S")
        csv_path = base_path.replace('.csv', f'{timestamp}.csv')
    else:
        csv_path = base_path
    print(f"CSVファイルパス: {csv_path}")
    return csv_path

def set_document_readonly_with_timeout(filepath, edit_password, timeout=30):
    """ドキュメントを読み取り専用に設定します。タイムアウト機能付き。"""
    def target():
        try:
            print(f"ドキュメント開始: {filepath}")
            word = win32.gencache.EnsureDispatch('Word.Application')
            word.Visible = False
            doc = word.Documents.Open(filepath)
            if doc.ProtectionType == win32.constants.wdNoProtection:
                doc.Protect(Type=win32.constants.wdAllowOnlyReading, NoReset=True, Password=edit_password)
                print(f"読み取り専用に設定: {filepath}")
                result_queue.put('OK')
            else:
                print(f"既に保護されています: {filepath}")
                result_queue.put('PASS')
            doc.Save()
            doc.Close(False)
        except Exception as e:
            print(f"エラーが発生しました: {e}, ファイル: {filepath}")
            result_queue.put('NG')
        finally:
            if 'word' in locals():
                word.Quit()

    result_queue = queue.Queue()
    thread = threading.Thread(target=target)
    thread.start()
    thread.join(timeout)
    if thread.is_alive():
        print(f"タイムアウト: 処理が長すぎます: {filepath}")
        # スレッドが終了するのを強制的に待たない（Wordが応答しない場合には処理を続行）
        return 'TIMEOUT'
    return result_queue.get()

def process_directory_for_docx(directory, edit_password):
    """ディレクトリ内のdocxファイルを処理します。"""
    results = []
    total_files = sum(len(files) for _, _, files in os.walk(directory) if any(file.endswith('.docx') for file in files))
    print(f"合計で処理する.docxファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in filter(lambda f: f.endswith('.docx'), files):
            filepath = os.path.join(root, file)
            result = set_document_readonly_with_timeout(filepath, edit_password)
            results.append({'NAME': os.path.basename(filepath), 'RESULT': result, 'PATH': filepath})
            file_count += 1
            print_progress(file, result, file_count, total_files)

    if not results:
        print("指定されたディレクトリに.docxファイルが見つかりません。")
    return results

def print_progress(file_name, result, file_count, total_files):
    """処理の進捗を表示します。"""
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"進捗: {file_count}/{total_files}, ファイル名: {file_name}, 結果: {result}, 時刻: {current_time}")

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
    edit_password = 'your_edit_password'
    current_directory = get_application_path()
    results = process_directory_for_docx(current_directory, edit_password)
    if results:
        output_csv_path = get_csv_path()
        write_results_to_csv(results, output_csv_path)
        print(f'処理結果は{output_csv_path}に保存されました。')
    else:
        print("処理するファイルが見つかりませんでした。")

    input("処理が完了しました。エンターキーを押して終了してください...")
