import boto3
import paramiko
import json
import os
import base64

def lambda_handler(event, context):
    # 環境変数からシークレット名を取得
    secret_name = os.environ['SECRET_NAME']
    log_file_path = os.environ['LOG_FILE_PATH']
    search_phrase = os.environ['SEARCH_PHRASE']
    region_name = os.environ['REGION']

    # Secrets Manager クライアントを初期化
    session = boto3.session.Session()
    client = session.client(
        service_name='secretsmanager',
        region_name=region_name
    )

    # Secrets Managerからシークレットを取得
    get_secret_value_response = client.get_secret_value(SecretId=secret_name)
    if 'SecretString' in get_secret_value_response:
        secret = json.loads(get_secret_value_response['SecretString'])
    else:
        decoded_binary_secret = base64.b64decode(get_secret_value_response['SecretBinary'])
        secret = json.loads(decoded_binary_secret.decode('utf-8'))

    # SSHクライアントの初期化
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    # 秘密鍵をデコード
    key = paramiko.RSAKey.from_private_key_file(secret['private_key_file'])

    # SSH接続の確立
    ssh.connect(hostname=secret['instance_ip'], username=secret['ssh_user'], pkey=key)

    # コマンドの設定と実行
    grep_phrase = ' -e '.join(["'{}'".format(phrase.strip()) for phrase in search_phrase.split(',')])
    command = f'grep -e {grep_phrase} {log_file_path}'
    stdin, stdout, stderr = ssh.exec_command(command)
    log_content = stdout.read().decode()

    # SSH接続の終了
    ssh.close()

    # ログの内容
    return {
        'statusCode': 200,
        'body': log_content 
    }
