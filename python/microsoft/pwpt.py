import os
import csv
import win32com.client
import pythoncom
from threading import Timer
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


def set_password_pptx(app, file_path, password):
    doc = app.Presentations.Open(file_path, WithWindow=False)
    try:
        doc.WritePassword = password
        doc.Save()
    finally:
        doc.Close()


def set_password_docx(app, file_path, password):
    doc = app.Documents.Open(file_path)
    try:
        doc.WritePassword = password
        doc.SaveAs2(file_path, WritePassword=password)
    finally:
        doc.Close()


def set_password_xlsx(app, file_path, password, result_queue):
    app.Visible = False
    doc = None
    try:
        doc = app.Workbooks.Open(file_path)
        if doc.WriteReserved:
            result_queue.put(
                ("PASS", "File is password protected or another error occurred")
            )
        else:
            doc.Password = password
            doc.SaveAs(file_path, Password="", WriteResPassword=password)
            result_queue.put(("OK", "Password set successfully"))
    except Exception as e:
        result_queue.put(("ERROR", f"Error setting password: {e}"))
    finally:
        if doc is not None:
            doc.Close()
        app.Quit()


def set_password_office(file_path, password, result_queue):
    pythoncom.CoInitialize()
    app = None
    try:
        extension = os.path.splitext(file_path)[1].lower()
        app = open_office_application(extension)
        if extension == ".pptx":
            set_password_pptx(app, file_path, password)
        elif extension == ".docx":
            set_password_docx(app, file_path, password)
        elif extension == ".xlsx":
            set_password_xlsx(app, file_path, password, result_queue)
        else:
            result_queue.put(("NG", "Unsupported file type"))
    except Exception as e:
        result_queue.put(("ERROR", f"Error setting password: {e}"))
    finally:
        if app is not None:
            app.Quit()
        pythoncom.CoUninitialize()


def process_file(file_path, password):
    result_queue = Queue()
    timer = Timer(10.0, lambda: result_queue.put(("TIMEOUT", "Operation timed out")))
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
                    results.append((os.path.basename(file_path), file_path, result[1]))
                    print(f"{file_path}: {result[1]}")
                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
                    continue
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
