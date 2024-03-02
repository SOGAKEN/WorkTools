import csv
from datetime import datetime
from wcwidth import wcswidth


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
    for i in range(max_length):
        section1 = sections1[i] if i < len(sections1) else ""
        section2 = sections2[i] if i < len(sections2) else ""
        if section1 != section2:
            differences.append((i + 1, section1, section2))
    return differences


def write_to_csv(filename, data):
    """CSVファイルにデータを書き込む"""
    with open(filename, "w", newline="", encoding="utf-8") as csvfile:
        writer = csv.writer(csvfile)
        # ヘッダーに'SectionName'を追加
        writer.writerow(["No.", "SectionName", "File1", "File2", "Result"])
        for row in data:
            writer.writerow(row)


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


def process_files_to_csv(file1, file2, keywords_with_options, output_csv):
    csv_data = []
    comparison_number = 1
    # セクション名の表示幅の最大値を計算
    max_section_display_length = max(
        calculate_display_length(options.get("section_name", keyword))
        for keyword, options in keywords_with_options.items()
    )

    for keyword, options in keywords_with_options.items():
        section_name = options.pop("section_name", keyword)
        sections_file1 = extract_sections(file1, keyword, **options)
        sections_file2 = extract_sections(file2, keyword, **options)
        differences = compare_sections(sections_file1, sections_file2)

        for diff in differences:
            section_number, section_file1, section_file2 = diff
            result = "NG"
            # セクション名の表示幅を調整
            adjusted_section_name = section_name + " " * (
                max_section_display_length -
                calculate_display_length(section_name)
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
                    section_name,
                    section_file1.strip(),
                    section_file2.strip(),
                    result,
                ]
            )
            comparison_number += 1

        common_length = min(len(sections_file1), len(sections_file2))
        for i in range(common_length):
            if (i + 1, sections_file1[i], sections_file2[i]) not in differences:
                result = "OK"
                # セクション名の表示幅を調整
                adjusted_section_name = section_name + " " * (
                    max_section_display_length -
                    calculate_display_length(section_name)
                )
                output_text = f"No. {comparison_number:<3} | セクション名: {adjusted_section_name} | Result: {result}"
                print(output_text)
                # CSVデータにセクション名を含める
                csv_data.append(
                    [
                        comparison_number,
                        section_name,
                        sections_file1[i].strip(),
                        sections_file2[i].strip(),
                        result,
                    ]
                )
                comparison_number += 1

    write_to_csv(output_csv, csv_data)


# 使用例とその他の関数定義は省略

# 使用例
keywords_with_options = {
    "test": {"lines_to_include": 2, "section_name": "テストセクション"},
    "uniqu": {
        "lines_to_include": 1,
        "comma_sections_to_compare": [0, 1],
        "section_name": "ユニークセクション",
    },
    "import": {"lines_to_include": 1, "section_name": "インポート"},
}
process_files_to_csv(
    "file1.log", "file2.log", keywords_with_options, "comparison_results.csv"
)
