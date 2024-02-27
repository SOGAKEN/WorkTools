import os
import csv
import win32com.client
import pythoncom
from threading import Thread, Timer
from queue import Queue


def open_office_application(extension):
    apps = {
        ".pptx": "PowerPoint.Application",
        ".docx": "Word.Application",
        ".xlsx": "Excel.Application",
    }
    app_name = apps.get(extension)
    if app_name:
        if extension == ".docx":
            return win32com.client.Dispatch(app_name)
        else:
            return win32com.client.gencache.EnsureDispatch(app_name)
    else:
        raise ValueError("Unsupported file type")


def set_password_office(file_path, password, result_queue):
    pythoncom.CoInitialize()
    app = None
    try:
        extension = os.path.splitext(file_path)[1].lower()
        app = open_office_application(extension)
        if extension == ".pptx":
            doc = app.Presentations.Open(file_path, WithWindow=False)
            doc.WritePassword = password
            doc.Save()
            doc.Close()
        elif extension == ".docx":
            doc = app.Documents.Open(file_path)
            doc.WritePassword = password
            doc.SaveAs2(file_path, WritePassword=password)
            doc.Close()
        elif extension == ".xlsx":
            app.Visible = False
            doc = app.Workbooks.Open(file_path)
            if not doc.WriteReserved:
                doc.Password = password
                doc.SaveAs(file_path, Password=password)
            doc.Close()
        else:
            raise ValueError("Unsupported file type")
        result_queue.put(("OK", "Password set successfully"))
    except Exception as e:
        result_queue.put(("ERROR", f"Error setting password: {e}"))
    finally:
        if app is not None:
            app.Quit()
        pythoncom.CoUninitialize()


def process_file_with_timeout(file_path, password, timeout_seconds=30):
    result_queue = Queue()
    thread = Thread(
        target=set_password_office, args=(file_path, password, result_queue)
    )
    thread.start()
    thread.join(timeout=timeout_seconds)
    if thread.is_alive():
        return ("TIMEOUT", "Operation timed out")
    else:
        return result_queue.get()


def process_directory_for_documents(directory, password):
    results = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith((".pptx", ".docx", ".xlsx")):
                file_path = os.path.join(root, file)
                print(f"Processing {file_path}...")
                result = process_file_with_timeout(file_path, password)
                results.append(
                    (os.path.basename(file_path), file_path, result))
                print(f"{file_path}: {result}")
    return results


def write_results_to_csv(results, output_csv_path):
    with open(output_csv_path, "w", newline="", encoding="utf-8-sig") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["FileName", "FilePath", "Result"])
        for result in results:
            writer.writerow([result[0], result[1], result[2][1]])


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
