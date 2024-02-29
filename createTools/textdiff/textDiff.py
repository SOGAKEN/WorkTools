import csv
import tkinter as tk
from tkinter import filedialog


def compare_files(file1_path, file2_path, output_csv):
    with open(file1_path, "r", encoding="utf-8") as file1, open(
        file2_path, "r", encoding="utf-8"
    ) as file2:
        file1_lines = file1.readlines()
        file2_lines = file2.readlines()

    differences = []
    max_lines = max(len(file1_lines), len(file2_lines))

    for i in range(max_lines):
        line1 = file1_lines[i].strip() if i < len(file1_lines) else ""
        line2 = file2_lines[i].strip() if i < len(file2_lines) else ""
        if line1 != line2:
            differences.append((i + 1, line1, line2))

    with open(output_csv, "w", newline="", encoding="utf-8") as csvfile:
        csvwriter = csv.writer(csvfile)
        csvwriter.writerow(["Line", "File1", "File2"])
        for diff in differences:
            csvwriter.writerow(diff)

    if differences:
        print("NG: 差分あり")
        for diff in differences:
            print(f"Line {diff[0]}: File1='{diff[1]}' vs File2='{diff[2]}'")
    else:
        print("OK: 差分なし")

    return differences


def select_file(title="ファイルを選択"):
    root = tk.Tk()
    root.withdraw()  # Tkのルートウィンドウを表示しない
    file_path = filedialog.askopenfilename(title=title)
    return file_path


if __name__ == "__main__":
    file1_path = select_file("ファイル1を選択してください")
    file2_path = select_file("ファイル2を選択してください")
    output_csv = "comparison_result.csv"

    compare_files(file1_path, file2_path, output_csv)

    input("終了するにはエンターキーを押してください...")
