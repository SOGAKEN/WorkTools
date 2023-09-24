import pandas as pd
import argparse
import os
from datetime import datetime

def process_excel(input_file):
    df = pd.read_excel(input_file)
    
    df['new_col'] = df.iloc[:,1] * 2 

    # 最後の列に追加
    df['new_col'] = df.pop('new_col')  
    
    # ファイル名から拡張子を除去し、先頭に"new_"と日付を追加
    base_name = os.path.basename(input_file)  # ファイル名のみ取得
    file_name_without_extension = os.path.splitext(base_name)[0]  # 拡張子を除去
    current_date = datetime.now().strftime('%Y%m%d')  # 現在の日付を"yyyymmdd"形式で取得
    output_file = f"new_{current_date}_{file_name_without_extension}.xlsx"

    df.to_excel(output_file, index=False)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process an Excel file.')
    parser.add_argument('input_file', type=str, help='Path to the input Excel file')
    
    args = parser.parse_args()
    
    process_excel(args.input_file)
