import win32com.client as win32
import os
import csv
import sys
from datetime import datetime

def get_application_path():
    if getattr(sys, 'frozen', False):
        return os.path.dirname(sys.executable)
    return os.path.dirname(os.path.abspath(__file__))

def get_csv_path(base_name='word_process_results.csv'):
    base_path = os.path.join(get_application_path(), base_name)
    if os.path.exists(base_path):
        timestamp = datetime.now().strftime("_%Y%m%d%H%M%S")
        return base_path.replace('.csv', f'{timestamp}.csv')
    return base_path

def set_document_readonly(filepath, edit_password):
    try:
        word = win32.gencache.EnsureDispatch('Word.Application')
        word.Visible = False
        doc = word.Documents.Open(filepath)
        if doc.ProtectionType == win32.constants.wdNoProtection:
            doc.Protect(Type=win32.constants.wdAllowOnlyReading, NoReset=True, Password=edit_password)
            result = 'OK'
        else:
            result = 'PASS'
        doc.Save()
        doc.Close(False)
    except Exception as e:
        print(f"エラーが発生しました: {e}")
        result = 'NG'
    finally:
        if 'word' in locals():
            word.Quit()
    return result

def process_directory_for_docx(directory, edit_password):
    results = []
    total_files = sum(len(files) for _, _, files in os.walk(directory) if any(file.endswith('.docx') for file in files))
    print(f"合計で処理する.docxファイルの数: {total_files}")

    file_count = 0
    for root, _, files in os.walk(directory):
        for file in filter(lambda f: f.endswith('.docx'), files):
            filepath = os.path.join(root, file)
            result = set_document_readonly(filepath, edit_password)
            results.append({'NAME': os.path.basename(filepath), 'RESULT': result, 'PATH': filepath})
            file_count += 1
            print_progress(file, result, file_count, total_files)

    if not results:
        print("指定されたディレクトリに.docxファイルが見つかりません。")
    return results

def print_progress(file_name, result, file_count, total_files):
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"進捗: {file_count}/{total_files}, ファイル名: {file_name}, 結果: {result}, 時刻: {current_time}")

def write_results_to_csv(results, csv_path):
    with open(csv_path, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = ['NAME', 'RESULT', 'PATH']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for result in results:
            writer.writerow(result)

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
