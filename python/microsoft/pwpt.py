import os
import csv
import win32com.client
import pythoncom
from threading import Timer
from queue import Queue


def set_password_office(file_path, password, result_queue):
    pythoncom.CoInitialize()
    extension = os.path.splitext(file_path)[1].lower()
    try:
        if extension == ".pptx":
            app = win32com.client.gencache.EnsureDispatch(
                "PowerPoint.Application")
            doc = app.Presentations.Open(file_path, WithWindow=False)
            doc.WritePassword = password
            doc.Save()
            doc.Close()
        elif extension == ".docx":
            app = win32com.client.gencache.EnsureDispatch("Word.Application")
            doc = app.Documents.Open(file_path)
            doc.WritePassword = password
            doc.SaveAs2(file_path, WritePassword=password)
            doc.Close()
        elif extension == ".xlsx":
            app = win32com.client.gencache.EnsureDispatch("Excel.Application")
            doc = app.Workbooks.Open(file_path)
            doc.Password = password
            doc.SaveAs(file_path, WriteResPassword=password)
            doc.Close()
        else:
            result_queue.put(("NG", "Unsupported file type"))
            return
        app.Quit()
        result_queue.put(("OK", "Password set successfully"))
    except Exception as e:
        result_queue.put(
            ("PASS", "File is password protected or another error occurred")
        )
    finally:
        pythoncom.CoUninitialize()


def process_file(file_path, password):
    result_queue = Queue()
    timer = Timer(10.0, lambda: result_queue.put(
        ("TIMEOUT", "Operation timed out")))
    try:
        timer.start()
        set_password_office(file_path, password, result_queue)
    finally:
        timer.cancel()
    return file_path, result_queue.get()


def process_directory_for_documents(directory, password):
    results = []
    file_count = 0
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith((".pptx", ".docx", ".xlsx")):
                file_path = os.path.join(root, file)
                print(f"Processing {file_path}...")
                try:
                    result = process_file(file_path, password)
                    results.append(
                        (os.path.basename(file_path), file_path, result[1]))
                    print(f"{file_path}: {result[1]}")
                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
                    continue  # Skip this file and continue with the next
                finally:
                    file_count += 1
    print(f"Total files processed: {file_count}")
    return results


def write_results_to_csv(results, output_csv_path):
    with open(output_csv_path, "w", newline="", encoding="utf-8-sig") as csvfile:
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
    print("Starting file processing...")
    results = process_directory_for_documents(current_directory, edit_password)
    if results:
        output_csv_path = get_csv_path()
        write_results_to_csv(results, output_csv_path)
        print(f"Results have been saved to {output_csv_path}.")
    else:
        print("No files were processed.")
    input("Processing complete. Press Enter to exit...")
