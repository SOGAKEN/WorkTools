import os
import csv
import win32com.client
import pythoncom
from threading import Timer


def set_password(file_path, password):
    try:
        pythoncom.CoInitialize()
        extension = os.path.splitext(file_path)[1].lower()
        if extension == ".pptx":
            app = win32com.client.Dispatch("PowerPoint.Application")
            presentation = app.Presentations.Open(file_path, WithWindow=False)
            presentation.SaveAs(file_path, Password=password)
            presentation.Close()
            app.Quit()
        elif extension == ".docx":
            app = win32com.client.Dispatch("Word.Application")
            document = app.Documents.Open(file_path)
            document.SaveAs2(file_path, WritePassword=password)
            document.Close()
            app.Quit()
        elif extension == ".xlsx":
            app = win32com.client.Dispatch("Excel.Application")
            workbook = app.Workbooks.Open(file_path)
            workbook.SaveAs(file_path, WriteResPassword=password)
            workbook.Close()
            app.Quit()
        else:
            return "NG", "Unsupported file type"
        return "OK", "Password set successfully"
    except Exception as e:
        return "NG", str(e)
    finally:
        pythoncom.CoUninitialize()


def process_file(file_path, password):
    timer = Timer(10.0, lambda: os._exit(1))
    try:
        timer.start()
        result, message = set_password(file_path, password)
    except SystemExit:
        return file_path, "TIMEOUT"
    finally:
        timer.cancel()
    return file_path, result


def process_directory_for_documents(directory, password):
    results = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith((".pptx", ".docx", ".xlsx")):
                file_path = os.path.join(root, file)
                try:
                    result = process_file(file_path, password)
                    results.append(
                        (os.path.basename(file_path), file_path, result[1]))
                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
                    continue  # Skip this file and continue with the next
    return results


def write_results_to_csv(results, output_csv_path):
    with open(output_csv_path, "w", newline="", encoding="utf-8") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["FileName", "FilePath", "Result"])
        for result in results:
            writer.writerow(result)


def get_application_path():
    return os.path.dirname(os.path.abspath(__file__))


def get_csv_path():
    return os.path.join(get_application_path(), "results.csv")


if __name__ == "__main__":
    edit_password = "your_edit_password"  # Set the appropriate password here
    current_directory = get_application_path()
    results = process_directory_for_documents(current_directory, edit_password)
    if results:
        output_csv_path = get_csv_path()
        write_results_to_csv(results, output_csv_path)
        print(f"Results have been saved to {output_csv_path}.")
    else:
        print("No files were processed.")
    input("Processing complete. Press Enter to exit...")
