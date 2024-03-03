import csv
import os
from datetime import datetime
from wcwidth import wcswidth
import sys


def get_log_files():
    # 実行ファイルのディレクトリを取得（PyInstaller対応）
    if getattr(sys, "frozen", False):
        directory = os.path.dirname(sys.executable)
    else:
        directory = os.path.dirname(os.path.abspath(__file__))

    files = os.listdir(directory)
    file1 = file2 = None
    for file in files:
        if "before" in file and file.endswith(".log"):
            file1 = os.path.join(directory, file)
        elif "after" in file and file.endswith(".log"):
            file2 = os.path.join(directory, file)
    return file1, file2


def extract_sections(
    filename, keyword, lines_to_include=1, comma_sections_to_compare=None
):
    sections = []
    current_section = []
    line_count = 0
    with open(filename, "r") as file:
        for line in file:
            if line_count > 0:
                line_count -= 1
            if keyword in line:
                if current_section:  # 新しいセクションの開始前に現在のセクションを保存
                    sections.append("".join(current_section))
                    current_section = []
                line_count = lines_to_include
            if line_count > 0:
                if comma_sections_to_compare is not None:
                    parts = line.split(",")
                    line = (
                        ",".join(
                            [
                                parts[i]
                                for i in comma_sections_to_compare
                                if i < len(parts)
                            ]
                        )
                        + "\n"
                    )
                current_section.append(line)
        if current_section:  # ファイルの最後のセクションを追加
            sections.append("".join(current_section))
    return sections


def compare_sections(sections1, sections2):
    differences = []
    max_length = max(len(sections1), len(sections2))
    section_counts = {}  # セクション名の出現回数を追跡する辞書
    for section in sections1 + sections2:
        section_name = section.split(" ", 1)[0]
        section_counts[section_name] = section_counts.get(section_name, 0) + 1
    for i in range(max_length):
        section1 = sections1[i] if i < len(sections1) else ""
        section2 = sections2[i] if i < len(sections2) else ""
        if section1 != section2:
            # 以下のロジックは削除または修正します
            differences.append((i + 1, section1, section2))
    return differences


def write_to_csv(filename, data):
    """CSVファイルにデータを書き込む。ファイル名に現在の日時を追加する"""
    datetime_str = datetime.now().strftime("_%Y%m%d%H%M%S")
    filename = f"{filename}{datetime_str}.csv"
    with open(filename, "w", newline="", encoding="utf-8-sig") as csvfile:
        writer = csv.writer(csvfile, quoting=csv.QUOTE_MINIMAL)
        writer.writerow(["No.", "SectionName", "File1", "File2", "Result"])
        for row in data:
            # 各セルが文字列かどうかを確認し、必要に応じて文字列に変換
            processed_row = [
                (str(cell).replace("\n", "\n") if "\n" in str(cell) else str(cell))
                for cell in row
            ]
            writer.writerow(processed_row)


def print_with_color(text, color):
    """テキストを指定された色で出力する"""
    colors = {
        "red": "\033[91m",
        "green": "\033[92m",
        "end": "\033[0m",
    }
    print(colors[color] + text + colors["end"])


def calculate_display_length(text):
    """テキストの表示幅を計算する"""
    return wcswidth(text)


def process_files_to_csv(keywords_with_options, output_csv):
    file1, file2 = get_log_files()
    if not file1 or not file2:
        print("必要なログファイルが見つかりません。")
        return

    csv_data = []
    comparison_number = 1
    section_appearances = {}  # 全セクションの出現回数を計算するための辞書

    # 全セクションの出現回数を計算
    for keyword, options in keywords_with_options.items():
        section_name = options.get("section_name", keyword)
        section_appearances[section_name] = 0  # 初期化

    for file in [file1, file2]:
        for keyword, options in keywords_with_options.items():
            section_name = options.get("section_name", keyword)
            temp_options = {k: v for k, v in options.items() if k != "section_name"}
            sections = extract_sections(file, keyword, **temp_options)
            # セクションが存在するたびにカウントアップ
            section_appearances[section_name] += len(sections)

    # セクション名の表示幅の最大値を計算（数字を追加する前）
    max_section_name_length = max(
        calculate_display_length(section_name + " 999")
        for section_name in section_appearances.keys()
    )

    for keyword, options in keywords_with_options.items():
        section_name = options.get("section_name", keyword)
        temp_options = {k: v for k, v in options.items() if k != "section_name"}
        sections_file1 = extract_sections(file1, keyword, **temp_options)
        sections_file2 = extract_sections(file2, keyword, **temp_options)
        differences = compare_sections(sections_file1, sections_file2)

        section_number = 1  # セクション番号をリセット
        for diff in differences:
            _, section_file1, section_file2 = diff
            result = "NG"

            # セクション名に数字を追加（複数回出現する場合のみ）
            if section_appearances[section_name] > 2:
                numbered_section_name = f"{section_name} {section_number}"
                section_number += 1  # セクション番号をインクリメント
            else:
                numbered_section_name = section_name

            # セクション名の表示幅を調整
            adjusted_section_name = numbered_section_name + " " * (
                max_section_name_length
                - calculate_display_length(numbered_section_name)
            )
            output_text = f"No. {comparison_number:<3} | セクション名: {adjusted_section_name} | Result: {result}"
            if result == "NG":
                print_with_color(output_text, "red")
            else:
                print(output_text)
            # CSVデータにセクション名を含める
            csv_data.append(
                [
                    comparison_number,
                    adjusted_section_name.strip(),
                    section_file1.strip(),
                    section_file2.strip(),
                    result,
                ]
            )
            comparison_number += 1

        # 差異がないセクションに対する処理
        common_length = min(len(sections_file1), len(sections_file2))
        for i in range(common_length):
            if (i + 1, sections_file1[i], sections_file2[i]) not in differences:
                result = "OK"
                if section_appearances[section_name] > 2:
                    numbered_section_name = f"{section_name} {section_number}"
                    section_number += 1
                else:
                    numbered_section_name = section_name

                adjusted_section_name = numbered_section_name + " " * (
                    max_section_name_length
                    - calculate_display_length(numbered_section_name)
                )
                output_text = f"No. {comparison_number:<3} | セクション名: {adjusted_section_name} | Result: {result}"
                print(output_text)
                csv_data.append(
                    [
                        comparison_number,
                        adjusted_section_name.strip(),
                        sections_file1[i].strip(),
                        sections_file2[i].strip(),
                        result,
                    ]
                )
                comparison_number += 1

    write_to_csv(output_csv, csv_data)


# end

# ========================================================================
# 使用例
# lines_to_include            : 比較行数(int)
# comma_sections_to_compare   :カンマ区切りの比較場所([])
# section_name                :セクションの名前(string)
# ========================================================================

keywords_with_options = {
    "test": {"lines_to_include": 2, "section_name": "テストセクション"},
    "uniqu": {
        "lines_to_include": 1,
        "comma_sections_to_compare": [0, 1, 3],
        "section_name": "ユニークセクション",
    },
    "import": {"lines_to_include": 1, "section_name": "インポート"},
}
output_csv = "comparison_results"
process_files_to_csv(keywords_with_options, output_csv)
