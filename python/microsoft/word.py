import win32com.client as win32
from win32com.client import constants

def set_document_readonly(filepath, edit_password):
    try:
        # Wordアプリケーションを開始
        word = win32.gencache.EnsureDispatch('Word.Application')
        word.Visible = False

        # ドキュメントを開く
        doc = word.Documents.Open(filepath)

        # 読み取り専用保護を適用（編集にはパスワードが必要）
        if doc.ProtectionType == constants.wdNoProtection:
            doc.Protect(Type=constants.wdAllowOnlyReading, NoReset=True, Password=edit_password)
            print("ドキュメントは読み取り専用に設定されました。")
        else:
            print("ドキュメントは既に保護されています。")

        # 変更を保存して閉じる
        doc.Save()
        doc.Close()

    except Exception as e:
        print(f"エラーが発生しました: {e}")

    finally:
        # Wordアプリケーションを閉じる
        word.Quit()

# 使用例
# 'path_to_your_document.docx'を文書のパスに置き換えてください。
# 'your_edit_password'を希望の編集パスワードに置き換えてください。
set_document_readonly('path_to_your_document.docx', 'your_edit_password')
