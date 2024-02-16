import win32com.client as win32
from win32com.client import constants
import os
import csv
from datetime import datetime

# 基本となるCSVファイルのパス
base_csv_path = 'word_process_results.csv'
# ファイルが既に存在するかチェック
if os.path.exists(base_csv_path):
    # 現在の日付と時間をファイル名に追加
    timestamp = datetime.now().strftime("_%Y%m%d%H%M%S")
    output_csv_path = base_csv_path.replace('.csv', f'{timestamp}.csv')
else:
    output_csv_path = base_csv_path

def set_document_readonly(filepath, edit_password):
    try:
        word = win32.gencache.EnsureDispatch('Word.Application')
        word.Visible = False
        doc = word.Documents.Open(filepath)
        
        if doc.ProtectionType == constants.wdNoProtection:
            doc.Protect(Type=constants.wdAllowOnlyReading, NoReset=True, Password=edit_password)
            result = 'OK'
        else:
            result = 'PASS'
        
        doc.Save()
        doc.Close()
        word.Quit()
        return result
    except Exception as e:
        word.Quit()
        return 'NG'

def process_files(directory, edit_password):
    results = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith('.docx'):
                filepath = os.path.join(root, file)
                result = set_document_readonly(filepath, edit_password)
                results.append({
                    'NAME': os.path.basename(filepath),
                    'RESULT': result,
                    'PATH': filepath
                })
    return results

def write_results_to_csv(results, csv_path):
    with open(csv_path, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = ['NAME', 'RESULT', 'PATH']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        writer.writeheader()
        for result in results:
            writer.writerow(result)

if __name__ == '__main__':
    edit_password = 'your_edit_password'
    current_directory = os.path.dirname(os.path.abspath(__file__))
    results = process_files(current_directory, edit_password)
    write_results_to_csv(results, output_csv_path)
    print(f'処理結果は{output_csv_path}に保存されました。'.encode('utf-8').decode('utf-8'))
