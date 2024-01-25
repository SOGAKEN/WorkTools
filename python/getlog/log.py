import boto3
import json
import os

def lambda_handler(event, context):
    # 環境変数から設定値を取得
    instance_id = os.environ['INSTANCE_ID']
    log_file_path = os.environ['LOG_FILE_PATH']
    search_phrase = os.environ['SEARCH_PHRASE']

    # SSM クライアントの初期化
    ssm_client = boto3.client('ssm')

    # コマンドの設定（ログファイルの内容を取得）
    commands = [f'cat {log_file_path}']

    # SSM Run Command の実行
    response = ssm_client.send_command(
        InstanceIds=[instance_id],
        DocumentName='AWS-RunShellScript',
        Parameters={'commands': commands},
    )

    # コマンドの実行 ID を取得
    command_id = response['Command']['CommandId']

    # コマンドの出力を取得
    output = ssm_client.get_command_invocation(
        CommandId=command_id,
        InstanceId=instance_id
    )

    # ログの内容
    log_content = output['StandardOutputContent']

    # 特定の文言を検索
    if search_phrase in log_content:
        print(f"'{search_phrase}' was found in the log.")
        return {
            'statusCode': 200,
            'body': json.dumps(f"'{search_phrase}' was found in the log.")
        }
    else:
        print(f"'{search_phrase}' was not found in the log.")
        return {
            'statusCode': 200,
            'body': json.dumps(f"'{search_phrase}' was not found in the log.")
        }
